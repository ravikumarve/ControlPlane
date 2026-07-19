package proxy

import (
	"fmt"
	"net"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/alert"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/audit"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/hitl"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/inject"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/policy"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/ratelimit"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/schema"
)

// Options configures the MCP Guard proxy.
type Options struct {
	Mode           string                // stdio | tcp | both
	Listen         string                // TCP listen address
	Upstream       string                // TCP upstream target
	Policy         *policy.Engine
	Audit          *audit.Logger
	SchemaPinner   *schema.Pinner
	HITL           *hitl.Handler
	InjectDetector  *inject.Detector      // nil = skip injection scan
	RateLimiter     *ratelimit.KeyedLimiter  // nil = no rate limit
	AlertDispatcher *alert.Dispatcher     // nil = no alerts
}

// Proxy is the main MCP Guard proxy instance.
type Proxy struct {
	opts     Options
	mu       sync.Mutex
	stopped  bool
	listener net.Listener
}

// New creates a new MCP Guard proxy.
func New(opts Options) *Proxy {
	return &Proxy{opts: opts}
}

// Start begins proxying MCP traffic based on the configured mode.
func (p *Proxy) Start() error {
	log.Info().Str("mode", p.opts.Mode).Msg("proxy starting")

	switch p.opts.Mode {
	case "stdio":
		return p.startStdio()
	case "tcp":
		return p.startTCP()
	default:
		return fmt.Errorf("unsupported proxy mode: %s", p.opts.Mode)
	}
}

// Stop signals the proxy to shut down gracefully.
func (p *Proxy) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stopped = true
	if p.listener != nil {
		p.listener.Close()
	}
	log.Info().Msg("proxy stopped")
}

// isStopped checks if a stop has been requested.
func (p *Proxy) isStopped() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stopped
}
