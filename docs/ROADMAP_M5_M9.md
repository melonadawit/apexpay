# ApexPay — Phased Roadmap Milestones M5-M9 (Full Platform v1.1.0)

> Senior Engineering Manager view - optimal code structure, excellent algorithms, quality data structures, best practices throughout.

---

## Overview Timeline

| Milestone | Name | Focus | Duration (2 senior eng) | Cumulative |
|---|---|---|---|---|
| **M1** | Vertical Slice | Ledger M1 + payments pending/succeeded + mock | 1.5w | 1.5w |
| **M2** | Merchant Demo + Onboarding | Fayda front/back + 6-step wizard outstanding + links + checkout | 3w | 4.5w |
| **M3** | AI + Exam | Agent rules + admin exam + evidence pack | 1w | 5.5w |
| M4 | MVP Old Freeze | Old MVP tag | - | 5.5w |
| **M5** | **Commercial Gateway Core** | Refunds FULL + 3 rails + smart routing | 3w | 8.5w |
| **M6** | **Recurring + Disbursement** | Subscriptions + Payouts bulk | 3w | 11.5w |
| **M7** | **Workforce** | Payroll ET tax/pension 7/11 OT | 2.5w | 14w |
| **M8** | **Intelligence** | RAG citations + Swarm + Recon | 4w | 18w |
| **M9** | **Mobile + Gold Polish** | Flutter + outstanding UI/UX gold + security hardening | 4w | **22-24w** |

1 eng: ~36-40w.

---

## Milestone M5 — Commercial Gateway Core (WP13-WP15)

**Goal:** Production-grade collection with refund policy and intelligent rail selection.

**Deliverables:**

### WP13 Refunds FULL
- Domain `refund/domain.go` + service with fee reversal policies (non_refundable, pro_rata, full) O(1) idempotency by (merchant_id, refund_ref) unique
- Ledger M2 posting: Dr merchant_payable (R-FR) + Dr platform_fee_due FR Cr clearing:connector R ; filter zero entries optimization
- API `POST /v1/refunds` idempotent, `GET /v1/refunds/{id}`, `GET /v1/payments/{id}/refunds`
- UI: payments detail -> refund button bottom sheet outstanding, amount slider, reason select, maker-checker badge >50k threshold
- Tests: refund_exceeded, duplicate_ref, partial twice, ledger balanced property
- Metrics: `refunds_total{status,type}`

### WP14 Connectors (Telebirr, CBE Birr, Bank IPS, EthSwitch QR, Card slot)
- Interface `connector.Connector` (Initialize, Verify, Refund, Health) Strategy pattern
- Implementations:
  - `mock` 50ms latency always up
  - `telebirr_sandbox`: HMAC SHA256 `X-APP-Key`, callback IP allowlist, latency inject 150-300ms, sandbox URL from config `telebirr_sandbox.sandbox.qa...`, retry with exponential backoff
  - `cbe_birr_sandbox`: similar 30% failure inject to test routing
  - `bank_ips`: ISO20022 mock pain.001 generator
  - `ethswitch`: QR code generator/validator per Ethiopian Interoperable QR standard spec
  - `card_acquirer`: token-only, no PAN log
- Table `connector_configs` encrypted secrets via AES-GCM key `CONNECTOR_ENCRYPTION_KEY` - optimal crypto package, never plain
- Health sampler worker: goroutine ticker 30s, calls Health() per connector, inserts `connector_health_samples`, updates Redis cache `health:{connector}` TTL 60s for fast routing O(1) lookup
- Circuit breaker: in-memory map + Redis backup, failures >=5 in 1m window => open 60s, metric `connector_circuit_open`

### WP15 Smart Routing Engine
- Data structures: `RoutingRule` sorted by priority O(n log n) once per request cached, `ConnectorHealth` aggregated 5m window success_rate, avg_latency
- Algorithm:
  ```
  candidates = filter(rules, amount,currency,method,enabled) -> sort priority asc O(n log n)
  selected = candidates[0]
  healthMap = getHealth(primary+fallbacks) O(k) k<=3
  if circuit open -> fallback1
  if strategy==success_rate and primary.success_rate <0.7 -> pick max success_rate fallback
  if strategy==latency -> pick min latency
  ```
