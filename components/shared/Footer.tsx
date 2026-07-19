import Link from "next/link";

export function Footer() {
  return (
    <footer className="container-content pb-12 pt-20">
      <div className="flex flex-col justify-between gap-16 lg:flex-row lg:items-start">
        <div className="max-w-xs">
          <Link href="/" className="flex items-center gap-3 font-mono text-sm font-extrabold tracking-tight no-underline text-white">
            <span className="flex h-4 w-4 items-center justify-center" style={{ background: "#06b6d4" }}>
              <span className="h-1.5 w-1.5" style={{ background: "#050505" }} />
            </span>
            CONTROLPLANE
          </Link>
          <p className="mt-6 text-sm leading-relaxed text-gray-400">
            Zero-trust security infrastructure for the AI agent ecosystem. Enforcing policy, preventing injection, and securing the bridge to your servers.
          </p>
        </div>

        <div className="flex gap-16">
          <div>
            <h5 className="mb-6 font-mono text-xs uppercase tracking-widest text-white">Company</h5>
            <ul className="flex flex-col gap-4">
              <li><Link href="/docs/COMPANY.md" className="text-sm text-gray-400 no-underline transition-colors hover:text-cyan-500">About</Link></li>
              <li><Link href="/docs/ARCHITECTURE.md" className="text-sm text-gray-400 no-underline transition-colors hover:text-cyan-500">Architecture</Link></li>
              <li><Link href="/docs/SECURITY.md" className="text-sm text-gray-400 no-underline transition-colors hover:text-cyan-500">Security</Link></li>
            </ul>
          </div>
          <div>
            <h5 className="mb-6 font-mono text-xs uppercase tracking-widest text-white">Product</h5>
            <ul className="flex flex-col gap-4">
              <li><a href="#gateway" className="text-sm text-gray-400 no-underline transition-colors hover:text-cyan-500">mcp-guard</a></li>
              <li><Link href="/docs/FRONTEND-PLAN.md" className="text-sm text-gray-400 no-underline transition-colors hover:text-cyan-500">Frontend Plan</Link></li>
              <li><Link href="/docs/CONTRIBUTING.md" className="text-sm text-gray-400 no-underline transition-colors hover:text-cyan-500">Contributing</Link></li>
            </ul>
          </div>
        </div>
      </div>

      <div className="mt-24 flex flex-col justify-between gap-4 border-t border-border-solid pt-8 font-mono text-xs uppercase tracking-widest text-gray-500 sm:flex-row">
        <div>© 2026 ControlPlane. MIT licensed.</div>
        <div>STATUS: <span className="text-cyan-500">ALL SYSTEMS NOMINAL</span></div>
      </div>
    </footer>
  );
}
