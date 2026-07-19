# MCP Guard

**Lightweight security sidecar for MCP agents.**  
Single binary. No Kubernetes. No SaaS. Deploy in 5 minutes.

MCP Guard sits between AI agents (Claude Code, Cursor, Copilot) and MCP servers, enforcing tool-level access control, **parameter-level policy matching**, injection detection, rate limiting, schema pinning, alert webhooks, and tamper-evident audit logging.

```bash
mcp-guard serve                    # Start proxy daemon
mcp-guard top                      # Live TUI dashboard
mcp-guard logs --tail              # View audit log
mcp-guard init --template github   # Generate config from template
```

---

## Features

| Layer | Capability |
|-------|-----------|
| **Access Control** | Tool-level + **parameter-level** RBAC via YAML policies (allow/block/human-in-the-loop) |
| **Injection Detection** | 10 prompt injection pattern categories, homoglyph detection, JSON depth bombs |
| **Rate Limiting** | Token bucket per identity, configurable `100/m` / `10/s` / `1000/h` |
| **Parameter Blocking** | Glob-match on tool parameters — block SQL by query content, filesystem access by path, git writes by branch |
| **Alert Webhooks** | Slack, Discord, or generic HTTP webhooks on block/rate-limit/injection events |
| **Policy Templates** | Quick-start configs for GitHub, PostgreSQL, Slack, Filesystem, and generic MCP servers |
| **Audit Log** | HMAC-SHA256 chained JSONL — tamper-evident, verifiable |
| **Schema Pinning** | SHA-256 drift detection against supply-chain poisoning |
| **HITL** | Human-in-the-loop approval via webhook (Slack, Teams, custom) |
| **Live Dashboard** | Full TUI — 2-column grid with traffic stats, per-tool breakdown, activity sparkline, scrollable feed |
| **Proxy Modes** | `stdio` (CLI agents) + `tcp` (remote servers) |

## Quick Start

```bash
# Build
go build -ldflags="-s -w" -o mcp-guard .

# Generate config from a template
./mcp-guard init --template github

# Edit config to match your needs
vim mcp-guard.yaml

# Run in stdio mode (e.g., with Claude Code)
claude --mcp-server "./mcp-guard serve"

# Or run in TCP mode
./mcp-guard serve

# Open the live dashboard in another terminal
./mcp-guard top
```

## CLI Commands

```
  approve     Approve or deny a pending HITL request
  init        Generate a default mcp-guard.yaml config file
              --template <name>   Use a template (github, postgres, slack, filesystem, generic)
  logs        View or verify the audit log
  pin         Manage schema pins
  policy      List or test security policies
  serve       Start the MCP Guard proxy daemon
  status      Show daemon status
  top         Live TUI dashboard
```

## Configuration

MCP Guard uses a single YAML config file (default: `./mcp-guard.yaml` or `/etc/mcp-guard/config.yaml`):

```yaml
version: "1"
proxy:
  mode: tcp           # stdio | tcp
  listen: ":8443"
  upstream: "localhost:3000"

policies:
  # Tool-level allow for safe operations
  - name: allow-reads
    action: allow
    match:
      identity: "*"
      tools: ["read_*", "get_*", "list_*"]

  # Parameter-level blocking — block destructive SQL by content
  - name: block-destructive-sql
    action: block
    alert: true
    match:
      identity: "*"
      tools: ["query", "run_sql", "execute"]
      params:
        sql: ["DROP *", "TRUNCATE *", "ALTER *", "DELETE *"]

  # Path-based filesystem blocking
  - name: block-system-paths
    action: block
    match:
      tools: ["read_*", "write_*"]
      params:
        path: ["/etc/**", "/var/**", "/proc/**", "/sys/**"]

  # Human-in-the-loop for sensitive operations
  - name: risky-operations
    action: hitl
    match:
      tools: ["execute_payout", "send_email_batch"]

ratelimit:
  default: 60/min

alert:
  webhook_url: "${SLACK_WEBHOOK_URL}"
  channel: slack          # slack | discord | generic
  on_block: true
  on_rate_limit: true

audit:
  path: /var/log/mcp-guard/audit.jsonl
  hmac_key: "${MCP_GUARD_HMAC_KEY}"
  rotation: 100MB

schema_pinning:
  enabled: true
  mode: warn              # warn | block
  store: ".mcp-guard/pins.json"
```

### Policy Templates

Generate a complete, production-ready config in one command:

