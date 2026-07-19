# ControlPlane AI — Product Roadmap

## Vision

ControlPlane AI builds the security layer for the AI agent ecosystem — starting with MCP (Model Context Protocol) and expanding to all agent communication protocols.

## North Star

Every AI agent deployed in production runs behind a ControlPlane security layer.

---

## Now — Active (Q3 2026)

### MCP Security Gateway v1.0 (mcp-guard)

**Status**: MVP complete, targeting stable release

**Capabilities**:
- Zero-trust sidecar proxy for MCP agents
- Identity-based RBAC with granular tool scoping
- Pre-execution policy evaluation (ALLOW / BLOCK / HITL)
- Human-in-the-Loop (HITL) approval via webhook
- Injection detection for tool parameters
- Immutable HMAC-chained audit log
- Rate limiting per identity
- Circuit breaker for upstream failures
- YAML-based configuration
- Single Go binary deployment (no runtime dependencies)

**Deliverables**:
- [x] Core proxy engine
- [x] Policy evaluation (allow/block/hitl)
- [x] HMAC audit chain
- [x] Rate limiting
- [ ] V1.0 release with documentation
- [ ] Quickstart guides and examples
- [ ] Community edition (MIT license)

---

## Next — Q4 2026

### Policy Management Console

- Web UI for visual policy editing
- Policy templates for common use cases
- Real-time policy testing and preview
- Audit log viewer with search and filtering

### Injection Detection Hardening

- Schema pinning — detect schema drift in real-time
- Anti-poisoning support for tool descriptions
- Pattern library for known injection vectors

### Observability

- Prometheus metrics endpoint
- Structured JSON logging
- Integration guides for Datadog, Grafana

### Pro Tier Launch

- **$199/license** — Premium features for professional teams
- HITL approval management dashboard
- Advanced policy rules (time-based, context-based)
- Priority support

---

## Later — Q1 2027

### Multi-Protocol Support

- Extend beyond MCP to other agent protocols
- Protocol adapter framework
- Consistent policy model across protocols

### Enterprise Features

- SSO / SAML authentication for policy management
- Audit dashboard with visual query builder
- Role-based access for policy management
- Custom audit retention policies

### Enterprise Tier Launch

- **$999/year** — Enterprise-grade features
- Audit dashboard with compliance reporting
- Enterprise SSO integration
- SLAs and dedicated support

### Kubernetes Operator

- Auto-injection of security proxy into agent pods
- Kubernetes-native policy management via CRDs
- Service mesh integration for sidecar injection

---

## Explicit Non-Goals

These areas are explicitly outside ControlPlane AI's product scope:

- Building or training LLMs
- Building vector databases or embedding stores
- Building a competing MCP server registry
- Building agent orchestration frameworks
- Building agent-to-agent communication protocols

---

## Milestone Summary

| Product | Q3 2026 | Q4 2026 | Q1 2027 |
|---------|---------|---------|---------|
| MCP Security Gateway | V1.0 stable | Maintenance | Maintenance |
| Policy Management Console | — | Beta | GA |
| SIEM Integration | — | Beta | GA |
| Multi-Protocol | — | Research | Beta |
| Kubernetes Operator | — | — | Beta |
| Community (MIT) | Available | Available | Available |
| Pro ($199) | — | Available | Available |
| Enterprise ($999/yr) | — | — | Available |

## How to Influence the Roadmap

This roadmap is driven by community and customer feedback. To influence priorities:

1. **Vote on GitHub Issues** — Add reactions to existing issues
2. **Submit feature requests** — Use the feature request template in [CONTRIBUTING.md](./CONTRIBUTING.md)
3. **Sponsor development** — Enterprise tier customers get roadmap prioritization
