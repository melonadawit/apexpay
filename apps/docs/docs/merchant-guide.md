# Merchant Guide — Outstanding Modern Dashboard + Tracking Timeline

For Meron — SMB owner tutoring — create link, share on Telegram, know when paid + NBE exam-ready.

## Dashboard • ዳሽቦርድ

After onboarding approved active + operating book created:

- TPV Today glass gradient emerald #0B6E4F → light #10A37A sparkline 7 bars 40/60/30/80/50/90/60 success rate 96.2% fallback 3 times per routing fallback_used + active links 12 QR thumbnails + recent payments telebirr/cbe/bank 2FA >5000 badge + quick actions Create Link Pay Vendor Run Payroll AI suggest + AI chat panel glassmorphic Swarm trace Tool get_tpv → ETB 125,430 + create_payment_link + RAG ask 5000 per ONPS/10/2025 citation

Outstanding UI: glassmorphic nav backdrop-blur-xl bg-white/70 border white/50 shadow-glass, card border 1px rgba(0,0,0,0.06) elevated hover medium, DonutProgress 64px progress 0-100% animated, stepper animated line pathLength motion.div, timeline Linear style vertical line + dot + card + timestamp

## Payments • ክፍያዎች

- List 7 cols Tx Ref mono Amount bold Method telebirr/cbe + connector_id small Status succeeded green failed red Routing telebirr primary fallback + 2FA verified pending + Action View • እይ link to detail exam timeline ledger M1 balanced + routing decision + webhook deliveries HMAC
- Detail Exam per SAD A1: lifecycle vertical timeline 4 steps created→pending routed via telebirr primary health success 96% latency 210ms score 0.88 chosen true reason primary healthy → pending→processing connector Initialize mock_ref checkout_url → processing→succeeded connector Verify succeeded ledger M1 journal posting_key payment_success:pay_01H balanced true → succeeded→webhook pending outbox payment.succeeded + webhook delivery success 200 attempt1 HMAC valid per secret prefix whsec_ + Ledger Journal M2 Dr payable R-FR + Dr fee_due FR Cr clearing R + Fee Policies non_refundable/pro_rata/full + Idempotency (merchant_id, refund_ref) unique 409 + Actions Resend Webhook Evidence Pack JSON NBE

## Links • Payment Links Outstanding

- Create Link: amount chips 100/500/1000/5000 selected bg primary, custom amount input, description AI suggest, expiry optional, Generate Link button → public_token tok_abc checkout_url https://checkout.apexpay.et/c/abc QR data URL base64 mock_qr + share telegram https://t.me/share/url?url=... whatsapp https://wa.me/?text=Pay ETB 500... + copy haptics + QR preview live
- List QR thumbnails status active/paid expired cancelled + token abc123 + share menu Telegram WhatsApp copy

## Refunds • ተመላሽ ክፍያ — FULL M2

- Create bottom sheet payment select txr_01H ETB 500 succeeded, Refund Ref unique (merchant_id, refund_ref), Amount partial allowed, Fee Policy non_refundable default/pro_rata reverse pro-rata fee*refund/pay Round2 ETB scale bankers rounding decimal precise/full full fee reversal on full refund, Reason, button Refund • M2 Dr payable R-FR + Dr fee_due FR Cr clearing R + Idempotency unique conflict 409 if different amount + Remaining refundable = amount - refunded
- List: Refund Ref Payment Amount Fee Rev Status Ledger M2 balanced + Payment status CASE WHEN sum(amount)>=amount THEN refunded ELSE partially_refunded

## Payouts • ክፍያዎች ለአቅራቢዎች

- Single payout maker-checker >50k ETB beneficiary select CBE ****1234 verified Awash ****5678, amount 10000, Payout Ref unique, Create Payout pending_approval if >50k, Ledger M3 Dr payable Cr clearing bank atomic per batch book
- Bulk CSV Upload 1000 rows preview outstanding GitHub Actions timeline: dropzone dashed 2px rounded-2xl border-2 pulse on drag scale 0.98 pulseGlow, preview table 5 cols Name Account Bank Amount Status ✓ valid ⚠ name mismatch Levenshtein 2 require override note bank_name CBE Commercial Bank of Ethiopia, total ETB 30k 3 beneficiaries maker-checker required balance check sufficient, Create Batch pbat_01H pending_approval dual approve >50k + Batches list pending_approval finance submitted admin approve needed + succeeded ledger M3 balanced + Real-time-ish SWR poll 5s

## Payroll • ደሞዝ • Workforce Money OS

- Total Gross 200k, Total Tax 20k ET brackets binary search O(log n) 0-600 0% etc effective_from versioned per migration 0008 seed, Total Net 150k Pension 7%/11% OT 1.25/1.5/2.0
- Employees 10 Fayda badge EMP001 Base 20k CBE ****1234 Sales, Import CSV 500 emps <2s p99 per NFR
- Runs table status pipeline visual stepper: Run Ref prun_July2026 Period 07/2026 Type regular Status pending_approval Total Net 150k Action View Calculate → Approve dual >100k → Disburse → payout batch
- Run Detail table 8 cols Employee Gross OT Taxable Income Tax ET Pension 7%/11% Net Status sticky footer totals + Approve dual + Disburse ledger M4 Dr salary 200k Cr payroll_payable 150k Cr et_income_tax_payable 20k Cr pension_payable 30k balanced + second journal Dr payroll_payable Cr bank + Payslip PDF modern QR verification breakdown pie + ET CSV ERCA JSON + ZIP download

## Developers • ገንቢዎች

- API Keys Test/Live Separate Scopes Reveal Once + Public Keys pk_test_ + Embedded SDK checkout.js outstanding + Quickstart 6 lines copy-paste like ApexPay + Webhook Endpoints HMAC + SSRF block + Retry exponential backoff + Bank List GET /v1/banks 14 ET banks + Methods Health GET /v1/methods?amount=1000 ranked score + OpenAPI Swagger Embedded Modern

## Tracking Timeline • የጊዜ መስመር

Outstanding vertical like Linear: Business info saved draft • 1 min ago + Fayda verification OTP sent 123456 mock face score 0.92 + Bank added CBE ****1234 verified + Docs 4/6 uploaded 66% hash integrity + Compliance checks 8 green/red risk 42 medium dual approval needed if high + Submitted pending_kyc → kyc_in_review → fayda_pending → compliance_check → pending_approval → active Progress 100% Test keys immediately live after pilot 30-60 days

Next: [Compliance Center RAG Citations](/docs/compliance-rag)
