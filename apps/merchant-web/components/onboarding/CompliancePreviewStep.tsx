"use client"
import * as React from "react"
import { Card } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"

export function CompliancePreviewStep({ data }: { data: any }) {
  const checks = [
    { type: "tin_validation", label: "TIN Validation • ቲን", status: "passed", score: 95, details: "Format 10 digits valid" },
    { type: "business_license_validation", label: "Business License • ፈቃድ", status: "passed", score: 90, details: "Not expired" },
    { type: "bank_account_validation", label: "Bank Account • ባንክ", status: "passed", score: 88, details: "Name match Levenshtein 1" },
    { type: "fayda_verification", label: "Fayda Verification • ፋይዳ", status: data.owners?.[0]?.fayda_verified ? "passed" : "pending", score: data.owners?.[0]?.fayda_verified ? 92 : 0, details: data.owners?.[0]?.fayda_verified ? "OTP verified, face 0.92" : "Awaiting OTP" },
    { type: "restricted_industry", label: "Restricted Industry • የተከለከለ ዘርፍ", status: "passed", score: 100, details: "E-commerce allowed" },
    { type: "website_policy_check", label: "Website Policy • ድህረ ገጽ ፖሊሲ", status: "needs_review", score: 70, details: "Refund/Privacy/Terms found via crawl" },
    { type: "aml_screening", label: "AML Screening • አደጋ", status: "passed", score: 85, details: "No sanctions hit" },
    { type: "risk_scoring", label: "Risk Scoring • አደጋ ነጥብ", status: "passed", score: 42, details: "Medium risk TPV 500k" },
  ]

  const avgScore = Math.round(checks.reduce((a, c) => a + c.score, 0) / checks.length)

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold">Compliance Preview • ተገዢነት ቅድመ-እይታ</h2>
        <p className="text-sm text-muted-foreground">Automated checks per NBE ONPS + PayAtlas + Chapa model, risk_score weighted sum optimal.</p>
      </div>

      <Card className="p-4 bg-gradient-to-br from-primary-50 to-white">
        <div className="flex items-center justify-between">
          <div><p className="text-sm text-muted-foreground">Risk Score • አደጋ</p><p className="text-3xl font-bold">{avgScore}/100</p><p className="text-xs">Medium • TPV 500k ETB</p></div>
          <div className="h-20 w-20 rounded-full border-8 border-primary flex items-center justify-center font-bold text-primary" style={{ borderColor: avgScore < 50 ? "#10B981" : avgScore < 75 ? "#F59E0B" : "#EF4444" }}>{avgScore}</div>
        </div>
        <Progress value={avgScore} className="mt-3" />
      </Card>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {checks.map(c => (
          <Card key={c.type} className={`p-3 border ${c.status === "passed" ? "border-green-200 bg-green-50/30" : c.status === "needs_review" ? "border-amber-200 bg-amber-50/30" : "border-black/10"}`}>
            <div className="flex items-center justify-between">
              <p className="text-sm font-medium">{c.label}</p>
              <span className={`text-xs px-2 py-0.5 rounded-full ${c.status === "passed" ? "bg-green-100 text-green-700" : c.status === "needs_review" ? "bg-amber-100 text-amber-700" : "bg-neutral-100"}`}>{c.status}</span>
            </div>
            <p className="text-xs text-muted-foreground mt-1">{c.details}</p>
            <Progress value={c.score} className="mt-2 h-1.5" />
          </Card>
        ))}
      </div>

      <div className="rounded-xl bg-blue-50 border border-blue-200 p-3 text-xs">
        <p className="font-semibold">Next Steps:</p>
        <ul className="list-disc list-inside mt-1">
          <li>If risk high (&gt;=70) or TPV &gt;1M ETB =&gt; dual approval required maker-checker</li>
          <li>After approval: merchant operating book created + ledger_accounts seeded (merchant_operating, liability:merchant_payable, etc)</li>
          <li>Test API keys auto-created, live keys after pilot 30-60 days per NBE analog</li>
          <li>Outbox merchant.activated event triggers email confetti outstanding</li>
        </ul>
      </div>
    </div>
  )
}
