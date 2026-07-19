import { Shield, Key, FileSearch, Gauge } from "lucide-react";

const features = [
  {
    icon: Shield,
    title: "Zero-Trust Policy",
    desc: "Granular RBAC at tool and parameter level. Allow, block, or route to human approval.",
  },
  {
    icon: FileSearch,
    title: "Injection Detection",
    desc: "Scan every tool call for prompt injection, homoglyphs, and parameter smuggling.",
  },
  {
    icon: Gauge,
    title: "Single Binary",
    desc: "Deploy in 5 minutes. No Kubernetes, no database, no SaaS dependency.",
  },
  {
    icon: Key,
    title: "Immutable Audit",
    desc: "HMAC-chained JSONL audit trail. Tamper-evident and verifiable via CLI.",
  },
];

export function Features() {
  return (
    <section className="section-padding bg-gray-50 dark:bg-surface-card/50">
      <div className="container-content">
        <h2 className="text-section-title text-center font-bold text-primary dark:text-white">
          Agent-aware security, not just API gateways
        </h2>
        <div className="mt-12 grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
          {features.map((f) => (
            <div key={f.title} className="group rounded-xl p-6 transition-colors hover:bg-white dark:hover:bg-surface-card">
              <f.icon className="h-8 w-8 text-accent" />
              <h3 className="mt-4 font-semibold">{f.title}</h3>
              <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">{f.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
