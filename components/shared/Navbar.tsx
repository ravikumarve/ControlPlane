import Link from "next/link";

export function Navbar() {
  return (
    <header className="sticky top-0 z-50 border-b border-surface-border bg-white/80 backdrop-blur-sm dark:bg-surface-dark/80">
      <nav className="container-content flex h-16 items-center justify-between">
        <Link
          href="/"
          className="text-lg font-bold text-primary"
        >
          ControlPlane AI
        </Link>

        <div className="flex items-center gap-6">
          <Link
            href="/pricing"
            className="text-sm text-gray-600 hover:text-primary dark:text-gray-400 dark:hover:text-accent transition-colors"
          >
            Pricing
          </Link>
          <Link
            href="/docs"
            className="text-sm text-gray-600 hover:text-primary dark:text-gray-400 dark:hover:text-accent transition-colors"
          >
            Docs
          </Link>
          <Link
            href="/security"
            className="text-sm text-gray-600 hover:text-primary dark:text-gray-400 dark:hover:text-accent transition-colors"
          >
            Security
          </Link>
        </div>
      </nav>
    </header>
  );
}
