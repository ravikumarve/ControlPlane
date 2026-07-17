package cmd

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/matrix/mcp-guard/internal/config"
	"github.com/matrix/mcp-guard/internal/tui"
	"github.com/matrix/mcp-guard/internal/version"
)

var auditPathFlag string

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Live TUI dashboard for MCP Guard",
	Long: `Opens a real-time terminal dashboard showing
traffic statistics and live audit log feed.

Supported features:
  • Live counters (total / allowed / blocked / HITL / rate limited / injection blocks)
  • Per-tool breakdown with mini bar charts
  • Per-identity breakdown
  • Activity sparkline (5s buckets, 2min rolling window)
  • Color-coded live audit log feed
  • Pause/resume with 'p', scroll with ↑↓

Controls:
  q / Ctrl+C  Quit
  p           Pause/resume live feed
  ↑ ↓         Scroll audit log
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		auditPath := auditPathFlag
		mode := "stdio"

		// If no explicit path, try config
		if auditPath == "" {
			auditPath = "/var/log/mcp-guard/audit.jsonl"
			cfg, err := config.Load()
			if err == nil && cfg != nil {
				if cfg.Audit.Path != "" {
					auditPath = cfg.Audit.Path
				}
				if cfg.Proxy.Mode != "" {
					mode = cfg.Proxy.Mode
				}
			}
		}

		// Check if the audit file exists
		if _, err := os.Stat(auditPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "⚠  Audit log not found at %s\n", auditPath)
			fmt.Fprintf(os.Stderr, "   Start mcp-guard serve first, or use --audit-path.\n")
			fmt.Fprintf(os.Stderr, "   Dashboard will connect when file appears.\n\n")
			time.Sleep(1 * time.Second)
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
	topCmd.Flags().StringVar(&auditPathFlag, "audit-path", "", "path to audit JSONL (default: from config or /var/log/mcp-guard/audit.jsonl)")
	topCmd.Flags().Bool("verbose", false, "verbose logging")
}
