"use client"
import * as React from "react"
import Link from "next/link"

export default function RefundDetailPage({ params }: { params: { id: string } }) {
  return (
    <div className="min-h-screen bg-neutral-50 p-6">
      <div className="max-w-4xl mx-auto space-y-4">
        <Link href="/refunds" className="text-sm text-primary">← Back to Refunds • ተመላሾች</Link>
        <h1 className="text-2xl font-bold">Refund Exam • {params.id} — M2 Ledger Balanced</h1>

        <div className="grid grid-cols-3 gap-4">
          <div className="col-span-2 space-y-4">
            <div className="rounded-2xl border bg-white p-4">
              <h3 className="font-semibold">Refund Lifecycle • የሂደት ዑደት</h3>
              <div className="mt-3 relative pl-6 border-l-2 border-neutral-200 space-y-3">
                <div className="relative text-xs"><div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-primary border-2 border-white" /><p>created → processing • refund_ref ref_001 amount 100.00 ETB payment txr_01H 500.00 fee_policy pro_rata fee_reversal 2.90 pro_rata = fee*refund/pay Round2 ETB scale bankers rounding decimal precise</p></div>
                <div className="relative text-xs"><div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-primary border-2 border-white" /><p>processing → succeeded • connector Mock refund mock_ref_ref_001 success • ledger M2 posting_key refund:{params.id} balanced true</p></div>
                <div className="relative text-xs"><div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-green-500 border-2 border-white" /><p>succeeded → webhook pending • outbox refund.succeeded • webhook delivery success 200 attempt1 HMAC X-ApexPay-Signature valid per secret prefix whsec_</p></div>
              </div>
            </div>

            <div className="rounded-2xl border bg-white p-4">
              <h3 className="font-semibold">Ledger Journal M2 • Dr payable R-FR + Dr fee_due FR Cr clearing R — Balanced ✓</h3>
              <div className="mt-3 space-y-2 text-xs font-mono">
                <div className="grid grid-cols-4 gap-2 bg-neutral-50 p-2 rounded"><span>Account</span><span>Direction</span><span>Amount ETB</span><span>Note</span></div>
                <div className="grid grid-cols-4 gap-2 p-2 border-b"><span>liability:merchant_payable</span><span className="text-green-600">debit</span><span>97.10</span><span>R-FR = 100-2.90 per pro_rata</span></div>
                <div className="grid grid-cols-4 gap-2 p-2 border-b"><span>liability:platform_fee_due</span><span className="text-green-600">debit</span><span>2.90</span><span>FR fee reversal pro_rata = fee*refund/pay Round2</span></div>
                <div className="grid grid-cols-4 gap-2 p-2"><span>asset:clearing:mock</span><span className="text-red-600">credit</span><span>100.00</span><span>R full refund amount</span></div>
                <p className="text-[11px] text-muted-foreground">Debit 100.00 = 97.10+2.90 == Credit 100.00 per ValidateBalanced O(n) + filtered zero entries if feeReversal zero + quality check SQL having sum(debit)!=sum(credit) expect 0 rows</p>
                <p className="text-[11px]">Payment status → partially_refunded / refunded CASE WHEN COALESCE(SUM(amount),0) FROM refunds WHERE payment_id=$1 AND status IN ('processing','succeeded') >= amount THEN refunded ELSE partially_refunded</p>
              </div>
            </div>
          </div>

          <div className="space-y-4">
            <div className="rounded-2xl border bg-white p-4">
              <h4 className="font-semibold text-sm">Idempotency • (merchant_id, refund_ref) Unique</h4>
              <p className="text-xs mt-2">refund_ref ref_001 unique index per merchant → duplicate with different amount 409 conflict duplicate_refund_ref stable code per SAD §11. Same amount returns same refund id/url per Idempotency-Key header.</p>
            </div>
            <div className="rounded-2xl border bg-white p-4">
              <h4 className="font-semibold text-sm">Fee Policies • Outstanding</h4>
              <ul className="text-xs mt-2 space-y-1 list-disc list-inside">
                <li>non_refundable: platform keeps fee FR=0 Dr payable R Cr clearing R</li>
                <li>pro_rata: FR = totalFee * (refund/pay) Round2 bankers rounding decimal precise</li>
                <li>full: FR = totalFee if refund==payment else 0</li>
              </ul>
            </div>
            <div className="rounded-2xl border bg-white p-4">
              <h4 className="font-semibold text-sm">Actions</h4>
              <div className="mt-2 grid grid-cols-2 gap-2">
                <button className="rounded-xl border h-10 text-xs">Resend Webhook • refund.succeeded</button>
                <button className="rounded-xl border h-10 text-xs">Evidence Pack JSON • NBE</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
