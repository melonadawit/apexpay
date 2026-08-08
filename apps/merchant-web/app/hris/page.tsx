"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type Team, type Contract, type Shift } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function HrisPage() {
  const { checking } = useRequireAuth()
  const { data: teams, refetch: refetchTeams } = useData(() => api.hris.teams(), [])
  const { data: contracts } = useData(() => api.hris.contracts(), [])
  const { data: shifts } = useData(() => api.hris.shifts(), [])
  const { data: attendance } = useData(() => api.hris.attendance(), [])

  const [teamName, setTeamName] = React.useState("")
  const [err, setErr] = React.useState("")
  const [saving, setSaving] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const addTeam = async () => {
    setSaving(true); setErr("")
    try { await api.hris.createTeam({ name: teamName }); setTeamName(""); refetchTeams() }
    catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Workforce OS • HRIS</h1>
        <p className="text-sm text-muted-foreground">Teams, contracts, shifts, attendance clocking, onboarding, reviews.</p>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Teams</h3>
            <div className="flex gap-2">
              <input value={teamName} onChange={(e) => setTeamName(e.target.value)} placeholder="Team name" className="flex-1 rounded-xl border h-11 px-3 text-sm" />
              <button onClick={addTeam} disabled={saving} className="rounded-xl bg-primary text-white px-4 text-sm disabled:opacity-50">Add</button>
            </div>
            {err && <p className="text-sm text-red-600">{err}</p>}
            <div className="space-y-2">
              {(teams ?? []).map((t) => (
                <div key={t.id} className="rounded-xl border p-3">
                  <p className="text-sm font-medium">{t.name}</p>
                  <p className="text-xs text-muted-foreground">{t.description || "—"}</p>
                </div>
              ))}
            </div>
            <h4 className="font-semibold pt-2">Shifts</h4>
            <div className="space-y-1">
              {(shifts ?? []).map((s) => <p key={s.id} className="text-xs">{s.name} • {s.start_time}-{s.end_time}</p>)}
            </div>
          </div>

          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Contracts</h3>
            <div className="space-y-2">
              {(contracts ?? []).map((c) => (
                <div key={c.id} className="rounded-xl border p-3">
                  <p className="text-sm font-medium">{c.contract_type} • {c.start_date}{c.end_date ? ` → ${c.end_date}` : ""}</p>
                  <p className="text-xs text-muted-foreground">ETB {c.salary_amount} • {c.status}</p>
                </div>
              ))}
            </div>
          </div>

          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Attendance Today</h3>
            <div className="space-y-2">
              {(attendance ?? []).slice(0, 10).map((a) => (
                <div key={a.id} className="rounded-xl border p-3 flex justify-between">
                  <div>
                    <p className="text-sm font-medium">{a.employee_id}</p>
                    <p className="text-xs text-muted-foreground">{a.clock_date} • {a.hours}h</p>
                  </div>
                  <span className={`text-[11px] px-2 py-0.5 rounded-full ${a.status === "present" ? "bg-green-500/15 text-green-700" : "bg-amber-500/15 text-amber-700"}`}>{a.status}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
