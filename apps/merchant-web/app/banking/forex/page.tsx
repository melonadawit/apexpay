"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type ForexRate, type ForexRequest } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function ForexPage() {
  const { checking } = useRequireAuth()
  const { data: rates, loading: ratesLoading } = useData(() => api.banking.forexRates(), [])
  const { data: requests, loading: reqLoading } = useData(() => api.banking.forexRequests(), [])

  if (checking) {
    return (
      <div className="min-h-screen flex items-center justify-center text-sm text-muted-foreground">
        Checking session…
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Forex • ውጭ ምንዛሪ</h1>
          <p className="text-sm text-muted-foreground mt-2">
            FDI transfers, 2.5% markup, NBE approval required. Rates cached 60s via Redis.
          </p>
        </div>

        <div className="rounded-2xl border bg-card overflow-hidden">
          <h3 className="font-semibold p-4">Live Rates (NBE)</h3>
          <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold">
            <span>Pair</span><span>Rate</span><span>Buy</span><span>Sell</span><span>Source</span><span>Updated</span>
          </div>
          {ratesLoading && <p className="p-4 text-sm text-muted-foreground">Loading…</p>}
          {(rates ?? []).map((r) => (
            <div key={r.from_currency + r.to_currency} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
              <span className="font-medium">{r.from_currency}→{r.to_currency}</span>
              <span>{r.rate}</span>
              <span>{r.buy_rate}</span>
              <span>{r.sell_rate}</span>
              <span>{r.source}</span>
              <span>{r.last_updated_at}</span>
            </div>
          ))}
        </div>

        <div className="rounded-2xl border bg-card overflow-hidden">
          <h3 className="font-semibold p-4">Forex Requests</h3>
          <div className="grid grid-cols-7 gap-2 bg-muted p-3 text-[11px] font-semibold">
            <span>ID</span><span>Pair</span><span>From</span><span>To</span><span>Fee %</span><span>Status</span><span>Created</span>
          </div>
          {reqLoading && <p className="p-4 text-sm text-muted-foreground">Loading…</p>}
          {(requests ?? []).map((r) => (
            <div key={r.id} className="grid grid-cols-7 gap-2 p-3 border-t text-xs hover:bg-muted/50">
              <span className="font-mono text-[10px]">{r.id}</span>
              <span>{r.from_currency}→{r.to_currency}</span>
              <span>ETB {r.from_amount}</span>
              <span>{r.to_amount}</span>
              <span>{r.forex_fee_percent}%</span>
              <span>
                <span className={`px-2 py-0.5 rounded-full text-[11px] ${r.status === "completed" ? "bg-green-500/15 text-green-700" : "bg-amber-500/15 text-amber-700"}`}>
                  {r.status}
                </span>
                {r.nbe_approval_status && <span className="ml-1 text-[10px] text-muted-foreground">NBE {r.nbe_approval_status}</span>}
              </span>
              <span>{r.created_at}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
