"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function AccountingPage() {
  const { checking } = useRequireAuth()
  const { data: accounts } = useData(() => api.accounting.accounts(), [])
  const { data: trial } = useData(() => api.accounting.trialBalance(), [])
  const { data: pnl } = useData(() => api.accounting.profitLoss(), [])
  const { data: bs } = useData(() => api.accounting.balanceSheet(), [])
  const { data: cashflow } = useData(() => api.accounting.cashFlow(), [])

  if (checking) return <Centered>Checking session…</Centered>

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Accounting & Bookkeeping • የሂሳብ አያያዝ</h1>
        <p className="text-sm text-muted-foreground">Chart of accounts, trial balance, P&L, balance sheet and cash flow — derived from the ledger.</p>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Statement title={pnl?.title} period={pnl?.period} lines={pnl?.lines} />
          <Statement title={bs?.title} period={bs?.period} lines={bs?.lines} />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card overflow-hidden lg:col-span-2">
            <h3 className="font-semibold p-4">Trial Balance</h3>
            <div className="grid grid-cols-4 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>Code</span><span>Account</span><span className="text-right">Debit</span><span className="text-right">Credit</span>
            </div>
            {(trial ?? []).map((t) => (
              <div key={t.code} className="grid grid-cols-4 gap-2 p-3 border-t text-xs">
                <span className="font-mono text-[10px]">{t.code}</span>
                <span>{t.name}</span>
                <span className="text-right">{t.debit !== "0" ? t.debit : ""}</span>
                <span className="text-right">{t.credit !== "0" ? t.credit : ""}</span>
              </div>
            ))}
          </div>

          <div className="space-y-4">
            <div className="rounded-2xl border bg-card p-4">
              <h3 className="font-semibold">Cash Flow</h3>
              <div className="mt-2 space-y-1">
                {(cashflow ?? []).map((c) => (
                  <div key={c.label} className="flex justify-between text-sm border-t pt-1">
                    <span>{c.label}</span><span className="font-medium">ETB {c.amount}</span>
                  </div>
                ))}
              </div>
            </div>
            <div className="rounded-2xl border bg-card p-4">
              <h3 className="font-semibold">Chart of Accounts</h3>
              <p className="text-xs text-muted-foreground mt-1">{accounts?.length ?? 0} accounts</p>
              <div className="mt-2 max-h-60 overflow-y-auto space-y-1">
                {(accounts ?? []).map((a) => (
                  <div key={a.code} className="flex justify-between text-xs border-t pt-1">
                    <span>{a.code}</span><span className="capitalize">{a.category}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function Statement({ title, period, lines }: { title?: string; period?: string; lines?: Array<{ label: string; amount: string; kind: string }> }) {
  return (
    <div className="rounded-2xl border bg-card p-5">
      <h3 className="font-semibold">{title || "Statement"}</h3>
      <p className="text-xs text-muted-foreground">{period}</p>
      <div className="mt-3 space-y-1">
        {(lines ?? []).map((l) => (
          <div key={l.label} className={`flex justify-between text-sm border-t pt-1 ${l.kind === "total" ? "font-bold" : ""}`}>
            <span>{l.label}</span><span>ETB {l.amount}</span>
          </div>
        ))}
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
