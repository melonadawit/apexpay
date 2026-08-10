"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function AnalyticsPage() {
  const { t } = useLanguage()
  const { checking } = useRequireAuth()
  const { data: revenue } = useData(() => api.analytics.revenue(), [])
  const { data: methods } = useData(() => api.analytics.methods(), [])
  const { data: cohorts } = useData(() => api.analytics.cohorts(), [])

  if (checking) return <Centered>Checking session…</Centered>

  const totalRev = (revenue ?? []).reduce((s, d) => s + Number(d.revenue || 0), 0)
  const totalPays = (revenue ?? []).reduce((s, d) => s + d.payment_count, 0)
  const totalSucc = (revenue ?? []).reduce((s, d) => s + d.success_count, 0)

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">{t("Analytics & Cohort","ትንታኔ")}</h1>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card label="Revenue (30d)" value={`ETB ${totalRev.toLocaleString()}`} />
          <Card label="Payments (30d)" value={String(totalPays)} />
          <Card label="Success Rate" value={`${totalPays ? ((totalSucc / totalPays) * 100).toFixed(1) : 0}%`} />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="rounded-2xl border bg-card overflow-hidden">
            <h3 className="font-semibold p-4">Revenue by Day</h3>
            <div className="space-y-1 px-4 pb-4">
              {(revenue ?? []).slice(0, 15).map((d) => (
                <div key={d.stat_date} className="flex justify-between text-xs border-t pt-1">
                  <span>{d.stat_date}</span>
                  <span className="font-medium">ETB {Number(d.revenue).toLocaleString()}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="rounded-2xl border bg-card overflow-hidden">
            <h3 className="font-semibold p-4">Success by Method</h3>
            <div className="space-y-1 px-4 pb-4">
              {(methods ?? []).map((m) => (
                <div key={m.method} className="flex justify-between text-xs border-t pt-1">
                  <span className="capitalize">{m.method}</span>
                  <span>{m.success}/{m.count} • ETB {Number(m.revenue).toLocaleString()}</span>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="rounded-2xl border bg-card overflow-hidden">
          <h3 className="font-semibold p-4">Subscription Cohort Retention</h3>
          <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold">
            <span>Cohort</span><span>Customers</span><span>M1</span><span>M2</span><span>M3</span><span>MRR</span>
          </div>
          {(cohorts ?? []).map((c) => (
            <div key={c.cohort_month} className="grid grid-cols-6 gap-2 p-3 border-t text-xs">
              <span>{c.cohort_month}</span><span>{c.customers}</span><span>{c.month1_retention}%</span><span>{c.month2_retention}%</span><span>{c.month3_retention}%</span><span>ETB {Number(c.mrr).toLocaleString()}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function Card({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-2xl border bg-card p-6">
      <p className="text-sm text-muted-foreground">{label}</p>
      <p className="text-2xl font-bold mt-2">{value}</p>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
