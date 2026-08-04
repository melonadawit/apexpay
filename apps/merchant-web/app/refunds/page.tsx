"use client"
import * as React from "react"

export default function RefundsPage() {
  return (
    <div className="min-h-screen bg-neutral-50 p-6">
      <div className="max-w-6xl mx-auto space-y-4">
        <h1 className="text-2xl font-bold">Refunds • ተመላሽ ክፍያ — FULL M2 Ledger</h1>
        <div className="grid grid-cols-3 gap-4">
          <div className="rounded-2xl border bg-white p-6 space-y-3">
            <h3 className="font-semibold">Create Refund • Outstanding Bottom Sheet</h3>
            <select className="w-full rounded-xl border h-12 px-3"><option>Payment txr_01H ETB 500.00 succeeded</option><option>Payment txr_02H ETB 6000 2FA verified</option></select>
            <input placeholder="Refund Ref unique (merchant_id, refund_ref)" className="w-full rounded-xl border h-12 px-3" defaultValue="ref_001" />
            <input placeholder="Amount ETB partial allowed" className="w-full rounded-xl border h-12 px-3" defaultValue="100.00" />
            <select className="w-full rounded-xl border h-12 px-3"><option>Fee Policy: non_refundable (platform keeps fee) — default</option><option>pro_rata — reverse pro-rata fee * refund/pay Round2 ETB scale</option><option>full — full fee reversal on full refund</option></select>
            <input placeholder="Reason" className="w-full rounded-xl border h-12 px-3" defaultValue="Customer request" />
            <button className="w-full rounded-xl bg-primary text-white h-12 font-semibold">Refund • M2 Dr payable R-FR + Dr fee_due FR Cr clearing R</button>
            <p className="text-[11px] text-muted-foreground">Idempotency by (merchant_id, refund_ref) unique conflict if different amount 409. Remaining refundable = amount - refunded amount. Fee reversal calc decimal precise bankers rounding.</p>
          </div>
          <div className="col-span-2 rounded-2xl border bg-white p-4">
            <h3 className="font-semibold">Refunds List • Ledger M2</h3>
            <div className="mt-3 grid grid-cols-6 gap-2 bg-neutral-50 p-3 text-xs font-semibold"><span>Refund Ref</span><span>Payment</span><span>Amount</span><span>Fee Rev</span><span>Status</span><span>Ledger</span></div>
            <div className="grid grid-cols-6 gap-2 p-3 border-t text-xs"><span>ref_001</span><span>txr_01H</span><span>ETB 100.00</span><span>ETB 2.90 pro_rata</span><span><span className="px-2 py-0.5 rounded-full bg-green-100">succeeded</span></span><span>M2 balanced</span></div>
            <div className="mt-4 rounded-xl bg-blue-50 border border-blue-200 p-3 text-xs">
              <p className="font-semibold">Ledger M2 Example 100 ETB fee 2.90 pro_rata 100% refund:</p>
              <p>Dr liability:merchant_payable 97.10 + Dr liability:platform_fee_due 2.90 Cr asset:clearing:mock 100.00 — Debit 100 == Credit 100 balanced ValidateBalanced O(n) + filter zero entries</p>
              <p>Payment status → partially_refunded / refunded CASE WHEN sum(amount) &gt;= amount THEN refunded ELSE partially_refunded</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
