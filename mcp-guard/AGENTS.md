# MCP Guard — Project Context

**Go module:** `github.com/matrix/mcp-guard`
**Build:** Go 1.24, single 8.5MB binary (no CGO, no Docker)
**Repo root:** `ravikumarve/ControlPlane` → `mcp-guard/`

## Architecture

```
cmd/              Cobra CLI — 9 commands: serve, top, logs, policy, init, pin, approve, status
internal/
  audit/          HMAC-SHA256 chained JSONL logger (tamper-evident)
  config/         Viper-based YAML loader + Validate() for 9+ error types
  hitl/           Webhook dispatcher for human-in-the-loop approval
  inject/         Prompt injection detector — 10 patterns, homoglyph map, depth/length limits
  policy/         First-match-wins RBAC engine (allow/block/HITL), glob matching
  proxy/          JSON-RPC proxy — stdio mode + TCP mode (scanner-based, shared context)
  ratelimit/      Token bucket per-identity (parses "100/m", "10/s", "1000/h")
  schema/         Schema pinning — SHA-256 drift detection
  testmcp/        Test MCP server + end-to-end TCP integration test
  tui/            Bubble Tea full dashboard — stats, tools, identities, sparkline, feed
  version/        Build version
```

## Request Flow (TCP mode)

```
Client request
  → Handshake?              Forward transparently
  → Rate limit check?       Exceeded → BLOCK + audit (Reason: "rate limit exceeded")
  → Injection detection?    Detected → BLOCK + audit (Reason: injection description)
  → Policy evaluation       ALLOW  → forward to upstream + audit
                            BLOCK  → reject with -32000 + audit + reason
                            HITL   → pending + webhook + audit
```

Inline checks run in this exact order. Each short-circuits on detection.

## Key Decisions

- **No Docker/K8s** — Latitude 3460 CPU constraint is the competitive moat. Target SMBs and solo devs with the same constraints.
- **Single binary** — Go + stdlib only for runtime. Bubble Tea for TUI, Cobra for CLI, viper for config.
- **ANSI strings in hot path** — TUI render loops use raw ANSI escape codes for stats/tables/timelines to avoid lipgloss allocation overhead. Lipgloss used only for panel wrappers and layout.
- **First-match-wins** policy evaluation. Default deny when no policy matches.
- **HMAC chains** — each audit entry includes `hmac` + `prev_hmac` for tamper detection.
- **Config validation** — `Validate()` runs on load, catches invalid mode/actions/duplicate names/bad rate limits/etc.

## Test Strategy

| Package | Tests | What it validates |
|---------|-------|-------------------|
| audit | 7 | HMAC chains, broken chain detection, concurrent writes |
| config | 15 | Env expansion, defaults, Validate() error checking |
| inject | 13 | 10 injection patterns, homoglyph map, depth bomb, length limit |
| policy | 9 | First-match-wins, glob, default deny, identity match |
| proxy | 13 | Parse request/response, blocked response, tool extraction |
| ratelimit | 14 | Token bucket, refill, keyed limiter, rate string parsing |
| testmcp | 1 (5 scns) | End-to-end TCP: allow/block/HITL + audit HMAC integrity |

**Total:** 58+ tests, all passing.

## CLI Reference

```
mcp-guard serve                    Start proxy daemon
mcp-guard top                      Live TUI dashboard
mcp-guard logs --tail              View audit log
mcp-guard policy list              List policies
mcp-guard policy test <tool>       Test policy decision
mcp-guard pin list                 List schema pins
mcp-guard pin verify               Verify pinned schemas
mcp-guard approve <id>             Approve HITL request
mcp-guard init                     Generate default config
mcp-guard status                   Daemon status
```

## Config (mcp-guard.yaml)

```yaml
version: "1"
proxy:
  mode: tcp           # stdio | tcp
  listen: ":8443"
  upstream: "localhost:3000"

policies:
  - name: allow-reads
    action: allow
    match:
      identity: "*"
      tools: ["read_*", "get_*"]

  - name: block-dangerous
    action: block
    match:
      tools: ["delete_*", "exec_*"]

  - name: risky-pay
    action: hitl
    match:
      tools: ["execute_payout"]

audit:
  path: /var/log/mcp-guard/audit.jsonl
  hmac_key: "${MCP_GUARD_HMAC_KEY}"
```

## Development

```bash
go build -ldflags="-s -w" -o mcp-guard .
go test -count=1 -timeout 30s ./...
go vet ./...
```

---

## 💾 Session Memory Ledger

### [2026-07-18 00:30] — Full Standalone TUI Dashboard Redesign
- **State**: Success — 7 suites pass, pushed to GitHub
- **Core Redesign**: Complete rewrite of `internal/tui/` — from basic stats to full analytics dashboard
- **Layout**: Responsive — side-by-side on wide (>80), stacked on narrow. Lipgloss panels + ANSI render hot path.
- **Dashboard Panels**: Header (version/mode/uptime/status), Traffic Summary (6 counters with mini bars), Top Tools (8), Top Identities (7), Activity Sparkline (24 buckets, 2min window), Live Feed (scrollable)
- **Binary**: 8.5MB, `go vet` clean, all 7 test suites passing

### [2026-07-18 00:12] — Core Hardening
- **State**: Success — injection detection, rate limiting, config validation, TCP proxy rewrite
- **New Packages**: `inject/` (10 pattern categories, homoglyph map), `ratelimit/` (token bucket)
- **TCP Rewrite**: `io.Copy` → line scanner + shared context for clean goroutine shutdown
- **Config**: `Validate()` catches 9+ error types

### [2026-07-17 23:47] — Integration Test + TUI
- **Integration Test**: End-to-end TCP — testmcp server → proxy → client, validates all 3 decision paths + HMAC chain
- **Arch Fix**: `ExtractToolName()` for proper tools/call RBAC
- **TUI**: First version with live audit log tailing

### [2026-07-17 22:35] — Initial TUI Dashboard
- **State**: Success — `mcp-guard top` with Bubble Tea, pause/resume, colored feed

### [2026-07-17 22:30] — GitHub Published
- **Repo**: https://github.com/ravikumarve/ControlPlane

### [2026-07-17 22:15] — Unit Tests Complete
- **30 tests** across 4 packages — policy, audit, config, proxy

### [2026-07-17 22:00] — Go Scaffold + PRD + Architecture
- **22 files**, fully compilable, 7.1MB binary
- **CLI**: init, serve, status, logs, policy, pin, approve

### [2026-07-17] — Idea Validation
- **Concept**: MCP Guard — lightweight Go-based security sidecar
- **Market**: 8.2/10 validation score, unique SMB/solo-dev niche
