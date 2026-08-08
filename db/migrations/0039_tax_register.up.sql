-- Tax register: append-only record of collected VAT/TOT/withholding per period, used to
-- build the tax schedule (VAT/TOT return) and to post the tax liability into the GL.
create table if not exists tax_register (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  period      text not null,          -- YYYY-MM
  tax_type    text not null check (tax_type in ('vat','tot','withholding')),
  source      text not null,          -- invoice | payment
  source_id   text not null,
  amount      numeric(20,8) not null,
  paid        numeric(20,8) not null default 0,
  created_at  timestamptz not null default now(),
  unique (merchant_id, tax_type, source, source_id)
);
create index if not exists tax_register_merchant_idx on tax_register (merchant_id, period);
