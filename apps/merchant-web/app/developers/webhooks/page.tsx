"use client"
import * as React from "react"
import Link from "next/link"
import { Loader2, Plus, RotateCcw } from "lucide-react"
import { api, WebhookEndpoint, WebhookDelivery } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function WebhooksPage() {
  const { t } = useLanguage()
  const { data: endpoints, loading: eload, error: eerr, refetch: erefetch } = useData<WebhookEndpoint[]>(
    () => api.developer.webhookEndpoints(), []
  )
  const { data: deliveries, loading: dload, refetch: drefetch } = useData<WebhookDelivery[]>(
    () => api.developer.webhookDeliveries(), []
  )

  const [url, setUrl] = React.useState("")
  const [secret, setSecret] = React.useState("")
  const [events, setEvents] = React.useState("payment.succeeded,payment.failed")
  const [creating, setCreating] = React.useState(false)
  const [created, setCreated] = React.useState("")

  const create = async () => {
    setCreating(true)
    setCreated("")
    try {
      const ev = events.split(",").map((s) => s.trim()).filter(Boolean)
      const res = await api.developer.createWebhookEndpoint({ url, secret, events: ev })
      setCreated(res.id)
      setUrl(""); setSecret("")
      await erefetch()
    } catch (e) {
      alert((e as Error).message || "Failed to create endpoint")
    } finally {
      setCreating(false)
    }
  }

  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-5xl mx-auto space-y-6">
        <Link href="/developers" className="text-sm text-primary">← Developers</Link>
        <div className="flex justify-between items-center">
          <h1 className="text-2xl font-bold">{t("Webhooks", "ዌብሁኮች")}</h1>
          <button onClick={() => { erefetch(); drefetch(); }} className="rounded-xl border bg-card px-3 h-9 text-xs flex items-center gap-1">
            <RotateCcw className="h-3 w-3" /> Refresh
          </button>
        </div>

        <div className="rounded-2xl border bg-card p-6 space-y-3">
          <h3 className="font-semibold">Add an endpoint</h3>
          {created && <p className="text-sm text-green-600">Created endpoint {created}. Secret shown to you at creation only.</p>}
          <input value={url} onChange={(e) => setUrl(e.target.value)} placeholder="https://merchant.example.et/webhook (https required, SSRF-safe)"
            className="w-full rounded-xl border h-10 px-3 text-sm" />
          <input value={secret} onChange={(e) => setSecret(e.target.value)} placeholder="Signing secret (>= 16 chars)"
            className="w-full rounded-xl border h-10 px-3 text-sm" />
          <input value={events} onChange={(e) => setEvents(e.target.value)} placeholder="event1,event2"
            className="w-full rounded-xl border h-10 px-3 text-sm" />
          <button onClick={create} disabled={creating || !url || secret.length < 16}
            className="rounded-xl bg-primary text-foreground px-5 h-10 text-sm font-semibold disabled:opacity-50 flex items-center gap-1">
            <Plus className="h-4 w-4" /> {creating ? "Creating…" : "Add Endpoint"}
          </button>
        </div>

        <div className="rounded-2xl border bg-card p-6">
          <h3 className="font-semibold mb-3">Endpoints</h3>
          {eload && <p className="text-xs text-muted-foreground">Loading…</p>}
          {eerr && <p className="text-xs text-red-600">{eerr}</p>}
          {!eload && !eerr && (endpoints ?? []).length === 0 && <p className="text-xs text-muted-foreground">No endpoints yet.</p>}
          <div className="space-y-2">
            {(endpoints ?? []).map((e) => (
              <div key={e.id} className="flex items-center justify-between rounded-xl border p-3 text-sm">
                <span className="font-mono text-xs">{e.url}</span>
                <span className={`text-xs px-2 py-0.5 rounded-full ${e.status === "active" ? "bg-green-500/20 text-green-700" : "bg-amber-500/20"}`}>{e.status}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="rounded-2xl border bg-card p-6">
          <h3 className="font-semibold mb-3">Recent deliveries</h3>
          {dload && <p className="text-xs text-muted-foreground">Loading…</p>}
          {!dload && (deliveries ?? []).length === 0 && <p className="text-xs text-muted-foreground">No deliveries yet.</p>}
          <div className="space-y-2">
            {(deliveries ?? []).map((d) => (
              <div key={d.id} className="flex items-center justify-between rounded-xl border p-3 text-xs">
                <span>{d.event_type} <span className="text-muted-foreground">• attempt {d.attempt_count} • HTTP {d.last_status_code}</span></span>
                <span className={`px-2 py-0.5 rounded-full ${d.status === "success" ? "bg-green-500/20 text-green-700" : d.status === "pending" ? "bg-amber-500/20" : "bg-red-100 text-red-700"}`}>{d.status}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
