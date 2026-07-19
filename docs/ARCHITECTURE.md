# ControlPlane AI — MCP Security Gateway: Architecture Document
**Version**: 1.0 | **Date**: 2026-07-17

## 1. High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     AI Agent (Claude Code, Cursor, etc.)    │
│                           (MCP Client)                       │
└──────────────────────────┬──────────────────────────────────┘
                           │ JSON-RPC (stdio or TCP)
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    MCP Guard Proxy                            │
│                                                               │
│  ┌──────────┐   ┌──────────┐   ┌────────────┐               │
│  │Transport │──▶│ Policy   │──▶│ Schema     │               │
│  │ Layer    │   │ Engine   │   │ Pinner     │               │
│  └──────────┘   └────┬─────┘   └─────┬──────┘               │
│                      │               │                       │
│                      ▼               ▼                       │
│  ┌──────────────────────────────────────┐                    │
│  │         Decision Router               │                    │
│  │  ┌────────┐  ┌─────────┐  ┌────────┐ │                    │
│  │  │ ALLOW  │  │  BLOCK  │  │ HITL   │ │                    │
│  │  └───┬────┘  └────┬────┘  └───┬────┘ │                    │
│  └──────┼────────────┼────────────┼──────┘                    │
│         │            │            │                            │
│         ▼            ▼            ▼                            │
│  ┌──────────────────────────────────────┐                    │
│  │         Audit Logger (JSONL + HMAC)   │                    │
│  └──────────────────────────────────────┘                    │
└──────────────────────────┬──────────────────────────────────┘
                           │ JSON-RPC (forwarded or blocked)
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    MCP Server (tool executor)                 │
└─────────────────────────────────────────────────────────────┘
```

## 2. Data Flow (Per Tool Call)

```
Agent ──JSON-RPC Request──▶ Transport Layer
                                │
                                ▼
                    1. Parse method & params
                    2. Identify caller identity
                                │
                                ▼
                    3. Policy Engine: evaluate
                       - Is identity allowed?
                       - Is tool allowed for identity?
                       - Are params within constraints?
                       - Rate limit check?
                                │
                    ┌───────────┼───────────┐
                    ▼           ▼           ▼
                 ALLOW       BLOCK       HITL
                    │           │           │
                    ▼           ▼           ▼
              Schema Pin   Log: BLOCK   Send Webhook
              Check Hash                     │
                    │                  ┌─────┘
                    ▼                  │
               Forward to MCP    ┌─────┘
               Server            │
                    │            ▼
                    ▼       Pending Approval
              Receive Response    │
                    │       ┌──── Approved?
                    ▼       ▼       ▼
               Log: ALLOW   YES     NO
                    │       │       │
                    ▼       ▼       ▼
              Return to   Forward  Return
              Agent       + Log    Error
```

## 3. Configuration (YAML Schema)

```yaml
# mcp-guard.yaml
version: "1"

proxy:
  mode: stdio          # stdio | tcp | both
  listen: ":8443"      # TCP listen address (mode: tcp)
  upstream: ""         # TCP upstream (mode: tcp)

policies:
  - name: read-only
    match:
      identity: "*"                 # all agents
      tools: ["read_*", "search_*", "get_*", "list_*"]
    action: allow
    rate_limit: 100/min

  - name: payment-agent
    match:
      identity: "payment-bot"
      tools: ["execute_payout", "refund"]
    action: hitl                     # requires human approval
    constraints:
      max_amount: 500
      channels: ["slack", "email"]

  - name: block-dangerous
    match:
      tools: ["drop_table", "rm_rf", "exec_shell"]
    action: block
    alert: true

schema_pinning:
  enabled: true
  mode: block                        # block | warn | log
  store: ".mcp-guard/pins.json"

audit:
  path: "/var/log/mcp-guard/audit.jsonl"
  hmac_key: "${MCP_GUARD_HMAC_KEY}"  # from env
  rotation: 100MB

hitl:
  webhook_url: "${SLACK_WEBHOOK_URL}"
  timeout: 5m                        # auto-deny after 5 min
  channels: [slack, email]
```

## 4. Core Data Structures (Go)

```go
// Policy definition (from YAML)
type Policy struct {
    Name        string            `yaml:"name"`
    Match       PolicyMatch       `yaml:"match"`
    Action      string            `yaml:"action"` // allow | block | hitl
    Constraints *Constraints      `yaml:"constraints,omitempty"`
    RateLimit   string            `yaml:"rate_limit,omitempty"`
    Alert       bool              `yaml:"alert,omitempty"`
}

type PolicyMatch struct {
    Identity string   `yaml:"identity"` // glob pattern
    Tools    []string `yaml:"tools"`    // glob patterns
}

type Constraints struct {
    MaxAmount float64 `yaml:"max_amount,omitempty"`
    // Future: allowed_ips, time_window, etc.
}

