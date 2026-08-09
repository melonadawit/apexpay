-- Fix payroll_claims: add the missing updated_at column and allow the intermediate
-- approved_by_manager status used by the manager -> finance approval flow.
alter table payroll_claims
  add column if not exists updated_at timestamptz not null default now();

alter table payroll_claims drop constraint if exists payroll_claims_status_check;
alter table payroll_claims add constraint payroll_claims_status_check
  check (status in ('pending','approved_by_manager','approved','rejected','paid'));
