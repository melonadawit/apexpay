"use client"
import * as React from "react"
import Link from "next/link"

export default function SubscriptionDetailPage({ params }: { params: { id: string } }) {
  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-5xl mx-auto space-y-4">
        <Link href="/subscriptions" className="text-sm text-primary">← Back to Subscriptions</Link>
        <h1 className="text-2xl font-bold">Subscription Detail • {params.id} — Dunning 1d/3d/5d + Customer Portal Magic Link JWT 24h</h1>

        <div className="grid grid-cols-3 gap-4">
          <div className="col-span-2 space-y-4">
            <div className="rounded-2xl border bg-card p-4">
              <h3 className="font-semibold">Subscription Lifecycle • FSM incomplete→trialing→active→past_due→canceled/paused</h3>
              <div className="mt-3 relative pl-6 border-l-2 border-neutral-200 space-y-3 text-xs">
                <div className="relative"><div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-blue-500 border-2 border-white" /><p>incomplete → trialing • plan Monthly Coffee 500 ETB interval month trial 7d • current_period_start now current_period_end trial_end 7d</p></div>
                <div className="relative"><div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-green-500 border-2 border-white" /><p>trialing → active • trial_end passed • invoice draft → open due current_period_end attempt 0 next +1d • dunning worker cron hourly SELECT ... FOR UPDATE SKIP LOCKED</p></div>
                <div className="relative"><div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-amber-500 border-2 border-white" /><p>active → past_due • payment attempt failed • attempt_count 1 next_attempt_at +3d (72h) email mock • after 3 fails past_due webhook subscription.past_due • customer portal magic link JWT 24h shows invoices outstanding list with pay button</p></div>
              </div>
            </div>

            <div className="rounded-2xl border bg-card p-4">
              <h3 className="font-semibold">Invoices • Dunning Schedule Optimal 1d/3d/5d</h3>
              <div className="mt-3 grid grid-cols-6 gap-2 bg-muted p-3 text-xs font-semibold"><span>Invoice</span><span>Amount</span><span>Status</span><span>Attempt</span><span>Next Attempt</span><span>Payment</span></div>
              <div className="grid grid-cols-6 gap-2 p-3 border-t text-xs"><span>sinv_01H</span><span>ETB 500</span><span>open</span><span>1</span><span>+3d 2026-08-07 09:00 EAT</span><span>pay_01H failed insufficient</span></div>
              <div className="grid grid-cols-6 gap-2 p-3 border-t text-xs"><span>sinv_02H</span><span>ETB 500</span><span>paid</span><span>0</span><span>-</span><span>pay_02H succeeded ledger M1 linked subscription_id</span></div>
              <p className="text-[11px] text-muted-foreground mt-2">NextDunningAttempt(attemptCount) 0→+24h 1→+72h 2→+120h per subscription_invoices next_attempt_at index where status open attempt_count&lt;3 • Worker hourly SELECT ... FOR UPDATE SKIP LOCKED • Ledger invoice success posts M1 linked subscription_id Dr clearing Cr payable Cr fee_due</p>
            </div>
          </div>

          <div className="space-y-4">
            <div className="rounded-2xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Customer • cust_01H • Fayda ****1234</h4>
              <p className="text-xs mt-2">Abebe Kebede • abebe@example.et • Fayda FIN hash sha256(salt+FIN)+last4 only • Bank CBE ****1234 verified</p>
              <p className="text-xs mt-1">Portal: magic link JWT 24h token via signed JWT `customer_portal:&#123;customer_id&#125;:&#123;expiry&#125;` • Hosted pages outstanding pay button modern • /customer-portal/&#123;token&#125; shows invoices outstanding list with pay button</p>
              <button className="mt-3 w-full rounded-xl border h-9 text-xs">Open Customer Portal • magic link 24h</button>
            </div>

            <div className="rounded-2xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Plan • splan_01H • Monthly Coffee</h4>
              <p className="text-xs mt-2">Amount ETB 500 interval month interval_count 1 trial_days 7 status active • addInterval day/week 7*count/month/year via time.AddDate optimal</p>
            </div>

            <div className="rounded-2xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Actions • Outstanding</h4>
              <div className="mt-2 grid grid-cols-2 gap-2">
                <button className="rounded-xl border h-10 text-xs">Cancel • subscription.canceled webhook</button>
                <button className="rounded-xl border h-10 text-xs">Pause • paused</button>
                <button className="rounded-xl border h-10 text-xs">Resume • active</button>
                <button className="rounded-xl border h-10 text-xs">Proration • upgrade</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
