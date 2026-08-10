"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type Loan } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function LendingPage() {
  const { t } = useLanguage()
  const { checking } = useRequireAuth()
  const { data, refetch } = useData(() => api.lending.loans(), [])

  const [creditLineId, setCreditLineId] = React.useState("")
  const [amount, setAmount] = React.useState("")
  const [purpose, setPurpose] = React.useState("working_capital")
  const [err, setErr] = React.useState("")
  const [saving, setSaving] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const create = async () => {
    setSaving(true); setErr("")
    try { await api.lending.create({ credit_line_id: creditLineId, amount, purpose, currency: "ETB" }); setAmount(""); refetch() }
    catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">{t("Lending & Micro-loans","ብድር")}</h1>
        <p className="text-sm text-muted-foreground">Collateral-free loans from your credit line, scored on TPV and payroll.</p>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Request Loan</h3>
            <input value={creditLineId} onChange={(e) => setCreditLineId(e.target.value)} placeholder="Credit line ID" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="Amount ETB" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <select value={purpose} onChange={(e) => setPurpose(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
              <option value="working_capital">Working capital</option>
              <option value="inventory">Inventory</option>
              <option value="payroll">Payroll</option>
              <option value="expansion">Expansion</option>
            </select>
            {err && <p className="text-sm text-red-600">{err}</p>}
            <button onClick={create} disabled={saving} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">Request Loan</button>
          </div>

          <div className="lg:col-span-2 rounded-2xl border bg-card overflow-hidden">
            <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>ID</span><span>Amount</span><span>Purpose</span><span>Rate</span><span>Outstanding</span><span>Status</span>
            </div>
            {(data ?? []).map((l) => (
              <div key={l.id} className="grid grid-cols-6 gap-2 p-3 border-t text-xs">
                <span className="font-mono text-[10px]">{l.id}</span>
                <span className="font-semibold">ETB {l.amount}</span>
                <span>{l.purpose}</span>
                <span>{l.interest_rate}%</span>
                <span>ETB {l.outstanding_amount}</span>
                <span className={`text-[11px] px-2 py-0.5 rounded-full w-fit ${l.status === "repaid" ? "bg-green-500/15 text-green-700" : l.status === "disbursed" ? "bg-blue-500/15 text-blue-700" : "bg-amber-500/15 text-amber-700"}`}>{l.status}</span>
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
