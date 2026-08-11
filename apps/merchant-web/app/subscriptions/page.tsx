"use client"
import * as React from "react"
import { Loader2, Inbox } from "lucide-react"
import { api } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

type Sub = {
  id?: string
  customer_id?: string
  customer_name?: string
  plan_id?: string
  plan_name?: string
  status?: string
  amount?: string
  currency?: string
  current_period_end?: string
  trial_end?: string
  created_at?: string
}

export default function SubscriptionsPage() {
  const { t } = useLanguage()
  const { data: subs, loading, error } = useData<Sub[]>(() => api.subscriptions.subscriptions() as Promise<Sub[]>, [])

  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold">{t("Subscriptions", "ደንበኝነት ምዝገባ")}</h1>

        <div className="rounded-2xl border bg-card p-4">
          <h3 className="font-semibold">Subscriptions</h3>
          {loading && (
            <div className="flex items-center gap-2 p-4 text-sm text-muted-foreground">
              <Loader2 className="h-4 w-4 animate-spin" /> Loading…
            </div>
          )}
          {error && <p className="p-4 text-sm text-red-600">{error}</p>}
          {!loading && !error && (subs ?? []).length === 0 && (
            <div className="flex flex-col items-center gap-2 p-8 text-center text-muted-foreground">
              <Inbox className="h-8 w-8 text-muted-foreground/50" />
              <p className="text-sm font-medium">No subscriptions yet.</p>
            </div>
          )}
          <div className="mt-3 grid grid-cols-6 gap-2 bg-muted p-3 text-xs font-semibold">
            <span>Customer</span><span>Plan</span><span>Amount</span><span>Status</span><span>Period End</span><span>Created</span>
          </div>
          {(subs ?? []).map((s, i) => (
            <div key={i} className="grid grid-cols-6 gap-2 p-3 border-t text-xs">
              <span className="font-mono">{s.customer_name || s.customer_id || "—"}</span>
              <span>{s.plan_name || s.plan_id || "—"}</span>
              <span>{s.amount ? `ETB ${s.amount}` : "—"}</span>
              <span>
                <span className={`px-2 py-0.5 rounded-full ${
                  s.status === "active" ? "bg-green-500/20" : s.status === "past_due" ? "bg-amber-500/20" : "bg-blue-500/20"
                }`}>
                  {s.status || "—"}
                </span>
              </span>
              <span>{s.current_period_end || s.trial_end || "—"}</span>
              <span>{s.created_at ? s.created_at.slice(0, 10) : "—"}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
