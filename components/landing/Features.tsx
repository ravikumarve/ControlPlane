export function Features() {
  return (
    <section id="gateway" className="section-padding section-border">
      <div className="container-content">
        <div className="mb-20">
          <p className="mb-2 font-mono text-xs uppercase tracking-widest text-cyan-500">
            By ControlPlane
          </p>
          <h2 className="text-section-title mb-6 font-bold leading-tight tracking-tighter">
            The MCP Security<br />Gateway.
          </h2>
          <p className="max-w-[650px] text-base leading-relaxed text-gray-400">
            mcp-guard is a zero-trust sidecar for Model Context Protocol (MCP) agents. It prevents injection attacks from becoming database disasters.
          </p>
        </div>

        {/* Bento Grid */}
        <div className="grid grid-cols-12 gap-px border border-border-solid" style={{ background: "#27272a" }}>
          {/* Architecture */}
          <div className="col-span-12 bg-surface p-16 transition-colors hover:bg-panel md:col-span-8">
            <span className="mb-8 flex items-center gap-2.5 font-mono text-xs uppercase tracking-wider text-gray-400 before:text-cyan-500 before:content-['['] after:text-cyan-500 after:content-[']']">
              Architecture
            </span>
            <h3 className="text-card-title mb-4 font-bold text-white">Single Binary Sidecar</h3>
            <p className="flex-1 text-base leading-relaxed text-gray-400">
              Written in Go. No Kubernetes required. No heavy SaaS dependencies. Deploy the 8.5MB static binary alongside your agent orchestrator, point your tools to it, and instantly secure the perimeter.
            </p>
          </div>

          {/* Injection */}
          <div className="col-span-12 bg-surface p-16 transition-colors hover:bg-panel md:col-span-4">
            <span className="mb-8 flex items-center gap-2.5 font-mono text-xs uppercase tracking-wider text-gray-400 before:text-orange-500 before:content-['['] after:text-orange-500 after:content-[']']">
              Threat Vector
            </span>
            <h3 className="text-card-title mb-4 font-bold text-white">Injection Defense</h3>
            <p className="flex-1 text-base leading-relaxed text-gray-400">
              Malicious prompts can hijack agent reasoning. We validate parameters against strict patterns before the tool is ever invoked — 10 detection categories including homoglyphs and depth bombs.
            </p>
          </div>

          {/* RBAC */}
          <div className="col-span-12 bg-surface p-16 transition-colors hover:bg-panel md:col-span-6">
            <span className="mb-8 flex items-center gap-2.5 font-mono text-xs uppercase tracking-wider text-gray-400 before:text-cyan-500 before:content-['['] after:text-cyan-500 after:content-[']']">
              Policy Engine
            </span>
            <h3 className="text-card-title mb-4 font-bold text-white">YAML-Based RBAC</h3>
            <p className="flex-1 text-base leading-relaxed text-gray-400">
              Define granular access controls with simple YAML. Allow, block, or route to human approval — at the tool level and parameter level. Glob matching on identities, tools, and parameter values.
            </p>
          </div>

          {/* HITL */}
          <div className="col-span-12 bg-surface p-16 transition-colors hover:bg-panel md:col-span-6">
            <span className="mb-8 flex items-center gap-2.5 font-mono text-xs uppercase tracking-wider text-gray-400 before:text-cyan-500 before:content-['['] after:text-cyan-500 after:content-[']']">
              Compliance
            </span>
            <h3 className="text-card-title mb-4 font-bold text-white">Human-In-The-Loop</h3>
            <p className="flex-1 text-base leading-relaxed text-gray-400">
              For destructive actions (financial transfers, database drops), configure policies that pause execution and send an approval request via Slack, Discord, or custom webhook.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
