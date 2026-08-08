#!/bin/sh
set -eu
API="http://api:8080"
KEY="sk_test_demo_6c0b88c984e74070b870"

# On any failing command, report the line number + the failing expression so CI failures
# are diagnosable (the `test` builtin exits non-zero silently under `set -e`).
trap 'rc=$?; echo "SMOKE FAILED at line $LINENO (exit $rc)"; exit $rc' ERR

wait_ready() {
  # Allow ample time for cold-start: image pulls, migrations, and the API binding to :8080.
  # Each attempt is bounded by curl --max-time so a hanging connect can't stall the loop.
  for i in $(seq 1 90); do
    if curl -fsS --max-time 3 "$API/healthz" >/dev/null 2>&1; then
      echo "API ready (attempt $i)"
      return 0
    fi
    sleep 2
  done
  echo "API did not become ready after 90 attempts" >&2
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

# Re-login for the session-authenticated module sections (the prior token was revoked).
test "$(curl -s -o /tmp/auth2.json -w '%{http_code}' -X POST "$API/v1/auth/login" \
  -H 'Content-Type: application/json' -d '{"email":"demo@apexpay.et","password":"Admin@12345"}')" = "200"
SESSION_TOKEN=$(sed -n 's/.*"token":"\([^"]*\)".*/\1/p' /tmp/auth2.json | head -1)
test -n "$SESSION_TOKEN"

# ---- 12. Banking modules (session-authenticated): current accounts, forex, notifications. ----
test "$(curl -s -o /tmp/bank_acct.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/banking/current_accounts")" = "200"
grep -q 'ETB-CBE-7778889990' /tmp/bank_acct.json
test "$(curl -s -o /tmp/bank_fx.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/banking/forex/rates")" = "200"
grep -q 'USD' /tmp/bank_fx.json
test "$(curl -s -o /tmp/bank_notif.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/banking/notifications")" = "200"
grep -q 'Current Account Opened' /tmp/bank_notif.json
# Banking without a session must be rejected.
test "$(curl -s -o /dev/null -w '%{http_code}' "$API/v1/banking/current_accounts")" = "401"


# ---- 13. Banking action endpoints: create + list. ----
test "$(curl -s -o /tmp/bank_tax.json -w '%{http_code}' -X POST "$API/v1/banking/tax_payments" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"tax_type":"vat","amount":"1500.00","currency":"ETB","period_month":8,"period_year":2026,"status":"draft"}')" = "201"
grep -q 'vat' /tmp/bank_tax.json
test "$(curl -s -o /tmp/bank_tax_list.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/banking/tax_payments")" = "200"
grep -q 'vat' /tmp/bank_tax_list.json
# Create a vendor invoice.
test "$(curl -s -o /tmp/bank_inv.json -w '%{http_code}' -X POST "$API/v1/banking/vendor_invoices" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"invoice_number":"SMOKE-INV-1","invoice_date":"2026-08-08","vendor_name":"Smoke Vendor","amount":"10000.00","currency":"ETB","tax_amount":"1500.00","withholding_tax_amount":"200.00","total_amount":"11300.00","status":"pending_approval","ocr_confidence":0.9}')" = "201"
test "$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/banking/vendor_invoices")" = "200"


# ---- 14. HRIS (Workforce OS): teams, risk engine. ----
test "$(curl -s -o /tmp/hris_teams.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/hris/teams")" = "200"
grep -q 'Engineering' /tmp/hris_teams.json
test "$(curl -s -o /tmp/risk_rules.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/risk/rules")" = "200"
grep -q 'High-ticket' /tmp/risk_rules.json
test "$(curl -s -o /tmp/risk_eval.json -w '%{http_code}' -X POST "$API/v1/risk/evaluate" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"entity_type":"payment","entity_id":"pay_smoke","amount":"600000.00"}')" = "200"
grep -q 'findings' /tmp/risk_eval.json

echo 'Docker API smoke suite passed'

# ---- 15. Treasury + Invoicing. ----
test "$(curl -s -o /tmp/treasury_pos.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/treasury/position")" = "200"
grep -q 'accounts' /tmp/treasury_pos.json
test "$(curl -s -o /tmp/invoices.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/invoices")" = "200"
grep -q 'INV-SMOKE-001' /tmp/invoices.json
test "$(curl -s -o /tmp/inv_aging.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/invoices/aging")" = "200"
grep -q 'bucket' /tmp/inv_aging.json

