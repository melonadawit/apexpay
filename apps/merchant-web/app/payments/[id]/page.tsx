"use client"
import * as React from "react"
import { Loader2 } from "lucide-react"
import { api } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"
import { formatETB } from "@/lib/api/use-data"

export default function PaymentDetailPage({ params }: { params: { id: string } }) {
  const { t } = useLanguage()
  const { data, loading, error } = useData(() => api.payment(params.id), [params.id])

  const p = data?.payment
  const journals = data?.journals ?? []

  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-4xl mx-auto space-y-4">
        <h1 className="text-2xl font-bold">
          {t("Payment Exam", "የክፍያ ምርመራ")} • {params.id}
        </h1>
        <p className="text-sm text-muted-foreground">
          NBE exam console — lifecycle + ledger journals + connector refs for a single transaction.
        </p>

        {loading && (
          <div className="flex items-center gap-2 p-6 text-sm text-muted-foreground">
            <Loader2 className="h-4 w-4 animate-spin" /> Loading…
          </div>
        )}
        {error && <p className="p-4 text-sm text-red-600">{error}</p>}

        {!loading && !error && p && (
          <div className="grid grid-cols-3 gap-4">
            <div className="col-span-2 space-y-4">
              <div className="rounded-2xl border bg-card p-4">
                <h3 className="font-semibold">{t("Lifecycle", "የህይወት ዑደት")}</h3>
                <div className="mt-3 relative pl-6 border-l-2 border-neutral-200 space-y-3">
                  <div className="relative text-xs">
                    <div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-primary border-2 border-white" />
                    <p>created → {p.status} • tx_ref {p.tx_ref} • via connector {p.connector_id || "—"} (ref {p.connector_ref || "—"})</p>
                  </div>
                  <div className="relative text-xs">
                    <div className="absolute -left-[29px] top-1 h-3 w-3 rounded-full bg-green-500 border-2 border-white" />
                    <p>amount {formatETB(p.amount)} {p.currency} • fee {formatETB(p.fee_amount)} • net {formatETB(p.net_amount)}</p>
                  </div>
                </div>
              </div>

              <div className="rounded-2xl border bg-card p-4">
                <h3 className="font-semibold">{t("Ledger Journal", "መዝገብ")}</h3>
                {journals.length === 0 && (
                  <p className="mt-2 text-xs text-muted-foreground">No ledger journals posted for this transaction yet.</p>
                )}
                {journals.map((j) => (
                  <div key={j.id} className="mt-3">
                    <p className="text-xs font-mono text-muted-foreground">
                      {j.posting_key} • {j.memo}
                    </p>
                    <div className="mt-2 space-y-1 text-xs font-mono">
                      <div className="grid grid-cols-4 gap-2 bg-muted p-2 rounded">
                        <span>Account</span><span>Direction</span><span>Amount ETB</span><span>Name</span>
                      </div>
                      {j.entries.map((e, i) => (
                        <div key={i} className="grid grid-cols-4 gap-2 p-2 border-b">
                          <span>{e.account_code}</span>
                          <span className={e.direction === "debit" ? "text-green-600" : "text-red-600"}>{e.direction}</span>
                          <span>{formatETB(e.amount)}</span>
                          <span className="text-muted-foreground">{e.account_name}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="space-y-4">
              <div className="rounded-2xl border bg-card p-4">
                <h4 className="font-semibold text-sm">Transaction</h4>
                <div className="mt-2 space-y-1 text-xs">
                  <p>Status: <span className="font-semibold">{p.status}</span></p>
                  <p>Method: {p.method || "—"}</p>
                  <p>Connector: {p.connector_id || "—"}</p>
                  <p>2FA: {p.requires_2fa ? (p.two_fa_verified ? "✓ verified" : "pending") : "not required"}</p>
                  <p>Checkout: <span className="font-mono break-all">{p.checkout_url}</span></p>
                </div>
              </div>

              <div className="rounded-2xl border bg-card p-4">
                <h4 className="font-semibold text-sm">{t("Actions", "እርምጃዎች")}</h4>
                <div className="mt-2 grid grid-cols-2 gap-2">
                  <a href={`/refunds/${params.id}`} className="rounded-xl border h-10 text-xs grid place-items-center">Refund • ተመላሽ</a>
                  <a href={`/payments`} className="rounded-xl border h-10 text-xs grid place-items-center">Evidence Pack JSON</a>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
