"use client";

import { useEffect, useState } from "react";
import { api, AuditEntry } from "@/lib/api";

function DecisionBadge({ decision }: { decision: string }) {
  const styles: Record<string, string> = {
    allow: "bg-green-950/30 text-green-400 border border-green-900/30",
    block: "bg-red-950/30 text-red-400 border border-red-900/30",
    pending: "bg-yellow-950/30 text-yellow-400 border border-yellow-900/30",
    parse_error: "bg-gray-950/30 text-gray-400 border border-gray-800",
  };
  const s = styles[decision] || styles.parse_error;
  return (
    <span className={`text-xs font-mono px-2 py-0.5 rounded-full ${s}`}>
      {decision}
    </span>
  );
}

export default function AuditPage() {
  const [entries, setEntries] = useState<AuditEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filter, setFilter] = useState("all");

  const fetchAudit = async () => {
    try {
      const data = await api.getAudit();
      setEntries(data.entries);
      setError("");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "failed to fetch audit";
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAudit();
  }, []);

  const filtered =
    filter === "all"
      ? entries
      : entries.filter((e) => e.decision === filter);

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Audit Log</h1>
          <p className="text-gray-500 text-mono text-sm mt-1">
            tamper-proof request history
          </p>
        </div>
        <div className="flex items-center gap-3">
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="bg-panel border border-border-solid rounded-lg px-3 py-1.5 text-sm font-mono text-white focus:outline-none focus:border-cyan/50"
          >
            <option value="all">all</option>
            <option value="allow">allow</option>
            <option value="block">block</option>
            <option value="pending">pending</option>
          </select>
          <button
            onClick={fetchAudit}
            className="px-3 py-1.5 text-sm border border-border-solid rounded-lg text-gray-400 hover:text-white hover:border-gray-500 transition-colors"
          >
            refresh
          </button>
          <span className="text-xs text-gray-500 text-mono">
            {entries.length} entries
          </span>
        </div>
      </div>

      {/* Error */}
      {error && (
        <div className="text-red-400 text-sm text-mono bg-red-950/30 border border-red-900/30 rounded-lg px-4 py-3">
          ⚠ {error}
        </div>
      )}

      {/* Loading */}
      {loading && (
        <div className="text-center py-20 text-gray-500 text-mono animate-pulse">
          loading audit log...
        </div>
      )}

      {/* Empty state */}
      {!loading && filtered.length === 0 && (
        <div className="text-center py-20 text-gray-500 text-mono">
          {entries.length === 0
            ? "no audit entries yet — make some proxy requests"
            : "no entries match the selected filter"}
        </div>
      )}

      {/* Audit entries */}
      {!loading && filtered.length > 0 && (
        <div className="space-y-2">
          {filtered.map((entry, i) => (
            <div
              key={i}
              className="bg-panel border border-border-solid rounded-lg px-5 py-4"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1 min-w-0 space-y-2">
                  <div className="flex items-center gap-3 flex-wrap">
                    <DecisionBadge decision={entry.decision} />
                    <span className="font-mono text-sm text-white truncate">
                      {entry.tool}
                    </span>
                    <span className="text-xs text-gray-500 text-mono">
                      {entry.identity || "unknown"}
                    </span>
                  </div>
                  {entry.reason && (
                    <div className="text-xs text-gray-400 text-mono">
                      {entry.reason}
                    </div>
                  )}
                </div>
                <div className="text-right shrink-0">
                  <div className="text-xs text-gray-500 text-mono">
                    {new Date(entry.timestamp).toLocaleTimeString()}
                  </div>
                  {entry.duration > 0 && (
                    <div className="text-xs text-gray-500 text-mono">
                      {entry.duration}ms
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* HMAC chain info */}
      {entries.length > 0 && (
        <div className="text-xs text-gray-600 text-mono bg-panel border border-border-solid rounded-lg px-4 py-3">
          <span className="text-gray-500">🔗 HMAC chain: </span>
          last entry hash:{" "}
          <span className="text-cyan">{entries[entries.length - 1]?.hmac?.slice(0, 16)}...</span>
        </div>
      )}
    </div>
  );
}