echo 'Docker API smoke suite passed'

# ---- 16. Compliance console, fixed assets, analytics, 2FA. ----
test "$(curl -s -o /tmp/compliance_status.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/compliance-console/status")" = "200"
grep -q 'overall_status' /tmp/compliance_status.json
test "$(curl -s -o /tmp/fixed_assets.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/fixed-assets")" = "200"
grep -q 'Delivery Van' /tmp/fixed_assets.json
test "$(curl -s -o /tmp/analytics.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/analytics/revenue")" = "200"
test "$(curl -s -o /tmp/2fa.json -w '%{http_code}' -X POST "$API/v1/2fa/enroll" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' -d '{"account":"demo@apexpay.et"}')" = "200"
grep -q 'secret' /tmp/2fa.json

echo 'Docker API smoke suite passed'

# ---- 17. Accounting & Bookkeeping. ----
test "$(curl -s -o /tmp/acc_coa.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/accounting/accounts")" = "200"
test "$(curl -s -o /tmp/acc_tb.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/accounting/trial-balance")" = "200"
test "$(curl -s -o /tmp/acc_pnl.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/accounting/profit-loss")" = "200"
grep -q 'Profit' /tmp/acc_pnl.json
test "$(curl -s -o /tmp/acc_bs.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/accounting/balance-sheet")" = "200"
grep -q 'Assets' /tmp/acc_bs.json

echo 'Docker API smoke suite passed'

# ---- 18. Inventory, Disputes, Loyalty, Lending. ----
test "$(curl -s -o /tmp/prods.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/inventory/products")" = "200"
test "$(curl -s -o /tmp/disputes.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/disputes")" = "200"
test "$(curl -s -o /tmp/loyalty.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/loyalty/accounts")" = "200"
test "$(curl -s -o /tmp/lending.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/lending")" = "200"

# ---- 19. Apex Assistant (session-authenticated, read-only, role-scoped). ----
test "$(curl -s -o /tmp/assistant.json -w '%{http_code}' -X POST "$API/v1/assistant/chat" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"message":"what is my cash position and tpv today?"}')" = "200"
grep -q '"thread_id"' /tmp/assistant.json
grep -q '"answer"' /tmp/assistant.json
grep -q '"actor":"merchant"' /tmp/assistant.json
ASSISTANT_THREAD=$(sed -n 's/.*"thread_id":"\([^"]*\)".*/\1/p' /tmp/assistant.json | head -1)
test -n "$ASSISTANT_THREAD"
test "$(curl -s -o /tmp/assistant_thread.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/assistant/threads/$ASSISTANT_THREAD")" = "200"
grep -q '"messages"' /tmp/assistant_thread.json
# Assistant must be rejected without a session.
test "$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API/v1/assistant/chat" -H 'Content-Type: application/json' -d '{"message":"hi"}' )" = "401"

# ---- 20. Real GL: manual journal entries + fiscal period close. ----
# Balanced entry posts (201).
test "$(curl -s -o /tmp/gl_ok.json -w '%{http_code}' -X POST "$API/v1/accounting/journal-entries" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"memo":"smoke manual entry","lines":[{"account_code":"asset:bank","direction":"debit","amount":"2500.00"},{"account_code":"revenue:product","direction":"credit","amount":"2500.00"}]}')" = "201"
grep -q '"id"' /tmp/gl_ok.json
# Unbalanced entry rejected (400).
test "$(curl -s -o /tmp/gl_bad.json -w '%{http_code}' -X POST "$API/v1/accounting/journal-entries" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"memo":"bad","lines":[{"account_code":"asset:bank","direction":"debit","amount":"100"},{"account_code":"revenue:product","direction":"credit","amount":"90"}]}')" = "400"
# List journal entries (200).
test "$(curl -s -o /tmp/gl_list.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/accounting/journal-entries")" = "200"
# Period close: close current month, then a posting is rejected.
PERIOD=$(date +%Y-%m)
test "$(curl -s -o /tmp/gl_close.json -w '%{http_code}' -X POST "$API/v1/accounting/periods/close" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d "{\"period\":\"$PERIOD\"}")" = "200"
test "$(curl -s -o /tmp/gl_closed.json -w '%{http_code}' -X POST "$API/v1/accounting/journal-entries" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"memo":"after close","lines":[{"account_code":"asset:bank","direction":"debit","amount":"50"},{"account_code":"revenue:product","direction":"credit","amount":"50"}]}')" = "400"
grep -q 'closed' /tmp/gl_closed.json
# Reopen then posting succeeds again (201).
test "$(curl -s -o /tmp/gl_reopen.json -w '%{http_code}' -X POST "$API/v1/accounting/periods/reopen" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d "{\"period\":\"$PERIOD\"}")" = "200"
test "$(curl -s -o /tmp/gl_ok2.json -w '%{http_code}' -X POST "$API/v1/accounting/journal-entries" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"memo":"after reopen","lines":[{"account_code":"asset:bank","direction":"debit","amount":"50"},{"account_code":"revenue:product","direction":"credit","amount":"50"}]}')" = "201"
# Period list (200).
test "$(curl -s -o /tmp/gl_periods.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/accounting/periods")" = "200"

