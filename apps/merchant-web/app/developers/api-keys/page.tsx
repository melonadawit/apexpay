"use client"
import * as React from "react"

export default function ApiKeysPage() {
  const [showSecret, setShowSecret] = React.useState(false)
  const [keys, setKeys] = React.useState([
    { id:"key_01H", name:"test key", prefix:"sk_test_51Hq", type:"secret", env:"test", status:"active", last_used:"2 min ago", scopes:["payments:read","payments:write"] },
    { id:"key_02H", name:"live key", prefix:"sk_live_51Hq", type:"secret", env:"live", status:"pending_activation", last_used:"never", scopes:["*"] },
  ])

  return (
    <div className="min-h-screen bg-neutral-50 p-6">
      <div className="max-w-5xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold">API Keys • Test/Live Separate Scopes Reveal Once — Outstanding</h1>

        <div className="rounded-2xl border bg-white p-6 space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="font-semibold">Secret Keys • sk_test_ / sk_live_ — hash at rest per DATABASE</h3>
            <button className="rounded-xl bg-primary text-white px-4 h-9 text-xs">Create New Key • test/live separate live only after KYC active</button>
          </div>

          <div className="space-y-3">
            {keys.map(k=>(
              <div key={k.id} className="rounded-xl border p-4 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="h-10 w-10 rounded-xl bg-primary/10 flex items-center justify-center font-mono text-xs">{k.prefix.slice(0,3)}</div>
                  <div>
                    <p className="font-medium text-sm">{k.name} • {k.prefix}•••••••• • {k.env} • {k.type}</p>
                    <p className="text-xs text-muted-foreground">Scopes {k.scopes.join(", ")} • Last used {k.last_used} • Status {k.status} • Prefix unique index api_keys_prefix_uidx O(1) lookup • secret_hash index where not null • scopes jsonb • last_used_at async best effort non-blocking Go routine</p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <span className={`text-xs px-2 py-0.5 rounded-full ${k.status==="active" ? "bg-green-100 text-green-700" : "bg-amber-100"}`}>{k.status}</span>
                  <button onClick={()=> setShowSecret(!showSecret)} className="rounded-lg border px-3 py-1 text-xs">{showSecret ? "Hide" : "Reveal"} • sk_live shown once</button>
                </div>
              </div>
            ))}
          </div>

          {showSecret && (
            <div className="rounded-xl bg-amber-50 border border-amber-200 p-3 font-mono text-sm">
              sk_test_51Hq...abc123xyz — <span className="text-red-600">shown once per security best practice — copy now!</span> — Stored as hash sha256(salt+secret) or bcrypt/argon2 at rest, prefix visible later for audit who used which key.
            </div>
          )}

          <div className="rounded-xl bg-blue-50 border border-blue-200 p-3 text-xs">
            <p className="font-semibold">Security Best Practice:</p>
            <ul className="list-disc list-inside mt-1 space-y-0.5">
              <li>Secret shown once — hash at rest per DATABASE + prefix unique index O(1) lookup • secret_hash index where not null • scopes jsonb • last_used_at async best effort non-blocking Go routine</li>
              <li>Test/Live separation: test keys immediately after registration draft, live keys only after KYC active + dual approval if high risk + pilot 30-60 days analogy NBE</li>
              <li>Prefix visible later for audit who used which key — audit_logs actor_type api_key actor_id key prefix</li>
              <li>Scopes: payments:read, payments:write, refunds:write, payouts:write, etc — RBAC map O(1) role check owner/admin/developer/finance</li>
            </ul>
          </div>
        </div>

        <div className="rounded-2xl border bg-white p-6">
          <h3 className="font-semibold">Public Keys • pk_test_ / pk_live_ + Embedded SDK checkout.js</h3>
          <p className="font-mono text-xs mt-2">pk_test_51Hq... + checkout.js embedded SDK `https://checkout.apexpay.et/sdk.js` — no secret, safe for frontend, tokenization only, no PAN storage</p>
          <div className="mt-3 rounded-xl bg-neutral-900 text-white p-4 font-mono text-[11px] overflow-auto">
            {`<script src="https://checkout.apexpay.et/sdk.js"></script>
<script>
  const apexpay = new ApexPay('pk_test_51Hq...');
  apexpay.checkout({
    amount: '500',
    currency: 'ETB',
    tx_ref: 'txr_01H_' + Date.now(),
    method: 'telebirr',
    customer_email: 'cust@example.et',
    return_url: 'https://example.et/return',
    callback_url: 'https://example.et/callback'
  });
</script>`}
          </div>
        </div>
      </div>
    </div>
  )
}
