package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/audit"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/config"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View or verify the audit log",
	Long: `Display recent audit log entries or verify the HMAC chain integrity.

Examples:
  mcp-guard logs --tail        # Watch live log entries
  mcp-guard logs --verify      # Verify HMAC chain integrity
  mcp-guard logs --last 50     # Show last 50 entries
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tail, _ := cmd.Flags().GetBool("tail")
		verify, _ := cmd.Flags().GetBool("verify")
		last, _ := cmd.Flags().GetInt("last")

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		auditPath := cfg.Audit.Path
		if auditPath == "" {
			auditPath = "/var/log/mcp-guard/audit.jsonl"
		}

		if verify {
			return verifyAuditChain(auditPath)
		}

		file, err := os.Open(auditPath)
		if err != nil {
			return fmt.Errorf("failed to open audit log: %w", err)
		}
		defer file.Close()

		if tail {
			// Follow the file (like tail -f)
			return tailAuditLog(auditPath)
		}

		scanner := bufio.NewScanner(file)
		lines := []string{}
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if last > 0 && last < len(lines) {
			lines = lines[len(lines)-last:]
		}

		for _, line := range lines {
			fmt.Println(line)
		}
		return nil
	},
}

func verifyAuditChain(path string) error {
	verifier := audit.NewVerifier(path)
	valid, err := verifier.Verify()
	if err != nil {
		return fmt.Errorf("verification error: %w", err)
	}
	if valid {
		fmt.Println("✅ Audit log HMAC chain is INTACT — no tampering detected")
	} else {
		fmt.Println("❌ Audit log HMAC chain is BROKEN — possible tampering!")
	}
	return nil
}

func tailAuditLog(path string) error {
	// Simple follow implementation using polling
	fmt.Printf("Tailing %s (Ctrl+C to stop)...\n", path)
	// In a real implementation, use fsnotify for inotify-based watching
	// For MVP, we just read the file and block
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Seek to end
	file.Seek(0, 2)

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return scanner.Err()
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().Bool("tail", false, "Follow log output (like tail -f)")
	logsCmd.Flags().Bool("verify", false, "Verify HMAC chain integrity")
	logsCmd.Flags().Int("last", 0, "Show last N entries")
}
