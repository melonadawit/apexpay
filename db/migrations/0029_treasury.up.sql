-- 0029: Treasury & Cash Management.
-- Internal account transfers/sweeps and cash-flow forecasts across the merchant's
-- accounts (current + virtual + escrow + cards).

create table treasury_transfers (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  from_account_id text not null references current_accounts(id),
  to_account_id   text not null references current_accounts(id),
  amount          numeric(20,8) not null check (amount > 0),
  currency        char(3) not null default 'ETB',
  purpose         text, -- concentration, sweep, funding
  status          text not null default 'pending' check (status in ('pending','completed','failed','cancelled')),
  ledger_journal_id text,
  created_by      text references users(id),
  created_at      timestamptz not null default now(),
  completed_at    timestamptz,
  updated_at      timestamptz not null default now(),
  check (from_account_id <> to_account_id)
);
create index treasury_transfers_merchant_idx on treasury_transfers (merchant_id, created_at desc);

create table cash_flow_forecasts (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  forecast_date date not null,
  horizon_days  int not null default 90,
  -- Buckets (numeric::text at read; stored as numeric):
  inflow_today    numeric(20,8) not null default 0,
  inflow_30d      numeric(20,8) not null default 0,
  inflow_60d      numeric(20,8) not null default 0,
  inflow_90d      numeric(20,8) not null default 0,
  outflow_today   numeric(20,8) not null default 0,
  outflow_30d     numeric(20,8) not null default 0,
  outflow_60d     numeric(20,8) not null default 0,
  outflow_90d     numeric(20,8) not null default 0,
  net_90d         numeric(20,8) not null default 0,
  generated_at    timestamptz not null default now()
);
create index cash_flow_forecasts_merchant_idx on cash_flow_forecasts (merchant_id, forecast_date desc);