- APIs: `GET /v1/methods?amount=1000&currency=ETB` ranked array with score 0.6*success+0.4*(1-latency/1000), `GET /admin/connectors/health` Recharts
- UI outstanding: checkout shows badge "Using best route: Telebirr (2% faster today)" with tooltip explaining fallback trail; admin dashboard charts latency line, success bar, circuit state chips green/yellow/red animated
- Tests: table-driven routing decision 100% coverage, circuit open integration

**Exit Gate M5 Checklist (see MVP §9.2)**
- Refund partial + full + fee reversal
- 3 connectors latency inject working
- Health samples every 30s populated
- `GET /v1/methods` ranked
- Failover demo: kill telebirr -> mock fallback audit `routing.fallback_used:true`

---

## Milestone M6 — Recurring + Disbursement (WP16-WP17)

**Goal:** Merchant monetizes subscriptions and disburses vendors.

### WP16 Subscriptions
- Tables: customers (fayda_fin_hash optional), subscription_plans, subscriptions FSM (incomplete,trialing,active,past_due,canceled,paused), subscription_invoices
- Services:
  - `subscription.CreatePlan` validates amount>0, interval_count default 1
  - `CreateSubscription` trial handling: if trialDays>0 status trialing, period = trialEnd else interval add (day/week/month/year) optimal `addInterval` using time.AddDate
  - Invoice draft/open/paid, dunning: `NextDunningAttempt` 1d,3d,5d after failure exponential-ish
- Workers: dunning worker Cron hourly scans invoices due + past_due, attempts payment via saved method mock, updates attempt_count, webhook `subscription.*`
- APIs: `POST /v1/subscription_plans`, `/customers`, `/subscriptions`, `/subscriptions/{id}/cancel`, customer portal magic link 24h token via signed JWT
- UI outstanding: subscriptions pages tabs Plans/Subscriptions/Customers/Invoices with cards glassmorphic, plan creation wizard, invoice overdue red badge, portal hosted pages with payment button modern
- Ledger: invoice success posts M1 linked subscription_id

### WP17 Payouts
- Tables: beneficiaries (bank verification name fuzzy match Levenshtein <3), payout_batches (book per batch), payouts FSM (created->pending_approval->queued->processing->succeeded|failed|returned)
- Service:
  - `CreateSingle` checks merchant balance decimal precise, ApprovalThreshold 50k ETB -> pending_approval else queued ; ledger M3 Dr payable Cr clearing_bank immediate if queued else on approve
  - `CreateBulk` validates bulk 1-1000 items O(n) sum, total <= balance, creates batch pending_approval (all bulk require approval per policy), payouts created status created, journal Dr payable total Cr clearing total per batch book
  - CSV parser for bulk: papaparse frontend + backend csv lib, preview table outstanding with validation icons green/red, amount sum, fees calc (MDR)
  - Payout links: escrow book creation, claim via OTP, move escrow->clearing on claim
  - Maker-checker: dual approval >50k, check approver != submitter
- APIs: `POST /v1/beneficiaries`, `POST /v1/payouts`, `POST /v1/payouts/bulk`, `POST /v1/payout_batches/{id}/approve`
- UI outstanding: payouts page bulk dropzone dashed pulse on drag, CSV row validation checklist animated, batch timeline like GitHub Actions steps (queued->processing->succeeded), beneficiary combobox with bank logos
- Metrics: `payouts_total{status}`, `payout_batch_amount`

**Exit Gate M6**
- Subscription trial -> past_due dunning 3 attempts visible
- Beneficiary name fuzzy match validation working
- Bulk CSV 5 rows upload preview OK -> batch -> approve dual -> all succeeded
- Ledger M3 balanced per batch
- Payout link claim flow with escrow

