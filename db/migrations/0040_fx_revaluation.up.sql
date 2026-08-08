-- FX revaluation register: append-only record of unrealized foreign-currency revaluation
-- of current accounts, with the FX gain/loss recognized and posted to the GL.
create table if not exists fx_revaluations (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  period      text not null,            -- YYYY-MM
  account_id  text not null,
  currency    text not null,
  amount_fx   numeric(20,8) not null,
  rate        numeric(20,8) not null,   -- local (ETB) per unit of foreign currency
  amount_etb  numeric(20,8) not null,
  fx_gain     numeric(20,8) not null,   -- positive = gain, negative = loss
  created_at  timestamptz not null default now()
);
create index if not exists fx_revaluations_merchant_idx on fx_revaluations (merchant_id, period);
