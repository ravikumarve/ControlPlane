# Brand Guidelines — ControlPlane AI

This document defines the visual and verbal identity of ControlPlane AI. Follow these guidelines when creating content, designing interfaces, or communicating about our products.

---

## Brand Name

The company name is **ControlPlane AI**.

- Always capitalize **C**, **P**, **A**, and **I**
- Never use "Control Plane AI" (separated), "Controlplane AI" (lowercase p), or "ControlPlaneAI" (no space)
- Use the full name "ControlPlane AI" on first reference in any document or communication
- The abbreviated form "CP AI" is acceptable only in internal contexts or space-constrained UIs
- Do not use "CPAI" or "cp-ai" in external communications

---

## Tagline

**"Security infrastructure for the AI agent ecosystem"**

Usage rules:
- Capitalize as shown — lowercase "security," "infrastructure," and "ecosystem"
- Use below or beside the logo as a subtitle
- May be omitted in contexts where the brand name alone is sufficient (e.g., social media avatars)
- Do not modify or extend the tagline

---

## Colors

### Primary Palette

| Role | Name | Hex | Usage |
|------|------|-----|-------|
| Primary | Deep Blue | `#1A237E` | Logo, headings, primary buttons, links |
| Accent | Cyan | `#00BCD4` | Interactive elements, agent-related highlights, active states |

Deep Blue conveys trust, authority, and security — the foundation of our brand. Cyan represents intelligence, agent communication, and the dynamic nature of the AI ecosystem.

### Semantic Palette

| Role | Name | Hex | Usage |
|------|------|-----|-------|
| Success | Green | `#4CAF50` | Approved actions, healthy status, positive indicators |
| Danger | Red | `#F44336` | Blocked actions, errors, security violations |
| Warning | Amber | `#FFC107` | Pending actions, rate limits, attention-required states |

### Background Palette

| Context | Hex | Usage |
|---------|-----|-------|
| Dark theme background | `#0D1117` | Documentation, CLI terminals, dashboard dark mode |
| Light theme background | `#FAFAFA` | Website, documentation light mode, marketing materials |

### Color Usage Guidelines

- **Do** use Deep Blue as the dominant color (approximately 60% of color surface)
- **Do** use Cyan sparingly for call-to-action elements (approximately 10%)
- **Do** ensure text-on-background contrast ratios meet WCAG AA standards (4.5:1 for normal text, 3:1 for large text)
- **Do not** use Cyan for error states — reserve Red for errors
- **Do not** use Deep Blue text on Dark Background — use white or light gray instead
- **Do not** apply gradient blends to the primary logo

---

## Typography

### Interface Fonts

Digital interfaces should use system-native font stacks for performance and consistency.

**Primary stack (UI):**

```css
font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
             Oxygen, Ubuntu, Cantarell, 'Helvetica Neue', sans-serif;
```

**Code font (monospace):**

```css
font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Consolas',
             'Courier New', monospace;
```

### Type scale

| Style | Size | Weight | Usage |
|-------|------|--------|-------|
| H1 | 2.5rem / 40px | 700 (Bold) | Page titles, hero headings |
| H2 | 2rem / 32px | 700 (Bold) | Section headings |
| H3 | 1.5rem / 24px | 600 (Semibold) | Subsection headings |
| Body | 1rem / 16px | 400 (Regular) | Paragraph text |
| Small | 0.875rem / 14px | 400 (Regular) | Captions, metadata, footnotes |
| Code | 0.875rem / 14px | 400 (Regular) | Inline code, code blocks |

### Typography Guidelines

- **Do** use sentence case for headings ("Security infrastructure for..." not "Security Infrastructure For...")
- **Do** limit line lengths to 70-80 characters for readability
- **Do** use code font for commands, file paths, code blocks, and protocol references (MCP, JSON-RPC)
- **Do not** use all-caps for emphasis — use bold or italic instead
- **Do not** use decorative or display fonts in any context

