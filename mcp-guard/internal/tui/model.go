package tui

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matrix/mcp-guard/internal/audit"
)

// bucketInterval is the time window for each activity bar (5 seconds).
const bucketInterval = 5 * time.Second

// bucketCount is how many bars in the activity timeline (24 = 2 minutes).
const bucketCount = 24

// statEntry tracks per-tool or per-identity statistics.
type statEntry struct {
	name   string
	calls  int
	blocks int
}

// Model is the standalone Bubble Tea dashboard.
type Model struct {
	ready     bool
	paused    bool
	width     int
	height    int
	scrollPos int // scroll offset for live feed

	// Runtime stats
	totalCalls  int
	totalAllowed int
	totalBlocked int
	totalHITL   int
	totalRateLimited int
	totalInjected    int
	startTime   time.Time

	// Per-tool stats (top N)
	toolStats map[string]*statEntry

	// Per-identity stats (top N)
	identityStats map[string]*statEntry

	// Activity timeline (rolling 2min window)
	activityBuckets [bucketCount]int
	lastBucketTime  time.Time

	// Audit log tailing
	auditPath    string
	fileOffset   int64
	fileExists   bool
	lastFileSize int64
	staleCount   int // seconds without update
	entries      []audit.AuditEntry

	// Meta
	version string
	mode    string
}

// NewModel creates the dashboard with sensible defaults.
func NewModel(auditPath, version, mode string) Model {
	return Model{
		startTime:     time.Now(),
		auditPath:     auditPath,
		version:       version,
		mode:          mode,
		entries:       make([]audit.AuditEntry, 0, 500),
		toolStats:     make(map[string]*statEntry),
		identityStats: make(map[string]*statEntry),
		lastBucketTime: time.Now(),
	}
}

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------

type tickMsg struct{}
type pollMsg struct{ path string }
type refreshMsg struct{}

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Second, func(time.Time) tea.Msg { return tickMsg{} }),
		tea.Tick(800*time.Millisecond, func(time.Time) tea.Msg {
			return pollMsg{path: m.auditPath}
		}),
	)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ready = true
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "p":
			m.paused = !m.paused
		case "down":
			if m.scrollPos < len(m.entries)-1 {
				m.scrollPos++
			}
		case "up":
			if m.scrollPos > 0 {
				m.scrollPos--
			}
		}

	case tickMsg:
		// Update stale detection (1 tick = 1 second)
		if m.fileExists {
			m.staleCount++
		}

		// Prune activity buckets if they expired
		m.pruneBuckets()

		return m, tea.Tick(time.Second, func(time.Time) tea.Msg { return tickMsg{} })

	case pollMsg:
		if !m.paused {
			m.processNewEntries(msg.path)
		}
		return m, tea.Tick(800*time.Millisecond, func(time.Time) tea.Msg {
			return pollMsg{path: m.auditPath}
		})
	}
	return m, nil
}

// ---------------------------------------------------------------------------
// Audit log processing
// ---------------------------------------------------------------------------

func (m *Model) processNewEntries(path string) {
	f, err := os.Open(path)
	if err != nil {
		m.fileExists = false
		return
	}
	defer f.Close()

	info, _ := f.Stat()
	if !m.fileExists {
		m.fileExists = true
		m.fileOffset = 0 // reset on new file (truncation or first open)
		m.staleCount = 0
	}

	if info.Size() < m.fileOffset {
		// File was truncated — reset
		m.fileOffset = 0
	}

	if info.Size() <= m.fileOffset {
		return // nothing new
	}

	m.staleCount = 0
	m.lastFileSize = info.Size()

	data := make([]byte, info.Size()-m.fileOffset)
	n, _ := f.ReadAt(data, m.fileOffset)
	m.fileOffset = info.Size()

	if n == 0 {
		return
	}

	start := 0
	now := time.Now()
	for i := 0; i < n; i++ {
		if data[i] == '\n' && i > start {
			entry := audit.ParseEntry(data[start:i])
			m.ingestEntry(entry, now)
			start = i + 1
		}
	}
}

