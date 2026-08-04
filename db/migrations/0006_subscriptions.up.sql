-- 0006_subscriptions: recurring billing with trials, dunning, customer portal

create table customers (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  email           text,
  phone           text,
  name            text,
  fayda_fin_hash  text, -- optional link to Fayda for customer KYC if needed
  metadata        jsonb not null default '{}'::jsonb,
  created_at      timestamptz not null default now()
);
create index customers_merchant_idx on customers (merchant_id);
create index customers_merchant_email_idx on customers (merchant_id, lower(email)) where email is not null;
create index customers_fayda_hash_idx on customers (fayda_fin_hash) where fayda_fin_hash is not null;

create table subscription_plans (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  name            text not null,
  description     text,
  amount          numeric(20,8) not null check (amount >0),
  currency        char(3) not null default 'ETB',
  interval_type   text not null check (interval_type in ('day','week','month','year')),
  interval_count  int not null default 1 check (interval_count >0),
  trial_days      int not null default 0 check (trial_days >=0),
  status          text not null default 'active' check (status in ('active','archived')),
  created_at      timestamptz not null default now()
);
create index sub_plans_merchant_idx on subscription_plans (merchant_id);

create table subscriptions (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  customer_id         text not null references customers(id) on delete cascade,
  plan_id             text not null references subscription_plans(id),
  status              text not null check (status in ('incomplete','trialing','active','past_due','canceled','paused')),
  current_period_start timestamptz not null,
  current_period_end   timestamptz not null,
  trial_end           timestamptz,
  cancel_at           timestamptz,
  created_at          timestamptz not null default now()
);
create index subscriptions_merchant_status_idx on subscriptions (merchant_id, status);
create index subscriptions_customer_idx on subscriptions (customer_id);
create index subscriptions_period_end_idx on subscriptions (current_period_end) where status in ('trialing','active','past_due');

create table subscription_invoices (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  subscription_id text not null references subscriptions(id) on delete cascade,
  payment_id      text references payments(id) on delete set null,
  amount          numeric(20,8) not null check (amount >=0),
  currency        char(3) not null default 'ETB',
  status          text not null check (status in ('draft','open','paid','uncollectible','void')),
  attempt_count   int not null default 0 check (attempt_count >=0),
  next_attempt_at timestamptz,
  due_at          timestamptz not null,
  created_at      timestamptz not null default now()
);
create index sub_invoices_subscription_idx on subscription_invoices (subscription_id);
create index sub_invoices_due_idx on subscription_invoices (due_at, status) where status in ('open','draft');
create index sub_invoices_next_attempt_idx on subscription_invoices (next_attempt_at) where status='open' and attempt_count <3;

comment on table subscription_invoices is 'Dunning schedule: attempt 0 -> +1d, 1 -> +3d, 2 -> +5d per ET business best practice';
