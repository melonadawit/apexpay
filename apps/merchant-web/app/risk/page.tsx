"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type RiskRule } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function RiskPage() {
  const { checking } = useRequireAuth()
  const { data: rules, refetch } = useData(() => api.risk.rules(), [])
  const { data: flags } = useData(() => api.risk.flags(), [])

  const [name, setName] = React.useState("")
  const [ruleType, setRuleType] = React.useState("threshold_amount")
  const [amountLimit, setAmountLimit] = React.useState("500000")
  const [action, setAction] = React.useState("flag")
  const [err, setErr] = React.useState("")
  const [saving, setSaving] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const addRule = async () => {
    setSaving(true); setErr("")
    try {
      await api.risk.createRule({
        name, rule_type: ruleType, action, severity: "high",
        parameters: { window_minutes: 60, amount_limit: amountLimit },
      })
      setName(""); refetch()
    } catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Risk & Fraud • አደጋ</h1>
        <p className="text-sm text-muted-foreground">Transaction monitoring, velocity checks, and review flags.</p>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Add Rule</h3>
            <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Rule name" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <select value={ruleType} onChange={(e) => setRuleType(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
              <option value="threshold_amount">Threshold amount</option>
              <option value="velocity_amount">Velocity amount</option>
              <option value="velocity_count">Velocity count</option>
              <option value="high_ticket">High ticket</option>
              <option value="high_failure_rate">High failure rate</option>
            </select>
            <input value={amountLimit} onChange={(e) => setAmountLimit(e.target.value)} placeholder="Amount limit" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <select value={action} onChange={(e) => setAction(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
              <option value="flag">Flag</option><option value="review">Review</option><option value="block">Block</option>
            </select>
            {err && <p className="text-sm text-red-600">{err}</p>}
            <button onClick={addRule} disabled={saving} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">
              {saving ? "Adding…" : "Add Rule"}
            </button>
          </div>

          <div className="rounded-2xl border bg-card overflow-hidden">
            <h3 className="font-semibold p-4">Rules</h3>
            {(rules ?? []).map((r) => (
              <div key={r.id} className="border-t p-3 text-sm flex justify-between">
                <span>{r.name} <span className="text-[10px] text-muted-foreground">{r.rule_type}</span></span>
                <span className={`text-[11px] px-2 py-0.5 rounded-full ${r.severity === "high" ? "bg-red-500/15 text-red-700" : "bg-amber-500/15 text-amber-700"}`}>{r.action} • {r.severity}</span>
              </div>
            ))}
          </div>

          <div className="rounded-2xl border bg-card overflow-hidden">
            <h3 className="font-semibold p-4">Flags • Review Queue</h3>
            {(flags ?? []).slice(0, 10).map((f) => (
              <div key={f.id} className="border-t p-3 text-xs">
                <p className="font-medium">{f.rule_name}</p>
                <p className="text-muted-foreground mt-1">{f.reason}</p>
                <p className="text-[10px] text-muted-foreground mt-1">{f.entity_type} • {f.entity_id} • {f.status}</p>
              </div>
            ))}
            {(flags ?? []).length === 0 && <p className="p-4 text-sm text-muted-foreground">No flags.</p>}
          </div>
        </div>
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
