"use client"
import * as React from "react"
import Link from "next/link"
import { useLanguage } from "@/components/providers/language-provider"

export default function DevelopersPage() {
  const { t } = useLanguage()
  const [showSecret, setShowSecret] = React.useState(false)
  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <div className="flex justify-between items-center">
          <h1 className="text-2xl font-bold">{t("Developers","ገንቢዎች")}</h1>
          <div className="flex gap-2">
            <Link href="/developers/api-keys" className="rounded-xl border bg-card px-4 h-10 text-xs grid place-items-center font-medium">API Keys</Link>
            <Link href="/developers/webhooks" className="rounded-xl border bg-card px-4 h-10 text-xs grid place-items-center font-medium">Webhooks</Link>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-6">
          <div className="rounded-2xl border bg-card p-6 space-y-4">
            <h3 className="font-semibold">API Keys • Test/Live Separate Scopes Reveal Once</h3>
            <div className="rounded-xl border p-3 bg-muted">
              <p className="text-xs text-muted-foreground">Secret Key • sk_test_ — shown once hash at rest per DATABASE</p>
              <p className="font-mono text-sm mt-1">{showSecret ? "sk_test_51Hq...abc123xyz" : "sk_test_••••••••••••••••"}</p>
              <div className="mt-2 flex gap-2"><button onClick={()=> setShowSecret(!showSecret)} className="rounded-lg border px-3 py-1 text-xs">{showSecret ? "Hide" : "Reveal"} • sk_live shown once</button><button className="rounded-lg bg-primary text-foreground px-3 py-1 text-xs">Create New Key • test/live separate live only after KYC active</button></div>
              <p className="text-[11px] text-muted-foreground mt-2">Prefix unique index api_keys_prefix_uidx O(1) lookup • secret_hash index • scopes [] • last_used_at async best effort non-blocking Go routine</p>
            </div>

            <div>
              <h4 className="font-semibold text-sm">Public Key • pk_test_</h4>
              <p className="font-mono text-xs mt-1">pk_test_51Hq... + checkout.js embedded SDK</p>
            </div>

            <div className="rounded-xl bg-blue-500/10 border border-blue-500/20 p-3 text-xs">
              <p className="font-semibold">Quickstart 6 lines copy-paste like ApexPay "simple as copy and paste":</p>
              <pre className="mt-2 bg-background text-foreground rounded-lg p-3 overflow-auto text-[11px]">{`npm install apexpay-js
import ApexPay from 'apexpay-js'
import { useLanguage } from "@/components/providers/language-provider"
const apexpay = new ApexPay('sk_test_...')
apexpay.initialize({ tx_ref: 'txr_01', amount: '500', currency: 'ETB', method: 'telebirr' })
  .then(res => console.log(res.checkout_url))
  .catch(err => console.error(err.code))
// Verify: GET /v1/transactions/verify/{tx_ref}
// Webhook: HMAC SHA256 X-ApexPay-Signature
// Fayda: POST /v1/onboarding/fayda/verify/init FIN 12-digit front/back <2MB selfie OTP mock 123456
// Banks: GET /v1/banks CBE/Awash/Dashen
// Methods: GET /v1/methods?amount=1000 ranked score 0.6*success+0.4*(1-latency/1000)
// RAG: POST /v1/compliance/ask query "When is 2FA required?" lang en/am top_k 5 threshold 0.65 guard
// Swarm: POST /v1/swarm/run goal "Create link 100 ETB for coffee if today TPV>0"
`}</pre>
            </div>
          </div>

          <div className="space-y-6">
            <div className="rounded-2xl border bg-card p-4">
              <h3 className="font-semibold">Webhook Endpoints • HMAC + SSRF Block + Retry Exponential Backoff</h3>
              <input placeholder="https://merchant.example.et/webhook" className="mt-3 w-full rounded-xl border h-10 px-3 text-sm" defaultValue="https://merchant.example.et/webhook" />
              <div className="mt-2 flex gap-2">
                <label className="text-xs flex items-center gap-1"><input type="checkbox" defaultChecked /> payment.succeeded</label>
                <label className="text-xs flex items-center gap-1"><input type="checkbox" defaultChecked /> refund.succeeded</label>
                <label className="text-xs flex items-center gap-1"><input type="checkbox" /> payout.succeeded</label>
              </div>
              <button className="mt-3 w-full rounded-xl bg-primary text-foreground h-10 text-sm">Save Endpoint • SSRF block private ranges Allowlist 10.0.0.0/8 172... 192.168 127.0.0.1</button>
              <div className="mt-3 space-y-1 text-xs">
                <div className="flex justify-between border-b py-1"><span>payment.succeeded → https://.../webhook</span><span className="px-2 py-0.5 rounded-full bg-green-500/20">success 200 attempt1 HMAC valid</span></div>
                <div className="flex justify-between border-b py-1"><span>refund.succeeded</span><span className="px-2 py-0.5 rounded-full bg-amber-500/20">failed retry backoff 2s 4s 8s</span></div>
              </div>
            </div>

            <div className="rounded-2xl border bg-card p-4">
              <h3 className="font-semibold">Banks List • GET /v1/banks + Methods Health • GET /v1/methods?amount=1000 ranked</h3>
              <div className="mt-3 grid grid-cols-2 gap-2 text-xs">
                {[
                  ["CBE", "Commercial Bank of Ethiopia", "latency 180ms success 0.92"],
                  ["AWASH", "Awash Bank", "latency 210ms success 0.95"],
                  ["DASHEN", "Dashen Bank", "latency 260ms success 0.89"],
                  ["ABYSSINIA", "Bank of Abyssinia", "latency 200ms success 0.96 circuit closed"],
                ].map(([code,name,health])=>(
                  <div key={code} className="rounded-xl border p-2 flex items-center gap-2"><span>🏦</span><div><p className="font-medium">{code}</p><p className="text-[11px] text-muted-foreground">{name} • {health}</p></div></div>
                ))}
              </div>
              <p className="text-[11px] text-muted-foreground mt-2">Methods ranked score 0.6*success+0.4*(1-latency/1000) sort desc • chosen true • Health snapshot telebirr 0.96 210ms CBE 0.89 260ms mock 1.0 45ms + fallback trail fallback_used false + Routing rules priority sort O(n log n) + circuit 5 fails open 60s map O(1)</p>
            </div>

            <div className="rounded-2xl border bg-card p-4">
              <h3 className="font-semibold">OpenAPI Swagger Embedded Modern • libs/openapi/openapi.yaml 21 paths</h3>
              <p className="text-xs text-muted-foreground mt-1">Title ApexPay API Full v1.1.0 • Servers https://api.apexpay.et/v1 + localhost:8080/v1 • Security bearerAuth sk_test_/sk_live_ prefix + secret hash at rest FIN never logged only last4 • FIN privacy note sha256+last4 + ETB + ONPS/10/2025 2FA &gt;5000 • Paths onboarding/kyc owners fayda verify init/confirm banks transactions/initialize verify payment_links refunds subscription_plans subscriptions beneficiaries payouts bulk employees payroll_runs calculate compliance/ask RAG citations mandatory no hallucination guard 0.65 swarm/run needs_confirmation &gt;100k methods ranked devices/register FCM unique • Contract test node scripts/audit/contract_test.js 21 paths present privacy NBE notes • k6 smoke 100 VUs p95&lt;300 ledger p99&lt;30</p>
              <button className="mt-2 rounded-xl border px-4 h-9 text-xs">Open Swagger UI • Outstanding</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
