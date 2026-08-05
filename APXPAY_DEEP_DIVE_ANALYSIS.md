# ApexPay Deep-Dive Analysis — Full Platform v1.1.0 Gold Review
**Date:** 2026-08-05 Africa/Addis_Ababa
**Repo:** https://github.com/melonadawit/apexpay.git
**Analyst:** Arena AI Agent Mode — Senior Engineer Audit
**Commit:** Analysis before any code change

---

## 0. Executive Summary

**ApexPay vision is correct and defensible:** Ethiopian-first Stripe + Modern Treasury + Gusto + AI Operator, built for NBE PO Gateway ONPS/02/2020, 09/2023 tiered KYC, ONPS/10/2025 2FA >5000 ETB, Fayda National ID verification.

**Codebase reality:** 150+ files, 12 migrations, Go chi API, 4 Next.js apps, Flutter 3.22, Python RAG. Money-path design is excellent (decimal, advisory locks, double-entry invariant), but commercial implementation is **~75-80% skeleton, ~40% of critical integrations are mocked**.

Docs claim 100% Gold — reality:
- Ledger, onboarding domain, payroll tax binary search, routing priority sort: **A+ done**.
- Connectors (Telebirr, CBE Birr, Bank IPS, EthSwitch, Card): **Only mock exists**.
- Ledger book resolution hardcoded `merchant_operating_default` instead of per-merchant lookup.
- Webhook SSRF block returns false always — security hole.
- RAG + Swarm are rule-based mocks, not real LLM/e5, Go proxy returns static answer.

**Verdict:** Strong foundation. Needs 2-weeks hardening to be pilot-ready with NBE, 8-10 weeks to be truly outstanding & powerful.

---

## 1. What ApexPay Is — Domain Understanding

### Ethiopian Context Nailed
- **Fayda:** FIN 12-digit hashed sha256(salt+FIN)+last4 only, FAN 16 alias, front/back <2MB, selfie liveness, OTP consent via id.gov.et VeriFayda 2.0 / OIDC eSignet, offline QR FaydaEncode, face match 0.85, MinIO encrypted vault presigned 15m, no plain FIN in logs.
- **NBE Directives:** ONPS/02/2020 gateway operator, ONPS/09/2023 tiered KYC level1 <=5000 balance, ONPS/10/2025 2FA PIN/OTP/biometric >5000 ETB, AML 200k ETB FIC reporting, PayAtlas PSP checklist 14 doc types.
- **Rails fragmentation:** Telebirr dominant, CBE Birr, Bank IPS ISO20022 pain.001, EthSwitch interoperable QR (EMVCo TLV), card token-only.
- **Smart routing:** `score = 0.6*success_rate_5m + 0.4*(1-latency/1000)`, health sampler 30s, circuit breaker 5 fails open 60s, fallback trail audited.
- **Ledgers M1-M4:** 
  - M1 Payment: Dr asset:clearing:{connector} Amount Cr liability:merchant_payable Net + Cr platform_fee_due Fee
  - M2 Refund: Dr payable (R-FR) + Dr fee_due FR Cr clearing R
  - M3 Payout: Dr payable Amount Cr clearing:bank Amount
  - M4 Payroll: Dr salary expense Gross Cr payroll_payable Net Cr et_income_tax Tax Cr pension_payable Pension — per-run book.

### Product Scope Full Spec v1.1.0
- Collect: one-off payments, payment links public_token QR, hosted checkout 420px glass, 2FA OTP, subscriptions plans/trial/dunning 1d/3d/5d, customer portal magic link JWT 24h
- Store: ledger 8 book types merchant_operating rail_clearing platform_revenue payroll_run payout_batch escrow suspense reserve
- Disburse: beneficiaries fuzzy Levenshtein <3, single/bulk 1000 CSV papaparse, payout links escrow, maker-checker dual >50k payout >100k payroll
- Workforce: employees Fayda hash bank masked, payroll_runs draft->calculating->pending_approval->approved->processing->completed, ET tax brackets 7 versioned binary search O(log n), pension 7%/11%, OT 1.25 weekday 1.5 weekend 2.0 holiday 1.3 night
- Intelligence: RAG pgvector 1024 chunk 800 overlap 100 embed batch 32 threshold 0.65 no hallucination guard, citations mandatory, Swarm planner RulesPlanner keyword critic confirmation >100k executor ToolExecutor O(1) registry create_payment_link create_payout calculate_payroll get_tpv ask_compliance
- Mobile Flutter: dashboard TPV glass gradient emerald + sparkline shimmer + quick actions FAB, create link bottom sheet draggable, QR scanner overlay 260 corner brackets pulse green + FaydaEncode + EthSwitch QR + vibration, approvals swipe + local_auth biometric + confetti Lottie + FCM + offline Hive draft_links offline_queue sync badge

---

## 2. Architecture Deep Dive — Real Implementation

### Backend Go API `services/api`
- `cmd/api/main.go` 451 lines REAL wiring: pgxpool, redis, minio EnsureBucket, 10 repos, 8 services, chi router RequestID RealIP Recoverer Timeout 15s Compress, /healthz version 1.1.0-full + fayda_mode, /readyz db ping, /metrics prometheus, public /v1/banks real query banks table fallback CBE/Awash/Dashen, /v1/payment_links/public/{token} real linkRepo.GetByToken, /v1/onboarding + fayda protected authMw.APIKeyAuth, /v1/transactions, /links, /refunds, /customers (stub), /subscription_plans (hack r2 router), /subscriptions, /beneficiaries (delegates payout handler), /payouts, /payout_batches (delegates), /employees, /payroll_runs, /webhooks, /compliance/ask STATIC answer "5000 ETB per ONPS/10/2025", /swarm/run static, /methods ranked, /devices/register static, /admin onboarding queue real query + review stub + connectors/health real query fallback + recon breaks [] + evidence bundle hashes.

