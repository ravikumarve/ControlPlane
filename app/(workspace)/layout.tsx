"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api } from "@/lib/api";

const navItems = [
  { href: "/dashboard", label: "Dashboard", icon: "◆" },
  { href: "/policies", label: "Policies", icon: "⬡" },
  { href: "/audit", label: "Audit Log", icon: "☰" },
  { href: "/settings", label: "Settings", icon: "⚙" },
];

export default function WorkspaceLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const [authenticated, setAuthenticated] = useState(false);
  const [loading, setLoading] = useState(true);

  // Login page doesn't require auth
  const isLoginPage = pathname === "/login";

  useEffect(() => {
    if (isLoginPage) {
      setLoading(false);
      return;
    }
    setAuthenticated(api.isAuthenticated());
    setLoading(false);
  }, [isLoginPage]);

  const handleLogout = () => {
    api.logout();
    router.push("/login");
  };

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center bg-void">
        <div className="text-cyan text-mono animate-pulse">connecting...</div>
      </div>
    );
  }

  if (!isLoginPage && !authenticated) {
    router.push("/login");
    return null;
  }

  return (
    <div className="flex h-screen bg-void text-white">
      {/* Sidebar */}
      <aside className="w-56 shrink-0 border-r border-border-solid bg-surface flex flex-col">
        <div className="px-5 py-6 border-b border-border-solid">
          <a href="/" className="text-lg font-bold tracking-tight">
            <span className="text-cyan">◆</span>{" "}
            <span className="text-white">Control</span>
            <span className="text-orange">Plane</span>
          </a>
          <div className="text-xs text-gray-500 text-mono mt-1">workspace</div>
        </div>

        <nav className="flex-1 px-3 py-4 space-y-1">
          {navItems.map((item) => {
            const active = pathname === item.href;
            return (
              <a
                key={item.href}
                href={item.href}
                className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-colors ${
                  active
                    ? "bg-orange-dim text-orange border border-orange/20"
                    : "text-gray-400 hover:text-white hover:bg-panel"
                }`}
              >
                <span className="text-mono">{item.icon}</span>
                {item.label}
              </a>
            );
          })}
        </nav>

        <div className="px-3 py-4 border-t border-border-solid">
          <button
            onClick={handleLogout}
            className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-gray-500 hover:text-white hover:bg-panel w-full transition-colors"
          >
            <span className="text-mono">↩</span>
            Logout
          </button>
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1 overflow-y-auto bg-surface">
        <div className="max-w-6xl mx-auto px-8 py-8">{children}</div>
      </main>
    </div>
  );
}