---

## Voice & Tone

### Core Principles

| Principle | Description |
|-----------|-------------|
| **Authoritative but approachable** | We know security deeply, but we explain it without condescension |
| **Technical but not academic** | We use precise terminology without unnecessary jargon or citations |
| **Security-conscious but not fear-mongering** | We state risks factually. We do not manufacture urgency or exploit anxiety |
| **Direct, clear, jargon-aware** | We say what we mean. When we use specialized terms, we define them or provide context |

### Voice Characteristics

- **Active voice** — "The proxy blocks unauthorized calls" not "Unauthorized calls are blocked by the proxy"
- **Concise sentences** — Prefer 15-25 word sentences. Break complex ideas into multiple sentences
- **Second-person ("you")** — Address the reader directly in documentation and guides
- **Specific over vague** — "Prevents prompt injection via NFKC normalization" not "Provides advanced security features"
- **Honest about trade-offs** — Acknowledge limitations. Security is about trade-offs, and we respect our users' ability to understand them

### Voice Anti-Patterns

| Avoid | Instead use |
|-------|-------------|
| "Revolutionary," "game-changing," "cutting-edge" | Precise descriptions of what the product does |
| "Don't worry," "trust us" | Evidence, code, architecture documentation |
| "We're the only ones who..." | "We focus on..." — acknowledge the ecosystem |
| "Simply," "just," "obviously" | Direct instructions without qualifiers |
| Fear-based urgency ("before it's too late") | Factual risk assessment |

### Examples

**Good:**
> "MCP has no built-in authentication or authorization. If your agent can reach a tool, it can call it. mcp-guard sits between the two and enforces policy at every request."

**Avoid:**
> "Your agents are completely exposed and hackers will steal all your data! You need our revolutionary security solution immediately."

---

## Logo Usage

### Primary Logo

The ControlPlane AI logo is a text-based wordmark:

- **Text**: "ControlPlane AI" set in bold weight
- **Color**: Primary Deep Blue (`#1A237E`) on light backgrounds
- **Color**: White or light gray on dark backgrounds
- **No icon or symbol accompany the primary logo at this time**

### Logo Placement

- Minimum clear space: the height of the letter "C" on all sides
- Minimum size: 120px wide (digital), 1.5 inches (print)
- Do not rotate, stretch, skew, or apply effects to the logo
- Do not place the logo on busy backgrounds or images without a solid backing

### Favicon

The favicon uses the initials "CP" centered in a Deep Blue square, with the "C" in white and "P" in Cyan. Minimum size: 16x16px.

---

## Product Names

| Name | Context | Usage Rules |
|------|---------|-------------|
| **MCP Security Gateway** | Product name | Capitalize each word. Use in marketing, documentation, and product descriptions. Always include "MCP" prefix on first reference. |
| **mcp-guard** | Repository / CLI name | Lowercase, hyphenated. Use when referring to the repository, the binary, or CLI commands. Code-font or inline-code formatting preferred. |

### Naming Rules

- Always refer to the company as **ControlPlane AI** — never use just "ControlPlane" in external contexts
- On second and subsequent references within a document, "the gateway" or "the proxy" may be used for the MCP Security Gateway
- The CLI binary is invoked as `mcp-guard` — e.g., `mcp-guard serve`, `mcp-guard logs --verify`
- Do not refer to the product as "ControlPlane AI Gateway" — this creates confusion between company and product
- Do not use "CP AI Gateway" or similar abbreviations in external communications

---

## Attribution

When using ControlPlane AI brand assets externally:

- Open-source projects using our code: attribute "ControlPlane AI" in your project's README or About page
- Articles or reviews mentioning us: use "ControlPlane AI" with a link to https://controlplane.dev
- Derivative works: clearly indicate they are not officially endorsed unless you have written permission

---

_These brand guidelines are maintained by the ControlPlane AI team and will be updated as the brand evolves._
