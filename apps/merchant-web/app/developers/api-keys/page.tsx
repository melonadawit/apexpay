"use client"
import * as React from "react"
import { Loader2, Plus, RotateCcw } from "lucide-react"
import { api, DeveloperApiKey } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"
import { useLanguage } from "@/components/providers/language-provider"

export default function ApiKeysPage() {
  const { t } = useLanguage()
  const { data: keys, loading, error, refetch } = useData<DeveloperApiKey[]>(
    () => api.developer.apiKeys(),
    []
  )
  const [revealedSecret, setRevealedSecret] = React.useState("")
  const [name, setName] = React.useState("")
  const [env, setEnv] = React.useState("test")
  const [creating, setCreating] = React.useState(false)
  const [busy, setBusy] = React.useState("")

  const createKey = async () => {
    if (!name.trim()) return
    setCreating(true)
    try {
      const res = await api.developer.createApiKey({ name: name.trim(), environment: env })
      setRevealedSecret(res.secret)
      setName("")
      await refetch()
    } catch (e) {
      alert((e as Error).message || "Failed to create key")
    } finally {
      setCreating(false)
    }
  }

  const revoke = async (id: string) => {
    if (!confirm("Revoke this API key? It will stop working immediately.")) return
    setBusy(id)
    try {
      await api.developer.revokeApiKey(id)
      await refetch()
    } catch (e) {
      alert((e as Error).message || "Failed to revoke")
    } finally {
      setBusy("")
    }
  }

  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-5xl mx-auto space-y-6">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold">{t("API Keys", "የ API ቁልፎች")}</h1>
            <p className="text-sm text-muted-foreground">Test/Live keys, scoped, hash-at-rest. Secrets shown once.</p>
          </div>
          <button onClick={() => refetch()} className="rounded-xl border bg-card px-3 h-9 text-xs flex items-center gap-1">
            <RotateCcw className="h-3 w-3" /> Refresh
          </button>
        </div>

        {revealedSecret && (
          <div className="rounded-2xl bg-amber-500/10 border border-amber-500/30 p-4">
            <p className="text-sm font-semibold">New key created — copy it now, it will not be shown again.</p>
            <code className="block mt-2 rounded-lg bg-background border p-3 font-mono text-sm break-all">{revealedSecret}</code>
          </div>
        )}

        <div className="rounded-2xl border bg-card p-6 space-y-4">
          <h3 className="font-semibold">Create a key</h3>
          <div className="flex gap-3 flex-wrap items-end">
            <div className="flex-1 min-w-[200px]">
              <label className="text-xs text-muted-foreground">Name</label>
              <input value={name} onChange={(e) => setName(e.target.value)} placeholder="e.g. mobile app"
                className="w-full rounded-xl border h-10 px-3 text-sm mt-1" />
            </div>
            <div>
              <label className="text-xs text-muted-foreground">Environment</label>
              <select value={env} onChange={(e) => setEnv(e.target.value)}
                className="rounded-xl border h-10 px-3 text-sm mt-1">
                <option value="test">Test</option>
                <option value="live">Live</option>
              </select>
            </div>
            <button onClick={createKey} disabled={creating || !name.trim()}
              className="rounded-xl bg-primary text-foreground px-5 h-10 text-sm font-semibold disabled:opacity-50 flex items-center gap-1">
              <Plus className="h-4 w-4" /> {creating ? "Creating…" : "Create Key"}
            </button>
          </div>
        </div>

        <div className="rounded-2xl border bg-card p-6">
          <h3 className="font-semibold mb-3">Your keys</h3>
          {loading && <div className="flex items-center gap-2 text-sm text-muted-foreground"><Loader2 className="h-4 w-4 animate-spin" /> Loading…</div>}
          {error && <p className="text-sm text-red-600">{error}</p>}
          {!loading && !error && (keys ?? []).length === 0 && (
            <p className="text-sm text-muted-foreground">No API keys yet.</p>
          )}
          <div className="space-y-3">
            {(keys ?? []).map((k) => (
              <div key={k.id} className="rounded-xl border p-4 flex items-center justify-between">
                <div>
                  <p className="font-medium text-sm">{k.name} <span className="font-mono text-xs text-muted-foreground">• {k.key_prefix}•••••••• • {k.environment} • {k.key_type}</span></p>
                  <p className="text-xs text-muted-foreground mt-0.5">Scopes: {k.scopes.join(", ") || "—"}{k.last_used_at ? ` • Last used ${new Date(k.last_used_at).toLocaleString()}` : ""}</p>
                </div>
                <div className="flex items-center gap-2">
                  <span className={`text-xs px-2 py-0.5 rounded-full ${k.status === "active" ? "bg-green-500/20 text-green-700" : "bg-red-100 text-red-700"}`}>{k.status}</span>
                  {k.status === "active" && (
                    <button onClick={() => revoke(k.id)} disabled={busy === k.id}
                      className="rounded-lg border px-3 py-1 text-xs text-red-600 disabled:opacity-50">
                      {busy === k.id ? "…" : "Revoke"}
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="rounded-2xl border bg-card p-6">
          <h3 className="font-semibold">Public Key & Embedded SDK</h3>
          <p className="text-xs text-muted-foreground mt-1">Use the public key with checkout.js — safe for the frontend, no secret.</p>
          <pre className="mt-3 rounded-xl bg-muted text-foreground p-4 font-mono text-[11px] overflow-auto">{`<script src="https://checkout.apexpay.et/sdk.js"></script>
<script>
  const apexpay = new ApexPay('pk_test_...');
  apexpay.checkout({ amount: '500', currency: 'ETB', tx_ref: 'txr_' + Date.now(),
    method: 'telebirr', customer_email: 'cust@example.et',
    return_url: 'https://example.et/return', callback_url: 'https://example.et/callback' });
</script>`}</pre>
        </div>
      </div>
    </div>
  )
}
