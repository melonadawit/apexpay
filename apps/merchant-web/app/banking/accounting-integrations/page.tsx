"use client"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function AccountingIntegrationsPage() {
  const { checking } = useRequireAuth()
  const { data, loading } = useData(() => api.banking.accountingIntegrations(), [])

  if (checking) return <Centered>Checking session…</Centered>

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Accounting Integrations • የሂሳብ ውህደት</h1>
        <p className="text-sm text-muted-foreground">Tally • Zoho • QuickBooks — two-way sync with CA access controls.</p>

        <div className="rounded-2xl border bg-card overflow-hidden">
          <div className="grid grid-cols-5 gap-2 bg-muted p-3 text-[11px] font-semibold">
            <span>Provider</span><span>Status</span><span>Last Sync</span><span>Sync Status</span><span>Error</span>
          </div>
          {loading && <p className="p-4 text-sm text-muted-foreground">Loading…</p>}
          {(data ?? []).map((i) => (
            <div key={i.id} className="grid grid-cols-5 gap-2 p-3 border-t text-xs hover:bg-muted/50">
              <span className="font-medium capitalize">{i.provider}</span>
              <span className={`px-2 py-0.5 rounded-full text-[11px] ${i.status === "connected" ? "bg-green-500/15 text-green-700" : "bg-amber-500/15 text-amber-700"}`}>{i.status}</span>
              <span>{i.last_sync_status ? i.created_at : "—"}</span>
              <span>{i.last_sync_status || "—"}</span>
              <span className="text-[10px] text-red-600">{i.last_sync_error || "—"}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
