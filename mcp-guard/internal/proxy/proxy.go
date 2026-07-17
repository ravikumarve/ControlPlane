package proxy

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/matrix/mcp-guard/internal/audit"
	"github.com/matrix/mcp-guard/internal/hitl"
	"github.com/matrix/mcp-guard/internal/policy"
	"github.com/matrix/mcp-guard/internal/schema"
)

// Options configures the MCP Guard proxy.
type Options struct {
	Mode         string           // stdio | tcp | both
	Listen       string           // TCP listen address
	Upstream     string           // TCP upstream target
	Policy       *policy.Engine
	Audit        *audit.Logger
	SchemaPinner *schema.Pinner
	HITL         *hitl.Handler
}

// Proxy is the main MCP Guard proxy instance.
type Proxy struct {
	opts    Options
	mu      sync.Mutex
	stopped bool
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
	log.Info().Msg("proxy stop signal sent")
}

// isStopped checks if a stop has been requested.
func (p *Proxy) isStopped() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stopped
}
