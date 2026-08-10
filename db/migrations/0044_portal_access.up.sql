-- Self-service portals: token-gated access for vendors and customers to view their own
-- invoices / payment status without a dashboard login. Tokens are stored as hashes.
create table if not exists portal_access (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  portal_type   text not null check (portal_type in ('vendor','customer')),
  entity_id     text not null,          -- vendor id or customer email (external)
  entity_name   text,
  token_hash    text not null,
  token_last4   text,
  expires_at    timestamptz not null,
  last_accessed_at timestamptz,
  access_count  integer not null default 0,
  is_revoked    boolean not null default false,
  created_at    timestamptz not null default now()
);
create index if not exists portal_access_merchant_idx on portal_access (merchant_id);
create index if not exists portal_access_token_hash_idx on portal_access (token_hash);
