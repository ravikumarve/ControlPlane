"use client";

const API_BASE = process.env.NEXT_PUBLIC_ADMIN_API_URL || "http://localhost:9090";

export interface StatsSnapshot {
  total_calls: number;
  allowed: number;
  blocked: number;
  hitl_pending: number;
  rate_limited: number;
  injection_blocked: number;
}

export interface AuditEntry {
  timestamp: string;
  decision: string;
  tool: string;
  identity: string;
  duration: number;
  reason: string;
  hmac: string;
  prev_hmac: string;
}

export interface Policy {
  name: string;
  action: string;
  match: {
    identity: string;
    tools: string[];
    params?: Record<string, unknown>;
  };
}

class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

function getToken(): string | null {
  return localStorage.getItem("cp_token");
}

async function request<T>(
  method: string,
  path: string,
  body?: unknown
): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  const token = getToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!res.ok) {
    throw new ApiError(`API error: ${res.statusText}`, res.status);
  }

  return res.json();
}

export const api = {
  login: async (apiKey: string) => {
    const data = await request<{ token: string; status: string }>("POST", "/api/login", {
      api_key: apiKey,
    });
    localStorage.setItem("cp_token", data.token);
    return data;
  },

  logout: () => {
    localStorage.removeItem("cp_token");
  },

  isAuthenticated: (): boolean => {
    return !!getToken();
  },

  getStatus: () =>
    request<{ status: string; stats: StatsSnapshot }>("GET", "/api/status"),

  listPolicies: () =>
    request<{ policies: Policy[] }>("GET", "/api/policies"),

  savePolicies: (policies: Policy[]) =>
    request<{ status: string }>("POST", "/api/policies", { policies }),

  getAudit: () =>
    request<{ entries: AuditEntry[] }>("GET", "/api/audit"),

  getConfig: () =>
    request<{ config: unknown }>("GET", "/api/config"),
};
