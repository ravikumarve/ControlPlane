"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";

export default function LoginPage() {
  const router = useRouter();
  const [apiKey, setApiKey] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await api.login(apiKey);
      router.push("/dashboard");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "login failed";
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-void">
      <div className="w-full max-w-sm">
        <div className="text-center mb-8">
          <div className="text-3xl font-bold tracking-tight">
            <span className="text-cyan">◆</span>{" "}
            <span className="text-white">Control</span>
            <span className="text-orange">Plane</span>
          </div>
          <p className="text-gray-500 text-mono text-sm mt-2">workspace login</p>
        </div>

        <form
          onSubmit={handleSubmit}
          className="bg-panel border border-border-solid rounded-xl p-6 space-y-4"
        >
          <div className="space-y-2">
            <label
              htmlFor="api-key"
              className="text-sm text-gray-400 text-mono block"
            >
              API Key
            </label>
            <input
              id="api-key"
              type="password"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              placeholder="Enter your API key"
              className="w-full bg-void border border-border-solid rounded-lg px-4 py-2.5 text-white font-mono text-sm placeholder:text-gray-500 focus:outline-none focus:border-cyan/50 transition-colors"
              autoFocus
            />
          </div>

          {error && (
            <div className="text-red-400 text-sm text-mono bg-red-950/30 border border-red-900/30 rounded-lg px-3 py-2">
              ✕ {error}
            </div>
          )}

          <button
            type="submit"
            disabled={loading || !apiKey}
            className="w-full bg-orange hover:bg-orange/90 text-white font-medium rounded-lg px-5 py-2.5 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "authenticating..." : "Sign in"}
          </button>

          <p className="text-xs text-gray-500 text-center">
            Default key: <code className="text-cyan text-mono">controlplane-dev-key</code>
          </p>
        </form>
      </div>
    </div>
  );
}
