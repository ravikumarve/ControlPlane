# ControlPlane — MCP Security Infrastructure

**ControlPlane** builds lightweight, self-hosted security tools for the Model Context Protocol (MCP) ecosystem. Single binary. No Kubernetes. No SaaS. 5-minute deploy.

## Projects

### [MCP Guard](./mcp-guard/) — Security Sidecar for MCP Agents

MCP Guard sits between AI agents (Claude Code, Cursor, Copilot) and MCP servers, enforcing tool-level access control, schema pinning, and tamper-evident audit logging.

```bash
# Install
go install github.com/ravikumarve/ControlPlane/mcp-guard@latest

# Initialize
mcp-guard init

# Run
mcp-guard serve
```

**Features**:
- 🔒 Tool-level RBAC via YAML policies
- 📝 Tamper-evident audit log (HMAC-SHA256 chain)
- 🔐 Schema pinning against supply-chain poisoning
- 👋 Human-in-the-loop approval workflows (Slack/Webhook)
- 🚦 Rate limiting per-identity, per-tool
- ⚡ stdio + TCP proxy modes
- 🖥️ CLI-only — single 7MB binary, no deps

**Status**: MVP — pass all tests, ready for production pilots.

## Documentation

| Doc | Description |
|-----|-------------|
| [PRD.md](./PRD.md) | Product requirements, MVP scope, pricing |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Proxy flow, data structures, component breakdown |
| [VALIDATION_REPORT.md](./VALIDATION_REPORT.md) | Market research, competitor analysis, go-to-market |

## Quick Start

```bash
# Generate config
mcp-guard init

# Edit the config
vim mcp-guard.yaml

# Run in stdio mode (e.g., with Claude Code)
claude --mcp-server "mcp-guard serve"

# List policies
mcp-guard policy list

# Test a policy
mcp-guard policy test read_database --identity my-agent

# View audit log
mcp-guard logs --tail

# Verify audit integrity
mcp-guard logs --verify
```

## License

MIT — see [LICENSE](./LICENSE)
