"use client";

import { useEffect, useState } from "react";
import { api, StatsSnapshot } from "@/lib/api";

interface StatCardProps {
  label: string;
  value: number;
  icon: string;
  color: "cyan" | "orange" | "red" | "green" | "yellow" | "gray";
  subtitle?: string;
}

function StatCard({ label, value, icon, color, subtitle }: StatCardProps) {
  const colorMap = {
    cyan: "border-cyan/20 bg-cyan-dim",
    orange: "border-orange/20 bg-orange-dim",
    red: "border-red-500/20 bg-red-950/30",
    green: "border-green-500/20 bg-green-950/30",
    yellow: "border-yellow-500/20 bg-yellow-950/30",
    gray: "border-gray-500/20 bg-gray-950/30",
  };

  return (
    <div
      className={`rounded-xl border p-5 ${colorMap[color]} transition-colors`}
    >
      <div className="flex items-center justify-between mb-3">
        <span className="text-mono text-lg">{icon}</span>
      </div>
      <div className="text-2xl font-bold font-mono">{value.toLocaleString()}</div>
      <div className="text-sm text-gray-400 text-mono mt-1">{label}</div>
      {subtitle && (
        <div className="text-xs text-gray-500 text-mono mt-1">{subtitle}</div>
      )}
    </div>
  );
}

export default function DashboardPage() {
  const [stats, setStats] = useState<StatsSnapshot | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [autoRefresh, setAutoRefresh] = useState(true);

  const fetchStats = async () => {
    try {
      const data = await api.getStatus();
      setStats(data.stats);
      setError("");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "failed to fetch stats";
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
    if (!autoRefresh) return;
    const interval = setInterval(fetchStats, 5000);
    return () => clearInterval(interval);
  }, [autoRefresh]);

  const totalRequests = stats
    ? stats.allowed + stats.blocked + stats.hitl_pending + stats.rate_limited
    : 0;

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Dashboard</h1>
          <p className="text-gray-500 text-mono text-sm mt-1">
            proxy status overview
          </p>
        </div>
        <div className="flex items-center gap-3">
          <label className="flex items-center gap-2 text-sm text-gray-400">
            <input
              type="checkbox"
              checked={autoRefresh}
              onChange={(e) => setAutoRefresh(e.target.checked)}
              className="rounded border-border-solid bg-void text-cyan focus:ring-cyan/30"
            />
            auto-refresh
          </label>
          <button
            onClick={fetchStats}
            className="px-3 py-1.5 text-sm border border-border-solid rounded-lg text-gray-400 hover:text-white hover:border-gray-500 transition-colors"
          >
            refresh
          </button>
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
          loading stats...
        </div>
      )}

      {/* Stats Grid */}
      {stats && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4">
          <StatCard
            label="Total Requests"
            value={totalRequests}
            icon="∑"
            color="cyan"
          />
          <StatCard
            label="Allowed"
            value={stats.allowed}
            icon="✓"
            color="green"
            subtitle={
              totalRequests > 0
                ? `${((stats.allowed / totalRequests) * 100).toFixed(1)}%`
                : "0%"
            }
          />
          <StatCard
            label="Blocked"
            value={stats.blocked}
            icon="✕"
            color="red"
            subtitle={
              totalRequests > 0
                ? `${((stats.blocked / totalRequests) * 100).toFixed(1)}%`
                : "0%"
            }
          />
          <StatCard
            label="HITL Pending"
            value={stats.hitl_pending}
            icon="◷"
            color="yellow"
          />
          <StatCard
            label="Rate Limited"
            value={stats.rate_limited}
            icon="⊘"
            color="orange"
          />
          <StatCard
            label="Injection Blocked"
            value={stats.injection_blocked}
            icon="⚡"
            color="gray"
          />
        </div>
      )}

      {/* Status indicator */}
      {stats && (
        <div className="flex items-center gap-2 text-sm text-gray-500 text-mono bg-panel border border-border-solid rounded-lg px-4 py-3">
          <span className="w-2 h-2 rounded-full bg-green-500 inline-block" />
          proxy is running
        </div>
      )}
    </div>
  );
}
