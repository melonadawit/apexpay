# ApexPay - FULL Platform Spec (v1 Commercial)
# Includes all previously EXCLUDED features

| Field | Value |
|---|---|
| **Document ID** | APX-FULL-001 |
| **Version** | 1.1.0-expanded |
| **Date** | 2026-08-03 |
| **Based on** | APX-MVP-001 + DATABASE.md + SAD.md |
| **Change** | ❌ subs, payouts, payroll, RAG, swarm, smart routing, Flutter, refunds full, sandbox rails -> ✅ INCLUDED |

---

## 0. New Goal Statement

**Original MVP:** Link -> pay via mock -> ledger -> webhook -> AI
**Full v1:** A merchant can collect (one-off, links, subscriptions), store (ledger books), refund, disburse to vendors/staff (payouts + full payroll with ET statutory packs), get AI swarm to operate it, get NBE compliance answers from RAG, with smart rail routing, real Telebirr/CBE/Bank connectors, and Flutter merchant app - all exam-ready.

---

## 1. Revised Functional Scope Matrix

| Capability | Before | Now | Details |
|---|---|---|---|
| Merchant auth, keys, webhooks, ledger, links, checkout, mock | ✅ | ✅ | Unchanged |
| **Refunds** | ⚪ stub | **✅ FULL** | Full lifecycle, partial/full, fee reversal policy |
| **Sandbox + Real Rails** | ⚪ optional | **✅ 3 rails** | mock + telebirr sandbox + CBE Birr + bank slot + card slot |
| **Subscriptions** | ❌ | **✅** | Plans, subscriptions, trials, dunning, customer portal, proration |
| **Payouts** | ❌ | **✅** | Single/bulk, payout links, vendor pay, bank list, approval workflow |
| **Payroll / Workforce** | ❌ | **✅** | HR lite -> full payroll OS: employees, salary structures, OT, claims, commissions, ET income tax, pension, cost centers |
| **RAG / Compliance** | ❌ | **✅** | NBE corpus ingest, embeddings, citations, Amharic/English, eval harness |
| **Swarm Multi-Agent** | ❌ | **✅** | Planner, Critic, Executor, Collect/Recover/Analyze/Advise/PayrollAssist crews |
| **Smart Routing** | ❌ | **✅** | Connector health, latency, success rate, cost, failover, A/B |
| **Flutter App** | ❌ | **✅** | Merchant app: dashboard, create link, scan QR, approve payout, push notifications |

---

## 2. Feature Expansion Detail

### 2.1 Refunds - FULL (was stub)

**User stories:**
- As Meron, I can refund a payment fully or partially from dashboard/API
- As Yonatan, refund via API triggers webhook `refund.succeeded`
- As Finance, I see fee handling rule (non-refundable vs pro-rata)

**FSM:** `created -> processing -> succeeded | failed`
Unique constraint: (merchant_id, payment_id, refund_ref)

**Ledger Model M2 - Refund:**
Amount R = refund amount, Fee reversal FR (policy)
```
Journal posting_key = refund:{refund_id}
Dr liability:merchant_payable   R-FR (merchant balance decreases)
Dr liability:platform_fee_due   FR
  Cr asset:clearing:mock        R
```
Must be in same transaction as payment status -> partially_refunded / refunded

**APIs:**
- `POST /v1/refunds`
- `GET /v1/refunds/{id}`
- `GET /v1/payments/{id}/refunds`

**UI:** Payments detail -> Refund button with amount, reason, maker-checker if > threshold

### 2.2 Subscriptions & Recurring

**Concepts:**
- `subscription_plans`: interval (day/week/month/year), trial_days, amount, currency, ETB-first
- `subscriptions`: status (incomplete,trialing,active,past_due,canceled), current_period_start/end
- `subscription_invoices` + `invoice_items`: billing attempts
- `customers`: for recurring

**Workers:** dunning worker (retry schedule: 1d,3d,5d), invoicing cron, webhook `subscription.*`

**Ledger:** subscription invoice success posts same as M1, linked to subscription_id

**APIs:**
- `POST /v1/subscription_plans`
- `POST /v1/subscriptions`
- `POST /v1/subscriptions/{id}/cancel|pause|resume`
- `GET /v1/customers` + customer portal magic link

**UI:** Subscription plans page, subscription list, customer portal (hosted)