- **Platform layers:**
  - `config/config.go` no silent defaults fail fast ✅
  - `crypto/crypto.go` ValidateFIN 12-digit regex, FAN 16, TIN 10, HashFIN sha256(salt+fin)+Last4+MaskAccount Encrypt AES-GCM ✅
  - `storage/minio.go` EnsureBucket PresignedPutUrl 15m PresignedGetUrl 15m UploadWithHash sha256 streaming O(n) ObjectKey merchants/{id}/kyc/{type}_{id}.pdf ✅
  - `http/response.go` success bool + data + error code/message + request_id correlation per SAD §11 ✅
  - `middleware/auth.go` APIKeyAuth lookup prefix index api_keys_prefix_uidx O(1) + RBAC map O(1) owner/admin/developer/finance + token bucket rate limit skeleton
  - `id/id.go` ULID prefixed mer_ kyc_ own_ fayda_ pay_ ref_ cust_ splan_ sub_ ben_ pbat_ pout_ emp_ prun_ lbk_ ljrn_ O(1) ✅

- **Core Domains:**
  - `ledger/domain.go` ValidateBalanced debit==credit O(n), repository PostJournalTx advisory lock pg_advisory_xact_lock(hashtext(book_id)) atomic journal+entries+balances upsert PK O(1) `GetBalance` PK O(1) ListJournalsByRef ref index, tests property 10k iter deterministic seed 42 + TestNoFloatMoney grep float64 amount + TestPostingKeyUniqueness map O(1) + Benchmark p99<30ms ✅ excellent
  - `onboarding/domain.go` BusinessType sole/PLC/share 5 shareholders min 10 multi max 40% single, Industry restricted gambling/crypto/adult map, Status FSM draft→submitted→in_review→fayda_pending→compliance_check→needs_more_info→approved, RequiredDocs() unique O(n) map, Service CreateKYC TIN 10 validation restricted 403 AddOwner ownership 0-100 SubmitKYC hasAuthSignatory + faydaVerifiedCount>=1 + settlement bank + docs missing check O(n) map + risk scoring weighted sum + PEP*30 + TPV high +20, Dual approval risk>=70 or TPV>1M ✅
  - `fayda/domain.go` VerificationMethod otp/face/fingerprint/offline_qr/oidc_esignet Status initiated/otp_sent/verified/failed, Service Init validates FIN/FAN regex hash immediately last4 requestID ULID verifier.RequestOTP 200ms creates otp_sent ConfirmOTP idempotent verified return OTP expiry 5m demographics_match bool face_score encrypted ref MinIO UpdateVerificationResult VerifyOfflineQR, Verifier mock/live Strategy pattern, MockVerifier RequestOTP mock_tx VerifyOTP 123456 always success 0.92 0.88 000000 fail VerifyOfflineQR split | format, LiveVerifier placeholder id.gov.et POST /auth/otp + /auth/verify + offline QR NIDP public key + OIDC eSignet exchange code id_token JWT verify ✅ mock gold, live TODO
  - `connector/domain.go` Connector interface Initialize/Verify/Refund/Health ID(), `mock.go` 50ms checkout_url mock_ref InitializeResponse VerifyResponse ✅ only mock
  - `routing/domain.go` ConnectorID mock/telebirr/cbe_birr/bank_ips/ethswitch/card_acquirer Strategy success_rate/latency/cost/round_robin RoutingRule priority HealthSample ConnectorHealth success_rate_5m avg_latency_5m uptime24h circuitState RoutingDecision chosen+fallbacks+reason+healthSnapshot Service Evaluate filter amount/currency/method/enabled sort priority asc O(n log n) selected candidates[0] getHealth list samples since 5m success rate avg latency circuitState from map circuit open check 60s strategy success_rate if primary <0.7 pick best fallback max success_rate latency min latency RankedMethods all connectors success*0.6+(1-latency/1000)*0.4 sort Score desc RecordFailure failures++ >=5 open 60s RecordSuccess half_open→closed ValidateNoCycle primary!=fallback1!=fallback2 ✅ excellent algorithm, integration missing health sampler worker ticker 30s described but not in main.go api only worker cmd
  - `payment/domain.go` Status etc, Service Initialize idempotency by key duplicate tx_ref unique (merchant_id,tx_ref) routing decision fallback mock Fee=amount*0.029 rounded ETB scale 2 Requires2FA if ETB>5000, connector.Initialize audit circuit RecordSuccess/Failure CreatePaymentTx with outbox payment.created IdempotencyKey save, Verify idempotent second success no-op per MVP B6 if requires 2FA not verified return pending connector.Verify VerifyResponse status succeeded Ledger M1 posting atomically journal posting_key payment_success:{id} entries filter zero fee optimization UpdateStatusTx ✅ but BookID hardcoded merchant_operating_default not per-merchant lookup
  - `refund/domain.go` Status created/processing/succeeded/failed FeePolicy non_refundable/pro_rata/full Service CreateRefund idempotency by (merchant_id,refund_ref) unique + remaining refundable check + ledger M2 Dr payable (R-FR) + Dr fee_due FR Cr clearing R filtered zero + mock connector success + UpdateRefundStatus ✅, payment status partially_refunded/refunded second step missing? need check
  - `subscription/domain.go` Plans amount currency interval_type trial_days status Customers email phone name fayda hash Subscriptions FSM incomplete,trialing,active,past_due,canceled paused current_period_start/end trial_end Invoices draft/open/paid attempt_count due_at Service CreatePlan amount>0 interval_count default 1 CreateSubscription trial handling if trialDays>0 status trialing period = trialEnd else interval add day/week/month/year optimal time.AddDate Invoice draft/open/paid dunning NextDunningAttempt 1d/3d/5d exponential-ish Worker dunning cron hourly scans invoices due + past_due attempts mock update attempt_count next_attempt_at webhook subscription.* ✅ skeleton
  - `payout/domain.go` Beneficiary name masked hash bank_code bank_name type verification_status PayoutBatch BookID BatchRef amount currency status approved_by Payout PayoutRef amount currency status method PayoutStatus created->pending_approval->queued->processing->succeeded|failed|returned Service CreateSingle checks merchant balance decimal precise GetMerchantBalance COALESCE sum net_amount succeeded - sum amount queued/processing/succeeded ApprovalThreshold 50k -> pending_approval else queued ledger M3 Dr payable Cr clearing_bank CreateBulk validates bulk 1-1000 O(n) sum total <= balance creates batch pending_approval (all bulk require approval per policy) payouts created status created journal Dr payable total Cr clearing total per batch book CSV parser papaparse frontend + backend csv preview outstanding validation icons amount sum fees calc MDR Payout links escrow book claim via OTP move escrow->clearing Maker-checker dual approval >50k check approver != submitter ✅ fuzzy Levenshtein claimed not found
  - `payroll/domain.go` Employee fin_hash bank masked hash base_salary employment_date cost_center EmploymentType permanent/contract/part_time RunType regular/off_cycle/bonus/adjustment RunStatus FSM draft->calculating->pending_approval->approved->processing->completed->failed PayrollRun BookID per run PeriodMonth/Year Type Status Totals ApprovedBy PayrollItem gross OT hours/amount commission taxable income_tax pension 7%/11% net TaxBracket Min Max Rate Deduction EffectiveFrom OT rates map 1.25/1.5/2.0/1.3 OTType weekday WeekdayWeekend Holiday Night Service CalculateTax binary search O(log n) over sorted brackets 7 ET 2024 0-600 0% 0 601-1650 10% -60 1651-3200 15% -140 3201-5250 20% -300 5251-7800 25% -565 7801-10900 30% -955 >10900 35% -1500 formula tax=taxable*rate-deduction rounded 2 decimals Pension 7% gross employee 11% employer OT hourly_rate = base/208 (26*8) ET standard CalculateRun loops employees O(n) active only gross=base+OT+commission pensionEmp taxable=gross-pensionEmp - exemptions incomeTax via binary search deductions = pensionEmp+incomeTax+other net=gross-deductions aggregates totals State machine draft->calculating bulk insert items Tx ->pending_approval->approved dual >100k net ->processing->completed Ledger M4 per run book creation ledger_books book_type payroll_run journal posting_key payroll_run:{id} Dr expense:salary totalGross Cr payroll_payable totalNet Cr et_income_tax_payable totalTax Cr pension_payable totalPension ValidateBalanced before posting O(k) Disburse creates payout batch for employees banks second journal Dr payroll_payable Cr clearing:bank totalNet via payouts ✅ calc excellent PDF pending
  - `swarm/domain.go` SessionStatus planning/executing/needs_confirmation/completed/failed/cancelled SwarmSession MerchantID UserID Goal Plan []PlanStep Status ConfirmationRequired ConfirmationData FinalOutput PlanStep Step Tool Description Args Status Result AgentRun MerchantID ThreadID SwarmSessionID InputText Intent ToolCalls []ToolCall OutputText Model rules_v1 ToolCall Tool Args Result LatencyMS Status ToolDefinition Name Description ArgsSchema Threshold RoleAllowed Registry map O(1) tool definitions create_payment_link threshold 100k etc planner RulesPlanner keyword critic confirmation_required >100k executor ToolExecutor JSON schema validation state machine planning→executing→needs_confirmation→completed Repository CreateSession planBytes json Marshal GetSession UpdateSession CreateAgentRun Service Run goal → needs_confirmation if >100k + Confirm() Executor ToolExecutorImpl service injection function map optimal O(1) paymentLinkCreator payoutCreator payrollCalculator tpvGetter tpv_today 125430 count 42 complianceAsker answer 5000 ETB per ONPS/10/2025 Execute switch tool create_payment_link amount float64/int handling currency default ETB etc latency recorded ✅ logic good but RulesPlanner only keyword, no LLM
  - `link`, `webhook`, `rag` domains exist: rag Document chunk embedding vector(1024) ivfflat lists=100 AskRequest Lang en/am MerchantID TopK Citation doc_title page score url Domain pgvector service embed batch 32 threshold 0.65 guard ✅
  - `webhook/service.go` Delivery MerchantID EndpointID EventType Payload URL Secret AttemptCount Repository ListPendingDeliveries limit 100 MarkSuccess MarkFailed statusCode errMsg nextAttempt client Timeout 10s Sign HMAC SHA256 hex Deliver SSRF block isPrivateURL 127.0.0.1 10.0.0.0/8 172.16.0.0/12 192.168.0.0/16 placeholder returns false always ❌ backoff 1s 2s 4s 8s 16s 32s max 1h ✅
  - `worker/cmd/worker/main.go` config Load, logger, pgxpool, redis client, goroutines outbox drain ticker 1s SELECT WHERE published_at NULL ORDER BY created_at LIMIT 100 FOR UPDATE SKIP LOCKED create webhook_deliveries for merchant active endpoints + mark published_at now, webhook retry ticker 2s ListPendingDeliveries status pending/failed where next_attempt_at <= now + Deliver HMAC + exponential backoff SSRF block, health sampler ticker 30s per spec for each connector config enabled call Health latency_ms success bool insert health_samples + Redis cache health:{connector} TTL 60s O(1), dunning ticker 1h SELECT subscription_invoices status open next_attempt_at <= now attempt_count<3 attempt payment via saved method mock update attempt_count next_attempt_at = NextDunningAttempt, recon daily 02:00 Africa/Addis_Ababa calculate next 02:00 sleep Until next, fetch bank statements MinIO SFTP mock parse MT940/csv insert recon_statements matching engine amount tolerance 0.01 ETB window 24h O(n+m) map connector_ref->journal ✅ per FINAL_BUILD_STATUS description, need verify file exists full implementation