---

## Milestone M7 — Workforce Payroll ET (WP18)

**Goal:** Ethiopian-compliant payroll.

### WP18 Payroll OS
- Tables: employees (fayda_fin_hash, bank masked/hash, base_salary, employment_date, cost_center), payroll_runs (book_id per run, period_month/year, type regular/off_cycle/bonus/adjustment, status FSM draft->calculating->pending_approval->approved->processing->completed->failed), payroll_items (gross, OT hours/amount, commission, taxable, income_tax, pension 7%/11%, net), payroll_claims
- Core Algorithm - optimal decimal precise, no float:
  - `CalculateTax` binary search O(log n) over sorted brackets table `payroll_tax_brackets` DB-driven effective_date versioned - placeholder 2024 brackets: 0-600 0% 0, 601-1650 10% -60, 1651-3200 15% -140, 3201-5250 20% -300, 5251-7800 25% -565, 7801-10900 30% -955, >10900 35% -1500 formula tax = taxable*rate - deduction rounded 2 decimals
  - Pension employee 7% gross, employer 11% gross, capped per rule (placeholder no cap)
  - OT: hourly_rate = base_salary / 208 (26 days *8h) ET standard, rates: weekday 1.25x, weekend 1.5x, holiday 2x, night 1.3x per labour law
  - `CalculateRun` loops employees O(n) active only, computes gross = base+OT+commission, pensionEmp, taxable = gross - pensionEmp (non taxable portion) - other exemptions, incomeTax via binary search, deductions = pensionEmp+incomeTax+other, net = gross - deductions, aggregates totals
- State machine: draft -> calculating bulk insert items Tx -> pending_approval -> approved (dual approval >100k net) -> processing -> completed
- Ledger M4: per run book creation `ledger_books` book_type payroll_run, journal posting_key `payroll_run:{id}`: Dr expense:salary totalGross Cr liability:payroll_payable totalNet Cr liability:et_income_tax_payable totalTax Cr liability:pension_payable totalPension ; ValidateBalanced before posting O(k) k=entries
- Disburse: creates payout batch for employees banks, second journal Dr payroll_payable Cr asset:clearing:bank totalNet via payouts
- APIs: `POST /v1/employees` (+ import CSV), `POST /v1/payroll_runs`, `POST /v1/payroll_runs/{id}/calculate`, `/{id}/approve`, `/{id}/disburse`, `GET /v1/payroll_runs/{id}/payslips/{employee_id}/pdf` generates PDF modern template with logo, QR verification, breakdown table, pie deductions
- ET reports: CSV ERCA with TIN, taxable, tax ; JSON for export
- UI outstanding: employees table with Fayda badge verified, avatars, cost_center chips ; runs table status pipeline visual stepper ; run detail table sticky footer totals, row expand breakdown chart, approval flow dual avatar, payslip preview drawer glassmorphic, download all zip
- Tests: tax bracket property vs known examples, rounding edge .005, payroll balanced invariant, 500 employees calc p99 <2s benchmark

**Exit Gate M7**
- 10 employees CRUD
- Run July regular calculate -> approval dual -> disburse -> payout batch created
- Payslip PDF download modern outstanding
- Tax matches expected per bracket
- Ledger M4 balanced per run book

---

## Milestone M8 — Intelligence (WP19-WP21)

**Goal:** NBE answers with citations, AI crews operate money, recon matches bank.

### WP19 RAG Compliance (Python rag-worker now required)
- Stack: PostgreSQL + pgvector extension vector(1024), optional Qdrant at scale 1M+, Python 3.12 rag-worker
- Pipeline: `rag_documents` status pending -> download PDF -> extract via PyMuPDF -> clean -> chunk 800 tokens overlap 100 for optimal retrieval -> embed via multilingual-e5-large or bge-m3 1024 dim (embed batch 32 for efficiency) -> upsert pgvector ivfflat lists=100 index + metadata title, source_type, page
- Services:
  - `Ask` query -> embed -> vector.Search topK 5 cosine O(log n) via ivfflat -> threshold 0.65 if top score < threshold return no answer guard (prevent hallucination) -> build prompt with context [1]..[n] + question + lang -> LLM generate (OpenAI compatible or local) -> return answer + citations
  - `IngestDocument` chunking function simple char-based skeleton real tiktoken