### 2.3 Payouts / Disbursement

**Capabilities:** single payout, bulk CSV upload, payout links (recipient claims), bank account validation ET bank list, approval rules.

**FSM:** `created -> pending_approval -> queued -> processing -> succeeded | failed | returned`
Maker-checker: > X ETB requires admin role.

**Ledger Model M3 - Payout:**
```
Journal posting_key = payout:{payout_id}
Dr liability:merchant_payable   A
  Cr asset:clearing:bank        A
And fee if applicable:
Dr liability:merchant_payable   fee
  Cr liability:platform_fee_due fee
```
Escrow book for payout links until claimed.

**Tables:** `payouts`, `payout_batches`, `beneficiaries`

**APIs:** `POST /v1/payouts`, `POST /v1/payouts/bulk`, `POST /v1/payout_batches/{id}/approve`

### 2.4 Payroll / Workforce Money OS

This is the biggest expansion - now included.

**Domains:**
- `employees`: name, tin, pension id, bank, salary, status
- `employments`: salary structure, grade, cost_center
- `payroll_runs`: month, type (regular, off-cycle, bonus), status (draft, pending_approval, processing, completed)
- `payroll_items`: gross, deductions, OT, taxable income, income tax ET, pension employee 7% + employer 11%, net pay
- `payroll_claims`, `payroll_commissions`, `payroll_ot`

**Ledger Model M4 - Payroll Run:**
Per your docs/PAYROLL-BUSINESS-POWER vision: payroll_run becomes a book.
```
For each employee:
Dr expense:salary         gross
  Cr liability:payroll_payable   net
  Cr liability:et_income_tax     tax
  Cr liability:pension_payable   pension total
Then disburse payroll_payable via Payout batch (M3)
```

**ET Stat Packs:** Formula engine for ET income tax brackets (configurable), pension calc, cost allocation.

**AI Payroll Assist:** "Run payroll for July, add 10k bonus for Sales" -> swarm creates draft run.

**APIs:** `POST /v1/employees`, `/v1/payroll_runs`, `/v1/payroll_runs/{id}/calculate|approve|disburse`

**UI:** Workforce > Employees, Runs, Approvals

### 2.5 RAG Compliance Layer (was out)

**Stack:** pgvector for MVP, Qdrant at scale, Python 3.12 rag-worker for ingest.

**Corpus:** NBE directives (Payment System Operator, AML/CFT), ApexPay policies, FAQ.

**Pipeline:** docs -> chunk (800 tokens overlap 100) -> embed (bge-m3 or multilingual-e5) -> pgvector -> query -> LLM with citations required.

**Tables:** `rag_documents` (id, source_type, title, url, hash, status), `rag_chunks` (document_id, content, embedding vector(1024), metadata)

**APIs:**
- `POST /v1/compliance/ask` -> answer + citations[] + source pages
- `POST /v1/admin/rag/ingest`

**Quality:** No answer without citation, hallucination eval harness, Amharic + English queries.

**UI:** Compliance center in admin-web + merchant chat: "What does NBE say about refund timelines?"

### 2.6 Swarm Multi-Agent System (was out)

**Roles:** Planner -> decomposes intent, Critic -> checks policy/ledger invariants, Executors (Collect, Recovery, Analytics, Payroll Crew)

**Architecture:** Go swarm executor (worker) orchestrator + Python optional LLM tool caller. Tools are Go API functions, not free-form SQL.

**Tool allowlist:** create_payment_link, list_payments, create_payout, create_refund, calculate_payroll (draft only), ask_compliance

**Safety:**
- Ledger tools require confirmation if amount > threshold
- No direct ledger_entries insert from LLM - must go via domain service
- All tool_calls logged to agent_runs (existing table) + new `swarm_sessions`

**APIs:** `POST /v1/agent/chat` enhanced to `POST /v1/swarm/run` with trace

**UI:** Merchant command center chat now shows plan/steps, not just response.

### 2.7 Smart Routing (was out)

**Problem:** Ethiopia has fragmented rails - Telebirr may be down, CBE slow, bank IPS cheaper.

**Solution:** Routing engine in connector mesh.

**Signals:** `connector_health_samples` (latency, success_rate, uptime last 5m), cost per rail, merchant preference, amount tier.