func (m *Model) ingestEntry(e audit.AuditEntry, now time.Time) {
	m.entries = append(m.entries, e)
	if len(m.entries) > 500 {
		m.entries = m.entries[len(m.entries)-500:]
	}

	m.totalCalls++

	// Update bucket
	m.addToCurrentBucket(now)

	// Update counters
	switch e.Decision {
	case "allow":
		m.totalAllowed++
	case "block":
		m.totalBlocked++
		// Check reason for sub-type
		if strings.Contains(e.Reason, "injection") {
			m.totalInjected++
		} else if strings.Contains(e.Reason, "rate limit") {
			m.totalRateLimited++
		}
	case "pending", "hitl":
		m.totalHITL++
	}

	// Update tool stats
	tool := e.Tool
	if tool == "" {
		tool = "(unknown)"
	}
	ts, ok := m.toolStats[tool]
	if !ok {
		ts = &statEntry{name: tool}
		m.toolStats[tool] = ts
	}
	ts.calls++
	if e.Decision == "block" {
		ts.blocks++
	}

	// Update identity stats
	identity := e.Identity
	if identity == "" {
		identity = "(anonymous)"
	}
	is, ok := m.identityStats[identity]
	if !ok {
		is = &statEntry{name: identity}
		m.identityStats[identity] = is
	}
	is.calls++
	if e.Decision == "block" {
		is.blocks++
	}
}

// ---------------------------------------------------------------------------
// Activity timeline (rolling bucket)
// ---------------------------------------------------------------------------

func (m *Model) addToCurrentBucket(now time.Time) {
	// If more than bucketInterval has passed since last, advance
	if now.Sub(m.lastBucketTime) >= bucketInterval {
		// How many buckets to advance?
		elapsed := now.Sub(m.lastBucketTime)
		skip := int(elapsed / bucketInterval)
		if skip > bucketCount {
			skip = bucketCount // catch up by clearing everything
		}
		// Shift and clear
		for i := 0; i < bucketCount-skip; i++ {
			m.activityBuckets[i] = m.activityBuckets[i+skip]
		}
		for i := bucketCount - skip; i < bucketCount; i++ {
			m.activityBuckets[i] = 0
		}
		m.lastBucketTime = now
	}
	m.activityBuckets[bucketCount-1]++
}

func (m *Model) pruneBuckets() {
	now := time.Now()
	elapsed := now.Sub(m.lastBucketTime)
	if elapsed >= bucketInterval*time.Duration(bucketCount) {
		// All buckets expired — zero them
		for i := range m.activityBuckets {
			m.activityBuckets[i] = 0
		}
		m.lastBucketTime = now
		return
	}
	if elapsed >= bucketInterval {
		skip := int(elapsed / bucketInterval)
		for i := 0; i < bucketCount-skip; i++ {
			m.activityBuckets[i] = m.activityBuckets[i+skip]
		}
		for i := bucketCount - skip; i < bucketCount; i++ {
			m.activityBuckets[i] = 0
		}
		m.lastBucketTime = now
	}
}

// ---------------------------------------------------------------------------
// Helpers for rendering
// ---------------------------------------------------------------------------

func (m Model) uptime() string {
	d := time.Since(m.startTime)
	h := int(d.Hours())
	mins := int(d.Minutes()) % 60
	secs := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, mins, secs)
	}
	return fmt.Sprintf("%dm %ds", mins, secs)
}

func (m Model) topTools(n int) []statEntry {
	var list []statEntry
	for _, se := range m.toolStats {
		list = append(list, *se)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].calls > list[j].calls
	})
	if len(list) > n {
		list = list[:n]
	}
	return list
}

func (m Model) topIdentities(n int) []statEntry {
	var list []statEntry
	for _, se := range m.identityStats {
		list = append(list, *se)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].calls > list[j].calls
	})
	if len(list) > n {
		list = list[:n]
	}
	return list
}

func (m Model) visibleEntries() []audit.AuditEntry {
	if len(m.entries) == 0 {
		return nil
	}
	// Show last N where N depends on available height
	n := 12
	if m.height > 40 {
		n = 20
	}
	if m.scrollPos > 0 {
		// Manual scroll mode — show from scroll position
		end := m.scrollPos + n
		if end > len(m.entries) {
			end = len(m.entries)
		}
		return m.entries[m.scrollPos:end]
	}
	// Auto-scroll to latest
	if len(m.entries) <= n {
		return m.entries
	}
	return m.entries[len(m.entries)-n:]
}

// ---------------------------------------------------------------------------
// Sparkline / bar rendering
// ---------------------------------------------------------------------------

