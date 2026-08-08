-- 0032: Compliance console, notification preferences, fixed-asset & depreciation.

-- Compliance console: tracks merchant compliance status, deadlines, and check results.
create table compliance_console (
  merchant_id      text primary key references merchants(id) on delete cascade,
  onboarding_status text not null default 'pending',
  kyc_expiry_date  date,
  license_expiry   date,
  fayda_verified   boolean not null default false,
  risk_tier        text not null default 'low',
  next_erca_due    date,
  next_pension_due date,
  annual_tax_filing_due date,
  aml_due          date,
  overall_status   text not null default 'attention' check (overall_status in ('good','attention','overdue')),
  notes            text,
  updated_at       timestamptz not null default now()
);

create table compliance_checks_log (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  check_type    text not null, -- kyc, license, tax_filing, pension, aml, fayda
  status        text not null, -- passed, failed, due, overdue, pending
  detail        text,
  checked_at    timestamptz not null default now()
);
create index compliance_checks_log_merchant_idx on compliance_checks_log (merchant_id, checked_at desc);

-- Notification preferences: which event types a user wants via which channels.
create table notification_preferences (
  merchant_id text not null references merchants(id) on delete cascade,
  user_id     text not null references users(id) on delete cascade,
  event_type  text not null, -- bulk_payouts_approval, payroll_run_completed, tax_payment_due, escrow_released, etc.
  email       boolean not null default true,
  sms         boolean not null default false,
  push        boolean not null default true,
  in_app      boolean not null default true,
  primary key (merchant_id, user_id, event_type)
);

-- Fixed assets & depreciation (Ethiopian businesses need this for tax).
create table fixed_assets (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  asset_name      text not null,
  category        text not null check (category in ('building','machinery','vehicle','equipment','furniture','computer','land','other')),
  acquisition_date date not null,
  cost            numeric(20,8) not null check (cost > 0),
  salvage_value   numeric(20,8) not null default 0,
  useful_life_years int not null check (useful_life_years > 0),
  depreciation_method text not null default 'straight_line' check (depreciation_method in ('straight_line','declining_balance')),
  depreciation_rate numeric(6,4), -- for declining balance
  accumulated_depreciation numeric(20,8) not null default 0,
  net_book_value  numeric(20,8) not null,
  status          text not null default 'active' check (status in ('active','disposed','written_off')),
  notes           text,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now()
);
create index fixed_assets_merchant_idx on fixed_assets (merchant_id);

create table depreciation_entries (
  id              text primary key,
  asset_id        text not null references fixed_assets(id) on delete cascade,
  merchant_id     text not null references merchants(id) on delete cascade,
  period          text not null, -- e.g. 2026-07
  amount          numeric(20,8) not null,
  book_value_after numeric(20,8) not null,
  created_at      timestamptz not null default now(),
  unique (asset_id, period)
);
create index depreciation_entries_merchant_idx on depreciation_entries (merchant_id, period);
