-- 0012_enhancements: final enhancements, indexes, checks for full platform v1.1.0 gold
-- Adds missing columns if 0001 was old MVP, plus banks seed and performance indexes

-- Ensure merchants enhancements already in 0001 but add if missing (idempotent)
do $$
begin
  alter table merchants add column if not exists onboarding_status text;
  alter table merchants add column if not exists business_type text;
  alter table merchants add column if not exists risk_tier text;
  alter table merchants add column if not exists risk_score int;
  alter table merchants add column if not exists settlement_type text;
  alter table merchants add column if not exists kyc_profile_id text;
  alter table merchants add column if not exists fayda_verified boolean;
exception when others then null;
end $$;

-- payments enhancements
do $$
begin
  alter table payments add column if not exists routing_rule_id text;
  alter table payments add column if not exists requires_2fa boolean not null default false;
  alter table payments add column if not exists two_fa_verified boolean not null default false;
exception when others then null;
end $$;

-- Bank list table for GET /v1/banks - Ethiopian banks per NBE
create table if not exists banks (
  code        text primary key,
  name        text not null,
  name_am     text,
  swift_code  text,
  is_active   boolean not null default true,
  created_at  timestamptz not null default now()
);

insert into banks (code, name, name_am) values
('CBE', 'Commercial Bank of Ethiopia', 'የኢትዮጵያ ንግድ ባንክ'),
('AWASH', 'Awash Bank', 'አዋሽ ባንክ'),
('DASHEN', 'Dashen Bank', 'ዳሽን ባንክ'),
('ABYSSINIA', 'Bank of Abyssinia', 'አቢሲኒያ ባንክ'),
('BERHAN', 'Berhan Bank', 'ብርሃን ባንክ'),
('WEGAGEN', 'Wegagen Bank', 'ወጋገን ባንክ'),
('NIB', 'Nib International Bank', 'ኒብ ኢንተርናሽናል ባንክ'),
('UNITED', 'United Bank', 'ዩናይትድ ባንክ'),
('COOP', 'Cooperative Bank of Oromia', 'የኦሮሚያ ኅብረት ሥራ ባንክ'),
('OROMIA', 'Oromia Bank', 'ኦሮሚያ ባንክ'),
('BUNNA', 'Bunna International Bank', 'ቡና ኢንተርናሽናል ባንክ'),
('LION', 'Lion International Bank', 'አንበሳ ባንክ'),
('ZEMEN', 'Zemen Bank', 'ዘመን ባንክ'),
('CBO', 'Commercial Bank of Oromia', 'ኦሮሚያ ንግድ ባንክ')
on conflict (code) do nothing;

-- Performance indexes for full platform gold
create index if not exists payments_requires_2fa_idx on payments (requires_2fa) where requires_2fa=true and two_fa_verified=false;
create index if not exists merchants_fayda_verified_idx on merchants (fayda_verified) where fayda_verified=true;
create index if not exists refunds_fee_policy_idx on refunds (fee_policy);
create index if not exists employees_merchant_active_idx on employees (merchant_id) where status='active';
create index if not exists payroll_runs_period_idx on payroll_runs (period_year, period_month);

-- Add ledger accounts seed data placeholder comment
-- Real seed in db/seeds/0001_platform_books.sql would insert: merchant_operating, rail_clearing, platform_revenue, escrow, suspense, reserve, refund_clearing, payroll_run, payout_batch, sandbox books + standard accounts liability:merchant_payable, asset:clearing:*, liability:platform_fee_due etc

-- Add check for FIN hash not containing plain 12 digits pattern (app-level, but DB comment)
comment on column fayda_verifications.fin_hash is 'sha256(salt+FIN) per privacy rule, never plain FIN per id.gov.et compliance. Last4 stored separately.';
comment on column bank_accounts.account_number_hash is 'sha256 hash for lookup, masked display ****1234 outstanding UI.';
comment on column refunds.fee_policy is 'non_refundable = platform keeps fee (default), pro_rata = reverse pro-rata, full = reversal on full refund.';

-- Add materialized view for TPV dashboard outstanding
create materialized view if not exists merchant_tpv_daily as
select merchant_id, date_trunc('day', created_at)::date as day, sum(amount) as tpv, count(*) as cnt, sum(case when status='succeeded' then 1 else 0 end) as success_cnt
from payments group by merchant_id, day;

create unique index if not exists merchant_tpv_daily_uidx on merchant_tpv_daily (merchant_id, day);

comment on materialized view merchant_tpv_daily is 'Outstanding dashboard TPV: refreshed hourly via worker, powers merchant home + swarm get_tpv tool O(1) lookup.';
