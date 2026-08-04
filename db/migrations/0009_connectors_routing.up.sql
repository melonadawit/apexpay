-- 0009_connectors_routing: smart routing + connector health + circuit breaker support

create table connector_configs (
  id              text primary key,
  connector_id    text not null check (connector_id in ('mock','telebirr','cbe_birr','bank_ips','ethswitch','card_acquirer')),
  merchant_id     text references merchants(id) on delete cascade,
  environment     text not null check (environment in ('test','live')) default 'test',
  config          jsonb not null, -- encrypted secrets via AES-GCM, never plain in logs
  enabled         boolean not null default true,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now(),
  unique (connector_id, merchant_id, environment)
);
create index connector_configs_merchant_idx on connector_configs (merchant_id);
create index connector_configs_connector_idx on connector_configs (connector_id, environment) where enabled=true;

create table connector_health_samples (
  id              text primary key,
  connector_id    text not null,
  environment     text not null default 'live' check (environment in ('test','live')),
  latency_ms      int not null check (latency_ms >=0),
  success         boolean not null,
  error_code      text,
  sampled_at      timestamptz not null default now()
);
create index health_sample_idx on connector_health_samples (connector_id, sampled_at desc);
create index health_sample_env_idx on connector_health_samples (environment, sampled_at desc);
create index health_sample_time_idx on connector_health_samples (sampled_at desc);

create table routing_rules (
  id                text primary key,
  merchant_id       text references merchants(id) on delete cascade,
  name              text not null,
  description       text,
  min_amount        numeric(20,8) check (min_amount >=0),
  max_amount        numeric(20,8) check (max_amount >=0),
  currency          char(3) not null default 'ETB',
  payment_method    text check (payment_method in ('telebirr','cbe_birr','bank','card','qr','mobile_money')),
  primary_connector text not null,
  fallback1         text,
  fallback2         text,
  strategy          text not null default 'success_rate' check (strategy in ('success_rate','latency','cost','round_robin')),
  enabled           boolean not null default true,
  priority          int not null default 100 check (priority >=0),
  created_at        timestamptz not null default now(),
  updated_at        timestamptz not null default now(),
  check (primary_connector != fallback1 and primary_connector != fallback2 and (fallback1 is null or fallback1 != fallback2))
);
create index routing_merchant_idx on routing_rules (merchant_id, priority asc) where enabled=true;
create index routing_global_idx on routing_rules (priority asc) where merchant_id is null and enabled=true;

-- Seed default routing rules for outstanding demo
insert into routing_rules (id, merchant_id, name, min_amount, max_amount, currency, primary_connector, fallback1, fallback2, strategy, priority) values
('route_small', null, 'Small amounts <1000 ETB telebirr primary', 0, 1000, 'ETB', 'telebirr', 'mock', 'cbe_birr', 'success_rate', 10),
('route_medium', null, 'Medium 1000-50000 ETB success_rate', 1000, 50000, 'ETB', 'telebirr', 'cbe_birr', 'mock', 'success_rate', 20),
('route_large', null, 'Large >50000 ETB bank primary', 50000, null, 'ETB', 'bank_ips', 'telebirr', 'mock', 'cost', 30),
('route_qr', null, 'QR EthSwitch primary', null, null, 'ETB', 'ethswitch', 'telebirr', 'mock', 'latency', 40)
on conflict (id) do nothing;

comment on table routing_rules is 'Smart routing: health sampler 30s, circuit breaker 5 fails open 60s, strategy success_rate/latency/cost/round_robin, priority sort O(n log n), fallback chain recorded in routing decision audit.';
