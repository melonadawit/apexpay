"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

export default function RelationshipManagersPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Relationship Managers • Dedicated RM • Priority Support • One RM Per Merchant • Outstanding • RazorpayX Parity • P0 • RM • SLA • Support Tickets</h1>
        <Card className="p-6">
          <p className="text-sm">Dedicated Relationship Manager per RazorpayX: Dedicated Relationship Manager to help you with a seamless banking experience, Priority Support Services Benefit from prompt issue resolution and priority support ensuring smooth banking operations. Merchant ID mer_01H • RM User ID user_rm_001 • Assigned At 2026-01-15 • Assigned By Admin • Status Active • One RM Per Merchant • Unique Merchant ID • Outstanding per RazorpayX.</p>
          <div className="mt-4 rounded-xl border p-4">
            <p className="font-medium text-sm">RM: Abebe Kebede • Senior RM • abebe@apextrading.et • Phone +251911111111 • Assigned At 2026-01-15 • Assigned By Admin • Status Active • One RM Per Merchant</p>
            <p className="text-[11px] text-muted-foreground mt-1">Merchant ID mer_01H • RM User ID user_rm_001 • Assigned At 2026-01-15T10:00:00Z • Assigned By Admin • Status Active • One RM Per Merchant • Unique Merchant ID • One RM Per Merchant • Outstanding per RazorpayX Dedicated Relationship Manager + Priority Support</p>
          </div>
        </Card>
      </div>
    </div>
  )
}
