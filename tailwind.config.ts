import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        void: "#050505",
        surface: "#0a0a0c",
        panel: "#121214",
        border: {
          faint: "#1f1f22",
          solid: "#27272a",
        },
        orange: {
          DEFAULT: "#f97316",
          dim: "rgba(249, 115, 22, 0.1)",
        },
        cyan: {
          DEFAULT: "#06b6d4",
          dim: "rgba(6, 182, 212, 0.1)",
        },
        gray: {
          100: "#ffffff",
          400: "#a1a1aa",
          500: "#52525b",
        },
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "-apple-system", "sans-serif"],
        mono: ["JetBrains Mono", "Fira Code", "monospace"],
      },
      fontSize: {
        "display-xl": ["4.5rem", { lineHeight: "1.1", letterSpacing: "-0.03em" }],
        display: ["3.75rem", { lineHeight: "1.1", letterSpacing: "-0.03em" }],
        hero: ["clamp(3rem, 5vw, 4.5rem)", { lineHeight: "1.1", letterSpacing: "-0.03em" }],
        "section-title": ["3rem", { lineHeight: "1.1", letterSpacing: "-0.03em" }],
        "card-title": ["2rem", { lineHeight: "1.1" }],
      },
      spacing: {
        section: "8rem",
      },
      maxWidth: {
        content: "1280px",
      },
      animation: {
        "fade-in": "fadeIn 0.5s ease-out",
      },
      keyframes: {
        fadeIn: {
          "0%": { opacity: "0" },
          "100%": { opacity: "1" },
        },
      },
    },
  },
  plugins: [],
};

export default config;
