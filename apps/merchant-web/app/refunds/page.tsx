"use client"
import * as React from "react"
import { Loader2, Inbox } from "lucide-react"
import { api } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

type Refund = {
  id?: string
  refund_ref?: string
  payment_id?: string
  amount?: string
  fee_reversal_amount?: string
  fee_policy?: string
  status?: string
  reason?: string
  created_at?: string
}

export default function RefundsPage() {
  const { t } = useLanguage()
  const { data: payments = [] } = useData<Array<{ id: string; tx_ref?: string; amount?: string; status?: string }>>(() => api.payments(25) as Promise<any[]>, [])
  const [paymentId, setPaymentId] = React.useState("")
  const { data: refunds, loading, error } = useData<Refund[]>(() => api.refunds.byPayment(paymentId || "none") as Promise<Refund[]>, [paymentId])
  const [msg, setMsg] = React.useState("")
  const [ref, setRef] = React.useState("")
  const [amount, setAmount] = React.useState("")
  const [reason, setReason] = React.useState("")
  const [submitting, setSubmitting] = React.useState(false)

  // Default to the first payment.
  React.useEffect(() => {
    if (!paymentId && (payments ?? []).length > 0) setPaymentId((payments ?? [])[0].id)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [payments])

  const create = async () => {
    if (!paymentId || !amount) return
    setSubmitting(true)
    setMsg("")
    try {
      await api.refunds.create({ payment_id: paymentId, refund_ref: ref || undefined, amount, reason })
      setMsg(t("Refund created.", "ተመላሽ ተፈጥሯል።"))
    } catch (e) {
      setMsg((e as Error).message || "Could not create refund.")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-6xl mx-auto space-y-4">
        <h1 className="text-2xl font-bold">{t("Refunds", "ተመላሽ ክፍያ")}</h1>
        <div className="grid grid-cols-3 gap-4">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Create Refund</h3>
            <select
              value={paymentId}
              onChange={(e) => setPaymentId(e.target.value)}
              className="w-full rounded-xl border h-12 px-3"
            >
              <option value="">Select payment…</option>
              {(payments ?? []).map((p) => (
                <option key={p.id} value={p.id}>
                  {p.tx_ref || p.id} • ETB {p.amount || "0"} • {p.status}
                </option>
              ))}
            </select>
            <input
              placeholder="Refund ref (unique)"
              value={ref}
              onChange={(e) => setRef(e.target.value)}
              className="w-full rounded-xl border h-12 px-3"
            />
            <input
              placeholder="Amount ETB (partial allowed)"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="w-full rounded-xl border h-12 px-3"
            />
            <input
              placeholder="Reason"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              className="w-full rounded-xl border h-12 px-3"
            />
            <button
              onClick={create}
              disabled={submitting || !paymentId || !amount}
              className="w-full rounded-xl bg-primary text-foreground h-12 font-semibold disabled:opacity-50"
            >
              {submitting ? "Refunding…" : t("Refund", "ተመላሽ")}
            </button>
            {msg && <p className="text-xs text-muted-foreground">{msg}</p>}
          </div>

          <div className="col-span-2 rounded-2xl border bg-card p-4">
            <h3 className="font-semibold">Refunds List</h3>
            {loading && (
              <div className="flex items-center gap-2 p-4 text-sm text-muted-foreground">
                <Loader2 className="h-4 w-4 animate-spin" /> Loading…
              </div>
            )}
            {error && <p className="p-4 text-sm text-red-600">{error}</p>}
            {!loading && !error && (refunds ?? []).length === 0 && (
              <div className="flex flex-col items-center gap-2 p-8 text-center text-muted-foreground">
                <Inbox className="h-8 w-8 text-muted-foreground/50" />
                <p className="text-sm font-medium">No refunds yet.</p>
              </div>
            )}
            <div className="mt-3 grid grid-cols-6 gap-2 bg-muted p-3 text-xs font-semibold">
              <span>Refund Ref</span><span>Payment</span><span>Amount</span><span>Policy</span><span>Status</span><span>Reason</span>
            </div>
            {(refunds ?? []).map((r, i) => (
              <div key={i} className="grid grid-cols-6 gap-2 p-3 border-t text-xs">
                <span className="font-mono">{r.refund_ref || r.id || "—"}</span>
                <span className="font-mono">{r.payment_id || "—"}</span>
                <span>ETB {r.amount || "0"}</span>
                <span>{r.fee_policy || "—"}</span>
                <span>
                  <span className={`px-2 py-0.5 rounded-full ${r.status === "succeeded" ? "bg-green-500/20" : "bg-amber-500/20"}`}>
                    {r.status || "pending"}
                  </span>
                </span>
                <span>{r.reason || "—"}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
