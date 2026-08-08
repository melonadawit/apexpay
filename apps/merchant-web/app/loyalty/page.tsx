"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type LoyaltyTier } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function LoyaltyPage() {
  const { checking } = useRequireAuth()
  const { data: tiers, refetch } = useData(() => api.loyalty.tiers(), [])
  const { data: accounts } = useData(() => api.loyalty.accounts(), [])
  const { data: txs } = useData(() => api.loyalty.transactions(), [])

  const [tName, setTName] = React.useState("")
  const [minSpend, setMinSpend] = React.useState("")
  const [pct, setPct] = React.useState("1")
  const [err, setErr] = React.useState("")
  const [saving, setSaving] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const addTier = async () => {
    setSaving(true); setErr("")
    try { await api.loyalty.createTier({ name: tName, min_spend: minSpend, cashback_percent: pct }); setTName(""); setMinSpend(""); refetch() }
    catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Loyalty & Cashback • ታማኝነት</h1>
        <p className="text-sm text-muted-foreground">Reward repeat customers with tiers and cashback.</p>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Add Tier</h3>
            <input value={tName} onChange={(e) => setTName(e.target.value)} placeholder="Tier name" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={minSpend} onChange={(e) => setMinSpend(e.target.value)} placeholder="Min spend ETB" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={pct} onChange={(e) => setPct(e.target.value)} placeholder="Cashback %" className="w-full rounded-xl border h-11 px-3 text-sm" />
            {err && <p className="text-sm text-red-600">{err}</p>}
            <button onClick={addTier} disabled={saving} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">Add Tier</button>
          </div>

          <div className="rounded-2xl border bg-card overflow-hidden">
            <h3 className="font-semibold p-4">Tiers</h3>
            {(tiers ?? []).map((t) => (
              <div key={t.id} className="border-t p-3 flex justify-between text-sm">
                <span className="font-medium">{t.name}</span>
                <span className="text-xs text-muted-foreground">ETB {t.min_spend} • {t.cashback_percent}%</span>
              </div>
            ))}
          </div>

          <div className="rounded-2xl border bg-card overflow-hidden">
            <h3 className="font-semibold p-4">Top Customers</h3>
            {(accounts ?? []).slice(0, 8).map((a) => (
              <div key={a.id} className="border-t p-3 flex justify-between text-sm">
                <span>{a.customer_email || "—"}</span>
                <span className="text-xs">{a.points} pts • {a.tier_name || "—"}</span>
              </div>
            ))}
            <h3 className="font-semibold p-4 border-t">Cashback History</h3>
            {(txs ?? []).slice(0, 6).map((t) => (
              <div key={t.id} className="border-t p-2 text-xs flex justify-between">
                <span>{t.type}</span><span>ETB {t.amount}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
