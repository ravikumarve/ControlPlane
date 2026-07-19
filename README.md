# ControlPlane AI

**Security infrastructure for the AI agent ecosystem.**

ControlPlane AI builds zero-trust security tools for organizations deploying autonomous AI agents. Our products sit between agent orchestration engines and the tools/servers they call — enforcing policy, preventing injection attacks, and providing tamper-evident audit trails.

## Why ControlPlane AI

AI agents now have direct access to production databases, APIs, and operational tools through protocols like MCP (Model Context Protocol). But these protocols have no built-in security. Existing API gateways only validate static HTTP routes and cannot inspect dynamic agent intent.

ControlPlane AI solves this with **agent-aware security layers** that understand tool semantics, not just transport protocols.

## Products

### MCP Security & Gateway Proxy
A zero-trust sidecar for MCP agents — RBAC, injection detection, circuit breaker, and HITL approval. Currently in MVP.

> **Repository**: [ravikumarve/mcp-guard](https://github.com/ravikumarve/mcp-guard)

**Capabilities:**
- **RBAC for Agents** — restrict agent identities to granular tool scopes
- **Pre-Execution Circuit Breaker** — evaluate tool parameters before execution
- **Human-in-the-Loop (HITL)** — pause and approve high-risk operations via webhook
- **Immutable Audit Trail** — HMAC-chained JSONL for compliance

### Roadmap
| Product | Status | Target |
|---------|--------|--------|
| MCP Security Gateway (mcp-guard) | 🟢 MVP | V1 Release Q3 2026 |
| Policy Management UI | 🟡 Design | Q4 2026 |
| SIEM Integration | 🔵 Research | Q1 2027 |

## Documentation

| Document | Purpose |
|----------|---------|
| [COMPANY.md](./docs/COMPANY.md) | Mission, vision, values, portfolio |
| [PRD.md](./docs/PRD.md) | Product requirements for MCP Security Gateway |
| [ARCHITECTURE.md](./docs/ARCHITECTURE.md) | System architecture and data flow |
| [PRODUCT-ROADMAP.md](./docs/PRODUCT-ROADMAP.md) | Phased delivery milestones |
| [VALIDATION_REPORT.md](./docs/VALIDATION_REPORT.md) | Market research and competitive analysis |
| [GOVERNANCE.md](./docs/GOVERNANCE.md) | Decision-making and maintainer roles |
| [CONTRIBUTING.md](./docs/CONTRIBUTING.md) | How to contribute |
| [CODE_OF_CONDUCT.md](./docs/CODE_OF_CONDUCT.md) | Community standards |
| [SECURITY.md](./docs/SECURITY.md) | Vulnerability disclosure |
| [BRAND.md](./docs/BRAND.md) | Logo, colors, voice, messaging |
| [E2E-TESTING-PIPELINE.md](./docs/E2E-TESTING-PIPELINE.md) | End-to-end testing strategy, CI pipeline, and quality gates |

## License

MIT — see [LICENSE](./LICENSE)
