# ControlPlane AI — Project Context

**ControlPlane AI** builds security infrastructure for the AI agent ecosystem.  
Our first product is an MCP Security & Gateway Proxy — a zero-trust sidecar for AI agents.

This repo holds the company foundation — mission, strategy, architecture, and product ideation.

## 💾 Session Memory Ledger

### [2026-07-19 17:30] — Split MCP Guard into Standalone Project
- **State**: Success — mcp-guard Go implementation moved to its own repo (`ravikumarve/mcp-guard`)
- **Actions**: Cleaned ControlPlane AI repo of all Go code. Now holds only company docs.
- **Next Turn Directive**: Build ControlPlane AI company foundation documents

### [2026-07-19 19:30] — Company Foundation Documents Complete (12/12)
- **State**: Success — All 12 company foundation documents written and branded
- **MCP Data Used**: grep_app (brand pattern research), code_tree (project structure validation)
- **Agency Agents Deployed**: @docs (2 parallel tasks — COMPANY, GOVERNANCE, CODE_OF_CONDUCT, SECURITY, BRAND, E2E-TESTING-PIPELINE)
- **Actions Performed**:
  - Rewrote `README.md` as ControlPlane AI company overview (removed mcp-guard product focus)
  - Updated `PRD.md`, `ARCHITECTURE.md` titles to use ControlPlane AI branding
  - Created 10 new documents: COMPANY.md, GOVERNANCE.md, CODE_OF_CONDUCT.md, CONTRIBUTING.md, PRODUCT-ROADMAP.md, E2E-TESTING-PIPELINE.md, SECURITY.md, BRAND.md (SECURITY and BRAND had partial content from prior session, full docs completed now)
  - One docs agent task silently failed (CONTRIBUTING.md, PRODUCT-ROADMAP.md) — written directly as fallback
- **Architectural Decision**: ControlPlane AI is the parent company; mcp-guard is the product (CLI/repo name). All docs use this naming convention per BRAND.md guidelines. Repo now holds only company docs — Go code lives in separate `ravikumarve/mcp-guard` repo.
- **Docs Created**: README, COMPANY, GOVERNANCE, CODE_OF_CONDUCT, CONTRIBUTING, PRODUCT-ROADMAP, SECURITY, BRAND, E2E-TESTING-PIPELINE, PRD, ARCHITECTURE, VALIDATION_REPORT,  idea.md, LICENSE = **14 files total**
- **Next Turn Directive**: Start mcp-guard v1.0 release prep, or begin work on another project (FounderLens, StudyAI, etc.)

### [2026-07-19 19:00] — Security & Brand Documentation Complete
- **State**: Success — SECURITY.md (139 lines) and BRAND.md (199 lines) created
- **MCP Data Used**: code_tree (project structure check), read (existing docs for style reference)
- **Actions Performed**:
  - Created `SECURITY.md` — vulnerability disclosure policy with reporting flow, PGP key section, scope/out-of-scope definitions, disclosure timeline table, Hall of Fame, and security practices (HMAC chain audit, default-deny policy, injection detection, minimal binary)
  - Created `BRAND.md` — brand guidelines with naming rules, tagline, color palette (primary/semantic/background), typography stack and scale, voice & tone principles with anti-patterns and examples, logo usage, product naming conventions, attribution rules
  - Both files validated for structure: proper H1/H2 hierarchy, code blocks with language tags, consistent formatting
- **Key Outputs**: SECURITY.md references mcp-guard security practices; BRAND.md codifies product naming (MCP Security Gateway / mcp-guard); both are referenced in README.md documentation table
- **Next Turn Directive**: Continue with COMPANY.md or PRODUCT-ROADMAP.md as planned in the project structure

### [2026-07-17] — MCP Security Gateway Idea Validation
- **State**: Success — Full validation report delivered
- **MCP Data Used**: websearch_tavily_search (3 rounds — MCP market, competitors, gateway comparison), websearch_tavily_search (competitor deep-dive), websearch_tavily_search (market sizing)
- **Actions Performed**:
  - Analyzed `idea.md` — MCP Security & Gateway Proxy concept
  - Conducted comprehensive market research across 15+ sources
  - Mapped competitive landscape (15+ open-source + commercial vendors)
  - Identified underserved SMB/self-hosted niche
  - Built SWOT analysis and go-to-market strategy
  - Wrote `VALIDATION_REPORT.md` with full findings and pivot recommendation
- **Key Insight**: CPU-only constraint (Latitude 3460) is a competitive advantage — building for the same constraints as target market
- **Next Turn Directive**: Build company foundation documents and roadmap

## Project Structure
```
ControlPlane AI/
├── README.md              ← Company overview (root)
├── LICENSE                ← MIT (root)
├── docs/
│   ├── AGENTS.md          ← This file — session context
│   ├── COMPANY.md         ← Mission, vision, values, portfolio
│   ├── PRODUCT-ROADMAP.md ← Product direction and milestones
│   ├── GOVERNANCE.md      ← Decision-making, roles
│   ├── CONTRIBUTING.md    ← How to contribute
│   ├── CODE_OF_CONDUCT.md ← Community standards
│   ├── SECURITY.md        ← Vulnerability disclosure
│   ├── BRAND.md           ← Logo, colors, voice
│   ├── E2E-TESTING-PIPELINE.md ← Test strategy
│   ├── PRD.md             ← Product requirements (MCP Security Gateway)
│   ├── ARCHITECTURE.md    ← System architecture
│   ├── VALIDATION_REPORT.md ← Market research
│   └── idea.md            ← Original concept
```
