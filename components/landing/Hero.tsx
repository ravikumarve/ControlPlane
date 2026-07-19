export function Hero() {
  return (
    <section className="grid min-h-[85vh] grid-cols-1 items-center gap-16 border-none lg:grid-cols-2">
      <div>
        <div className="mb-8 inline-flex items-center gap-2.5 border border-border-solid bg-surface px-4 py-1.5 font-mono text-xs uppercase tracking-wider text-gray-400">
          Product: <span className="text-cyan-500">mcp-guard v0.1.0</span>
        </div>

        <p className="mb-2 font-mono text-sm uppercase tracking-widest text-gray-400">
          ControlPlane
        </p>

        <h1 className="text-hero mb-6 font-bold leading-tight tracking-tighter">
          Security for the<br />Agent Ecosystem.
        </h1>

        <p className="mb-12 max-w-[500px] text-base leading-relaxed text-gray-400">
          We build security infrastructure for organizations deploying autonomous AI agents. Our MCP Security Gateway sits between orchestration engines and their tools — enforcing policy, blocking injection, and auditing every call.
        </p>

        <div className="flex gap-4">
          <a
            href="#gateway"
            className="inline-flex items-center justify-center border border-white bg-white px-6 py-3 font-mono text-xs font-bold uppercase tracking-wider text-[#050505] no-underline transition-all hover:border-cyan-500 hover:bg-cyan-500 hover:shadow-lg"
            style={{ boxShadow: "0 0 20px rgba(6,182,212,0.1)" }}
          >
            Explore mcp-guard
          </a>
          <a
            href="docs/ARCHITECTURE.md"
            className="inline-flex items-center justify-center border border-border-solid bg-surface px-6 py-3 font-mono text-xs font-bold uppercase tracking-wider text-white no-underline transition-all hover:border-white"
          >
            Read the Spec
          </a>
        </div>
      </div>

      {/* Terminal Mockup */}
      <div className="flex w-full flex-col border border-border-solid shadow-2xl" style={{ background: "#050505", boxShadow: "0 20px 40px rgba(0,0,0,0.8)" }}>
        <div className="flex items-center justify-between border-b border-border-solid bg-surface px-6 py-3 font-mono text-xs text-gray-400">
          <span>mcp-guard :: controlplane</span>
          <span className="text-cyan-500">ACTIVE</span>
        </div>
        <div className="overflow-x-auto px-6 py-8 font-mono text-sm leading-loose text-white">
          <span className="text-gray-500">[08:42:11]</span> <span className="text-cyan-500">INFO</span> Starting proxy on :8080<br />
          <span className="text-gray-500">[08:42:12]</span> <span className="text-cyan-500">INFO</span> Loaded 12 policies (YAML).<br />
          <span className="text-gray-500">[08:42:15]</span> <span className="text-gray-400">REQ</span>  Agent_01 → Tool: <span className="font-mono">read_repo</span><br />
          <span className="text-gray-500">[08:42:15]</span> <span className="text-cyan-500">ALLOW</span> Policy match (3ms).<br /><br />
          <span className="text-gray-500">[08:43:02]</span> <span className="text-gray-400">REQ</span>  Agent_04 → Tool: <span className="font-mono">drop_table</span><br />
          <span className="text-gray-500">[08:43:02]</span> <span className="font-bold text-orange-500">BLOCK</span> Unauthorized tool access.<br />
          <span className="text-gray-500">[08:43:02]</span> <span className="font-bold text-orange-500">TRIP</span> Circuit breaker activated.<br />
        </div>
      </div>
    </section>
  );
}
