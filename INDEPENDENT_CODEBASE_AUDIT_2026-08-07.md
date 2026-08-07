# Independent ApexPay codebase audit

**Repository reviewed:** `melonadawit/apexpay`, shallow HEAD `45cc1be`  
**Review date:** 7 August 2026  
**Method:** static, repository-wide review of the Go/Python/TypeScript/Dart code, schema, API contract, CI and Docker configuration. Per request, I did **not** install Go, download Go modules, or create a local Go cache. Docker is not available in this workspace (`docker: command not found`), so container builds and runtime tests could not be executed here.

## Executive verdict

ApexPay is an **ambitious Ethiopia-focused payment, merchant onboarding and payroll platform prototype**. It has a thoughtfully structured database and many well-designed domain concepts. However, it is **not currently deployable as a secure production payment platform, nor demonstrably runnable as the documented full stack**.

The repository is strongest as a product specification / UI prototype plus partial backend implementation. Its documentation says “100% Gold” and “all TODOs done,” but this is contradicted by the code and configuration: core authentication does not validate the secret; most request handlers cannot obtain the merchant from context; production rails and Fayda are mocks; critical jobs are ticker-only skeletons; frontend compose images do not have Dockerfiles; and migrations are not run by Compose.

**Recommended release decision:** do not process real money, identity data, payroll disbursements, or regulated customer data yet. Freeze feature expansion, repair the platform vertical slice, and independently verify it in CI before making further feature claims.

---

## What I understand the product to be

ApexPay aims to be an Ethiopian financial platform similar in breadth to a payments processor plus RazorpayX-style business banking/payroll tooling:

- merchant KYC/onboarding, document vault and Fayda identity verification;
- payment initialization, hosted checkout/payment links, refunds, subscriptions, payouts and smart connector routing;
- double-entry ledger and outbox/webhook model;
- Ethiopian payroll (tax, pension, leave, loans, claims, payslips and compliance reports);
- merchant/admin Next.js applications, Flutter mobile client, a Python RAG service and a Go worker;
- PostgreSQL/pgvector, Redis and MinIO, packaged as Docker Compose.

The architectural intent is a modular monolith: `services/api/internal/<domain>` uses domain/service/repository/handler layers; schema evolves through numbered SQL migrations; API code uses Chi + pgx; money mostly uses `shopspring/decimal` / `numeric(20,8)`.

---

## What has been done well

### Product and domain modelling

1. **Clear Ethiopia-specific product direction.** The scope accounts for ETB, local banks, Fayda, pension/payroll requirements, Amharic UI and Addis Ababa time. The effort to make the platform locally relevant is visible throughout the domains and documentation.
2. **Strong initial relational model.** The schema includes useful database constraints, foreign keys, status checks, composite indexes, idempotency/outbox tables, ledger books/journals/entries/balances, webhook tables and progressive migrations (0001–0016).
3. **Correct money direction.** Payment/payout/refund/payroll models generally use decimal values and PostgreSQL `numeric(20,8)`, avoiding normal float use for money. This is the right baseline for a financial system.
4. **Good ledger intent and some real tests.** The ledger has balance validation and tests, including a property-style 10k invariant test. Payroll formula tests are also present.
5. **Good separation of concerns in many domains.** The domain/service/repository/handler arrangement is more maintainable than placing all logic in HTTP controllers. Interfaces for connectors and Fayda make later integration possible.
6. **Security design goals are often sound.** Hashed FIN handling, masked account numbers, MinIO presigned URLs, audit/outbox concepts, rate-limiting intentions, least-privilege non-root final images and avoiding plaintext secrets in the repository are all appropriate goals.
7. **API operational basics exist.** The Go API configures request IDs, recovery, timeouts, compression, health/readiness and Prometheus metrics. The worker is separated from the API process.
8. **High-quality UI effort.** Merchant UI pages/components, reusable UI primitives, an onboarding wizard, payroll screens, docs and Flutter screen structure show significant UX investment.
9. **The repo is small and organized enough to recover.** The shallow checkout is only ~4.1 MB, with a conventional top-level layout. This is repairable if scope is controlled.

---

## Critical blockers — repair before any pilot

### P0-1. API-key authentication accepts an arbitrary token with a known prefix

**Evidence:** `services/api/internal/platform/middleware/auth.go`, lines around 43–60.

The middleware fetches the API-key row by `key_prefix` and checks only `status`. It explicitly comments that verification of `secret_hash` is a placeholder and never compares the bearer secret to a hash. Anybody who learns/guesses a valid 12-character prefix can authenticate as that merchant. This is an account takeover / payment authorization vulnerability.

