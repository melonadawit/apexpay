-- Down migration 0013 payroll comprehensive
drop table if exists payroll_audit_logs;
drop table if exists payroll_employee_portal_access;
drop table if exists payroll_final_settlements;
drop table if exists payroll_compliance_reports;
drop table if exists payroll_loan_repayments;
drop table if exists payroll_loans;
drop table if exists payroll_variable_inputs;
drop table if exists payroll_attendance_inputs;
drop table if exists payroll_salary_revisions;
-- Alter employees drop added columns (safe idempotent)
do $$ begin
  alter table employees drop column if exists department_id;
  alter table employees drop column if exists designation_id;
  alter table employees drop column if exists grade_id;
  alter table employees drop column if exists branch_id;
  alter table employees drop column if exists reporting_manager_id;
  alter table employees drop column if exists ctc_annual;
  alter table employees drop column if exists ctc_monthly;
  alter table employees drop column if exists salary_structure_id;
  alter table employees drop column if exists probation_end_date;
  alter table employees drop column if exists confirmation_status;
  alter table employees drop column if exists employment_status_ext;
  alter table employees drop column if exists date_of_joining;
  alter table employees drop column if exists bank_account_name;
  alter table employees drop column if exists joining_date;
  alter table employees drop column if exists documents;
  alter table employees drop column if exists nationality;
  alter table employees drop column if exists gender;
  alter table employees drop column if exists address;
  alter table employees drop column if exists city;
  alter table employees drop column if exists region;
  alter table employees drop column if exists is_fayda_verified;
  alter table employees drop column if exists fayda_verified_at;
  alter table employees drop column if exists employment_history;
  alter table employees drop column if exists updated_at;
exception when others then null; end $$;

do $$ begin
  alter table payroll_runs drop column if exists pay_calendar_id;
  alter table payroll_runs drop column if exists cutoff_date;
  alter table payroll_runs drop column if exists disbursal_date;
  alter table payroll_runs drop column if exists payroll_data;
  alter table payroll_runs drop column if exists variance_report;
  alter table payroll_runs drop column if exists bank_file_key;
  alter table payroll_runs drop column if exists bank_file_hash;
  alter table payroll_runs drop column if exists employer_total_pension;
  alter table payroll_runs drop column if exists total_employer_cost;
  alter table payroll_runs drop column if exists total_employees_paid;
  alter table payroll_runs drop column if exists total_employees_failed;
  alter table payroll_runs drop column if exists locked_at;
  alter table payroll_runs drop column if exists updated_at;
exception when others then null; end $$;

do $$ begin
  alter table payroll_items drop column if exists earnings_breakdown;
  alter table payroll_items drop column if exists deductions_breakdown;
  alter table payroll_items drop column if exists employer_contributions_breakdown;
  alter table payroll_items drop column if exists ytd;
  alter table payroll_items drop column if exists paid_days;
  alter table payroll_items drop column if exists lop_days;
  alter table payroll_items drop column if exists proration_factor;
  alter table payroll_items drop column if exists ctc_monthly;
  alter table payroll_items drop column if exists is_on_hold;
  alter table payroll_items drop column if exists hold_reason;
  alter table payroll_items drop column if exists updated_at;
exception when others then null; end $$;

drop table if exists payroll_structure_components;
drop table if exists payroll_salary_structures;
drop table if exists payroll_branches;
drop table if exists payroll_grades;
drop table if exists payroll_designations;
drop table if exists payroll_departments;
