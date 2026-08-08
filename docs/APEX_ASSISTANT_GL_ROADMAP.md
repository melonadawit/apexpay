# ApexPay — Assistant, Comprehensive GL & Management-Suite Design & Roadmap

> **Status:** Part A (P1, read-only assistant) ✅ · Part B (real GL: journal entries + period
> close + depreciation + inventory COGS + **tax schedules**) ✅ · Part C (1) Procurement/AP
> ✅ · Part C (2) Expenses largely covered by payroll claims; Part B remainder
> (multi-currency revaluation) — next.
>
> Senior Engineering Manager view — grounded in what already exists in the codebase.
> Goal: deliver the *Apex Assistant* (role-scoped conversational agents for merchants,
> vendors, employees, employers), turn the accounting engine into a real, light-weight
> double-entry General Ledger that beats legacy desktop products (QuickBooks, Peachtree),
> and add the management modules that make ApexPay a true embedded finance operating system.
> Core constraint honored throughout: the product must stay **API-first, Postgres-native,
> real-time, and light** — never a heavyweight desktop app.

---

## 0. Why this is the right move (and what already exists)

ApexPay already contains most of the *hard primitives* we'd otherwise have to build:

| Capability | Where it exists today | What's missing |
|---|---|---|
| Double-entry ledger engine | `internal/ledger` — `PostJournalTx`, `GetBalance`, `ListJournalsByRef`, `ValidateBalanced`, idempotent posting keys, advisory-lock balance updates. Payments/refunds/payouts/payroll already post through it. | A user-facing way to *create* journal entries; period close; inventory costing. |
| Accounting reports | `internal/accounting` — CoA, trial balance, P&L, balance sheet, cash flow. | **Read-only** — no journal posting, no period close, no fiscal-period awareness in reports. |
| Conversational assistant core | `internal/rag` — embed → vector search → rerank → LLM with **mandatory citations** + anti-hallucination threshold. | Scoped only to compliance ("ask_compliance"); no structured-data grounding. |
| Agent orchestrator | `internal/swarm` — tool registry with **role + amount-threshold gating** (`create_payment_link`, `list_payments`, `create_refund`, `create_payout`, `calculate_payroll_draft`, `get_tpv`, `ask_compliance`). | Very small tool set; not wired to a general assistant UI or per-actor scoping. |
| Data-security posture | RBAC roles, `internal/team` approvals, `internal/auth` sessions. | Agents must *inherit* this scoping, not bypass it. |

**Conclusion:** we are not building from scratch. We are *productizing* a strong core. This
is a moat — legacy desktop accounting (QuickBooks, Peachtree) is batch-close, single-currency,
desktop-bound. We are real-time, API-first, multi-currency (ETB edge), and audit-append-only
(already a passing CI check). Our advantage is **integration + real-time + embedded UX**, not
re-inventing GAAP.

---

## 1. Part A — The Apex Assistant (conversational agents)

### 1.1 Product concept
A single assistant, "Apex Assistant," that appears as a chat surface for **four actor types**,
each with its own role-scoped view:

| Actor | Sees | Typical questions |
|---|---|---|
| **Merchant / Owner** | Whole-org finance, inventory, payroll, accounting, treasury | "What's my gross margin this month?", "Which invoices are overdue?", "Run the payroll draft." |
| **Employee** | Own payslips, leave, expenses, own reimbursements | "When was I last paid and how much?", "What's my remaining annual leave?", "Submit my travel expense." |
| **Vendor / Supplier** | Only their own invoices, purchase orders, payment status, AP aging | "Which invoices are unpaid?", "When will invoice INV-123 be paid?" |
| **Employer / Admin** | Org-level (subset of merchant owner, dept-scoped) | "Approve the pending refunds.", "Cost center P&L for sales." |

