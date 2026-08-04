/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ["class"],
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}",
  ],
  theme: {
    container: { center: true, padding: "2rem", screens: { "2xl": "1400px" } },
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        primary: {
          DEFAULT: "#0B6E4F", // ET Green - outstanding
          light: "#10A37A",
          dark: "#094E38",
          50: "#ECFDF5",
          foreground: "#FFFFFF",
        },
        accent: {
          gold: "#EAB308",
          yellow: "#FEF08A",
          DEFAULT: "#EAB308",
        },
        neutral: {
          50: "#FAFAFA",
          100: "#F4F4F5",
          200: "#E4E4E7",
          800: "#27272A",
          900: "#18181B",
        },
      },
      borderRadius: {
        lg: "16px",
        xl: "24px",
        "2xl": "32px",
      },
      fontFamily: {
        sans: ["Inter", "Noto Sans Ethiopic", "system-ui", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"],
      },
      keyframes: {
        "accordion-down": { from: { height: "0" }, to: { height: "var(--radix-accordion-content-height)" } },
        "accordion-up": { from: { height: "var(--radix-accordion-content-height)" }, to: { height: "0" } },
        shimmer: { "0%": { transform: "translateX(-100%)" }, "100%": { transform: "translateX(100%)" } },
        confetti: { "0%": { transform: "translateY(-100vh) rotate(0deg)", opacity: "1" }, "100%": { transform: "translateY(100vh) rotate(720deg)", opacity: "0" } },
        pulseGlow: { "0%,100%": { boxShadow: "0 0 20px rgba(11,110,79,0.2)" }, "50%": { boxShadow: "0 0 30px rgba(11,110,79,0.4)" } },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
        shimmer: "shimmer 2s infinite",
        confetti: "confetti 3s ease-out forwards",
        pulseGlow: "pulseGlow 2s ease-in-out infinite",
      },
      boxShadow: {
        soft: "0 10px 30px rgba(0,0,0,0.07)",
        medium: "0 20px 40px rgba(0,0,0,0.10)",
        glass: "0 8px 32px rgba(0,0,0,0.08), inset 0 1px 0 rgba(255,255,255,0.5)",
      },
      backdropBlur: { xs: "2px", xl: "20px", "2xl": "40px" },
    },
  },
  plugins: [require("tailwindcss-animate")],
}
