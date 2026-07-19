# Security Policy

ControlPlane AI builds security infrastructure for the AI agent ecosystem. We take the security of our products seriously and appreciate the community's help in disclosing vulnerabilities responsibly.

## Reporting a Vulnerability

If you believe you have found a security vulnerability in any ControlPlane AI project or product, please report it to us immediately.

**Email**: security@controlplane.dev

We aim to acknowledge receipt of your report within **48 hours**. Please include as much detail as possible to help us triage and reproduce the issue:

- Product, version, and environment details
- Steps to reproduce the vulnerability
- Proof-of-concept code (if applicable)
- Potential impact assessment

Encrypted reports are preferred. Our PGP key is available below.

## Scope

This security policy covers all ControlPlane AI open-source projects and products, including but not limited to:

- **mcp-guard** — MCP Security & Gateway Proxy
- ControlPlane AI documentation and websites
- All officially maintained repositories under the ControlPlane AI GitHub organization

Third-party dependencies integrated into our projects are considered **in scope only if the vulnerability arises from our usage or configuration** of the dependency, not from the dependency itself.

## Out of Scope

The following are considered out of scope for our security disclosure program:

- Vulnerabilities in third-party dependencies that are used as-is (report to the upstream maintainer)
- Theoretical attacks without a working proof of concept
- Social engineering attacks against ControlPlane AI maintainers
- Physical attacks against infrastructure we do not control
- Issues requiring unlikely or impractical attacker prerequisites
- Self-XSS or similar attacks that cannot be chained with other vulnerabilities
- Missing HTTP security headers that do not result in a demonstrable vulnerability
- Rate limiting bypasses without demonstrated impact
- Automated tool scan results without manual verification
- Vulnerabilities in pre-release or alpha software that are already documented as known issues

If you are unsure whether a vulnerability is in scope, please report it anyway. We prefer over-reporting to missing a critical issue.

## Disclosure Policy

We follow a coordinated disclosure process to protect users while acknowledging researchers' contributions.

| Stage | Timeline | Description |
|-------|----------|-------------|
| **Acknowledgment** | Within 48 hours of report | We confirm receipt and provide a tracking identifier |
| **Triage** | Within 5 business days | We verify the vulnerability, assess severity using CVSS 4.0, and determine impact |
| **Fix Development** | 30 days (Critical/High) / 60 days (Medium/Low) | We develop, test, and prepare a patch. Timeline may be extended with researcher consent for complex issues |
| **Release & Disclosure** | Coordinated date | We publish a security advisory, release the fix, and credit the researcher (with their consent) |

We will keep you informed throughout the process and will not publish details before a fix is available.

## PGP Key

All security reports should be encrypted using our PGP key.

```
Fingerprint:  (placeholder — generate and publish actual key)
Key ID:       (placeholder)
Download:     https://controlplane.dev/.well-known/pgp-key.asc
```

To fetch the key:

```shell
curl https://controlplane.dev/.well-known/pgp-key.asc | gpg --import
```

## Security Advisories

We publish all confirmed security advisories through **GitHub Security Advisories** for each affected repository. Each advisory includes:

- CVE identifier (when assigned)
- Severity rating (CVSS 4.0)
- Affected versions and fixed versions
- Description and impact
- Mitigation guidance
- Researcher credit

You can watch individual repositories or monitor the ControlPlane AI organization on GitHub for advisory notifications.

## Hall of Fame

We maintain a Hall of Fame to recognize researchers who have helped us improve our security. If your verified disclosure leads to a code change, you may choose to be credited in:

- The published security advisory
- Our SECURITY.md Hall of Fame section
- A dedicated acknowledgment in the affected project's release notes

Credit is entirely optional. Some researchers prefer to remain anonymous, and we respect that choice.

### Acknowledgments

_No researchers have been credited yet. This section will be updated as disclosures are received._

## Security Practices

ControlPlane AI products are designed with security as a core requirement, not an afterthought. The following practices are baked into our architecture:

### HMAC-Chain Audit Logs

All audit events are written to an append-only JSONL log where each entry includes a hash of the previous entry, creating an immutable chain. Tampering with any log entry breaks the chain for all subsequent entries, making unauthorized modification detectable. Verification is built into the CLI:

```shell
mcp-guard logs --verify
```

### Default-Deny Policy Engine

The policy engine operates on a **default-deny** principle. No tool, resource, or action is permitted unless explicitly allowed by a policy rule. This eliminates entire classes of authorization bypass vulnerabilities and ensures that misconfiguration results in denied access rather than unintended access.

### Injection Detection at the Proxy Layer

All tool parameters are inspected at the proxy layer before reaching the target server. Detection includes:

- **NFKC normalization** to catch homoglyph and Unicode-based attacks
- **Pattern matching** against known injection signatures (prompt injection, shell injection, SQL injection)
- **Schema validation** to ensure parameters match expected types and constraints
- **Parameter length limits** to prevent buffer-overflow-style attacks

### Minimal Binary with No Runtime Dependencies

The mcp-guard binary is compiled as a single static Go binary with no runtime dependencies. This reduces the attack surface by:

- Eliminating dependency-confusion vulnerabilities
- Removing the need for a runtime environment (no Python, Node.js, or JVM)
- Enabling reproducible builds for verifiable integrity
- Simplifying updates to a single binary replacement

---

_This security policy is maintained by the ControlPlane AI team and will be reviewed quarterly._
