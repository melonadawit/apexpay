"use client"
import * as React from "react"
import { motion } from "framer-motion"
import { Card, GlassCard } from "@/components/ui/card"
import { DonutProgress } from "@/components/ui/progress"
import { TPVRecharts, HealthRecharts } from "./recharts"
import { api } from "@/lib/api/client"
import { useData, formatETB } from "@/lib/api/use-data"
import { useRequireAuth } from "@/lib/api/require-auth"

export default function DashboardPage() {
  const { checking } = useRequireAuth()
  const { data: summary, loading: summaryLoading } = useData(() => api.summary(), [])
  const { data: payments } = useData(() => api.payments(5), [])

  if (checking) {
    return <div className="min-h-screen flex items-center justify-center text-sm text-muted-foreground">Checking session…</div>
  }

  const tpvToday = formatETB(summary?.tpv_today)
  const successRate = summary ? (summary.success_rate_7_days * 100).toFixed(1) : "—"
  const activeLinks = summary?.active_links ?? "—"
  const recent = (payments ?? []).slice(0, 5)

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50 p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold">Dashboard • ዳሽቦርድ</h1>
            <p className="text-sm text-muted-foreground">Welcome • live from ApexPay API</p>
          </div>
          <DonutProgress value={summary ? Math.round(successRateAsNum(summary)) : 78} />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="rounded-2xl bg-gradient-to-br from-primary to-primary-light p-6 text-foreground shadow-medium"
          >
            <p className="text-foreground/70 text-sm">TPV Today • ዛሬ</p>
            <p className="text-3xl font-bold mt-2">
              {summaryLoading ? "…" : `ETB ${tpvToday}`}
            </p>
            <p className="text-xs mt-2 bg-card/20 inline-block px-2 py-0.5 rounded-full">
              7-day TPV: ETB {formatETB(summary?.tpv_7_days)}
            </p>
            <div className="mt-4">
              <TPVRecharts />
            </div>
          </motion.div>

          <Card className="p-6">
            <p className="text-sm text-muted-foreground">Success Rate (7d) • ስኬት</p>
            <p className="text-2xl font-bold mt-2">{summaryLoading ? "…" : `${successRate}%`}</p>
            <p className="text-xs text-green-600 mt-1">
              {summary
                ? `${summary.success_count_7_days} succeeded / ${summary.total_count_7_days} total`
                : ""}
            </p>
            <div className="mt-2">
              <HealthRecharts />
            </div>
          </Card>

          <Card className="p-6">
            <p className="text-sm text-muted-foreground">Active Links • ሊንኮች</p>
            <p className="text-2xl font-bold mt-2">{summaryLoading ? "…" : activeLinks}</p>
            <p className="text-xs text-muted-foreground mt-1">Live payment links</p>
          </Card>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 space-y-4">
            <Card className="p-4">
              <h3 className="font-semibold">Recent Payments • የቅርብ ክፍያዎች</h3>
              <div className="mt-3 space-y-2">
                {recent.length === 0 && !summaryLoading && (
                  <p className="text-sm text-muted-foreground">No payments yet.</p>
                )}
                {recent.map((p) => (
                  <div key={p.id} className="flex items-center justify-between rounded-xl border p-3 hover:bg-muted">
                    <div className="flex items-center gap-3">
                      <div className="h-8 w-8 rounded-full bg-green-500/20 text-green-700 flex items-center justify-center text-xs">✓</div>
                      <div>
                        <p className="text-sm font-medium">
                          ETB {formatETB(p.amount)} • {p.method || "—"}{" "}
                          {p.requires_2fa && "• 2FA"}
                        </p>
                        <p className="text-xs text-muted-foreground">{p.tx_ref}</p>
                      </div>
                    </div>
                    <span
                      className={`text-xs px-2 py-0.5 rounded-full ${
                        p.status === "succeeded"
                          ? "bg-green-500/20 text-green-700"
                          : "bg-red-100 text-red-700"
                      }`}
                    >
                      {p.status}
                    </span>
                  </div>
                ))}
              </div>
            </Card>

            <Card className="p-4">
              <h3 className="font-semibold">Quick Actions • ፈጣን እርምጃዎች</h3>
              <div className="mt-3 grid grid-cols-3 gap-2">
                <button className="rounded-xl border p-3 text-sm font-medium hover:bg-muted">Create Link • ሊንክ</button>
                <button className="rounded-xl border p-3 text-sm font-medium hover:bg-muted">Pay Vendor • አቅራቢ ክፍያ</button>
                <button className="rounded-xl border p-3 text-sm font-medium hover:bg-muted">Run Payroll • ደሞዝ</button>
              </div>
            </Card>
          </div>

          <GlassCard className="p-4 space-y-3">
            <h3 className="font-semibold">AI Chat • Swarm 🤖</h3>
            <div className="rounded-xl bg-card border p-3 text-sm">
              <p className="text-muted-foreground">Goal: Create link 100 ETB for coffee if today TPV&gt;0</p>
              <p className="mt-2 text-xs font-medium">
                Final: Created link https://checkout.apexpay.et/c/coffee100
              </p>
            </div>
            <div className="rounded-xl bg-muted p-3 text-xs">
              <p className="font-semibold">RAG Ask • ተገዢነት</p>
              <p className="mt-1">Q: When is 2FA required? A: Transactions above 5000 ETB require 2FA per ONPS/10/2025</p>
            </div>
          </GlassCard>
        </div>
      </div>
    </div>
  )
}

function successRateAsNum(s: { success_rate_7_days: number }): number {
  return s.success_rate_7_days * 100
}
