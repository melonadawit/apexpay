#!/bin/sh
set -eu
API="http://api:8080"
KEY="sk_test_demo_6c0b88c984e74070b870"
for i in $(seq 1 60); do curl -fsS "$API/healthz" >/dev/null && break; sleep 1; done
curl -fsS "$API/healthz" >/dev/null
# Authentication: absence and wrong secret must both fail; valid test key must pass.
test "$(curl -s -o /dev/null -w '%{http_code}' "$API/v1/methods")" = "401"
test "$(curl -s -o /dev/null -w '%{http_code}' -H 'Authorization: Bearer sk_test_demo_wrong_secret_xxxxxxxxx' "$API/v1/methods")" = "401"
test "$(curl -s -o /tmp/methods.json -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/methods?amount=100&currency=ETB")" = "200"
grep -q 'success' /tmp/methods.json
# Explicit role: seeded role:ops key may access reconciliation operations.
test "$(curl -s -o /tmp/recon.json -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/admin/payment-reconciliation")" = "200"
grep -q 'success' /tmp/recon.json
# Idempotency + persisted 2FA + ledger verification path.
PAYLOAD='{"tx_ref":"docker-smoke-001","amount":"6000.00","currency":"ETB","method":"mock","return_url":"https://merchant.example/return"}'
test "$(curl -s -o /tmp/payment.json -w '%{http_code}' -X POST "$API/v1/transactions/initialize" -H "Authorization: Bearer $KEY" -H 'Content-Type: application/json' -H 'Idempotency-Key: docker-smoke-idem-001' -d "$PAYLOAD")" = "201"
PAYMENT_ID=$(sed -n 's/.*"id":"\([^"]*\)".*/\1/p' /tmp/payment.json | head -1)
test -n "$PAYMENT_ID"
# Same key + same input is safe; changed input must conflict.
test "$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API/v1/transactions/initialize" -H "Authorization: Bearer $KEY" -H 'Content-Type: application/json' -H 'Idempotency-Key: docker-smoke-idem-001' -d "$PAYLOAD")" = "201"
test "$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API/v1/transactions/initialize" -H "Authorization: Bearer $KEY" -H 'Content-Type: application/json' -H 'Idempotency-Key: docker-smoke-idem-001' -d '{"tx_ref":"docker-smoke-001","amount":"6001.00","currency":"ETB","method":"mock"}')" = "409"
test "$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API/v1/transactions/$PAYMENT_ID/2fa/verify" -H "Authorization: Bearer $KEY" -H 'Content-Type: application/json' -d '{"otp":"123456"}')" = "200"
test "$(curl -s -o /tmp/verify.json -w '%{http_code}' -H "Authorization: Bearer $KEY" "$API/v1/transactions/verify/docker-smoke-001")" = "200"
grep -q 'succeeded' /tmp/verify.json
echo 'Docker API smoke suite passed'
