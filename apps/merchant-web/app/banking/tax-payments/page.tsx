"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type TaxPayment } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

const TAX_TYPES = ["vat", "tot", "withholding", "paye", "pension"]

export default function TaxPaymentsPage() {
  const { checking } = useRequireAuth()
  const { data, loading, refetch } = useData(() => api.banking.taxPayments(), [])
  const [taxType, setTaxType] = React.useState("vat")
  const [amount, setAmount] = React.useState("")
  const [periodMonth, setPeriodMonth] = React.useState("")
  const [periodYear, setPeriodYear] = React.useState("")
  const [dueDate, setDueDate] = React.useState("")
  const [error, setError] = React.useState("")
  const [creating, setCreating] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const create = async () => {
    setCreating(true)
    setError("")
    try {
      await api.banking.createTaxPayment({
        tax_type: taxType,
        amount,
        currency: "ETB",
        period_month: periodMonth ? Number(periodMonth) : null,
        period_year: periodYear ? Number(periodYear) : null,
        due_date: dueDate,
        status: "draft",
      })
      setAmount("")
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
          <h1 className="text-3xl font-bold">Tax Payments • የግብር ክፍያዎች</h1>
          <p className="text-sm text-muted-foreground mt-2">
            VAT 15% • TOT 2%/10% • Withholding 2% • PAYE • Pension — automated pre-filled forms.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Create Tax Payment</h3>
            <select value={taxType} onChange={(e) => setTaxType(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
              {TAX_TYPES.map((t) => <option key={t} value={t}>{t}</option>)}
            </select>
            <input value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="Amount ETB" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <div className="grid grid-cols-2 gap-2">
              <input value={periodMonth} onChange={(e) => setPeriodMonth(e.target.value)} placeholder="Month (1-12)" className="w-full rounded-xl border h-11 px-3 text-sm" />
              <input value={periodYear} onChange={(e) => setPeriodYear(e.target.value)} placeholder="Year" className="w-full rounded-xl border h-11 px-3 text-sm" />
            </div>
            <input value={dueDate} onChange={(e) => setDueDate(e.target.value)} type="date" className="w-full rounded-xl border h-11 px-3 text-sm" />
            {error && <p className="text-sm text-red-600">{error}</p>}
            <button onClick={create} disabled={creating} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">
              {creating ? "Creating…" : "Create Tax Payment"}
            </button>
          </div>

          <div className="lg:col-span-2 rounded-2xl border bg-card overflow-hidden">
            <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>Type</span><span>Amount</span><span>Period</span><span>Due</span><span>Status</span><span>Ref</span>
            </div>
            {loading && <p className="p-4 text-sm text-muted-foreground">Loading…</p>}
            {(data ?? []).map((t) => (
              <div key={t.id} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                <span>{t.tax_type}</span>
                <span className="font-semibold">ETB {t.amount}</span>
                <span>{t.period_month ? `${t.period_month}/${t.period_year}` : "—"}</span>
                <span>{t.due_date || "—"}</span>
                <span className={`px-2 py-0.5 rounded-full text-[11px] ${t.status === "paid" ? "bg-green-500/15 text-green-700" : "bg-amber-500/15 text-amber-700"}`}>{t.status}</span>
                <span className="text-[10px]">{t.payment_reference || "—"}</span>
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
