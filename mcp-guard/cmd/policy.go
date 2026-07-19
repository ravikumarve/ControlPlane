package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/config"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/policy"
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Manage security policies",
	Long:  `List, apply, or test MCP Guard policies.`,
}

var policyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List loaded policies",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if len(cfg.Policies) == 0 {
			fmt.Println("No policies defined in config.")
			return nil
		}

		fmt.Printf("%-20s %-10s %-10s %s\n", "NAME", "ACTION", "IDENTITY", "TOOLS")
		fmt.Println(strings.Repeat("-", 80))
		for _, p := range cfg.Policies {
			toolsStr := strings.Join(p.Match.Tools, ", ")
			if len(toolsStr) > 40 {
				toolsStr = toolsStr[:37] + "..."
			}
			fmt.Printf("%-20s %-10s %-10s %s\n", p.Name, p.Action, p.Match.Identity, toolsStr)
		}
		return nil
	},
}

var policyTestCmd = &cobra.Command{
	Use:   "test [tool-name]",
	Short: "Dry-run a policy against a sample tool call",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		identity, _ := cmd.Flags().GetString("identity")
		if identity == "" {
			identity = "test-agent"
		}

		engine := policy.NewEngine(cfg.Policies)
		decision := engine.Evaluate(identity, args[0], nil)

		fmt.Printf("  Identity: %s\n", identity)
		fmt.Printf("  Tool:     %s\n", args[0])
		fmt.Printf("  Decision: %s\n", decision.Action)
		if decision.PolicyName != "" {
			fmt.Printf("  Matched:  %s\n", decision.PolicyName)
		}
		if decision.Reason != "" {
			fmt.Printf("  Reason:   %s\n", decision.Reason)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(policyCmd)
	policyCmd.AddCommand(policyListCmd)
	policyCmd.AddCommand(policyTestCmd)

	policyTestCmd.Flags().String("identity", "", "Agent identity to test (default: test-agent)")
}
