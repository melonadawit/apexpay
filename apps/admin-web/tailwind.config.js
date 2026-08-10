/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}", "./lib/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: { DEFAULT: "#0B6E4F", light: "#10A37A", dark: "#094E38", foreground: "#FFFFFF" },
        accent: { DEFAULT: "#EAB308", gold: "#EAB308" },
      },
      boxShadow: {
        soft: "0 10px 30px rgba(0,0,0,0.07)",
        medium: "0 20px 40px rgba(0,0,0,0.10)",
      },
    },
  },
  plugins: [],
}
