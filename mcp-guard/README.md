# MCP Guard

**Lightweight security sidecar for MCP agents.**  
Single binary. No Kubernetes. No SaaS. Deploy in 5 minutes.

MCP Guard sits between AI agents (Claude Code, Cursor, Copilot) and MCP servers, enforcing tool-level access control, injection detection, rate limiting, schema pinning, and tamper-evident audit logging.

```bash
mcp-guard serve                    # Start proxy daemon
mcp-guard top                      # Live TUI dashboard
mcp-guard logs --tail              # View audit log
```

---

## Features

| Layer | Capability |
|-------|-----------|
| **Access Control** | Tool-level RBAC via YAML policies (allow/block/human-in-the-loop) |
| **Injection Detection** | 10 prompt injection pattern categories, homoglyph detection, JSON depth bombs |
| **Rate Limiting** | Token bucket per identity, configurable `100/m` / `10/s` / `1000/h` |
| **Audit Log** | HMAC-SHA256 chained JSONL — tamper-evident, verifiable |
| **Schema Pinning** | SHA-256 drift detection against supply-chain poisoning |
| **HITL** | Human-in-the-loop approval via webhook (Slack, Teams, custom) |
| **Live Dashboard** | Full TUI — stats, per-tool breakdown, activity sparkline, scrollable log |
| **Proxy Modes** | `stdio` (CLI agents) + `tcp` (remote servers) |

## Quick Start

```bash
# Build
go build -o mcp-guard .

# Generate default config
./mcp-guard init

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
  - name: allow-reads
    action: allow
    match:
      identity: "*"
      tools: ["read_*", "get_*", "list_*"]

  - name: block-admin
    action: block
    match:
      tools: ["delete_*", "exec_*", "drop_*"]

  - name: risky-operations
    action: hitl
    match:
      tools: ["execute_payout", "send_email_batch"]
    alert: true

audit:
  path: /var/log/mcp-guard/audit.jsonl
  hmac_key: "${MCP_GUARD_HMAC_KEY}"
```

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
│   AI     │ ──────────────▶│  MC P      │ ──────────────▶│   MCP    │
│  Agent   │ ◀──────────────│  Guard     │ ◀──────────────│  Server  │
└──────────┘                │  Proxy     │                └──────────┘
                            │            │
                            │  ┌──────┐  │
                            │  │Policy│  │  allow/block/HITL
                            │  ├──────┤  │
                            │  │Inject│  │  prompt injection scan
                            │  ├──────┤  │
                            │  │Rate  │  │  token bucket /identity
                            │  ├──────┤  │
                            │  │Audit │  │  HMAC-chained JSONL
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
  audit/          HMAC-chained JSONL audit logger
  config/         YAML config loader + validator
  hitl/           Human-in-the-loop webhook dispatcher
  inject/         Prompt injection detection engine
  policy/         First-match-wins RBAC engine
  proxy/          JSON-RPC proxy (stdio + TCP)
  ratelimit/      Token bucket rate limiter
  schema/         Schema pinning (SHA-256 drift detection)
  testmcp/        Test MCP server + integration tests
  tui/            Bubble Tea live dashboard
  version/        Build version and metadata
```

## Build & Test

```bash
# Build (8.5MB binary)
go build -ldflags="-s -w" -o mcp-guard .

# Run all tests (7 suites, 58+ tests)
go test -count=1 -timeout 30s ./...

# End-to-end integration test
go test -count=1 -v ./internal/testmcp/...

# Verify
go vet ./...
```

## License

MIT — see [LICENSE](../LICENSE)
