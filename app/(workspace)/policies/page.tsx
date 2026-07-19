"use client";

import { useEffect, useState } from "react";
import { api, Policy } from "@/lib/api";

export default function PoliciesPage() {
  const [policies, setPolicies] = useState<Policy[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const [saveMsg, setSaveMsg] = useState("");
  const [editing, setEditing] = useState(false);
  const [editPolicies, setEditPolicies] = useState<Policy[]>([]);

  const fetchPolicies = async () => {
    try {
      const data = await api.listPolicies();
      setPolicies(data.policies);
      setEditPolicies(JSON.parse(JSON.stringify(data.policies)));
      setError("");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "failed to fetch policies";
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPolicies();
  }, []);

  const handleSave = async () => {
    setSaving(true);
    setSaveMsg("");
    try {
      await api.savePolicies(editPolicies);
      setPolicies(JSON.parse(JSON.stringify(editPolicies)));
      setEditing(false);
      setSaveMsg("policies saved successfully");
      setTimeout(() => setSaveMsg(""), 3000);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "failed to save";
      setSaveMsg(`error: ${msg}`);
    } finally {
      setSaving(false);
    }
  };

  const addPolicy = () => {
    setEditPolicies([
      ...editPolicies,
      {
        name: `policy-${editPolicies.length + 1}`,
        action: "allow",
        match: { identity: "*", tools: ["*"] },
      },
    ]);
  };

  const removePolicy = (index: number) => {
    setEditPolicies(editPolicies.filter((_, i) => i !== index));
  };

  const updatePolicy = (
    index: number,
    field: string,
    value: string | string[]
  ) => {
    const updated = [...editPolicies];
    const policy = updated[index];
    if (field === "name") {
      updated[index] = { ...policy, name: value as string };
    } else if (field === "action") {
      updated[index] = { ...policy, action: value as string };
    } else if (field === "identity") {
      updated[index] = {
        ...policy,
        match: { ...policy.match, identity: value as string },
      };
    } else if (field === "tools") {
      updated[index] = {
        ...policy,
        match: { ...policy.match, tools: value as string[] },
      };
    }
    setEditPolicies(updated);
  };

  if (loading) {
    return (
      <div className="text-center py-20 text-gray-500 text-mono animate-pulse">
        loading policies...
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Policies</h1>
          <p className="text-gray-500 text-mono text-sm mt-1">
            manage access control rules
          </p>
        </div>
        <div className="flex items-center gap-3">
          {!editing ? (
            <>
              <button
                onClick={() => setEditing(true)}
                className="px-4 py-2 text-sm bg-cyan/10 border border-cyan/20 text-cyan rounded-lg hover:bg-cyan/20 transition-colors"
              >
                edit policies
              </button>
            </>
          ) : (
            <>
              <button
                onClick={addPolicy}
                className="px-4 py-2 text-sm border border-border-solid text-gray-400 rounded-lg hover:text-white hover:border-gray-500 transition-colors"
              >
                + add policy
              </button>
              <button
                onClick={() => {
                  setEditing(false);
                  setEditPolicies(JSON.parse(JSON.stringify(policies)));
                }}
                className="px-4 py-2 text-sm border border-border-solid text-gray-400 rounded-lg hover:text-white hover:border-gray-500 transition-colors"
              >
                cancel
              </button>
              <button
                onClick={handleSave}
                disabled={saving}
                className="px-4 py-2 text-sm bg-orange text-white rounded-lg hover:bg-orange/90 transition-colors disabled:opacity-50"
              >
                {saving ? "saving..." : "save"}
              </button>
            </>
          )}
        </div>
      </div>

      {/* Messages */}
      {error && (
        <div className="text-red-400 text-sm text-mono bg-red-950/30 border border-red-900/30 rounded-lg px-4 py-3">
          ⚠ {error}
        </div>
      )}
      {saveMsg && (
        <div
          className={`text-sm text-mono rounded-lg px-4 py-3 ${
            saveMsg.startsWith("error")
              ? "text-red-400 bg-red-950/30 border border-red-900/30"
              : "text-green-400 bg-green-950/30 border border-green-900/30"
          }`}
        >
          {saveMsg.startsWith("error") ? "⚠ " : "✓ "}
          {saveMsg}
        </div>
      )}

      {/* Policy list */}
      <div className="space-y-3">
        {(editing ? editPolicies : policies).length === 0 && (
          <div className="text-center py-16 text-gray-500 text-mono">
            no policies configured
          </div>
        )}

        {(editing ? editPolicies : policies).map((policy, i) => (
          <div
            key={i}
            className={`rounded-xl border p-5 ${
              editing
                ? "border-orange/20 bg-orange-dim"
                : "border-border-solid bg-panel"
            }`}
          >
            <div className="flex items-start justify-between gap-4">
              <div className="flex-1 space-y-3">
                {/* Name + Action */}
                <div className="flex items-center gap-3 flex-wrap">
                  {editing ? (
                    <>
                      <input
                        value={policy.name}
                        onChange={(e) => updatePolicy(i, "name", e.target.value)}
                        className="bg-void border border-border-solid rounded-lg px-3 py-1.5 text-sm font-mono text-white w-48 focus:outline-none focus:border-cyan/50"
                      />
                      <select
                        value={policy.action}
                        onChange={(e) => updatePolicy(i, "action", e.target.value)}
                        className="bg-void border border-border-solid rounded-lg px-3 py-1.5 text-sm font-mono text-white focus:outline-none focus:border-cyan/50"
                      >
                        <option value="allow">allow</option>
                        <option value="block">block</option>
                        <option value="hitl">hitl</option>
                      </select>
                    </>
                  ) : (
                    <>
                      <span className="text-sm font-mono text-white">
                        {policy.name}
                      </span>
                      <span
                        className={`text-xs font-mono px-2 py-0.5 rounded-full ${
                          policy.action === "allow"
                            ? "bg-green-950/30 text-green-400 border border-green-900/30"
                            : policy.action === "block"
                            ? "bg-red-950/30 text-red-400 border border-red-900/30"
                            : "bg-yellow-950/30 text-yellow-400 border border-yellow-900/30"
                        }`}
                      >
                        {policy.action}
                      </span>
                    </>
                  )}
                </div>

                {/* Match: Identity */}
                <div className="flex items-center gap-2 text-sm">
                  <span className="text-gray-500 text-mono">identity:</span>
                  {editing ? (
                    <input
                      value={policy.match.identity}
                      onChange={(e) => updatePolicy(i, "identity", e.target.value)}
                      className="bg-void border border-border-solid rounded-lg px-3 py-1 text-xs font-mono text-cyan w-48 focus:outline-none focus:border-cyan/50"
                    />
                  ) : (
                    <span className="font-mono text-xs text-cyan">
                      {policy.match.identity}
                    </span>
                  )}
                </div>

                {/* Match: Tools */}
                <div className="flex items-center gap-2 text-sm">
                  <span className="text-gray-500 text-mono">tools:</span>
                  {editing ? (
                    <input
                      value={policy.match.tools.join(", ")}
                      onChange={(e) =>
                        updatePolicy(
                          i,
                          "tools",
                          e.target.value.split(",").map((s) => s.trim())
                        )
                      }
                      className="bg-void border border-border-solid rounded-lg px-3 py-1 text-xs font-mono text-cyan flex-1 focus:outline-none focus:border-cyan/50"
                      placeholder="tool1, tool2, *"
                    />
                  ) : (
                    <span className="font-mono text-xs text-cyan">
                      {policy.match.tools.join(", ")}
                    </span>
                  )}
                </div>
              </div>

              {/* Delete button (edit mode only) */}
              {editing && (
                <button
                  onClick={() => removePolicy(i)}
                  className="text-gray-500 hover:text-red-400 transition-colors p-1"
                  title="remove policy"
                >
                  ✕
                </button>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
