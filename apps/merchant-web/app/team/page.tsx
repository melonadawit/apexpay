"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type Member, type Approval } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

const ROLES = ["owner", "admin", "developer", "finance", "support", "ops", "compliance", "viewer"]

export default function TeamPage() {
  const { checking } = useRequireAuth()
  const { data: members, refetch: refetchMembers } = useData(() => api.team.members(), [])
  const { data: approvals, refetch: refetchApprovals } = useData(() => api.team.approvals(), [])

  const [email, setEmail] = React.useState("")
  const [name, setName] = React.useState("")
  const [role, setRole] = React.useState("viewer")
  const [err, setErr] = React.useState("")
  const [saving, setSaving] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const invite = async () => {
    setSaving(true); setErr("")
    try { await api.team.invite({ email, name, role, permissions: [] }); setEmail(""); setName(""); refetchMembers() }
    catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  const decide = async (id: string, decision: string) => {
    await api.team.decide(id, decision); refetchApprovals()
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Team & Approvals • ቡድን</h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Invite Member</h3>
            <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Name" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Email" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <select value={role} onChange={(e) => setRole(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
              {ROLES.map((r) => <option key={r} value={r}>{r}</option>)}
            </select>
            {err && <p className="text-sm text-red-600">{err}</p>}
            <button onClick={invite} disabled={saving} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">
              {saving ? "Inviting…" : "Invite"}
            </button>
          </div>

          <div className="rounded-2xl border bg-card overflow-hidden">
            <h3 className="font-semibold p-4">Members</h3>
            {(members ?? []).map((m) => (
              <div key={m.user_id} className="border-t p-3 flex justify-between text-sm">
                <div>
                  <p className="font-medium">{m.name}</p>
                  <p className="text-xs text-muted-foreground">{m.email}</p>
                </div>
                <span className="text-[11px] px-2 py-0.5 rounded-full bg-primary/10 text-primary">{m.role}</span>
              </div>
            ))}
          </div>

          <div className="rounded-2xl border bg-card overflow-hidden">
            <h3 className="font-semibold p-4">Approvals Inbox</h3>
            {(approvals ?? []).slice(0, 10).map((a) => (
              <div key={a.id} className="border-t p-3 text-xs space-y-1">
                <p className="font-medium">{a.summary || `${a.resource_type} ${a.resource_id}`}</p>
                <p className="text-muted-foreground">ETB {a.amount} • {a.action} • {a.status} • {a.approvals.length}/{a.required_approvals}</p>
                {a.status === "pending" && (
                  <div className="flex gap-2 mt-1">
                    <button onClick={() => decide(a.id, "approve")} className="rounded-xl bg-green-600 text-white px-3 h-7 text-[11px]">Approve</button>
                    <button onClick={() => decide(a.id, "reject")} className="rounded-xl bg-red-600 text-white px-3 h-7 text-[11px]">Reject</button>
                  </div>
                )}
              </div>
            ))}
            {(approvals ?? []).length === 0 && <p className="p-4 text-sm text-muted-foreground">No pending approvals.</p>}
          </div>
        </div>
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
