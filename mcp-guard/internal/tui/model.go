package tui

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/audit"
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
	totalCalls       int
	totalAllowed     int
	totalBlocked     int
	totalHITL        int
	totalRateLimited int
	totalInjected    int
	startTime        time.Time

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
	version       string
	mode          string
	targetAddr    string
	syncStatus    string // "NOMINAL", "DEGRADED", "OFFLINE"
	listenerReady bool
}

// NewModel creates the dashboard with sensible defaults.
func NewModel(auditPath, version, mode, targetAddr string) Model {
	return Model{
		startTime:     time.Now(),
		auditPath:     auditPath,
		version:       version,
		mode:          mode,
		targetAddr:    targetAddr,
		syncStatus:    "NOMINAL",
		listenerReady: true,
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
		m.fileOffset = 0
		m.staleCount = 0
	}

	if info.Size() < m.fileOffset {
		m.fileOffset = 0
	}

	if info.Size() <= m.fileOffset {
		return
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
	if now.Sub(m.lastBucketTime) >= bucketInterval {
		elapsed := now.Sub(m.lastBucketTime)
		skip := int(elapsed / bucketInterval)
		if skip > bucketCount {
			skip = bucketCount
		}
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
		end := m.scrollPos + n
		if end > len(m.entries) {
			end = len(m.entries)
		}
		return m.entries[m.scrollPos:end]
	}
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
	if width < 2 {
		width = 2
	}
	if maxVal == 0 || val == 0 {
		// Empty state: use dots to maintain bracket spacing
		if width > 12 {
			return strings.Repeat(".", 8) + strings.Repeat(" ", width-8)
		}
		return strings.Repeat(".", width-2) + " "
	}
	filled := int(math.Round(float64(val) / float64(maxVal) * float64(width)))
	if filled > width {
		filled = width
	}
	if filled < 1 && val > 0 {
		filled = 1
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// miniBarBracketed renders a mini bar with bracket boundaries for clarity.
func miniBarBracketed(val, maxVal, width int) string {
	if width < 4 {
		width = 4
	}
	bar := miniBar(val, maxVal, width-2)
	return "[" + bar + "]"
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
// padToWidth ensures a string (stripped of ANSI) is at least width chars.
func padToWidth(s string, width int) string {
	stripped := stripANSI(s)
	if len(stripped) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(stripped))
}

// padContent pads every line in a multi-line content string to the given width.
// This ensures lipgloss panels fill their allocated column width instead of shrink-wrapping.
func padContent(content string, width int) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = padToWidth(line, width)
	}
	return strings.Join(lines, "\n")
}

// ---------------------------------------------------------------------------
// View
// ---------------------------------------------------------------------------

func (m Model) View() string {
	if !m.ready {
		return centerText("Loading MCP Guard Dashboard...", 40)
	}

	header := m.renderHeader()
	footer := m.renderFooter()

	if m.width < 80 {
		// Narrow: stack everything vertically
		body := lipgloss.JoinVertical(lipgloss.Top,
			m.renderStats(m.width-8),
			m.renderActivity(m.width-8),
			m.renderTopTools(8, m.width-8),
			m.renderTopIdentities(7, m.width-8),
		)
		feed := m.renderFeed(m.width - 6)
		return lipgloss.JoinVertical(lipgloss.Top, header, body, feed, footer)
	}

	// ── Wide: 2-column grid ──────────────────────────────────────────
	// Each panel has border (2) + padding (4) = 6 overhead.
	// left_panel and right_panel each get (width - margin - 6*2) / 2 of content
	colContentWidth := (m.width - 16) / 2
	if colContentWidth < 30 {
		colContentWidth = 30
	}

	// Left column: Traffic Summary + Activity
	leftCol := lipgloss.JoinVertical(lipgloss.Top,
		m.renderStats(colContentWidth),
		m.renderActivity(colContentWidth),
	)

	// Right column: Top Tools + Top Identities
	rightCol := lipgloss.JoinVertical(lipgloss.Top,
		m.renderTopTools(8, colContentWidth),
		m.renderTopIdentities(7, colContentWidth),
	)

	// Middle row: side by side with 2-char gap
	leftCol = lipgloss.NewStyle().MarginRight(2).Render(leftCol)
	midRow := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	// Full-width bottom feed
	feed := m.renderFeed(m.width - 8)

	return lipgloss.JoinVertical(lipgloss.Top, header, midRow, feed, footer)
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

	leftBits := fmt.Sprintf("MCP GUARD  %s  %s  %s",
		m.version,
		headerModeStyle("["+m.mode+"]"),
		styled(statusLabel, statusColor),
	)
	rightBits := fmt.Sprintf("Uptime: %s", m.uptime())
	// Header content must fill m.width - 6 to match Padding(0,3) on each side.
	// This aligns header text indent with panel content (border 1 + padding 2 = 3).
	space := m.width - 6 - len(stripANSI(leftBits)) - len(stripANSI(rightBits))
	if space < 2 {
		space = 2
	}
	filler := headerFillerStyle(strings.Repeat(" ", space))

	return headerBarStyle(leftBits + filler + rightBits)
}

func (m Model) renderStats(width int) string {
	var b strings.Builder

	maxVal := m.totalCalls
	if maxVal < 1 {
		maxVal = 1
	}

	// title row
	b.WriteString(panelTitleStyle("Traffic Summary & Metrics"))
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

	// Bar width: line = label(16) + space + value(6) + space + pct(6) + space + bar = 31 + bar
	barWidth := width - 31
	if barWidth < 10 {
		barWidth = 10
	}
	if barWidth > 30 {
		barWidth = 30
	}

	for _, row := range rows {
		pct := 0.0
		if m.totalCalls > 0 {
			pct = float64(row.value) / float64(m.totalCalls) * 100
		}
		bar := miniBarBracketed(row.value, maxVal, barWidth)
		barStyled := styled(bar, row.color)
		labelStyled := styled(fmt.Sprintf("%-16s", row.label), colorLabelStr)
		valStyled := styled(fmt.Sprintf("%6d", row.value), row.color)
		pctStyled := styled(fmt.Sprintf("%5.1f%%", pct), colorDimStr)

		b.WriteString(fmt.Sprintf("%s %s %s %s\n",
			labelStyled, valStyled, pctStyled, barStyled))
	}

	return panelStyle(padContent(b.String(), width))
}

func (m Model) renderTopTools(n, width int) string {
	tools := m.topTools(n)
	var b strings.Builder
	b.WriteString(panelTitleStyle("Top Tools & High-Hit Metrics"))
	b.WriteString("\n\n")

	if len(tools) == 0 {
		b.WriteString(m.renderEmptyPanel(width, "LISTENER ACTIVE", "waiting for incoming token traffic..."))
		return panelStyle(padContent(b.String(), width))
	}

	maxCalls := tools[0].calls
	if maxCalls < 1 {
		maxCalls = 1
	}

	// Line: name(18) + space + value(5) + space + bar = 25 + bar
	barWidth := width - 25
	if barWidth < 8 {
		barWidth = 8
	}
	if barWidth > 20 {
		barWidth = 20
	}

	for _, t := range tools {
		bar := miniBarBracketed(t.calls, maxCalls, barWidth)
		barStyled := styled(bar, colorCyanStr)
		nameStyled := styled(fmt.Sprintf("%-18s", t.name), colorBrightStr)
		countStyled := styled(fmt.Sprintf("%5d", t.calls), colorWhiteStr)
		blockInfo := ""
		if t.blocks > 0 {
			blockInfo = styled(fmt.Sprintf("  %d blk", t.blocks), colorRedStr)
		}
		b.WriteString(fmt.Sprintf("%s %s %s%s\n", nameStyled, countStyled, barStyled, blockInfo))
	}

	return panelStyle(padContent(b.String(), width))
}

func (m Model) renderTopIdentities(n, width int) string {
	idents := m.topIdentities(n)
	var b strings.Builder
	b.WriteString(panelTitleStyle("Top Active Agent Identities"))
	b.WriteString("\n\n")

	if len(idents) == 0 {
		b.WriteString(m.renderEmptyPanel(width, "PROXY STANDBY", "no active agent identities detected..."))
		return panelStyle(padContent(b.String(), width))
	}

	maxCalls := idents[0].calls
	if maxCalls < 1 {
		maxCalls = 1
	}

	// Line: name(20) + space + value(5) + space + bar = 27 + bar
	barWidth := width - 27
	if barWidth < 8 {
		barWidth = 8
	}
	if barWidth > 20 {
		barWidth = 20
	}

	for _, id := range idents {
		bar := miniBarBracketed(id.calls, maxCalls, barWidth)
		barStyled := styled(bar, colorCyanStr)
		nameStyled := styled(fmt.Sprintf("%-20s", id.name), colorBrightStr)
		countStyled := styled(fmt.Sprintf("%5d", id.calls), colorWhiteStr)
		blockInfo := ""
		if id.blocks > 0 {
			blockInfo = styled(fmt.Sprintf("  %d blk", id.blocks), colorRedStr)
		}
		b.WriteString(fmt.Sprintf("%s %s %s%s\n", nameStyled, countStyled, barStyled, blockInfo))
	}

	return panelStyle(padContent(b.String(), width))
}

// renderEmptyPanel creates a structural placeholder visual for empty data states.
// The inner box is sized to exactly fit within `width` chars:
//
//   ┌── STATE LABEL ────────────┐
//   │  ·  ·  ·  ·  ·  ·  ·    │
//   │                           │
//   │     help text centered   │
//   └───────────────────────────┘
func (m Model) renderEmptyPanel(width int, stateLabel, helpText string) string {
	if width < 24 {
		width = 24
	}

	// Line layout: "  │..." or "  ┌──...──┐" or "  └──...──┘"
	// Left margin: "  " (2) + frame char (1) + spacer (0/1) = 2-4
	// Right: frame char (1) + implicit padContent padding
	//
	// Top:    "  ┌── LABEL ──...──┐"  → overhead: 2+3+1+3 = 9  (spaces + ┌── + space + ──┐)
	// Middle: "  │ ... │"           → overhead: 2+1+1+1 = 5   (spaces + │ + space + │)
	// Bottom: "  └──...──┘"         → overhead: 2+1+2   = 5   (spaces + └ + ──┘)
	innerW := width - 6 // content between │ and │ in middle rows

	var sb strings.Builder
	sb.WriteString("\n")

	// Top:   ┌── LABEL ──────────────────────┐
	dashTop := width - 9 - len(stateLabel)
	if dashTop < 2 {
		dashTop = 2
	}
	topBorder := fmt.Sprintf("  ┌── %s %s──┐", stateLabel, strings.Repeat("─", dashTop))
	sb.WriteString(dimStyle(padToWidth(topBorder, width)))
	sb.WriteString("\n")

	// Dots row:  │ ·  ·  ·  ·  ·  ·  ·  ·  │
	dotsCount := innerW / 3
	if dotsCount > 30 {
		dotsCount = 30
	}
	dotsRow := fmt.Sprintf("  │ %s│", strings.Repeat("·  ", dotsCount))
	sb.WriteString(dimStyle(padToWidth(dotsRow, width)))
	sb.WriteString("\n")

	// Blank line
	blankLine := fmt.Sprintf("  │ %s │", strings.Repeat(" ", innerW))
	sb.WriteString(dimStyle(padToWidth(blankLine, width)))
	sb.WriteString("\n")

	// Help text centered
	helpLine := fmt.Sprintf("  │ %s │", centerText(helpText, innerW))
	sb.WriteString(dimStyle(padToWidth(helpLine, width)))
	sb.WriteString("\n")

	// Blank line
	blankLine2 := fmt.Sprintf("  │ %s │", strings.Repeat(" ", innerW))
	sb.WriteString(dimStyle(padToWidth(blankLine2, width)))
	sb.WriteString("\n")

	// Bottom border:   └──────────────────────────┘
	dashBot := width - 5
	if dashBot < 2 {
		dashBot = 2
	}
	botBorder := fmt.Sprintf("  └%s┘", strings.Repeat("─", dashBot))
	sb.WriteString(dimStyle(padToWidth(botBorder, width)))

	return sb.String()
}

func (m Model) renderActivity(width int) string {
	var b strings.Builder
	b.WriteString(panelTitleStyle("Activity Timeline (5s buckets · 2min window)"))
	b.WriteString("\n\n")

	spark := sparkline(m.activityBuckets)

	// Calculate padding for sparkline if it's narrower than available width
	maxSparkWidth := width - 4
	if len(spark) < maxSparkWidth {
		pad := maxSparkWidth - len(spark)
		spark += strings.Repeat(" ", pad)
	}
	b.WriteString(fmt.Sprintf("  %s\n", spark))

	// Label the time axis
	now := time.Now()
	endLabel := now.Format("15:04:05")
	startLabel := now.Add(-time.Duration(bucketCount) * bucketInterval).Format("15:04:05")
	// Spread the labels across the full width
	timeAxis := fmt.Sprintf("  %s%*s%s",
		startLabel,
		maxSparkWidth-len(startLabel)-len(endLabel),
		"",
		endLabel)
	b.WriteString(dimStyle(timeAxis))
	b.WriteString("\n")

	return panelStyle(padContent(b.String(), width))
}

func (m Model) renderFeed(width int) string {
	entries := m.visibleEntries()
	var b strings.Builder

	// Build a combined header with entry count
	entryInfo := fmt.Sprintf("  %d entries  |  %d calls tracked",
		len(m.entries), m.totalCalls)
	b.WriteString(feedHeaderStyle("Live Feed"))
	b.WriteString(dimStyle(entryInfo))
	b.WriteString("\n")

	if len(entries) == 0 {
		fillWidth := width - 4
		if fillWidth < 20 {
			fillWidth = 20
		}
		dashLine := fmt.Sprintf("  %s", strings.Repeat("─", fillWidth))
		b.WriteString("\n")
		b.WriteString(dimStyle(dashLine))
		b.WriteString("\n")
		b.WriteString(dimStyle(fmt.Sprintf("  %s", centerText("[PROXY IDLE]  waiting for JSON-RPC traffic...", fillWidth))))
		b.WriteString("\n")
		b.WriteString(dimStyle(fmt.Sprintf("  %s", centerText("start mcp-guard serve to begin proxying", fillWidth))))
		b.WriteString("\n")
		b.WriteString(dimStyle(dashLine))
		return feedPanelStyle(padContent(b.String(), width))
	}

	// Show the "stopped scrolling" indicator
	showScrollHint := m.scrollPos > 0

	for _, e := range entries {
		ts := styled(e.Timestamp.Format("15:04:05"), colorDimStr)
		decision := colorDecisionStr(e.Decision)
		tool := styled(fmt.Sprintf("%-20s", e.Tool), colorBrightStr)
		identity := styled(fmt.Sprintf("%-16s", e.Identity), colorDimStr)

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

	return feedPanelStyle(padContent(b.String(), width))
}

func (m Model) renderFooter() string {
	// Left: controls
	controls := styled("q:quit", colorWhiteStr) + dimStyle(" · ") +
		styled("p:pause", colorWhiteStr) + dimStyle(" · ") +
		styled("↑↓:scroll", colorDimStr)

	// Status dot and entries
	statusDot := "●"
	statusColor := colorGreenStr
	if !m.fileExists {
		statusDot = "○"
		statusColor = colorRedStr
	} else if m.paused {
		statusDot = "⏸"
		statusColor = colorYellowStr
	}
	statusEntry := styled(statusDot, statusColor) + " " + dimStyle(fmt.Sprintf("%d entries", len(m.entries)))

	// Determine target display
	targetDisplay := m.mode
	if m.targetAddr != "" {
		targetDisplay = m.targetAddr
	}
	targetInfo := dimStyle(fmt.Sprintf("Target: %s", targetDisplay))

	// Sync status
	var syncColor string
	switch m.syncStatus {
	case "NOMINAL":
		syncColor = colorGreenStr
	case "DEGRADED":
		syncColor = colorYellowStr
	default:
		syncColor = colorRedStr
	}
	syncInfo := styled("Sync: "+m.syncStatus, syncColor)

	// Right side meta-data
	rightInfo := fmt.Sprintf("%s │ %s", targetInfo, syncInfo)

	// Build footer line spread across full width
	leftPart := fmt.Sprintf("%s  %s", controls, statusEntry)
	footerText := fmt.Sprintf("%s  │  %s", leftPart, rightInfo)

	return footerStyle(footerText)
}

// ---------------------------------------------------------------------------
// String constants (used in rendering to avoid lipgloss allocations in hot path)
// ---------------------------------------------------------------------------

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
