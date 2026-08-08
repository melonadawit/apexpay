# ApexPay — Stack Runbook

How to stand up the full ApexPay platform (API + worker + RAG + merchant dashboard +
hosted checkout) on your own machine, verify it end-to-end, and run the tests.

> **TL;DR for Docker users**
> ```bash
> cp .env.example .env          # optional: set secrets
> docker compose -f deploy/docker/docker-compose.yml up -d --build
> docker compose -f deploy/docker/docker-compose.yml -f deploy/docker/docker-compose.integration.yml up --abort-on-container-exit --exit-code-from integration integration
> ```
> Then open http://localhost:3000 (merchant dashboard) and http://localhost:3001 (checkout).

---

## 1. Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| Docker + Compose v2 | recent | Required for the full stack (Postgres, Redis, MinIO, API, worker, web apps) |
| Node.js | >= 18 | To build the Next.js apps (`merchant-web`, `checkout-web`) |
| Go | 1.22 | Only needed to build/run the API outside Docker or run Go tests locally |
| Python | 3.10+ | Only for the RAG service (`services/rag`) |
| Flutter | 3.22+ | Only for the mobile app (`apps/mobile`) — **optional** |

The fastest, most reliable path is **Docker Compose** — it handles everything including
migrations. The Go/Python/Flutter toolchains are only needed if you run services outside
Docker or run local tests.

---

## 2. What you're standing up

```
                        ┌────────────────────────────┐
   Browser ───► :3000  │ merchant-web (Next.js 14)   │  dashboard (login + data)
                        └──────────┬─────────────────┘      │ /api/proxy (server-side)
                                   │ APEXPAY_API_URL=api:8080│  injects session cookie / API key
                        ┌──────────▼─────────────────┐
   Browser ───► :3001  │ checkout-web (Next.js 14)   │  hosted payment page
                        └──────────┬─────────────────┘      │ /api/proxy (server-side)
                                   │ APEXPAY_API_URL=api:8080
                        ┌──────────▼─────────────────┐
                        │ api (Go) :8080             │◄── worker (Go, background jobs)
                        │  + admin + checkout + auth │
                        └──┬──────┬──────┬──────┬────┘
                           │      │      │      │
                     postgres  redis  minio  (rag :8001)
```

- **api** — the Go HTTP API (chi router, pgx, Redis, MinIO, Prometheus `/metrics`).
- **worker** — background jobs: outbox→webhooks, idempotency reconciliation, forex rate
  cache, credit scoring, accounting sync, FCM notifications, connector health, dunning.
- **rag** — Python FastAPI compliance/RAG worker (`/v1/compliance/ask`), used by
  `/v1/compliance/ask`.
- **merchant-web** — merchant dashboard (real API-backed after the wiring work).
- **checkout-web** — hosted checkout (real API-backed).
- **mobile** — Flutter app (optional, not part of compose).

---

## 3. Quick start (Docker Compose)

### 3.1 Create your env file

```bash
cp deploy/docker/.env.example .env
```

Fill in at least:
- `CONNECTOR_ENCRYPTION_KEY` — 32+ char random key (used to derive per-purpose keys).
- `JWT_SECRET` — random secret.
- `APEXPAY_MERCHANT_API_KEY` — a merchant **test** API key for the dashboard (optional once
  login is used; the dashboard falls back to this when no session exists).

Generate random values:
```bash
openssl rand -hex 32   # connector key
openssl rand -hex 32   # jwt secret
```

### 3.2 Build and start

```bash
docker compose -f deploy/docker/docker-compose.yml up -d --build
```

Wait for health, then verify:

```bash
curl -s http://localhost:8080/healthz
curl -s http://localhost:8080/readyz
```

Expected:
- `healthz` → `{"status":"ok","version":"1.1.0-full",...}`
- `readyz` → `{"status":"ready",...}` (DB reachable)

The `migrate` service applies `db/migrations/*.up.sql` in order automatically on startup.

### 3.3 Run the DB-backed smoke suite (validates everything end-to-end)

