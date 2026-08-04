"use client"
import * as React from "react"
import { Button } from "@/components/ui/button"

export default function LinksPage() {
  const [amount, setAmount] = React.useState("500")
  const [desc, setDesc] = React.useState("Tutoring • አስተማሪ")
  const [link, setLink] = React.useState("")
  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-6xl mx-auto grid grid-cols-3 gap-6">
        <div className="col-span-1 rounded-2xl border bg-card p-6 space-y-4">
          <h2 className="font-bold">Create Link • ሊንክ ፍጠር — Outstanding</h2>
          <div className="flex flex-wrap gap-2">{[100,500,1000,5000].map(a=> <button key={a} onClick={()=> setAmount(a.toString())} className={`px-3 py-1 rounded-full text-xs border ${amount===a.toString() ? "bg-primary text-foreground border-primary" : "bg-card border-border"}`}>ETB {a}</button>)}</div>
          <input value={amount} onChange={e=> setAmount(e.target.value)} placeholder="Amount ETB" className="w-full rounded-xl border h-12 px-3" />
          <input value={desc} onChange={e=> setDesc(e.target.value)} placeholder="Description" className="w-full rounded-xl border h-12 px-3" />
          <Button className="w-full" onClick={()=> setLink(`https://checkout.apexpay.et/c/abc${Date.now()}`)}>Generate Link • አመንጭ</Button>
          <p className="text-[11px] text-muted-foreground">✓ QR preview live • ✓ Share Telegram/WhatsApp system share • ✓ Public token unique • ✓ Expires_at optional • ✓ Haptics</p>
        </div>
        <div className="col-span-2 space-y-4">
          {link && <div className="rounded-2xl border bg-card p-6 text-center"><div className="mx-auto h-40 w-40 bg-muted/80 rounded-xl flex items-center justify-center">QR</div><p className="mt-3 font-mono text-sm">{link}</p><div className="mt-3 flex gap-2 justify-center"><Button variant="outline">Copy • ቅዳ</Button><Button onClick={()=> (navigator as any).share && (navigator as any).share({ title: desc, text: `Pay ETB ${amount} for ${desc}`, url: link })}>Share • አጋራ via Telegram/WhatsApp</Button></div></div>}
          <div className="rounded-2xl border bg-card p-4">
            <h3 className="font-semibold">Links List • {3} • QR thumbnails</h3>
            <div className="mt-3 space-y-2">
              {[
                {id:"pl_01H", amount:"500", desc:"Tutoring", status:"active", token:"abc123"},
                {id:"pl_02H", amount:"1000", desc:"Coffee", status:"paid", token:"def456"},
              ].map(l=>(
                <div key={l.id} className="flex items-center justify-between rounded-xl border p-3">
                  <div className="flex items-center gap-3"><div className="h-12 w-12 bg-muted/80 rounded-lg flex items-center justify-center">QR</div><div><p className="text-sm font-medium">ETB {l.amount} • {l.desc}</p><p className="text-xs text-muted-foreground">{l.token} • {l.status}</p></div></div>
                  <span className={`text-xs px-2 py-0.5 rounded-full ${l.status==="paid" ? "bg-green-500/20 text-green-700" : "bg-blue-500/20 text-blue-700"}`}>{l.status}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
