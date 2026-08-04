-- 0001_init: MVP core + platform extensions
-- PostgreSQL 16+ per DATABASE v1.1.0
-- Best practice: forward-only, no destructive prod down

create extension if not exists pgcrypto;
create extension if not exists "uuid-ossp";

-- merchants enhanced but keep backward compatible with MVP 1.0.0
create table if not exists merchants (
  id                text primary key,
  legal_name        text not null,
  display_name      text not null,
  email             text not null,
  country_code      char(2) not null default 'ET',
  default_currency  char(3) not null default 'ETB',
  status            text not null check (status in ('draft','pending_kyc','kyc_in_review','fayda_pending','compliance_check','pending_approval','active','suspended','closed','rejected','pending','suspended','closed','active')),
  onboarding_status text not null default 'not_started' check (onboarding_status in ('not_started','draft','submitted','in_review','fayda_verification','compliance_review','approved','rejected','needs_more_info','fayda_pending','active')),
  business_type     text check (business_type in ('sole_proprietorship','plc','share_company','partnership','cooperative','government','ngo','other')),
  risk_tier         text check (risk_tier in ('low','medium','high','prohibited')) default 'low',
  risk_score        int not null default 0,
  mdr_rate          numeric(10,6) not null default 0.029000,
  settlement_type   text not null default 'T+1' check (settlement_type in ('T+0','T+1','T+2','weekly')),
  branding          jsonb not null default '{}'::jsonb,
  metadata          jsonb not null default '{}'::jsonb,
  kyc_profile_id    text,
  fayda_verified    boolean not null default false,
  created_at        timestamptz not null default now(),
  updated_at        timestamptz not null default now()
);
create unique index if not exists merchants_email_uidx on merchants (lower(email));
create index if not exists merchants_status_idx on merchants (status);
create index if not exists merchants_onboarding_idx on merchants (onboarding_status);

create table if not exists users (
  id              text primary key,
  email           text not null,
  password_hash   text,
  name            text not null,
  phone           text,
  status          text not null check (status in ('active','disabled','pending_verification')),
  email_verified  boolean not null default false,
  phone_verified  boolean not null default false,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now()
);
create unique index if not exists users_email_uidx on users (lower(email));

create table if not exists merchant_members (
  merchant_id  text not null references merchants(id),
  user_id      text not null references users(id),
  role         text not null check (role in ('owner','admin','developer','finance','support','ops','compliance','viewer')),
  permissions  jsonb not null default '[]'::jsonb,
  created_at   timestamptz not null default now(),
  primary key (merchant_id, user_id)
);

create table if not exists api_keys (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  name            text not null,
  key_type        text not null check (key_type in ('secret','public')),
  key_prefix      text not null,
  secret_hash     text,
  environment     text not null check (environment in ('test','live')),
  status          text not null check (status in ('active','revoked','pending_activation')),
  scopes          jsonb not null default '[]'::jsonb,
  last_used_at    timestamptz,
  created_at      timestamptz not null default now(),
  revoked_at      timestamptz
);
create index if not exists api_keys_merchant_idx on api_keys (merchant_id);
create unique index if not exists api_keys_prefix_uidx on api_keys (key_prefix);
create index if not exists api_keys_secret_hash_idx on api_keys (secret_hash) where secret_hash is not null;

-- payments enhanced with routing + 2FA
create table if not exists payments (
  id                text primary key,
  merchant_id       text not null references merchants(id),
  tx_ref            text not null,
  amount            numeric(20,8) not null check (amount > 0),
  currency          char(3) not null,
  status            text not null check (status in ('created','pending','processing','succeeded','failed','canceled','refunded','partially_refunded')),
  method            text,
  description       text,
  customer_email    text,
  customer_name     text,
  customer_phone    text,
  connector_id      text not null,
  connector_ref     text,
  routing_rule_id   text,
  checkout_url      text,
  return_url        text,
  callback_url      text,
  metadata          jsonb not null default '{}'::jsonb,
  fee_amount        numeric(20,8),
  net_amount        numeric(20,8),
  failure_code      text,
  failure_message   text,
  requires_2fa      boolean not null default false,
  two_fa_verified   boolean not null default false,
  succeeded_at      timestamptz,
  failed_at         timestamptz,
  created_at        timestamptz not null default now(),
  updated_at        timestamptz not null default now(),
  unique (merchant_id, tx_ref)
);
create index if not exists payments_merchant_created_idx on payments (merchant_id, created_at desc);
create index if not exists payments_merchant_status_idx on payments (merchant_id, status);
create index if not exists payments_connector_ref_idx on payments (connector_id, connector_ref);

create table if not exists payment_links (
  id             text primary key,
  merchant_id    text not null references merchants(id),
  payment_id     text references payments(id),
  amount         numeric(20,8) not null check (amount > 0),
  currency       char(3) not null,
  description    text,
  status         text not null check (status in ('active','paid','expired','cancelled')),
  public_token   text not null unique,
  expires_at     timestamptz,
  created_at     timestamptz not null default now(),
  updated_at     timestamptz not null default now()
);
create index if not exists payment_links_merchant_idx on payment_links (merchant_id, created_at desc);