```bash
docker compose -f deploy/docker/docker-compose.yml \
  -f deploy/docker/docker-compose.integration.yml up \
  --build --abort-on-container-exit --exit-code-from integration integration
```

This seeds test data and runs `scripts/integration/api_smoke.sh`, which asserts:
1. Authn/authz (missing + wrong API key → 401; valid → 200)
2. RBAC admin routes
3. Admin onboarding queue / merchant exam / recon breaks / evidence / connectors health
4. Admin review action (maker-checker)
5. **Hosted checkout** (public token → load → initialize → requires-2FA → poll)
6. Idempotency + persisted 2FA + ledger verification
7. **Dashboard session auth** (login → me → logout → revoked)

On success it prints `Docker API smoke suite passed` and exits 0.

---

## 4. Logging into the merchant dashboard

1. Open http://localhost:3000
2. You'll be redirected to `/login`.
3. Use the seeded demo account (from `tests/integration/seed.sql`):
   - **Email:** `demo@apexpay.et`
   - **Password:** `Admin@12345`
4. On success the browser receives an **httpOnly** `apexpay_session` cookie; the dashboard
   loads live data (TPV, success rate, active links, recent payments).

> If you didn't run the integration seed, create a user + membership manually (see §7).

---

## 5. Running services outside Docker (for development)

### 5.1 API only

```bash
cd services/api
export DATABASE_URL="postgres://apexpay:apexpay_dev@localhost:5432/apexpay?sslmode=disable"
export REDIS_URL="localhost:6379"
export MINIO_ENDPOINT="localhost:9000" MINIO_ACCESS_KEY=minioadmin MINIO_SECRET_KEY=minioadmin
export CONNECTOR_ENCRYPTION_KEY="$(openssl rand -hex 32)"
export JWT_SECRET="$(openssl rand -hex 32)"
go run ./cmd/api        # :8080
```

### 5.2 Worker

```bash
cd services/api
export DATABASE_URL=... REDIS_URL=... CONNECTOR_ENCRYPTION_KEY=...
go run ./cmd/worker     # background jobs
```

### 5.3 Next.js apps (point the proxy at localhost API)

```bash
cd apps/merchant-web
APEXPAY_API_URL=http://localhost:8080 npm run dev        # :3000

cd apps/checkout-web
APEXPAY_API_URL=http://localhost:8080 npm run dev -- -p 3001
```

The browser still talks to the Next proxy (`/api/proxy`), so no CORS issues — only the
server needs to reach the Go API.

### 5.4 RAG service

```bash
cd services/rag
python -m venv .venv && . .venv/bin/activate
pip install -r requirements.txt
DATABASE_URL=... uvicorn app.api:app --port 8001
```

Set `RAG_SERVICE_URL=http://localhost:8001` on the API. `/v1/compliance/ask` returns **503**
(not a fake answer) when the RAG service is unreachable.

---

## 6. Running the tests

### 6.1 Go unit tests (fast, no DB)

```bash
cd services/api
go test ./... -count=1          # ledger, payment, payroll, crypto, middleware, webhook, reconciliation
```

### 6.2 Go DB-integration tests (need a reachable Postgres)

```bash
cd services/api
# require a Postgres at the DSN below (or set DATABASE_URL). Skips cleanly if absent.
go test -tags integration ./internal/integration/...
```

These cover the onboarding→Fayda→payment→ledger→webhook chain and the admin
review (maker-checker) + auth login flows.

### 6.3 TypeScript parse/syntax check (merchant-web)

```bash
cd apps/merchant-web && npm ci && npm run build
```

### 6.4 Full end-to-end (Docker)

See §3.3. This is the highest-value check — it exercises migrations + workers + all wired
API flows + seed data in one shot.

### 6.5 Performance smoke

```bash
k6 run scripts/k6/smoke.js
k6 run scripts/k6/payroll_comprehensive.js
```

---

## 7. Creating your own dashboard user (without the seed)

```sql
-- Insert a user with an argon2id hash. Generate the hash with the app itself:
--   (cd services/api && go run ./scripts/genhash "YourPassword")
INSERT INTO users (id, email, name, status, password_hash, email_verified)
VALUES ('user_x', 'you@company.et', 'You', 'active', '<hash>', true);

INSERT INTO merchant_members (merchant_id, user_id, role)
SELECT id, 'user_x', 'owner' FROM merchants WHERE legal_name='Your Merchant';
```

