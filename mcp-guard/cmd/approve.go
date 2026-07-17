package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/matrix/mcp-guard/internal/config"
	"github.com/matrix/mcp-guard/internal/hitl"
)

var approveCmd = &cobra.Command{
	Use:   "approve [request-id]",
	Short: "Approve or deny a pending HITL request",
	Long: `Approve or deny a human-in-the-loop approval request.

Examples:
  mcp-guard approve abc-123              # Approve a request
  mcp-guard approve abc-123 --deny       # Deny a request
  mcp-guard approve list                 # List pending requests
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("request-id or 'list' required")
		}

		deny, _ := cmd.Flags().GetBool("deny")

		if args[0] == "list" {
			// List pending
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if cfg.HITL == nil || cfg.HITL.WebhookURL == "" {
				fmt.Println("HITL not configured (no webhook_url)")
				return nil
			}
			h := hitl.NewHandler(cfg.HITL)
			pending := h.ListPending()
			if len(pending) == 0 {
				fmt.Println("No pending approval requests.")
				return nil
			}
			for _, req := range pending {
				fmt.Printf("  %s: %s called %s (risk: %.2f)\n",
					req.ID, req.Identity, req.Tool, req.RiskScore)
			}
			return nil
		}

		action := "approved"
		if deny {
			action = "denied"
		}

		fmt.Printf("Request %s %s\n", args[0], action)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(approveCmd)
	approveCmd.Flags().Bool("deny", false, "Deny the request instead of approving")
}
