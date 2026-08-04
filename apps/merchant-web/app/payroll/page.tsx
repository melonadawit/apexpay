"use client"
import * as React from "react"
import Link from "next/link"

export default function PayrollPage() {
  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold">Payroll • ደሞዝ • Workforce Money OS</h1>

        <div className="grid grid-cols-3 gap-4">
          <div className="rounded-2xl border bg-card p-4"><p className="text-sm text-muted-foreground">Total Gross • አጠቃላይ</p><p className="text-2xl font-bold">ETB 200,000</p><p className="text-xs">10 employees • July 2026 regular</p></div>
          <div className="rounded-2xl border bg-card p-4"><p className="text-sm text-muted-foreground">Total Tax • ግብር</p><p className="text-2xl font-bold">ETB 20,000</p><p className="text-xs">ET brackets binary search O(log n) • 0-600 0% etc</p></div>
          <div className="rounded-2xl border bg-card p-4"><p className="text-sm text-muted-foreground">Total Net • የተጣራ</p><p className="text-2xl font-bold">ETB 150,000</p><p className="text-xs">Pension 7%/11% • OT 1.25/1.5/2.0</p></div>
        </div>

        <div className="grid grid-cols-4 gap-6">
          <div className="rounded-2xl border bg-card p-4">
            <h3 className="font-semibold">Employees • 10 • Fayda badge</h3>
            <div className="mt-3 space-y-2 text-xs">
              {[
                { code: "EMP001", name: "Abebe Kebede", base: "20000", fayda: true, bank: "CBE ****1234", cost: "Sales" },
                { code: "EMP002", name: "Almaz Tadesse", base: "25000", fayda: true, bank: "Awash ****5678", cost: "Eng" },
              ].map(e => (
                <div key={e.code} className="flex items-center justify-between rounded-xl border p-2"><div><p className="font-medium">{e.name} • {e.code}</p><p className="text-[11px] text-muted-foreground">Base {e.base} • {e.bank} • {e.cost} {e.fayda && "• Fayda ✓"}</p></div></div>
              ))}
            </div>
            <button className="mt-3 w-full rounded-xl border h-10 text-xs">Import CSV • 500 employees &lt;2s p99</button>
          </div>

          <div className="col-span-3 rounded-2xl border bg-card p-4">
            <div className="flex justify-between items-center"><h3 className="font-semibold">Payroll Runs • Runs table status pipeline visual stepper</h3><button className="rounded-xl bg-primary text-foreground px-4 h-9 text-xs">Create Run July regular</button></div>
            <div className="mt-4 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-xs font-semibold"><span>Run Ref</span><span>Period</span><span>Type</span><span>Status</span><span>Total Net</span><span>Action</span></div>
              <div className="grid grid-cols-6 gap-2 p-3 border-t text-xs"><span>prun_July2026</span><span>07/2026</span><span>regular</span><span><span className="px-2 py-0.5 rounded-full bg-amber-500/20">pending_approval</span></span><span>ETB 150,000</span><Link href="/payroll/prun_July2026" className="text-primary">View • Calculate → Approve dual &gt;100k → Disburse → payout batch</Link></div>
            </div>

            <div className="mt-6 rounded-xl bg-blue-500/10 border border-blue-500/20 p-3 text-xs">
              <p className="font-semibold">Ledger M4 per run book: Dr expense:salary 200k Cr payroll_payable 150k Cr et_income_tax_payable 20k Cr pension_payable 30k balanced ValidateBalanced + second journal Dr payroll_payable Cr bank 150k via payouts</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
