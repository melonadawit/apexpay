"use client"
import * as React from "react"
import { useLanguage } from "@/components/providers/language-provider"
import { ThemeToggle } from "@/components/providers/theme-provider"
import { Globe, Bell } from "lucide-react"

export function Header() {
  const { language, setLanguage } = useLanguage()

  return (
    <header className="h-16 border-b border-border sticky top-0 z-30 flex items-center justify-between px-6
      bg-[hsl(var(--background-card))]/80 backdrop-blur-xl
      transition-colors duration-300">
      {/* Left: env badge */}
      <div className="flex items-center gap-2">
        <div className="h-6 w-6 rounded-md bg-accent-gold flex items-center justify-center text-black text-[10px] font-bold">
          LIVE
        </div>
        <span className="text-sm font-medium text-[hsl(var(--foreground-muted))]">production-et</span>
      </div>

      {/* Right: controls */}
      <div className="flex items-center gap-2">
        {/* Language toggle */}
        <button
          onClick={() => setLanguage(language === "en" ? "am" : "en")}
          className="flex items-center gap-1.5 text-sm font-medium px-3 py-1.5 rounded-full 
            transition-colors duration-200
            dark:hover:bg-white/10 dark:text-white/70
            hover:bg-black/5 text-[hsl(var(--foreground-muted))]"
        >
          <Globe size={15} />
          {language === "en" ? "EN" : "አማ"}
        </button>

        {/* Theme toggle */}
        <ThemeToggle />

        {/* Notifications */}
        <button className="h-9 w-9 rounded-full flex items-center justify-center transition-colors duration-200
          dark:hover:bg-white/10 dark:text-white/70
          hover:bg-black/5 text-[hsl(var(--foreground-muted))] relative">
          <Bell size={17} />
          <span className="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-red-500 border-2
            border-[hsl(var(--background-card))]" />
        </button>
      </div>
    </header>
  )
}