### Frontend
- `merchant-web`: Next14 package.json framer-motion tailwind-merge clsx cva radix-ui dialog dropdown label select progress toast tabs avatar lucide recharts react-dropzone qrcode.react sonner zod react-hook-form. `app/page.tsx` outstanding landing glass sticky header Start Onboarding, hero AI-native payment gateway for Ethiopia, 3 GlassCards Fayda Smart Routing Payroll ET, Card NBE checklist. `app/dashboard/page.tsx` client motion GlassCard DonutProgress TPVRecharts HealthRecharts TPV today 125430 success rate 96.2% active links 12 recent payments 3 items succeeded 2FA routing fallback_used glass AI Chat Swarm trace 2 steps tool get_tpv create_payment_link RAG Ask. `app/dashboard/recharts.tsx` AreaChart/Line. `app/onboarding/page.tsx` wizard 6 steps BusinessInfoStep OwnersFaydaStep BankAccountStep DocumentsVaultStep CompliancePreviewStep ReviewSubmitStep. Components onboard: BankAccountStep combobox bank list CBE/Awash/Dashen logos account name auto-check must match legal, BusinessInfoStep smart form industry search restricted gambling/crypto/adult blocked per NBE, FaydaCapture modal corner guides animated pulse glare detection, OnboardingWizard stepper animated line, DocumentsVaultStep dropzone dashed pulse on drag preview thumbs hash integrity progress donut, CompliancePreview risk gauge chart green/red timeline, ReviewSubmit confetti terms consent NBE. `lib/papaparse-bulk.ts` bulk CSV, `lib/pdf/payslip.ts` jsPDF payslip modern. `app/payments/page.tsx` mockPayments array hardcoded 3 items succeeded/failed 2FA not API ❌, `app/payments/[id]/page.tsx` likely same, `app/links/page.tsx` ?, `app/refunds`, `app/subscriptions`, `app/payouts`, `app/payroll` exists per find but need to check content is mock.
- `checkout-web`: `app/c/[token]/page.tsx` token page + `app/page.tsx` likely outstanding mobile 420px glass method selector radio icons Telebirr/CBE best route badge tooltip fallback trail 2FA OTP input processing lottie success confetti receipt PDF — placeholder per docs.
- `admin-web`: `app/page.tsx` placeholder — need Kanban drag-drop merchant exam split pane KYC left checks right Fayda images blurred docs viewer side-by-side OCR.
- `docs`: docusaurus config, intro, quickstart, onboarding, fayda, payments, compliance-rag, swarm, comparison, merchant-guide, api-reference — real docs ✅.
- `mobile`: pubspec.yaml flutter 3.22 riverpod_annotation go_router dio hive hive_flutter secure_storage camera qr_code_scanner mobile_scanner share_plus firebase_core firebase_messaging local_auth lottie shimmer cached_network_image intl google_fonts flutter_svg equatable json_annotation collection. `main.dart`, `core/api/api_client.dart` dio interceptor, `fcm_service.dart`, `router/app_router.dart` go_router, `storage/hive_boxes.dart` HiveBoxes init draft_links offline_queue sync badge, `sync/offline_sync.dart` TODO, `theme/app_theme.dart` ET Green #0B6E4F gold #EAB308. Features: `auth/presentation/login_page.dart` glass card token secure_storage refresh interceptor biometric, `dashboard/presentation/dashboard_page.dart` glass gradient emerald sparkline recent payments shimmer pull-to-refresh empty state coffee illustration quick actions FAB, `links/presentation/create_link_sheet.dart` bottom sheet draggable half/expanded amount chips 100/500/1000 AI suggest QR preview live share system copy haptics, `qr/presentation/qr_scanner_page.dart` camera permission dialog overlay rounded 260 guides corner brackets pulse green supports FaydaEncode offline QR + EthSwitch QR + vibration, `approvals/presentation/approvals_page.dart` inbox pending payouts/payroll runs swipe approve/reject local_auth biometric confetti Lottie push FCM token registration POST /v1/devices/register, `profile/onboarding/presentation/onboarding_wizard_page.dart` 6-step PageView dot indicator Fayda capture modals camera overlay OTP pin animated FIN/FAN mask docs dropzone thumbs compliance gauge review NBE consent + camera overlay glare detection brightness >200, `offline_queue` box sync on reconnect idempotency key TODO, startup <2s profile deferred components.

