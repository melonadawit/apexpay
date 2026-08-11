"use client"
import * as React from "react"
import Link from "next/link"
import { Check, Copy, Wrench, Webhook, Zap, Link2, TerminalSquare } from "lucide-react"
import { useLanguage } from "@/components/providers/language-provider"

function CopyBlock({ code, label }: { code: string; label: string }) {
  const [copied, setCopied] = React.useState(false)
  const copy = async () => {
    try {
      await navigator.clipboard.writeText(code)
      setCopied(true)
      setTimeout(() => setCopied(false), 1500)
    } catch {
      /* clipboard unavailable in some contexts */
    }
  }
  return (
    <div className="rounded-xl bg-neutral-950 text-neutral-100 overflow-hidden">
      <div className="flex items-center justify-between px-4 py-2 border-b border-neutral-800">
        <span className="text-[11px] font-mono text-neutral-400">{label}</span>
        <button onClick={copy} className="flex items-center gap-1 text-[11px] text-neutral-300 hover:text-white">
          {copied ? <Check className="h-3 w-3" /> : <Copy className="h-3 w-3" />}
          {copied ? "Copied" : "Copy"}
        </button>
      </div>
      <pre className="px-4 py-3 text-[12px] leading-relaxed overflow-auto">{code}</pre>
    </div>
  )
}

const SNIPPETS: { mode: string; icon: React.ReactNode; title: string; when: string; code: string }[] = [
  {
    mode: "1",
    icon: <Link2 className="h-4 w-4" />,
    title: "Hosted Checkout / Payment Link",
    when: "Share a payment link on WhatsApp/Telegram or an invoice. No code.",
    code: `# Create a payment link (one-off or recurring)
curl -X POST https://api.apexpay.et/v1/payment_links \\
  -H "Authorization: Bearer sk_test_..." \\
  -H "Content-Type: application/json" \\
  -d '{"amount":"1500.00","currency":"ETB","description":"Order #1001"}'

# → checkout_url: https://checkout.apexpay.et/c/tkn_...
# Share it, or use the "share.whatsapp" / "share.telegram" URLs.`,
  },
  {
    mode: "2",
    icon: <Zap className="h-4 w-4" />,
    title: "Embedded JS SDK (checkout.js)",
    when: "Embed a pay button in a single-page or mobile-web store. Public key only.",
    code: `<script src="https://checkout.apexpay.et/sdk.js"></script>
<script>
  const apexpay = new ApexPay('pk_test_...');
  apexpay.checkout({
    amount: '500', currency: 'ETB',
    tx_ref: 'txr_' + Date.now(),
    method: 'telebirr',            // telebirr | cbe_birr | bank | card_acquirer | ethswitch
    customer_email: 'cust@example.et',
    return_url: 'https://store.example/return',
    callback_url: 'https://store.example/api/payments/webhook'
  });
</script>`,
  },
  {
    mode: "3",
    icon: <TerminalSquare className="h-4 w-4" />,
    title: "Direct REST API",
    when: "Full cart / checkout integration with your own backend.",
    code: `// 1. Initialize a payment
curl -X POST https://api.apexpay.et/v1/transactions/initialize \\
  -H "Authorization: Bearer sk_live_..." \\
  -H "Content-Type: application/json" \\
  -H "Idempotency-Key: order-1001" \\
  -d '{"tx_ref":"order-1001","amount":"2500.00","currency":"ETB",
       "method":"telebirr","customer_email":"buyer@example.com",
       "return_url":"https://store.example/pay/return",
       "callback_url":"https://store.example/api/payments/callback"}'
# → 201 { "id":"pay_...", "checkout_url":"...", "status":"created" }

# 2. Redirect the customer to checkout_url

# 3. Verify on return (server-side)
curl https://api.apexpay.et/v1/transactions/verify/order-1001 \\
  -H "Authorization: Bearer sk_live_..."
# → { "status":"succeeded" }`,
  },
  {
    mode: "4",
    icon: <Webhook className="h-4 w-4" />,
    title: "Webhook",
    when: "Get notified authoritatively when a payment resolves (fulfil orders).",
    code: `// ApexPay signs each delivery with an HMAC.
POST https://store.example/api/payments/webhook
Headers:
  X-ApexPay-Signature: <HMAC-SHA256(signingSecret, rawBody)>

{ "event_type": "payment.succeeded",
  "payment_id": "pay_...",
  "tx_ref": "order-1001",
  "status": "succeeded" }

// 1. Verify the signature before trusting the event.
// 2. Mark the order paid (idempotently — events can be retried).`,
  },
]

