# ApexPay — Payments-Only Integration Guide

> **Audience:** Merchants & developers who want to use **only the payment gateway**
> with their e-commerce / online store — and do **not** want payroll, HR, accounting,
> or any other module.
>
> **How it works:** ApexPay is modular. You get an API key, pick one of the four
> integration modes below, and connect it to your store the same way you would
> Chapa, Arif Pay, or Telebirr. You never touch `/payroll/*`, `/accounting/*`, etc.
> Those modules stay available on your account if you ever want them — but nothing
> forces you to use them.

---

## 0. One-time setup

1. **Sign up** and complete onboarding (KYC / Fayda where required).
2. **Get an API key** from the Developer portal:
   - Test key: `sk_test_...` (works immediately)
   - Live key: `sk_live_...` (available once the account is active)
   - A **public key** `pk_...` for frontend/browser use (safe to expose).
3. (Recommended) **Add a webhook endpoint** with a signing secret so your server
   receives payment events. Webhook endpoints are `https`-only and SSRF-protected.

That's the entire setup. No payroll, no HR, no ledger configuration required.

---

## Mode 1 — Hosted Checkout (no code / shareable)

Best for: invoices, WhatsApp/Telegram orders, or "send a payment link".

```bash
curl -X POST http://localhost:8080/v1/payment_links \
  -H "Authorization: Bearer sk_test_..." \
  -H "Content-Type: application/json" \
  -d '{"amount":"1500.00","currency":"ETB","description":"Order #1001"}'
```

Response:

```json
{
  "id": "pl_...",
  "amount": "1500.00",
  "currency": "ETB",
  "status": "active",
  "public_token": "tkn_...",
  "checkout_url": "https://checkout.apexpay.et/c/tkn_...",
  "expires_at": "...",
  "share": {
    "whatsapp": "https://wa.me/?text=Pay%20ETB%201500.00%20https://checkout.apexpay.et/c/tkn_...",
    "telegram": "https://t.me/share/url?url=..."
  }
}
```

Share `checkout_url` with the customer. ApexPay runs Telebirr / CBE Birr / bank /
card / EthSwitch QR and the 2FA step. Optionally set `callback_url` to receive a
webhook when the customer pays.

---

## Mode 2 — Embedded JS SDK (`checkout.js`)

Best for: a single-page app or mobile-web store. Public key only — no secret in the
browser.

```html
<script src="https://checkout.apexpay.et/sdk.js"></script>
<script>
  const apexpay = new ApexPay('pk_test_...');
  apexpay.checkout({
    amount: '500',
    currency: 'ETB',
    tx_ref: 'txr_' + Date.now(),
    method: 'telebirr',            // or cbe_birr, bank, card_acquirer, ethswitch
    customer_email: 'cust@example.et',
    return_url: 'https://store.example/return',
    callback_url: 'https://store.example/api/payments/webhook'
  });
</script>
```

---

## Mode 3 — Direct REST API (for e-commerce carts / custom backend)

Best for: WooCommerce-style carts, custom storefronts, headless commerce, or a
server-side integration where you control the checkout UI. This is the classic
gateway flow and matches the `/v1/transactions` API exactly.

### 3a. Initialize the payment

```bash
curl -X POST http://localhost:8080/v1/transactions/initialize \
  -H "Authorization: Bearer sk_live_..." \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: order-1001-retry" \
  -d '{
    "tx_ref": "order-1001",
    "amount": "2500.00",
    "currency": "ETB",
    "method": "telebirr",
    "customer_email": "buyer@example.com",
    "return_url": "https://store.example/pay/return",
    "callback_url": "https://store.example/api/payments/callback"
  }'
```

Response (`201`):

```json
{
  "id": "pay_...",
  "tx_ref": "order-1001",
  "amount": "2500.00",
  "currency": "ETB",
  "status": "created",
  "checkout_url": "https://checkout.apexpay.et/c/...",
  "requires_2fa": true,
  "fee_amount": "72.50",
  "net_amount": "2427.50",
  "connector_id": "telebirr"
}
```

- `tx_ref` is **your** unique reference — use the order id. Reuse the same
  `Idempotency-Key` for safe retries.
- Send the customer to `checkout_url`, or embed it.

### 3b. Verify on return (server-side)

```bash
curl http://localhost:8080/v1/transactions/verify/{tx_ref} \
  -H "Authorization: Bearer sk_live_..."
```

Returns the payment with `status` (`created` / `processing` / `succeeded` /
`failed`). If `requires_2fa` is true and not yet verified, status stays `created`.

### 3c. Receive the webhook (authoritative)

When the payment resolves, ApexPay sends an event to your `callback_url`:

```
POST https://store.example/api/payments/callback
Headers:
  X-ApexPay-Signature: <HMAC-SHA256(signing_secret, raw_body)>
  Content-Type: application/json

{ "event_type": "payment.succeeded", "payment_id": "pay_...", "tx_ref": "order-1001", "status": "succeeded" }
```

**Best practice for e-commerce:**
1. **Verify the HMAC signature** before trusting the event.
2. Use the webhook (not the client redirect) to mark the order paid.
3. Make fulfillment **idempotent** — events can be retried.

---

## Mode 4 — Headless / server-to-server

For a mobile app or backend that builds its own UI, use the API key directly
server-to-server: `initialize` → redirect/embed → `verify` + webhook. Same shape as
Mode 3; you don't use the hosted page.

---

## Which mode should I use?

| Need | Mode |
|---|---|
| Share a payment link (WhatsApp/invoice) | **1 — Hosted checkout / payment link** |
| Embed a pay button in a SPA/mobile-web | **2 — JS SDK** |
| Full cart/checkout integration, custom backend | **3 — REST API** |
| Headless or server-driven checkout | **4 — REST API, headless** |

All modes support Telebirr, CBE Birr, bank/IPS, card, and EthSwitch QR, with 2FA
above the NBE threshold and ledger settlement in the same transaction.

---

## What you are NOT required to use

ApexPay's payroll, HR, general ledger, procurement, inventory, tax, and budgeting
modules are **optional**. A payments-only merchant simply never calls those
endpoints. You get a clean gateway experience — like Chapa / Arif Pay / Telebirr —
and can enable the deeper modules later on the same account if you grow into them.

---

## Status of the plug-in story (honest note)

**Implemented and working today:** API keys (test/live/public), `initialize` /
`verify` / list / detail, hosted checkout, payment links, webhooks with HMAC signing,
and the developer-portal screens. The exact request/response shapes above are from
the running API.

**Not yet built as one-click drop-ins:**
- Ready-made **WooCommerce / Shopify / Magento plugins** (integrate via the REST API
  today; no install-and-go plugin exists in the repo yet).
- A **published npm/PHP SDK package** (the `checkout.js` snippet is documented and
  works, but there's no published package in the repo).
- **Live** Telebirr / CBE Birr / bank connectors are wired via the connector registry
  but run in test/mock mode — a real go-live needs the live connector credentials.

These are the natural next artifacts if you want to turn this guide into a
merchant-facing "Connect your store" onboarding flow.
