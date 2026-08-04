-- 0003_onboarding_docs_compliance: document vault + compliance checks + reviews
-- MinIO encrypted file vault per DATABASE v1.1.0

create table merchant_documents (
  id                text primary key,
  merchant_id       text not null references merchants(id) on delete cascade,
  kyc_profile_id    text not null references merchant_kyc_profiles(id) on delete cascade,
  owner_id          text references merchant_beneficial_owners(id) on delete set null,
  doc_type          text not null check (doc_type in (
                      'company_registration','tin_certificate','business_license','vat_certificate',
                      'memorandum_articles','board_resolution','shareholder_list','audit_report',
                      'ubo_id_front','ubo_id_back','fayda_card_front','fayda_card_back',
                      'passport_front','proof_of_address','bank_letter','cancelled_cheque',
                      'website_screenshot','refund_policy_doc','trade_name_certificate','other'
                    )),
  file_key          text not null, -- merchants/{merchant_id}/kyc/{doc_type}_{id}.pdf in MinIO
  file_hash         text not null, -- sha256 integrity
  mime_type         text not null check (mime_type in ('application/pdf','image/jpeg','image/png','image/jpg')),
  file_size_bytes   int not null check (file_size_bytes >0 and file_size_bytes <= 5242880), -- 5MB max, Fayda 2MB app enforced
  status            text not null default 'pending' check (status in ('pending','uploaded','ocr_done','verified','rejected','expired')),
  ocr_data          jsonb not null default '{}'::jsonb,
  verified_by       text references users(id),
  rejection_reason  text,
  expires_at        timestamptz,
  created_at        timestamptz not null default now(),
  updated_at        timestamptz not null default now()
);
create index mdoc_merchant_type_idx on merchant_documents (merchant_id, doc_type);
create index mdoc_kyc_idx on merchant_documents (kyc_profile_id);
create unique index mdoc_file_hash_uidx on merchant_documents (merchant_id, file_hash);
create index mdoc_status_idx on merchant_documents (status);

create table compliance_checks (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  kyc_profile_id  text references merchant_kyc_profiles(id) on delete set null,
  check_type      text not null check (check_type in (
                    'tin_validation','business_license_validation','bank_account_validation',
                    'aml_screening','pep_check','restricted_industry','website_policy_check',
                    'fayda_verification','document_authenticity','risk_scoring','sanctions'
                  )),
  status          text not null check (status in ('pending','passed','failed','needs_review')),
  score           int check (score >=0 and score <=100),
  provider        text, -- internal, fayda, manual
  details         jsonb not null default '{}'::jsonb,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now()
);
create index compliance_merchant_idx on compliance_checks (merchant_id, check_type);
create index compliance_kyc_idx on compliance_checks (kyc_profile_id);

create table onboarding_reviews (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  kyc_profile_id  text references merchant_kyc_profiles(id) on delete cascade,
  reviewer_id     text references users(id),
  reviewer_type   text not null check (reviewer_type in ('system','ops','compliance','admin')),
  from_status     text not null,
  to_status       text not null,
  action          text not null check (action in ('submit','approve','reject','request_info','escalate','assign')),
  comments        text,
  internal_notes  text,
  created_at      timestamptz not null default now()
);
create index onboarding_reviews_merchant_idx on onboarding_reviews (merchant_id, created_at desc);
create index onboarding_reviews_kyc_idx on onboarding_reviews (kyc_profile_id, created_at desc);
