"use client"
import * as React from "react"
import { Moon, Sun } from "lucide-react"

type Theme = "dark" | "light"

interface ThemeContextType {
  theme: Theme
  setTheme: (t: Theme) => void
  toggleTheme: () => void
}

const ThemeContext = React.createContext<ThemeContextType>({
  theme: "dark",
  setTheme: () => {},
  toggleTheme: () => {},
})

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = React.useState<Theme>("dark")
  const [mounted, setMounted] = React.useState(false)

  // Read saved preference after mount (avoids SSR mismatch)
  React.useEffect(() => {
    const saved = localStorage.getItem("apexpay_theme") as Theme
    const preferred = saved || (window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light")
    setThemeState(preferred)
    setMounted(true)
  }, [])

  // Apply class to <html> element
  React.useEffect(() => {
    if (!mounted) return
    const root = document.documentElement
    root.classList.remove("dark", "light")
    root.classList.add(theme)
    localStorage.setItem("apexpay_theme", theme)
  }, [theme, mounted])

  const setTheme = (t: Theme) => setThemeState(t)
  const toggleTheme = () => setThemeState(prev => (prev === "dark" ? "light" : "dark"))

  // Prevent flash: render children immediately, class applied via effect
  return (
    <ThemeContext.Provider value={{ theme, setTheme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  )
}

export const useTheme = () => React.useContext(ThemeContext)

/** Drop-in toggle button — use anywhere */
export function ThemeToggle({ className }: { className?: string }) {
  const { theme, toggleTheme } = useTheme()
  return (
    <button
      onClick={toggleTheme}
      aria-label="Toggle theme"
      className={`relative h-9 w-9 flex items-center justify-center rounded-full 
        transition-all duration-300 hover:scale-105 active:scale-95
        dark:bg-white/10 dark:hover:bg-white/20 dark:text-white
        bg-black/5 hover:bg-black/10 text-neutral-700
        ${className ?? ""}`}
    >
      <Sun
        size={17}
        className="absolute transition-all duration-300 dark:opacity-0 dark:rotate-90 dark:scale-0 opacity-100 rotate-0 scale-100"
      />
      <Moon
        size={17}
        className="absolute transition-all duration-300 dark:opacity-100 dark:rotate-0 dark:scale-100 opacity-0 -rotate-90 scale-0"
      />
    </button>
  )
}
