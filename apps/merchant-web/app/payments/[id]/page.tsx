"use client"
import * as React from "react"

export default function PaymentDetailPage({ params }: { params: { id: string } }) {
  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-4xl mx-auto space-y-4">
        <h1 className="text-2xl font-bold">Payment Exam • {params.id}</h1>
        <p className="text-sm text-muted-foreground">NBE exam console reconstruct any tx_ref &lt;60s per SAD A1 — lifecycle + journals + connector refs + webhooks + agent actions</p>

        <div className="grid grid-cols-3 gap-4">
          <div className="col-span-2 space-y-4">
            <div className="rounded-2xl border bg-card p-4">
              <h3 className="font-semibold">Lifecycle • የህይወት ዑደት</h3>
              <div className="mt-3 relative pl-6 border-l-2 border-neutral-200 space-y-3">
                {[
                  "created → pending • routed via telebirr primary (health success 96% latency 210ms) score 0.88 chosen true reason primary healthy",
                  "pending → processing • connector Initialize mock_ref_txr_01H checkout_url https://checkout.apexpay.et/mock/txr_01H",
                  "processing → succeeded • connector Verify succeeded amount 500.00 • ledger M1 journal posting_key payment_success:pay_01H balanced true",
                  "succeeded → webhook pending • outbox payment.succeeded published_at 2026-08-04T09:10:00Z • webhook delivery success 200 attempt 1",
                ].map((t,i)=>(
                  <div key={i} className="relative text-xs"><div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-primary border-2 border-white" /><p>{t}</p></div>
                ))}
              </div>
            </div>

            <div className="rounded-2xl border bg-card p-4">
              <h3 className="font-semibold">Ledger Journal • መዝገብ — posting_key payment_success:{params.id} — balanced ✓</h3>
              <div className="mt-3 space-y-2 text-xs font-mono">
                <div className="grid grid-cols-3 gap-2 bg-muted p-2 rounded"><span>Account</span><span>Direction</span><span>Amount ETB</span></div>
                <div className="grid grid-cols-3 gap-2 p-2 border-b"><span>asset:clearing:telebirr</span><span className="text-green-600">debit</span><span>500.00</span></div>
                <div className="grid grid-cols-3 gap-2 p-2 border-b"><span>liability:merchant_payable</span><span className="text-red-600">credit</span><span>485.50</span></div>
                <div className="grid grid-cols-3 gap-2 p-2"><span>liability:platform_fee_due</span><span className="text-red-600">credit</span><span>14.50</span></div>
                <p className="text-[11px] text-muted-foreground">Debit 500.00 == Credit 500.00 (485.50+14.50) per ValidateBalanced O(n) + quality check SQL having sum(debit)!=sum(credit) expect 0 rows</p>
              </div>
            </div>
          </div>

          <div className="space-y-4">
            <div className="rounded-2xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Routing Decision • ዘመናዊ መንገድ</h4>
              <p className="text-xs mt-2">Rule: Medium 1000-50000 ETB success_rate priority 20 — primary telebirr fallback cbe_birr fallback2 mock</p>
              <p className="text-xs mt-1">Chosen: telebirr • Reason: primary healthy success 96% latency 210ms • Fallback trail: none • fallback_used false</p>
              <p className="text-xs mt-1">Health snapshot: telebirr success 0.96 latency 210ms circuit closed, cbe_birr success 0.89 latency 260ms, mock success 1.0 latency 45ms</p>
            </div>

            <div className="rounded-2xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Webhook Deliveries • ዌብሁክ</h4>
              <div className="mt-2 text-xs space-y-1">
                <p>payment.succeeded → https://merchant.example.et/webhook • 200 • attempt 1 • 2026-08-04 09:10:01Z • HMAC valid per webhook secret prefix</p>
                <p className="text-muted-foreground">Retry exponential backoff 1s, 2s, 4s • SSRF block private ranges • idempotent receiver documented</p>
              </div>
            </div>

            <div className="rounded-2xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Actions • እርምጃዎች</h4>
              <div className="mt-2 grid grid-cols-2 gap-2">
                <button className="rounded-xl border h-10 text-xs">Refund • ተመላሽ</button>
                <button className="rounded-xl border h-10 text-xs">Resend Webhook</button>
                <button className="rounded-xl border h-10 text-xs col-span-2">Evidence Pack JSON • For NBE</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
