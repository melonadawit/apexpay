# ApexPay

**The all-in-one financial operating system for Ethiopia.**

ApexPay is an AI-native, API-first platform that unifies **payments**, a real-time **double-entry
General Ledger**, **workforce & payroll**, **procurement/AP**, **inventory**, **tax**, **FX**,
**budgeting**, and a role-scoped **AI assistant** — all under one compliant, Ethiopia-native core.

> **Setup & operations:** see [`docs/RUNBOOK.md`](docs/RUNBOOK.md).
> **Architecture:** see [`docs/SAD.md`](docs/SAD.md).
> **Design:** see [`docs/SDD.md`](docs/SDD.md).
> **Investor overview:** see the ApexPay investor deck (PDF).

---

## Why ApexPay

- **Complete platform, not a feature.** 43 business modules share one ledger, so a sale, a payroll
  run, a vendor invoice, and a tax payment all reconcile automatically.
- **AI assistant built in.** A role-scoped Apex Assistant (RAG + tool-calling agents) answers finance,
  inventory, and payroll questions for merchants, employees, vendors, and employers — in English and
  Amharic.
- **Lighter than legacy accounting.** Real-time, API-first, Postgres-native GL — not a heavyweight
  desktop QuickBooks-style product.
- **Ethiopia-native compliance.** NBE ONPS/10/2025, ET income-tax brackets, pension 7%/11%, VAT/TOT,
  Fayda ID verification, Ethiopian calendar, labour proclamation 1156/2019.
- **Production-grade.** 10/10 CI gates green, including a 29-section end-to-end smoke suite.

---

## Repository Layout

```
apexpay/
├── services/api/            # Go API core (single money engine)
│   └── internal/
│       ├── ledger/          # double-entry GL (ValidateBalanced, advisory locks)
│       ├── accounting/      # journal entries, fiscal close, reports
│       ├── payment/         # initialize / 2FA / verify + ledger post
│       ├── payroll/         # workforce OS (runs, leave, claims, loans, reports)
│       ├── assistant/       # role-scoped conversational AI
│       ├── swarm/           # agent orchestrator
│       ├── rag/             # grounded answers with citations
│       ├── procurement/     # vendors, POs, receipts, AP
│       ├── tax/  fxreval/  budget/  treasury/   # finance ops
│       ├── portal/          # vendor & customer self-service
│       └── platform/        # config, crypto, middleware, i18n, storage, http
├── apps/
│   ├── merchant-web/        # Next.js merchant dashboard
│   ├── admin-web/           # Next.js admin / compliance console
│   ├── checkout-web/        # Next.js hosted checkout
│   ├── mobile/              # Flutter mobile app
│   └── docs/                # documentation site
├── db/migrations/           # 44 versioned PostgreSQL migrations
├── deploy/docker/           # Dockerfiles + docker-compose
├── scripts/                 # integration smoke suite, seeders
└── tests/integration/       # seed + integration tests
```

---

## Quick Start (development)

The API, workers, and web apps run via Docker Compose (PostgreSQL, Redis, MinIO, API, worker,
merchant-web, checkout-web). See [`docs/RUNBOOK.md`](docs/RUNBOOK.md) for full setup.

```bash
docker compose -f deploy/docker/docker-compose.yml up --build
```

**Run the end-to-end smoke suite** (auth, payments, payroll, GL, tax, FX, procurement, portals,
assistant, i18n):

```bash
docker compose -f deploy/docker/docker-compose.yml \
  -f deploy/docker/docker-compose.integration.yml \
  up --build --abort-on-container-exit --exit-code-from integration integration
```

---

## Key Stakeholder Flows

| Stakeholder | What they see / do |
|---|---|
| **Merchant / Owner** | Whole business: payments, cash, P&L, balance sheet, inventory, invoices, budget, tax. Approves payroll, closes periods, posts journals. AI assistant answers finance questions. |
| **Employee** | Own YTD pay, leave balance, expense claims, payslips. AI answers only their own data. |
| **Vendor / Supplier** | Self-service portal: their AP invoices and payment status only. |
| **Customer** | Hosted checkout: Telebirr, CBE Birr, bank, card, EthSwitch QR; 2FA above NBE threshold. |
| **Admin / Compliance** | Admin dashboard: onboarding queue, KYC/Fayda exam, compliance checks, approve/reject. |

---

## Architecture Highlights

- **One Go monolith API core** — a single money engine shares the ledger (ADR-001).
- **PostgreSQL as source of truth** (44 migrations); Redis as cache; MinIO as document vault.
- **Double-entry GL** with advisory-lock balance updates and idempotent posting.
- **Outbox pattern** + workers for reliable async jobs (webhooks, dunning, recon, notifications).
- **Append-only audit** and append-only assistant history for full traceability.
- **i18n** — per-user English/Amharic across API and UI.

See [`docs/SAD.md`](docs/SAD.md) for the full architecture and `docs/SDD.md` for component design.

---

## Quality & Security

10/10 CI gates green:

- Build & type-check (`go build`, `go vet`, `gofmt`)
- Static security: `gosec`, `trivy`, `gitleaks`
- Correctness lints: `no-float-money`, `fin-privacy`, `audit-append-only`, `sql-param-cast`
- Frontend: `lighthouse-axe`
- End-to-end: `docker-smoke` (29 sections)

Money is always `decimal.Decimal`. Financial journals must balance. Audit logs cannot be updated or
deleted. No plain FIN is ever logged or returned.

---

## Compliance (Ethiopia)

- **NBE** ONPS/10/2025 (2FA above threshold), payment system oversight.
- **Income tax** — ET brackets, pension employee 7% / employer 11%.
- **VAT / TOT** and withholding schedules.
- **Fayda** national ID verification (hash + last4, OTP consent).
- **Labour Proclamation 1156/2019** — leave (annual/sick/maternity), severance, final settlement.

See `APXPAY_ETHIOPIA_LAW_COMPLIANCE.md` and `FULL_PLATFORM_SPEC.md`.

---

## License

Proprietary. © ApexPay.