var barChars = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

func sparkline(buckets [bucketCount]int) string {
	maxVal := 0
	for _, v := range buckets {
		if v > maxVal {
			maxVal = v
		}
	}
	if maxVal == 0 {
		return strings.Repeat("▁", bucketCount)
	}
	var sb strings.Builder
	for _, v := range buckets {
		idx := int(math.Round(float64(v) / float64(maxVal) * float64(len(barChars)-1)))
		if idx >= len(barChars) {
			idx = len(barChars) - 1
		}
		sb.WriteString(barChars[idx])
	}
	return sb.String()
}

func miniBar(val, maxVal, width int) string {
	if maxVal == 0 {
		return strings.Repeat(" ", width)
	}
	filled := int(math.Round(float64(val) / float64(maxVal) * float64(width)))
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// Status returns connection status string.
func (m Model) status() string {
	if !m.fileExists {
		return "disconnected"
	}
	if m.staleCount > 15 {
		return "stale"
	}
	return "live"
}

// ---------------------------------------------------------------------------
// View
// ---------------------------------------------------------------------------

func (m Model) View() string {
	if !m.ready {
		return centerText("Loading MCP Guard Dashboard...", 40)
	}

	// Build all the sections
	header := m.renderHeader()
	statsPanel := m.renderStats()
	topToolsPanel := m.renderTopTools(8)
	topIdentPanel := m.renderTopIdentities(7)
	activityPanel := m.renderActivity()
	feedPanel := m.renderFeed()
	footer := m.renderFooter()

	// Layout depends on width
	if m.width < 80 {
		// Narrow: stack everything vertically
		body := lipglossJoinVertical(
			statsPanel,
			activityPanel,
			topToolsPanel,
			topIdentPanel,
		)
		return lipglossJoinVertical(header, body, feedPanel, footer)
	}

	// Wide: side-by-side layout
	// Left column: stats + activity
	leftCol := lipglossJoinVertical(statsPanel, activityPanel)

	// Right column: top tools + top identities
	rightCol := lipglossJoinVertical(topToolsPanel, topIdentPanel)

	// Middle section: side by side
	midRow := lipglossJoinHorizontal(leftCol, "  ", rightCol)

	return lipglossJoinVertical(header, midRow, feedPanel, footer)
}

// ---------------------------------------------------------------------------
// Render sections
// ---------------------------------------------------------------------------

func (m Model) renderHeader() string {
	status := m.status()
	var statusColor, statusLabel string
	switch status {
	case "live":
		statusColor = colorGreenStr
		statusLabel = "● LIVE"
	case "stale":
		statusColor = colorYellowStr
		statusLabel = "● STALE"
	default:
		statusColor = colorRedStr
		statusLabel = "○ DISCONNECTED"
	}

	// Build header with padding on right
	leftBits := fmt.Sprintf("MCP GUARD  %s  %s  %s",
		m.version,
		headerModeStyle("["+m.mode+"]"),
		styled(statusLabel, statusColor),
	)
	rightBits := fmt.Sprintf("Uptime: %s", m.uptime())

	// Fill middle with dots if wide enough
	filler := ""
	if m.width > 70 {
		space := m.width - len(stripANSI(leftBits)) - len(stripANSI(rightBits)) - 4
		if space > 0 {
			filler = headerFillerStyle(strings.Repeat(" ", space))
		}
	}

	return headerBarStyle(leftBits + filler + rightBits)
}

func (m Model) renderStats() string {
	var b strings.Builder

	maxVal := m.totalCalls
	if maxVal < 1 {
		maxVal = 1
	}

	b.WriteString(panelTitleStyle("Traffic Summary"))
	b.WriteString("\n\n")

	rows := []struct {
		label string
		value int
		color string
		icon  string
	}{
		{"Total Calls", m.totalCalls, colorWhiteStr, ""},
		{"Allowed", m.totalAllowed, colorGreenStr, "✓"},
		{"Blocked", m.totalBlocked, colorRedStr, "●"},
		{"HITL Pending", m.totalHITL, colorYellowStr, "⏳"},
		{"Rate Limited", m.totalRateLimited, colorYellowStr, "⚡"},
		{"Injection Blocks", m.totalInjected, colorRedStr, "🛡"},
	}

	for _, row := range rows {
		pct := 0.0
		if m.totalCalls > 0 {
			pct = float64(row.value) / float64(m.totalCalls) * 100
		}
		bar := miniBar(row.value, maxVal, 20)
		barStyled := styled(bar, row.color)
		labelStyled := styled(fmt.Sprintf("%-16s", row.label), colorLabelStr)
		valStyled := styled(fmt.Sprintf("%6d", row.value), row.color)
		pctStyled := styled(fmt.Sprintf("%5.1f%%", pct), colorDimStr)

		b.WriteString(fmt.Sprintf("%s %s %s %s\n",
			labelStyled, valStyled, pctStyled, barStyled))
	}

	return panelStyle(b.String())
}

func (m Model) renderTopTools(n int) string {
	tools := m.topTools(n)
	var b strings.Builder
	b.WriteString(panelTitleStyle("Top Tools"))
	b.WriteString("\n\n")

	if len(tools) == 0 {
		b.WriteString(dimStyle("  waiting for data..."))
		return panelStyle(b.String())
	}

	maxCalls := tools[0].calls
	if maxCalls < 1 {
		maxCalls = 1
	}

	for _, t := range tools {
		bar := miniBar(t.calls, maxCalls, 12)
		barStyled := styled(bar, colorCyanStr)
		nameStyled := styled(fmt.Sprintf("%-16s", t.name), colorBrightStr)
		countStyled := styled(fmt.Sprintf("%5d", t.calls), colorWhiteStr)
		blockInfo := ""
		if t.blocks > 0 {
			blockInfo = styled(fmt.Sprintf("  %d blocks", t.blocks), colorRedStr)
		}
		b.WriteString(fmt.Sprintf("%s %s %s%s\n", nameStyled, countStyled, barStyled, blockInfo))
	}
	return panelStyle(b.String())
}

func (m Model) renderTopIdentities(n int) string {
	idents := m.topIdentities(n)
	var b strings.Builder
	b.WriteString(panelTitleStyle("Top Identities"))
	b.WriteString("\n\n")

	if len(idents) == 0 {
		b.WriteString(dimStyle("  waiting for data..."))
		return panelStyle(b.String())
	}

	maxCalls := idents[0].calls
	if maxCalls < 1 {
		maxCalls = 1
	}

	for _, id := range idents {
		bar := miniBar(id.calls, maxCalls, 12)
		barStyled := styled(bar, colorCyanStr)
		nameStyled := styled(fmt.Sprintf("%-18s", id.name), colorBrightStr)
		countStyled := styled(fmt.Sprintf("%5d", id.calls), colorWhiteStr)
		blockInfo := ""
		if id.blocks > 0 {
			blockInfo = styled(fmt.Sprintf("  %d blocks", id.blocks), colorRedStr)
		}
		b.WriteString(fmt.Sprintf("%s %s %s%s\n", nameStyled, countStyled, barStyled, blockInfo))
	}
	return panelStyle(b.String())
}

func (m Model) renderActivity() string {
	var b strings.Builder
	b.WriteString(panelTitleStyle("Activity (5s buckets · 2min window)"))
	b.WriteString("\n\n")

	spark := sparkline(m.activityBuckets)
	b.WriteString(fmt.Sprintf("  %s\n", spark))

	// Label the time axis
	now := time.Now()
	endLabel := now.Format("15:04:05")
	startLabel := now.Add(-time.Duration(bucketCount) * bucketInterval).Format("15:04:05")
	timeAxis := fmt.Sprintf("  %s%*s%s", startLabel, bucketCount-2*len(startLabel), "", endLabel)
	b.WriteString(dimStyle(timeAxis))
	b.WriteString("\n")

	return panelStyle(b.String())
}

func (m Model) renderFeed() string {
	entries := m.visibleEntries()
	var b strings.Builder
	b.WriteString(feedHeaderStyle("Live Feed"))
	b.WriteString("\n")

	if len(entries) == 0 {
		b.WriteString("\n" + dimStyle("  waiting for traffic..."))
		return feedPanelStyle(b.String())
	}

	// Show the "stopped scrolling" indicator
	showScrollHint := m.scrollPos > 0

	for _, e := range entries {
		ts := styled(e.Timestamp.Format("15:04:05"), colorDimStr)
		decision := colorDecisionStr(e.Decision)
		tool := styled(fmt.Sprintf("%-18s", e.Tool), colorBrightStr)
		identity := styled(fmt.Sprintf("%-14s", e.Identity), colorDimStr)

		detail := ""
		switch e.Decision {
		case "block":
			detail = dimStyle(e.Reason)
		case "pending", "hitl":
			detail = dimStyle("awaiting approval")
		}

		b.WriteString(fmt.Sprintf("  %s %s  %s %s %s\n", ts, decision, identity, tool, detail))
	}

	if showScrollHint {
		b.WriteString("\n" + dimStyle("  ↑↓ scrolling — tail follows at bottom"))
	}

	return feedPanelStyle(b.String())
}

func (m Model) renderFooter() string {
	controls := "q:quit  p:pause" + styled("  ↑↓:scroll", colorDimStr)
	statusDot := "●"
	statusColor := colorGreenStr
	if !m.fileExists {
		statusDot = "○"
		statusColor = colorRedStr
	} else if m.paused {
		statusDot = "⏸"
		statusColor = colorYellowStr
	}

	status := styled(statusDot, statusColor)
	entryCount := fmt.Sprintf("%d entries", len(m.entries))
	return footerStyle(fmt.Sprintf("%s  %s    %s", controls, entryCount, status))
}

// ---------------------------------------------------------------------------
// String constants (used in rendering to avoid lipgloss allocations in hot path)
// ---------------------------------------------------------------------------

// Temporary string-based styles for section rendering.
// We do NOT render lipgloss styles inside loops — we use simple ANSI strings.
var (
	colorGreenStr  = "\033[32m"
	colorRedStr    = "\033[31m"
	colorYellowStr = "\033[33m"
	colorBlueStr   = "\033[34m"
	colorCyanStr   = "\033[36m"
	colorDimStr    = "\033[2m"
	colorWhiteStr  = "\033[37m"
	colorBrightStr = "\033[1;37m"
	colorLabelStr  = "\033[2;37m"
	resetStr       = "\033[0m"
)

func styled(text, color string) string {
	return color + text + resetStr
}

// stripANSI removes ANSI escape sequences for width calculation.
func stripANSI(s string) string {
	var out strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		out.WriteRune(r)
	}
	return out.String()
}

