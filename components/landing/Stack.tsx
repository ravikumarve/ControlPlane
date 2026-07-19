export function Stack() {
  return (
    <section id="stack" className="section-padding section-border">
      <div className="container-content grid grid-cols-1 items-center gap-24 lg:grid-cols-2">
        <div>
          <h2 className="text-section-title mb-6 font-bold leading-tight tracking-tighter">
            Air-Gapped<br />Resilience.
          </h2>
          <p className="mb-8 text-base leading-relaxed text-gray-400">
            ControlPlane is engineered for defense, healthcare, and enterprise environments. The product runs locally. Your secrets never leave the wire.
          </p>

          <div className="mt-12 grid grid-cols-2 gap-px border border-border-solid" style={{ background: "#27272a" }}>
            <div className="bg-surface p-10">
              <div className="mb-2 flex items-center gap-2.5 font-mono text-4xl font-bold text-white">
                <span className="text-cyan-500">[</span> &lt;4ms <span className="text-cyan-500">]</span>
              </div>
              <div className="font-mono text-xs uppercase tracking-wider text-gray-400">Added Latency</div>
            </div>
            <div className="bg-surface p-10">
              <div className="mb-2 flex items-center gap-2.5 font-mono text-4xl font-bold text-white">
                <span className="text-cyan-500">[</span> 8.5MB <span className="text-cyan-500">]</span>
              </div>
              <div className="font-mono text-xs uppercase tracking-wider text-gray-400">Binary Size</div>
            </div>
          </div>
        </div>

        <div>
          <ul className="border-t border-border-solid">
            <li className="flex items-center justify-between border-b border-border-solid py-6 text-sm text-gray-400">
              Company <span className="font-mono text-xs text-white">ControlPlane</span>
            </li>
            <li className="flex items-center justify-between border-b border-border-solid py-6 text-sm text-gray-400">
              Security Gateway <span className="font-mono text-xs text-white">Go 1.24 (mcp-guard)</span>
            </li>
            <li className="flex items-center justify-between border-b border-border-solid py-6 text-sm text-gray-400">
              Policy Engine <span className="font-mono text-xs text-white">YAML-based RBAC</span>
            </li>
            <li className="flex items-center justify-between border-b border-border-solid py-6 text-sm text-gray-400">
              Protocol <span className="font-mono text-xs text-white">MCP 2026-07-28</span>
            </li>
            <li className="flex items-center justify-between border-b border-border-solid py-6 text-sm text-gray-400">
              Audit Trail <span className="font-mono text-xs text-white">HMAC-chained JSONL</span>
            </li>
            <li className="flex items-center justify-between border-b border-border-solid py-6 text-sm text-gray-400">
              Landing Page <span className="font-mono text-xs text-white">Next.js 14 + Tailwind</span>
            </li>
          </ul>
        </div>
      </div>
    </section>
  );
}
