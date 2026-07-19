package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/admin"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/alert"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/audit"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/config"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/hitl"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/policy"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/proxy"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/schema"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP Guard proxy daemon",
	Long: `Starts the MCP Guard proxy daemon, intercepting JSON-RPC traffic
between AI agents and MCP servers based on configured policies.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize logger
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		// Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		log.Info().Str("mode", cfg.Proxy.Mode).Msg("starting mcp-guard proxy")

		// Initialize components
		policyEngine := policy.NewEngine(cfg.Policies)
		auditLogger, err := audit.NewLogger(cfg.Audit)
		if err != nil {
			return fmt.Errorf("failed to initialize audit logger: %w", err)
		}
		defer auditLogger.Close()

		var schemaPinner *schema.Pinner
		if cfg.SchemaPinning.Enabled {
			schemaPinner = schema.NewPinner(cfg.SchemaPinning)
		}

		var hitlHandler *hitl.Handler
		if cfg.HITL != nil && cfg.HITL.WebhookURL != "" {
			hitlHandler = hitl.NewHandler(cfg.HITL)
		}

		var alertDispatcher *alert.Dispatcher
		if cfg.Alert != nil && cfg.Alert.WebhookURL != "" {
			channel := cfg.Alert.Channel
			if channel == "" {
				channel = "generic"
			}
			alertDispatcher = alert.NewDispatcher(cfg.Alert.WebhookURL, channel)
			log.Info().Str("channel", channel).Str("webhook", cfg.Alert.WebhookURL).Msg("alert dispatcher configured")
		}

		// Initialize proxy stats
		proxyStats := &proxy.Stats{}

		// Build proxy
		p := proxy.New(proxy.Options{
			Mode:            cfg.Proxy.Mode,
			Listen:          cfg.Proxy.Listen,
			Upstream:        cfg.Proxy.Upstream,
			Policy:          policyEngine,
			Audit:           auditLogger,
			SchemaPinner:    schemaPinner,
			HITL:            hitlHandler,
			AlertDispatcher: alertDispatcher,
			Stats:           proxyStats,
		})

		// Handle shutdown
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		// Start admin server if enabled
		var adminServer *admin.Server
		if cfg.Admin.Enabled {
			adminServer = admin.New(
				cfg.Admin.APIKey,
				cfg.Admin.Listen,
				policyEngine,
				auditLogger,
				proxyStats,
				func() interface{} { return cfg },
			)
			go func() {
				if err := adminServer.Start(); err != nil {
					log.Error().Err(err).Msg("admin server stopped")
				}
			}()
		}

		// Serve
		if err := p.Start(); err != nil {
			return fmt.Errorf("proxy error: %w", err)
		}

		// Shutdown: stop admin then proxy
		<-sigCh
		log.Info().Msg("shutting down...")
		if adminServer != nil {
			adminServer.Stop()
		}
		p.Stop()

		log.Info().Msg("proxy stopped")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().BoolP("verbose", "v", false, "verbose logging")
	serveCmd.Flags().String("config", "", "path to config file")
}
