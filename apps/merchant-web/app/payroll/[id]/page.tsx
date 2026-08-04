"use client"
import * as React from "react"

export default function PayrollRunDetail({ params }: { params: { id: string } }) {
  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-6xl mx-auto space-y-4">
        <h1 className="text-2xl font-bold">Payroll Run Detail • {params.id}</h1>
        <div className="flex items-center gap-2 text-xs">
          <span className="px-3 py-1 rounded-full bg-neutral-200">draft → calculating → pending_approval • current</span>
          <span className="px-3 py-1 rounded-full bg-amber-500/20">pending_approval • Needs dual if &gt;100k net</span>
        </div>

        <div className="rounded-2xl border bg-card overflow-hidden">
          <div className="grid grid-cols-8 gap-2 bg-muted p-3 text-xs font-semibold">
            <span>Employee</span><span>Gross</span><span>OT</span><span>Taxable</span><span>Income Tax ET</span><span>Pension 7%/11%</span><span>Net</span><span>Status</span>
          </div>
          {[
            {name:"Abebe Kebede", gross:"20000", ot:"1250", taxable:"18750", tax:"1800", pension:"1400/2200", net:"16800"},
            {name:"Almaz Tadesse", gross:"25000", ot:"0", taxable:"23250", tax:"2500", pension:"1750/2750", net:"20750"},
          ].map((r,i)=>(
            <div key={i} className="grid grid-cols-8 gap-2 p-3 border-t text-xs hover:bg-muted">
              <span>{r.name}</span><span>{r.gross}</span><span>{r.ot}</span><span>{r.taxable}</span><span>{r.tax}</span><span>{r.pension}</span><span className="font-bold">{r.net}</span><span>calculated</span>
            </div>
          ))}
          <div className="grid grid-cols-8 gap-2 p-3 bg-muted font-bold text-xs sticky bottom-0">
            <span>Total 10 emps</span><span>200,000</span><span>5,000</span><span>185,000</span><span>20,000</span><span>30,000 pension total</span><span>150,000 net</span><span></span>
          </div>
        </div>

        <div className="flex gap-2">
          <button className="rounded-xl bg-primary text-foreground px-6 h-12">Approve • dual finance+admin</button>
          <button className="rounded-xl border px-6 h-12">Disburse → payout batch • ledger M4</button>
          <button className="rounded-xl border px-6 h-12">Download Payslips PDF ZIP • outstanding modern template QR</button>
          <button className="rounded-xl border px-6 h-12">ET Report CSV ERCA • JSON</button>
        </div>

        <div className="rounded-2xl border bg-card p-4">
          <h3 className="font-semibold text-sm">Payslip PDF Preview • Outstanding Modern</h3>
          <div className="mt-3 rounded-xl border p-4 bg-gradient-to-br from-white to-neutral-50 max-w-[400px]">
            <div className="flex justify-between"><span className="font-bold">Apex Trading PLC</span><span className="text-xs">July 2026</span></div>
            <div className="mt-2 text-xs space-y-1">
              <p>Employee: Abebe Kebede • EMP001 • Sales • Fayda ****1234 ✓</p>
              <p>Base 20,000 + OT 1,250 (5h weekday 1.25x) = Gross 21,250</p>
              <p>Taxable 19,850 • Tax 1,800 (bracket 1651-3200 15%-142.5) • Pension Emp 7% 1,400 • Employer 11% 2,200</p>
              <p className="font-bold text-base">Net Pay ETB 16,800</p>
              <div className="mt-2 h-10 bg-muted/80 rounded flex items-center justify-center text-[10px]">Pie chart deductions • QR verification</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