```bash
# GitHub MCP — blocks destructive writes on main/production branches
mcp-guard init --template github

# PostgreSQL — allows SELECT only, blocks DDL/DML by SQL content
mcp-guard init --template postgres

# Slack — blocks admin actions (kick, ban, delete)
mcp-guard init --template slack

# Filesystem — blocks access to system paths (/etc, /var, /proc)
mcp-guard init --template filesystem

# Generic — sensible defaults for any MCP server
mcp-guard init --template generic
```

Each template includes pre-configured policies, rate limits, alerts, audit, and schema pinning.

## Live Dashboard

```
┌──────────────────────────────────┬──────────────────────────┐
│  ▲ MCP GUARD  v0.1.0  [tcp]  ● LIVE  Uptime: 12m 34s      │
├──────────────────────────────────┼──────────────────────────┤
│  Traffic Summary                 │  Top Tools               │
│  Total Calls           1,234     │  read_db       892 ████  │
│  Allowed               1,112  37 │  list_users    289 ██    │
│  Blocked                 122     │  delete_db      45 █     │
│  HITL Pending             12     │  execute_pay     8 █     │
│  Rate Limited              5     │                          │
│  Injection Blocks          2     │  Top Identities          │
│                                  │  admin         987 ███   │
│  Activity (2min window)          │  ci/pipeline   245 █     │
│  ▁▂▃▄▅▆▇█▇▆▅▄▃▂▁               │  bot/scraper    12 █     │
│  14:23:31       14:25:31        │                          │
├──────────────────────────────────┴──────────────────────────┤
│  Live Feed                              q:quit  p:pause     │
│  14:23:01  ALLOW  admin     read_db                         │
│  14:23:05  BLOCK  admin     delete_db  [blocked by policy]  │
│  14:24:12  HITL   bot       execute_payout                  │
└─────────────────────────────────────────────────────────────┘
```

Controls: `q` quit · `p` pause · `↑↓` scroll · `--audit-path` custom location

## Architecture

```
┌──────────┐    JSON-RPC    ┌────────────┐    JSON-RPC    ┌──────────┐
│   AI     │ ──────────────▶│  MCP       │ ──────────────▶│   MCP    │
│  Agent   │ ◀──────────────│  Guard     │ ◀──────────────│  Server  │
└──────────┘                │  Proxy     │                └──────────┘
                            │            │
                            │  ┌──────┐  │
                            │  │Policy│  │  tool + parameter RBAC
                            │  ├──────┤  │
                            │  │Inject│  │  prompt injection scan
                            │  ├──────┤  │
                            │  │Rate  │  │  token bucket /identity
                            │  ├──────┤  │
                            │  │Audit │  │  HMAC-chained JSONL
                            │  ├──────┤  │
                            │  │Alert │  │  Slack/Discord webhooks
                            │  └──────┘  │
                            └────────────┘
                                  │
                            ┌─────▼─────┐
                            │  TUI Dash │  mcp-guard top
                            └───────────┘
```

## Package Layout

```
cmd/              CLI commands (9 commands)
internal/
  alert/          Webhook dispatcher — Slack, Discord, generic HTTP
  audit/          HMAC-chained JSONL audit logger
  config/         YAML config loader + validator (9+ error types)
  hitl/           Human-in-the-loop webhook dispatcher
  inject/         Prompt injection detection engine (10 pattern categories)
  policy/         First-match-wins RBAC engine (tool + parameter matching)
  proxy/          JSON-RPC proxy (stdio + TCP scanner-based)
  ratelimit/      Token bucket rate limiter (parse "100/m", "10/s", "1000/h")
  schema/         Schema pinning (SHA-256 drift detection)
  templates/      Embedded policy YAMLs — github, postgres, slack, filesystem, generic
  testmcp/        Test MCP server + end-to-end TCP integration test
  tui/            Bubble Tea live dashboard (2-column grid, ANSI hot path)
  version/        Build version and metadata
```

## Build & Test

```bash
# Build (8.5MB binary, no CGO)
go build -ldflags="-s -w" -o mcp-guard .

# Run all tests (10 suites, 68+ tests)
go test -count=1 -timeout 30s ./...

# End-to-end integration test (allow/block/HITL + HMAC chain)
go test -count=1 -v ./internal/testmcp/...

# Verify
go vet ./...
```

## Compliance

MCP Guard's tamper-evident audit chains and injection detection help meet **EU AI Act** requirements (enforcement: August 2, 2026 — penalties up to 7% global revenue). The HMAC-chained JSONL log provides verifiable evidence of all agent actions for audit trails.

## License

MIT — see [LICENSE](LICENSE)
