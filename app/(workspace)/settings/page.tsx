"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";

export default function SettingsPage() {
  const [config, setConfig] = useState<Record<string, unknown> | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  const fetchConfig = async () => {
    try {
      const data = await api.getConfig();
      setConfig(data.config as Record<string, unknown>);
      setError("");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "failed to fetch config";
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConfig();
  }, []);

  const copyToClipboard = () => {
    if (!config) return;
    navigator.clipboard.writeText(JSON.stringify(config, null, 2));
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  if (loading) {
    return (
      <div className="text-center py-20 text-gray-500 text-mono animate-pulse">
        loading settings...
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Settings</h1>
          <p className="text-gray-500 text-mono text-sm mt-1">
            current mcp-guard configuration
          </p>
        </div>
        <button
          onClick={copyToClipboard}
          className="px-3 py-1.5 text-sm border border-border-solid rounded-lg text-gray-400 hover:text-white hover:border-gray-500 transition-colors"
        >
          {copied ? "copied!" : "copy config"}
        </button>
      </div>

      {/* Error */}
      {error && (
        <div className="text-red-400 text-sm text-mono bg-red-950/30 border border-red-900/30 rounded-lg px-4 py-3">
          ⚠ {error}
        </div>
      )}

      {/* Config display */}
      {config && (
        <div className="bg-panel border border-border-solid rounded-xl overflow-hidden">
          <div className="px-5 py-3 border-b border-border-solid flex items-center justify-between">
            <span className="text-sm text-gray-400 text-mono">mcp-guard.yaml</span>
            <span className="text-xs text-gray-500">read-only</span>
          </div>
          <pre className="p-5 text-sm font-mono text-gray-300 overflow-x-auto leading-relaxed">
            {JSON.stringify(config, null, 2)}
          </pre>
        </div>
      )}

      {/* Quick info cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-panel border border-border-solid rounded-xl p-5">
          <div className="text-xs text-gray-500 text-mono mb-1">Admin API</div>
          <div className="text-sm font-mono text-cyan">:9090</div>
        </div>
        <div className="bg-panel border border-border-solid rounded-xl p-5">
          <div className="text-xs text-gray-500 text-mono mb-1">Proxy Mode</div>
          <div className="text-sm font-mono text-white">
            {(config as Record<string, unknown>)?.proxy ? (
              ((config as Record<string, unknown>).proxy as Record<string, unknown>).mode as string
            ) : (
              "—"
            )}
          </div>
        </div>
        <div className="bg-panel border border-border-solid rounded-xl p-5">
          <div className="text-xs text-gray-500 text-mono mb-1">Audit Log</div>
          <div className="text-sm font-mono text-white">
            {(config as Record<string, unknown>)?.audit ? (
              ((config as Record<string, unknown>).audit as Record<string, unknown>).path as string
            ) : (
              "—"
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
