-- 0010_swarm_recon: multi-agent swarm sessions + recon + push devices for Flutter FCM

create table swarm_sessions (
  id                    text primary key,
  merchant_id           text not null references merchants(id) on delete cascade,
  user_id               text references users(id) on delete set null,
  goal                  text not null,
  plan                  jsonb not null default '[]'::jsonb, -- array of PlanStep
  status                text not null check (status in ('planning','executing','needs_confirmation','completed','failed','cancelled')) default 'planning',
  confirmation_required boolean not null default false,
  confirmation_data     jsonb not null default '{}'::jsonb,
  final_output          text,
  created_at            timestamptz not null default now(),
  updated_at            timestamptz not null default now()
);
create index swarm_merchant_idx on swarm_sessions (merchant_id, created_at desc);
create index swarm_status_idx on swarm_sessions (status);

-- agent_runs already exists, add swarm_session_id fk if not exists (migration safe)
do $$ begin
  alter table agent_runs add column if not exists swarm_session_id text references swarm_sessions(id) on delete set null;
exception when others then null; end $$;

create table recon_statements (
  id              text primary key,
  connector_id    text not null,
  statement_date  date not null,
  environment     text not null default 'live' check (environment in ('test','live')),
  raw_file_ref    text not null, -- MinIO
  raw_file_hash   text not null,
  parsed_json     jsonb not null,
  total_amount    numeric(20,8) not null check (total_amount >=0),
  total_count     int not null check (total_count >=0),
  status          text not null check (status in ('pending','parsed','matched','has_breaks')) default 'pending',
  created_at      timestamptz not null default now()
);
create index recon_statements_connector_date_idx on recon_statements (connector_id, statement_date desc);
create unique index recon_statements_file_hash_uidx on recon_statements (raw_file_hash);

create table recon_breaks (
  id              text primary key,
  statement_id    text references recon_statements(id) on delete cascade,
  ledger_book_id  text references ledger_books(id) on delete set null,
  reference_type  text not null check (reference_type in ('payment','refund','payout','payroll')),
  reference_id    text not null,
  expected_amount numeric(20,8) not null,
  actual_amount   numeric(20,8) not null,
  difference      numeric(20,8) not null,
  status          text not null check (status in ('open','investigating','resolved','written_off')) default 'open',
  assigned_to     text references users(id),
  resolution_notes text,
  created_at      timestamptz not null default now(),
  resolved_at     timestamptz
);
create index recon_breaks_statement_idx on recon_breaks (statement_id);
create index recon_breaks_status_idx on recon_breaks (status) where status='open';
create index recon_breaks_reference_idx on recon_breaks (reference_type, reference_id);

create table push_devices (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  user_id         text not null references users(id) on delete cascade,
  platform        text not null check (platform in ('android','ios','web')),
  fcm_token       text not null,
  device_info     jsonb not null default '{}'::jsonb,
  last_active_at  timestamptz not null default now(),
  created_at      timestamptz not null default now()
);
create index push_devices_user_idx on push_devices (user_id);
create index push_devices_merchant_idx on push_devices (merchant_id);
create unique index push_devices_fcm_token_uidx on push_devices (fcm_token);

comment on table swarm_sessions is 'Swarm multi-agent: planner/critic/executor, confirmation thresholds >100k, tool allowlist, audited agent_runs';
comment on table recon_breaks is 'Recon matching engine: tolerance 0.01 ETB, window 24h, suspense posting for breaks';
