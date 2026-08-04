# API Reference — 21 Paths Full v1.1.0 + 2FA >5000 ETB per ONPS/10/2025

OpenAPI spec: `libs/openapi/openapi.yaml` — title ApexPay API Full v1.1.0, servers https://api.apexpay.et/v1 + localhost:8080/v1, security bearerAuth sk_test_/sk_live_ prefix secret hash at rest FIN never logged only last4.

All paths per `scripts/audit/contract_test.js` 21 paths present privacy NBE notes.

## Onboarding

- `POST /v1/onboarding/kyc` — Create KYC profile NBE checklist legal_name business_type registration_number TIN 10-digit industry_category region city office_address_full contact_person_name role email phone expected_monthly_tpv. Returns id kyc_01H status draft.
- `POST /v1/onboarding/owners` — Add beneficial owner UBO >10% + authorized signatory per NBE min 5 shareholders. full_name full_name_am role owner/shareholder/director/authorized_rep/ubo ownership_percentage 0-100 id_type fayda/passport phone is_authorized_signatory is_pep.
- `POST /v1/onboarding/bank-accounts` — Bank account settlement per PayAtlas Bank Statement: account_name must == legal fuzzy Levenshtein <3, bank_code CBE/AWASH/DASHEN etc 14 ET banks, account_number masked ****1234 hash sha256, is_settlement_default.
- `POST /v1/onboarding/documents/presign` — MinIO presigned PUT 15m TTL direct upload no server buffering O(n) streaming sha256 merchants/{merchant_id}/kyc/{doc_type}_{id}.pdf.
- `POST /v1/onboarding/documents` — file_key file_hash mime whitelist pdf/jpg/png size <5MB Fayda <2MB per NIDP, status uploaded.
- `POST /v1/onboarding/submit` — Submit KYC completeness O(n) docMap + hasAuthSignatory + faydaVerifiedCount>=1 + settlement bank + risk weighted sum TPV high +20 PEP*30.

## Fayda ID Verification Front/Back <2MB + Selfie + OTP Consent id.gov.et

- `POST /v1/onboarding/fayda/verify/init` — FIN 12-digit plain never stored plain hashed sha256(salt+FIN)+last4 only, FAN 16 alias alternative, method otp/face/fingerprint/offline_qr/oidc_esignet, front_file_key MinIO key <2MB per NIDP, back_file_key, selfie_file_key, consent IP logged. Returns request_id fin_last4 ****1234 otp_sent mock 123456 fayda_transaction_id mock_tx.
- `POST /v1/onboarding/fayda/verify/confirm` — OTP 6-digit mock 123456 always success face 0.92 per MockVerifier, 000000 fail. Returns status verified face_match true face_score 0.92 >0.85 threshold demographics_match true fin_last4 1234. FIN stored as hash + last4 only per DATABASE privacy rule.
- `POST /v1/onboarding/fayda/verify/qr` — Offline QR FaydaEncode scan without network per id.gov.et "You can get your Fayda verified without the need for network connectivity, simply by either using the QR code (we also call the QR FaydaEncode) on your Fayda Credential". QRData FIN_LAST4|NAME|DOB|SIG signature valid.

## Banks + Methods Ranked

- `GET /v1/banks` — Ethiopian banks list 14 ET banks CBE/Awash/Dashen/Abyssinia/Berhan/Wegagen/NIB/United/Coop/Oromia/Bunna/Lion/Zemen/CBO per migration 0012 seed. Code name name_am.
- `GET /v1/methods?amount=1000&currency=ETB` — Ranked methods smart routing success_rate/latency/cost score 0.6*success+0.4*(1-latency/1000) sort desc O(n log n) priority sort + circuit breaker 5 fails open 60s map O(1) + health 5m success_rate latency. Returns connector telebirr success_rate_5m 0.95 latency_ms 210 score 0.88 chosen true reason primary healthy fallback_used false health_snapshot.

## Payments + 2FA

