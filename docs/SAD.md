# ApexPay — Software Architecture Document (SAD)

> **Document type:** Software Architecture Document (SAD)
> **Version:** 1.0
> **Status:** Approved (matches committed `master`)
> **Audience:** Architects, engineers, reviewers, technical investors

---

## 1. Introduction

### 1.1 Purpose
This document describes the software architecture of **ApexPay**, a multi-tenant, API-first financial
operating system for Ethiopia. It defines the system's structure, components, key interfaces, data
model, security model, deployment model, and the architectural decisions that shape it.

### 1.2 Scope
ApexPay unifies payments, a real-time double-entry General Ledger (GL), workforce & payroll,
procurement/AP, inventory, tax, FX, budgeting, and an AI assistant into one composable platform. This
SAD covers the backend (`services/api`), the four client surfaces, the data stores, and the CI/quality
gates. It does **not** cover the operational runbook in depth (see `docs/RUNBOOK.md`).

### 1.3 Definitions / Acronyms
- **API** — Application Programming Interface.
- **GL** — General Ledger (double-entry accounting).
- **NBE** — National Bank of Ethiopia.
- **RAG** — Retrieval-Augmented Generation (grounded AI answers).
- **ONPS** — NBE's Oversight of Payment Systems directive.
- **VAT / TOT** — Value Added Tax / Turnover Tax.
- **COGS** — Cost of Goods Sold.

---

## 2. Architectural Goals & Constraints

### 2.1 Quality attributes (priorities)
| Attribute | Priority | How it is achieved |
|---|---|---|
| **Correctness / money safety** | Highest | All money is `decimal` (never float); double-entry journals validated (`ValidateBalanced`); append-only audit; SQL-param-cast CI lint. |
| **Security** | Highest | gosec + trivy + gitleaks gates; API-key hashes (SHA-256), sessions, RBAC roles, per-actor scoping; FIN privacy (last4 only). |
| **Auditability** | High | Append-only `audit_logs` with DB trigger; append-only assistant threads/messages. |
| **Maintainability** | High | One Go core, 43 cohesive internal packages, narrow interfaces, gofmt/vet/tests gated in CI. |
| **Localization** | High | Per-user English/Amharic via an `i18n` catalog + locale middleware. |
| **Performance** | Medium | Postgres-native queries, report caching (Redis), advisory locks for balance updates. |
| **Deployability** | High | Docker multi-stage builds, distroless runtime, compose for dev/integration. |

### 2.2 Constraints
- **Ethiopia regulatory compliance** — NBE ONPS/10/2025, ET income-tax brackets, pension 7%/11%,
  VAT/TOT, labour proclamation 1156/2019.
- **Local rails** — Telebirr, CBE Birr, EthSwitch QR, bank/IPS multi-bank disbursal.
- **No float money** — enforced by a CI lint over payment/ledger/refund/payout/payroll/worker packages.
- **Bilingual** — English + Amharic across API and UI.
- **No heavyweight desktop accounting** — real-time, API-first, light.

---

## 3. System Context

```
                     ┌─────────────────────────────────────────────┐
                     │                 ApexPay                      │
   Merchant / Admin  │                                             │
   (dashboard users) │  ┌─────────┐ ┌─────────┐ ┌─────────┐        │
        │            │  │Merchant │ │  Admin  │ │ Checkout│        │
        │            │  │  Web    │ │   Web   │ │   Web   │        │
        ▼            │  └────┬────┘ └────┬────┘ └────┬────┘        │
                     │       │           │           │             │
   Employee / Vendor │  ┌────▼───────────▼───────────▼────────┐     │
   (portal / token)  │  │        /v1 REST API (Go core)       │     │
        │            │  └────┬────────┬────────┬────────┬────┘     │
        ▼            │       │        │        │        │          │
                     │  ┌────▼──┐ ┌───▼──┐ ┌───▼───┐ ┌─▼────┐      │
                     │  │Postgres│ │Redis │ │MinIO  │ │Workers│     │
                     │  └───────┘ └──────┘ └───────┘ └──────┘      │
                     └─────────────────────────────────────────────┘
```

The system is a **monolith API** (by design — one cohesive money core) surrounded by thin clients and
supporting stores. Workers run the same codebase for background jobs (outbox drain, webhooks, dunning,
recon, notifications, forex cache, credit scoring).

---

## 4. Architecture Views

### 4.1 Layered / logical view
1. **Presentation** — Merchant Web, Admin Web, Checkout Web, Mobile, Employee Portal.
2. **AI & Intelligence** — Apex Assistant (RAG + tool-calling agents), Swarm orchestrator,
   rules/LLM routing, i18n.
3. **Core Business Engine** — payments, ledger, payroll, accounting, procurement, inventory,
   tax, FX, budget, risk, reconciliation.
4. **Data & Platform** — PostgreSQL, Redis, MinIO, outbox/events, workers.

### 4.2 Module (package) view — `services/api/internal`
43 business/platform packages, e.g.:
- **Money core:** `ledger`, `accounting`, `payment`, `refund`, `payout`, `payroll`, `reconciliation`
- **Commercial:** `checkout`, `link`, `subscription`, `invoicing`, `inventory`, `procurement`, `lending`, `loyalty`, `dispute`
- **Workforce:** `payroll`, `hris`, `team`
- **Finance ops:** `treasury`, `budget`, `fxreval`, `tax`, `fixedasset`
- **Trust & identity:** `auth`, `onboarding`, `fayda`, `compliance`, `admin`, `bankverification`, `risk`
- **Intelligence:** `assistant`, `swarm`, `rag`
- **Platform:** `config`, `crypto`, `dbpool`, `errors`, `http`, `logger`, `middleware`, `secrets`, `storage`, `twofa`, `i18n`