**Routing rules:** DB table `routing_rules` (merchant_id nullable, method, amount_range, primary, fallback1,2, strategy=latency|cost|success_rate)

**FSM:** initialize tries primary; on timeout/circuit open -> fallback with same payment_id idempotent.

**APIs:** `GET /v1/methods?amount=1000&currency=ETB` returns ranked, `GET /v1/admin/connectors/health`

### 2.8 Flutter Mobile App (was out)

**Scope now:**
- Auth (magic link + password)
- Dashboard: TPV today, recent payments
- Create payment link + share to Telegram/WhatsApp
- Scan interoperable QR to verify payment
- Approve payout/payroll if ops role
- Push notifications via FCM: payment succeeded, payout needs approval
- Offline queue: create draft link offline, sync later

**Tech:** Flutter + go_router + Riverpod + Hive local cache

### 2.9 Real Connectors - Telebirr, CBE, Bank, EthSwitch, Card

**Interface:** (already in SAD) Initialize, Verify, Refund, HealthCheck

**Must implement:**
- `mock` (existing) + simulation of timelines
- `telebirr_sandbox`: HMAC signing, callback IP allowlist
- `cbe_birr_sandbox`: similar
- `bank_ips`: ISO20022 mock + recon file (MT940/csv)
- `ethswitch` slot: QR interoperable
- `card_acquirer` stub: no PAN storage, token only

**Config:** `connectors` table or env config with merchant_id override.

---

## 3. Database Expansion - DDL

Add to DATABASE.md §5:

