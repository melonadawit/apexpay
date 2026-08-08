-- 0026: payroll maker-checker support.
-- Tracks who created a run and each approval, so the >100k-net dual-approval
-- control can be enforced (approver != submitter, second distinct approver).

alter table payroll_runs add column if not exists created_by text references users(id);

create table if not exists payroll_approvals (
  id            text primary key,
  run_id        text not null references payroll_runs(id) on delete cascade,
  merchant_id   text not null references merchants(id),
  approver_id   text references users(id),
  approver_type text not null default 'finance',
  from_status   text not null,
  to_status     text not null,
  comments      text,
  created_at    timestamptz not null default now()
);
create index if not exists payroll_approvals_run_idx on payroll_approvals (run_id, created_at);
create index if not exists payroll_approvals_merchant_idx on payroll_approvals (merchant_id);