- `POST /v1/transactions/initialize` — Idempotency-Key header + amount numeric decimal precise + currency ETB + method telebirr/cbe_birr/bank/card/qr + description + customer_email + return_url + callback_url. Fee = amount*0.029 Round2 net = amount-fee. Requires_2FA true if ETB>5000 per ONPS/10/2025 + transaction exceeding 5,000 Birr must now use two-factor authentication including PIN, OTP, biometric. Returns id tx_ref checkout_url https://checkout.apexpay.et/c/tok_abc qr_data_url data:image/png;base64,mock_qr + checkout_session_id expires_at 24h + share telegram/whatsapp.
- `GET /v1/transactions/verify/{tx_ref}` — Verify pending → succeeded + ledger M1 balanced true per ValidateBalanced O(n) + quality check SQL. Idempotent second success no-op single journal posting_key per DATABASE unique (book_id, posting_key). Poll every 2s in checkout-web real polling Day 3.
- `POST /v1/transactions/{id}/2fa/verify` — OTP 6-digit mock 123456 for 2FA >5000 ETB.

## Payment Links + QR

- `POST /v1/payment_links` — Amount currency description expires_at optional. Returns id amount currency description status active public_token checkout_url https://checkout.apexpay.et/c/tok_abc qr_data_url.
- `GET /v1/payment_links` — List by merchant order created_at desc limit 100.
- `GET /v1/payment_links/public/{token}` — Public no auth outstanding for checkout-web mobile 420px centered amount large 32px bold method selector radio cards icons Telebirr/CBE + best route badge + 2FA OTP + processing lottie + success confetti.

## Refunds FULL M2

- `POST /v1/refunds` — payment_id refund_ref unique (merchant_id, refund_ref) Idempotency-Key + amount partial allowed + currency + reason + fee_policy non_refundable/pro_rata/full. Fee reversal calc decimal precise bankers rounding pro_rata = totalFee * (refund/pay) Round2 ETB scale. Ledger M2 Dr payable R-FR + Dr fee_due FR Cr clearing R filter zero. Returns id refund_ref amount fee_reversal status payment_id ledger_model M2.

## Subscriptions + Dunning

- `POST /v1/customers` — email phone name Fayda hash optional per customers table.
- `POST /v1/subscription_plans` — name description amount currency interval_type day/week/month/year interval_count trial_days.
- `POST /v1/subscriptions` — customer_id plan_id trial 7d status trialing period trial_end + invoice draft → open due current_period_end attempt 0 next +1d. Dunning worker cron hourly SELECT ... FOR UPDATE SKIP LOCKED attempt 0→+24h 1→+72h 2→+120h.
- `GET /v1/subscriptions` — List merchant status filter.
- `POST /v1/subscriptions/{id}/cancel` — canceled.

## Payouts + Bulk + Escrow

- `POST /v1/beneficiaries` — name account_no bank_code bank_name type individual/business verification fuzzy Levenshtein <3 hash+masked.
- `POST /v1/payouts` — beneficiary_id payout_ref amount currency method bank/mobile_money/payout_link — maker-checker >50k pending_approval else queued M3 Dr payable Cr clearing bank atomic per batch book.
- `POST /v1/payouts/bulk` — batch_ref currency items beneficiary_id amount payout_ref bulk 1-1000 total sum balance check batch pending_approval all bulk require approval journal Dr payable total Cr clearing total per batch book.
- `POST /v1/payout_batches/{id}/approve` — dual approval finance+admin if >50k.

## Payroll ET Tax + Pension

- `POST /v1/employees` — employee_code name bank masked hash base_salary employment_date cost_center Fayda hash optional status active.
- `GET /v1/employees` — List active.
- `POST /v1/payroll_runs` — run_ref period_month/year type regular/off_cycle/bonus/adjustment status draft.
- `POST /v1/payroll_runs/{id}/calculate` — binary search O(log n) tax brackets 7 ET 2024 0-600 0% etc effective_from versioned per migration 0008 seed + pension employee 7% employer 11% + OT rates map O(1) weekday 1.25 weekend 1.5 holiday 2.0 night 1.3 + totals gross/net/tax/pension update pending_approval ledger M4 draft Dr salary totalGross Cr payroll_payable net Cr tax payable tax Cr pension payable totalPension ValidateBalanced per run book.
- `POST /v1/payroll_runs/{id}/approve` — dual >100k net.
- `POST /v1/payroll_runs/{id}/disburse` — per run book + payout batch second journal Dr payroll_payable Cr bank totalNet.
- `GET /v1/payroll_runs/{id}/items` — List items gross taxable income_tax pension_employee employer net_pay.

## RAG Compliance Citations Mandatory Threshold 0.65 No Hallucination + AM/EN

