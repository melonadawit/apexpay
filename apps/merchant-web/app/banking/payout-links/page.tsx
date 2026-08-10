"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type PayoutLink } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function PayoutLinksPage() {
  const { t } = useLanguage()
  const { checking } = useRequireAuth()
  const { data, loading, refetch } = useData(() => api.banking.payoutLinks(), [])

  const [amount, setAmount] = React.useState("")
  const [recipientName, setRecipientName] = React.useState("")
  const [recipientPhone, setRecipientPhone] = React.useState("")
  const [purpose, setPurpose] = React.useState("")
  const [created, setCreated] = React.useState<PayoutLink | null>(null)
  const [error, setError] = React.useState("")
  const [creating, setCreating] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const create = async () => {
    setCreating(true); setError("")
    try {
      const res = await api.banking.createPayoutLink({
        amount, currency: "ETB",
        recipient_name: recipientName, recipient_phone: recipientPhone,
        purpose,
      })
      setCreated(res)
      setAmount(""); setRecipientName(""); setRecipientPhone(""); setPurpose("")
      refetch()
    } catch (e) { setError((e as Error).message) } finally { setCreating(false) }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold">{t("Payout Links","የክፍያ ሊንኮች")}</h1>
          <p className="text-sm text-muted-foreground mt-2">
            QR + public links — recipients claim without entering bank details; escrow-backed.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Create Payout Link</h3>
            <input value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="Amount ETB" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={recipientName} onChange={(e) => setRecipientName(e.target.value)} placeholder="Recipient name" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={recipientPhone} onChange={(e) => setRecipientPhone(e.target.value)} placeholder="Recipient phone" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={purpose} onChange={(e) => setPurpose(e.target.value)} placeholder="Purpose (refund/cashback)" className="w-full rounded-xl border h-11 px-3 text-sm" />
            {error && <p className="text-sm text-red-600">{error}</p>}
            <button onClick={create} disabled={creating} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">
              {creating ? "Creating…" : "Create Payout Link"}
            </button>
            {created && (
              <div className="rounded-xl bg-muted p-3 text-center">
                <p className="text-xs text-muted-foreground">Link created</p>
                <p className="mt-1 font-mono text-xs break-all">Token: {created.public_token}</p>
              </div>
            )}
          </div>

          <div className="lg:col-span-2 rounded-2xl border bg-card overflow-hidden">
            <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>Token</span><span>Amount</span><span>Recipient</span><span>Purpose</span><span>Status</span><span>Expires</span>
            </div>
            {loading && <p className="p-4 text-sm text-muted-foreground">Loading…</p>}
            {(data ?? []).map((p) => (
              <div key={p.id} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                <span className="font-mono text-[10px]">{p.public_token}</span>
                <span className="font-semibold">ETB {p.amount}</span>
                <span>{p.recipient_name || "—"}</span>
                <span>{p.purpose || "—"}</span>
                <span className={`px-2 py-0.5 rounded-full text-[11px] ${p.status === "active" ? "bg-blue-500/15 text-blue-700" : "bg-green-500/15 text-green-700"}`}>{p.status}</span>
                <span>{p.expires_at}</span>
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
