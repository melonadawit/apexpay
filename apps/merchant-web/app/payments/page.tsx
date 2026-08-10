"use client"
import * as React from "react"
import Link from "next/link"
import { api } from "@/lib/api/client"
import { useData, formatETB } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function PaymentsPage() {
  const { t } = useLanguage()
  const { data: payments, loading } = useData(() => api.payments(50), [])

  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-6xl mx-auto space-y-4">
        <h1 className="text-2xl font-bold">{t("Payments","ክፍያዎች")}</h1>
        <div className="rounded-2xl border bg-card overflow-hidden">
          <div className="grid grid-cols-7 gap-4 p-4 bg-muted text-xs font-semibold text-muted-foreground">
            <span>Tx Ref</span><span>Amount</span><span>Method</span><span>Status</span><span>Connector</span><span>2FA</span><span>Action</span>
          </div>
          {loading && <p className="p-4 text-sm text-muted-foreground">Loading…</p>}
          {!loading && (payments ?? []).length === 0 && (
            <p className="p-4 text-sm text-muted-foreground">No payments found.</p>
          )}
          {(payments ?? []).map((p) => (
            <div key={p.id} className="grid grid-cols-7 gap-4 p-4 border-t hover:bg-muted text-sm">
              <span className="font-mono text-xs">{p.tx_ref}</span>
              <span className="font-semibold">ETB {formatETB(p.amount)}</span>
              <span>{p.method || "—"}</span>
              <span>
                <span
                  className={`px-2 py-0.5 rounded-full text-xs ${
                    p.status === "succeeded"
                      ? "bg-green-500/20 text-green-700"
                      : p.status === "failed"
                      ? "bg-red-100 text-red-700"
                      : "bg-amber-100 text-amber-700"
                  }`}
                >
                  {p.status}
                </span>
              </span>
              <span className="text-xs">{p.connector_id}</span>
              <span className="text-xs">
                {p.requires_2fa ? (p.two_fa_verified ? "✓ verified" : "pending") : "-"}
              </span>
              <Link href={`/payments/${p.id}`} className="text-primary text-xs font-medium">View • እይ</Link>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