### 1.2 Architecture (reuse, don't rebuild)
- **Orchestrator:** extend the existing `internal/swarm` registry instead of writing a new agent loop.
- **Grounding in structured data via TOOLS (primary):** the assistant answers exact money questions by
  calling tools that query the ledger/repos (exact, auditable numbers). **RAG on documents (secondary):**
  for policies, contracts, compliance, and knowledge-base answers with citations.
- **Per-actor scoping layer:** every tool call resolves the actor → merchant → resource scope before
  executing. No tool may query data outside the actor's scope. This is a hard invariant.
- **LLM gateway (pluggable):** abstract `Embedder`, `VectorStore`, `LLM` behind the existing
  `internal/rag` interfaces so we can swap providers (local/cloud) without touching domain code.

### 1.3 Safety & control (non-negotiable)
- **Read-first rollout.** Phase 1 = read-only Q&A ("what is…"). Phase 2 = *proposals* ("I recommend
  approving refund X"). Phase 3 = hand-off to the **existing approval workflow** — the assistant
  *recommends*; a human *approves*. The assistant never auto-executes money movement.
- **Reuse the existing threshold + role gates** in the swarm tool registry; extend coverage to every
  new tool.
- **Full audit trail:** every assistant turn and tool call is append-only logged (already the norm).
- **Citations mandatory** for any document-derived answer (existing `rag` behavior). For structured
  answers, return the exact tool result + source (journal/report id).

### 1.4 Tool registry to build (grouped by module)
- **Finance/payments:** list payments, refund, payout, TPV, reconciliation status.
- **Accounting:** trial balance, P&L, balance sheet, cash flow, aging, expense-by-category trend.
- **Inventory:** stock levels, low-stock, valuation, reorder suggestions.
- **Payroll/HR:** own payslip, payroll draft, leave balance, expense submission.
- **Vendor:** own invoices, PO status, AP aging.
- **Compliance:** cited policy answers (already exists).

### 1.5 Phased plan (Apex Assistant)
- **A1 — Read-only assistant core** (merchant/owner + employee own-data): assistant session + chat
  endpoint, actor resolution, ~10 scoped tools, structured-grounding pipeline. **Definition of done:**
  new integration smoke section where an employee query returns only their data and a cross-scope
  query is rejected.
- **A2 — Vendor + employer actors** and NL→tool routing improvements, more tools.
- **A3 — Action proposals + approval hand-off** using the existing team/approval workflow.
- **A4 — Assistant observability** (turn logs, token cost, feedback thumbs, refusal analytics).

---

## 2. Part B — Comprehensive, light-weight General Ledger

### 2.1 What we keep
Double-entry journals, balanced-entries validation, idempotent posting, advisory locks, live reports.
This is *already* lighter and more modern than QuickBooks/Peachtree.

### 2.2 What we add to make it a complete GL
1. **Journal-entry creation (user-facing).** `POST /v1/accounting/journal-entries` with balanced
   debit/credit lines, memo, optional reference; validate balance before commit; post via the ledger
   service. Manual adjusting entries become possible.
2. **Fiscal period close.** Open/closed periods (month/quarter/year). A closed period rejects new
   postings; closes are append-only and reversible only via a reversing entry (never a delete). Reports
   accept `from`/`to` periods and respect locks.
3. **Inventory costing → COGS.** On product (re)valuation and sale, cost by **FIFO and weighted-average
   (AVCO)**; post COGS to the ledger. This is what makes "inventory" a real balance-sheet asset and
   gives P&L a true cost of sales line.
4. **Depreciation posting.** We already have `fixed_assets` + a depreciation schedule; post monthly
   depreciation entries into the ledger so accumulated depreciation flows to the balance sheet.
5. **Multi-currency revaluation.** Journal entries already carry `currency`; add revaluation of
   FX-denominated balances at period close with an FX gain/loss entry (reuse `internal/banking` forex).
6. **Tax schedules.** VAT/TOT/withholding payable from invoice/payment flow into tax liability accounts
   (Ethiopian-law compliance docs already exist; wire them into the GL).

### 2.3 "Light, not heavy" engineering notes
- Keep the ledger as a **thin, correct core** (primitive is commodity) — no over-engineering.
- Use **report caching** (already present in `accounting/repository.go`) so P&L/BS don't recompute on
  every request; invalidate on new postings.
- Everything is **real-time**: a posting updates balances immediately; no nightly batch close is
  required (period locks are logical, not batch jobs).
- Design for **scale with indexes**, not by turning into a desktop monolith.

### 2.4 Definition of done (Part B)
- New integration smoke section: create an unbalanced entry → rejected; create a balanced entry → posts
  and moves trial balance; close a period → post rejected; reopen via reversing entry → allowed and
  audited; inventory sale → COGS reflects in P&L.

---

## 3. Part C — Management systems to build into ApexPay

Ranked by leverage and fit with the existing ledger + approvals:

1. **Procurement / Accounts Payable (highest fit).** Vendors, purchase orders, receipts, **3-way match**
   (PO ↔ receipt ↔ invoice), AP aging, payment scheduling. Feeds the ledger and dovetails with the
   existing `vendor_invoices` + reconciliation. Direct SME pain point.
2. **Expense Management & Reimbursement.** Employee expense claims (categories, receipts), approval
   flow (existing `team`), posting to ledger, and reimbursement via existing payout/payroll. Pairs with
   the assistant's employee actor.
3. **Inventory costing + reorder** (feeds Part B item 3): stock valuation, reorder points, low-stock
   alerts, vendor reorder suggestions.
4. **Budgeting / FP&A.** Budgets per cost-center, budget vs actual variance, forecasts on top of the
   existing cash-flow forecasting.
5. **Time & Attendance → Payroll.** Clock-in/out, leave, overtime that feed the existing payroll engine
   (already has ET tax/pension, OT).
6. **Vendor & customer self-service portals.** Read-only portals over the Part-A agents.

> Do these **one at a time** behind the same rigor as existing modules: build + vet + test +
> integration-vet + gofmt + a smoke section + security lints, before starting the next.

---

## 4. Prioritized roadmap (recommended sequence)

| Phase | Scope | Why first |
|---|---|---|
| **P1** | Part A (A1): read-only Apex Assistant for merchant + employee | Highest moat; reuses swarm/rag; read-only = safe; directly useful. |
| **P2** | Part B: journal-entry creation + period close | Makes accounting a *real* GL (currently read-only); unlocks inventory/expense posting. |
| **P3** | Part C (1) Procurement/AP + (2) Expenses | Plug straight into ledger + approvals; classic SME pain; feeds assistant tools. |
| **P4** | Part B remainder: inventory costing, depreciation, multi-currency, tax | Completes financial statements (COGS, FX, tax). |
| **P5** | Part A (A2–A4) + Part C (3–6) | Vendor/employer actors, proposals+approvals, budgeting, attendance, portals. |

**Guardrail against scope sprawl:** each phase must finish (all checks + smoke + docs) before the next
begins. We harden and complete the core rather than adding many half-finished modules.

---

## 5. Engineering principles (locked)

- **No reinventing GAAP/IFRS.** Journaling is a commodity primitive; value = integration, real-time,
  embedded UX, multi-currency/ETB.
- **Structured data via tools, documents via RAG-with-citations.** No hallucinated money numbers.
- **Per-actor scoping is a hard invariant** — the assistant inherits RBAC; it never bypasses it.
- **Human-in-the-loop for actions.** Assistant proposes; approvals execute.
- **Append-only + auditable everywhere** (consistent with the passing `audit-append-only` CI check).
- **Light by construction:** thin correct core, report caching, indexes over complexity, no batch
  desktop-style processing.
- **Environment discipline:** no Go toolchain/caches in the workspace; all builds/CI run containerized
  or with toolchain + caches external to `/home/user`.