Then log in at http://localhost:3000/login.

---

## 8. Configuration reference

| Env | Service | Default | Purpose |
|-----|---------|---------|---------|
| `DATABASE_URL` | api, worker | — | Postgres DSN |
| `REDIS_URL` | api, worker | — | Redis addr |
| `CONNECTOR_ENCRYPTION_KEY` | api, worker | — | Master key (derived per-purpose) |
| `JWT_SECRET` | api | — | Reserved for future JWT features |
| `FAYDA_MODE` | api | `mock` | `mock` or `live` Fayda verification |
| `FAYDA_PARTNER_CODE/KEY` | api | `APEXPAY_TEST` | Fayda partner creds |
| `RAG_SERVICE_URL` | api | `http://rag:8001` | RAG/compliance worker |
| `APEXPAY_API_URL` | merchant-web, checkout-web | `http://api:8080` | Go API base (server-side) |
| `APEXPAY_API_KEY` | merchant-web | empty | Fallback merchant key when no session |
| `MINIO_*` | api | minioadmin/minioadmin | Object storage (KYC docs) |

---

## 9. Troubleshooting

**`readyz` fails (DB not ready)** — the `migrate` service may still be running. Check
`docker compose logs migrate`. Migration errors are printed with the failing file.

**Login works but dashboard is empty** — the merchant has no payments/links yet. Create a
payment link on the Links page or seed data; the dashboard reads live from `/v1/dashboard`.

**`/v1/compliance/ask` returns 503** — the RAG service is down/not running. Start it or set
`RAG_SERVICE_URL` to a reachable instance. The API deliberately does not fabricate answers.

**Dashboard shows "gateway_unreachable"** — `merchant-web` can't reach the Go API. Ensure
`APEXPAY_API_URL` is correct from inside the app's network (`http://api:8080` in compose).

**Go tests can't find DB** — the integration tests skip when `DATABASE_URL` isn't reachable;
that's expected. To run them, point `DATABASE_URL` at a Postgres.

**Ports already in use** — adjust the `ports:` mapping in `deploy/docker/docker-compose.yml`.

---

## 10. Production hardening (before real money)

1. **Secrets**: move `CONNECTOR_ENCRYPTION_KEY`, `JWT_SECRET`, DB/Redis/MinIO creds into a
   real secret manager; never commit `.env`.
2. **HTTPS + cookies**: set `Secure` on the `apexpay_session` cookie (already done in
   production builds) and serve behind TLS.
3. **RAG live**: connect a real embedding/vector store; the Go client is a thin proxy.
4. **Connectors**: only `mock` is registered today. Add Telebirr/CBE/Amole/EthSwitch adapters
   behind the existing `connector.Connector` interface.
5. **2FA provider**: the demo OTP is local-only; wire a real challenge provider before live.
6. **Rate limiting / auth**: apply stricter login rate limiting and add session rotation.
7. **CI is live**: the `security.yml` workflow runs `go build/vet/test`, integration compile
   vet, gosec, trivy, gitleaks, and a docker-smoke job on every push to `master`.

---

## 11. Repository map (where things live)

```
services/api/cmd/api        HTTP API entrypoint
services/api/cmd/worker     background worker entrypoint
services/api/internal/      domain packages (admin, auth, checkout, payment, payroll, ...)
services/api/internal/integration  DB integration tests (build tag: integration)
db/migrations/              0025 sequential SQL migrations
services/rag/               Python RAG/compliance service
apps/merchant-web/          Next.js merchant dashboard (+ /api/proxy, /api/auth)
apps/checkout-web/          Next.js hosted checkout (+ /api/proxy)
apps/mobile/                Flutter app (optional)
deploy/docker/              compose files + Dockerfiles
tests/integration/seed.sql  smoke-test seed data
scripts/integration/api_smoke.sh  end-to-end smoke assertions
scripts/k6/                 load tests
```
