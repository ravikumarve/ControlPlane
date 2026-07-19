import { ThreatGrid } from "@/components/shared/ThreatGrid";
import { Navbar } from "@/components/shared/Navbar";
import { Footer } from "@/components/shared/Footer";

export default function MarketingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <ThreatGrid />
      <div className="absolute left-0 right-0 top-0 z-0 h-[70vh] w-[70vw] -translate-y-20 translate-x-[-10vw] pointer-events-none"
        style={{
          background: "radial-gradient(ellipse at center, rgba(6,182,212,0.05) 0%, transparent 60%)",
          filter: "blur(80px)",
        }}
      />
      <div className="relative z-10 mx-auto max-w-content px-4 sm:px-6 lg:px-8">
        <Navbar />
        <main>{children}</main>
        <Footer />
      </div>
    </>
  );
}
