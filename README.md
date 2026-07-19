# ControlPlane

**Security infrastructure for the AI agent ecosystem.**

ControlPlane builds zero-trust security tools for organizations deploying autonomous AI agents. Our products sit between agent orchestration engines and the tools/servers they call — enforcing policy, preventing injection attacks, and providing tamper-evident audit trails.

---

## Products

### mcp-guard — MCP Security Gateway

A zero-trust sidecar for MCP agents — RBAC, injection detection, rate limiting, circuit breaker, HITL approval, and HMAC-chained audit. Single Go binary, no Kubernetes, no SaaS.

```bash
mcp-guard serve                    # Start proxy daemon
mcp-guard top                      # Live TUI dashboard
mcp-guard logs --tail              # View audit log
mcp-guard init --template github   # Generate config from template
```

### Roadmap

| Product | Status | Target |
|---------|--------|--------|
| mcp-guard V1 | ✅ MVP | Q3 2026 |
| Policy Management Console | 🟡 Design | Q4 2026 |
| SIEM Integration | 🔵 Research | Q1 2027 |

---

## Repository Structure

```
├── mcp-guard/           ← Go backend (the actual product)
│   ├── cmd/             ← CLI: serve, top, logs, policy, init, pin, approve
│   ├── internal/        ← Packages: proxy, policy, inject, audit, hitl, ratelimit, ...
│   └── main.go
├── app/                 ← Next.js landing page (controlplane.ai)
├── components/          ← React components (Button, Card, Hero, Features, ...)
├── docs/                ← Company foundation documents
├── lib/                 ← Frontend utilities
├── tailwind.config.ts   ← Brand colors and design tokens
└── controlplane_ai_landing_page.html  ← Design mockup
```

### Stack

| Layer | Technology |
|-------|-----------|
| Security Gateway | Go 1.24 (mcp-guard) |
| Landing Page | Next.js 14 + Tailwind CSS (static export) |
| Policy Engine | YAML-based RBAC |
| Audit | HMAC-chained JSONL |
| TUI | Bubble Tea |

---

## Development

```bash
# Build the Go gateway
cd mcp-guard && go build -ldflags="-s -w" -o mcp-guard .

# Run tests (68+ tests across 8 packages)
cd mcp-guard && go test -count=1 -timeout 30s ./...

# Landing page dev server (from root)
npm run dev

# Landing page build
npm run build
```

No database required. No environment variables needed for development.

---

## Documentation

| Document | Purpose |
|----------|---------|
| [COMPANY.md](./docs/COMPANY.md) | Mission, vision, values |
| [PRD.md](./docs/PRD.md) | Product requirements |
| [ARCHITECTURE.md](./docs/ARCHITECTURE.md) | System architecture |
| [FRONTEND-PLAN.md](./docs/FRONTEND-PLAN.md) | Frontend architecture |
| [PRODUCT-ROADMAP.md](./docs/PRODUCT-ROADMAP.md) | Milestones |
| [VALIDATION_REPORT.md](./docs/VALIDATION_REPORT.md) | Market research |
| [E2E-TESTING-PIPELINE.md](./docs/E2E-TESTING-PIPELINE.md) | Test strategy |
| [GOVERNANCE.md](./docs/GOVERNANCE.md) | Decision-making, roles |
| [CONTRIBUTING.md](./docs/CONTRIBUTING.md) | How to contribute |
| [CODE_OF_CONDUCT.md](./docs/CODE_OF_CONDUCT.md) | Community standards |
| [SECURITY.md](./docs/SECURITY.md) | Vulnerability disclosure |
| [BRAND.md](./docs/BRAND.md) | Logo, colors, voice |

---

## License

MIT — see [LICENSE](./LICENSE)
