import { Button } from "@/components/ui/Button";

export function CTA() {
  return (
    <section className="section-padding">
      <div className="container-content">
        <div className="rounded-2xl bg-primary px-8 py-16 text-center text-white">
          <h2 className="text-section-title font-bold">Ready to secure your agents?</h2>
          <p className="mt-4 text-primary-100">
            Deploy the MCP Security Gateway today. Community edition is free and open source.
          </p>
          <div className="mt-8 flex items-center justify-center gap-4">
            <Button variant="secondary" size="lg">
              Get Started
            </Button>
            <Button variant="outline" size="lg" className="border-white/20 text-white hover:bg-white/10">
              View on GitHub
            </Button>
          </div>
        </div>
      </div>
    </section>
  );
}
