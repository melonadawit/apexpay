"use client"
import * as React from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { motion } from "framer-motion"
import {
  LayoutDashboard, Receipt, RefreshCcw, Landmark,
  Users, HandCoins, KeySquare, ShieldCheck, LogOut,
  Building2, ShieldAlert, Wallet, FileText, UserCog, TrendingUp, Bell, Settings, BookOpenCheck
} from "lucide-react"
import { useLanguage } from "@/components/providers/language-provider"
import { cn } from "@/lib/utils"

export function Sidebar({ isCollapsed }: { isCollapsed?: boolean }) {
  const pathname = usePathname()
  const { t } = useLanguage()

  const links = [
    { href: "/dashboard",    icon: LayoutDashboard, label: t("Dashboard",     "ዳሽቦርድ") },
    { href: "/payments",     icon: Receipt,          label: t("Payments",      "ክፍያዎች") },
    { href: "/refunds",      icon: RefreshCcw,       label: t("Refunds",       "ተመላሽ ገንዘብ") },
    { href: "/payouts",      icon: Landmark,          label: t("Payouts",       "ወጪ ክፍያዎች") },
    { href: "/payroll",      icon: Users,             label: t("Payroll",       "የደሞዝ ክፍያ") },
    { href: "/hris",         icon: Building2,         label: t("Workforce",     "የሰው ሃይል") },
    { href: "/subscriptions",icon: HandCoins,         label: t("Subscriptions", "ምዝገባዎች") },
    { href: "/banking",      icon: Wallet,            label: t("Banking",       "ባንክ") },
    { href: "/treasury",     icon: TrendingUp,        label: t("Treasury",      "ግምጃ") },
    { href: "/invoices",     icon: FileText,          label: t("Invoices",      "ኢንቮይስ") },
    { href: "/risk",         icon: ShieldAlert,       label: t("Risk & Fraud",  "አደጋ") },
    { href: "/analytics",    icon: TrendingUp,        label: t("Analytics",     "ትንታኔ") },
    { href: "/accounting",   icon: BookOpenCheck,     label: t("Accounting",    "ሂሳብ") },
    { href: "/fixed-assets", icon: Building2,         label: t("Fixed Assets",  "ቋሚ ንብረት") },
    { href: "/compliance-console", icon: ShieldCheck, label: t("Compliance",   "ተገዢነት") },
    { href: "/team",         icon: UserCog,           label: t("Team",          "ቡድን") },
    { href: "/settings/notifications", icon: Bell,    label: t("Notifications", "ማሳወቂያ") },
    { href: "/settings/2fa", icon: ShieldCheck,       label: t("2FA Security",  "ደህንነት") },
    { href: "/developers",   icon: KeySquare,         label: t("Developers",    "አልሚዎች") },
  ]

  return (
    <aside className={cn(
      "fixed left-0 top-0 h-screen z-40 hidden md:flex flex-col justify-between pb-6",
      "border-r border-border bg-card/90 backdrop-blur-2xl transition-all duration-300",
      isCollapsed ? "w-[80px]" : "w-[260px]"
    )}>

      {/* Logo */}
      <div>
        <div className={cn("h-16 flex items-center border-b border-border", isCollapsed ? "justify-center px-0" : "px-6")}>
          <Link
            href="/dashboard"
            className="font-bold text-xl tracking-tight bg-gradient-to-br from-primary to-emerald-500 bg-clip-text text-transparent"
          >
            {isCollapsed ? "A" : "ApexPay"}
          </Link>
        </div>

        {/* Nav */}
        <nav className={cn("p-4 space-y-1", isCollapsed ? "px-2" : "p-4")}>
          {links.map((link) => {
            const active = pathname.startsWith(link.href)
            return (
              <Link key={link.href} href={link.href} className="block relative" title={isCollapsed ? link.label : undefined}>
                {active && (
                  <motion.div
                    layoutId="sidebar-active"
                    className="absolute inset-0 bg-primary/15 dark:bg-primary/20 rounded-xl"
                    transition={{ type: "spring", stiffness: 300, damping: 30 }}
                  />
                )}
                <div className={cn(
                  "relative flex items-center rounded-xl text-sm font-medium transition-colors",
                  isCollapsed ? "justify-center py-3" : "gap-3 px-3 py-2.5",
                  active
                    ? "text-primary"
                    : "text-muted-foreground hover:text-foreground hover:bg-muted"
                )}>
                  <link.icon
                    size={18}
                    className={cn(active ? "text-primary" : "text-muted-foreground")}
                  />
                  {!isCollapsed && <span>{link.label}</span>}
                </div>
              </Link>
            )
          })}
        </nav>
      </div>

      {/* Bottom user card */}
      <div className={cn(isCollapsed ? "px-2" : "px-4")}>
        <div className={cn("rounded-2xl border border-border bg-muted/50 transition-colors duration-300", isCollapsed ? "p-2 space-y-2 flex flex-col items-center" : "p-4 space-y-3")}>
          <div className={cn("flex items-center", isCollapsed ? "justify-center" : "gap-3")}>
            <div className="h-8 w-8 rounded-full bg-primary/20 flex shrink-0 items-center justify-center text-primary font-bold text-xs">
              MT
            </div>
            {!isCollapsed && (
              <div className="text-sm overflow-hidden">
                <p className="font-semibold leading-none text-foreground truncate">Apex Trading</p>
                <p className="text-xs text-muted-foreground mt-1 truncate">ID: 0092</p>
              </div>
            )}
          </div>
          <Link href="/" className="w-full">
            <button className={cn("flex items-center justify-center text-xs font-medium text-red-500 dark:text-red-400 hover:bg-red-500/10 transition-colors", isCollapsed ? "p-2 w-full rounded-xl" : "gap-2 py-2 w-full rounded-xl")}>
              <LogOut size={14} /> {!isCollapsed && <span>{t("Sign Out", "ውጣ")}</span>}
            </button>
          </Link>
        </div>
      </div>
    </aside>
  )
}

