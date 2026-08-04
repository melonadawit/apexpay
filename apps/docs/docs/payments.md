# Payments — Initialize/Verify + 2FA >5000 ETB + Smart Routing + Ledger M1 Balanced

Per payment/domain.go + service.go + repository.go + handler.go + ledger M1

## Initialize

- Request: tx_ref unique (merchant_id, tx_ref) per DATABASE unique index O(1) lookup, amount numeric(20,8) + currency char(3) per money rules never float, method telebirr/cbe_birr/bank/card/qr, description, customer_email, return_url, callback_url, Idempotency-Key header X-Idempotency-Key per DATABASE idempotency_keys PK (merchant_id,key) + request_hash
- Idempotency check O(1) lookup via idempotency_keys table resource_type payment resource_id
- Duplicate tx_ref 409 duplicate_tx_ref stable code per SAD §11
- Routing evaluation O(n log n) priority sort + health 5m success_rate latency + circuit breaker 5 fails open 60s map O(1) + score 0.6*success+0.4*(1-latency/1000) + fallback trail fallback_used false + health snapshot
- Fee calc: fee = amount * mdrRate 0.029 Round2 ETB scale 2 per ETB =2 app enforces scale, net = amount - fee
- 2FA check per NBE ONPS/10/2025 >5000 ETB requires_2fa true
- Connector Initialize 50ms latency inject per mock + telebirr sandbox HMAC X-APP-Key + callback IP allowlist + latency 150-300ms + CBE sandbox 30% failure inject + bank IPS ISO20022 + EthSwitch QR + card token no PAN + config encrypted AES-GCM + health sampler 30s + Redis cache
- Create payment pending checkout_url https://checkout.apexpay.et/c/tok_abc + fee_amount net_amount + routing_rule_id + return_url callback_url + created_at + updated_at now()
- CreatePaymentTx atomic: INSERT payments + outbox payment.created per outbox pattern ADR-005 no dual-write loss + SaveIdempotency

## Verify

- GetByTxRef merchant_id tx_ref O(1) unique index
- If status succeeded return idempotent no-op per MVP B6 single journal posting_key per DATABASE unique (book_id, posting_key)
- If requires_2FA true && !two_fa_verified return pending
- Call connector Verify connector_ref tx_ref → status succeeded amount failure_code
- If verify status != succeeded return pending worker will retry
- Ledger M1 posting atomically with status update per DATABASE transaction boundary NEVER commit payment success without ledger post in same Tx: journal ID NewLedgerJournal BookID merchant_operating_default PostingKey payment_success:{pay_id} Memo payment success ReferenceType payment ReferenceID pay_id TransferGroup pay_{id} + entries Dr asset:clearing:{connector_id} amount Cr liability:merchant_payable net Cr liability:platform_fee_due fee + filtered zero entries if feeReversal zero + ValidateBalanced debit==credit + balances upsert ledger_balances amount updated_at PK (book_id,account_id) + outbox payment.succeeded SELECT merchant_id FROM payments WHERE id=$2
- UpdateStatusTx tx begin UPDATE payments status succeeded_at + INSERT ledger_journals ON CONFLICT DO NOTHING + INSERT ledger_entries + INSERT ledger_balances ON CONFLICT DO UPDATE SET amount = ledger_balances.amount + EXCLUDED.amount + outbox payment.succeeded + tx commit

## 2FA Verify

- POST /v1/transactions/{id}/2fa/verify OTP 6-digit mock 123456 per spec per NBE ONPS/10/2025 2FA mandatory >5000 ETB + transaction exceeding 5,000 Birr must now use two-factor authentication including PIN, OTP, biometric per Addis Insight May 2025
- Real: SMS/Email 2FA service + OTP verify + update payment two_fa_verified true
- Mock: OTP 123456 passes + can_verify_now true + polling verify every 2s O(n) backoff per checkout-web real polling Day 3

## Methods Ranked

- GET /v1/methods?amount=1000&currency=ETB ranked array with latency_ms success_rate_5m score chosen boolean + health snapshot per routing service RankedMethods all connectors success*0.6+(1-latency/1000)*0.4 sort Score desc O(n log n)

## Checkout Outstanding

- Checkout-web mobile-first max 420px centered merchant logo name small amount large 32px bold method selector radio cards with icons Telebirr/CBE/bank/card/QR interoperable selected border primary 2px + subtle bg primary/5 + check icon + trust badges NBE Licensed PSO Gateway + secure lock + Using best route badge telebirr 2% faster today tooltip fallback trail fallback_used false health snapshot + 2FA OTP input if >5000 + processing lottie Ethiopia pattern polling verify every 2s O(n) backoff + success full-screen check animation Lottie scale + confetti canvas-confetti 3s + haptics + receipt actions download PDF jsPDF invoice outstanding template with QR email back to merchant link

Next: [Refunds FULL M2](/docs/refunds) + [Payouts Bulk 1000 + Escrow](/docs/payouts)
