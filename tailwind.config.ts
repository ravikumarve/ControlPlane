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
        primary: {
          DEFAULT: "#1A237E",
          50: "#E8EAF6",
          100: "#C5CAE9",
          200: "#9FA8DA",
          300: "#7986CB",
          400: "#5C6BC0",
          500: "#3F51B5",
          600: "#3949AB",
          700: "#303F9F",
          800: "#283593",
          900: "#1A237E",
        },
        accent: {
          DEFAULT: "#00BCD4",
          50: "#E0F7FA",
          100: "#B2EBF2",
          200: "#80DEEA",
          300: "#4DD0E1",
          400: "#26C6DA",
          500: "#00BCD4",
          600: "#00ACC1",
          700: "#0097A7",
          800: "#00838F",
          900: "#006064",
        },
        success: {
          DEFAULT: "#4CAF50",
          50: "#E8F5E9",
          500: "#4CAF50",
          900: "#1B5E20",
        },
        danger: {
          DEFAULT: "#F44336",
          50: "#FFEBEE",
          500: "#F44336",
          900: "#B71C1C",
        },
        warning: {
          DEFAULT: "#FFC107",
          50: "#FFF8E1",
          500: "#FFC107",
          900: "#FF6F00",
        },
        surface: {
          dark: "#0D1117",
          light: "#FAFAFA",
          card: "#161B22",
          border: "#30363D",
        },
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "-apple-system", "sans-serif"],
        mono: ["JetBrains Mono", "Fira Code", "monospace"],
      },
      fontSize: {
        "display-xl": ["4.5rem", { lineHeight: "1.1", letterSpacing: "-0.02em" }],
        display: ["3.75rem", { lineHeight: "1.1", letterSpacing: "-0.02em" }],
        hero: ["3rem", { lineHeight: "1.15", letterSpacing: "-0.01em" }],
        "section-title": ["2.25rem", { lineHeight: "1.2", letterSpacing: "-0.01em" }],
        "card-title": ["1.25rem", { lineHeight: "1.4" }],
      },
      spacing: {
        section: "5rem",
        "section-lg": "8rem",
      },
      maxWidth: {
        content: "72rem",
      },
      animation: {
        "fade-in": "fadeIn 0.5s ease-out",
        "slide-up": "slideUp 0.5s ease-out",
      },
      keyframes: {
        fadeIn: {
          "0%": { opacity: "0" },
          "100%": { opacity: "1" },
        },
        slideUp: {
          "0%": { opacity: "0", transform: "translateY(20px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
      },
    },
  },
  plugins: [],
};

export default config;
