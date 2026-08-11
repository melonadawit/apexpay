# Quickstart — 6 Lines Copy-Paste Like ApexPay "Simple as Copy and Paste"

ApexPay says: "By adding six lines of code — which ApexPay refers to as 'simple as a copy and paste job' — to their website or app, users can accept payments." [GxCamp](https://gxcamp.com/blog-details/the-first-ethiopias-online-payment-gateway-service--ApexPay-payment--)

ApexPay is even more powerful — same DX but with Fayda ID, smart routing, payouts, payroll, RAG, Swarm.

## 1. Install SDK

```bash
npm install apexpay-js
# or
yarn add apexpay-js
```

## 2. 6 Lines Copy-Paste — Outstanding

```javascript
import ApexPay from 'apexpay-js'

const apexpay = new ApexPay('sk_test_51Hq...') // test mode immediately after registration draft, live after KYC active per NBE

apexpay.initialize({
  tx_ref: 'txr_' + Date.now(), // unique (merchant_id, tx_ref) per DATABASE
  amount: '500.00', // ETB — numeric(20,8) + currency char(3) per DATABASE money rules never float
  currency: 'ETB',
  method: 'telebirr', // telebirr, cbe_birr, bank, card, qr — smart routing ranking score 0.6*success+0.4*(1-latency/1000) per routing service
  customer_email: 'cust@example.et',
  return_url: 'https://example.et/return',
  callback_url: 'https://example.et/callback',
  description: 'Tutoring • አስተማሪ'
}).then(res => {
  console.log(res.checkout_url) // https://checkout.apexpay.et/c/tok_abc QR outstanding Telegram/WhatsApp share
  window.location.href = res.checkout_url
}).catch(err => {
  console.error(err.code) // duplicate_tx_ref 409 etc per SAD §11 stable error codes
})
```

## 3. Verify Payment (Idempotent Second Success No-Op per MVP B6)

```javascript
// Poll verify every 2s O(n) with backoff per checkout-web real polling Day 3 spec
const verify = async (tx_ref) => {
  const res = await fetch(`https://api.apexpay.et/v1/transactions/verify/${tx_ref}`, {
    headers: { 'Authorization': 'Bearer sk_test_...' }
  })
  const data = await res.json()
  // data.status: succeeded, data.ledger_journal_balanced true per ValidateBalanced O(n) + quality check SQL
  return data
}
```

## 4. 2FA Mandatory >5000 ETB per ONPS/10/2025

```javascript
// If initialize returns requires_2fa true
if (res.requires_2fa) {
  // Show OTP input 6-digit
  await fetch(`https://api.apexpay.et/v1/transactions/${res.id}/2fa/verify`, {
    method: 'POST',
    headers: { 'Authorization': 'Bearer sk_test_...', 'Content-Type': 'application/json' },
    body: JSON.stringify({ payment_id: res.id, otp: '123456' }) // mock OTP 123456 per spec
  })
  // Then verify again
}
```

## 5. Webhooks HMAC SHA256 + SSRF Block + Retry Exponential Backoff

```javascript
// Node.js example
const crypto = require('crypto')
const express = require('express')
const app = express()

app.post('/webhook', express.json(), (req, res) => {
  const sig = req.headers['x-apexpay-signature'] // HMAC SHA256 hex per webhook service Sign()
  const payload = JSON.stringify(req.body)
  const expected = crypto.createHmac('sha256', 'whsec_...').update(payload).digest('hex')
  if (sig !== expected) return res.status(400).send('Invalid signature')
  
  // Idempotent receiver documented per SAD outbox pattern ADR-005
  // event_type: payment.succeeded, refund.succeeded, payout.succeeded, subscription.*
  console.log(req.body.event_type, req.body.data)
  res.sendStatus(200)
})
```

## 6. Fayda ID Front/Back OTP Consent per id.gov.et

```javascript
// Step 1: Init — FIN 12-digit + front/back images <2MB + selfie + OTP consent via id.gov.et VeriFayda 2.0 / OIDC eSignet
const faydaInit = await fetch('https://api.apexpay.et/v1/onboarding/fayda/verify/init', {
  method: 'POST',
  headers: { 'Authorization': 'Bearer sk_test_...', 'Content-Type': 'application/json' },
  body: JSON.stringify({
    merchant_id: 'mer_01H',
    owner_id: 'own_01H',
    kyc_profile_id: 'kyc_01H',
    fin: '123456789012', // FIN 12-digit plain never stored plain hashed sha256(salt+FIN)+last4 only
    fan: '', // FAN 16 alias alternative
    method: 'otp', // otp, face, fingerprint, offline_qr, oidc_esignet
    front_file_key: 'merchants/mer_01H/kyc/fayda_front_own_01H.jpg', // MinIO key <2MB per NIDP
    back_file_key: 'merchants/mer_01H/kyc/fayda_back_own_01H.jpg',
    selfie_file_key: 'merchants/mer_01H/kyc/selfie_own_01H.jpg'
  })
})
// Returns request_id fin_last4 ****1234 + otp_sent mock 123456

// Step 2: Confirm OTP 6-digit
await fetch('https://api.apexpay.et/v1/onboarding/fayda/verify/confirm', {
  method: 'POST',
  headers: { 'Authorization': 'Bearer sk_test_...', 'Content-Type': 'application/json' },
  body: JSON.stringify({ request_id: 'fayda_01H', otp: '123456' }) // mock OTP 123456 always success face 0.92
})
// Returns status verified face_match true face_score 0.92 >0.85 threshold demographics_match true fin_last4 1234
// FIN stored as hash + last4 only per DATABASE privacy rule, never plain in logs
```

Done! You now accept payments, verify Fayda, handle webhooks HMAC, and comply with NBE 2FA >5000 ETB per ONPS/10/2025.

Next: [Merchant Onboarding NBE + Fayda](/docs/onboarding) with 6-step wizard outstanding screenshots.
