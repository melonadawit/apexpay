-- Apex Assistant: conversation threads + append-only message log.
create table assistant_threads (
  id           text primary key,
  user_id      text not null references users(id) on delete cascade,
  merchant_id  text not null references merchants(id) on delete cascade,
  actor        text not null default 'merchant' check (actor in ('merchant','employee')),
  title        text not null default '',
  created_at   timestamptz not null default now()
);

create index assistant_threads_user_idx on assistant_threads (user_id, created_at desc);

-- Append-only: messages are never updated or deleted.
create table assistant_messages (
  id          text primary key,
  thread_id   text not null references assistant_threads(id) on delete cascade,
  role        text not null check (role in ('user','assistant')),
  content     text not null,
  intent      text not null default '',
  tools_used  jsonb not null default '[]'::jsonb,
  data        text not null default '',
  created_at  timestamptz not null default now()
);

create index assistant_messages_thread_idx on assistant_messages (thread_id, created_at);
