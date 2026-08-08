-- Fiscal periods: enable closing books so postings are locked by month.
create table if not exists fiscal_periods (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  period      text not null, -- YYYY-MM
  status      text not null default 'open' check (status in ('open','closed')),
  closed_at   timestamptz,
  closed_by   text references users(id),
  created_at  timestamptz not null default now(),
  unique (merchant_id, period)
);

create index if not exists fiscal_periods_merchant_idx on fiscal_periods (merchant_id, period desc);
