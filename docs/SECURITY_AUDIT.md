# ApexPay — Formal Security Audit

> **Document type:** Security Audit Report
> **Version:** 1.0
> **Status:** Internal / review
> **Audience:** Engineering, security reviewers, technical investors
> **Scope:** `services/api` (Go core), CI security gates, key data-at-rest and identity controls.
> **Method:** Static review of committed code + live verification against a seeded Postgres.

---

## 1. Executive Summary

ApexPay's security posture is strong for the money-critical core: money is never stored as float,
API-key secrets are hashed at rest and verified by prefix+digest, financial identity (FIN) data is
kept to a salted hash plus last-4, webhook endpoints are SSRF-protected, and CI runs `gosec`,
`trivy`, and `gitleaks`. The audit confirmed the money-safety, FIN-privacy, and secret-handling
controls hold up to live inspection.

Two findings are worth addressing before production hardening (see §5):
- **Audit trail is not actually append-only.** The SAD and the marketing deck claim `audit_logs` is
  "guarded by a DB trigger that prevents UPDATE/DELETE", but **no such trigger exists** in the
  migrations or in the live schema. The CI `audit-append-only` job is a soft check that only echoes a
  reminder — it never fails the build.
- **No application-level secrets are baked into images**, but dev credentials live in
  `docker-compose.yml` with clearly-labelled placeholder values. This is fine for local dev; the
  production path must inject these via environment/secret managers only (never bake into the image).

Both are low-cost to remediate and are the recommended next actions.

| Area | Rating |
|---|---|
| Money safety (no float, double-entry validation) | **Strong** — verified |
| Secret handling at rest (API keys, connector keys) | **Strong** — verified |
| Financial identity privacy (FIN last-4/hash) | **Strong** — verified |
| Webhook / network SSRF protection | **Strong** — verified |
| CI security gates (gosec/trivy/gitleaks) | **Strong** |
| Audit-trail integrity (append-only) | **Needs attention** — gap found |
| Secrets in source/config | **Good** — no real secrets; dev placeholders only |

---

## 2. Scope & Method

**In scope:** authentication/authorization, money representation, secret handling, financial-identity
privacy, webhook network security, audit/append-only integrity, CI gates.

**Method performed in this audit:**
1. Static grep/`gosec`-style review of the money core (`payment`, `ledger`, `refund`, `payout`,
   `payroll`, `worker`).
2. Live queries against a seeded Postgres to confirm schema-level claims (triggers, tables).
3. Reproduction of each CI security gate locally where possible.

---

## 3. Verified Controls (with evidence)

### 3.1 Money is always `decimal.Decimal` — verified
```
$ grep -R "float64.*amount\|amount.*float64" ... payment ledger refund payout payroll worker
# → no matches  (PASS)
```
All money columns are `numeric(20,8)`; every journal must balance
(`ledger.ValidateBalanced`) before insert, and posting is idempotent via
`(book_id, posting_key)` uniqueness.

### 3.2 API-key secrets hashed at rest, prefix+digest verification — verified
- `api_keys.secret_hash = sha256(full token)`; only the first 12 chars (`key_prefix`) are stored
  plaintext and used for O(1) lookup.
- Auth middleware re-hashes the presented bearer and compares to `secret_hash` (§`platform/middleware/auth.go`).
- New keys created through the developer portal follow the same scheme and were verified to work as
  Bearer tokens (HTTP 200 on `/v1/methods`).

### 3.3 Financial identity (FIN) privacy — verified
```
$ grep -RE "l.(Debug|Info|Warn|Error|Fatal)([^)]*FIN|fmt.(Print|Printf|Sprintf)([^)]*FIN" fayda
# → no matches excluding fin_last4 / fin_hash  (PASS)
```
FIN is stored as a salted hash (`sha256(salt+fin)`); only the last-4 is ever returned/logged.

### 3.4 Webhook SSRF protection — verified
`POST /webhooks/endpoints` requires an `https://` URL, resolves the host, and rejects private /
loopback IP ranges before registering an endpoint. Webhook secrets are AES-GCM encrypted at rest and
require ≥16 characters.

### 3.5 CI security gates — present and non-failing
- `gosec` (severity high, confidence high) — hard fail on any high finding.
- `trivy` (FS scan, HIGH/CRITICAL) — SARIF uploaded to the security tab.
- `gitleaks` (with `.gitleaks.toml`) — hard fail on leaked secrets.
- `no-float-money-lint`, `fin-privacy-lint`, `sql-param-cast-lint` — hard fail.

### 3.6 No real secrets in tracked source/config — verified
A high-signal scan of `services/`, `apps/` found no private keys, `sk_live_` secrets, or real
passwords. Compose only carries clearly-labelled **dev placeholders**
(`dev_jwt_secret_change_in_prod`, `0123456789abcdef0123456789abcdef`).

---

## 4. Findings

### F-1 (Medium) — `audit_logs` is not truly append-only
**Claim vs reality:** The SAD (§5.3, §6.3) and investor deck state audit logs are "guarded by a DB
trigger that prevents UPDATE/DELETE". Verification shows **no trigger exists** on `audit_logs`
(either in migrations or live schema), and the CI `audit-append-only` job only echoes a reminder and
never fails the build.

**Risk:** A compromised or buggy writer could update/delete historical audit rows without detection,
undermining the regulatory/audit story that the platform markets.

**Evidence:**
```
SELECT tgname FROM pg_trigger WHERE tgrelid='audit_logs'::regclass AND NOT tgisinternal;
# → (no rows)
```

### F-2 (Low) — Dev credentials in `docker-compose.yml` are visible
`JWT_SECRET` and `CONNECTOR_ENCRYPTION_KEY` are hardcoded dev placeholders in the committed compose
file. Acceptable for local dev, but they must never be used in production and should be injected via
secrets management in any real deployment.

**Evidence:** `deploy/docker/docker-compose.yml` lines 61–62, 77.

### F-3 (Info) — No explicit rate limiting on login observed
Rate limiting middleware exists for the API generally, but the audit did not find a dedicated
credential-stuffing guard (e.g., per-IP login throttle) beyond the general middleware. Confirm or
add login throttling for production.

---

## 5. Recommendations (prioritized)

1. **Add a real append-only trigger on `audit_logs`** and make the CI check fail on absence.
   - Migration:
     ```sql
     CREATE OR REPLACE FUNCTION prevent_audit_mutation() RETURNS trigger AS $$
     BEGIN RAISE EXCEPTION 'audit_logs is append-only'; END $$ LANGUAGE plpgsql;
     CREATE TRIGGER audit_logs_no_mutation
       BEFORE UPDATE OR DELETE ON audit_logs
       FOR EACH ROW EXECUTE FUNCTION prevent_audit_mutation();
     ```
   - Change the CI `audit-append-only` job to `grep -q` (fail on missing trigger) rather than `echo`.
   - Update the SAD wording to match reality once applied.

2. **Keep dev secrets out of any production build path.** Confirm production deploys inject
   `JWT_SECRET`, `CONNECTOR_ENCRYPTION_KEY`, DB creds, and MinIO creds via a secret manager only.

3. **Add login throttling / lockout** (per-account + per-IP) and re-confirm the general rate limiter
   covers auth endpoints.

4. **Add an integration test that asserts the audit trigger** (attempt an UPDATE/DELETE in a test and
   expect a rejection) so the guarantee can't silently regress.

---

## 6. Sign-off

Verified against commit `HEAD` on `master`. The money-safety, FIN-privacy, and secret-handling
controls passed; F-1 is the key remediation to close the gap between the marketed audit guarantee
and the actual enforcement.
