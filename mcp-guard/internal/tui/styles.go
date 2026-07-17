package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color definitions
var (
	colorGreen  = lipgloss.AdaptiveColor{Light: "#00AA00", Dark: "#22DD22"}
	colorRed    = lipgloss.AdaptiveColor{Light: "#CC0000", Dark: "#FF4444"}
	colorYellow = lipgloss.AdaptiveColor{Light: "#CCAA00", Dark: "#FFDD00"}
	colorBlue   = lipgloss.AdaptiveColor{Light: "#0044CC", Dark: "#4488FF"}
	colorCyan   = lipgloss.AdaptiveColor{Light: "#007777", Dark: "#44DDDD"}
	colorGray   = lipgloss.AdaptiveColor{Light: "#777777", Dark: "#999999"}
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#3366CC")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorGray).
			Padding(1, 2)

	fullPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorGray).
			Padding(1, 2)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCyan)

	footerStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(0, 0)
)

// colorStat returns a colored stat label.
func colorStat(label string) string {
	return lipgloss.NewStyle().
		Width(16).
		Foreground(colorGray).
		Render(label)
}

// colorDecision returns a colored decision string.
func colorDecision(d string) string {
	var c lipgloss.AdaptiveColor
	switch d {
	case "allow":
		c = colorGreen
	case "block":
		c = colorRed
	case "pending", "hitl":
		c = colorYellow
	default:
		c = colorGray
	}
	return lipgloss.NewStyle().
		Foreground(c).
		Bold(true).
		Width(6).
		Render(d)
}

// fmtTimeStyle returns a timestamp in the default style.
func fmtTimeStyle(t string) string {
	return lipgloss.NewStyle().Foreground(colorGray).Render(t)
}

// fmtToolStyle returns a colored tool name.
func fmtToolStyle(t string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#000", Dark: "#EEE"}).Render(t)
}
