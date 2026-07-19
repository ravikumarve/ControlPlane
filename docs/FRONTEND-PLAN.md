# ControlPlane AI — Frontend Architecture Plan

**Status**: Planning only — no code written yet
**Design System Origin**: Landing page (built first, establishes all UI tokens)

---

## 1. Design System Origin

The **landing page** is the source of truth for all frontend UI at ControlPlane AI.

### Why
- The landing page defines the public face of the company — colors, typography, spacing, tone
- Every subsequent page (Policy Console, dashboard, docs) reuses the same design tokens
- Building one page first forces real design decisions before scaling to multiple views

### What the landing page establishes
- Tailwind theme config (colors, fonts, spacing, breakpoints)
- Component primitives (Button, Card, Nav, Footer, Container, Heading hierarchy)
- Layout patterns (hero section, feature grid, call-to-action, navigation)
- Dark/light mode strategy
- Animation patterns via `motion/react` (entry animations, hover states, page transitions)

### Rule
The landing page repo IS the frontend monorepo seed. Do not create a separate "design system" package — the landing page IS the design system. Policy Console imports from it.

---

## 2. Stack

| Layer | Choice | Why |
|-------|--------|-----|
| Framework | Next.js 14+ (App Router) | Required by project standards; static export for landing page |
| Language | TypeScript (strict) | Type safety, catches config errors at build time |
| Styling | Tailwind CSS v3+ | Utility-first, design tokens in `tailwind.config.ts`, no runtime CSS |
| Animation | motion/react (framer) | Page transitions, scroll animations, micro-interactions |
| Icons | lucide-react | Tree-shakeable, consistent style, MIT license |
| Deployment | GitHub Pages (landing) / Vercel (console) | Static export for marketing, full Next.js for authenticated routes |

### What we explicitly skip
- No component library (shadcn/ui, MUI, Chakra) — own components from landing page
- No CSS-in-JS (emotion, styled-components) — Tailwind only
- No state management library (Redux, Zustand) — Next.js server components + React context for auth only
- No GraphQL — REST API from the Go backend
- No Storybook — landing page IS the component showcase

---

## 3. Route Map

### Landing Page (Next.js static export) — Q3 2026

```
/                    ← Hero, features, pricing, CTA
/pricing             ← Pricing tiers (Community / Pro / Enterprise)
/docs                ← Quickstart, installation, configuration
/docs/*              ← Sub-pages (routing, policies, audit)
/blog                ← Dev.to cross-post archive
/blog/*              ← Individual posts
/security            ← Security.txt style disclosure page
```

### Policy Management Console (authenticated SPA) — Q4 2026

```
/app/login           ← Auth0 / GitHub OAuth login
/app/dashboard       ← Overview: requests, blocks, HITL approvals
/app/policies        ← Policy list, editor, test runner
/app/policies/:id    ← Single policy detail + edit
/app/audit           ← Audit log viewer with search/filter
/app/audit/:id       ← Single audit entry detail
/app/settings        ← Proxy config, rate limits, webhook URLs
/app/settings/api-keys ← API key management
```

### Shared Layout
Landing and Console share the same `tailwind.config.ts` and component primitives — just different layouts (marketing vs. app shell).

---

## 4. Component Hierarchy

```
src/
├── app/                    ← Next.js App Router pages
│   ├── (marketing)/        ← Landing route group
│   │   ├── layout.tsx      ← Marketing layout (nav, footer)
│   │   ├── page.tsx        ← Home hero
│   │   ├── pricing/
│   │   └── docs/
│   └── (app)/              ← Console route group (Q4)
│       ├── layout.tsx      ← App shell (sidebar, topbar)
│       ├── dashboard/
│       └── policies/
├── components/
│   ├── ui/                 ← Primitives (Button, Card, Badge, Input, Modal)
│   │   ├── Button.tsx
│   │   ├── Card.tsx
│   │   ├── Badge.tsx
│   │   └── ...
│   ├── landing/            ← Marketing-specific (Hero, Features, PricingCard, CTA)
│   ├── console/            ← App-specific (PolicyTable, AuditLog, HITLModal)
│   └── shared/             ← Cross-cutting (Navbar, Footer, ThemeToggle, SEO)
├── lib/
│   ├── utils.ts            ← cn() helper, formatters
│   └── api.ts              ← API client for Go backend (console only)
├── hooks/
│   ├── use-auth.ts         ← Auth session
│   └── use-policies.ts     ← SWR/React Query for policy data
├── styles/
│   └── globals.css         ← Tailwind directives + CSS custom properties
├── tailwind.config.ts      ← THE source of truth for design tokens
└── motion.config.ts        ← Shared animation variants
```

### Component design rules
1. **Primitives are composable** — `Button` accepts `variant`, `size`, `asChild` props. Never hardcode styles in page components.
2. **No prop drilling beyond 2 levels** — Use React context for theme/auth, server components for data fetching.
3. **Server components by default** — `"use client"` only for interactivity (forms, animations, toggles).
4. **motion/react only in client components** — Leaf components only, not layout wrappers.

---

## 5. Design Token Mapping (from BRAND.md)

| BRAND.md Token | Tailwind Config | CSS Variable |
|----------------|-----------------|--------------|
| Deep Blue #1A237E | `primary` | `--color-primary` |
| Cyan #00BCD4 | `accent` | `--color-accent` |
| Green #4CAF50 | `success` | `--color-success` |
| Red #F44336 | `danger` | `--color-danger` |
| Amber #FFC107 | `warning` | `--color-warning` |
| Dark #0D1117 | `surface-dark` | `--color-surface-dark` |
| Light #FAFAFA | `surface-light` | `--color-surface-light` |
| Inter / System UI | `fontFamily.sans` | `--font-sans` |
| JetBrains Mono | `fontFamily.mono` | `--font-mono` |

These are defined in `tailwind.config.ts` and exported as CSS custom properties. The landing page is where these get finalized before any building begins.

---

## 6. Build & Deployment

### Landing page
- `next build && next export` (static HTML)
- Deploy to GitHub Pages or Vercel (free tier)
- No server runtime needed — pure static assets

### Policy Console (Q4 2026)
- Full Next.js deployment to Vercel
- API routes proxy to the Go backend (mcp-guard)
- Auth middleware for route protection
- ISR for audit log pages (revalidate on new data)

---

## 7. What NOT to Do

- Do not build any frontend code until the landing page design is finalized
- Do not start with the Policy Console — the design decisions will be wrong without the landing page foundation
- Do not use a pre-built template or theme — the landing page IS the theme
- Do not create a separate `/components` library published to npm — keep it in-repo
- Do not add a CSS preprocessor (Sass, Less) — Tailwind handles everything
- Do not optimize for Lighthouse score at the expense of design quality (the audience is technical)
