export function CTA() {
  return (
    <section className="section-padding">
      <div className="container-content">
        <div
          className="mx-0 my-16 border border-border-solid bg-surface p-16 text-center sm:p-24"
          style={{
            backgroundImage: "radial-gradient(circle at 50% 0%, rgba(6,182,212,0.05) 0%, transparent 70%)",
          }}
        >
          <h2 className="text-section-title mb-6 font-bold leading-tight tracking-tighter">
            Fortify your agent layer.
          </h2>
          <p className="mx-auto mb-12 max-w-[600px] text-base leading-relaxed text-gray-400">
            Stop trusting autonomous agents with bare API keys. Deploy ControlPlane&apos;s mcp-guard and enforce zero-trust policies today.
          </p>
          <div className="flex items-center justify-center gap-4">
            <a
              href="https://github.com/ravikumarve/ControlPlane"
              className="inline-flex items-center justify-center border border-white bg-white px-6 py-3 font-mono text-xs font-bold uppercase tracking-wider text-[#050505] no-underline transition-all hover:border-cyan-500 hover:bg-cyan-500"
            >
              View on GitHub
            </a>
            <a
              href="docs/PRD.md"
              className="inline-flex items-center justify-center border border-border-solid bg-surface px-6 py-3 font-mono text-xs font-bold uppercase tracking-wider text-white no-underline transition-all hover:border-white"
            >
              Read PRD
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}
