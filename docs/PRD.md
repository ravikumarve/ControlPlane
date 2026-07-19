# ControlPlane AI — MCP Security Gateway: Product Requirements Document (MVP)
**Version**: 1.0 | **Date**: 2026-07-17 | **Status**: Draft

## 1. Problem Statement
MCP (Model Context Protocol) gives AI agents direct access to tools, databases, and APIs — but the protocol has no built-in security. Independent audits score average MCP server security at 34/100. Existing gateway solutions require Kubernetes, SaaS subscriptions, or enterprise budgets. Solo devs and SMBs need something **lightweight, self-hosted, and deployable in 5 minutes**.

## 2. Target Persona
- **Primary**: Solo developer running Claude Code / Cursor with custom MCP servers
- **Secondary**: SMB DevOps team (2-10 people) with self-hosted infra
- **Archetype**: Technical, terminal-native, budget-conscious, values simplicity over features

## 3. MVP Scope (Ship in 2 weeks)

| Priority | Feature | Description |
|----------|---------|-------------|
| **P0** | JSON-RPC Proxy | Intercepts stdio and TCP MCP traffic transparently |
| **P0** | YAML Policy Engine | Tool-level RBAC: `agent X can call tools Y, Z` |
| **P0** | Audit Logging | Flat-file JSONL with HMAC chain (tamper-evident) |
| **P0** | CLI Lifecycle | `init`, `serve`, `status`, `logs`, `stop` |
| **P1** | Schema Pinning | SHA-256 hash tool definitions; alert on drift (anti-poisoning) |
| **P1** | Injection Detection | NFKC normalization + regex pattern scan |
| **P1** | HITL Webhooks | Slack/Discord/email approval for high-risk actions |
| **P2** | Policy `test` | Dry-run mode to validate policy before deploying |
| **P2** | Prometheus Metrics | `/metrics` endpoint for basic observability |

## 4. Non-Goals (Explicitly Out of MVP Scope)
- ❌ OAuth 2.1 / SSO provider (delegate to existing IdP)
- ❌ Kubernetes operator / Helm chart
- ❌ Multi-tenant SaaS control plane
- ❌ Web dashboard / UI
- ❌ TUI monitoring (Phase 2)
- ❌ Custom MCP server registry
- ❌ Database dependency (SQLite optional in Phase 2)

## 5. Success Criteria
- **Time-to-deploy**: < 5 minutes from `go install` to `mcp-guard serve`
- **Performance**: < 100µs added latency per call (Go binary benchmark)
- **Binary size**: < 15MB
- **Audit integrity**: HMAC chain verifiable with `mcp-guard logs --verify`
- **Security score**: Blocks 100% of unauthorized tool calls (by policy)

## 6. User Stories

### Core Flow
```
As a solo dev running Claude Code with 3 MCP servers,
I want to drop mcp-guard in front of them with a YAML policy,
So that my agent can only call the tools I explicitly allow.

As a DevOps engineer,
I want tamper-evident audit logs of every tool call,
So that I can pass a SOC 2 audit without a SIEM.
```

### Approval Flow
```
As a team lead,
I want to approve database write operations via Slack,
So that no agent mutates production data without human review.
```

## 7. Pricing Model (Post-MVP)
| Tier | Price | Features |
|------|-------|----------|
| Community | Free (MIT) | Core RBAC, basic audit, 3 policies |
| Pro | $199/license | HITL, schema pinning, unlimited policies |
| Enterprise | $999/yr | SSO, SIEM forwarding, hash-chain audit, priority support |