- Evals: precision@3 measured via eval script 0.8 threshold
- APIs: `POST /v1/compliance/ask` {q, lang en/am, topK} -> {answer, citations[{docTitle, page, chunkId, score, sourceURL}]} streaming SSE outstanding, `GET /v1/admin/rag/docs`, `POST /v1/admin/rag/ingest`
- UI outstanding: compliance center chat like Perplexity.ai - input + AM/EN toggle, history sidebar, citations as rounded badges with hover PDF preview, answer streaming markdown, source list with icons NBE directive vs policy, empty state "Ask about NBE refund timeframe, 2FA limit, AML reporting"
- Guard: no answer without citation - prompt engineering rule "If answer not in context, say Not in compliance corpus"
- Sample ingests: ONPS/10/2025 excerpt about 2FA >5000 ETB, refund policy doc, AML directive ETB 200k reporting

### WP20 Swarm Multi-Agent
- Registry: `swarm/domain.go` tool definitions map O(1) lookup: create_payment_link threshold 100k, create_payout 50k, calculate_payroll_draft 100k, refund, get_tpv, ask_compliance, list_payments - role allowed owner/admin/developer/finance
- Planner: `RulesPlanner` MVP keyword matching deterministic (if "link" -> create_payment_link, if "payroll" -> calculate_payroll, if "payout" -> create_payout, if "tpv" -> get_tpv) ; prod LLM planner via OpenAI function calling
- Critic: checks policy/hallucination/amount thresholds, sets confirmation_required boolean if total amount > threshold or destructive, builds confirmation_data
- Executor: `ToolExecutor` interface calls Go services (createPaymentLink via payment service, calculatePayrollDraft via payroll service etc) with JSON schema validation via validator
- Session state machine: planning -> executing -> needs_confirmation (outstanding modal UI) -> completed / failed / cancelled ; each step status pending/executing/succeeded/failed
- Audit: `agent_runs` + `swarm_sessions` tables, every tool call logged with latency_ms, args, result, request_id correlation
- APIs: `POST /v1/swarm/run` {goal} -> session, `POST /v1/swarm/{id}/confirm` {confirmed bool}, `GET /v1/swarm/{id}` + `GET /v1/agent/runs`
- UI outstanding: merchant command center chat right panel glassmorphic, swarm run shows stepper timeline with tool call cards like Vercel AI SDK - each card icon + description + args preview + status spinner check + result link, confirmation modal outstanding with breakdown + biometric option, final output summary with links
- Safety: no direct ledger_entries insert from LLM - must via domain service; ledger invent blocked by critic
- Metrics: `swarm_run_total{status}`, `swarm_tool_call_duration_seconds{tool}`

### WP21 Recon
- Tables: recon_statements (connector_id, statement_date, raw_file_ref MinIO, raw_file_hash, parsed_json, total_amount, total_count, status pending/parsed/matched/has_breaks), recon_breaks (statement_id, ledger_book_id, reference_type, reference_id, expected, actual, difference, status open/investigating/resolved/written_off)
- Parser: CSV/MT940/JSON - detects headers amount, connector_ref, date ; stores parsed_json array
- Matching engine: for each statement row, find ledger journal by (connector_id, connector_ref) + amount tolerance 0.01 ETB + date window 24h ; if mismatch or not found -> break open O(n*m) optimized via map connector_ref->journal O(n+m)
- Auto-resolve rules: if diff <1 ETB fee rounding -> resolve; else open
- Suspense posting: break difference -> suspense book journal until resolved
- Worker cron daily 02:00 Africa/Addis_Ababa
- APIs: `POST /v1/admin/recon/statements/upload`, `GET /v1/admin/recon/breaks?status=open`, `POST /v1/admin/recon/breaks/{id}/resolve`
- UI outstanding: recon dashboard chart matched vs breaks over time, breaks table with difference red, assign dropdown, resolve button with adjustment reason, statement upload dropzone preview

