-- 0028: Risk & Fraud Engine.
-- Rule-based transaction monitoring, velocity checks, and manual-review flags.

create table risk_rules (
  id            text primary key,
  merchant_id   text references merchants(id) on delete cascade, -- NULL = platform-global rule
  name          text not null,
  description   text,
  rule_type     text not null check (rule_type in (
                  'velocity_amount','velocity_count','threshold_amount','threshold_count',
                  'new_device','new_ip','high_ticket','high_failure_rate','manual'
                )),
  parameters    jsonb not null default '{}'::jsonb, -- {window_minutes, amount_limit, count_limit, ...}
  action        text not null default 'flag' check (action in ('flag','review','block')),
  severity      text not null default 'medium' check (severity in ('low','medium','high','critical')),
  enabled       boolean not null default true,
  created_at    timestamptz not null default now(),
  updated_at    timestamptz not null default now()
);
create index risk_rules_merchant_idx on risk_rules (merchant_id);
create index risk_rules_enabled_idx on risk_rules (merchant_id, enabled);

create table risk_flags (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  entity_type   text not null check (entity_type in ('payment','refund','payout','merchant','device','ip','customer')),
  entity_id     text not null,
  rule_id       text references risk_rules(id) on delete set null,
  rule_name     text,
  severity      text not null default 'medium',
  action        text not null default 'flag',
  reason        text,
  details       jsonb not null default '{}'::jsonb,
  status        text not null default 'open' check (status in ('open','investigating','resolved','dismissed')),
  assigned_to   text references users(id) on delete set null,
  created_at    timestamptz not null default now(),
  resolved_at   timestamptz,
  resolved_by   text references users(id) on delete set null,
  updated_at    timestamptz not null default now()
);
create index risk_flags_status_idx on risk_flags (merchant_id, status, created_at);
create index risk_flags_entity_idx on risk_flags (entity_type, entity_id);
