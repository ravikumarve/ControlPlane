package proxy

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/matrix/mcp-guard/internal/audit"
	"github.com/matrix/mcp-guard/internal/hitl"
	"github.com/matrix/mcp-guard/internal/policy"
)

// connContext carries the shared state for a single proxy connection.
type connContext struct {
	client   net.Conn
	upstream net.Conn
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// startTCP runs the proxy in TCP mode.
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

// handleConnection manages a full duplex proxy session.
func (p *Proxy) handleConnection(client net.Conn) {
	defer client.Close()

	upstream, err := net.DialTimeout("tcp", p.opts.Upstream, 10*time.Second)
	if err != nil {
		log.Error().Err(err).Msg("dial upstream failed")
		return
	}
	defer upstream.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cc := &connContext{
		client:   client,
		upstream: upstream,
		ctx:      ctx,
		cancel:   cancel,
	}

	log.Debug().Str("remote", client.RemoteAddr().String()).Msg("new connection")

	cc.wg.Add(2)
	go p.forwardClient(cc)
	go p.forwardUpstream(cc)
	cc.wg.Wait()
}

// forwardClient reads requests from the client, applies policy, and forwards allowed calls.
func (p *Proxy) forwardClient(cc *connContext) {
	defer cc.wg.Done()
	defer cc.cancel() // signal upstream goroutine to stop

	scanner := bufio.NewScanner(cc.client)
	buf := make([]byte, 0, 256*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		select {
		case <-cc.ctx.Done():
			return
		default:
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		req, err := ParseRequest([]byte(line))
		if err != nil {
			fmt.Fprintln(cc.upstream, line)
			continue
		}

		toolName := p.handleClientRequest(cc, req, line)
		_ = toolName
	}

	if err := scanner.Err(); err != nil {
		log.Debug().Err(err).Msg("client scanner error")
	}
}

// handleClientRequest processes a single client request and returns the decision.
func (p *Proxy) handleClientRequest(cc *connContext, req *JSONRPCRequest, rawLine string) string {
	identity := detectIdentity(req)
	toolName := ExtractToolName(req)
	start := time.Now()

	// Protocol handshake — forward transparently
	if req.Method == "initialize" || req.Method == "notifications/initialized" {
		fmt.Fprintln(cc.upstream, rawLine)
		return toolName
	}

	// Rate limit check — per-identity token bucket
	if rl := p.opts.RateLimiter; rl != nil && !rl.Allow(identity) {
		entry := audit.AuditEntry{
			ID:        uuid.New().String(),
			Timestamp: time.Now(),
			Identity:  identity,
			Tool:      toolName,
			Params:    req.Params,
			Decision:  "block",
			Duration:  time.Since(start).Milliseconds(),
			Reason:    "rate limit exceeded",
		}
		p.logAudit(entry)
		fmt.Fprintln(cc.client, string(NewBlockedResponse(req.ID, "rate limit exceeded")))
		log.Warn().Str("identity", identity).Str("tool", toolName).Msg("rate limited")
		return toolName
	}

	// Injection detection — run before policy evaluation
	if injector := p.opts.InjectDetector; injector != nil {
		if sr := injector.ScanParams(toolName, req.Params); sr.Injection.Detected {
			entry := audit.AuditEntry{
				ID:        uuid.New().String(),
				Timestamp: time.Now(),
				Identity:  identity,
				Tool:      toolName,
				Params:    req.Params,
				Decision:  "block",
				Duration:  time.Since(start).Milliseconds(),
				Reason:    sr.Injection.Reason,
			}
			p.logAudit(entry)
			fmt.Fprintln(cc.client, string(NewBlockedResponse(req.ID, "injection detected: "+sr.Injection.Reason)))
			log.Warn().Str("tool", toolName).Str("reason", sr.Injection.Reason).Msg("injection blocked")
			return toolName
		}
	}

	// Policy evaluation
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
		fmt.Fprintln(cc.client, string(NewBlockedResponse(req.ID, decision.Reason)))
		log.Warn().Str("tool", toolName).Str("reason", decision.Reason).Msg("blocked")

	case policy.ActionHITL:
		entry.Decision = "pending"
		if p.opts.HITL != nil {
			p.opts.HITL.Submit(hitl.Request{
				ID:       entry.ID,
				Identity: identity,
				Tool:     toolName,
				Params:   req.Params,
				RawData:  rawLine,
			})
		}
		msg := "requires human approval"
		fmt.Fprintln(cc.client, string(NewBlockedResponse(req.ID, msg)))
		log.Info().Str("id", entry.ID).Str("tool", toolName).Msg("pending HITL")

	case policy.ActionAllow:
		entry.Decision = "allow"
		entry.Duration = time.Since(start).Milliseconds()
		p.logAudit(entry)
		fmt.Fprintln(cc.upstream, rawLine)
	}

	return toolName
}

// forwardUpstream reads responses from the upstream server and relays them to the client.
// This also captures tools/list responses for schema pinning.
func (p *Proxy) forwardUpstream(cc *connContext) {
	defer cc.wg.Done()

	scanner := bufio.NewScanner(cc.upstream)
	buf := make([]byte, 0, 256*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		select {
		case <-cc.ctx.Done():
			return
		default:
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Try to parse as JSON-RPC response
		resp, err := ParseResponse([]byte(line))
		if err == nil && resp.Result != nil {
			// Could inspect response for tools/list — future schema pinning
			_ = resp
		}

		fmt.Fprintln(cc.client, line)
	}

	if err := scanner.Err(); err != nil {
		log.Debug().Err(err).Msg("upstream scanner error")
	}
}
