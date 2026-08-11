# @apexpay/node

Official ApexPay **Node.js** SDK — a payments-only client for merchants who want just
the payment gateway for their e-commerce / online store (like Chapa, Arif Pay, or
Telebirr). No payroll, no HR, no ledger required.

## Install

```bash
npm install @apexpay/node
```

## Quickstart

```ts
import { ApexPay } from "@apexpay/node";

const apexpay = new ApexPay({ apiKey: "sk_test_..." });

// 1. Initialize a payment and send the customer to checkout_url
const payment = await apexpay.initialize({
  tx_ref: "order-1001",
  amount: "2500.00",
  currency: "ETB",
  method: "telebirr", // telebirr | cbe_birr | bank | card_acquirer | ethswitch
  customer_email: "buyer@example.com",
  return_url: "https://store.example/pay/return",
  callback_url: "https://store.example/api/payments/callback",
});
// payment.checkout_url -> redirect the customer here

// 2. (On return) verify server-side
const verified = await apexpay.verify("order-1001");
if (verified.status === "succeeded") { /* mark order paid */ }

// 3. Or create a shareable payment link (hosted checkout, no code)
const link = await apexpay.createPaymentLink({ amount: "1500.00", description: "Order #1001" });
// link.checkout_url -> share via WhatsApp/Telegram
```

## Webhooks

```ts
import { verifyWebhookSignature } from "@apexpay/node";

// In your HTTP handler:
const rawBody = req.rawBody; // exact bytes ApexPay sent
const signature = req.headers["x-apexpay-signature"];
if (!verifyWebhookSignature("your-signing-secret", rawBody, signature)) {
  // reject — signature invalid
}
// trusted event — mark the order paid (idempotently)
```

## API

- `initialize(req)` → `Payment` with `checkout_url`
- `verify(txRef)` → `Payment` with `status`
- `getPayment(id)` → `Payment`
- `createPaymentLink({ amount, currency?, description? })` → `PaymentLink`
- `verifyWebhookSignature(secret, rawBody, signature)` → `boolean`

See [`docs/PAYMENTS_ONLY_GUIDE.md`](../../docs/PAYMENTS_ONLY_GUIDE.md) for the full
integration modes and setup.
