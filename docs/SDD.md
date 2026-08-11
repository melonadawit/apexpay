# ApexPay — Software Design Document (SDD)

> **Document type:** Software Design Document (SDD)
> **Version:** 1.0
> **Status:** Approved (matches committed `master`)
> **Audience:** Engineers implementing or extending ApexPay; reviewers.

---

## 1. Introduction

### 1.1 Purpose
This document details the concrete design of ApexPay's components: service/repository/handler
patterns, key flows (payments, ledger, assistant, payroll), interfaces, data structures, and the
algorithms used. It complements the SAD (`docs/SAD.md`) with implementation-level design.

### 1.2 Conventions
- **Layering:** `handler` (HTTP) → `service` (business rules) → `repository` (SQL/persistence).
- **Money:** `github.com/shopspring/decimal` everywhere money appears.
- **Errors:** `platform/errors` typed errors (`validation`, `not_found`, `unauthorized`).
- **Context:** request context carries authenticated identity via `platform/middleware`.

---

## 2. Request Lifecycle (HTTP → Handler → Service → Repository)

```
HTTP request
  -> middleware (RequestID, RealIP, Recoverer, Timeout, Compress, RateLimit, Auth, Locale)
  -> chi route -> Handler (decode JSON, read identity from context)
  -> Service (validate, apply business rules, orchestrate)
  -> Repository (pgxpool SQL, transactions)
  -> Response (WriteJSON / WriteError)
```

**Identity in context** (set by auth middleware):
- `CtxMerchantID`, `CtxUserID`, `CtxAPIKeyID`, `CtxRole`.
- `CtxLocale` (set by locale middleware): `en` | `am`.

---

## 3. Payment Flow (initialize → verify → ledger)

### 3.1 Initialize
1. `PaymentService.Initialize(ctx, req)` validates amount/currency, resolves routing, reserves an
   idempotency key, creates a `Payment` row.
2. Returns the payment with `requires_2fa` if amount > NBE threshold (2FA required).

### 3.2 2FA
- `Verify2FA(ctx, merchantID, paymentID, otp)` validates a TOTP/OTP against the payment.

### 3.3 Verify (ledger post)
- `Verify(ctx, req)` checks the connector result; on success it calls
  `UpdateStatusTx` which, in **one transaction**, updates the payment status **and** posts the
  double-entry journal (debit clearing / credit merchant payable + platform fee) to the operating
  ledger. This guarantees "never commit payment success without a ledger post in the same tx."

### 3.4 Idempotency
- `ReserveIdempotency` / `MarkConnectorStarted` / `FailIdempotency` manage the
  `idempotency_keys` state machine: `in_progress → connector_started → completed`.
- Same key + same input is safe; changed input returns 409 conflict.

---

## 4. General Ledger Design

### 4.1 Core types
- `Journal` — book id, posting key (idempotent), memo, reference type/id.
- `Entry` — account id, direction (debit/credit), amount, currency.
- `ledger_balances` — `(book_id, account_id)` amount, updated under advisory lock.

### 4.2 PostJournalTx algorithm
1. Validate `debits == credits` (`ValidateBalanced`) else reject.
2. `pg_advisory_xact_lock(hashtext(book_id))`.
3. Insert journal idempotently (`ON CONFLICT (book_id, posting_key) DO NOTHING`).
4. Insert entries idempotently (`ON CONFLICT (id) DO NOTHING`).
5. Upsert balances (`ON CONFLICT (book_id, account_id) DO UPDATE SET amount = amount + excluded`).
6. Commit.

### 4.3 Operating book + chart of accounts
`EnsureOperatingBook` creates the merchant's `merchant_operating` book with a standard chart
(assets, liabilities, equity, revenue, expense) on first use. Account codes are resolved to IDs
before posting (FK requirement).

### 4.4 Downstream posting modules
| Module | Journal posted |
|---|---|
| Payment success | debit clearing, credit merchant_payable + platform_fee |
| Refund | debit merchant_payable, credit clearing |
| Depreciation | debit `expense:depreciation`, credit `asset:accumulated_depreciation` |
| Inventory sale (COGS) | debit `expense:cost_of_sales`, credit `asset:inventory` |
| Tax liability | debit `asset:receivable`, credit `liability:tax` |
| FX revaluation | debit `asset:bank`/credit `revenue:fx_gain` (or loss) |
| Expense claim | debit `expense:operating`, credit `liability:payable` |
| Manual journal | user-supplied balanced lines |

### 4.5 Fiscal period close
- `fiscal_periods(merchant_id, period, status)` — `open` | `closed`.
- A closed period rejects new postings (`PostJournalEntry` checks `PeriodStatus`).
- Reopen is an explicit operator action (append-only lifecycle).

---

## 5. Apex Assistant Design

### 5.1 Request flow
`POST /v1/assistant/chat` (session-authenticated):
1. Middleware resolves locale → context.
2. Handler resolves **actor**: if the user maps to an employee record → `ActorEmployee`, else `ActorMerchant`.
3. `AssistantService.Chat(ctx, scope, text)`:
   a. `routeIntent(text)` — deterministic keyword rules → candidate intents.
   b. For each allowed tool (actor-gated), call `tool.Run(ctx, scope)` against live data.
   c. `compose(...)` — build a natural-language answer in the caller's locale.
   d. Persist append-only `assistant_messages` (user + assistant) under a `thread`.

### 5.2 Tools (read-only, actor-scoped)
- **Merchant:** summary (TPV), payments, invoices/aging, inventory, treasury/cash, loans,
  P&L, balance sheet.
- **Employee:** own YTD pay, own leave balance, own expense claims.

