# ControlPlane AI — Project Context

**ControlPlane AI** builds security infrastructure for the AI agent ecosystem.  
Our first product is an MCP Security & Gateway Proxy — a zero-trust sidecar for AI agents.

This repo holds the company foundation — mission, strategy, architecture, and product ideation.

## 💾 Session Memory Ledger

### [2026-07-19 22:45] — Workspace UI Pages (Login, Dashboard, Policies, Audit, Settings)
- **State**: Success — 5 workspace pages built, all compiled
- **Files Created**:
  - `lib/api.ts` — API client with `api.login()`, `api.getStatus()`, `api.listPolicies()`, `api.savePolicies()`, `api.getAudit()`, `api.getConfig()`
  - `app/(workspace)/layout.tsx` — workspace layout with sidebar nav (◆ Dashboard, ⬡ Policies, ☰ Audit, ⚙ Settings), auth guard, logout
  - `app/(workspace)/login/page.tsx` — login form with API key input, error handling, default key hint
  - `app/(workspace)/dashboard/page.tsx` — 6 stat cards grid (total, allowed, blocked, hitl, rate limited, injection blocked), auto-refresh toggle, percentage breakdowns
  - `app/(workspace)/policies/page.tsx` — policy list with inline editing (name, action, identity, tools), add/remove, save/cancel
  - `app/(workspace)/audit/page.tsx` — audit log viewer with decision filter (all/allow/block/pending), HMAC chain display
  - `app/(workspace)/settings/page.tsx` — config viewer (JSON dump), copy to clipboard, quick info cards
- **Files Modified**: `components/shared/Navbar.tsx` — added orange "Workspace" link to landing page nav
- **URL Structure**: `/login`, `/dashboard`, `/policies`, `/audit`, `/settings` (flat, clean URLs via route groups)
- **Auth Flow**: API key → POST /api/login → store token in localStorage → Bearer auth on all subsequent requests → auto-redirect to /login if unauthenticated
- **Build Status**: Next.js build ✅ (5 new routes, 87.6 kB static export), Go build ✅, Go tests ✅
- **Next Turn Directive**: Wire up the full stack end-to-end — run `mcp-guard serve` with admin enabled, verify workspace login works against live API, then bundle for Gumroad

### [2026-07-19 22:00] — Admin API + Proxy Stats + Workspace Infrastructure
- **State**: Success — Admin HTTP server built and wired
- **Files Created**:
  - `mcp-guard/internal/proxy/stats.go` — atomic counters (total, allow, block, hitl, rate_limit, injection)
  - `mcp-guard/internal/admin/server.go` — HTTP server with auth (Bearer API key), CORS, 6 endpoints
- **Files Modified**:
  - `mcp-guard/internal/proxy/proxy.go` — added `Stats *Stats` field to Options, `Stats()` accessor
  - `mcp-guard/internal/proxy/tcp.go` — stats increments at all decision points
  - `mcp-guard/internal/config/config.go` — added `AdminConfig` (enabled, listen, api_key, origins)
  - `mcp-guard/internal/policy/types.go` — added `RawPolicy`, `Engine.List()`, `Engine.Replace()`
  - `mcp-guard/internal/audit/logger.go` — added `Logger.Recent(n)` method
  - `mcp-guard/cmd/serve.go` — wired admin server start/stop, proxyStats, config closure
- **Admin API Endpoints**: POST /api/login, GET /api/status, GET /api/policies, POST /api/policies, GET /api/audit, GET /api/config
- **Config Defaults**: admin.listen = ":9090", admin.api_key = "controlplane-dev-key"
- **Build Status**: Go build ✅, all 68+ tests pass ✅, Next.js build ✅ (87.5 kB static export)
- **Next Turn Directive**: Build workspace UI pages (login, dashboard, policies, audit, settings)

### [2026-07-19 17:30] — Split MCP Guard into Standalone Project
- **State**: Success — mcp-guard Go implementation moved to its own repo (`ravikumarve/mcp-guard`)
- **Actions**: Cleaned ControlPlane AI repo of all Go code. Now holds only company docs.
- **Next Turn Directive**: Build ControlPlane AI company foundation documents

### [2026-07-19 21:00] — Merged mcp-guard back into ControlPlane monorepo
- **State**: Success — everything unified under one roof
- **Actions**:
  - Moved `/media/matrix/DATA/opencode_projects/mcp-guard/` back into `ControlPlane AI/mcp-guard/`
  - Fixed Go module path from `github.com/matrix/mcp-guard` → `github.com/ravikumarve/ControlPlane/mcp-guard` (20 files updated)
  - Updated .gitignore to include mcp-guard binary patterns
  - Rewrote landing page HTML to sell **ControlPlane** as the company/product, with mcp-guard positioned as the implementation (12 mcp-guard references kept intentionally as the product name)
  - Updated README with unified structure showing mcp-guard/ under ControlPlane
  - Verified: Go build ✅, all 68+ tests pass ✅, Next.js build ✅
- **Architectural Decision**: ControlPlane is the company AND product. mcp-guard is the Go implementation package name. Everything lives in one monorepo. Two separate directories was the source of all confusion.
- **Current Repo Structure**:
  ```
  ControlPlane AI/
  ├── mcp-guard/           ← Go backend (product code, 68+ tests)
  ├── app/                 ← Next.js landing page
  ├── components/          ← React components
  ├── docs/                ← Company docs
  ├── controlplane_ai_landing_page.html  ← Design mockup
  └── README.md
  ```
- **Next Turn Directive**: Adapt the landing page HTML design into the Next.js scaffold, or start building the product features.

### [2026-07-19 20:30] — Next.js project scaffold complete + README updated
- **State**: Success — full frontend scaffold built and verified
- **Actions**:
  - Created Next.js 14 project with App Router, Tailwind, TypeScript strict
  - Tailwind config with all BRAND.md color tokens (primary, accent, success, danger, warning, surface)
  - Component primitives: Button (5 variants x 3 sizes), Card (3 variants), Badge (5 variants)
  - Shared components: Navbar (sticky, glassmorphism), Footer (3-column)
  - Landing sections: Hero, Features (4-column grid with lucide icons), CTA
  - `lib/utils.ts` with cn() helper
  - Static export config — no server, no database, no backend
  - Build verified: compiles clean, exports to `out/` at 87.5 kB first load
  - README.md rewritten with current structure, stack table, and dev commands
- **Architectural Decisions**:
  - No backend, no database — static export only. Landing page needs zero ops.
  - No component library — Button/Card/Badge are hand-written, match brand exactly.
  - Main branch development — branching is overhead for a solo dev.
  - Conventional commits (feat/fix/style/docs/chore).
- **Next Turn Directive**: Start landing page design, or move to another project.

### [2026-07-19 20:00] — Frontend Architecture Plan (design deferred)
- **State**: Success — FRONTEND-PLAN.md created, no code written
- **Architectural Decision**: Landing page is the design system origin. No frontend code until landing page is designed. Policy Console (Q4) inherits all tokens from the landing page. Stack: Next.js 14+ App Router / Tailwind / TypeScript / motion/react / lucide-react. No component library (shadcn/MUI), no CSS-in-JS, no state management library.
- **Key Outputs**: `docs/FRONTEND-PLAN.md` — route map, component hierarchy, design token mapping, build/deployment strategy
- **Next Turn Directive**: Move to another project (FounderLens, StudyAI, VAJRA) or wait for landing page design phase

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
