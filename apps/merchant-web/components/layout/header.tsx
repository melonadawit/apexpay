"use client"
import * as React from "react"
import { useRouter } from "next/navigation"
import { useLanguage } from "@/components/providers/language-provider"
import { ThemeToggle } from "@/components/providers/theme-provider"
import { logout } from "@/lib/api/auth"
import { Globe, Bell, Menu, LogOut } from "lucide-react"

export function Header({ onToggleSidebar }: { onToggleSidebar?: () => void }) {
  const { language, setLanguage } = useLanguage()
  const router = useRouter()
  const [signingOut, setSigningOut] = React.useState(false)

  const signOut = async () => {
    setSigningOut(true)
    try {
      await logout()
    } finally {
      router.push("/login")
      router.refresh()
    }
  }

  return (
    <header className="h-16 border-b border-border sticky top-0 z-30 flex items-center justify-between px-6
      bg-card/80 backdrop-blur-xl
      transition-colors duration-300">
      {/* Left: env badge & toggle */}
      <div className="flex items-center gap-4">
        {onToggleSidebar && (
          <button onClick={onToggleSidebar} className="text-muted-foreground hover:text-foreground transition-colors">
            <Menu size={20} />
          </button>
        )}
        <div className="flex items-center gap-2">
          <div className="h-6 w-6 rounded-md bg-accent-gold flex items-center justify-center text-foreground text-[10px] font-bold">
            LIVE
          </div>
          <span className="text-sm font-medium text-muted-foreground hidden sm:inline">production-et</span>
        </div>
      </div>

      {/* Right: controls */}
      <div className="flex items-center gap-2">
        {/* Language toggle */}
        <button
          onClick={() => setLanguage(language === "en" ? "am" : "en")}
          className="flex items-center gap-1.5 text-sm font-medium px-3 py-1.5 rounded-full 
            transition-colors duration-200
            dark:hover:bg-card/10 dark:text-foreground/70
            hover:bg-background/5 text-muted-foreground"
        >
          <Globe size={15} />
          {language === "en" ? "EN" : "አማ"}
        </button>

        {/* Theme toggle */}
        <ThemeToggle />

        {/* Notifications */}
        <button className="h-9 w-9 rounded-full flex items-center justify-center transition-colors duration-200
          dark:hover:bg-card/10 dark:text-foreground/70
          hover:bg-background/5 text-muted-foreground relative">
          <Bell size={17} />
          <span className="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-red-500 border-2
            border-card" />
        </button>

        {/* Sign out */}
        <button
          onClick={signOut}
          disabled={signingOut}
          title="Sign out"
          className="h-9 px-3 rounded-full flex items-center gap-1.5 text-sm font-medium transition-colors
            text-muted-foreground hover:bg-background/5 hover:text-foreground disabled:opacity-50"
        >
          <LogOut size={15} />
          <span className="hidden sm:inline">{signingOut ? "…" : "Sign out"}</span>
        </button>
      </div>
    </header>
  )
}
