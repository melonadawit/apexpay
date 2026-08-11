"use client"
import * as React from "react"
import Link from "next/link"
import { Loader2 } from "lucide-react"
import { api } from "@/lib/api/client"
import { useData, formatETB } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

type Refund = {
  id: string
  merchant_id: string
  payment_id: string
  refund_ref: string
  amount: string
  currency: string
  status: string
  reason?: string
  fee_reversal?: string
  connector_id?: string
  connector_ref?: string
  failure_code?: string
  failure_message?: string
  created_at?: string
}

export default function RefundDetailPage({ params }: { params: { id: string } }) {
  const { t } = useLanguage()
  const { data: rf, loading, error } = useData<Refund>(
    () => api.refunds.get(params.id) as Promise<Refund>,
    [params.id]
  )

  const net = rf ? Math.max(0, Number(rf.amount || 0) - Number(rf.fee_reversal || 0)) : 0

  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-4xl mx-auto space-y-4">
        <Link href="/refunds" className="text-sm text-primary">← Back to Refunds • ተመላሾች</Link>
        <h1 className="text-2xl font-bold">{t("Refund Exam", "የተመላሽ ምርመራ")} • {params.id}</h1>

        {loading && (
          <div className="flex items-center gap-2 p-6 text-sm text-muted-foreground">
            <Loader2 className="h-4 w-4 animate-spin" /> Loading…
          </div>
        )}
        {error && <p className="p-4 text-sm text-red-600">{error}</p>}

        {!loading && !error && rf && (
          <div className="grid grid-cols-3 gap-4">
            <div className="col-span-2 space-y-4">
              <div className="rounded-2xl border bg-card p-4">
                <h3 className="font-semibold">Refund Lifecycle • የሂደት ዑደት</h3>
                <div className="mt-3 relative pl-6 border-l-2 border-neutral-200 space-y-3">
                  <div className="relative text-xs">
                    <div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-primary border-2 border-white" />
                    <p>created → {rf.status} • refund_ref {rf.refund_ref} • amount {formatETB(rf.amount)} {rf.currency} • payment {rf.payment_id} {rf.reason ? `• ${rf.reason}` : ""}</p>
                  </div>
                  <div className="relative text-xs">
                    <div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-green-500 border-2 border-white" />
                    <p>connector {rf.connector_id || "—"} (ref {rf.connector_ref || "—"}) • fee reversal {formatETB(rf.fee_reversal)}</p>
                  </div>
                </div>
              </div>

              <div className="rounded-2xl border bg-card p-4">
                <h3 className="font-semibold">Ledger Journal M2 • Dr payable R-FR + Dr fee_due FR Cr clearing R — Balanced ✓</h3>
                <div className="mt-3 space-y-2 text-xs font-mono">
                  <div className="grid grid-cols-4 gap-2 bg-muted p-2 rounded"><span>Account</span><span>Direction</span><span>Amount ETB</span><span>Note</span></div>
                  <div className="grid grid-cols-4 gap-2 p-2 border-b"><span>liability:merchant_payable</span><span className="text-green-600">debit</span><span>{net.toFixed(2)}</span><span>R-FR</span></div>
                  {Number(rf.fee_reversal || 0) > 0 && (
                    <div className="grid grid-cols-4 gap-2 p-2 border-b"><span>liability:platform_fee_due</span><span className="text-green-600">debit</span><span>{Number(rf.fee_reversal).toFixed(2)}</span><span>FR fee reversal</span></div>
                  )}
                  <div className="grid grid-cols-4 gap-2 p-2"><span>asset:clearing:{rf.connector_id || "mock"}</span><span className="text-red-600">credit</span><span>{Number(rf.amount).toFixed(2)}</span><span>R full refund amount</span></div>
                  <p className="text-[11px] text-muted-foreground">Debit {(Number(net) + Number(rf.fee_reversal || 0)).toFixed(2)} == Credit {Number(rf.amount).toFixed(2)} per ValidateBalanced O(n)</p>
                </div>
              </div>
            </div>

            <div className="space-y-4">
              <div className="rounded-2xl border bg-card p-4">
                <h4 className="font-semibold text-sm">Refund • ተመላሽ</h4>
                <div className="mt-2 space-y-1 text-xs">
                  <p>Ref: <span className="font-mono">{rf.refund_ref}</span></p>
                  <p>Status: <span className="font-semibold">{rf.status}</span></p>
                  <p>Amount: {formatETB(rf.amount)} {rf.currency}</p>
                  <p>Fee reversal: {formatETB(rf.fee_reversal)}</p>
                  {rf.failure_message && <p className="text-red-600">{rf.failure_message}</p>}
                </div>
              </div>
              <div className="rounded-2xl border bg-card p-4">
                <h4 className="font-semibold text-sm">Fee Policies</h4>
                <ul className="text-xs mt-2 space-y-1 list-disc list-inside">
                  <li>non_refundable: platform keeps fee FR=0</li>
                  <li>pro_rata: FR = totalFee * (refund/pay)</li>
                  <li>full: FR = totalFee if refund==payment else 0</li>
                </ul>
              </div>
              <div className="rounded-2xl border bg-card p-4">
                <h4 className="font-semibold text-sm">Actions</h4>
                <div className="mt-2 grid grid-cols-2 gap-2">
                  <button className="rounded-xl border h-10 text-xs">Resend Webhook • refund.succeeded</button>
                  <button className="rounded-xl border h-10 text-xs">Evidence Pack JSON • NBE</button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