### Database
- `db/migrations/0001_init.up.sql`: extension pgcrypto uuid-ossp, merchants status draft..active onboarding_status, users, merchant_users role owner/admin/developer/finance/support/ops/compliance/viewer, api_keys key_prefix secret_hash environment test/live status active/revoked scopes last_used_at async best effort, payments merchant_id tx_ref unique amount numeric(20,8) currency status created/pending/processing/succeeded/failed/canceled/refunded/partially_refunded method description customer_email connector_id connector_ref routing_rule_id checkout_url return_url callback_url metadata jsonb fee_amount net_amount failure_code requires_2fa two_fa_verified succeeded_at created_at updated_at indexes merchant_created, status, connector_ref, payment_links id merchant_id payment_id amount status active/paid/expired/cancelled public_token unique expires_at, checkout_sessions public_token unique, ledger_books id merchant_id book_type merchant_operating/rail_clearing/platform_revenue/payroll_run/payout_batch/escrow/suspense/reserve/refund_clearing/sandbox name currency status open/closed, ledger_accounts book_id code name normal_balance debit/credit unique book_id code, ledger_journals id book_id posting_key unique book_id posting_key memo transfer_group reference_type reference_id, ledger_entries id journal_id book_id account_id direction debit/credit amount currency meta jsonb, ledger_balances book_id account_id primary key amount updated_at, outbox_events id merchant_id aggregate_type aggregate_id event_type payload created_at published_at index unpublished, webhook_endpoints id merchant_id url secret_hash secret_prefix status active/disabled events default ["*"], webhook_deliveries id merchant_id endpoint_id outbox_event_id event_type payload status pending/success/failed/dead attempt_count next_attempt_at last_status_code last_error, idempotency_keys merchant_id key primary key request_hash response_code response_body resource_type resource_id, agent_runs id merchant_id thread_id swarm_session_id input_text intent tool_calls jsonb output_text model rules_v1 created_at merchant_created, audit_logs id merchant_id actor_type actor_id action resource_type resource_id ip inet request_id data jsonb merchant_created resource.
- `0002_onboarding_kyc.up.sql`: merchant_kyc_profiles id merchant_id version unique merchant_id version legal_name trade_name business_type sole/plc/share/partnership/cooperative/government/ngo/other registration_number registration_date tin_number check char_length 10 vat_number business_license_no license_expiry industry_category business_description website_url app_url annual_turnover expected_monthly_tpv avg_ticket_amount employee_count region city sub_city woreda kebele house_no office_address_full gps_lat/lng contact_person_name role contact_email phone has_refund_policy has_privacy_policy has_terms_and_conditions refund_url privacy_url terms_url onboarding_status draft/submitted/in_review/fayda_pending/compliance_check/needs_more_info/approved/rejected kyc_level level1/2/3 risk_notes submitted_at reviewed_at created_at updated_at indexes merchant_status. merchant_beneficial_owners id merchant_id kyc_profile_id full_name full_name_am role owner/shareholder/director/authorized_rep/contact_person/ubo ownership_percentage 5,2 0-100 nationality char2 default ET id_type fayda/passport/driving_license/kebele_id/other id_number_hash id_number_last4 fayda_fin_hash fayda_fan fayda_verified date_of_birth gender phone email address is_pep is_authorized_signatory verification_status pending/fayda_verified/document_verified/verified/rejected. bank_accounts id merchant_id account_name masked hash bank_code bank_name branch account_type current/saving is_settlement_default verification_status pending/verified/failed verification_method bank_letter/micro_deposit/manual verified_at created_at.
- `0003_onboarding_docs_compliance.up.sql`: documents? compliance checks? etc — not read but per spec 14 doc types + hash integrity + file_key MinIO.
- `0004_fayda_verification.up.sql`: fayda_verifications table request_id otp expiry demographics_match face_score.
- `0005_refunds.up.sql`: refunds id merchant_id payment_id refund_ref unique merchant_id refund_ref amount currency status created/processing/succeeded/failed reason fee_reversal connector_id connector_ref.
- `0006_subscriptions.up.sql`: customers, subscription_plans, subscriptions FSM, subscription_invoices.
- `0007_payouts.up.sql`: beneficiaries, payout_batches book per batch, payouts FSM.
- `0008_payroll.up.sql`: employees, payroll_runs book per run, payroll_items, payroll_tax_brackets 7 seeded 0-600 0% etc versioned effective_from, OT rates weekday 1.25 weekend 1.5 holiday 2.0 night 1.3.
- `0009_connectors_routing.up.sql`: connector_configs encrypted secrets AES-GCM env CONNECTOR_ENCRYPTION_KEY, routing_rules merchant_id nullable method amount_range primary fallback1 fallback2 strategy latency/cost/success_rate priority sort no cycle validation, connector_health_samples latency_ms success bool sampled_at.
- `0010_swarm_recon.up.sql`: swarm_sessions plan json status confirmation_required confirmation_data final_output, recon_statements parser MT940/csv/json + matching engine amount tolerance 0.01 ETB window 24h O(n+m) map + suspense posting + cron daily 02:00 Africa/Addis_Ababa.
- `0011_rag_pgvector.up.sql`: pgvector extension, rag_documents source_type title url hash status, rag_chunks document_id content embedding vector(1024) ivfflat lists=100 metadata.
- `0012_enhancements.up.sql`: banks seed 14 ET banks CBE/Awash/Dashen/Abyssinia/Berhan/Wegagen/NIB/United/Coop/Oromia/Bunna/Lion/Zemen/CBO, materialized view merchant_tpv_daily, indexes per access pattern, etc.

