import Link from "next/link";

export function Footer() {
  return (
    <footer className="border-t border-surface-border bg-surface-dark text-gray-400">
      <div className="container-content py-12">
        <div className="grid gap-8 sm:grid-cols-3">
          <div>
            <span className="font-bold text-white">ControlPlane AI</span>
            <p className="mt-2 text-sm">
              Security infrastructure for the AI agent ecosystem.
            </p>
          </div>

          <div>
            <h3 className="mb-3 text-sm font-semibold text-white">Product</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link href="/pricing" className="hover:text-white transition-colors">
                  Pricing
                </Link>
              </li>
              <li>
                <Link href="/docs" className="hover:text-white transition-colors">
                  Documentation
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <h3 className="mb-3 text-sm font-semibold text-white">Company</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link href="/security" className="hover:text-white transition-colors">
                  Security
                </Link>
              </li>
              <li>
                <a
                  href="https://github.com/ravikumarve/ControlPlane"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="hover:text-white transition-colors"
                >
                  GitHub
                </a>
              </li>
            </ul>
          </div>
        </div>

        <div className="mt-8 border-t border-gray-800 pt-8 text-center text-xs">
          &copy; {new Date().getFullYear()} ControlPlane AI. MIT licensed.
        </div>
      </div>
    </footer>
  );
}
