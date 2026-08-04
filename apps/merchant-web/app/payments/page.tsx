"use client"
import * as React from "react"
import Link from "next/link"

const mockPayments = [
  { id:"pay_01H", tx_ref:"txr_01H", amount:"500.00", currency:"ETB", status:"succeeded", method:"telebirr", connector_id:"telebirr", fee:"14.50", net:"485.50", routing:"telebirr primary", checkout_url:"https://checkout.apexpay.et/c/abc", created_at:"2 min ago" },
  { id:"pay_02H", tx_ref:"txr_02H", amount:"6000.00", currency:"ETB", status:"succeeded", method:"bank", connector_id:"bank_ips", fee:"174.00", net:"5826.00", routing:"bank_ips primary fallback telebirr", requires_2fa:true, two_fa_verified:true },
  { id:"pay_03H", tx_ref:"txr_03H", amount:"1000.00", currency:"ETB", status:"failed", method:"cbe_birr", failure_code:"insufficient_balance" },
]

export default function PaymentsPage() {
  return (
    <div className="min-h-screen bg-neutral-50 p-6">
      <div className="max-w-6xl mx-auto space-y-4">
        <h1 className="text-2xl font-bold">Payments • ክፍያዎች</h1>
        <div className="rounded-2xl border bg-white overflow-hidden">
          <div className="grid grid-cols-7 gap-4 p-4 bg-neutral-50 text-xs font-semibold text-muted-foreground">
            <span>Tx Ref</span><span>Amount</span><span>Method</span><span>Status</span><span>Routing</span><span>2FA</span><span>Action</span>
          </div>
          {mockPayments.map(p=> (
            <div key={p.id} className="grid grid-cols-7 gap-4 p-4 border-t hover:bg-neutral-50 text-sm">
              <span className="font-mono text-xs">{p.tx_ref}</span>
              <span className="font-semibold">ETB {p.amount}</span>
              <span>{p.method} <span className="text-[10px] text-muted-foreground">{p.connector_id}</span></span>
              <span><span className={`px-2 py-0.5 rounded-full text-xs ${p.status==="succeeded" ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700"}`}>{p.status}</span></span>
              <span className="text-xs">{p.routing}</span>
              <span className="text-xs">{p.requires_2fa ? (p.two_fa_verified ? "✓ verified" : "pending") : "-"}</span>
              <Link href={`/payments/${p.id}`} className="text-primary text-xs font-medium">View • እይ</Link>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