### 5.3 Interfaces
```go
type Tool struct {
    Name, Description string
    Actors []ActorType
    Run    func(ctx context.Context, scope Scope) (ToolResult, error)
}
type Scope struct { MerchantID, UserID, EmployeeID string; Actor ActorType; Locale i18n.Locale }
```
Wiring (`cmd/api/assistant_wiring.go`) adapts concrete repositories into a narrow `Readers`
interface (payments, invoices, inventory, treasury, loans, accounting, employee) via JSON-based
adapters.

### 5.4 i18n
`internal/i18n` exposes a `Catalog` with EN + Amharic message maps and a `Get(locale, key)`
with fallback. The assistant, payroll handlers, onboarding, fayda, banking, webhook, and portals
all localize user-facing output.

---

## 6. Payroll / Workforce Design (summary)

- `payroll.NewPgRepository(pool, ledger)` — full workforce OS: departments, designations, grades,
  branches, salary structures, employees, payroll runs, attendance, variable inputs, loans, leave,
  claims, final settlement, reports, employee portal.
- **Payroll run lifecycle:** create → attendance/variable inputs → calculate (V2 formula engine with
  ET tax brackets, proration, OT) → approve → disburse (posts ledger + generates bank/pension/ERCA files).
- **Expense claims:** manager → finance approval; on finance approval posts
  debit `expense:operating` / credit `liability:payable` (idempotent per claim).

---

## 7. Reconciliation, Tax, FX, Budget

### 7.1 Reconciliation
- `reconciliation.Service.Decide` — approves/rejects a reconciliation case, writes an audit log and an
  outbox event (with `::text` casts for jsonb params).

### 7.2 Tax
- `tax.Service.Schedule` — builds per-period VAT/TOT/withholding schedule from `tax_register`.
- `tax.Service.PostToLedger` — posts outstanding tax (debit receivable, credit `liability:tax`).
- Invoicing records collected VAT/withholding into the register via a `SetTaxRecorder` hook.

### 7.3 FX revaluation
- `fxreval.Service.Revalue` — revalues non-ETB current accounts at current rate; gain/loss posted to
  GL; baseline established on first run.

### 7.4 Budget / FP&A
- `budget.Service.Variance` — budget vs actual (from ledger) per category with variance %.
- `budget.Service.SetBudget` — upsert per `(merchant, period, category)`.

---

## 8. Procurement & Portals

### 8.1 Procurement / AP
- Vendors, purchase orders (with line items + tax), goods-received receipts, AP invoices with a
  PO 3-way-style match (`matched`/`mismatch`), and AP aging buckets.

### 8.2 Self-service portals
- `portal_access` table stores hashed, expiring tokens for vendors/customers.
- `POST /v1/portal/token` (session-gated) issues a token; `GET /v1/portal/me` (X-Portal-Token)
  returns the party's own invoices scoped to them — no dashboard login.

---

## 9. Data Structures & Algorithms

| Concern | Choice |
|---|---|
| Money | `decimal.Decimal` (shopspring). |
| Balance integrity | advisory lock + upsert, double-entry validated. |
| Intent routing | deterministic keyword rules (ADR-006). |
| Idempotency | unique `(merchant, key)` state machine. |
| Report caching | Redis read-through (cacheGet). |
| Aging | SQL CASE buckets `current/30/60/90/90plus`, ordered. |
| Tax brackets | lookup table; P&L via ledger category sums. |
| i18n | string catalog with fallback to English. |

---

## 10. Testing Strategy
- **Unit tests** per package (service/repository/domain), including error paths.
- **CI gates:** build, vet, `vet -tags integration`, gofmt, plus security/correctness lints.
- **End-to-end:** `docker-smoke` runs 29 sections against real Postgres via compose
  (auth, payments, payroll, GL, tax, FX, procurement, portals, assistant, i18n).
- **No float money / fin-privacy / audit append-only / SQL-param-cast** are lint-enforced.

---

## 10bis. Verified Live-Data Endpoints (recent additions)

The following handlers were previously stubs (returned empty/hardcoded) or missing routes;
they are now backed by real DB reads and exercised by the merchant dashboard and/or smoke suite.

| Endpoint | Module | What it now returns |
|---|---|---|
| `GET /v1/transactions/{id}` | `payment` | payment row + its ledger journals/entries (lifecycle) |
| `GET /v1/payout_batches`, `GET /v1/beneficiaries` | `payout` | list batches (with payout count) / beneficiaries |
| `GET /v1/subscriptions`, `GET /v1/subscriptions/{id}` | `subscription` | list; detail enriched with plan, customer, invoices |
| `GET /v1/refunds/{id}` | `refund` | real refund row (was hardcoded `succeeded`) |
| `GET /v1/payroll/payroll_runs` | `payroll` | real `ListRuns` DB query (was empty) |
| `GET /v1/payroll/final_settlements` | `payroll` | real `ListFinalSettlements` DB query (was empty) |
| `GET /v1/payroll/tax_brackets` | `payroll` | active ET brackets from `payroll_tax_brackets` |
| `GET /v1/payroll/payroll_reports/cost_center` | `payroll` | stored cost-center report metadata (falls back to mock) |

Related domain types gained snake_case JSON tags (`Payment`, `PayrollRun`/`Item`/`Calendar`/
`FinalSettlement`, `PayoutBatch`/`Payout`/`Beneficiary`, `Subscription`/`Plan`/`Customer`/
`Invoice`, `Refund`, `TaxBracket`) so the frontend `lib/api/client.ts` reads them directly.

---

## 11. Extension Points
- **New payment connector** — implement the `connector` adapter + register it.
- **New assistant tool** — add a `Tool` with actor gate and a `Readers` adapter.
- **New report** — add an `accounting`/`payroll` repo method + handler route.
- **New worker** — add under `internal/worker/<name>` and wire scheduling.
- **New locale** — extend the `i18n` catalog maps.