**Exit Gate M8**
- RAG ask "What is 2FA threshold?" -> answer "5000 ETB per ONPS/10/2025" + citation page 3 score 0.92 clickable PDF viewer
- RAG irrelevant question -> "Not in compliance corpus"
- Swarm "Create link 100 ETB for coffee and run payroll July" -> planner shows 2 steps, needs_confirmation modal for payroll 150k net outstanding, confirm -> executes link + payroll draft, audited
- Recon upload statement CSV -> 0 breaks initially, inject mismatch -> 1 break open -> resolve writes suspense adjustment

---

## Milestone M9 — Mobile + Gold Polish (WP22-WP23)

### WP22 Flutter Merchant App 1.1.0
- Setup: Flutter 3.22 + Riverpod + go_router + dio + hive + flutter_secure_storage + camera + mobile_scanner + share_plus + firebase_core/messaging + local_auth + lottie + shimmer + cached_network_image
- Layers: `lib/src/core/{api, router, storage/theme}` + `features/{auth, dashboard, links, qr, approvals, profile/onboarding}` + `features/profile/onboarding/presentation/onboarding_wizard_page.dart` outstanding 6-step wizard matching web
- Features breakdown:
  - Auth: email+password+OTP 2FA UI outstanding glass card, token secure_storage, refresh interceptor, biometric login optional
  - Dashboard: glassmorphic TPV card gradient emerald + sparkline Recharts-like via fl_chart, recent payments list shimmer loading skeleton, pull-to-refresh, empty state illustration Ethiopian coffee ceremony, quick actions FAB
  - Create Link: bottom sheet draggable half/expanded, amount input ETB formatting with `_,___`, chips 100/500/1000, AI suggest description via API, QR preview live using qr_flutter, share button system share (Telegram, WhatsApp) via share_plus, copy link haptics vibration
  - Scan QR: camera permission request outstanding dialog with illustration, scanner overlay rounded square 260dp guides + corner brackets animated pulse green, supports FaydaEncode offline QR (decode FIN last4 + name) + EthSwitch interoperable QR (parse merchant + amount + verify via GET /v1/payments/qr/{code}), vibration on scan, result card glassmorphic
  - Approvals inbox: list pending payouts/payroll runs amount badge warning, swipe right approve left reject (flutter_slidable), approve triggers biometric prompt local_auth then API call, success confetti Lottie 3s, failure snackbar
  - Onboarding Wizard Mobile: 6 steps PageView dot indicator, step1 business info smart form, step2 owners & Fayda capture modals with camera overlay corner guides glare detection via canvas brightness check (if glare > threshold show "Move to shade"), front/back/selfie capture + OTP 6-digit pin input animated, FIN/FAN input mask, consent checkbox AM/EN, step3 bank combobox with bank logos from GET /v1/banks, step4 docs vault dropzone with thumbnails PDF/JPG preview, required checklist donut progress, step5 compliance gauge chart, step6 review summary cards + terms
  - Push: FCM token registration `POST /v1/devices/register` + `push_devices` table, local notification mock on payment succeeded (foreground), background handler
  - Offline: Hive box `draft_links` stores offline create link drafts when no internet, sync badge count appBar, sync on reconnect using idempotency key same as web
  - Theming: dark/light ThemeData same tokens as web: primary ET Green #0B6E4F, accent gold #EAB308, radius xl 24, Inter + NotoSansEthiopic, 8pt grid, glassmorphic nav, empty states illustrations
- Build: `flutter build apk` + iOS TestFlight internal

