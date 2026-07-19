# ControlPlane — MCP Security Infrastructure

**ControlPlane** builds lightweight, self-hosted security tools for the Model Context Protocol (MCP) ecosystem. Single binary. No Kubernetes. No SaaS. Deploy in 5 minutes.

## [MCP Guard](./mcp-guard/) — Security Sidecar for MCP Agents

MCP Guard sits between AI agents (Claude Code, Cursor, Copilot) and MCP servers, enforcing tool-level access control, injection detection, rate limiting, and tamper-evident audit logging.

```bash
cd mcp-guard/
go build -o mcp-guard .
./mcp-guard init
./mcp-guard serve
./mcp-guard top          # Live TUI dashboard
```

### Features
- 🔒 Tool-level RBAC via YAML policies (allow / block / human-in-the-loop)
- 🛡️ Prompt injection detection (10 pattern categories, homoglyph, depth bomb)
- 🚦 Rate limiting — token bucket per identity
- 📝 Tamper-evident audit log (HMAC-SHA256 chained JSONL)
- 🔐 Schema pinning against supply-chain poisoning
- 👋 Human-in-the-loop approval workflows (Slack/Webhook)
- ⚡ stdio + TCP proxy modes
- 🖥️ Live TUI dashboard — stats, per-tool breakdown, activity sparkline
- 💾 Single 8.5MB binary, zero external deps at runtime

### Quick Start

```bash
cd mcp-guard/
go build -ldflags="-s -w" -o mcp-guard .
./mcp-guard init                     # Generate mcp-guard.yaml
vim mcp-guard.yaml                   # Edit policies
./mcp-guard serve -v                 # Start proxy

# In another terminal:
./mcp-guard top --audit-path /tmp/mcp-guard/audit.jsonl
```

### Documentation

| Doc | Description |
|-----|-------------|
| [mcp-guard/README.md](./mcp-guard/README.md) | Full project README (build, config, usage, architecture) |
| [PRD.md](./PRD.md) | Product requirements, MVP scope, pricing |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Proxy flow, data structures, component breakdown |
| [VALIDATION_REPORT.md](./VALIDATION_REPORT.md) | Market research, competitor analysis, go-to-market |

### Status

**MVP** — 7 test suites passing, 58+ tests including end-to-end integration. Ready for production pilots.

## License

MIT — see [LICENSE](./LICENSE)
