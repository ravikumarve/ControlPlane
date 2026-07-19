import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: {
    default: "ControlPlane | Security for the AI Agent Ecosystem",
    template: "%s | ControlPlane",
  },
  description:
    "ControlPlane builds zero-trust security infrastructure for organizations deploying autonomous AI agents. MCP Security Gateway, RBAC, injection detection, and audit.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
