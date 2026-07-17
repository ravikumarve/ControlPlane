package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/matrix/mcp-guard/internal/audit"
	"github.com/matrix/mcp-guard/internal/hitl"
	"github.com/matrix/mcp-guard/internal/policy"
)

// startTCP runs the proxy in TCP mode:
//   - Listens on a TCP port for MCP client connections
//   - Forwards JSON-RPC to upstream MCP server
//   - Intercepts both directions for policy + audit
func (p *Proxy) startTCP() error {
	listener, err := net.Listen("tcp", p.opts.Listen)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", p.opts.Listen, err)
	}
	p.mu.Lock()
	p.listener = listener
	p.mu.Unlock()
	defer listener.Close()

	log.Info().Str("addr", p.opts.Listen).Str("upstream", p.opts.Upstream).Msg("TCP proxy listening")

	for {
		client, err := listener.Accept()
		if err != nil {
			if p.isStopped() {
				break
			}
			log.Error().Err(err).Msg("accept error")
			continue
		}

		go p.handleConnection(client)
	}

	return nil
}

// handleConnection handles a single TCP client connection.
func (p *Proxy) handleConnection(client net.Conn) {
	defer client.Close()

	upstream, err := net.DialTimeout("tcp", p.opts.Upstream, 10*time.Second)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to upstream")
		return
	}
	defer upstream.Close()

	log.Debug().Str("remote", client.RemoteAddr().String()).Msg("new connection")

	// Client -> Proxy -> Upstream
	go p.forwardClientToUpstream(client, upstream)

	// Upstream -> Proxy -> Client
	p.forwardUpstreamToClient(upstream, client)
}

// forwardClientToUpstream reads from client, applies policy, forwards allowed calls.
func (p *Proxy) forwardClientToUpstream(client, upstream net.Conn) {
	scanner := bufio.NewScanner(client)
	buf := make([]byte, 0, 256*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		req, err := ParseRequest([]byte(line))
		if err != nil {
			// Not a JSON-RPC request, forward transparently
			fmt.Fprintln(upstream, line)
			continue
		}

		start := time.Now()
		identity := detectIdentity(req)
		toolName := ExtractToolName(req)

		// Skip handshake
		if req.Method == "initialize" || req.Method == "notifications/initialized" {
			fmt.Fprintln(upstream, line)
			continue
		}

		decision := p.opts.Policy.Evaluate(identity, toolName)

		entry := audit.AuditEntry{
			ID:        uuid.New().String(),
			Timestamp: time.Now(),
			Identity:  identity,
			Tool:      toolName,
			Params:    req.Params,
		}

		switch decision.Action {
		case policy.ActionBlock:
			entry.Decision = "block"
			entry.Duration = time.Since(start).Milliseconds()
			p.logAudit(entry)
			blocked := NewBlockedResponse(req.ID, decision.Reason)
			fmt.Fprintln(client, string(blocked))
			log.Warn().Str("tool", toolName).Str("reason", decision.Reason).Msg("blocked")

		case policy.ActionHITL:
			entry.Decision = "pending"
			if p.opts.HITL != nil {
				p.opts.HITL.Submit(hitl.Request{
					ID:       entry.ID,
					Identity: identity,
					Tool:     toolName,
					Params:   req.Params,
					RawData:  line,
				})
			}
			blocked := NewBlockedResponse(req.ID, "requires human approval")
			fmt.Fprintln(client, string(blocked))
			log.Info().Str("id", entry.ID).Str("tool", toolName).Msg("sent for approval")

		case policy.ActionAllow:
			entry.Decision = "allow"
			entry.Duration = time.Since(start).Milliseconds()
			p.logAudit(entry)
			fmt.Fprintln(upstream, line)
		}
	}
}

// forwardUpstreamToClient relays upstream responses back to client.
func (p *Proxy) forwardUpstreamToClient(upstream, client net.Conn) {
	_, err := io.Copy(client, upstream)
	if err != nil {
		log.Debug().Err(err).Msg("upstream->client copy done")
	}
}
