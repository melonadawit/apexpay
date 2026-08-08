"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type Dispute } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function DisputesPage() {
  const { checking } = useRequireAuth()
  const { data, refetch } = useData(() => api.disputes.list(), [])

  const [paymentId, setPaymentId] = React.useState("")
  const [amount, setAmount] = React.useState("")
  const [reason, setReason] = React.useState("fraud")
  const [err, setErr] = React.useState("")
  const [saving, setSaving] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const create = async () => {
    setSaving(true); setErr("")
    try { await api.disputes.create({ payment_id: paymentId, amount, reason_code: reason, currency: "ETB" }); setPaymentId(""); setAmount(""); refetch() }
    catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  const decide = async (id: string, decision: string) => { await api.disputes.decide(id, decision, "merchant resolution", "0"); refetch() }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Disputes & Chargebacks • ክርክር</h1>
        <p className="text-sm text-muted-foreground">File disputes, submit evidence, and resolve chargebacks.</p>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Open Dispute</h3>
            <input value={paymentId} onChange={(e) => setPaymentId(e.target.value)} placeholder="Payment ID" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="Amount ETB" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <select value={reason} onChange={(e) => setReason(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
              <option value="fraud">Fraud</option>
              <option value="service_not_received">Service not received</option>
              <option value="duplicate">Duplicate charge</option>
              <option value="refund_requested">Refund requested</option>
            </select>
            {err && <p className="text-sm text-red-600">{err}</p>}
            <button onClick={create} disabled={saving} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">Open Dispute</button>
          </div>

          <div className="lg:col-span-2 rounded-2xl border bg-card overflow-hidden">
            <div className="grid grid-cols-5 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>Payment</span><span>Amount</span><span>Reason</span><span>Status</span><span>Action</span>
            </div>
            {(data ?? []).map((d) => (
              <div key={d.id} className="grid grid-cols-5 gap-2 p-3 border-t text-xs">
                <span className="font-mono text-[10px]">{d.payment_id || d.id}</span>
                <span>ETB {d.amount}</span>
                <span>{d.reason_code}</span>
                <span className={`text-[11px] px-2 py-0.5 rounded-full w-fit ${d.status === "won" ? "bg-green-500/15 text-green-700" : d.status === "lost" ? "bg-red-500/15 text-red-700" : "bg-amber-500/15 text-amber-700"}`}>{d.status}</span>
                <span>
                  {d.status === "open" || d.status === "evidence_submitted" ? (
                    <div className="flex gap-1">
                      <button onClick={() => decide(d.id, "won")} className="text-green-600">Won</button>
                      <button onClick={() => decide(d.id, "lost")} className="text-red-600">Lost</button>
                    </div>
                  ) : d.resolution || "—"}
                </span>
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
