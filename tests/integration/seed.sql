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