## 3. What Works Excellent

- Ledger invariant, decimal precise, advisory locks, ULID, clean arch Strategy Factory Singleton pooling.
- Onboarding domain RequiredDocs unique set O(n) map + risk scoring weighted sum.
- Payroll tax binary search O(log n) optimal, pension 7/11, OT rates map O(1).
- Routing priority sort O(n log n) + health 5m window O(k) k<=3 + circuit breaker map O(1) lookup + score 0.6*success+0.4*(1-latency/1000) + no cycle validation.
- Chunking tiktoken cl100k_base token count O(n) sliding window 800 overlap 100 sentence boundary, embedding batch 32 L2 normalized cosine e5 prefix query:/passage: best practice.
- Vector search pgvector ivfflat lists=100 O(log n) cosine, threshold 0.65 guard prevents hallucination topK 5.
- Swarm registry map O(1) RulesPlanner keyword O(n) confirmation_required >100k state machine.
- Fayda glare detection canvas 100x100 sampled every 400ms brightness average + bright pixel ratio >15% triggers warning O(n/4) 16-step sampling — outstanding spec.
- Design tokens final `libs/ui/tokens.json` ET Green #0B6E4F gold #EAB308 radius xl 24 motion ease [0.22,1,0.36,1] shared Tailwind + Flutter ThemeData same palette + Framer Motion variants fade+slide+scale stepper progress animated pathLength AnimatePresence file upload dropzone pulse border skeleton shimmer 2s staggered 50ms*index confetti Lottie.

## 4. Critical Gaps & Missing

