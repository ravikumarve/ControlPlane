package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Lipgloss-based styles for the dashboard layout.
// These wrap the section placeholders defined in model.go.
var (
	// Color palette
	colorHeaderBg = lipgloss.Color("#1E3A5F")
	colorHeaderFg = lipgloss.Color("#FFFFFF")
	colorBorder   = lipgloss.AdaptiveColor{Light: "#AAAAAA", Dark: "#555555"}
	colorPanelBg  = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#1A1A2E"}
	colorFeedBg   = lipgloss.AdaptiveColor{Light: "#F5F5F5", Dark: "#0D0D1A"}
	colorAccent   = lipgloss.AdaptiveColor{Light: "#3366CC", Dark: "#5588EE"}
	colorTitleFg  = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}
	colorDim      = lipgloss.AdaptiveColor{Light: "#888888", Dark: "#777777"}
)

func init() {
	// Override placeholder functions with real lipgloss styles
	panelStyle = func(s string) string {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2).
			MaxWidth(56).
			Render(s)
	}

	panelTitleStyle = func(s string) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			Render(s)
	}

	headerBarStyle = func(s string) string {
		return lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorHeaderFg).
			Bold(true).
			Padding(0, 1).
			Width(120).
			Render(s)
	}

	headerModeStyle = func(s string) string {
		return lipgloss.NewStyle().
			Foreground(colorDim).
			Render(s)
	}

	headerFillerStyle = func(s string) string {
		return lipgloss.NewStyle().
			Background(colorHeaderBg).
			Render(s)
	}

	feedHeaderStyle = func(s string) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			Render(s)
	}

	feedPanelStyle = func(s string) string {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1).
			Render(s)
	}

	dimStyle = func(s string) string {
		return lipgloss.NewStyle().
			Foreground(colorDim).
			Render(s)
	}

	footerStyle = func(s string) string {
		return lipgloss.NewStyle().
			Foreground(colorDim).
			Padding(0, 1).
			Render(s)
	}
}

// ---------------------------------------------------------------------------
// Legacy ANSI helpers (kept for rendering hot path in stats/tables)
// ---------------------------------------------------------------------------

// Ensure all the string-based style vars are accessible.
// These are defined in model.go and used in render methods.
