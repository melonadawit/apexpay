-- 0002_onboarding_kyc: NBE onboarding KYC profiles, beneficial owners, bank accounts
-- Per PayAtlas ET PSP requirements + NBE ONPS/02/2020 shareholder min 5 / 10 multi-system

create table merchant_kyc_profiles (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  version             int not null default 1,
  legal_name          text not null,
  trade_name          text,
  business_type       text not null check (business_type in ('sole_proprietorship','plc','share_company','partnership','cooperative','government','ngo','other')),
  registration_number text not null,
  registration_date   date,
  tin_number          text not null check (char_length(tin_number)=10),
  vat_number          text,
  business_license_no text,
  license_expiry      date,
  industry_category   text not null,
  business_description text not null,
  website_url         text,
  app_url             text,
  annual_turnover     numeric(20,8),
  expected_monthly_tpv numeric(20,8),
  avg_ticket_amount   numeric(20,8),
  employee_count      int,
  region              text not null,
  city                text not null,
  sub_city            text,
  woreda              text,
  kebele              text,
  house_no            text,
  office_address_full text not null,
  gps_lat             numeric(10,7),
  gps_lng             numeric(10,7),
  contact_person_name text not null,
  contact_person_role text not null,
  contact_email       text not null,
  contact_phone       text not null,
  has_refund_policy   boolean not null default false,
  has_privacy_policy  boolean not null default false,
  has_terms_and_conditions boolean not null default false,
  refund_policy_url   text,
  privacy_policy_url  text,
  terms_url           text,
  onboarding_status   text not null default 'draft' check (onboarding_status in ('draft','submitted','in_review','fayda_pending','compliance_check','needs_more_info','approved','rejected')),
  kyc_level           text not null default 'level2' check (kyc_level in ('level1','level2','level3')),
  risk_notes          text,
  submitted_at        timestamptz,
  reviewed_at         timestamptz,
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now(),
  unique (merchant_id, version)
);
create index kyc_merchant_idx on merchant_kyc_profiles (merchant_id, created_at desc);
create index kyc_status_idx on merchant_kyc_profiles (onboarding_status);

create table merchant_beneficial_owners (
  id                    text primary key,
  merchant_id           text not null references merchants(id) on delete cascade,
  kyc_profile_id        text not null references merchant_kyc_profiles(id) on delete cascade,
  full_name             text not null,
  full_name_am          text,
  role                  text not null check (role in ('owner','shareholder','director','authorized_rep','contact_person','ubo')),
  ownership_percentage  numeric(5,2) check (ownership_percentage >=0 and ownership_percentage <=100),
  nationality           char(2) not null default 'ET',
  id_type               text not null check (id_type in ('fayda','passport','driving_license','kebele_id','other')),
  id_number_hash        text,
  id_number_last4       text,
  fayda_fin_hash        text,
  fayda_fan             text,
  fayda_verified        boolean not null default false,
  date_of_birth         date,
  gender                text check (gender in ('male','female','other')),
  phone                 text not null,
  email                 text,
  address               text,
  is_pep                boolean not null default false,
  is_authorized_signatory boolean not null default false,
  verification_status   text not null default 'pending' check (verification_status in ('pending','fayda_verified','document_verified','verified','rejected')),
  created_at            timestamptz not null default now(),
  updated_at            timestamptz not null default now()
);
create index mbo_merchant_idx on merchant_beneficial_owners (merchant_id);
create index mbo_fayda_hash_idx on merchant_beneficial_owners (fayda_fin_hash) where fayda_fin_hash is not null;
create index mbo_kyc_idx on merchant_beneficial_owners (kyc_profile_id);

create table bank_accounts (
  id                    text primary key,
  merchant_id           text not null references merchants(id) on delete cascade,
  account_name          text not null,
  account_number_masked text not null,
  account_number_hash   text not null,
  bank_code             text not null,
  bank_name             text not null,
  branch                text,
  account_type          text not null check (account_type in ('current','saving')) default 'current',
  is_settlement_default boolean not null default false,
  verification_status   text not null default 'pending' check (verification_status in ('pending','verified','failed')),
  verification_method   text check (verification_method in ('bank_letter','micro_deposit','manual')),
  verified_at           timestamptz,
  created_at            timestamptz not null default now()
);
create index bank_merchant_idx on bank_accounts (merchant_id);
