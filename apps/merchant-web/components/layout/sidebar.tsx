"use client"
import * as React from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { motion } from "framer-motion"
import {
  LayoutDashboard, Receipt, RefreshCcw, Landmark,
  Users, HandCoins, KeySquare, ShieldCheck, LogOut
} from "lucide-react"
import { useLanguage } from "@/components/providers/language-provider"
import { cn } from "@/lib/utils"

export function Sidebar() {
  const pathname = usePathname()
  const { t } = useLanguage()

  const links = [
    { href: "/dashboard",    icon: LayoutDashboard, label: t("Dashboard",     "ዳሽቦርድ") },
    { href: "/payments",     icon: Receipt,          label: t("Payments",      "ክፍያዎች") },
    { href: "/refunds",      icon: RefreshCcw,       label: t("Refunds",       "ተመላሽ ገንዘብ") },
    { href: "/payouts",      icon: Landmark,          label: t("Payouts",       "ወጪ ክፍያዎች") },
    { href: "/payroll",      icon: Users,             label: t("Payroll",       "የደሞዝ ክፍያ") },
    { href: "/subscriptions",icon: HandCoins,         label: t("Subscriptions", "ምዝገባዎች") },
    { href: "/developers",   icon: KeySquare,         label: t("Developers",    "አልሚዎች") },
    { href: "/compliance",   icon: ShieldCheck,       label: t("Compliance",    "ተገዢነት") },
  ]

  return (
    <aside className="fixed left-0 top-0 h-screen w-[260px] z-40 hidden md:flex flex-col justify-between pb-6
      border-r border-border
      bg-[hsl(var(--background-card))]/90 backdrop-blur-2xl
      transition-colors duration-300">

      {/* Logo */}
      <div>
        <div className="h-16 flex items-center px-6 border-b border-border">
          <Link
            href="/dashboard"
            className="font-bold text-xl tracking-tight bg-gradient-to-br from-primary to-emerald-500 bg-clip-text text-transparent"
          >
            ApexPay
          </Link>
        </div>

        {/* Nav */}
        <nav className="p-4 space-y-1">
          {links.map((link) => {
            const active = pathname.startsWith(link.href)
            return (
              <Link key={link.href} href={link.href} className="block relative">
                {active && (
                  <motion.div
                    layoutId="sidebar-active"
                    className="absolute inset-0 bg-primary/15 dark:bg-primary/20 rounded-xl"
                    transition={{ type: "spring", stiffness: 300, damping: 30 }}
                  />
                )}
                <div className={cn(
                  "relative flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-colors",
                  active
                    ? "text-primary"
                    : "text-[hsl(var(--foreground-muted))] hover:text-[hsl(var(--foreground))] hover:bg-[hsl(var(--surface-2))]"
                )}>
                  <link.icon
                    size={18}
                    className={cn(active ? "text-primary" : "text-[hsl(var(--foreground-muted))]")}
                  />
                  {link.label}
                </div>
              </Link>
            )
          })}
        </nav>
      </div>

      {/* Bottom user card */}
      <div className="px-4">
        <div className="rounded-2xl border border-border bg-[hsl(var(--surface-2))] p-4 space-y-3 transition-colors duration-300">
          <div className="flex items-center gap-3">
            <div className="h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold text-xs">
              MT
            </div>
            <div className="text-sm">
              <p className="font-semibold leading-none text-[hsl(var(--foreground))]">Apex Trading</p>
              <p className="text-xs text-[hsl(var(--foreground-muted))] mt-1">Merchant ID: 0092</p>
            </div>
          </div>
          <Link href="/">
            <button className="w-full flex items-center justify-center gap-2 text-xs font-medium text-red-500 dark:text-red-400 hover:bg-red-500/10 py-2 rounded-xl transition-colors">
              <LogOut size={14} /> {t("Sign Out", "ውጣ")}
            </button>
          </Link>
        </div>
      </div>
    </aside>
  )
}
