package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/config"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/schema"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Manage schema pins",
	Long:  `Schema pinning hashes tool definitions to detect supply-chain poisoning.`,
}

var pinListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pinned server schemas",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if !cfg.SchemaPinning.Enabled {
			fmt.Println("Schema pinning is disabled in config.")
			return nil
		}

		pinner := schema.NewPinner(cfg.SchemaPinning)
		pins, err := pinner.LoadPins()
		if err != nil {
			return fmt.Errorf("failed to load pins: %w", err)
		}

		if len(pins) == 0 {
			fmt.Println("No schemas pinned yet. Connect to an MCP server to auto-pin.")
			return nil
		}

		for _, pin := range pins {
			fmt.Printf("  Server: %s\n", pin.ServerURL)
			fmt.Printf("  Tools:  %d hashes stored\n", len(pin.ToolHashes))
			fmt.Printf("  Pinned: %s\n", pin.PinnedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}
		return nil
	},
}

var pinVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Check all pinned schemas for drift",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		pinner := schema.NewPinner(cfg.SchemaPinning)
		drifting, err := pinner.VerifyAll()
		if err != nil {
			return fmt.Errorf("verification failed: %w", err)
		}

		if len(drifting) == 0 {
			fmt.Println("✅ All pinned schemas match — no drift detected")
		} else {
			fmt.Println("❌ Schema drift detected:")
			for _, d := range drifting {
				fmt.Printf("  - %s: %s\n", d.Server, d.Tool)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pinCmd)
	pinCmd.AddCommand(pinListCmd)
	pinCmd.AddCommand(pinVerifyCmd)
}
