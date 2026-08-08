-- 0031: Team (multi-user roles) + generic approvals inbox.
-- merchant_members (from 0001) already stores members/roles; this adds a unified
-- approval workflow table that any module (payouts, payroll, refunds, banking) can use.

create table approval_requests (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  resource_type text not null, -- payout, payroll_run, refund, transfer, invoice, expense, card, ...
  resource_id   text not null,
  action        text not null, -- approve, reject, disbursal, ...
  summary       text,          -- human-readable "Approve payout batch X for ETB Y"
  amount        numeric(20,8) not null default 0,
  currency      char(3) not null default 'ETB',
  status        text not null default 'pending' check (status in ('pending','approved','rejected','cancelled')),
  requested_by  text references users(id),
  -- Maker-checker: a dual-approval request needs 2 distinct approvers.
  required_approvals int not null default 1,
  approvals     jsonb not null default '[]'::jsonb, -- [{user_id, name, role, decided_at, decision}]
  metadata      jsonb not null default '{}'::jsonb,
  created_at    timestamptz not null default now(),
  decided_at    timestamptz,
  updated_at    timestamptz not null default now()
);
create index approval_requests_merchant_status_idx on approval_requests (merchant_id, status, created_at);
create index approval_requests_resource_idx on approval_requests (resource_type, resource_id);
