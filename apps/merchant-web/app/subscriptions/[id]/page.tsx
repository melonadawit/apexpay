"use client"
import * as React from "react"
import Link from "next/link"
import { Loader2 } from "lucide-react"
import { api, SubscriptionDetail } from "@/lib/api/client"
import { useData, formatETB } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function SubscriptionDetailPage({ params }: { params: { id: string } }) {
  const { t } = useLanguage()
  const { data, loading, error } = useData<SubscriptionDetail>(
    () => api.subscriptions.get(params.id),
    [params.id]
  )

  const sub = data?.subscription
  const plan = data?.plan
  const cust = data?.customer
  const invoices = data?.invoices ?? []

  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-5xl mx-auto space-y-4">
        <Link href="/subscriptions" className="text-sm text-primary">← Back to Subscriptions</Link>
        <h1 className="text-2xl font-bold">
          {t("Subscription Detail", "የደንበኝነት ምዝገባ ዝርዝር")} • {params.id}
        </h1>

        {loading && (
          <div className="flex items-center gap-2 p-6 text-sm text-muted-foreground">
            <Loader2 className="h-4 w-4 animate-spin" /> Loading…
          </div>
        )}
        {error && <p className="p-4 text-sm text-red-600">{error}</p>}

        {!loading && !error && sub && (
          <div className="grid grid-cols-3 gap-4">
            <div className="col-span-2 space-y-4">
              <div className="rounded-2xl border bg-card p-4">
                <h3 className="font-semibold">Subscription Lifecycle • FSM incomplete→trialing→active→past_due→canceled/paused</h3>
                <div className="mt-3 relative pl-6 border-l-2 border-neutral-200 space-y-3 text-xs">
                  <div className="relative">
                    <div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-primary border-2 border-white" />
                    <p>
                      Status: <span className="font-semibold">{sub.status}</span> • current period {sub.current_period_start?.slice(0, 10)} → {sub.current_period_end?.slice(0, 10)} {sub.trial_end ? `• trial ends ${sub.trial_end.slice(0, 10)}` : ""}
                    </p>
                  </div>
                  {plan && (
                    <div className="relative">
                      <div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-green-500 border-2 border-white" />
                      <p>Plan {plan.name} • {formatETB(plan.amount)} {plan.currency} interval {plan.interval_count} {plan.interval_type}{plan.trial_days ? ` trial ${plan.trial_days}d` : ""}</p>
                    </div>
                  )}
                </div>
              </div>

              <div className="rounded-2xl border bg-card p-4">
                <h3 className="font-semibold">Invoices • Dunning Schedule Optimal 1d/3d/5d</h3>
                <div className="mt-3 grid grid-cols-5 gap-2 bg-muted p-3 text-xs font-semibold">
                  <span>Invoice</span><span>Amount</span><span>Status</span><span>Attempt</span><span>Due</span>
                </div>
                {invoices.length === 0 && <p className="p-3 text-xs text-muted-foreground">No invoices yet.</p>}
                {invoices.map((inv) => (
                  <div key={inv.id} className="grid grid-cols-5 gap-2 p-3 border-t text-xs">
                    <span className="font-mono">{inv.id}</span>
                    <span>{formatETB(inv.amount)} {inv.currency}</span>
                    <span>
                      <span className={`px-2 py-0.5 rounded-full ${
                        inv.status === "paid" ? "bg-green-500/20 text-green-700" : inv.status === "open" ? "bg-amber-500/20 text-amber-700" : "bg-blue-500/20"
                      }`}>{inv.status}</span>
                    </span>
                    <span>{inv.attempt_count}</span>
                    <span>{inv.due_at?.slice(0, 10)}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="space-y-4">
              <div className="rounded-2xl border bg-card p-4">
                <h4 className="font-semibold text-sm">Customer</h4>
                <p className="text-xs mt-2 font-medium">{cust?.name || "—"}</p>
                <p className="text-xs text-muted-foreground">{cust?.email || ""} {cust?.phone || ""}</p>
                <button className="mt-3 w-full rounded-xl border h-9 text-xs">Open Customer Portal • magic link 24h</button>
              </div>

              <div className="rounded-2xl border bg-card p-4">
                <h4 className="font-semibold text-sm">Plan</h4>
                {plan ? (
                  <p className="text-xs mt-2">
                    {plan.name} • ETB {formatETB(plan.amount)} {plan.currency} interval {plan.interval_type} count {plan.interval_count} trial {plan.trial_days}d
                  </p>
                ) : <p className="text-xs mt-2 text-muted-foreground">No plan.</p>}
              </div>

              <div className="rounded-2xl border bg-card p-4">
                <h4 className="font-semibold text-sm">Actions</h4>
                <div className="mt-2 grid grid-cols-2 gap-2">
                  <button className="rounded-xl border h-10 text-xs">Cancel</button>
                  <button className="rounded-xl border h-10 text-xs">Pause</button>
                  <button className="rounded-xl border h-10 text-xs">Resume</button>
                  <button className="rounded-xl border h-10 text-xs">Proration</button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
