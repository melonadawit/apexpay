# ApexPay — Full Platform v1.1.0 (NBE Onboarding + Fayda + Outstanding UI)

> **🚀 Full setup, login, and troubleshooting instructions: see [`docs/RUNBOOK.md`](docs/RUNBOOK.md).**

This repository implements the **full platform** per `docs/ROADMAP_M5_M9.md` and updated `DATABASE.md` + `MVP.md` v1.1.0-full.

## What Changed vs MVP 1.0.0

Previously ❌ subs, payouts, payroll, RAG, swarm, smart routing, Flutter, refunds stub → now ✅ ALL included plus:

- **NBE-Grade Merchant Onboarding** (ONPS/02/2020, ONPS/09/2023, ONPS/10/2025) — 6-step wizard outstanding modern UI like Mercury/Linear
- **Fayda National ID Verification** — 12-digit FIN hash + FAN, front/back images <2MB, selfie, OTP consent via id.gov.et VeriFayda 2.0 / OIDC eSignet, offline QR FaydaEncode, face match 0.85 threshold, privacy hashed storage, MinIO encrypted vault
- Required docs checklist per PayAtlas ET PSP: company registration (notarized), TIN 10-digit, business license, VAT, memorandum, board resolution, shareholder list, UBO front/back, proof of address, bank letter, website refund/privacy/terms screenshots

## Repository Layout (Target)

```
apexpay/
├── services/api/internal/
│   ├── onboarding/ (domain, service, repository - NBE checklist O(n) validation, risk scoring weighted sum)
│   ├── fayda/ (domain, service, verifier_mock/live Strategy pattern, privacy hash)
│   ├── refund/ (domain, service, fee reversal pro_rata/non_refundable/full decimal precise)
│   ├── subscription/ (plans, customers, subs FSM, dunning 1d/3d/5d)
│   ├── payout/ (beneficiaries fuzzy match Levenshtein, batches, bulk 1000, payout links escrow book)
│   ├── payroll/ (employees, runs per-book, ET tax brackets binary search O(log n), pension 7/11, OT 1.25/1.5/2.0)
│   ├── routing/ (rules priority sort O(n log n), health 5m success_rate, circuit breaker 5 fails -> open 60s)
│   ├── rag/ (pgvector 1024, chunk 800/100, embed batch 32, threshold 0.65 no hallucination guard)
│   ├── swarm/ (planner RulesPlanner keyword, critic threshold, executor ToolExecutor O(1) registry, confirmation modal)
│   ├── ledger/ (ValidateBalanced debit==credit invariant)
│   └── connector/ (interface mock/telebirr/cbe/bank/ethswitch/card)
├── apps/mobile/ Flutter 3.22+ Riverpod + outstanding UI (glass, confetti, QR scanner overlay, onboarding 6-step)
├── libs/ui/tokens.json — design tokens ET Green #0B6E4F, gold #EAB308, radius xl 24, motion ease [0.22,1,0.36,1]
├── db/migrations/ 0001..0012 full
├── deploy/docker/docker-compose.yml postgres+pgvector+redis+minio+api+worker+merchant-web+checkout-web
└── docs/ROADMAP_M5_M9.md
```

## Quick Start

```bash
docker compose -f deploy/docker/docker-compose.yml up -d
go run ./services/api/cmd/api  # :8080
# Merchant web
cd apps/merchant-web && npm install && npm run dev # :3000
# Checkout web
cd apps/checkout-web && npm install && npm run dev -- -p 3001
# Flutter
cd apps/mobile && flutter run
```

## Onboarding Flow Outstanding Modern

1. Business Info — smart form industry search, restricted gambling/crypto/adult blocked per NBE
2. Owners & Fayda — Fayda capture modal with corner guides animated pulse, glare detection, front/back/selfie, FIN/FAN + OTP 6-digit
3. Bank — combobox bank list CBE/Awash/Dashen with logos, account name auto-check must match legal
4. Docs Vault — dropzone dashed pulse on drag, preview thumbs, hash integrity, progress donut 0-100%
5. Compliance Preview — risk gauge chart, checks green/red, timeline
6. Review Submit — confetti, terms consent NBE

Tracking Timeline vertical like Linear, Kanban admin board drag-drop, Fayda verification card front blurred + face_score 0.92 + OTP verified check.

## Security Best Practices Implemented

- FIN hashed sha256(salt+fin) + last4 only, never plain in logs; presigned MinIO URLs 15m TTL; file hash integrity; encrypted at rest AES-256 SSE-S3
- Rate limit Fayda OTP 5/hour/IP via Redis token bucket
- 2FA mandatory >5000 ETB per ONPS/10/2025
- Maker-checker >50k payout, >100k payroll
- No float money: shopspring/decimal + numeric(20,8)
- Circuit breaker optimal in-memory + Redis
- Audit immutable + outbox pattern

## Outstanding UI/UX Design System

- Tokens: primary ET Green #0B6E4F, gold #EAB308, neutral zinc, radius lg 16 xl 24, shadow soft, font Inter + NotoSansEthiopic
- Components: Radix + shadcn/ui, glassmorphic nav backdrop-blur-xl, card border 1px rgba(0,0,0,0.06) elevated hover, stepper animated line, camera overlay corner brackets animated, file upload progress slim top
- Motion: Framer Motion 200-300ms ease-out, staggered list 50ms*index, confetti Lottie
- Empty states illustrations Ethiopian coffee ceremony, Axum obelisk subtle
- Lighthouse 90+ target, WCAG AA, keyboard nav, Amharic i18n checkout 100% dashboard 80%

## M5-M9 Roadmap

See `docs/ROADMAP_M5_M9.md` — 22-24 weeks 2 senior eng.

M5 Commercial core: refunds full + 3 rails + smart routing ranking
M6 Recurring: subscriptions dunning + payouts bulk
M7 Workforce: payroll ET tax binary search + pension + payslip PDF outstanding modern
M8 Intelligence: RAG citations mandatory + Swarm planner/critic/executor + Recon breaks
M9 Mobile + Gold: Flutter glass UI + QR FaydaEncode + approvals biometric + FCM + polish

## Testing

```bash
make test       # ledger invariant property 10k iter (via Docker)
make lint       # gosec
k6 run scripts/k6/smoke.js
flutter test

# DB-backed API smoke suite (real Postgres/Redis/MinIO via Docker compose)
docker compose -f deploy/docker/docker-compose.yml \
  -f deploy/docker/docker-compose.integration.yml up \
  --build --abort-on-container-exit --exit-code-from integration integration
```

Go unit tests live under `services/api/internal/**/*_test.go`. The **DB integration tests** live
under `services/api/internal/integration/` behind the `integration` build tag — they compile
against a real Postgres (`DATABASE_URL`) and skip when one is unavailable:

```bash
cd services/api
go test -tags integration ./internal/integration/...   # requires a reachable Postgres
```

*End of README v1.1.0-full*
