-- Down migration 0014
drop table if exists payroll_loan_emi_schedule;
drop table if exists payroll_leave_requests;
drop table if exists payroll_leave_balances;
drop table if exists payroll_calendars;

do $$ begin
  alter table payroll_claims drop column if exists receipt_file_hash;
  alter table payroll_claims drop column if exists approved_by_manager;
  alter table payroll_claims drop column if exists approved_by_finance;
  alter table payroll_claims drop column if exists manager_approved_at;
  alter table payroll_claims drop column if exists finance_approved_at;
  alter table payroll_claims drop column if exists rejection_reason;
  alter table payroll_claims drop column if exists is_taxable;
  alter table payroll_claims drop column if exists is_pensionable;
exception when others then null; end $$;

do $$ begin
  alter table payroll_final_settlements drop column if exists clearance_items_detailed;
  alter table payroll_final_settlements drop column if exists assets_returned;
  alter table payroll_final_settlements drop column if exists exit_interview;
exception when others then null; end $$;

do $$ begin
  alter table payroll_runs drop constraint if exists payroll_runs_calendar_fk;
exception when others then null; end $$;