### P0 Blockers for NBE Pilot
1. **Connector registry only mock** `map[string]Connector{"mock": NewMock()}` in main.go — no telebirr_sandbox HMAC X-APP-Key + callback IP allowlist + latency inject 150-300ms + retry exponential backoff + CBE Birr 30% failure inject + Bank IPS ISO20022 mock pain.001 generator + EthSwitch QR validator + card_acquirer token-only. Config `connector_configs` encrypted AES-GCM key `CONNECTOR_ENCRYPTION_KEY` exists but not used.
2. **Ledger book hardcoded** `BookID: "merchant_operating_default"` in payment/service.go:110 refund/service.go payroll/repository.go — must resolve via `ledgerRepo.GetOperatingBook(merchantID)` query `ledger_books where merchant_id=? and book_type=merchant_operating and status=open order created_at desc limit 1` + Redis cache health:{connector} TTL 60s O(1) already but book cache missing.
3. **Webhook SSRF** `isPrivateURL()` placeholder returns false always — allow attacker webhook to `http://169.254.169.254/latest/meta-data/` or `10.0.0.0/8` internal. Need deny `127.0.0.0/8 10/8 172.16/12 192.168/16 169.254/16 ::1 fe80::/10` + DNS rebinding check resolve IP after.
4. **Auth secret compare** `api_keys` secret_hash bcrypt/argon2 hash stored prefix visible but `Auth` middleware lookup by prefix index O(1) but secret_hash compare missing? Best effort last_used_at async via goroutine okay but need constant-time compare.
5. **Idempotency race** `CreatePaymentTx` + `SaveIdempotency` not same Tx — parallel same Idempotency-Key can create duplicate payments. Need pgx Tx `BEGIN; SELECT pg_advisory_xact_lock(hashtext(merchant_id+key)); INSERT idempotency_keys ON CONFLICT DO NOTHING RETURNING; INSERT payments; COMMIT`.
6. **Public link enumeration** `/payment_links/public/{token}` no expiry check `expires_at < now() => 410 Gone` + rate limit token bucket per IP 100/min.
7. **File hash integrity** `UploadWithHash` sha256 streaming O(n) exists but documents table unique index per merchant file_hash? Need dedup + ClamAV stub VirusScanner interface clean.

### P1 High for Commercial Core
8. Fayda live verifier `LiveVerifier` placeholder — need real HTTP client to id.gov.et `/api/ekyc/otp` + `/api/ekyc/verify` + offline QR NIDP public key PEM EC P-256 verify RS256/ES256 signature + OIDC eSignet exchange code id_token JWT verify. OTP rate limit 5/hour/IP via Redis token bucket Lua `fayda:otp:{owner}` TTL 1h not enforced.
9. Bank verification fuzzy Levenshtein — claim Levenshtein <3 but not implemented. Need Go `levenshtein.Distance(name, legalName)` + bank letter OCR.
10. Payroll disburse atomicity — `DisburseRun` creates payout batch second journal Dr payroll_payable Cr clearing:bank totalNet via payouts but payrollRepo `CreateRunBookTx` and `BulkCreateItems` Tx begin loop insert but `DisburseRun` ledger M4 + payout batch not same Tx? Should be `BEGIN; INSERT ledger_journal; INSERT ledger_entries; INSERT payout_batches; INSERT payouts; COMMIT`.
11. RAG Go client not wired — main.go `/compliance/ask` returns static JSON not proxy to Python `http://rag:8001/v1/compliance/ask` embedding e5-large 1024 dim. Need `rag/service.go` Go client http client with timeout 5s + Redis cache query hash TTL 5m.
12. Mobile offline sync — `HiveBoxes` draft_links offline_queue sync badge init but `offline_sync.dart` TODO sync on reconnect idempotency key ULID + Dio retry interceptor.
13. Subscription saved payment method missing — dunning worker attempts payment via saved method mock but no `payment_methods` table tokenized.
14. Reconciliation parser TODO skeleton — MT940/csv/json parser + matching engine amount tolerance 0.01 ETB window 24h O(n+m) map connector_ref->journal + suspense posting + ops dashboard list assign resolve.
15. Maker-checker dual approval — risk>=70 or TPV>1M requires second approver for onboarding, payout >50k approver != submitter, payroll >100k — enforced in service but admin review handler stub returns static approved without checking duality + timeline.

### P2 Medium Outstanding -> Exceptional
16. USSD fallback — Ethiopia low internet, Telebirr USSD *127# flow — massive differentiator, checkout-web if offline show USSD code.
17. EthSwitch interoperable QR spec real — EMVCo TLV parser 00-99 tags.
18. Checkout-web outstanding UI — need build 420px glass + method selector icons Telebirr/CBE health dots green/yellow/red animated, 2FA OTP 6-digit input animated, processing Lottie, success confetti canvas-confetti full-screen 3s + haptics vibrate.
19. Admin-web Kanban — drag-drop dnd-kit, merchant exam split pane KYC left checks right Fayda images blurred + docs viewer side-by-side OCR + compliance checklist timeline.
20. i18n Amharic 100% checkout, 80% dashboard — language-provider.tsx exists but pages hardcoded mix.
21. Observability — prometheus client imported but metrics not instrumented connector_circuit_open gauge, payment_init_latency histogram, ledger_post_latency p99<30ms, payroll_calc_duration, webhook_delivery_latency, rag_query_latency. Add OpenTelemetry trace request_id merchant_id.
22. Tests — only ledger test property 10k iter, no table-driven routing 100% coverage, refund_exceeded, duplicate_ref, tax bracket known ET examples rounding .005, payroll balanced invariant 500 employees calc <2s bench, contract test OpenAPI diff.
23. Dispute/chargeback, escrow payout links claim via OTP move escrow->clearing, payslip QR verification endpoint GET /v1/payroll/payslips/{id}/verify QR + ERCA CSV.

---

## 5. How to Make More Comprehensive Outstanding Powerful

### Vision to Execution — 5 Phases 22-24w 2 Senior Eng (as per ROADMAP_M5-M9.md)

