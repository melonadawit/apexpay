"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type TreasuryPosition, type Forecast } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function TreasuryPage() {
  const { t } = useLanguage()
  const { checking } = useRequireAuth()
  const { data: position, refetch: refetchPos } = useData(() => api.treasury.position(), [])
  const { data: transfers, refetch: refetchXfer } = useData(() => api.treasury.transfers(), [])
  const { data: forecast, refetch: refetchForecast } = useData(() => api.treasury.forecast(), [])

  const [fromAcc, setFromAcc] = React.useState("")
  const [toAcc, setToAcc] = React.useState("")
  const [amount, setAmount] = React.useState("")
  const [err, setErr] = React.useState("")
  const [saving, setSaving] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const doTransfer = async () => {
    setSaving(true); setErr("")
    try {
      await api.treasury.createTransfer({ from_account_id: fromAcc, to_account_id: toAcc, amount, currency: "ETB", purpose: "concentration" })
      setFromAcc(""); setToAcc(""); setAmount("")
      refetchPos(); refetchXfer()
    } catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  const genForecast = async () => { setSaving(true); try { await api.treasury.generateForecast(); refetchForecast() } finally { setSaving(false) } }

  const accounts = (position?.accounts ?? []).map((a) => ({ id: a.account_id, label: `${a.account_name} (${a.account_number}) • ${a.balance}` }))

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Treasury • Cash Management • ግምጃ</h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6">
            <h3 className="font-semibold">Cash Position</h3>
            <p className="text-3xl font-bold mt-3">ETB {position?.total_balance ?? "—"}</p>
            <p className="text-xs text-muted-foreground">Available: ETB {position?.total_available ?? "—"}</p>
            <div className="mt-4 space-y-2">
              {(position?.accounts ?? []).map((a) => (
                <div key={a.account_id} className="rounded-xl border p-3">
                  <p className="text-sm font-medium">{a.account_name}</p>
                  <p className="text-xs text-muted-foreground">{a.account_number} • {a.account_type}</p>
                  <p className="text-sm mt-1">ETB {a.balance}</p>
                </div>
              ))}
            </div>
          </div>

          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Internal Transfer</h3>
            <select value={fromAcc} onChange={(e) => setFromAcc(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
              <option value="">From account…</option>
              {accounts.map((a) => <option key={a.id} value={a.id}>{a.label}</option>)}
            </select>
            <select value={toAcc} onChange={(e) => setToAcc(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
              <option value="">To account…</option>
              {accounts.map((a) => <option key={a.id} value={a.id}>{a.label}</option>)}
            </select>
            <input value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="Amount ETB" className="w-full rounded-xl border h-11 px-3 text-sm" />
            {err && <p className="text-sm text-red-600">{err}</p>}
            <button onClick={doTransfer} disabled={saving} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">
              {saving ? "Moving…" : "Transfer"}
            </button>
            <div className="space-y-1 pt-2">
              {(transfers ?? []).slice(0, 5).map((t) => (
                <p key={t.id} className="text-xs text-muted-foreground">{t.from_account_id} → {t.to_account_id} • ETB {t.amount} • {t.status}</p>
              ))}
            </div>
          </div>

          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <div className="flex justify-between items-center">
              <h3 className="font-semibold">Cash-flow Forecast</h3>
              <button onClick={genForecast} disabled={saving} className="rounded-xl border px-3 h-9 text-xs">Regenerate</button>
            </div>
            <Row label="Inflow (90d)" value={forecast?.inflow_90d} />
            <Row label="Outflow (90d)" value={forecast?.outflow_90d} />
            <Row label="Net (90d)" value={forecast?.net_90d} strong />
            <p className="text-[11px] text-muted-foreground pt-2">Based on payments, vendor invoices, payroll & tax obligations.</p>
          </div>
        </div>
      </div>
    </div>
  )
}

function Row({ label, value, strong }: { label: string; value?: string; strong?: boolean }) {
  return (
    <div className="flex justify-between text-sm">
      <span className="text-muted-foreground">{label}</span>
      <span className={strong ? "font-bold" : ""}>ETB {value ?? "—"}</span>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
