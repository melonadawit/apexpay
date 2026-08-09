alter table payroll_claims drop constraint if exists payroll_claims_status_check;
alter table payroll_claims add constraint payroll_claims_status_check
  check (status in ('pending','approved','rejected','paid'));
alter table payroll_claims drop column if exists updated_at;
