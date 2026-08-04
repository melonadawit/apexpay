-- 0004_fayda_verification: Fayda National ID verification per id.gov.et spec
-- FIN stored as sha256(salt+fin) + last4 only, never plain, response encrypted ref

create table fayda_verifications (
  id                      text primary key,
  merchant_id             text not null references merchants(id) on delete cascade,
  owner_id                text references merchant_beneficial_owners(id) on delete set null,
  kyc_profile_id          text references merchant_kyc_profiles(id) on delete set null,
  fin_hash                text not null,
  fin_last4               char(4) not null,
  fan                     text,
  partner_code            text not null,
  request_id              text not null unique,
  fayda_transaction_id    text,
  verification_method     text not null check (verification_method in ('otp','face','fingerprint','offline_qr','oidc_esignet','demographic')),
  otp_requested_at        timestamptz,
  otp_verified            boolean not null default false,
  otp_verified_at         timestamptz,
  consent_timestamp       timestamptz not null,
  consent_ip              inet,
  status                  text not null check (status in ('initiated','otp_sent','pending_consent','verified','failed','expired','revoked')),
  demographics_match      boolean,
  demographics_score      int check (demographics_score >=0 and demographics_score <=100),
  face_match              boolean,
  face_match_score        numeric(5,2) check (face_match_score >=0 and face_match_score <=1),
  response_encrypted_ref  text, -- MinIO key fayda_responses/{request_id}.enc
  response_hash           text,
  failure_code            text,
  failure_message         text,
  front_doc_id            text references merchant_documents(id) on delete set null,
  back_doc_id             text references merchant_documents(id) on delete set null,
  selfie_doc_id           text references merchant_documents(id) on delete set null,
  verified_at             timestamptz,
  expires_at              timestamptz,
  created_at              timestamptz not null default now(),
  updated_at              timestamptz not null default now()
);
create index fayda_merchant_idx on fayda_verifications (merchant_id, created_at desc);
create index fayda_owner_idx on fayda_verifications (owner_id);
create index fayda_fin_hash_idx on fayda_verifications (fin_hash);
create index fayda_status_idx on fayda_verifications (status);
create unique index fayda_request_uidx on fayda_verifications (request_id);

create table fayda_qr_verifications (
  id                text primary key,
  verification_id   text not null references fayda_verifications(id) on delete cascade,
  qr_data_hash      text not null,
  scanned_at        timestamptz not null default now(),
  offline_verified  boolean not null,
  created_at        timestamptz not null default now()
);
create index fayda_qr_verification_idx on fayda_qr_verifications (verification_id);

-- Add fayda_verified flag index already exists but ensure
comment on table fayda_verifications is 'Fayda FIN 12-digit per NIDP Ethiopia: stored as hash + last4 only, OTP consent required, front/back images <2MB, offline QR FaydaEncode, OIDC eSignet per id.gov.et/api';
