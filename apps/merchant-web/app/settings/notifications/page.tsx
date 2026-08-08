"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type NotifyPref } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function NotificationPrefsPage() {
  const { checking } = useRequireAuth()
  const { data, refetch } = useData(() => api.notificationsPrefs.list(), [])

  if (checking) return <Centered>Checking session…</Centered>

  const toggle = async (p: NotifyPref, key: "email" | "sms" | "push" | "in_app") => {
    await api.notificationsPrefs.update({ ...p, [key]: !p[key] })
    refetch()
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-4xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Notification Preferences • ማሳወቂያ</h1>
        <p className="text-sm text-muted-foreground">Choose which channels you receive for each event type.</p>

        <div className="rounded-2xl border bg-card overflow-hidden">
          <div className="grid grid-cols-5 gap-2 bg-muted p-3 text-[11px] font-semibold">
            <span>Event</span><span className="text-center">Email</span><span className="text-center">SMS</span><span className="text-center">Push</span><span className="text-center">In-app</span>
          </div>
          {(data ?? []).map((p) => (
            <div key={p.event_type} className="grid grid-cols-5 gap-2 p-3 border-t text-sm">
              <span className="text-xs">{p.event_type.replace(/_/g, " ")}</span>
              {(["email", "sms", "push", "in_app"] as const).map((ch) => (
                <span key={ch} className="text-center">
                  <button
                    onClick={() => toggle(p, ch)}
                    className={`h-6 w-6 rounded-full inline-block ${p[ch] ? "bg-primary" : "bg-neutral-200"}`}
                    aria-label={ch}
                  />
                </span>
              ))}
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