**Required fix:**
- Store API key material as a random public identifier/prefix plus an Argon2id or bcrypt hash of the complete secret.
- Query candidate key by a non-secret key identifier/prefix, then run constant-time password-hash verification on the supplied complete token.
- Enforce key type/environment/scopes and merchant status.
- Add an authenticated principal to context and test: valid key, wrong secret with valid prefix, revoked key, malformed key, scope denial and cross-tenant denial.
- Rotate/revoke every existing key schema/data before launch, because the previous semantics are insecure.

### P0-2. Tenant identity is not propagated to handlers

**Evidence:** authentication writes context with the private typed key `middleware.CtxMerchantID`, while handlers nearly everywhere read `r.Context().Value("merchant_id")` using a plain string. Examples include payment, onboarding, payout, payroll, subscription, routing, refund, webhook and Fayda handlers.

In Go, those keys are different values/types. The string lookup returns nil. Most handlers silently receive an empty merchant ID; `onboarding/handler.go` additionally performs a direct type assertion and can panic.

This breaks normal requests and creates severe tenant-isolation risk if any repository query handles an empty ID unexpectedly.

**Required fix:** expose accessor functions in the middleware package, e.g. `MerchantIDFromContext(ctx) (string, bool)`, and change **every** handler to use them. Return 401/403 when no principal is present. Add multi-tenant integration tests and make a static lint/search rule reject string context keys.

### P0-3. The documented Docker Compose stack cannot be built/run as provided

**Evidence:**

- `deploy/docker/docker-compose.yml` builds `apps/merchant-web` and `apps/checkout-web` with a `Dockerfile`, but neither directory contains one.
- `apps/checkout-web` lacks `package.json`; it only has `app/...` source files. It has no declared Node dependencies or build command.
- `Dockerfile.api` includes `COPY ../.. ./../..`, a source path outside its declared Compose build context. It is unnecessary and likely rejected/invalid depending on Docker path normalization.
- `Dockerfile.worker` attempts to build from `/app/services/worker`, but the sole `go.mod` lives under `/app/services/api`. Therefore the worker is outside the Go module and its imports of `apexpay/internal/...` cannot be built in that location.
- Compose starts database/API/worker but includes no migration job, and API startup does not apply migrations. A fresh database will not contain required tables.
- Browser containers use `NEXT_PUBLIC_API_URL=http://localhost:8080/v1`. In a user’s browser, `localhost` is the user’s machine, not the Compose API container. The frontend must use a same-origin reverse proxy or a publicly reachable API URL.

**Required fix:** make container startup the first acceptance criterion. Choose either (a) a root Go module with `services/api` and `services/worker` as packages, or (b) separate modules with a correctly versioned shared module; then repair Docker build contexts. Add Dockerfiles and package manifests/lock files for both web apps. Add a one-shot migration service or startup migration policy with a locking strategy. Route browser API calls through Nginx/Next rewrites using relative `/api` paths. In CI run `docker compose build`, then an isolated `docker compose up` smoke test.

### P0-4. Payments/ledger settlement cannot be trusted yet

**Evidence:** `payment/service.go` creates a ledger journal using fixed IDs such as `merchant_operating_default` and account IDs like `asset:clearing:<connector>`. The schema’s `ledger_accounts.id` is an ID, while these values look like account codes, and there is no demonstrated seed creating these IDs. It also only registers the mock connector. 2FA verification simply accepts `123456` and returns JSON; it does not update the payment record.

The link service similarly assigns a payment-link ID as a placeholder payment ID. The refund flow and numerous payroll flows have “demo”, “mock” or TODO behavior. This is unacceptable for financial posting.

**Required fix:** create actual merchant/book/account provisioning and resolve account IDs in one transaction. Make payment state transitions explicit and idempotent. Bind 2FA challenge to payment, user, expiry, purpose and attempt count; persist successful verification before verify/capture. Add real connector adapters only after sandbox contract tests. Reconcile every external result to a unique immutable ledger posting.

---

## High-priority correctness, security and operations gaps

### P1 — authorization and administrative access

- RBAC defaults a missing role to `owner` (`auth.go`). API keys do not load a user role, and the `/v1/admin` route permits `owner`. This makes admin authorization ill-defined.
- User/JWT authentication is advertised in documentation but no actual authentication service/route was found. `JWTSecret` is required but not used for JWT validation in the reviewed API path.
- Many handlers use `user_id` but authentication never sets it.
- Admin review, evidence, merchant exam, connector health and device registration endpoints return hard-coded demo responses instead of performing authorized operations.

**Fix:** model principals separately (user session vs machine API key), establish scopes and memberships, require real role assignment for admin endpoints, and add authorization tests for each resource and tenant.

### P1 — placeholders presented as production features

The code contains many explicit mock/demo/skeleton implementations despite status documents declaring “100% Gold” and zero TODOs:

- live Fayda verifier is a placeholder; mock OTP is `123456`;
- all payment routing falls back to a mock connector; Telebirr/CBE/bank/EthSwitch/card adapters are absent;
- compliance/RAG API returns a fixed legal answer and citation rather than calling the Python RAG service;
- agent chat returns a fixed fake payment-link output;
- customer/list APIs return static objects or empty lists;
- worker ticks only log messages; no outbox publication, webhook delivery, health sampling, dunning, reconciliation, payroll processing or shutdown signal handling occurs;
- payroll disbursement transaction is explicitly TODO and several reports/files use placeholder data;
- mobile offline synchronization is TODO;
- frontend onboarding and several screens show mock values and OTPs.

**Fix:** label the repository honestly as prototype/sandbox until behavior is real. Remove fake production responses or gate them behind an explicit `DEMO_MODE` unavailable in staging/production. Build one complete vertical slice before retaining breadth.

### P1 — test/CI claims do not match coverage

Only two Go test files, one Python test and one integration test directory/file were found for a ~24,422-line multi-service source tree. There are no visible tests for authentication, tenant access, handlers, migrations, connectors, webhooks, 2FA or end-to-end payment flow.

CI is mostly security scanning. It does not run `go test`, frontend build/test, Flutter analyze/test, RAG tests, OpenAPI validation, migrations, Compose build or real Lighthouse/Axe checks. The Lighthouse job prints “Would run.” The Gitleaks job references `.gitleaks.toml`, but that file is absent, so the workflow may fail.

**Fix:** CI should fail fast on formatting/static analysis, unit tests, API integration tests against a disposable Postgres/Redis/MinIO stack, migration-up from empty DB, OpenAPI contract validation, frontend builds/types/tests and container builds. Add coverage thresholds only after meaningful test suites exist.

### P1 — migrations and database hygiene

- Migrations are numerous and use `if not exists`, but no automatic migrator exists in the deployment path.
- `0001_init` has a duplicate/redundant merchant status enum list (`suspended`, `closed`, `active` appear twice), suggesting the schema has not been cleanly validated.
- Audit-log append-only is claimed, but the migration creates the table and indexes only; no immutability trigger was found. The CI check only prints a suggestion and succeeds.
- Down migrations are destructive, contrary to their “forward-only” commentary. They should never be exposed as a normal production Make target.
- Later schema scope is very broad (forex, lending, cards, accounts) without accompanying production-grade service/API implementations. This creates regulatory and maintenance exposure.

**Fix:** run migrations against a blank database in CI; add an actual audit `BEFORE UPDATE OR DELETE` prevention trigger and test it; use a migration tool image in Compose; treat down migrations as local-only; audit/reduce unimplemented regulated tables.

### P1 — external-data and regulated-data readiness

- The README cites specific NBE directives and operational thresholds as facts, but the code/docs do not show validated legal sources, approvals or audit evidence. Do not treat these numbers as legally verified.
- Fayda’s live integration, OIDC/JWT verification and signing are not implemented.
- MinIO server-side encryption is claimed in comments but `storage/minio.go` shows it as a placeholder; bucket/policy/retention/antivirus implementation needs validation.
- Webhook delivery and SSRF prevention are not implemented in the worker.
- No production secret manager, key rotation, backup/restore test, retention/deletion policy, DPIA/privacy process or incident response procedure is evident.

**Fix:** obtain legal/compliance review directly from qualified Ethiopian regulatory counsel and each partner (Fayda/banks/rails). Treat compliance documentation in the repo as a product requirement draft, not evidence of compliance.

---

## Important implementation concerns

1. **OpenAPI drift.** `libs/openapi/openapi.yaml` lists a much smaller/older route set than the API mounts, while several mounted API paths and request/response shapes are placeholders. Generate or validate the spec from the implementation and test it.
2. **Frontend build hygiene.** `apps/merchant-web/tsconfig.tsbuildinfo` is committed, despite being a generated artifact; it contains compiler diagnostics. Add `*.tsbuildinfo` to `.gitignore`, remove the tracked file, and resolve TypeScript errors rather than committing the build cache.
3. **Frontend dependency inconsistency.** The merchant PDF helper references jsPDF-related functionality but the package manifest does not include it. Checkout is not independently packageable.
4. **Health endpoint semantics.** `/readyz` reports Redis and MinIO “ok” regardless of their real availability; it only pings Postgres. Check each dependency or report degraded state accurately.
5. **Docker health check.** API Dockerfile invokes `/api -healthcheck`, but the program has no `-healthcheck` flag. The final image is distroless and the health-check syntax/command should be replaced with a valid health probe strategy.
6. **Unsafe async database use.** API-key last-used updates use detached background goroutines with `context.Background()` per request. This loses request lifetime/cancellation and can create unbounded work under load. Use a bounded queue or update synchronously/batched.
7. **Validation/idempotency gaps.** Payment idempotency saves `TxRef` as a “request hash” and repository lookup is marked placeholder. It must hash a canonical request and reject reuse with a changed body.
8. **No rate-limiting implementation confirmation.** The rate-limit middleware describes a skeleton; apply a tested Redis/token-bucket limiter to public and credential endpoints.
9. **Broad scope, weak delivery focus.** Banking, lending, forex, accounting, cards, RAG and multi-agent functions should not be developed alongside the security-critical initial payment slice.