// ---------------------------------------------------------------------------
// Temporary layout helpers (until we move to lipgloss join)
// ---------------------------------------------------------------------------

var joinVertical = func(parts ...string) string {
	return strings.Join(parts, "\n")
}

var joinHorizontal = func(parts ...string) string {
	return strings.Join(parts, "")
}

// Placeholder lipgloss wrappers — updated in styles.go
var (
	panelStyle      = func(s string) string { return "╭────\n" + s + "\n╰────" }
	panelTitleStyle = func(s string) string { return "\033[1;36m" + s + "\033[0m" }
	headerBarStyle  = func(s string) string { return "\033[44;37m " + s + " \033[0m" }
	headerModeStyle = func(s string) string { return "\033[2;37m" + s + "\033[0m" }
	headerFillerStyle = func(s string) string { return "\033[44m" + s + "\033[0m" }
	feedHeaderStyle  = func(s string) string { return "\033[1;36m" + s + "\033[0m" }
	feedPanelStyle   = func(s string) string { return s }
	dimStyle         = func(s string) string { return "\033[2m" + s + "\033[0m" }
	footerStyle      = func(s string) string { return "\033[2;37m" + s + "\033[0m" }
)

func lipglossJoinVertical(parts ...string) string {
	return strings.Join(parts, "\n")
}

func lipglossJoinHorizontal(parts ...string) string {
	var sb strings.Builder
	first := true
	for _, p := range parts {
		if !first {
			sb.WriteString("  ")
		}
		sb.WriteString(p)
		first = false
	}
	return sb.String()
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
		}
	pad := (width - len(text)) / 2
	if pad < 0 {
		pad = 0
		}
	return strings.Repeat(" ", pad) + text
}

func colorDecisionStr(d string) string {
	switch d {
	case "allow":
		return styled("ALLOW", colorGreenStr)
	case "block":
		return styled("BLOCK", colorRedStr)
	case "pending", "hitl":
		return styled("HITL ", colorYellowStr)
	default:
		return styled(fmt.Sprintf("%-5s", d), colorDimStr)
	}
}
