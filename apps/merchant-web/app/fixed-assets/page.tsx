"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type FixedAsset } from "@/lib/api/client"
import { useData } from "@/lib/api/use-data"

export default function FixedAssetsPage() {
  const { checking } = useRequireAuth()
  const { data, refetch } = useData(() => api.fixedAssets.list(), [])

  const [name, setName] = React.useState("")
  const [category, setCategory] = React.useState("equipment")
  const [cost, setCost] = React.useState("")
  const [life, setLife] = React.useState("5")
  const [method, setMethod] = React.useState("straight_line")
  const [err, setErr] = React.useState("")
  const [saving, setSaving] = React.useState(false)

  if (checking) return <Centered>Checking session…</Centered>

  const today = new Date().toISOString().slice(0, 10)

  const create = async () => {
    setSaving(true); setErr("")
    try {
      await api.fixedAssets.create({ asset_name: name, category, acquisition_date: today, cost, salvage_value: "0", useful_life_years: Number(life), depreciation_method: method })
      setName(""); setCost(""); refetch()
    } catch (e) { setErr((e as Error).message) } finally { setSaving(false) }
  }

  const depreciate = async (id: string) => { await api.fixedAssets.depreciate(id); refetch() }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Fixed Assets & Depreciation • ቋሚ ንብረቶች</h1>
        <p className="text-sm text-muted-foreground">Straight-line and declining-balance depreciation for tax.</p>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <h3 className="font-semibold">Add Asset</h3>
            <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Asset name" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <select value={category} onChange={(e) => setCategory(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
              {["building", "machinery", "vehicle", "equipment", "furniture", "computer", "land", "other"].map((c) => <option key={c}>{c}</option>)}
            </select>
            <input value={cost} onChange={(e) => setCost(e.target.value)} placeholder="Cost ETB" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <input value={life} onChange={(e) => setLife(e.target.value)} placeholder="Useful life (years)" className="w-full rounded-xl border h-11 px-3 text-sm" />
            <select value={method} onChange={(e) => setMethod(e.target.value)} className="w-full rounded-xl border h-11 px-3 text-sm">
              <option value="straight_line">Straight line</option>
              <option value="declining_balance">Declining balance</option>
            </select>
            {err && <p className="text-sm text-red-600">{err}</p>}
            <button onClick={create} disabled={saving} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold disabled:opacity-50">Add Asset</button>
          </div>

          <div className="lg:col-span-2 rounded-2xl border bg-card overflow-hidden">
            <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold">
              <span>Asset</span><span>Category</span><span>Cost</span><span>NBV</span><span>Method</span><span>Action</span>
            </div>
            {(data ?? []).map((a) => (
              <div key={a.id} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                <span className="font-medium">{a.asset_name}</span>
                <span>{a.category}</span>
                <span>ETB {a.cost}</span>
                <span className="font-semibold">ETB {a.net_book_value}</span>
                <span>{a.depreciation_method}</span>
                <span><button onClick={() => depreciate(a.id)} className="text-primary">Depreciate</button></span>
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
