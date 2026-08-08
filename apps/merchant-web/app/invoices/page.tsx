"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type Invoice, type AgingBucket } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function InvoicesPage() {
  const { checking } = useRequireAuth()
  const { data: invoices, refetch } = useData(() => api.invoices.list(), [])
  const { data: aging } = useData(() => api.invoices.aging(), [])

  const [customerName, setCustomerName] = React.useState("")
  const [customerEmail, setCustomerEmail] = React.useState("")
  const [invNumber, setInvNumber] = React.useState("")
  const [amount, setAmount] = React.useState("")
  const [taxPct, setTaxPct] = React.useState("15")
  const [withholdPct, setWithholdPct] = React.useState("2")
  const [err, setErr] = React.useState("")
  const [saving, setSaving] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const today = new Date().toISOString().slice(0, 10)
  const due = new Date(Date.now() + 30 * 86400000).toISOString().slice(0, 10)

  const create = async () => {
    setSaving(true); setErr("")
    try {
      await api.invoices.create({
        invoice_number: invNumber, customer_name: customerName, customer_email: customerEmail,
        issue_date: today, due_date: due, currency: "ETB", tax_percent: taxPct, withholding_percent: withholdPct,
        line_items: [{ description: "Service", quantity: "1", unit_price: amount }],
      })
      setInvNumber(""); setCustomerName(""); setAmount(""); refetch()
    } catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Invoices • Receivables • ኢንቮይስ</h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Create Invoice</h3>
            <input value={invNumber} onChange={(e) => setInvNumber(e.target.value)} placeholder="Invoice #" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={customerName} onChange={(e) => setCustomerName(e.target.value)} placeholder="Customer name" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={customerEmail} onChange={(e) => setCustomerEmail(e.target.value)} placeholder="Customer email" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="Amount ETB" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <div className="grid grid-cols-2 gap-2">
              <input value={taxPct} onChange={(e) => setTaxPct(e.target.value)} placeholder="Tax % (VAT)" className="w-full rounded-xl border h-11 px-3 text-sm" />
              <input value={withholdPct} onChange={(e) => setWithholdPct(e.target.value)} placeholder="Withholding %" className="w-full rounded-xl border h-11 px-3 text-sm" />
            </div>
            {err && <p className="text-sm text-red-600">{err}</p>}
            <button onClick={create} disabled={saving} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">
              {saving ? "Creating…" : "Create Invoice"}
            </button>
          </div>

          <div className="lg:col-span-2 rounded-2xl border bg-card overflow-hidden">
            <div className="flex justify-between items-center p-4 border-b">
              <h3 className="font-semibold">Invoices</h3>
              <div className="flex gap-2 text-[11px]">
                {(aging ?? []).map((b) => (
                  <span key={b.bucket} className="px-2 py-1 rounded-full bg-muted">{b.bucket}: {b.count} • ETB {b.amount}</span>
                ))}
              </div>
            </div>
            <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>Number</span><span>Customer</span><span>Amount</span><span>Due</span><span>Status</span><span>Action</span>
            </div>
            {(invoices ?? []).map((inv) => (
              <div key={inv.id} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                <span className="font-mono text-[10px]">{inv.invoice_number}</span>
                <span>{inv.customer_name}</span>
                <span className="font-semibold">ETB {inv.total_amount}</span>
                <span>{inv.due_date}</span>
                <span className={`px-2 py-0.5 rounded-full text-[11px] ${inv.status === "paid" ? "bg-green-500/15 text-green-700" : inv.status === "overdue" ? "bg-red-500/15 text-red-700" : "bg-amber-500/15 text-amber-700"}`}>{inv.status}</span>
                <span><button onClick={async () => { await api.invoices.issue(inv.id); refetch() }} className="text-primary text-[11px]">Issue link</button></span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
