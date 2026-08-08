-- 0025: merchant dashboard session auth.
-- Opaque, revocable, DB-backed sessions for dashboard users (not the server-to-server
-- API keys). Only the SHA-256 digest of a session token is stored, never the raw token.

create table auth_sessions (
  id              text primary key,
  user_id         text not null references users(id) on delete cascade,
  merchant_id     text not null references merchants(id) on delete cascade,
  token_hash      text not null,           -- sha256 hex of the opaque token
  user_agent      text,
  ip              inet,
  expires_at      timestamptz not null,
  last_active_at  timestamptz not null default now(),
  revoked_at      timestamptz,
  created_at      timestamptz not null default now()
);

create unique index auth_sessions_token_hash_idx on auth_sessions (token_hash);
create index auth_sessions_user_idx on auth_sessions (user_id, expires_at);
create index auth_sessions_merchant_idx on auth_sessions (merchant_id);