```sql
-- REFUNDS
create table refunds (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  payment_id      text not null references payments(id),
  refund_ref      text not null,
  amount          numeric(20,8) not null check (amount > 0),
  currency        char(3) not null,
  status          text not null check (status in ('created','processing','succeeded','failed')),
  reason          text,
  fee_reversal    numeric(20,8) not null default 0,
  connector_id    text not null,
  connector_ref   text,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now(),
  unique (merchant_id, refund_ref)
);
create index refunds_payment_idx on refunds (payment_id);

-- SUBSCRIPTIONS
create table customers (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  email           text, phone text, name text,
  metadata        jsonb not null default '{}',
  created_at      timestamptz not null default now()
);
create table subscription_plans (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  name            text not null,
  amount          numeric(20,8) not null,
  currency        char(3) not null default 'ETB',
  interval_type   text not null check (interval_type in ('day','week','month','year')),
  interval_count  int not null default 1,
  trial_days      int not null default 0,
  status          text not null default 'active',
  created_at      timestamptz not null default now()
);
create table subscriptions (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  customer_id     text not null references customers(id),
  plan_id         text not null references subscription_plans(id),
  status          text not null check (status in ('incomplete','trialing','active','past_due','canceled','paused')),
  current_period_start timestamptz not null,
  current_period_end   timestamptz not null,
  trial_end       timestamptz,
  cancel_at       timestamptz,
  created_at      timestamptz not null default now()
);
create table subscription_invoices (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  subscription_id text not null references subscriptions(id),
  payment_id      text references payments(id),
  amount          numeric(20,8) not null,
  status          text not null,
  attempt_count   int not null default 0,
  due_at          timestamptz not null,
  created_at      timestamptz not null default now()
);

-- PAYOUTS
create table beneficiaries (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  name            text not null,
  account_no      text not null,
  bank_code       text not null,
  type            text not null check (type in ('individual','business')),
  created_at      timestamptz not null default now()
);
create table payout_batches (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  amount          numeric(20,8) not null,
  currency        char(3) not null,
  status          text not null,
  approved_by     text,
  created_at      timestamptz not null default now()
);
create table payouts (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  batch_id        text references payout_batches(id),
  beneficiary_id  text references beneficiaries(id),
  payout_ref      text not null,
  amount          numeric(20,8) not null,
  currency        char(3) not null,
  status          text not null check (status in ('created','pending_approval','queued','processing','succeeded','failed','returned')),
  method          text not null,
  connector_id    text,
  failure_code    text,
  created_at      timestamptz not null default now(),
  unique (merchant_id, payout_ref)
);

-- PAYROLL
create table employees (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  employee_code   text not null,
  name            text not null,
  email           text, phone text,
  tin             text,
  pension_no      text,
  bank_account    text,
  base_salary     numeric(20,8) not null,
  status          text not null check (status in ('active','inactive','terminated')),
  metadata        jsonb not null default '{}',
  created_at      timestamptz not null default now(),
  unique (merchant_id, employee_code)
);
create table payroll_runs (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  book_id         text references ledger_books(id), -- each run has its own book
  run_ref         text not null,
  period_month    int not null check (period_month between 1 and 12),
  period_year     int not null,
  type            text not null check (type in ('regular','off_cycle','bonus')),
  status          text not null check (status in ('draft','calculating','pending_approval','approved','processing','completed','failed')),
  total_gross     numeric(20,8) not null default 0,
  total_deductions numeric(20,8) not null default 0,
  total_net       numeric(20,8) not null default 0,
  created_at      timestamptz not null default now(),
  unique (merchant_id, run_ref)
);
create table payroll_items (
  id              text primary key,
  run_id          text not null references payroll_runs(id),
  employee_id     text not null references employees(id),
  gross           numeric(20,8) not null,
  ot_amount       numeric(20,8) not null default 0,
  commission      numeric(20,8) not null default 0,
  taxable_income  numeric(20,8) not null,
  income_tax      numeric(20,8) not null,
  pension_employee numeric(20,8) not null,
  pension_employer numeric(20,8) not null,
  other_deductions numeric(20,8) not null default 0,
  net_pay         numeric(20,8) not null,
  status          text not null
);

-- RAG
create extension if not exists vector;
create table rag_documents (
  id            text primary key,
  title         text not null,
  source_type   text not null check (source_type in ('nbe_directive','policy','faq','evidence')),
  source_url    text,
  content_hash  text not null,
  status        text not null check (status in ('pending','indexed','failed')),
  created_at    timestamptz not null default now()
);
create table rag_chunks (
  id            text primary key,
  document_id   text not null references rag_documents(id) on delete cascade,
  chunk_index   int not null,
  content       text not null,
  embedding     vector(1024),
  metadata      jsonb not null default '{}',
  created_at    timestamptz not null default now()
);
create index rag_chunks_embedding_idx on rag_chunks using ivfflat (embedding vector_cosine_ops);

-- SMART ROUTING + HEALTH
create table connector_configs (
  id              text primary key,
  connector_id    text not null, -- telebirr, cbe, bank, ethswitch, card, mock
  merchant_id     text references merchants(id), -- null = global
  environment     text not null check (environment in ('test','live')),
  config          jsonb not null, -- keys, urls
  enabled         boolean not null default true,
  created_at      timestamptz not null default now(),
  unique (connector_id, merchant_id, environment)
);
create table connector_health_samples (
  id              text primary key,
  connector_id    text not null,
  latency_ms      int not null,
  success         boolean not null,
  sampled_at      timestamptz not null default now()
);
create index health_sample_time_idx on connector_health_samples (connector_id, sampled_at desc);
create table routing_rules (
  id              text primary key,
  merchant_id     text references merchants(id),
  min_amount      numeric(20,8),
  max_amount      numeric(20,8),
  currency        char(3) not null default 'ETB',
  primary_connector text not null,
  fallback1       text,
  fallback2       text,
  strategy        text not null default 'success_rate',
  enabled         boolean not null default true,
  created_at      timestamptz not null default now()
);

-- SWARM + RECON
create table swarm_sessions (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  user_id         text,
  goal            text not null,
  plan            jsonb not null,
  status          text not null,
  created_at      timestamptz not null default now()
);
create table recon_statements (
  id              text primary key,
  connector_id    text not null,
  statement_date  date not null,
  raw_file_ref    text not null,
  parsed_json     jsonb not null,
  created_at      timestamptz not null default now()
);
create table recon_breaks (
  id              text primary key,
  statement_id    text references recon_statements(id),
  ledger_book_id  text references ledger_books(id),
  difference      numeric(20,8) not null,
  status          text not null check (status in ('open','investigating','resolved')),
  created_at      timestamptz not null default now()
);
```

### New Books per Feature
- `payroll_run` book per run
- `payout_batch` book
- `escrow` for payout links / subscriptions pending
- `refund_clearing`

---

## 4. Revised API Surface

