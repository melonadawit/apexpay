"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type VendorInvoice } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function VendorPaymentsPage() {
  const { checking } = useRequireAuth()
  const { data, loading, refetch } = useData(() => api.banking.vendorInvoices(), [])

  const [invoiceNumber, setInvoiceNumber] = React.useState("")
  const [vendorName, setVendorName] = React.useState("")
  const [amount, setAmount] = React.useState("")
  const [taxAmount, setTaxAmount] = React.useState("")
  const [withholding, setWithholding] = React.useState("")
  const [error, setError] = React.useState("")
  const [creating, setCreating] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const total = (Number(amount || 0) + Number(taxAmount || 0) - Number(withholding || 0)).toFixed(2)

  const create = async () => {
    setCreating(true)
    setError("")
    try {
      await api.banking.createVendorInvoice({
        invoice_number: invoiceNumber,
        invoice_date: new Date().toISOString().slice(0, 10),
        vendor_name: vendorName,
        amount,
        currency: "ETB",
        tax_amount: taxAmount || "0",
        withholding_tax_amount: withholding || "0",
        total_amount: total,
        status: "pending_approval",
        ocr_confidence: 0,
      })
      setInvoiceNumber(""); setVendorName(""); setAmount(""); setTaxAmount(""); setWithholding("")
      refetch()
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setCreating(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Vendor Payments • አቅራቢ ክፍያዎች</h1>
          <p className="text-sm text-muted-foreground mt-2">
            Accounts payable with VAT 15% + withholding 2% auto-calculation and approval flow.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Upload Invoice</h3>
            <input value={invoiceNumber} onChange={(e) => setInvoiceNumber(e.target.value)} placeholder="Invoice #" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={vendorName} onChange={(e) => setVendorName(e.target.value)} placeholder="Vendor name" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="Amount ETB" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <div className="grid grid-cols-2 gap-2">
              <input value={taxAmount} onChange={(e) => setTaxAmount(e.target.value)} placeholder="VAT 15%" className="w-full rounded-xl border h-11 px-3 text-sm" />
              <input value={withholding} onChange={(e) => setWithholding(e.target.value)} placeholder="Withholding 2%" className="w-full rounded-xl border h-11 px-3 text-sm" />
            </div>
            <p className="text-xs text-muted-foreground">Total: ETB {total}</p>
            {error && <p className="text-sm text-red-600">{error}</p>}
            <button onClick={create} disabled={creating} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">
              {creating ? "Creating…" : "Upload Invoice"}
            </button>
          </div>

          <div className="lg:col-span-2 rounded-2xl border bg-card overflow-hidden">
            <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>Invoice</span><span>Vendor</span><span>Amount</span><span>VAT</span><span>Withholding</span><span>Status</span>
            </div>
            {loading && <p className="p-4 text-sm text-muted-foreground">Loading…</p>}
            {(data ?? []).map((v) => (
              <div key={v.id} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                <span className="font-mono text-[10px]">{v.invoice_number}</span>
                <span>{v.vendor_name || "—"}</span>
                <span className="font-semibold">ETB {v.amount}</span>
                <span>ETB {v.tax_amount}</span>
                <span>ETB {v.withholding_tax_amount}</span>
                <span className={`px-2 py-0.5 rounded-full text-[11px] ${v.status === "paid" ? "bg-green-500/15 text-green-700" : v.status === "approved" ? "bg-blue-500/15 text-blue-700" : "bg-amber-500/15 text-amber-700"}`}>{v.status}</span>
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
