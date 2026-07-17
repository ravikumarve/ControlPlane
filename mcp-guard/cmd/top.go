package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/matrix/mcp-guard/internal/config"
	"github.com/matrix/mcp-guard/internal/tui"
	"github.com/matrix/mcp-guard/internal/version"
)

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Live TUI dashboard for MCP Guard",
	Long: `Opens a real-time terminal dashboard showing
traffic statistics and live audit log feed.

Controls:
  q / Ctrl+C  Quit
  p           Pause/resume live feed
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config to find audit log path
		auditPath := "/var/log/mcp-guard/audit.jsonl"
		mode := "stdio"

		cfg, err := config.Load()
		if err == nil && cfg != nil {
			if cfg.Audit.Path != "" {
				auditPath = cfg.Audit.Path
			}
			if cfg.Proxy.Mode != "" {
				mode = cfg.Proxy.Mode
			}
		}

		// Check if the audit file exists or the daemon is running
		if _, err := os.Stat(auditPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "⚠  Audit log not found at %s\n", auditPath)
			fmt.Fprintf(os.Stderr, "   Start mcp-guard serve first, or check your config path.\n")
			fmt.Fprintf(os.Stderr, "   Using anyway — will connect when file appears.\n")
		}

		model := tui.NewModel(auditPath, version.Version, mode)
		p := tea.NewProgram(model, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(topCmd)
	topCmd.Flags().Bool("verbose", false, "verbose logging")
}
