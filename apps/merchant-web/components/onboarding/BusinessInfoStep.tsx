"use client"
import * as React from "react"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"

export function BusinessInfoStep({ data, onChange }: { data: any, onChange: (d:any)=>void }) {
  const update = (k:string,v:any)=> onChange({ ...data, [k]: v })
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold">Business Info • የንግድ መረጃ</h2>
        <p className="text-sm text-muted-foreground">Per NBE ONPS/02/2020 § capital & shareholder requirements + PayAtlas KYC</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2"><Label>Legal Name • ህጋዊ ስም *</Label><Input placeholder="Apex Trading PLC" value={data.legal_name||""} onChange={e=>update("legal_name", e.target.value)} /></div>
        <div className="space-y-2"><Label>Trade Name • የንግድ ስም</Label><Input placeholder="ApexPay" value={data.trade_name||""} onChange={e=>update("trade_name", e.target.value)} /></div>
        <div className="space-y-2"><Label>Business Type • የንግድ አይነት *</Label>
          <select className="w-full rounded-xl border border-black/10 h-12 px-3" value={data.business_type||"plc"} onChange={e=>update("business_type", e.target.value)}>
            <option value="sole_proprietorship">Sole Proprietorship</option><option value="plc">PLC</option><option value="share_company">Share Company (min 5 shareholders)</option><option value="partnership">Partnership</option>
          </select>
        </div>
        <div className="space-y-2"><Label>TIN Number • ቲን ቁጥር 10-digit *</Label><Input placeholder="0023456789" maxLength={10} value={data.tin_number||""} onChange={e=>update("tin_number", e.target.value)} /></div>
        <div className="space-y-2"><Label>Registration No • ምዝገባ ቁጥር *</Label><Input placeholder="MT/AA/..." value={data.registration_number||""} onChange={e=>update("registration_number", e.target.value)} /></div>
        <div className="space-y-2"><Label>Industry • ዘርፍ *</Label>
          <select className="w-full rounded-xl border border-black/10 h-12 px-3" value={data.industry||"e-commerce"} onChange={e=>update("industry", e.target.value)}>
            <option value="e-commerce">E-commerce</option><option value="education">Education</option><option value="tech">Tech</option><option value="food">Food</option><option value="gambling" disabled>🚫 Gambling (Prohibited per NBE)</option><option value="crypto" disabled>🚫 Crypto (Prohibited)</option><option value="adult" disabled>🚫 Adult (Prohibited)</option>
          </select>
        </div>
        <div className="md:col-span-2 space-y-2"><Label>Business Description • የንግድ መግለጫ *</Label><textarea className="w-full rounded-xl border border-black/10 p-3 min-h-[80px]" placeholder="We sell ..." value={data.business_description||""} onChange={e=>update("business_description", e.target.value)} /></div>
        <div className="space-y-2"><Label>Website URL • ድህረ ገጽ</Label><Input placeholder="https://example.et" value={data.website_url||""} onChange={e=>update("website_url", e.target.value)} /></div>
        <div className="space-y-2"><Label>Expected Monthly TPV • ወርሃዊ ግምት ETB *</Label><Input type="number" placeholder="500000" value={data.expected_monthly_tpv||""} onChange={e=>update("expected_monthly_tpv", e.target.value)} /></div>
        <div className="space-y-2"><Label>Region • ክልል *</Label><Input placeholder="Addis Ababa" value={data.region||""} onChange={e=>update("region", e.target.value)} /></div>
        <div className="space-y-2"><Label>City • ከተማ *</Label><Input placeholder="Addis Ababa" value={data.city||""} onChange={e=>update("city", e.target.value)} /></div>
        <div className="md:col-span-2 space-y-2"><Label>Full Office Address • ሙሉ አድራሻ *</Label><Input placeholder="Bole, Woreda 03, House No..." value={data.address_full||""} onChange={e=>update("address_full", e.target.value)} /></div>
      </div>

      <div className="rounded-xl bg-amber-50 border border-amber-200 p-3 text-xs">
        <p className="font-semibold">NBE Requirements • መስፈርቶች:</p>
        <ul className="list-disc list-inside mt-1 space-y-0.5">
          <li>Company Registration notarized Amharic/English</li>
          <li>Owner max 40% single person per directive</li>
          <li>Website must have refund, privacy, terms pages (PayAtlas)</li>
          <li>Risk scoring if TPV &gt;1M ETB =&gt; dual approval</li>
        </ul>
      </div>
    </div>
  )
}

function Input(props:any){ return <input {...props} className="w-full rounded-xl border border-black/10 h-12 px-3 focus:ring-2 focus:ring-primary focus:outline-none" /> }
function Label(props:any){ return <label {...props} className="text-sm font-medium" /> }
function Select(props:any){ return <select {...props} /> }