### 4.3 Runtime / concurrency view
- HTTP server (chi router) under `/v1`, with rate limiting, request ID, real IP, timeout, recoverer.
- `pgxpool` connection pool (tuned).
- Advisory locks (`pg_advisory_xact_lock`) for concurrent balance updates in the ledger.
- Background worker processes for async jobs (outbox, webhooks, dunning, recon, notifications, forex cache).
- Redis-backed report caching with read-through invalidation.

### 4.4 Deployment view
- **API + Worker** — Go multi-stage Docker builds → distroless static runtime (non-root).
- **PostgreSQL 17** — source of truth; 44 versioned migrations (0011 pgvector is optional; skipped in sandbox).
- **Redis** — cache / sessions (fallback to in-memory when unavailable).
- **MinIO** — document vault (presigned URLs).
- **docker-compose** for dev + a compose-based integration smoke runner.

---

## 5. Data Architecture

### 5.1 Database
- **PostgreSQL** is the single source of truth. 44 migrations (`db/migrations/0001…0044`).
- Money columns are `numeric(20,8)`; amounts are carried in Go as `decimal.Decimal`.
- Key tables: `merchants`, `users`, `api_keys`, `payments`, `payment_links`, `ledger_books`,
  `ledger_accounts`, `ledger_journals`, `ledger_entries`, `ledger_balances`, `payroll_*`,
  `invoices`, `ap_invoices`, `vendors`, `purchase_orders`(+items), `receipts`, `products`,
  `tax_register`, `fx_revaluations`, `budgets`, `fixed_assets`, `depreciation_entries`,
  `fiscal_periods`, `portal_access`, `assistant_threads`, `assistant_messages`, `audit_logs`.

### 5.2 The General Ledger (core invariant)
- **Double-entry:** every journal must balance (debits == credits) — validated by
  `ledger.ValidateBalanced` before insert.
- **Idempotent posting:** `(book_id, posting_key)` unique so re-posting is safe.
- **Balance integrity:** `ledger_balances` upserted under an advisory lock per book.
- **Append-only:** journal entries are never updated/deleted.
- Modules post into the same operating book: payments, refunds, payouts, payroll, depreciation,
  inventory COGS, tax liability, FX revaluation, expense claims, manual journal entries.

### 5.3 Append-only audit
- `audit_logs` table guarded by a DB trigger that prevents UPDATE/DELETE.
- Assistant threads and messages are also append-only for a full conversational audit trail.

---

## 6. Security Architecture

### 6.1 Authentication & authorization
- **API keys** — full secret hashed (SHA-256); only the 12-char prefix is queried; RBAC scopes.
- **Dashboard sessions** — opaque tokens stored as hashes, with expiry and touch.
- **RBAC roles** — `owner`, `admin`, `developer`, `finance`, `support`, `ops`, `compliance`, `viewer`.
- **Per-actor scoping** — the assistant and portals resolve actor (merchant/employee/vendor/customer)
  and refuse cross-scope access.

### 6.2 Data protection
- **FIN privacy** — only last4 ever returned/logged; FIN stored as salted hash.
- **Account masking** — `account_number_masked`, hash for verification.
- **Secrets** — derived keys via `platform/crypto`, connector encryption key, JWT secrets.
- **Transport** — HTTPS in production; Compress middleware; no plain secrets in images.

### 6.3 CI quality gates (10/10 green)
`go-build`, `gosec`, `trivy` (SCA), `gitleaks`, `no-float-money-lint`, `fin-privacy-lint`,
`audit-append-only`, `sql-param-cast-lint`, `lighthouse-axe`, `docker-smoke` (29 e2e sections).

---

## 7. Architectural Decisions (ADRs)

| ADR | Decision | Rationale |
|---|---|---|
| ADR-001 | One Go monolith API core | Single money core, shared ledger, no distributed-transaction complexity for money. |
| ADR-002 | Postgres as source of truth; Redis only cache | Durability & ACID for money; cache is disposable. |
| ADR-003 | `decimal.Decimal` for all money | No float error; CI lint enforces it. |
| ADR-004 | Double-entry GL with advisory locks | Balance integrity under concurrency. |
| ADR-005 | Outbox pattern for async jobs | Reliable, at-least-once events; workers idempotent. |
| ADR-006 | AI off the hot path | Assistant/Swarm route via deterministic rules first; LLM optional/later. |
| ADR-007 | RAG with mandatory citations + threshold | No hallucinated money answers; no answer below score. |
| ADR-008 | Per-user i18n via catalog + locale middleware | Clean single-language output; EN + Amharic. |
| ADR-009 | Append-only audit + assistant history | Regulatory + conversational auditability. |
| ADR-010 | SQL-param-cast CI lint | Prevents a whole class of pgx "cannot infer type" errors. |

---

## 8. Risks & Mitigations
- **Scope sprawl** — mitigated by phased roadmap and a finish-before-add guardrail.
- **LLM quality on money** — mitigated by read-only tools grounded in the ledger, human-in-the-loop
  for actions, citations mandatory.
- **Multi-tenancy isolation** — every query is scoped by `merchant_id`; portals by entity.
- **Operational** — Docker/distroless, health endpoints, worker idempotency, reconciliation.

---

## 9. References
- `docs/RUNBOOK.md`
- `docs/APEX_ASSISTANT_GL_ROADMAP.md`
- `FULL_PLATFORM_SPEC.md`
- `APXPAY_ETHIOPIA_LAW_COMPLIANCE.md`
