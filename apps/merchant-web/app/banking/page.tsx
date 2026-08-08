"use client"
import Link from "next/link"

const modules = [
  { href: "/banking/current-accounts", title: "Current Accounts", desc: "Partner-bank accounts & balances" },
  { href: "/banking/credit-lines", title: "Credit Lines", desc: "Collateral-free credit" },
  { href: "/banking/forex", title: "Forex", desc: "Rates & FDI transfers" },
  { href: "/banking/escrow", title: "Escrow", desc: "Marketplace escrow" },
  { href: "/banking/corporate-cards", title: "Corporate Cards", desc: "Virtual + physical cards" },
  { href: "/banking/virtual-accounts", title: "Virtual Accounts", desc: "Smart collect" },
  { href: "/banking/vendor-payments", title: "Vendor Payments", desc: "Accounts payable" },
  { href: "/banking/tax-payments", title: "Tax Payments", desc: "VAT/TOT/withholding" },
  { href: "/banking/petty-cash", title: "Petty Cash", desc: "Budgets & expenses" },
  { href: "/banking/payout-links", title: "Payout Links", desc: "QR payout links" },
  { href: "/banking/bank-verification", title: "Bank Verification", desc: "Penny testing" },
  { href: "/banking/support-tickets", title: "Support", desc: "Priority support & SLA" },
]

export default function BankingIndexPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-5xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Banking • ባንክ</h1>
        <p className="text-sm text-muted-foreground">Business banking modules.</p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {modules.map((m) => (
            <Link key={m.href} href={m.href} className="rounded-2xl border bg-card p-5 hover:shadow-medium transition-shadow">
              <h3 className="font-semibold">{m.title}</h3>
              <p className="text-xs text-muted-foreground mt-1">{m.desc}</p>
            </Link>
          ))}
        </div>
      </div>
    </div>
  )
}
