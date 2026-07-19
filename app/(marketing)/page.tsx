import { Hero } from "@/components/landing/Hero";
import { Features } from "@/components/landing/Features";
import { Stack } from "@/components/landing/Stack";
import { Roadmap } from "@/components/landing/Roadmap";
import { CTA } from "@/components/landing/CTA";

export default function HomePage() {
  return (
    <>
      <Hero />
      <Features />
      <Stack />
      <Roadmap />
      <CTA />
    </>
  );
}