export default function ConnectStorePage() {
  const { t } = useLanguage()
  const nodeSnippet = `npm install @apexpay/node

import { ApexPay } from "@apexpay/node";
const apexpay = new ApexPay({ apiKey: "sk_test_..." });

const payment = await apexpay.initialize({
  tx_ref: "order-1001",
  amount: "2500.00",
  currency: "ETB",
  method: "telebirr",
  callback_url: "https://store.example/api/payments/webhook",
});
// payment.checkout_url -> redirect the customer here

// On return, verify server-side:
const verified = await apexpay.verify("order-1001");
if (verified.status === "succeeded") { /* mark order paid */ }`

  return (
    <div className="min-h-screen bg-muted p-6">
      <div className="max-w-5xl mx-auto space-y-6">
        <Link href="/developers" className="text-sm text-primary">← Developers</Link>
        <div>
          <h1 className="text-2xl font-bold">{t("Connect your store", "መደብርዎን ያገናኙ")}</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Use ApexPay&apos;s payment gateway only — no payroll or HR required. Pick a mode, copy, and go.
          </p>
        </div>

        <div className="rounded-2xl border bg-card p-6">
          <div className="flex items-center gap-2 mb-3">
            <Wrench className="h-4 w-4 text-primary" />
            <h3 className="font-semibold">One-time setup</h3>
          </div>
          <ol className="list-decimal list-inside text-sm space-y-1 text-muted-foreground">
            <li>Get a test key <code className="font-mono">sk_test_...</code> (live key after KYC) from <Link href="/developers/api-keys" className="text-primary">API Keys</Link>.</li>
            <li>Optionally add a <Link href="/developers/webhooks" className="text-primary">webhook endpoint</Link> with a signing secret.</li>
            <li>Pick a mode below and connect to your store.</li>
          </ol>
        </div>

        {SNIPPETS.map((s) => (
          <div key={s.mode} className="rounded-2xl border bg-card overflow-hidden">
            <div className="p-5 pb-3">
              <div className="flex items-center gap-2">
                <span className="inline-flex items-center justify-center h-6 w-6 rounded-md bg-primary/10 text-primary">{s.icon}</span>
                <h3 className="font-semibold">Mode {s.mode} — {s.title}</h3>
              </div>
              <p className="text-xs text-muted-foreground mt-1">{s.when}</p>
            </div>
            <div className="px-5 pb-5">
              <CopyBlock code={s.code} label={`Mode ${s.mode}`} />
            </div>
          </div>
        ))}

        <div className="rounded-2xl border bg-card overflow-hidden">
          <div className="p-5 pb-3">
            <div className="flex items-center gap-2">
              <span className="inline-flex items-center justify-center h-6 w-6 rounded-md bg-primary/10 text-primary">
                <TerminalSquare className="h-4 w-4" />
              </span>
              <h3 className="font-semibold">Node.js SDK (official)</h3>
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              The <code className="font-mono">@apexpay/node</code> package wraps initialize / verify / payment links and webhook
              signature verification. (Scaffold in <code className="font-mono">sdk/node</code>.)
            </p>
          </div>
          <div className="px-5 pb-5">
            <CopyBlock code={nodeSnippet} label="Node SDK" />
          </div>
        </div>

        <div className="rounded-2xl bg-blue-500/10 border border-blue-500/20 p-4 text-xs text-muted-foreground">
          <p className="font-semibold text-foreground">Payments-only — nothing else required.</p>
          <p className="mt-1">
            You never call payroll, HR, or accounting endpoints. Those modules stay available on your account if you
            grow into them, but you can run a pure payment gateway today — just like Chapa, Arif Pay, or Telebirr.
          </p>
        </div>
      </div>
    </div>
  )
}