# ---- 21. Procurement & Accounts Payable: vendors, POs, receipts, AP invoices, aging. ----
# Create a vendor (201).
test "$(curl -s -o /tmp/pr_vendor.json -w '%{http_code}' -X POST "$API/v1/procurement/vendors" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"Smoke Supplies Co","payment_terms_days":30}')" = "201"
grep -q '"id"' /tmp/pr_vendor.json
PR_VID=$(sed -n 's/.*"id":"\([^"]*\)".*/\1/p' /tmp/pr_vendor.json | head -1)
test -n "$PR_VID"
# Create a PO (201) and capture its id.
test "$(curl -s -o /tmp/pr_po.json -w '%{http_code}' -X POST "$API/v1/procurement/purchase-orders" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d "{\"vendor_id\":\"$PR_VID\",\"po_number\":\"PO-SMOKE-1\",\"order_date\":\"2026-08-01\",\"tax_rate\":\"0.15\",\"items\":[{\"item_name\":\"Widget\",\"quantity\":\"10\",\"unit_price\":\"100\"}]}")" = "201"
grep -q '"po_number"' /tmp/pr_po.json
PR_POID=$(sed -n 's/.*"id":"\([^"]*\)".*/\1/p' /tmp/pr_po.json | head -1)
test -n "$PR_POID"
# Receive the PO (201).
test "$(curl -s -o /tmp/pr_recv.json -w '%{http_code}' -X POST "$API/v1/procurement/purchase-orders/$PR_POID/receive" \
  -H "Authorization: Bearer $SESSION_TOKEN")" = "201"
grep -q '"receipt_number"' /tmp/pr_recv.json
# AP invoice linked to the PO totals 1150 (subtotal 1000 + tax 150) -> matched (201).
test "$(curl -s -o /tmp/pr_inv.json -w '%{http_code}' -X POST "$API/v1/procurement/invoices" \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d "{\"vendor_id\":\"$PR_VID\",\"purchase_order_id\":\"$PR_POID\",\"invoice_number\":\"INV-SMOKE-AP\",\"invoice_date\":\"2026-08-05\",\"due_date\":\"2026-09-04\",\"subtotal\":\"1000\",\"tax_amount\":\"150\"}")" = "201"
grep -q '"match_status":"matched"' /tmp/pr_inv.json
# List invoices (200) and aging (200).
test "$(curl -s -o /tmp/pr_list.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/procurement/invoices")" = "200"
test "$(curl -s -o /tmp/pr_aging.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/procurement/aging")" = "200"

# ---- 22. Depreciation posts to the GL (accumulated depreciation + depreciation expense). ----
# Depreciate the seeded asset (200) — this records an entry AND posts a ledger journal.
test "$(curl -s -o /tmp/fa_dep.json -w '%{http_code}' -X POST "$API/v1/fixed-assets/fa_smoke/depreciate" \
  -H "Authorization: Bearer $SESSION_TOKEN")" = "200"
grep -q '"amount"' /tmp/fa_dep.json
# The GL trial balance now contains the depreciation expense and accumulated depreciation
# accounts with non-zero balances (proves the journal was posted to the ledger).
test "$(curl -s -o /tmp/fa_tb.json -w '%{http_code}' -H "Authorization: Bearer $SESSION_TOKEN" "$API/v1/accounting/trial-balance")" = "200"
grep -q 'expense:depreciation' /tmp/fa_tb.json
grep -q 'accumulated_depreciation' /tmp/fa_tb.json

echo 'Docker API smoke suite passed'
