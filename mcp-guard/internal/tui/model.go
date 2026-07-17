package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matrix/mcp-guard/internal/audit"
)

// Model is the Bubble Tea dashboard model.
type Model struct {
	ready      bool
	paused     bool
	width      int
	height     int

	// Stats
	totalCalls  int
	allowed     int
	blocked     int
	hitlPending int
	startTime   time.Time

	// Audit log tailing
	auditPath  string
	fileOffset int64
	entries    []audit.AuditEntry

	// Meta
	version string
	mode    string
}

// NewModel creates the dashboard model.
func NewModel(auditPath, version, mode string) Model {
	return Model{
		startTime: time.Now(),
		auditPath: auditPath,
		version:   version,
		mode:      mode,
		entries:   make([]audit.AuditEntry, 0, 200),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), pollCmd(m.auditPath))
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg{} })
}

func pollCmd(path string) tea.Cmd {
	return tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
		return pollMsg{path: path}
	})
}

type tickMsg struct{}
type pollMsg struct{ path string }

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
		}

	case tickMsg:
		return m, tickCmd()

	case pollMsg:
		if !m.paused {
			newEntries := readEntries(msg.path, &m.fileOffset)
			for _, e := range newEntries {
				m.entries = append(m.entries, e)
				m.totalCalls++
				switch e.Decision {
				case "allow":
					m.allowed++
				case "block":
					m.blocked++
				case "pending", "hitl":
					m.hitlPending++
				}
			}
			if len(m.entries) > 200 {
				m.entries = m.entries[len(m.entries)-200:]
			}
		}
		return m, pollCmd(m.auditPath)
	}
	return m, nil
}

func readEntries(path string, offset *int64) []audit.AuditEntry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	info, _ := f.Stat()
	if info.Size() <= *offset {
		return nil
	}

	f.Seek(*offset, 0)
	data := make([]byte, info.Size()-*offset)
	n, _ := f.Read(data)
	*offset = info.Size()

	var entries []audit.AuditEntry
	start := 0
	for i := 0; i < n; i++ {
		if data[i] == '\n' && i > start {
			entries = append(entries, audit.ParseEntry(data[start:i]))
			start = i + 1
		}
	}
	return entries
}

func (m Model) View() string {
	if !m.ready {
		return "Loading MCP Guard dashboard..."
	}

	uptime := time.Since(m.startTime).Round(time.Second).String()

	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render(" MCP GUARD DASHBOARD "))
	b.WriteString("  ") // pad
	if m.version != "" {
		b.WriteString(m.version)
	}
	b.WriteString(fmt.Sprintf("  Mode: %s  ", m.mode))
	if m.paused {
		b.WriteString(" ⏸ PAUSED")
	}
	b.WriteString("\n\n")

	// Stats panel
	b.WriteString(panelStyle.Render(
		fmt.Sprintf("%s %d\n%s %d\n%s %d\n%s %d\n%s %d\n\nUptime: %s",
			colorStat("Total Calls:"), m.totalCalls,
			colorStat("Allowed:"), m.allowed,
			colorStat("Blocked:"), m.blocked,
			colorStat("HITL Pending:"), m.hitlPending,
			colorStat("Rate Limited:"), 0,
			uptime,
		),
	))
	b.WriteString("\n")

	// Log feed panel - last 15 entries
	var feedLines []string
	start := 0
	if len(m.entries) > 15 {
		start = len(m.entries) - 15
	}
	for _, e := range m.entries[start:] {
		ts := fmtTimeStyle(e.Timestamp.Format("15:04:05"))
		decision := colorDecision(e.Decision)
		tool := fmtToolStyle(e.Tool)
		feedLines = append(feedLines, fmt.Sprintf("%s  %s  %s", ts, decision, tool))
	}
	if len(feedLines) > 0 {
		b.WriteString(fullPanelStyle.Render(
			headerStyle.Render("Live Feed") + "\n" +
				strings.Join(feedLines, "\n"),
		))
	} else {
		b.WriteString(fullPanelStyle.Render(
			headerStyle.Render("Live Feed") + "\nWaiting for traffic...",
		))
	}
	b.WriteString("\n\n")
	b.WriteString(footerStyle.Render(" q:quit  p:pause  ↑↓:scroll"))
	b.WriteString("\n")

	return b.String()
}