// Audit log entry (JSONL)
type AuditEntry struct {
    ID        string    `json:"id"`
    Timestamp time.Time `json:"ts"`
    Identity  string    `json:"identity"`
    Tool      string    `json:"tool"`
    Params    any       `json:"params"`
    Decision  string    `json:"decision"` // allow | block | hitl | pending | denied
    Duration  int64     `json:"duration_ms"`
    HMAC      string    `json:"hmac"`
    PrevHMAC  string    `json:"prev_hmac"` // chain
}

// HITL approval request
type HITLRequest struct {
    ID        string    `json:"id"`
    Identity  string    `json:"identity"`
    Tool      string    `json:"tool"`
    Params    any       `json:"params"`
    RiskScore float64   `json:"risk_score"`
    Status    string    `json:"status"` // pending | approved | denied
    CreatedAt time.Time `json:"created_at"`
    ExpiresAt time.Time `json:"expires_at"`
}

// Schema pin
type SchemaPin struct {
    ServerURL   string            `json:"server_url"`
    ToolHashes  map[string]string `json:"tool_hashes"` // toolName -> SHA256
    PinnedAt    time.Time         `json:"pinned_at"`
    LastChecked time.Time         `json:"last_checked"`
}
```

## 5. Component Breakdown

### 5.1 Transport Layer (`internal/proxy/`)
- **stdio proxy**: Wraps `io.ReadWriteCloser` — reads JSON-RPC from stdin, forwards to MCP server, intercepts both directions
- **TCP proxy**: `net.Listener` — accepts client connections, TLS optional, multiplexes sessions
- **Message parser**: Streaming JSON decoder for JSON-RPC 2.0 (handles partial frames)

### 5.2 Policy Engine (`internal/policy/`)
- **Matcher**: Glob-pattern matching on identity + tool name (Go `path.Match` or `glob`)
- **Evaluator**: Walks policy list in order, first match wins
- **Rate limiter**: Token bucket per policy (per-identity, per-tool)
- **Constraints**: Custom check functions (e.g., `params.amount <= 500`)

### 5.3 Schema Pinner (`internal/schema/`)
- On first `tools/list` response, compute SHA-256 of each tool definition
- Store in `pins.json`
- On subsequent connections, recompute and compare — block/warn on drift

### 5.4 Audit Logger (`internal/audit/`)
- Appends JSONL to file
- HMAC-SHA256 chain: `entry.HMAC = HMAC(entry.JSON + prevEntry.HMAC)`
- Verification: `mcp-guard logs --verify` re-computes chain

### 5.5 HITL Handler (`internal/hitl/`)
- Manages pending approval requests
- Sends webhook payload to configured URL
- Listens for callback (webhook reverse or polling)
- Auto-expires after timeout

## 6. Package Layout

```
mcp-guard/
├── main.go
├── go.mod / go.sum
├── cmd/
│   ├── root.go         # Cobra root, --config flag
│   ├── serve.go        # Run proxy daemon
│   ├── init.go         # Generate config template
│   ├── status.go       # Show daemon status
│   ├── logs.go         # Tail/inspect audit log
│   ├── policy.go       # Policy subcommands (list, apply, test)
│   ├── pin.go          # Schema pin subcommands
│   └── approve.go      # Approve/deny HITL requests
├── internal/
│   ├── proxy/
│   │   ├── proxy.go         # Main proxy orchestrator
│   │   ├── stdio.go         # stdio transport
│   │   ├── tcp.go           # TCP transport
│   │   └── message.go       # JSON-RPC types & parser
│   ├── config/
│   │   └── config.go        # YAML loading, env interpolation
│   ├── policy/
│   │   ├── types.go         # Policy structs
│   │   ├── engine.go        # Evaluation logic
│   │   └── matcher.go       # Glob matching
│   ├── audit/
│   │   ├── logger.go        # JSONL writer with HMAC
│   │   └── verifier.go      # Chain verification
│   ├── schema/
│   │   └── pinner.go        # Tool hash + drift detection
│   ├── hitl/
│   │   ├── handler.go       # Approval request lifecycle
│   │   └── webhook.go       # Slack/email dispatcher
│   └── version/
│       └── version.go       # Build info
```

## 7. Security Considerations
- **HMAC key**: Must be set via `MCP_GUARD_HMAC_KEY` env var, never in YAML
- **Config file**: Should be `chmod 600` — may contain sensitive patterns
- **Proxy blind spot**: stdio mode trusts the parent process; TCP mode should use mTLS
- **Policy order**: First-match-wins — put deny rules before allow rules
- **Log rotation**: Built-in size-based rotation; no log loss on restart (append-only)

## 8. Performance Targets
| Operation | Target | Method |
|-----------|--------|--------|
| Per-call latency | < 100µs | Zero-alloc message parsing, no reflection |
| Memory per connection | < 64KB | Streaming JSON decoder |
| Audit throughput | > 10K entries/sec | Async writer with buffer |
| Binary size | < 15MB | Static Go binary, no cgo |
