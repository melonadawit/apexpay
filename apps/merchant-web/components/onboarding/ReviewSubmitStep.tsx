"use client"
import * as React from "react"
import { Button } from "@/components/ui/button"
import { Card, GlassCard } from "@/components/ui/card"
import { Check, ShieldCheck, FileText, Building, CreditCard } from "lucide-react"

export function ReviewSubmitStep({ data }: { data: any }) {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold flex items-center gap-2"><ShieldCheck className="text-primary" /> Review & Submit • ግምገማ እና አስገባ</h2>
        <p className="text-sm text-muted-foreground">Confirm per NBE ONPS/02/2020 and consent Fayda verification via id.gov.et with OTP consent.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card className="p-4 space-y-2">
          <div className="flex items-center gap-2 font-semibold"><Building size={16} /> Business</div>
          <p className="text-sm">Legal: {data.legal_name || "Apex Trading PLC"}</p>
          <p className="text-sm">TIN: {data.tin_number || "0023456789"}</p>
          <p className="text-sm">Industry: {data.industry || "e-commerce"} • Region: {data.region || "Addis Ababa"}</p>
        </Card>
        <Card className="p-4 space-y-2">
          <div className="flex items-center gap-2 font-semibold"><Check size={16} /> Owners & Fayda</div>
          <p className="text-sm">{data.owners?.[0]?.full_name || "Abebe Kebede"} • 100% • Auth Signatory</p>
          <p className="text-sm">FIN ****-{data.owners?.[0]?.fin_last4 || "1234"} • Verified {data.owners?.[0]?.fayda_verified ? "✓ 0.92" : "pending"}</p>
        </Card>
        <Card className="p-4 space-y-2">
          <div className="flex items-center gap-2 font-semibold"><CreditCard size={16} /> Bank</div>
          <p className="text-sm">{data.bank_code || "CBE"} • {data.account_name || "Apex Trading PLC"}</p>
          <p className="text-sm">****{data.account_number?.slice(-4) || "1234"} • Default settlement</p>
        </Card>
        <Card className="p-4 space-y-2">
          <div className="flex items-center gap-2 font-semibold"><FileText size={16} /> Documents</div>
          <p className="text-sm">{Object.keys(data.uploadedDocs || {}).length} docs uploaded • {Object.keys(data.uploadedDocs || {}).join(", ") || "company_registration, tin_certificate, ..."}</p>
        </Card>
      </div>

      <GlassCard className="p-4 bg-amber-500/10/70">
        <h4 className="font-semibold text-sm">Consent & Declarations • ስምምነት</h4>
        <ul className="mt-2 space-y-2 text-xs">
          <li className="flex gap-2"><input type="checkbox" defaultChecked /> <span>I confirm business info true per NBE ONPS/02/2020 directive and capacity to manage payment gateway per capital ETB 3M requirement</span></li>
          <li className="flex gap-2"><input type="checkbox" defaultChecked /> <span>I consent Fayda verification via id.gov.et VeriFayda 2.0 / OIDC eSignet with OTP consent, understanding FIN stored as sha256 hash + last4 only, front/back images encrypted at rest AES-256, presigned URLs 15m TTL, no plain PII in logs</span></li>
          <li className="flex gap-2"><input type="checkbox" defaultChecked /> <span>I agree website has refund, privacy, terms pages per PayAtlas ET PSP requirement, and business not in restricted industries gambling/crypto/adult</span></li>
          <li className="flex gap-2"><input type="checkbox" defaultChecked /> <span>I understand risk scoring medium, dual approval if high risk or TPV&gt;1M ETB, test keys immediately, live keys after pilot 30-60 days analogy NBE, 2FA mandatory &gt;5000 ETB per ONPS/10/2025</span></li>
        </ul>
      </GlassCard>

      <div className="rounded-xl bg-muted p-3 text-xs text-muted-foreground">
        <p>After submit: compliance team reviews in Kanban board outstanding (Submitted → In Review → Fayda Pending → Compliance Check → Approved). Timeline vertical like Linear. Email with confetti animation. Ledger: merchant operating book created + accounts seeded. Outbox merchant.activated.</p>
      </div>
    </div>
  )
}