#### Phase M5.5 Commercial Hardening (2w)
- Fix ledger book resolver + cache Redis TTL 60s
- Implement telebirr_sandbox.go HMAC SHA256 sign = HMAC(secret, timestamp+body) header X-APP-Key X-Timestamp X-Signature callback IP allowlist CIDR env TELEBIRR_ALLOWED_IPS 10.0.0.0/24 + latency inject 150-300ms
- CBE Birr sandbox 30% failure inject configurable CBE_FAILURE_RATE=0.3 to test routing fallback
- Bank IPS ISO20022 pain.001 XML generator + recon file MT940 parser `func ParseMT940(r io.Reader) []Statement`
- SSRF hardening isPrivateURL parse ip + deny private ranges + DNS resolve check
- Auth fix bcrypt compare + last_used_at async
- Idempotency advisory lock hashtext(merchant_id+key) transactional
- Public link expiry 410 + rate limit 100/min IP token bucket Redis
- Beneficiary Levenshtein implementation `func Levenshtein(a,b string) int` O(n*m) DP if <3 accept else reject bank verification

#### Phase M6.5 NBE Gold Compliance (3w)
- Fayda live client: `fayda/live_verifier.go` RequestOTP POST /api/ekyc/otp Header PartnerCode+APIKey+UseCase FIN FAN file keys front/back/selfie + Response fayda_transaction_id otp_sent, VerifyOTP POST /api/ekyc/verify + demographics_match face_score 0.92, VerifyOfflineQR split | format FINLast4 NAME DOB SIG verify ECDSA with NIDP public key PEM file env FAYDA_NIDP_PUBLIC_KEY path, OIDC eSignet exchange code id_token JWT verify RS256 via id.gov.et/.well-known/jwks.json. Rate limit 5/hour/IP+owner Redis Lua.
- OCR pipeline Python rag-worker extended `ocr.py` PyMuPDF pdf -> text -> regex TIN 10-digit `^\\d{10}$` validate via NBE checksum? + business license expiry check + bank letter account name extract via NER? Save ocr_raw jsonb.
- Evidence pack export `GET /v1/admin/evidence?tx_ref=` bundle JSON merchant docs hashes not files payment lifecycle ledger entries fayda verification existence flag webhook deliveries audit_logs — per NBE exam.
- ERCA payroll report CSV `func GenerateERCACSV(items []PayrollItem) []byte` columns TIN taxable tax pension 7% 11% net cost_center + JSON export.
- Payslip PDF gold Go server-side `github.com/jung-kurt/gofpdf` + QR `boombuler/barcode/qr` verification URL `https://merchant.apexpay.et/verify/payslip/{hash}` breakdown table pie chart deductions.

#### Phase M7.5 Ethiopia Differentiators (2w)
- USSD fallback checkout-web `if (!navigator.onLine) show USSD *127*${amount}*${merchantCode}#` + copy button.
- EthSwitch QR real TLV parser per Ethiopian QR standard spec — implement `ethswitch/qr.go` `GenerateQR(merchantID, amount, ref)` `ParseQR(payload)`.
- Social commerce links Telegram/WhatsApp share — checkout-web outstanding share intents `telegram: whatsapp://send?text=` + `navigator.share` + haptics `navigator.vibrate(50)` + QR live preview canvas.
- Telegram bot: webhook `/webhooks/telegram` command `/create 100 coffee` -> creates payment_link via linkSvc.
- Employee self-service portal magic link JWT 24h view payslip claim expenses.

#### Phase M8.5 AI Moat Real (3w)
- RAG: Replace MockEmbedder with real model `sentence-transformers/multilingual-e5-large` 1024dim via `github.com/henomis/lingoose` or Python `sentence_transformers` batch 32 L2 normalized cosine prefix `query:` `passage:` per best practice. Use Qdrant option for scale but keep pgvector ivfflat lists=100 for MVP. Real LLM `OpenAI compatible` local `Mistral 7B` or `Claude` proxy with system prompt citation mandatory + no hallucination guard. Go client `rag/client.go` http Client timeout 5s cache Redis query hash TTL 5m. Streaming SSE for merchant chat Perplexity-like citations badges clickable PDF viewer. Eval harness 5 cases Amharic/English citation precision 0.8.
- Swarm: Planner LLM function calling constrained grammar JSON schema tools: create_payment_link threshold 100k desclang: create link for customer, create_payout threshold 50k, calculate_payroll_draft threshold 100k, get_tpv no threshold, ask_compliance, list_payments. Critic threshold checks policy ledger invariant ValidateBalanced, PEP count, amount > threshold requires confirmation. Executor ToolExecutor real injection domain services not mock. UI trace view outstanding steps timeline tool call cards Vercel AI SDK confirmation modal breakdown + biometric.
- Voice Amharic TTS RAG answer for merchant app.

