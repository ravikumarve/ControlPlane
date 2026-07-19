# ControlPlane AI — Company Overview

> Security infrastructure for the AI agent ecosystem.

---

## Mission

Secure the AI agent ecosystem with zero-trust, agent-aware security layers that protect production systems from unintended or malicious agent actions.

## Vision

Every AI agent should be safe to deploy in production — whether built by a solo developer in a weekend or deployed by a Fortune 500 platform team.

## Values

- **Simplicity** — Single-binary deployments, minimal configuration, sane defaults. Security should not require a PhD in infrastructure.
- **Transparency** — Open-core model, auditable source code, clear documentation. No black boxes.
- **Security-First** — Every feature is evaluated through the lens of threat modeling. We ship protections, not just features.
- **CPU-Conscious** — Designed to run on modest hardware. Our own development happens on a Latitude 3460. Bloat is a security risk.
- **Open-Core** — The core proxy is MIT-licensed. Community and enterprise editions add progressive capabilities without trapping users.

---

## Problem

AI agents now have direct access to production databases, APIs, and operational tools through protocols like MCP (Model Context Protocol). These protocols provide structured access to tools — but they have no built-in security model.

Existing security solutions fall short:

- **API gateways** validate static HTTP routes and cannot inspect dynamic agent intent or tool parameters.
- **Network firewalls** operate at the transport layer and cannot distinguish between benign and malicious tool calls.
- **IAM systems** provide coarse access control but cannot enforce context-aware policies (e.g., "only allow read queries during business hours").
- **Most teams rely on agent honesty** — hoping the LLM will not call dangerous tools with dangerous parameters. This is not a security strategy.

The result: organizations are hesitant to deploy autonomous agents in production, or worse, they deploy them without protection and accept the risk.

---

## Solution

ControlPlane AI provides **agent-aware security proxies** that understand tool semantics, not just transport protocols. Our products sit between agent orchestration engines and the tools/servers they call, inspecting every invocation before it reaches the target.

Core capabilities:

- **Parameter-aware policy enforcement** — Evaluate tool call arguments against policies before execution.
- **Pre-execution circuit breaker** — Block or flag dangerous invocations in real time.
- **Human-in-the-loop (HITL) approval** — Route high-risk operations to a human reviewer via webhook.
- **Immutable audit trails** — HMAC-chained, tamper-evident logs for compliance and post-incident analysis.
- **Agent identity and RBAC** — Restrict agent identities to granular tool scopes with attestable credentials.

---

## Products

### Current: MCP Security & Gateway Proxy

A zero-trust sidecar for MCP agents. Deployed alongside agent runtimes, it intercepts all tool invocations and enforces policy before execution reaches the target server. Single Go binary — no Kubernetes or sidecar injector required.

| Capability | Community | Pro | Enterprise |
|---|---|---|---|
| RBAC for agents | Up to 3 agents | Up to 25 agents | Unlimited |
| Pre-execution circuit breaker | Static rules | Regex + pattern matching | Full semantic analysis |
| HITL approval webhooks | — | 1 webhook | Unlimited webhooks |
| Immutable audit trail | JSONL stdout | JSONL file rotation | Forward to SIEM |
| Policy Management UI | — | — | Included |

### Planned: Policy Management Console

Web-based UI for authoring, testing, and deploying policies. Includes policy versioning, dry-run mode, and audit log viewer.

- **Target**: Q4 2026
- **Status**: Design phase

### Planned: SIEM Integration

Forward audit logs and policy events to Splunk, Datadog, Grafana Loki, or any Elastic-compatible endpoint. Enables integration with existing security operations workflows.

- **Target**: Q1 2027
- **Status**: Research phase

---

## Target Audience

| Persona | Need | Fit |
|---|---|---|
| **Solo developer** | Quick setup, minimal resource usage, protect a personal project | Community edition — runs on a $5 VPS |
| **SMB DevOps** | Agent guardrails without hiring a security team | Pro edition — self-serve, email support |
| **Enterprise platform team** | Compliance, audit, RBAC at scale | Enterprise edition — SSO, SLA, dedicated support |

---

## Business Model

**Open-core licensing** — The core proxy is MIT-licensed and freely available. Paid editions add progressive capabilities for teams that need them.

| Edition | License | Price | Includes |
|---|---|---|---|
| Community | MIT | Free | Core proxy, basic RBAC, static rules |
| Pro | Commercial | $199/license | Enhanced RBAC, HITL, file audit, regex rules |
| Enterprise | Commercial | $999/year | Unlimited agents, SIEM, Policy UI, SSO, SLA |

All editions include access to the public issue tracker and documentation. Pro and Enterprise include direct email support.

---

## Competitive Landscape

The AI agent security space is emerging. Key differentiators for ControlPlane AI:

| Dimension | ControlPlane AI | API Gateways (Kong, APISIX) | Cloud WAF (Cloudflare, AWS) | Do-It-Yourself |
|---|---|---|---|---|
| Agent-aware policy | Native — inspects tool parameters | None — HTTP route only | None — HTTP only | Requires custom middleware |
| Deployment | Single Go binary | K8s-heavy, DB required | Cloud proxy — adds latency | Custom, no standard |
| CPU footprint | < 20 MB idle, < 5% CPU | 500 MB+ per pod | N/A (cloud) | Depends on implementation |
| HITL approval | Built-in | Requires custom plugin | Not available | Must build from scratch |
| Audit trail | HMAC-chained, tamper-proof | Standard logs | Request logs only | Must implement |
| Open source core | Yes (MIT) | Some (Apache 2.0) | No | Yes, but scattered |

**Key differentiators:**

- **Lightweight single-binary** — No Kubernetes, no database, no sidecar injector. Run it on a Raspberry Pi or a $5 VPS.
- **CPU-conscious design** — Built on the same resource-constrained hardware our target market uses.
- **HITL approval workflows** — Pause dangerous operations and route them to a human reviewer. No other solution offers this out of the box.
- **Agent-aware semantics** — Understands the difference between `read_file("/etc/passwd")` and `read_file("/logs/app.log")` at the policy level.

---

## Roadmap

| Timeframe | Milestone |
|---|---|
| Q3 2026 | MCP Security Gateway v1.0 — core proxy, RBAC, circuit breaker, HITL, audit trail |
| Q4 2026 | Policy Management Console beta — web UI, policy versioning, dry-run |
| Q1 2027 | SIEM Integration — Splunk, Datadog, Loki forwarding |
| Q2 2027 | v2.0 — Policy-as-code (Rego), multi-protocol support (beyond MCP) |

---

## Contact

- **Maintainer**: Ravi Kumar
- **Website**: [controlplaneai.com](https://controlplaneai.com) (coming soon)
- **GitHub**: [github.com/ravikumarve](https://github.com/ravikumarve)
- **Email**: ravi@controlplaneai.com (coming soon)

---

*ControlPlane AI — Security infrastructure for the AI agent ecosystem.*
