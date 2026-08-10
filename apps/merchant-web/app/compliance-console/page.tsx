"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type ComplianceStatus } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function ComplianceConsolePage() {
  const { t } = useLanguage()
  const { checking } = useRequireAuth()
  const { data: status } = useData(() => api.compliance.status(), [])
  const { data: checks } = useData(() => api.compliance.checks(), [])

  if (checking) return <Centered>Checking session…</Centered>

  const rows = (s: ComplianceStatus | null): Array<{ label: string; value: string; status: string }> => {
    if (!s) return []
    const overdue = (d?: string) => (d && d < new Date().toISOString().slice(0, 10)) ? "overdue" : "ok"
    return [
      { label: "KYC Expiry", value: s.kyc_expiry_date || "—", status: overdue(s.kyc_expiry_date) },
      { label: "Business License", value: s.license_expiry || "—", status: overdue(s.license_expiry) },
      { label: "ERCA Filing", value: s.next_erca_due || "—", status: overdue(s.next_erca_due) },
      { label: "Pension", value: s.next_pension_due || "—", status: overdue(s.next_pension_due) },
      { label: "Annual Tax Filing", value: s.annual_tax_filing_due || "—", status: overdue(s.annual_tax_filing_due) },
      { label: "AML", value: s.aml_due || "—", status: overdue(s.aml_due) },
    ]
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">{t("Compliance Console","ተገዢነት")}</h1>
        <p className="text-sm text-muted-foreground">KYC, license, tax, pension and AML obligations with expiry tracking.</p>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6">
            <h3 className="font-semibold">Overall Status</h3>
            <p className={`text-2xl font-bold mt-3 ${status?.overall_status === "good" ? "text-green-600" : status?.overall_status === "overdue" ? "text-red-600" : "text-amber-600"}`}>
              {(status?.overall_status ?? "—").toUpperCase()}
            </p>
            <p className="text-xs text-muted-foreground mt-1">Onboarding: {status?.onboarding_status}</p>
            <p className="text-xs text-muted-foreground">Risk tier: {status?.risk_tier} • Fayda: {status?.fayda_verified ? "✓" : "✗"}</p>
          </div>

          <div className="lg:col-span-2 rounded-2xl border bg-card overflow-hidden">
            <div className="grid grid-cols-3 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>Obligation</span><span>Due Date</span><span>Status</span>
            </div>
            {rows(status ?? null).map((r) => (
              <div key={r.label} className="grid grid-cols-3 gap-2 p-3 border-t text-sm">
                <span>{r.label}</span>
                <span>{r.value}</span>
                <span className={`text-[11px] px-2 py-0.5 rounded-full w-fit ${r.status === "ok" ? "bg-green-500/15 text-green-700" : "bg-red-500/15 text-red-700"}`}>{r.status === "ok" ? "OK" : "OVERDUE"}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="rounded-2xl border bg-card p-4">
          <h3 className="font-semibold">Recent Compliance Checks</h3>
          <div className="mt-3 space-y-2">
            {(checks ?? []).slice(0, 10).map((c) => (
              <div key={c.id} className="rounded-xl border p-3 flex justify-between text-sm">
                <span>{c.check_type}</span>
                <span className={`text-[11px] px-2 py-0.5 rounded-full ${c.status === "passed" ? "bg-green-500/15 text-green-700" : c.status === "failed" ? "bg-red-500/15 text-red-700" : "bg-amber-500/15 text-amber-700"}`}>{c.status}</span>
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
