import Link from "next/link";

export function Navbar() {
  return (
    <nav className="container-content flex items-center justify-between border-b border-border-solid py-6">
      <Link href="/" className="flex items-center gap-3 font-mono text-sm font-extrabold tracking-tight no-underline text-white">
        <span className="flex h-4 w-4 items-center justify-center shadow-lg" style={{ background: "#06b6d4", boxShadow: "0 0 15px rgba(6,182,212,0.1)" }}>
          <span className="h-1.5 w-1.5" style={{ background: "#050505" }} />
        </span>
        ControlPlane
      </Link>

      <div className="flex gap-12">
        <a href="#gateway" className="font-mono text-xs uppercase tracking-wider text-gray-400 no-underline transition-colors hover:text-cyan-500">
          mcp-guard
        </a>
        <a href="#stack" className="font-mono text-xs uppercase tracking-wider text-gray-400 no-underline transition-colors hover:text-cyan-500">
          Stack
        </a>
        <a href="#roadmap" className="font-mono text-xs uppercase tracking-wider text-gray-400 no-underline transition-colors hover:text-cyan-500">
          Roadmap
        </a>
        <Link href="/login" className="font-mono text-xs uppercase tracking-wider text-orange no-underline transition-colors hover:text-orange/80">
          Workspace
        </Link>
      </div>

      <a
        href="https://github.com/ravikumarve/ControlPlane"
        className="inline-flex items-center justify-center border border-border-solid bg-surface px-5 py-3 font-mono text-xs font-bold uppercase tracking-wider text-white no-underline transition-all hover:border-white"
      >
        GitHub
      </a>
    </nav>
  );
}
