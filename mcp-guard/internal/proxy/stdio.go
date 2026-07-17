package proxy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/matrix/mcp-guard/internal/audit"
	"github.com/matrix/mcp-guard/internal/hitl"
	"github.com/matrix/mcp-guard/internal/policy"
)

// startStdio runs the proxy in stdio mode:
//   - Reads JSON-RPC from stdin
//   - Intercepts and evaluates against policy
//   - Forwards allowed calls via stdout to parent process
//   - Blocks/reports disallowed calls
func (p *Proxy) startStdio() error {
	scanner := bufio.NewScanner(os.Stdin)
	// Use a large buffer for potentially large JSON-RPC messages
	buf := make([]byte, 0, 256*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		if p.isStopped() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if err := p.handleMessage([]byte(line)); err != nil {
			log.Error().Err(err).Msg("failed to handle message")
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stdin scanner error: %w", err)
	}

	return nil
}

// handleMessage processes a single JSON-RPC message.
func (p *Proxy) handleMessage(data []byte) error {
	req, err := ParseRequest(data)
	if err != nil {
		// Not a valid JSON-RPC request — forward transparently
		log.Debug().Err(err).Msg("forwarding non-request message")
		fmt.Println(string(data))
		return nil
	}

	start := time.Now()
	identity := detectIdentity(req)
	toolName := ExtractToolName(req)

	// Special handling for tools/list (needed for schema pinning)
	if req.Method == "tools/list" {
		return p.handleToolsList(req, data)
	}

	// Skip internal methods
	if req.Method == "notifications/initialized" || req.Method == "initialize" {
		fmt.Println(string(data))
		return nil
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
		fmt.Println(string(NewBlockedResponse(req.ID, "rate limit exceeded")))
		log.Warn().Str("identity", identity).Str("tool", toolName).Msg("rate limited")
		return nil
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
			fmt.Println(string(NewBlockedResponse(req.ID, "injection detected: "+sr.Injection.Reason)))
			log.Warn().Str("tool", toolName).Str("reason", sr.Injection.Reason).Msg("injection blocked")
			return nil
		}
	}

	// Evaluate policy using the extracted tool name (for tools/call, this is the actual tool name)
	decision := p.opts.Policy.Evaluate(identity, toolName)

	entry := audit.AuditEntry{
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Identity:  identity,
		Tool:      toolName,
		Params:    req.Params,
		Duration:  time.Since(start).Milliseconds(),
	}

	switch decision.Action {
	case policy.ActionBlock:
		entry.Decision = "block"
		p.logAudit(entry)
		blocked := NewBlockedResponse(req.ID, decision.Reason)
		fmt.Println(string(blocked))
		log.Warn().
			Str("identity", identity).
			Str("tool", toolName).
			Str("reason", decision.Reason).
			Msg("blocked tool call")

	case policy.ActionHITL:
		if p.opts.HITL != nil {
			entry.Decision = "pending"
			p.logAudit(entry)
			p.opts.HITL.Submit(hitl.Request{
				ID:        entry.ID,
				Identity:  identity,
				Tool:      toolName,
				Params:    req.Params,
				RawData:   string(data),
				RiskScore: decision.RiskScore,
			})
			blocked := NewBlockedResponse(req.ID, "requires human approval — request sent")
			fmt.Println(string(blocked))
			log.Info().
				Str("id", entry.ID).
				Str("identity", identity).
				Str("tool", toolName).
				Msg("sent for human approval")
		} else {
			// HITL not configured, block by default
			entry.Decision = "block"
			p.logAudit(entry)
			blocked := NewBlockedResponse(req.ID, "requires human approval but HITL not configured")
			fmt.Println(string(blocked))
		}

	case policy.ActionAllow:
		// Schema pinning: update hashes if enabled
		if p.opts.SchemaPinner != nil {
			go p.opts.SchemaPinner.CheckAndPin(toolName, data)
		}

		entry.Decision = "allow"
		entry.Duration = time.Since(start).Milliseconds()
		p.logAudit(entry)
		// Forward the original request
		fmt.Println(string(data))
		log.Debug().
			Str("identity", identity).
			Str("tool", toolName).
			Msg("allowed tool call")
	}

	return nil
}

// handleToolsList intercepts tools/list responses for schema pinning.
func (p *Proxy) handleToolsList(req *JSONRPCRequest, data []byte) error {
	// Forward to get the tool definitions
	fmt.Println(string(data))

	// After forwarding, the response comes back on stdout
	// In stdio mode, we can't easily intercept the response without modifying the protocol.
	// For TCP mode, we can. For MVP, schema pinning works best in TCP mode.
	// In stdio mode, we rely on the first tools/list response capture.

	log.Debug().Msg("forwarded tools/list request")
	return nil
}

// detectIdentity extracts the caller identity from the request context.
// For MVP, this is a simple heuristic. In production, it would use OAuth tokens
// or session info from the MCP initialization.
func detectIdentity(req *JSONRPCRequest) string {
	if req.Params == nil {
		return "anonymous"
	}

	// Try to extract identity from params
	params, ok := req.Params.(map[string]any)
	if !ok {
		return "anonymous"
	}

	if identity, ok := params["_identity"]; ok {
		if s, ok := identity.(string); ok {
			return s
		}
	}

	return "anonymous"
}

// logAudit writes an audit entry asynchronously.
func (p *Proxy) logAudit(entry audit.AuditEntry) {
	if p.opts.Audit != nil {
		if err := p.opts.Audit.Write(entry); err != nil {
			log.Error().Err(err).Msg("failed to write audit entry")
		}
	}
}