#### Phase M9 Mobile + Gold Polish (4w)
- Design tokens final libs/ui/tokens.json primary ET Green #0B6E4F light #10A37A dark #094E38 50 #ECFDF5 accent gold #EAB308 yellow #FEF08A neutral zinc 50-900 semantic success #10B981 warning #F59E0B error #EF4444 glass white70 rgba(255,255,255,0.7) radius md 12 lg 16 xl 24 2xl 32 shadows soft 0 10 30 rgba(0,0,0,0.07) medium 0 20 40 large 0 30 60 font sans Inter Ethiopic Noto mono JetBrains Mono motion ease [0.22,1,0.36,1] duration fast 200 medium 300 slow 500 spring stiffness 300 damping 30. Sync Tailwind + Flutter ThemeData same palette.
- Motion Framer Motion variants fade+slide+scale stepper progress animated line pathLength motion.div onboarding wizard AnimatePresence Fayda overlay corner brackets animated pulse scale 1->1.1 infinite file upload dropzone pulse border checkout success confetti canvas-confetti full-screen 3s + haptic skeleton shimmer 2s infinite stagger list 50ms*index shimmer::after gradient 90deg.
- Merchant-web fix all mock pages SWR hooks useSWR('/v1/payments?merchant_id=...') empty states illustrations Ethiopian coffee ceremony Axum obelisk subtle data viz Recharts TPV success rate connector health latency line.
- Checkout-web outstanding 420px glass + method selector icons health dots green/yellow/red animated badge "Using best route: Telebirr (2% faster today)" tooltip fallback trail  2FA OTP 6-digit input auto-focus processing Lottie success confetti PDF receipt.
- Admin-web Kanban drag-drop dnd-kit + exam split pane.
- Flutter gold: offline sync Hive + FCM token registration + onboarding 6-step camera overlay glare detection brightness >200 warning Move to shade.
- Performance: k6 100 VUs 5m p95 init <300ms local ex-rail p99 <150ms staging payroll calc <2s for 500 emp RAG <1.5s payout bulk 1000 queued <10s. Lighthouse 90+ Perf 95 checkout A11y 100 BP 100 SEO 90. axe audit 0 serious. Flutter startup <2s.
- Security: gosec high 0 govulncheck trivy nancy no plain FIN in logs grep logs test CI PII redact zerolog field filter presigned 15m TTL ClamAV stub VirusScanner hash integrity file_hash unique per merchant encrypted SSE-S3 MinIO versioning retention 7y per NBE TLS everywhere secrets vault/env never git API keys pk_/sk_ hashed at rest bcrypt secrets at rest prefix visible later idempotency keys on payment and payout mutations maker-checker high-risk audit log including AI tool calls.

### 6. Immediate Action Backlog P0/P1/P2

**P0 (1w):**
- [ ] ledger_books resolver + cache
- [ ] SSRF fix
- [ ] Auth bcrypt compare
- [ ] Idempotency Tx
- [ ] Public link expiry + rate limit
- [ ] Telebirr + CBE mock with failure inject + HMAC structure

**P1 (2w):**
- [ ] Fayda live verifier HTTP client + NIDP public key
- [ ] Bank name Levenshtein
- [ ] Payroll disburse atomic Tx
- [ ] RAG Go client proxy to Python rag:8001
- [ ] Evidence pack JSON
- [ ] ERCA CSV + payslip PDF QR

**P2 (Outstanding):**
- [ ] Checkout-web outstanding 420px
- [ ] Admin-web Kanban + exam split pane
- [ ] Mobile offline sync + FCM
- [ ] USSD + EthSwitch QR TLV
- [ ] Swarm LLM planner + trace UI
- [ ] Observability metrics instrumentation
- [ ] k6 + Lighthouse axe CI green

---

## 7. How to Run Gold Demo Now

```bash
docker compose -f deploy/docker/docker-compose.yml up -d # postgres+pgvector 5432 redis 6379 minio 9000:9001 api 8080 worker merchant-web 3000 checkout-web 3001
make migrate-up # 0001-0012 clean from zero
go run ./scripts/seed/main.go # banks 14 + platform books 6 accounts + compliance user + mock Fayda + routing 4 + tax brackets 7 + RAG 3 samples
go run ./services/api/cmd/api # :8080 healthz readyz metrics
go run ./services/worker/cmd/worker # outbox drain 1s webhook retry 2s health sampler 30s dunning 1h recon 02:00 EAT
cd apps/merchant-web && npm i && npm run dev # :3000
cd apps/checkout-web && npm i && npm run dev -- -p 3001
cd apps/mobile && flutter pub get && flutter run
```

Demo 35 steps per FINAL_100_GOLD.md — create merchant register, KYC, Fayda init/confirm OTP 123456, bank, docs, submit, compliance queue approve dual, tracking timeline 100%, exam view, API key live sk_live shown once, payment_link amount 250 QR, checkout method selector Telebirr latency 200ms, initialize 6000 ETB requires_2fa true per NBE ONPS/10/2025 returns 2fa_session, verify OTP, connector success callback payment succeeded ledger balanced, GET /v1/methods ranked, kill telebirr health 0% fallback mock audit fallback_used, refund 100 partial ledger M2, subscription plans Monthly Coffee 500 interval month + customer + subs trial 7d invoice dunning, beneficiary + payout 1000 pending approval admin approve queued bank success ledger M3, employees x10 base 20000 + payroll_runs July regular calculate income_tax per ET bracket pension 7% 1400 employer 11% 2200 net approve disburse payout batch ledger M4 balanced, RAG ingest NBE directive ask When is 2FA required answer 5000 ETB per ONPS/10/2025 citations, swarm run goal Create link 100 ETB for coffee and run payroll July plan 2 steps needs_confirmation true for payroll >100k confirmation modal, Flutter login dashboard TPV create link bottom sheet share Telegram scan QR camera permission approvals inbox swipe approve biometric FCM token POST /v1/devices/register offline sync, recon upload statement csv GET /v1/admin/recon/breaks, evidence pack bundle, k6 smoke p95<300 ledger p99<30 payroll<2s RAG<1.5s.

---

## 8. Recommendation

You have built the hardest parts correctly: ledger, onboarding domain, payroll tax, routing algorithm, Fayda privacy. The gap is integration of real rails + security hardening + frontend wiring from mock to real API.

**Next step I propose:** I start fixing P0 now — ledger book resolver, SSRF, auth bcrypt, idempotency Tx, telebirr & cbe_birr real connectors, Levenshtein, and wiring merchant-web payments page from mockPayments hardcoded to real SWR fetch + RAG proxy + admin health.

Your repo is ready to become outstanding Ethiopian fintech unicorn architecture.

*ApexPay — AI-native payment gateway for Ethiopia • ለኢትዮጵያ outstanding modern UI Mercury/Linear + glassmorphic*

End of Deep Dive.