---

## Documentation assessment

The repository documentation is detailed and helpful for intended product direction. It explains terms, user flows and desired algorithms well. But it is not a reliable engineering status report.

For example, `FINAL_100_GOLD.md` reports 101 done, 0 partial, 0 TODO and “100% Gold”; `GO_100_GOLD_STATUS.md` describes the worker as complete. Code inspection finds explicit TODOs and placeholders (notably live Fayda integration and atomic payroll disbursement), missing frontend Docker/package files, worker logging skeletons and static API responses. The status documents should be replaced by an evidence-based readiness matrix with columns for: implemented, unit tested, integration tested, sandbox certified, security reviewed and production approved.

---

## Recommended next plan

### Phase 0 — stop unsafe claims and establish a baseline (1–3 days)

1. Mark the project **prototype / non-production**; remove 100%-complete claims.
2. Protect the main branch; require CI and review.
3. Remove fake live-looking values and public demo endpoints, or protect them with a development-only demo flag.
4. Inventory every endpoint as real, partial, mocked or absent. Make this the source of truth.

### Phase 1 — make one local stack reproducible (3–7 days)

1. Repair module layout and both Go Dockerfiles; do all module download/build inside Docker as requested.
2. Add missing Dockerfiles and `package.json`/lock file for checkout web; implement a browser-safe reverse proxy/relative API path.
3. Add a migration container/job and a seeded non-sensitive development fixture.
4. Make `docker compose build` and `up` deterministic with health checks that reflect actual state.
5. Add a smoke script that proves health, migration state and one unauthenticated/public route.

### Phase 2 — security and multi-tenancy first (1–2 weeks)

1. Fix API-key validation and key lifecycle; remove default-owner RBAC.
2. Repair typed context access across every handler.
3. Implement user authentication/session model or explicitly defer it; do not leave a dead JWT configuration.
4. Add tenant-isolation, auth-negative, API-key, rate-limit, secret-redaction and audit-immutability tests.
5. Have an independent security review before external users.

### Phase 3 — complete a narrow sandbox payment vertical slice (2–4 weeks)

Limit scope to merchant onboarding (without live Fayda until formally approved), API-key issuance, payment initialization, one certified sandbox connector, webhook receipt/delivery, idempotency, ledger posting, reconciliation and merchant payment visibility.

Every state transition must be persisted, idempotent, authenticated and reconciled. Use real database account IDs/books, not account-code strings. Add an end-to-end Compose test that verifies duplicate requests, webhook retry and ledger balance.

### Phase 4 — payroll and compliance after the payment core

Implement worker jobs rather than tick logs, make payroll disbursement atomic, validate payroll rules with a qualified local tax/payroll expert, build report fixtures, and certify output with pilot employers. Complete Fayda integration only after partner sandbox approval and formal security/privacy review.

### Phase 5 — defer optional modules

Keep RAG, “swarm” agent functionality, lending, forex, cards, accounting sync and broad banking UI behind flags or move them to a future roadmap. They materially increase attack surface and regulatory scope without proving the core payment service.

---

## Practical priority checklist

- [ ] **P0:** authenticate complete API secrets using Argon2id/bcrypt; rotate keys.
- [ ] **P0:** replace all raw-string context access with typed accessors; add tenant tests.
- [ ] **P0:** make Compose buildable; add missing web Dockerfiles/package metadata and migrations.
- [ ] **P0:** make payment 2FA/state/ledger posting real and transactional, not mock responses.
- [ ] **P1:** implement worker jobs or disable their advertised capabilities.
- [ ] **P1:** remove hard-coded production-looking compliance/admin/agent responses.
- [ ] **P1:** establish meaningful CI and integration coverage.
- [ ] **P1:** add actual audit immutability and verified dependency readiness.
- [ ] **P1:** reconcile OpenAPI, implementation, docs and deployment configuration.
- [ ] **P2:** clean generated artifacts (`*.tsbuildinfo`) and resolve frontend compilation diagnostics.

## Bottom line

The project demonstrates strong product vision and meaningful engineering groundwork, particularly its schema, modular organization, money-type choices and Ethiopia-oriented UX. The next success is **not adding more pages or feature domains**. It is proving a small, secure, tenant-safe, Docker-reproducible payment flow with real tests and accurate documentation. Only then should ApexPay expand toward payroll, rails, Fayda, banking and AI features.
