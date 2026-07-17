package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/matrix/mcp-guard/internal/version"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "mcp-guard",
	Short: "MCP Guard — Lightweight security sidecar for MCP agents",
	Long: `MCP Guard sits between AI agents and MCP servers, enforcing
tool-level access control, schema pinning, and audit logging.

Single binary. No Kubernetes. No SaaS. 5-minute deploy.`,
	Version: version.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./mcp-guard.yaml or /etc/mcp-guard/config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/mcp-guard/")
		viper.SetConfigName("mcp-guard")
	}

	viper.SetEnvPrefix("MCP_GUARD")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config: %s\n", err)
		}
		// Config file is optional for subcommands like `init`, `status`
	}
}
