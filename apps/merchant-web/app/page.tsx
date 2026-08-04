"use client"
import Link from "next/link"
import { motion } from "framer-motion"
import { Button } from "@/components/ui/button"
import { Card, GlassCard } from "@/components/ui/card"

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 via-white to-primary-50">
      <header className="sticky top-0 z-20 glass border-b">
        <div className="max-w-6xl mx-auto px-6 py-4 flex justify-between items-center">
          <div className="font-bold text-xl">ApexPay • አፔክስፔይ</div>
          <div className="flex gap-2">
            <Link href="/onboarding"><Button>Start Onboarding • መዝገብ</Button></Link>
          </div>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-6 py-12 space-y-12">
        <motion.div initial={{ opacity:0, y:20 }} animate={{ opacity:1, y:0 }} className="text-center space-y-6">
          <h1 className="text-5xl font-bold tracking-tight">AI-native payment gateway for Ethiopia • ለኢትዮጵያ</h1>
          <p className="text-muted-foreground max-w-2xl mx-auto">Collect via Telebirr, CBE Birr, Bank IPS, EthSwitch QR, Cards — with Fayda ID verification front/back + OTP, NBE-grade onboarding, smart routing, payouts, payroll ET tax, RAG compliance, Swarm AI.</p>
          <div className="flex justify-center gap-3">
            <Link href="/onboarding"><Button size="lg">Outstanding Onboarding Wizard • ጀምር</Button></Link>
            <Button size="lg" variant="outline">View Demo • ተመልከት</Button>
          </div>
        </motion.div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <GlassCard className="p-6"><h3 className="font-semibold">Fayda ID Verification • ፋይዳ</h3><p className="text-sm text-muted-foreground mt-2">Front/back images &lt;2MB, selfie liveness, OTP consent via id.gov.et VeriFayda 2.0, offline QR FaydaEncode, OIDC eSignet, FIN hashed sha256 + last4 only.</p></GlassCard>
          <GlassCard className="p-6"><h3 className="font-semibold">Smart Routing • ዘመናዊ መንገድ</h3><p className="text-sm text-muted-foreground mt-2">Health sampler 30s, circuit breaker 5 fails open 60s, success_rate/latency/cost strategies, fallback audit trail outstanding.</p></GlassCard>
          <GlassCard className="p-6"><h3 className="font-semibold">Payroll ET • ደሞዝ</h3><p className="text-sm text-muted-foreground mt-2">Binary search O(log n) tax brackets, pension 7%/11%, OT 1.25/1.5/2.0, payslip PDF modern, ledger M4 per run book balanced.</p></GlassCard>
        </div>

        <Card className="p-6">
          <h3 className="font-semibold">NBE Compliance Checklist • የተሟላ</h3>
          <ul className="list-disc list-inside text-sm mt-3 space-y-1 text-muted-foreground">
            <li>Company Registration notarized, TIN 10-digit, Business License not expired, VAT, Memorandum, Board Resolution, Shareholder list (PayAtlas + NBE ONPS/02/2020)</li>
            <li>Fayda front/back + selfie + OTP consent, encrypted MinIO vault, presigned 15m, hash integrity, no plain FIN in logs</li>
            <li>Settlement bank account name == legal name fuzzy Levenshtein &lt;3, bank letter required</li>
            <li>Website refund/privacy/terms required, 2FA mandatory &gt;5000 ETB per ONPS/10/2025, AML ETB 200k reporting FIC</li>
          </ul>
        </Card>
      </main>
    </div>
  )
}
