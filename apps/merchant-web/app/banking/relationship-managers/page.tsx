"use client"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function RelationshipManagersPage() {
  const { t } = useLanguage()
  const { checking } = useRequireAuth()
  const { data, loading } = useData(() => api.banking.relationshipManagers(), [])

  if (checking) return <Centered>Checking session…</Centered>

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-4xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">{t("Relationship Managers","የግንኙነት አስተዳዳሪዎች")}</h1>
        <div className="rounded-2xl border bg-card overflow-hidden">
          <div className="grid grid-cols-4 gap-2 bg-muted p-3 text-[11px] font-semibold">
            <span>ID</span><span>RM User</span><span>Status</span><span>Assigned</span>
          </div>
          {loading && <p className="p-4 text-sm text-muted-foreground">Loading…</p>}
          {(data ?? []).map((rm) => (
            <div key={rm.id} className="grid grid-cols-4 gap-2 p-3 border-t text-xs hover:bg-muted/50">
              <span className="font-mono text-[10px]">{rm.id}</span>
              <span>{rm.rm_user_id}</span>
              <span className={`px-2 py-0.5 rounded-full text-[11px] ${rm.status === "active" ? "bg-green-500/15 text-green-700" : "bg-amber-500/15 text-amber-700"}`}>{rm.status}</span>
              <span>{rm.assigned_at}</span>
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
