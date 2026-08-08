"use client"
import * as React from "react"
import { Button } from "@/components/ui/button"
import { api } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

type LinkResult = {
  id: string
  amount: string
  currency: string
  description?: string
  status: string
  public_token: string
  checkout_url: string
  share?: { telegram: string; whatsapp: string }
}

export default function LinksPage() {
  const [amount, setAmount] = React.useState("500")
  const [desc, setDesc] = React.useState("Tutoring • አስተማሪ")
  const [created, setCreated] = React.useState<LinkResult | null>(null)
  const [creating, setCreating] = React.useState(false)
  const [error, setError] = React.useState("")

  const { data: links, loading, refetch } = useData(
    () => api.links(),
    []
  )

  const create = async () => {
    setCreating(true)
    setError("")
    try {
      const res = (await api.createLink({ amount, currency: "ETB", description: desc })) as LinkResult
      setCreated(res)
      refetch()
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setCreating(false)
    }
  }

  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-6xl mx-auto grid grid-cols-3 gap-6">
        <div className="col-span-1 rounded-2xl border bg-card p-6 space-y-4">
          <h2 className="font-bold">Create Link • ሊንክ ፍጠር</h2>
          <div className="flex flex-wrap gap-2">
            {[100, 500, 1000, 5000].map((a) => (
              <button
                key={a}
                onClick={() => setAmount(a.toString())}
                className={`px-3 py-1 rounded-full text-xs border ${
                  amount === a.toString() ? "bg-primary text-foreground border-primary" : "bg-card border-border"
                }`}
              >
                ETB {a}
              </button>
            ))}
          </div>
          <input
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            placeholder="Amount ETB"
            className="w-full rounded-xl border h-12 px-3"
          />
          <input
            value={desc}
            onChange={(e) => setDesc(e.target.value)}
            placeholder="Description"
            className="w-full rounded-xl border h-12 px-3"
          />
          <Button className="w-full" onClick={create} disabled={creating}>
            {creating ? "Creating…" : "Generate Link • አመንጭ"}
          </Button>
          {error && <p className="text-sm text-red-600">{error}</p>}

          {created && (
            <div className="rounded-2xl border bg-card p-4 text-center">
              <p className="text-xs text-muted-foreground">Link created</p>
              <p className="mt-1 font-mono text-xs break-all">{created.checkout_url}</p>
              <div className="mt-3 flex gap-2 justify-center">
                <Button variant="outline" onClick={() => navigator.clipboard.writeText(created.checkout_url)}>
                  Copy • ቅዳ
                </Button>
                {created.share && (
                  <Button
                    onClick={() => {
                      navigator.clipboard.writeText(created.share!.whatsapp)
                      window.open(created.share!.whatsapp, "_blank")
                    }}
                  >
                    Share • አጋራ
                  </Button>
                )}
              </div>
            </div>
          )}
        </div>

        <div className="col-span-2 space-y-4">
          <div className="rounded-2xl border bg-card p-4">
            <h3 className="font-semibold">Links List • {(links ?? []).length}</h3>
            <div className="mt-3 space-y-2">
              {loading && <p className="text-sm text-muted-foreground">Loading…</p>}
              {!loading && (links ?? []).length === 0 && (
                <p className="text-sm text-muted-foreground">No payment links yet.</p>
              )}
              {(links ?? []).map((l) => (
                <div key={l.id} className="flex items-center justify-between rounded-xl border p-3">
                  <div>
                    <p className="text-sm font-medium">ETB {l.amount} • {l.description || "—"}</p>
                    <p className="text-xs text-muted-foreground">{l.public_token}</p>
                  </div>
                  <span
                    className={`text-xs px-2 py-0.5 rounded-full ${
                      l.status === "paid"
                        ? "bg-green-500/20 text-green-700"
                        : l.status === "active"
                        ? "bg-blue-500/20 text-blue-700"
                        : "bg-neutral-200 text-neutral-700"
                    }`}
                  >
                    {l.status}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
