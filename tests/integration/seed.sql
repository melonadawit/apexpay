-- Integration smoke-test seed.
-- Provides a merchant, its KYC/onboarding record, an API key with ops scope,
-- a payment link for hosted checkout, and a reconciliation case — so the smoke
-- script can exercise the real admin + checkout endpoints end-to-end.

INSERT INTO merchants (id, legal_name, display_name, email, status, onboarding_status, risk_score, risk_tier)
VALUES ('mer_docker_smoke', 'Docker Smoke PLC', 'Docker Smoke', 'docker-smoke@example.et', 'active', 'submitted', 42, 'medium')
ON CONFLICT (id) DO NOTHING;

-- Reviewer user for onboarding_reviews FK.
INSERT INTO users (id, email, name, status)
VALUES ('user_docker_ops', 'ops@example.et', 'Ops Reviewer', 'active')
ON CONFLICT (id) DO NOTHING;

-- KYC profile in a reviewable state (feeds /admin/onboarding/queue and /exam).
INSERT INTO merchant_kyc_profiles (
  id, merchant_id, version, legal_name, trade_name, business_type, registration_number, tin_number,
  industry_category, business_description, region, city, contact_person_name, contact_person_role,
  contact_email, contact_phone, office_address_full, onboarding_status, kyc_level
) VALUES (
  'kyc_docker_smoke', 'mer_docker_smoke', 1, 'Docker Smoke PLC', 'Docker Smoke', 'plc', 'MT/AA/123456',
  '0023456789', 'ecommerce', 'Smoke test merchant', 'Addis Ababa', 'Addis Ababa', 'Abebe Kebede',
  'owner', 'abebe@example.et', '0911111111', 'Bole, Woreda 03', 'submitted', 'level2'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO merchant_beneficial_owners (
  id, merchant_id, kyc_profile_id, full_name, role, ownership_percentage, nationality, id_type, phone, verification_status
) VALUES (
  'own_docker_smoke', 'mer_docker_smoke', 'kyc_docker_smoke', 'Abebe Kebede', 'owner', 100, 'ET', 'fayda', '0911111111', 'fayda_verified'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO merchant_documents (id, merchant_id, kyc_profile_id, doc_type, file_key, file_hash, mime_type, file_size_bytes, status)
VALUES
  ('doc_docker_reg', 'mer_docker_smoke', 'kyc_docker_smoke', 'company_registration', 'merchants/mer_docker_smoke/kyc/reg.pdf', 'hash_reg', 'application/pdf', 102400, 'verified'),
  ('doc_docker_tin', 'mer_docker_smoke', 'kyc_docker_smoke', 'tin_certificate', 'merchants/mer_docker_smoke/kyc/tin.pdf', 'hash_tin', 'application/pdf', 102400, 'verified')
ON CONFLICT (id) DO NOTHING;

INSERT INTO compliance_checks (id, merchant_id, kyc_profile_id, check_type, status, score, provider)
VALUES
  ('cc_docker_tin', 'mer_docker_smoke', 'kyc_docker_smoke', 'tin_validation', 'passed', 100, 'internal'),
  ('cc_docker_fayda', 'mer_docker_smoke', 'kyc_docker_smoke', 'fayda_verification', 'passed', 100, 'fayda')
ON CONFLICT (id) DO NOTHING;

INSERT INTO onboarding_reviews (id, merchant_id, kyc_profile_id, reviewer_id, reviewer_type, from_status, to_status, action, comments)
VALUES ('rev_docker_1', 'mer_docker_smoke', 'kyc_docker_smoke', 'user_docker_ops', 'ops', 'draft', 'submitted', 'submit', 'initial submission')
ON CONFLICT (id) DO NOTHING;

INSERT INTO bank_accounts (id, merchant_id, account_name, account_number_masked, account_number_hash, bank_code, bank_name, is_settlement_default, verification_status, verification_method)
VALUES ('bank_docker_smoke', 'mer_docker_smoke', 'Docker Smoke PLC', '****1234', 'hash_acct', 'CBE', 'Commercial Bank of Ethiopia', true, 'verified', 'bank_letter')
ON CONFLICT (id) DO NOTHING;

INSERT INTO ledger_books (id, merchant_id, book_type, name, currency, status)
VALUES ('book_docker_smoke', 'mer_docker_smoke', 'merchant_operating', 'Docker Smoke operating', 'ETB', 'open')
ON CONFLICT (merchant_id, book_type) DO NOTHING;

INSERT INTO ledger_accounts (id, book_id, code, name, normal_balance)
VALUES
  ('la_docker_clear', 'book_docker_smoke', 'asset:clearing:mock', 'Mock clearing', 'debit'),
  ('la_docker_payable', 'book_docker_smoke', 'liability:merchant_payable', 'Merchant payable', 'credit'),
  ('la_docker_fee', 'book_docker_smoke', 'liability:platform_fee_due', 'Platform fee', 'credit'),
  ('la_docker_bank', 'book_docker_smoke', 'asset:bank', 'Cash & Bank', 'debit'),
  ('la_docker_rev', 'book_docker_smoke', 'revenue:product', 'Product Revenue', 'credit'),
  ('la_docker_exp', 'book_docker_smoke', 'expense:operating', 'Operating Expenses', 'debit'),
  ('la_docker_equity', 'book_docker_smoke', 'equity:owner', 'Owner''s Equity', 'credit'),
  ('la_docker_dep_exp', 'book_docker_smoke', 'expense:depreciation', 'Depreciation', 'debit'),
  ('la_docker_accum_dep', 'book_docker_smoke', 'asset:accumulated_depreciation', 'Accumulated Depreciation', 'credit'),
  ('la_docker_cogs', 'book_docker_smoke', 'expense:cost_of_sales', 'Cost of Sales', 'debit'),
  ('la_docker_inv', 'book_docker_smoke', 'asset:inventory', 'Inventory', 'debit'),
  ('la_docker_tax', 'book_docker_smoke', 'liability:tax', 'Tax Payable', 'credit'),
  ('la_docker_ar', 'book_docker_smoke', 'asset:receivable', 'Accounts Receivable', 'debit'),
  ('la_docker_fxgain', 'book_docker_smoke', 'revenue:fx_gain', 'FX Gain', 'credit'),
  ('la_docker_fxloss', 'book_docker_smoke', 'expense:fx_loss', 'FX Loss', 'debit'),
  ('la_docker_ap_payable', 'book_docker_smoke', 'liability:payable', 'Accounts Payable', 'credit')
ON CONFLICT (book_id, code) DO NOTHING;

-- Vendor so the vendor self-service portal can be smoke-tested.
INSERT INTO vendors (id, merchant_id, name, email, phone, tin, payment_terms_days, status)
VALUES ('vend_smoke', 'mer_docker_smoke', 'Smoke Supplies Co', 'vendor@example.et', '+251911000000', 'ET-100000', 30, 'active')
ON CONFLICT (id) DO NOTHING;

-- Employee so the expense-claim -> GL reimbursement flow can be smoke-tested.
INSERT INTO employees (id, merchant_id, employee_code, name, base_salary, employment_date, employment_type, status, metadata)
VALUES ('emp_smoke', 'mer_docker_smoke', 'E-SMOKE', 'Smoke Employee', 50000, current_date - interval '365 days', 'permanent', 'active', '{}'::jsonb)
ON CONFLICT (id) DO NOTHING;

-- Product with a known cost price so order COGS posts to the GL in the smoke suite.
INSERT INTO products (id, merchant_id, name, sku, price, cost_price, currency, vat_category, stock_qty, low_stock_threshold, status)
VALUES ('prod_smoke', 'mer_docker_smoke', 'Smoke Widget', 'SKU-SMOKE', 2000, 1200, 'ETB', 'standard', 50, 5, 'active')
ON CONFLICT (id) DO NOTHING;

-- API key with explicit ops scope (RBAC admin routes).
INSERT INTO api_keys (id, merchant_id, name, key_type, key_prefix, secret_hash, environment, status, scopes)
VALUES (
  'key_docker_smoke', 'mer_docker_smoke', 'integration', 'secret',
  'sk_test_demo',
  'cc081db7827589a0f9d454de130b981be34315e01c9a440943bf74bb176b2b45',
  'test', 'active', '["role:ops"]'
)
ON CONFLICT (id) DO NOTHING;

-- Payment link for hosted checkout (/v1/checkout/{token}).
INSERT INTO payment_links (id, merchant_id, amount, currency, description, status, public_token, expires_at)
VALUES ('pl_docker_smoke', 'mer_docker_smoke', 6000.00, 'ETB', 'Docker Smoke Order', 'active', 'tkn_docker_smoke', now() + interval '24 hours')
ON CONFLICT (id) DO NOTHING;

-- Open reconciliation case (feeds /admin/recon/breaks and /admin/evidence).
INSERT INTO payment_reconciliation_cases (merchant_id, idempotency_key, tx_ref, status)
VALUES ('mer_docker_smoke', 'idem_recon_smoke', 'txr_recon_smoke', 'open')
ON CONFLICT (merchant_id, idempotency_key) DO NOTHING;

-- Dashboard user for session auth. Email demo@apexpay.et / password Admin@12345
-- (argon2id hash, random salt per hash).
INSERT INTO users (id, email, name, status, email_verified)
VALUES ('user_docker_admin', 'demo@apexpay.et', 'Demo Admin', 'active', true)
ON CONFLICT (id) DO NOTHING;

UPDATE users SET password_hash = '$argon2id$v=19$m=65536,t=1,p=4$KO6YyM3nOeOgcEmBGietGg$VndoW28CBk6wfmj9cQLHJ8BBU7P19Lqyxdpfk1A5wYI'
WHERE id = 'user_docker_admin' AND (password_hash IS NULL OR password_hash = '');

INSERT INTO merchant_members (merchant_id, user_id, role)
VALUES ('mer_docker_smoke', 'user_docker_admin', 'owner')
ON CONFLICT (merchant_id, user_id) DO NOTHING;

-- Banking module seed rows for smoke-testing the real banking read endpoints.
INSERT INTO forex_rates (id, from_currency, to_currency, rate, buy_rate, sell_rate, source)
VALUES ('fx_smoke_usd', 'ETB', 'USD', 57.50, 56.80, 58.20, 'nbe')
ON CONFLICT (from_currency, to_currency) DO NOTHING;

INSERT INTO notifications (id, merchant_id, user_id, type, title, message, data, is_read, action_url)
VALUES ('notif_smoke', 'mer_docker_smoke', 'user_docker_admin', 'current_account_opened', 'Current Account Opened',
        'Your current account is ready', '{"account":"smoke"}', false, '/banking/current-accounts')
ON CONFLICT (id) DO NOTHING;

INSERT INTO current_accounts (id, merchant_id, account_number, account_name, account_type, currency, bank_code, partner_bank_name, status, balance, available_balance)
VALUES ('ca_smoke', 'mer_docker_smoke', 'ETB-CBE-7778889990', 'Docker Smoke PLC', 'current', 'ETB', 'CBE', 'Commercial Bank of Ethiopia', 'active', 125000.00, 125000.00)
ON CONFLICT (id) DO NOTHING;

-- Foreign-currency account so multi-currency revaluation can be exercised in smoke.
-- Held USD: 10,000 at booked cost; the ETB-equivalent is revalued against forex_rates.
INSERT INTO current_accounts (id, merchant_id, account_number, account_name, account_type, currency, bank_code, partner_bank_name, status, balance, available_balance)
VALUES ('ca_smoke_usd', 'mer_docker_smoke', 'USD-CBE-1000001', 'Docker Smoke USD', 'current', 'USD', 'CBE', 'Commercial Bank of Ethiopia', 'active', 10000.00, 10000.00)
ON CONFLICT (id) DO NOTHING;

-- Prior FX revaluation baseline for the USD account at an older rate (50.0 ETB/USD) so the
-- current-period revaluation (57.5) produces a recognized unrealized FX gain in the smoke.
INSERT INTO fx_revaluations (id, merchant_id, period, account_id, currency, amount_fx, rate, amount_etb, fx_gain)
VALUES ('fxr_baseline_usd', 'mer_docker_smoke', '2026-07', 'ca_smoke_usd', 'USD', 10000, 50.00, 500000.00, 0)
ON CONFLICT (id) DO NOTHING;

-- HRIS + Risk smoke seed.
INSERT INTO hris_teams (id, merchant_id, name) VALUES ('team_smoke', 'mer_docker_smoke', 'Engineering')
ON CONFLICT (id) DO NOTHING;

INSERT INTO risk_rules (id, merchant_id, name, rule_type, parameters, action, severity, enabled)
VALUES ('rule_smoke', 'mer_docker_smoke', 'High-ticket', 'threshold_amount',
        '{"amount_limit":"500000","window_minutes":60}', 'flag', 'high', true)
ON CONFLICT (id) DO NOTHING;

-- Treasury + Invoicing smoke seed.
INSERT INTO invoices (id, merchant_id, invoice_number, customer_name, customer_email, issue_date, due_date, currency, subtotal, tax_amount, withholding_amount, total_amount, status)
VALUES ('inv_smoke', 'mer_docker_smoke', 'INV-SMOKE-001', 'Smoke Customer', 'customer@example.et', current_date, current_date + interval '30 days', 'ETB', 1000.00, 150.00, 20.00, 1130.00, 'sent')
ON CONFLICT (id) DO NOTHING;

-- Compliance + fixed asset + analytics smoke seed.
INSERT INTO compliance_console (merchant_id, onboarding_status, fayda_verified, risk_tier, overall_status)
VALUES ('mer_docker_smoke', 'approved', true, 'medium', 'good')
ON CONFLICT (merchant_id) DO NOTHING;

INSERT INTO fixed_assets (id, merchant_id, asset_name, category, acquisition_date, cost, useful_life_years, depreciation_method, net_book_value)
VALUES ('fa_smoke', 'mer_docker_smoke', 'Delivery Van', 'vehicle', current_date - interval '365 days', 1200000, 5, 'straight_line', 1200000)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Real payroll + payout data so the merchant UI shows live rows
-- (payroll calendar, salary structure, run + items, final
-- settlement, payout batch + beneficiary).
-- ============================================================

-- A few more employees across cost centers so reports have real headcount.
INSERT INTO employees (id, merchant_id, employee_code, name, base_salary, employment_date, employment_type, cost_center, status, metadata)
VALUES
  ('emp_eng_01', 'mer_docker_smoke', 'E-ENG-001', 'Abebe Bekele', 60000, current_date - interval '400 days', 'permanent', 'Engineering', 'active', '{}'::jsonb),
  ('emp_eng_02', 'mer_docker_smoke', 'E-ENG-002', 'Sara Tesfaye', 45000, current_date - interval '300 days', 'permanent', 'Engineering', 'active', '{}'::jsonb),
  ('emp_sales_01', 'mer_docker_smoke', 'E-SAL-001', 'Mekdes Ali', 40000, current_date - interval '250 days', 'permanent', 'Sales', 'active', '{}'::jsonb)
ON CONFLICT (id) DO NOTHING;
UPDATE employees SET cost_center='Engineering' WHERE id='emp_smoke';

-- Salary structure (referenced by payroll settings)
INSERT INTO payroll_salary_structures (id, merchant_id, name, ctc_annual, ctc_monthly, currency, effective_from, status, is_default)
VALUES ('ss_smoke', 'mer_docker_smoke', 'Docker Smoke Band G3', 600000, 50000, 'ETB', '2026-01-01', 'active', true)
ON CONFLICT (id) DO NOTHING;

-- Payroll calendar 2026 (Ethiopia business practice: cutoff 25th, disbursal 30th, pay last day).
INSERT INTO payroll_calendars (id, merchant_id, name, pay_frequency, year, month, cutoff_day, disbursal_day, pay_day, cutoff_date, disbursal_date, pay_date, is_locked, locked_at, locked_by, created_by)
VALUES
  ('cal_2026_07', 'mer_docker_smoke', 'Monthly Payroll Calendar 2026-07', 'monthly', 2026, 7, 25, 30, 31, '2026-07-25', '2026-07-30', '2026-07-31', true, '2026-07-30T22:00:00Z', 'user_docker_admin', 'user_docker_admin'),
  ('cal_2026_08', 'mer_docker_smoke', 'Monthly Payroll Calendar 2026-08', 'monthly', 2026, 8, 25, 30, 31, '2026-08-25', '2026-08-30', '2026-08-31', false, NULL, NULL, 'user_docker_admin')
ON CONFLICT (id) DO NOTHING;

-- Payroll run 2026-07 (completed) with per-employee items.
INSERT INTO payroll_runs (id, merchant_id, book_id, run_ref, period_month, period_year, type, status,
  total_gross, total_deductions, total_net, total_tax, total_pension, employer_total_pension,
  total_employer_cost, total_count, total_employees_paid, total_employees_failed, created_by)
VALUES ('run_2026_07', 'mer_docker_smoke', 'book_docker_smoke', 'RUN-2026-07', 7, 2026, 'regular', 'completed',
  195000, 42500, 152500, 32000, 13650, 21450, 216450, 4, 4, 0, 'user_docker_admin')
ON CONFLICT (id) DO NOTHING;

INSERT INTO payroll_items (id, run_id, employee_id, gross, ot_hours, ot_amount, commission, bonus, other_allowances,
  taxable_income, income_tax, pension_employee, pension_employer, other_deductions, net_pay, status)
VALUES
  ('pi_1', 'run_2026_07', 'emp_smoke',   50000, 0, 0, 0, 2000, 0, 46500,  9200, 3500, 5500, 3000, 37300, 'paid'),
  ('pi_2', 'run_2026_07', 'emp_eng_01', 60000, 8, 1500, 0, 0, 0, 54300,  12000, 4200, 6600, 3000, 46800, 'paid'),
  ('pi_3', 'run_2026_07', 'emp_eng_02', 45000, 0, 0, 0, 0, 0, 41850,  8000,  3150, 4950, 2500, 35850, 'paid'),
  ('pi_4', 'run_2026_07', 'emp_sales_01', 40000, 0, 0, 1500, 0, 0, 37200, 6500, 2800, 4400, 2000, 32500, 'paid')
ON CONFLICT (id) DO NOTHING;

-- Stored cost-center compliance report so /payroll_reports/cost_center returns live data.
INSERT INTO payroll_compliance_reports (id, merchant_id, period_month, period_year, report_type, status, metadata)
VALUES ('rpt_cc_202607', 'mer_docker_smoke', 7, 2026, 'cost_center_report', 'generated',
  '{"cost_centers":[{"cost_center":"Engineering","total_gross":"155000","total_net":"119950","headcount":3,"employer_cost":"171050","paid_days":90,"lop_days":0},{"cost_center":"Sales","total_gross":"40000","total_net":"32500","headcount":1,"employer_cost":"44400","paid_days":30,"lop_days":0}]}'::jsonb)
ON CONFLICT (id) DO NOTHING;

-- Final settlement (F&F) for an existing employee.
INSERT INTO payroll_final_settlements (id, merchant_id, employee_id, resignation_date, last_working_date,
  notice_period_days, notice_served_days, notice_shortfall_days, leave_encashment_days, leave_encashment_amount,
  severance_amount, gratuity_amount, bonus_pro_rata, outstanding_loans, outstanding_advances, other_earnings,
  other_deductions, total_payable, total_deductions, net_payable, status,
  clearance_checklist, clearance_items_detailed, assets_returned, exit_interview)
VALUES ('fnf_smoke', 'mer_docker_smoke', 'emp_smoke', '2026-06-15', '2026-07-15',
  30, 30, 0, 5, 2500, 15000, 0, 2000, 5000, 0, 0, 0, 19500, 5000, 14500, 'pending_approval',
  '[{"item":"Laptop","status":"pending","checked_by":"","notes":"MacBook Pro"}]'::jsonb,
  '[{"item":"Laptop LP001","category":"IT","status":"pending","required":true,"checked_by":"","checked_at":"","notes":"MacBook Pro 14 inch"},{"item":"ID Card ID-EMP007","category":"HR","status":"done","required":true,"checked_by":"HR Manager","checked_at":"2026-07-14","notes":"Returned"}]'::jsonb,
  '[{"asset_type":"laptop","asset_id":"LP001","returned":false,"condition":"good","returned_at":""}]'::jsonb,
  '{"conducted":false,"conducted_by":"","date":"","feedback":""}'::jsonb)
ON CONFLICT (id) DO NOTHING;

-- Payout beneficiary + batch + a payout so the payouts screen shows a live batch.
INSERT INTO beneficiaries (id, merchant_id, name, account_no_masked, account_no_hash, bank_code, bank_name, type, verification_status)
VALUES ('ben_smoke', 'mer_docker_smoke', 'Abebe Bekele', '****6789', 'hash_ben_1', 'CBE', 'Commercial Bank of Ethiopia', 'individual', 'verified')
ON CONFLICT (id) DO NOTHING;

INSERT INTO payout_batches (id, merchant_id, book_id, batch_ref, amount, currency, status, total_count, success_count, failed_count)
VALUES ('pbat_smoke', 'mer_docker_smoke', 'book_docker_smoke', 'PBATCH-SMOKE-001', 50000, 'ETB', 'completed', 1, 1, 0)
ON CONFLICT (id) DO NOTHING;

INSERT INTO payouts (id, merchant_id, batch_id, beneficiary_id, payout_ref, amount, currency, status, method)
VALUES ('pout_smoke', 'mer_docker_smoke', 'pbat_smoke', 'ben_smoke', 'POUT-SMOKE-001', 50000, 'ETB', 'succeeded', 'bank')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Subscription + refund seed data for the detail pages.
-- ============================================================

-- A settled payment to anchor the refund (FK to payments).
INSERT INTO payments (id, merchant_id, tx_ref, amount, currency, status, method, connector_id, connector_ref, fee_amount, net_amount, requires_2fa, two_fa_verified)
VALUES ('pay_smoke_ref', 'mer_docker_smoke', 'txr_refund_smoke', 250.00, 'ETB', 'succeeded', 'mock', 'mock', 'mock_ref_refund_smoke', 7.25, 242.75, false, false)
ON CONFLICT (id) DO NOTHING;

-- Subscription customer + plan + subscription + invoice.
INSERT INTO customers (id, merchant_id, email, phone, name)
VALUES ('cust_smoke', 'mer_docker_smoke', 'abebe@example.et', '+251911000001', 'Abebe Kebede')
ON CONFLICT (id) DO NOTHING;

INSERT INTO subscription_plans (id, merchant_id, name, description, amount, currency, interval_type, interval_count, trial_days, status)
VALUES ('splan_smoke', 'mer_docker_smoke', 'Monthly Coffee', 'Premium coffee subscription', 500, 'ETB', 'month', 1, 7, 'active')
ON CONFLICT (id) DO NOTHING;

INSERT INTO subscriptions (id, merchant_id, customer_id, plan_id, status, current_period_start, current_period_end, trial_end)
VALUES ('sub_smoke', 'mer_docker_smoke', 'cust_smoke', 'splan_smoke', 'active', now() - interval '20 days', now() + interval '10 days', now() - interval '13 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO subscription_invoices (id, merchant_id, subscription_id, payment_id, amount, currency, status, attempt_count, due_at)
VALUES
  ('sinv_smoke_1', 'mer_docker_smoke', 'sub_smoke', 'pay_smoke_ref', 500, 'ETB', 'paid', 0, now() - interval '10 days'),
  ('sinv_smoke_2', 'mer_docker_smoke', 'sub_smoke', NULL, 500, 'ETB', 'open', 1, now() + interval '10 days')
ON CONFLICT (id) DO NOTHING;

-- A refund anchored to the payment above.
INSERT INTO refunds (id, merchant_id, payment_id, refund_ref, amount, currency, status, reason, fee_reversal, connector_id, connector_ref)
VALUES ('ref_smoke', 'mer_docker_smoke', 'pay_smoke_ref', 'RFD-SMOKE-001', 100.00, 'ETB', 'succeeded', 'customer requested partial refund', 0, 'mock', 'mock_ref_refund_smoke')
ON CONFLICT (id) DO NOTHING;
