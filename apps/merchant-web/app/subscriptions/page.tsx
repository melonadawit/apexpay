"use client"
import * as React from "react"
import { useLanguage } from "@/components/providers/language-provider"

export default function SubscriptionsPage() {
  const { t } = useLanguage()
  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold">{t("Subscriptions","ደንበኝነት ምዝገባ")}</h1>

        <div className="grid grid-cols-4 gap-4">
          <div className="rounded-2xl border bg-card p-4"><p className="text-sm text-muted-foreground">MRR • ወርሃዊ</p><p className="text-2xl font-bold">ETB 25,000</p><p className="text-xs">12 active • 2 trialing • 1 past_due</p></div>
          <div className="rounded-2xl border bg-card p-4"><p className="text-sm text-muted-foreground">Churn • መውጣት</p><p className="text-2xl font-bold">2.1%</p><p className="text-xs">Dunning retry 1d/3d/5d</p></div>
          <div className="rounded-2xl border bg-card p-4"><p className="text-sm text-muted-foreground">Trials • ሙከራ</p><p className="text-2xl font-bold">5</p><p className="text-xs">7 days trial</p></div>
          <div className="rounded-2xl border bg-card p-4"><p className="text-sm text-muted-foreground">Invoices Overdue</p><p className="text-2xl font-bold">3</p><p className="text-xs">Need dunning</p></div>
        </div>

        <div className="grid grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-4">
            <h3 className="font-semibold">Plans • እቅዶች</h3>
            <div className="mt-3 space-y-2">
              <div className="rounded-xl border p-3"><p className="font-medium">Monthly Coffee ETB 500 • month 1 trial 7d</p><p className="text-xs text-muted-foreground">Interval month • Status active • Plan splan_01H</p></div>
              <input placeholder="Plan Name" className="w-full rounded-xl border h-10 px-3 text-sm" defaultValue="Monthly Coffee" />
              <input placeholder="Amount ETB" className="w-full rounded-xl border h-10 px-3 text-sm" defaultValue="500" />
              <div className="grid grid-cols-2 gap-2"><select className="rounded-xl border h-10 px-3 text-xs"><option>month</option><option>week</option></select><input placeholder="Trial Days" className="rounded-xl border h-10 px-3 text-xs" defaultValue="7" /></div>
              <button className="w-full rounded-xl bg-primary text-foreground h-10 text-sm">Create Plan • splan_ + amount decimal precise</button>
            </div>
          </div>

          <div className="col-span-2 rounded-2xl border bg-card p-4">
            <h3 className="font-semibold">Subscriptions • Customer Portal Magic Link JWT 24h</h3>
            <div className="mt-3 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-xs font-semibold"><span>Customer</span><span>Plan</span><span>Status</span><span>Period End</span><span>Invoice</span><span>Action</span></div>
              <div className="grid grid-cols-6 gap-2 p-3 border-t text-xs"><span>cust_01H • Abebe • Fayda ****1234</span><span>Monthly Coffee 500</span><span><span className="px-2 py-0.5 rounded-full bg-blue-500/20">trialing</span></span><span>2026-08-11 trial_end 7d</span><span>Draft • due 2026-08-11 • attempt 0 next +1d</span><span>Portal magic link • webhook sub.*</span></div>
              <div className="grid grid-cols-6 gap-2 p-3 border-t text-xs"><span>cust_02H • Almaz</span><span>Monthly Coffee 500</span><span><span className="px-2 py-0.5 rounded-full bg-amber-500/20">past_due</span></span><span>2026-07-30</span><span>Open • attempt 1 next +3d • email mock • after 3 fails past_due webhook</span><span>Cancel • Pause • Resume</span></div>
            </div>
            <div className="mt-3 rounded-xl bg-muted p-3 text-xs">
              <p className="font-semibold">Dunning Schedule Optimal:</p>
              <p>attempt 0 → +24h, 1 → +72h, 2 → +120h per subscription_invoices next_attempt_at index where status open attempt_count&lt;3 • Worker cron hourly SELECT ... FOR UPDATE SKIP LOCKED • Ledger invoice success posts M1 linked subscription_id</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
