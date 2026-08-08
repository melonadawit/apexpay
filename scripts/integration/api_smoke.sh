#!/bin/sh
set -eu
API="http://api:8080"
KEY="sk_test_demo_6c0b88c984e74070b870"

wait_ready() {
  for i in $(seq 1 60); do
    curl -fsS "$API/healthz" >/dev/null 2>&1 && return 0
    sleep 1
  done
  echo "API did not become ready" >&2
  exit 1
}
wait_ready
curl -fsS "$API/healthz" >/dev/null

# ---- 1. Authn/authz: absence and wrong secret must fail; valid test key must pass. ----
test "$(curl -s -o /dev/null -w '%{http_code}' "$API/v1/methods")" = "401"
test "$(curl -s -o /dev/null -w '%{http_code}' -H 'Authorization: Bearer sk_test_demo_wrong_secret_xxxxxxxxx' "$API/v1/methods")" = "401"
test "$(curl -s -o /tmp/methods.json -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/methods?amount=100&currency=ETB")" = "200"
grep -q 'success' /tmp/methods.json

# ---- 2. RBAC: seeded role:ops key reaches admin reconciliation operations. ----
test "$(curl -s -o /tmp/recon.json -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/admin/payment-reconciliation")" = "200"
grep -q 'success' /tmp/recon.json

# ---- 3. Admin: onboarding review queue returns the seeded reviewable merchant. ----
test "$(curl -s -o /tmp/queue.json -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/admin/onboarding/queue")" = "200"
grep -q 'mer_docker_smoke' /tmp/queue.json

# ---- 4. Admin: merchant compliance exam (PII-safe) returns seeded data. ----
test "$(curl -s -o /tmp/exam.json -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/admin/merchants/mer_docker_smoke/exam")" = "200"
grep -q 'kyc_docker_smoke' /tmp/exam.json
grep -q 'tin_validation' /tmp/exam.json

# ---- 5. Admin: open reconciliation breaks. ----
test "$(curl -s -o /tmp/breaks.json -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/admin/recon/breaks")" = "200"
grep -q 'txr_recon_smoke' /tmp/breaks.json

# ---- 6. Admin: compliance evidence by tx_ref. ----
test "$(curl -s -o /tmp/evidence.json -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/admin/evidence?tx_ref=txr_recon_smoke")" = "200"
grep -q 'txr_recon_smoke' /tmp/evidence.json

# ---- 7. Admin: connector health (may be empty in a fresh run — still 200). ----
test "$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/admin/connectors/health")" = "200"

# ---- 8. Admin: onboarding review action persists a real review record. ----
test "$(curl -s -o /tmp/review.json -w '%{http_code}' -X POST "$API/v1/admin/onboarding/mer_docker_smoke/review" \
  -H "Authorization: Bearer $KEY" -H 'Content-Type: application/json' \
  -d '{"action":"request_info","comment":"please update business license"}')" = "200"
grep -q 'request_info\|needs_more_info' /tmp/review.json

# ---- 9. Hosted checkout (public, token-as-capability): load, initialize, poll, 2FA. ----
test "$(curl -s -o /tmp/link.json -w '%{http_code}' "$API/v1/checkout/tkn_docker_smoke")" = "200"
grep -q 'tkn_docker_smoke' /tmp/link.json

test "$(curl -s -o /tmp/chk_init.json -w '%{http_code}' -X POST "$API/v1/checkout/tkn_docker_smoke/initialize" \
  -H 'Content-Type: application/json' -d '{"method":"mock","customer_email":"cust@example.et"}')" = "201"
CHK_TX=$(sed -n 's/.*"tx_ref":"\([^"]*\)".*/\1/p' /tmp/chk_init.json | head -1)
test -n "$CHK_TX"
grep -q '"requires_2fa":true' /tmp/chk_init.json

test "$(curl -s -o /tmp/chk_status.json -w '%{http_code}' "$API/v1/checkout/tkn_docker_smoke/status/$CHK_TX")" = "200"
grep -q 'succeeded\|pending' /tmp/chk_status.json

# ---- 10. Idempotency + persisted 2FA + ledger verification path (merchant API). ----
PAYLOAD='{"tx_ref":"docker-smoke-001","amount":"6000.00","currency":"ETB","method":"mock","return_url":"https://merchant.example/return"}'
test "$(curl -s -o /tmp/payment.json -w '%{http_code}' -X POST "$API/v1/transactions/initialize" \
  -H "Authorization: Bearer $KEY" -H 'Content-Type: application/json' -H 'Idempotency-Key: docker-smoke-idem-001' -d "$PAYLOAD")" = "201"
PAYMENT_ID=$(sed -n 's/.*"id":"\([^"]*\)".*/\1/p' /tmp/payment.json | head -1)
test -n "$PAYMENT_ID"
# Same key + same input is safe; changed input must conflict.
test "$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API/v1/transactions/initialize" \
  -H "Authorization: Bearer $KEY" -H 'Content-Type: application/json' -H 'Idempotency-Key: docker-smoke-idem-001' -d "$PAYLOAD")" = "201"
test "$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API/v1/transactions/initialize" \
  -H "Authorization: Bearer $KEY" -H 'Content-Type: application/json' -H 'Idempotency-Key: docker-smoke-idem-001' \
  -d '{"tx_ref":"docker-smoke-001","amount":"6001.00","currency":"ETB","method":"mock"}')" = "409"
test "$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API/v1/transactions/$PAYMENT_ID/2fa/verify" \
  -H "Authorization: Bearer $KEY" -H 'Content-Type: application/json' -d '{"otp":"123456"}')" = "200"
test "$(curl -s -o /tmp/verify.json -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/transactions/verify/docker-smoke-001")" = "200"
grep -q 'succeeded' /tmp/verify.json

# ---- 11. Dashboard session auth: login → me → logout. ----
test "$(curl -s -o /tmp/auth_fail.json -w '%{http_code}' -X POST "$API/v1/auth/login" \
  -H 'Content-Type: application/json' -d '{"email":"demo@apexpay.et","password":"wrong-password"}')" = "401"
test "$(curl -s -o /tmp/auth.json -w '%{http_code}' -X POST "$API/v1/auth/login" \
  -H 'Content-Type: application/json' -d '{"email":"demo@apexpay.et","password":"Admin@12345"}')" = "200"
SESSION_TOKEN=$(sed -n 's/.*"token":"\([^"]*\)".*/\1/p' /tmp/auth.json | head -1)
test -n "$SESSION_TOKEN"
test "$(curl -s -o /tmp/auth_me.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/auth/me")" = "200"
grep -q 'demo@apexpay.et' /tmp/auth_me.json
test "$(curl -s -o /tmp/auth_logout.json -w '%{http_code}' -X POST -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/auth/logout")" = "200"
# Revoked session must now be rejected.
test "$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/auth/me")" = "401"

echo 'Docker API smoke suite passed'