### WP23 Outstanding UI/UX Polish + Security Hardening Gold
- Design tokens finalization: `libs/ui/tokens.json` + Tailwind config + Flutter ThemeData same palette - extract into shared `design-tokens` package
  ```json
  {
    "colors": {"primary": {"default": "#0B6E4F", "light": "#10A37A"}, "accent": {"gold": "#EAB308"}},
    "radius": {"lg": "16px", "xl": "24px", "2xl": "32px"},
    "shadows": {"soft": "0 10px 30px rgba(0,0,0,0.07)"},
    "motion": {"ease": [0.22,1,0.36,1], "duration": {"fast": "200ms", "medium": "300ms"}},
    "typography": {"fontSans": "Inter", "fontEthiopic": "Noto Sans Ethiopic"}
  }
  ```
- Motion: Framer Motion variants fade + slide + scale, stepper progress animated line pathLength using motion.div, onboarding wizard steps AnimatePresence, Fayda capture overlay corner brackets pulse scale 1->1.1 infinite, file upload dropzone pulse border on drag, checkout success confetti canvas-confetti full-screen 3s + haptic, skeleton shimmer for lists, stagger list reveals delay 50ms * index
- Components outstanding:
  - Merchant web shell: glassmorphic nav `backdrop-blur-xl bg-white/70` border bottom, sidebar collapsible with icons + tooltips, cards hover elevated shadow soft -> medium transition 200ms
  - Dropzone: dashed border 2px, bg zinc-50 on idle, bg primary/5 on drag, icon upload cloud with bounce, file previews thumbs PDF icon + JPG thumb 64px, upload progress slim top bar, status chip check animate
  - Fayda camera: overlay guides rounded 16, corners brackets 30px L shape greenAccent, glare detection via canvas getImageData brightness average >200 => show warning outstanding helper "Move to shade • ጥላ ውስጥ ይሂዱ" with icon, auto-capture when corners detected (edge detection placeholder OpenCV.js light)
  - Document viewer: split pane left list right preview, PDF thumbs via react-pdf, image zoom pan, status badge verified yellow dot
  - Timeline: vertical line like Linear Timeline, dot for each status done/current/upcoming, card for each step content, timestamp
- Checkout outstanding:
  - Mobile-first max 420px centered, merchant logo + name small, amount large 32px bold, method selector radio cards with icons Telebirr, CBE, bank, card, QR interoperable - selected border primary 2px + subtle bg primary/5 + check icon, trust badges "NBE Licensed PSO Gateway" + secure lock
  - Processing lottie Ethiopia pattern looping, polling verify every 2s
  - Success: full-screen check animation Lottie scale + confetti + receipt actions row: download PDF (jsPDF invoice outstanding template with QR), email receipt, back to merchant link
- Accessibility: axe DevTools audit 0 serious violations on merchant-web + checkout-web + admin-web, keyboard nav logical tab order, focus rings visible primary outline 2px offset, color contrast AA 4.5:1 checked via tool, screen reader labels aria-label for file inputs, AM/EN lang attributes
- Security hardening:
  - Rate limit: Fayda OTP 5/hour per IP + per owner via Redis token bucket `fayda:otp:{owner}` TTL 1h
  - Document upload: presigned POST via MinIO 15m TTL, file type whitelist MIME pdf/jpg/png, size <5MB Fayda <2MB per NIDP, ClamAV stub interface `VirusScanner` returns clean, hash check sha256 integrity, encrypted SSE-S3 MinIO
  - Secrets: env vault, `gosec` high issues 0, `govulncheck`, dependency scan via `nancy` or `trivy`
  - No plain FIN/logs: grep logs test in CI ensuring no 12-digit pattern, PII redact middleware `zerolog` field filter, FIN only last4 in responses
  - 2FA >5000 ETB mandatory enforced in payment service, test in CI
  - Maker-checker: payout >50k, payroll >100k net, approval count check, approver != submitter enforced
