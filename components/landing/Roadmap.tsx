export function Roadmap() {
  return (
    <section id="roadmap" className="section-padding section-border">
      <div className="container-content">
        <div className="mb-20">
          <h2 className="text-section-title mb-6 font-bold leading-tight tracking-tighter">
            Product Roadmap.
          </h2>
          <p className="max-w-[650px] text-base leading-relaxed text-gray-400">
            Executing a phased rollout to secure the rapidly expanding autonomous AI ecosystem.
          </p>
        </div>

        <div className="grid grid-cols-1 gap-px border border-border-solid md:grid-cols-3" style={{ background: "#27272a" }}>
          {/* Q3 2026 — Active */}
          <div className="flex flex-col bg-surface p-12" style={{ boxShadow: "inset 0 2px 0 #06b6d4" }}>
            <div className="mb-4 font-mono text-xs uppercase tracking-wider text-cyan-500">Q3 2026</div>
            <h3 className="mb-4 text-2xl font-bold text-white">mcp-guard V1</h3>
            <p className="flex-1 text-sm leading-relaxed text-gray-400">
              Launch of the core MCP Security Gateway. RBAC, injection detection, rate limiting, circuit breaker, and tamper-evident audit trail.
            </p>
            <div className="mt-auto flex items-center gap-2 pt-8 font-mono text-xs uppercase tracking-wider text-gray-500">
              <span className="h-1.5 w-1.5 rounded-full bg-cyan-500" />
              MVP Released
            </div>
          </div>

          {/* Q4 2026 */}
          <div className="flex flex-col bg-surface p-12">
            <div className="mb-4 font-mono text-xs uppercase tracking-wider text-cyan-500">Q4 2026</div>
            <h3 className="mb-4 text-2xl font-bold text-white">Policy Console</h3>
            <p className="flex-1 text-sm leading-relaxed text-gray-400">
              Web-based management console to author, version, and monitor policies across fleets of mcp-guard instances.
            </p>
            <div className="mt-auto pt-8 font-mono text-xs uppercase tracking-wider text-gray-500">In Design</div>
          </div>

          {/* Q1 2027 */}
          <div className="flex flex-col bg-surface p-12">
            <div className="mb-4 font-mono text-xs uppercase tracking-wider text-cyan-500">Q1 2027</div>
            <h3 className="mb-4 text-2xl font-bold text-white">SIEM Integration</h3>
            <p className="flex-1 text-sm leading-relaxed text-gray-400">
              Forward tamper-evident audit trails to Splunk, Datadog, or ELK via standard protocols.
            </p>
            <div className="mt-auto pt-8 font-mono text-xs uppercase tracking-wider text-gray-500">Researching</div>
          </div>
        </div>
      </div>
    </section>
  );
}
