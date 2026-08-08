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
  ('la_docker_fee', 'book_docker_smoke', 'liability:platform_fee_due', 'Platform fee', 'credit')
ON CONFLICT (book_id, code) DO NOTHING;

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