- Performance:
  - Lighthouse merchant: Performance 90, Accessibility 100, Best Practices 100, SEO 90 ; checkout: Perf 95 etc ; first contentful paint <1.8s via Next.js optimization images, dynamic import Lottie
  - Flutter startup <2s profiled via `flutter run --profile`, deferred components
  - k6: 100 VUs 5m no errors p95 <300ms, payroll calc benchmark `go test -bench=. -benchtime=3s` <2s for 500 emps
- Docs: Updated README with onboarding sequence diagram mermaid, Fayda flow diagram, routing diagram, outstanding UI screenshots placeholders `docs/screenshots/onboarding-wizard.png`, demo video script `scripts/demo.sh` 15min gold

**Exit Gate M9 Full Gold**

Full demo script (see MVP §7 script 1-35) passes in ≤15 min by new engineer:
- Clone, docker compose up postgres+redis+minio+pgvector, migrate up 0001..0012 clean
- Seed banks list + platform books + compliance user + mock Fayda partner
- Register merchant email+phone OTP -> onboarding wizard outstanding 6 steps desktop mobile responsive: capture front/back Fayda via mock images, OTP 123456 mock, upload docs drag-drop, submit compliance auto checks pass
- Switch to admin web Kanban board drag merchant card submitted -> pending_approval dual approval second user -> active + operating book + test keys
- Back to merchant: dashboard empty state illustration "Welcome Meron, let's collect first ETB" + CTA outstanding + AI chat panel glassmorphic
- Create link 500 ETB tutoring: QR displayed, share Telegram button web share fallback, copy haptics
- Open checkout mobile: method selector outstanding, pay mock, success confetti PDF receipt
- Refund 100 partial outstanding flow timeline
- Subscription plan monthly 500 trial 7d + subscription + customer portal magic link
- Payout bulk CSV 5 rows preview outstanding -> approve via Flutter app swipe biometric -> succeeded exam ledger M3
- Payroll run July 10 emp OT + bonus -> calculate -> approve dual -> disburse -> payout batch -> download payslip PDF outstanding modern template QR verification
- RAG ask "What is 2FA limit?" -> answer 5000 ETB citation badges clickable opening PDF viewer highlight
- Swarm chat "Create link 100 for coffee if today TPV>0" -> planner shows steps stepper, confirmation modal outstanding if >100k, executes
- Admin exam search tx_ref -> shows lifecycle + ledger + routing trail fallback_used + webhook deliveries attempts + Fayda verified badge link to KYC profile versions + docs viewer blurred front/back until click auth + audit logs chain
- Evidence pack export JSON download hash verification
- k6 100 VUs 5m no errors
- Lighthouse 90+ + axe 0 serious + gosec 0 high
- Tag `v1.0.0-full` — ready for pilot with NBE.

---

## Risks Mitigation M5-M9 (Expanded)

See MVP §10 but add:
- Fayda API approval delay: use mock interface behind Verifier strategy, mock OTP 123456 always works, face_score mock 0.92, pluggable live via ENV toggle FAYDA_MODE=live
- ET tax brackets change: DB table payroll_tax_brackets versioned effective_from, config UI admin
- Photo fraud front/back: glare detection + hash + manual review + face match 0.85 threshold + liveness via selfie angle check, compliance reviewer sees originals with blur by default
- Scope big 22-24w: incremental deploys after each M, feature flags routing_rules.enabled, connector_configs.enabled, gifts demo
- Flutter store: PWA fallback checkout-web installable, internal TestFlight first

---

## Metrics M5-M9 (see MVP §11 full)

Add:
- onboarding_*, fayda_*, document_*, refunds, payouts bulk, payroll calc duration, rag query latency, swarm runs, connector health, routing fallback
- Prometheus + Grafana dashboards per milestone outstanding

---

## Next: Post Full (Network)

- KYC Level3 physical verification + geo agent assignment
- White-label + SDK npm `apexpay-js` + embedded checkout
- Capital lending based TPV + ledger receivables factoring algorithm credit scoring
- International corridors UPI
- Payroll anomaly detection ML

*End of ROADMAP M5-M9*
