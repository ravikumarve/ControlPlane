package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the mcp-guard daemon status",
	Long:  `Displays whether mcp-guard is running, uptime, and basic stats.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if the health/metrics endpoint is reachable
		port := ":8443"
		client := &http.Client{Timeout: 2 * time.Second}

		resp, err := client.Get(fmt.Sprintf("http://localhost%s/health", port))
		if err != nil {
			fmt.Println("❌ mcp-guard is NOT running")
			fmt.Println("  Start it with: mcp-guard serve")
			return nil
		}
		defer resp.Body.Close()

		fmt.Println("✅ mcp-guard is RUNNING")
		fmt.Printf("  Health endpoint: localhost%s/health\n", port)
		fmt.Println("  Logs: mcp-guard logs --tail")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