- `POST /v1/compliance/ask` — query When is 2FA required? lang en/am top_k 5. Pipeline: query → embed multilingual-e5-large 1024 dim normalized L2 + batch 32 optimal → pgvector ivfflat lists=100 cosine O(log n) <=> → threshold 0.65 guard if top score < threshold → no answer "Not in compliance corpus" prevents hallucination → build prompt context [1]..[n] + question + lang → LLM mock returns answer with citations mandatory per eval harness. Returns answer Transactions above 5000 ETB require 2FA per ONPS/10/2025 [1] + citations document_id chunk_id content score 0.92 page + no_answer false query.

## Swarm Multi-Agent Planner/Critic/Executor Confirmation >100k

- `POST /v1/swarm/run` — goal Create link 100 ETB for coffee and run payroll July with bonus → planner RulesPlanner keyword + critic threshold >100k + state machine planning→executing→needs_confirmation→completed → plan array step tool create_payment_link description args amount currency description status pending → confirmation_required true confirmation_data total_amount steps → needs_confirmation outstanding modal.
- `POST /v1/swarm/{id}/confirm` — confirmed bool true → executing all pending steps tool registry O(1) payment_link, payout, payroll, tpv, compliance + tool_calls audit latency + final_output Created link https://... + payroll run.
- `GET /v1/swarm/{id}` — Get session plan status final_output.

## Webhooks + Security

- `POST /v1/webhooks/endpoints` — URL Secret Events payment.succeeded refund.succeeded payout.succeeded subscription.* — SSRF block private ranges Allowlist 10.0.0.0/8 172... 192.168 127.0.0.1 + secret_hash hash + secret_prefix whsec_ + status active events ["*"].
- `GET /v1/webhooks/endpoints` — List endpoints id url status.
- `GET /v1/webhooks/deliveries` — List deliveries id event_type status attempt_count last_status_code.
- `POST /v1/webhooks/deliveries/{id}/resend` — status pending next_attempt_at now → queued.
- Delivery: HMAC SHA256 X-ApexPay-Signature sig + X-ApexPay-Event + X-Request-Id + Content-Type JSON + retry exponential backoff 1s 2s 4s 8s 16s 32s max 1h + circuit breaker.

## Methods + Devices + Admin

- `GET /v1/methods` — Ranked methods smart routing score.
- `POST /v1/devices/register` — FCM token unique push_devices table platform android/ios/web + device_info model os + last_active_at now.
- `GET /v1/admin/onboarding/queue` — Queue submitted/in_review/fayda_pending/compliance_check order created_at limit 50 + risk_score fayda_verified.
- `POST /v1/admin/onboarding/{id}/review` — action approve/reject/request_info + comments + internal_notes + merchant active + operating book 6 accounts seeded + outbox merchant.activated.
- `GET /v1/admin/connectors/health` — AVG latency 5m + success_rate 5m GROUP BY connector_id + circuit closed.
- `GET /v1/admin/recon/breaks` — status open.
- `GET /v1/admin/evidence?tx_ref` — tx_ref ledger_journals fayda_verified docs_hashes onboarding_reviews_chain audit_logs webhook_deliveries.
- `GET /v1/admin/merchants/{id}/exam` — KYC profiles owners Fayda badge face 0.92 OTP docs viewer hashes compliance checks onboarding reviews timeline banks ledger books.

## Error Codes Stable per SAD §11

- duplicate_tx_ref 409, duplicate_refund_ref 409, invalid_fayda_fin 400, fayda_otp_failed 400, fayda_not_verified 400, onboarding_not_ready 400, document_required 400, bank_verification_failed 400, insufficient_balance 400, refund_exceeded 400, not_found 404, unauthorized 401, forbidden 403, validation_error 400, conflict 409, rate_limited 429, connector_unavailable 502.

## Authentication

- Bearer sk_test_ / sk_live_ per DATABASE api_keys_prefix_uidx O(1) secret_hash hash at rest FIN masked ****1234 only last4 never plain FIN in logs + scopes + last_used_at async best effort non-blocking Go routine + prefix visible later for audit who used which key + test/live separate live only after KYC active + pilot 30-60 days analogy NBE + maker-checker dual approval risk>=70 or TPV>1M.

Full spec: `libs/openapi/openapi.yaml` + Postman collection `libs/postman/ApexPay.postman_collection.json`.
