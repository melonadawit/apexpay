"use client"
import * as React from "react"
import { usePathname } from "next/navigation"
import { Sidebar } from "./sidebar"
import { Header } from "./header"
import { Toaster } from "sonner"

export function Shell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const isLandingPage = pathname === "/"
  const isOnboarding = pathname === "/onboarding"
  const hideSidebar = isLandingPage || isOnboarding
  const [isCollapsed, setIsCollapsed] = React.useState(false)

  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-[hsl(var(--background))] transition-colors duration-300">
      {!hideSidebar && <Sidebar isCollapsed={isCollapsed} />}

      <div className={`flex-1 flex flex-col transition-all duration-300 ${!hideSidebar ? (isCollapsed ? "md:ml-[80px]" : "md:ml-[260px]") : ""}`}>
        {!hideSidebar && <Header onToggleSidebar={() => setIsCollapsed(!isCollapsed)} />}

        <main className="flex-1 overflow-x-hidden">
          {children}
        </main>
      </div>

      <Toaster
        position="top-right"
        richColors
        toastOptions={{
          classNames: {
            toast: "dark:!bg-[#1a1a1a] dark:!border-border dark:!text-foreground",
          },
        }}
      />
    </div>
  )
}
