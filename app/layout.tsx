import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: {
    default: "ControlPlane AI — Security infrastructure for the AI agent ecosystem",
    template: "%s | ControlPlane AI",
  },
  description:
    "ControlPlane AI builds zero-trust security tools for organizations deploying autonomous AI agents. MCP Security Gateway, RBAC, injection detection, and audit.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="min-h-screen">{children}</body>
    </html>
  );
}