```
# Existing
POST /v1/transactions/initialize
GET /v1/transactions/verify/{tx_ref}
POST /v1/payment_links

# NEW - Refunds
POST /v1/refunds
GET /v1/refunds/{id}

# NEW - Subscriptions
POST /v1/customers
POST /v1/subscription_plans
POST /v1/subscriptions
POST /v1/subscriptions/{id}/cancel

# NEW - Payouts
GET /v1/banks
POST /v1/beneficiaries
POST /v1/payouts
POST /v1/payouts/bulk
POST /v1/payout_batches/{id}/approve

# NEW - Payroll
POST /v1/employees
POST /v1/payroll_runs
POST /v1/payroll_runs/{id}/calculate
POST /v1/payroll_runs/{id}/approve
POST /v1/payroll_runs/{id}/disburse

# NEW - RAG + Swarm
POST /v1/compliance/ask
POST /v1/swarm/run
GET /v1/methods (now with smart routing rank)
GET /v1/admin/connectors/health
GET /v1/admin/recon/breaks

# Mobile push registration
POST /v1/devices/register
```

---

## 5. Revised Delivery Plan (Adds WP11-20)

| WP | Deliverable | Est (2 eng) |
|---|---|---|
| **WP11** | Refunds full + fee reversal policy + exam | 4-5 d |
| **WP12** | Second & third connectors (telebirr/cbe sandbox) + connector_configs + health sampler + circuit breaker | 7-10 d |
| **WP13** | Smart routing engine + routing_rules + /methods ranking + failover | 5-7 d |
| **WP14** | Subscriptions: plans, customers, subscriptions, invoices, dunning worker, customer portal | 10-14 d |
| **WP15** | Payouts: beneficiaries, batches, single/bulk, escrow book, maker-checker approvals | 8-10 d |
| **WP16** | Payroll: employees, runs, ET tax/pension engine, run books, posting M4, approvals | 14-18 d |
| **WP17** | RAG: pgvector setup, ingestion pipelines, embeddings, /compliance/ask with citations, eval | 10-12 d |
| **WP18** | Swarm: planner/critic, tool registry, swarm_sessions, confirmation thresholds, UI trace | 12-15 d |
| **WP19** | Recon: statement parser, breaks, suspense posting | 6-8 d |
| **WP20** | Flutter merchant app: auth, dashboard, link create/share, scan, approvals, FCM | 18-22 d |

**Original MVP 6-8 wks + Full = additional 12-16 wks**
**Total for 2 engineers: ~20-24 weeks to full commercial gateway + business OS.**
**Total for 1 engineer: ~32-40 weeks.**

**Milestones Revisited:**
- M4 MVP freeze (as before)
- M5 Refunds + 3 rails + Smart routing
- M6 Subscriptions + Payouts live
- M7 Payroll + ET packs
- M8 RAG + Swarm + Recon
- M9 Flutter beta + full E2E

---

## 6. NFR Impacts

- Payroll calc p99 < 2s for 500 employees
- RAG query p95 < 1.5s with pgvector
- Connector health sample every 30s
- Payout bulk: 1000 payouts queued < 10s
- Flutter: cold start < 2s, offline draft support

---

## 7. What to Change in Existing Docs

1. **MVP.md**: Change Scope Matrix table -> all ✅, update acceptance tests to include scripts 13-25 (refund, sub, payout, payroll, rag, routing)
2. **DATABASE.md**: Move §5 Post-MVP tables into §4 physical schema with DDL above
3. **SAD.md**: Add containers: rag-worker required (not optional), flutter app. Add component payroll/routing. Add ADR-009 Smart Routing, ADR-010 RAG with citations, ADR-011 Swarm safety.

---

## 8. Sign-off Questions for You

- Do you want all features in ONE v1 release, or keep MVP tag but start building full track in parallel (recommended: phased M5-M9)?
- For ET payroll, do you have the 2026 tax brackets and pension rules to encode, or should I use last published NBE/MoF brackets as placeholder?
- For Telebirr/CBE sandbox, do you already have sandbox API keys, or should we build with mock that matches expected real interface?

---

*This file is ready to replace the scope table in MVP.md. If you approve, I can generate:*
- *Updated DATABASE.md with full DDL*
- *Updated MVP.md / ROADMAP.md*
- *Go code skeleton with new modules: refund, subscription, payout, payroll, rag, swarm, routing, connector*
- *Flutter skeleton*

Tell me which to generate first.
