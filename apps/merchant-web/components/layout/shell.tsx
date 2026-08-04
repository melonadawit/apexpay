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

  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-[hsl(var(--background))] transition-colors duration-300">
      {!hideSidebar && <Sidebar />}

      <div className={`flex-1 flex flex-col ${!hideSidebar ? "md:ml-[260px]" : ""}`}>
        {!hideSidebar && <Header />}

        <main className="flex-1 overflow-x-hidden">
          {children}
        </main>
      </div>

      <Toaster
        position="top-right"
        richColors
        toastOptions={{
          classNames: {
            toast: "dark:!bg-[#1a1a1a] dark:!border-white/10 dark:!text-white",
          },
        }}
      />
    </div>
  )
}
