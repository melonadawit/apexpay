-- 0033: Analytics/cohort + real TOTP 2FA challenge storage.

-- Analytics: revenue, success by method/device, subscription cohort/retention.
create table analytics_daily (
  merchant_id     text not null references merchants(id) on delete cascade,
  stat_date       date not null,
  revenue         numeric(20,8) not null default 0,
  tpv             numeric(20,8) not null default 0,
  payment_count   int not null default 0,
  success_count   int not null default 0,
  failed_count    int not null default 0,
  refund_amount   numeric(20,8) not null default 0,
  method_breakdown jsonb not null default '{}'::jsonb, -- {telebirr: {count, revenue}, ...}
  device_breakdown jsonb not null default '{}'::jsonb,
  created_at      timestamptz not null default now(),
  primary key (merchant_id, stat_date)
);

create table subscription_cohorts (
  cohort_month   date not null,           -- first month customers joined
  merchant_id    text not null references merchants(id) on delete cascade,
  customers      int not null default 0,
  month1_retention numeric(5,2) not null default 0, -- % retained after 1 month
  month2_retention numeric(5,2) not null default 0,
  month3_retention numeric(5,2) not null default 0,
  mrr            numeric(20,8) not null default 0,
  primary key (cohort_month, merchant_id)
);

-- 2FA challenges: server-side challenge state for real TOTP/OOB flows.
create table twofa_challenges (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  user_id       text references users(id) on delete set null,
  payment_id    text references payments(id) on delete cascade,
  provider      text not null default 'totp', -- totp | sms | email | oob
  secret        text,                         -- encrypted TOTP secret / OOB code hash
  challenge     text,                         -- HMAC-signed challenge for TOTP
  expires_at    timestamptz not null,
  verified_at   timestamptz,
  created_at    timestamptz not null default now()
);
create index twofa_challenges_payment_idx on twofa_challenges (payment_id);
