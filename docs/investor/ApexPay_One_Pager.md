# ApexPay — Investor One-Pager

**The all-in-one Financial Operating System for Ethiopia.**

ApexPay is an AI-native, API-first platform that unifies **payments**, a real-time **double-entry
General Ledger**, **workforce & payroll**, **procurement/AP**, **inventory**, **tax**, **FX**,
**budgeting**, and a role-scoped **AI assistant** — one compliant, Ethiopia-native core.

---

## What was verified this cycle (all green on `master`)

- **Merchant dashboard on live data.** Payments (list + transaction detail with ledger journals),
  subscriptions (+ detail with plan/customer/invoices), refunds, payouts, and payroll (runs,
  calendars, final settlements, reports, tax brackets) all render from a **real seeded Postgres**,
  not mockups.
- **Developer portal.** Real, DB-backed API-key management (list/create/revoke) — created keys are
  immediately usable as Bearer tokens — plus webhooks endpoints & deliveries.
- **Embedded finance.** Lending, escrow, corporate cards, credit lines, and virtual accounts all
  return live seeded data.
- **Mobile (Flutter).** Tests pass 3/3; analyzer clean; Android APK built as a CI artifact on a
  large runner.
- **Security audit** (`docs/SECURITY_AUDIT.md`). Money-safety, FIN-privacy, secret-handling, and
  webhook SSRF controls verified. A real gap was found and closed: `audit_logs` is now **truly
  append-only** via a DB trigger, and CI now fails if the trigger is ever removed.

---

## Why ApexPay wins

- **Complete platform, not a feature** — 43 modules share one ledger, so a sale, a payroll run, a
  vendor invoice, and a tax payment all reconcile automatically.
- **AI assistant built in** — role-scoped, grounded in live ledger/data, English + Amharic.
- **Ethiopia-native compliance** — NBE ONPS/10/2025, ET income-tax brackets, pension 7%/11%,
  VAT/TOT, Fayda ID, labour proclamation 1156/2019.
- **Production-grade** — 10/10 CI gates green, including a 29-section end-to-end smoke on real
  Postgres; money is always `decimal.Decimal`.

---

## The opportunity

Defensible moat (integrated ledger + AI assistant), land-and-expand from payments to the full
finance stack, and **live proof** — the entire data path is running against real data, not mockups.

*References: `docs/SAD.md`, `docs/SDD.md`, `docs/SECURITY_AUDIT.md`, the full investor deck.*