create table if not exists checkout_sessions (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  payment_id      text not null references payments(id),
  payment_link_id text references payment_links(id),
  public_token    text not null unique,
  status          text not null check (status in ('open','completed','expired')),
  selected_method text,
  expires_at      timestamptz not null,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now()
);
create index if not exists checkout_sessions_payment_idx on checkout_sessions (payment_id);

-- ledger enhanced 8 book types
create table if not exists ledger_books (
  id            text primary key,
  merchant_id   text references merchants(id),
  book_type     text not null check (book_type in ('merchant_operating','rail_clearing','platform_revenue','payroll_run','payout_batch','escrow','suspense','reserve','refund_clearing','sandbox')),
  name          text not null,
  currency      char(3) not null default 'ETB',
  status        text not null default 'open' check (status in ('open','closed')),
  created_at    timestamptz not null default now()
);
create index if not exists ledger_books_merchant_idx on ledger_books (merchant_id);
create index if not exists ledger_books_type_idx on ledger_books (book_type);

create table if not exists ledger_accounts (
  id            text primary key,
  book_id       text not null references ledger_books(id),
  code          text not null,
  name          text not null,
  normal_balance text not null check (normal_balance in ('debit','credit')),
  created_at    timestamptz not null default now(),
  unique (book_id, code)
);

create table if not exists ledger_journals (
  id              text primary key,
  book_id         text not null references ledger_books(id),
  posting_key     text not null,
  memo            text,
  transfer_group  text,
  reference_type  text,
  reference_id    text,
  created_at      timestamptz not null default now(),
  unique (book_id, posting_key)
);
create index if not exists ledger_journals_ref_idx on ledger_journals (reference_type, reference_id);
create index if not exists ledger_journals_transfer_idx on ledger_journals (transfer_group);

create table if not exists ledger_entries (
  id            text primary key,
  journal_id    text not null references ledger_journals(id),
  book_id       text not null references ledger_books(id),
  account_id    text not null references ledger_accounts(id),
  direction     text not null check (direction in ('debit','credit')),
  amount        numeric(20,8) not null check (amount > 0),
  currency      char(3) not null,
  meta          jsonb not null default '{}'::jsonb,
  created_at    timestamptz not null default now()
);
create index if not exists ledger_entries_book_created_idx on ledger_entries (book_id, created_at);
create index if not exists ledger_entries_account_idx on ledger_entries (account_id);

create table if not exists ledger_balances (
  book_id     text not null references ledger_books(id),
  account_id  text not null references ledger_accounts(id),
  amount      numeric(20,8) not null default 0,
  updated_at  timestamptz not null default now(),
  primary key (book_id, account_id)
);

create table if not exists outbox_events (
  id              text primary key,
  merchant_id     text references merchants(id),
  aggregate_type  text not null,
  aggregate_id    text not null,
  event_type      text not null,
  payload         jsonb not null,
  created_at      timestamptz not null default now(),
  published_at    timestamptz
);
create index if not exists outbox_unpublished_idx on outbox_events (created_at) where published_at is null;

create table if not exists webhook_endpoints (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  url             text not null,
  secret_hash     text not null,
  secret_prefix   text not null,
  status          text not null check (status in ('active','disabled')),
  events          jsonb not null default '["*"]'::jsonb,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now()
);

create table if not exists webhook_deliveries (
  id              text primary key,
  merchant_id     text not null references merchants(id),
  endpoint_id     text not null references webhook_endpoints(id),
  outbox_event_id text not null references outbox_events(id),
  event_type      text not null,
  payload         jsonb not null,
  status          text not null check (status in ('pending','success','failed','dead')),
  attempt_count   int not null default 0,
  next_attempt_at timestamptz,
  last_status_code int,
  last_error      text,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now()
);
create index if not exists webhook_deliveries_poll_idx on webhook_deliveries (status, next_attempt_at) where status in ('pending','failed');

create table if not exists idempotency_keys (
  merchant_id     text not null references merchants(id),
  key             text not null,
  request_hash    text not null,
  response_code   int not null,
  response_body   jsonb not null,
  resource_type   text,
  resource_id     text,
  created_at      timestamptz not null default now(),
  primary key (merchant_id, key)
);

create table if not exists agent_runs (
  id            text primary key,
  merchant_id   text not null references merchants(id),
  thread_id     text,
  swarm_session_id text,
  input_text    text not null,
  intent        text,
  tool_calls    jsonb not null default '[]'::jsonb,
  output_text   text not null,
  model         text not null default 'rules_v1',
  created_at    timestamptz not null default now()
);
create index if not exists agent_runs_merchant_idx on agent_runs (merchant_id, created_at desc);

create table if not exists audit_logs (
  id            text primary key,
  merchant_id   text,
  actor_type    text not null,
  actor_id      text,
  action        text not null,
  resource_type text,
  resource_id   text,
  ip            inet,
  request_id    text,
  data          jsonb not null default '{}'::jsonb,
  created_at    timestamptz not null default now()
);
create index if not exists audit_logs_merchant_created_idx on audit_logs (merchant_id, created_at desc);
create index if not exists audit_logs_resource_idx on audit_logs (resource_type, resource_id);
