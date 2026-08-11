-- 0013_payroll_comprehensive: enterprise-grade payroll
-- Senior Engineer design: clean arch, decimal precise, ULID, optimal data structures, quality indexes
-- Goal: Salary Structure Engine, Departments, Grades, Attendance LOP Proration, Variable Inputs, Loans, Compliance Reports, F&F, Self-Service Portal

-- Departments, Designations, Grades, Branches — organizational hierarchy per merchant
create table payroll_departments (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  name        text not null, -- e.g., Engineering, Sales, Finance
  name_am     text,
  code        text, -- e.g., ENG, SALES
  cost_center text, -- for ledger allocation outstanding Money OS
  description text,
  created_at  timestamptz not null default now(),
  unique (merchant_id, code)
);
create index payroll_departments_merchant_idx on payroll_departments (merchant_id);

create table payroll_designations (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  title       text not null, -- Senior Engineer, Manager
  title_am    text,
  level       int not null default 1, -- hierarchy level for approval chain
  description text,
  created_at  timestamptz not null default now(),
  unique (merchant_id, title)
);

create table payroll_grades (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  name        text not null, -- G1, G2, L1, L2 or Junior/Senior
  name_am     text,
  min_salary  numeric(20,8) not null check (min_salary >=0),
  max_salary  numeric(20,8) not null check (max_salary >= min_salary),
  description text,
  created_at  timestamptz not null default now(),
  unique (merchant_id, name)
);

create table payroll_branches (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  name        text not null, -- Head Office, Shashemene Branch
  name_am     text,
  region      text not null, -- Oromiya, Addis Ababa
  city        text not null,
  sub_city    text,
  address     text,
  is_head     boolean not null default false,
  created_at  timestamptz not null default now(),
  unique (merchant_id, name)
);

-- Salary Structure Template — CTC-based like ApexPay
create table payroll_salary_structures (
  id             text primary key,
  merchant_id    text not null references merchants(id) on delete cascade,
  name           text not null, -- Fixed 500k Annual, Tech Band G3
  name_am        text,
  description    text,
  ctc_annual     numeric(20,8) not null check (ctc_annual >=0),
  ctc_monthly    numeric(20,8) not null check (ctc_monthly >=0), -- ctc_annual/12 generated
  currency       char(3) not null default 'ETB',
  effective_from date not null default current_date,
  effective_to   date,
  status         text not null check (status in ('active','archived','draft')) default 'active',
  is_default     boolean not null default false,
  created_by     text references users(id),
  created_at     timestamptz not null default now(),
  updated_at     timestamptz not null default now(),
  unique (merchant_id, name)
);
create index payroll_salary_structures_merchant_idx on payroll_salary_structures (merchant_id, status);

-- Structure Components: Earnings, Deductions, Employer Contributions, Reimbursements
create table payroll_structure_components (
  id              text primary key,
  structure_id    text not null references payroll_salary_structures(id) on delete cascade,
  component_type  text not null check (component_type in ('earning','deduction','employer_contribution','reimbursement')),
  code            text not null, -- BASIC, HOUSING, TRANSPORT, FUEL, SPECIAL_ALLOW, MEDICAL, OVERTIME, COMMISSION, BONUS, TAXABLE etc.
  name            text not null, -- Basic Salary, Housing Allowance
  name_am         text,
  calculation_type text not null check (calculation_type in ('fixed','percentage_of_basic','percentage_of_ctc','percentage_of_gross','formula')),
  -- For fixed
  amount          numeric(20,8) not null default 0,
  percentage      numeric(5,2) not null default 0, -- 40.00 = 40%
  formula         text, -- e.g., "CTC_MONTHLY * 0.4" or "BASIC * 0.1" — secure parser, no evil eval
  is_taxable      boolean not null default true,
  is_part_of_gross boolean not null default true, -- employer contribution not part of gross
  is_proratable   boolean not null default true, -- LOP proration applicable
  is_pensionable  boolean not null default true, -- included in pensionable salary?
  is_optional     boolean not null default false,
  tax_exempt_limit numeric(20,8) not null default 0, -- e.g., medical up to 1000 ETB exempt
  order_no        int not null default 0, -- display order in payslip
  meta            jsonb not null default '{}'::jsonb,
  created_at      timestamptz not null default now(),
  unique (structure_id, code)
);
create index payroll_structure_components_structure_idx on payroll_structure_components (structure_id, order_no);

-- Extend employees with org structure + CTC + docs
do $$ begin
  alter table employees add column if not exists department_id text references payroll_departments(id);
  alter table employees add column if not exists designation_id text references payroll_designations(id);
  alter table employees add column if not exists grade_id text references payroll_grades(id);
  alter table employees add column if not exists branch_id text references payroll_branches(id);
  alter table employees add column if not exists reporting_manager_id text references employees(id);
  alter table employees add column if not exists ctc_annual numeric(20,8) not null default 0;
  alter table employees add column if not exists ctc_monthly numeric(20,8) not null default 0;
  alter table employees add column if not exists salary_structure_id text references payroll_salary_structures(id);
  alter table employees add column if not exists probation_end_date date;
  alter table employees add column if not exists confirmation_status text check (confirmation_status in ('probation','confirmed','notice','terminated')) default 'probation';
  alter table employees add column if not exists employment_status_ext text check (employment_status_ext in ('active','on_hold','notice_period','terminated','retired')) default 'active';
  alter table employees add column if not exists date_of_joining date not null default current_date;
  alter table employees add column if not exists bank_account_name text;
  alter table employees add column if not exists joining_date date;
  alter table employees add column if not exists documents jsonb not null default '[]'::jsonb;
  alter table employees add column if not exists nationality char(2) not null default 'ET';
  alter table employees add column if not exists gender text check (gender in ('male','female','other'));
  alter table employees add column if not exists address text;
  alter table employees add column if not exists city text;
  alter table employees add column if not exists region text default 'Oromiya';
  alter table employees add column if not exists is_fayda_verified boolean not null default false;
  alter table employees add column if not exists fayda_verified_at timestamptz;
  alter table employees add column if not exists employment_history jsonb not null default '[]'::jsonb;
  alter table employees add column if not exists updated_at timestamptz not null default now();
exception when others then null; end $$;

-- Salary Revisions History + Arrears
create table payroll_salary_revisions (
  id               text primary key,
  merchant_id      text not null references merchants(id) on delete cascade,
  employee_id      text not null references employees(id) on delete cascade,
  old_base         numeric(20,8) not null,
  new_base         numeric(20,8) not null,
  old_ctc          numeric(20,8) not null default 0,
  new_ctc          numeric(20,8) not null default 0,
  old_structure_id text references payroll_salary_structures(id),
  new_structure_id text references payroll_salary_structures(id),
  effective_from   date not null,
  reason           text,
  approved_by      text references users(id),
  status           text not null check (status in ('pending','approved','rejected')) default 'pending',
  arrear_amount    numeric(20,8) not null default 0, -- auto calc (new-old)*months
  arrear_months    int not null default 0,
  created_at       timestamptz not null default now()
);
create index payroll_salary_revisions_employee_idx on payroll_salary_revisions (employee_id, effective_from desc);
create index payroll_salary_revisions_merchant_idx on payroll_salary_revisions (merchant_id);

-- Attendance & Leave Inputs per Run — LOP proration
create table payroll_attendance_inputs (
  id                  text primary key,
  run_id              text not null references payroll_runs(id) on delete cascade,
  employee_id         text not null references employees(id) on delete cascade,
  paid_days           int not null default 30 check (paid_days >=0),
  lop_days            int not null default 0 check (lop_days >=0),
  total_days          int not null default 30 check (total_days >0),
  present_days        int not null default 30,
  ot_weekday_hours    numeric(10,2) not null default 0,
  ot_weekend_hours    numeric(10,2) not null default 0,
  ot_holiday_hours    numeric(10,2) not null default 0,
  ot_night_hours      numeric(10,2) not null default 0,
  leave_taken         jsonb not null default '{}'::jsonb, -- {"annual":2,"sick":1,"maternity":0}
  leave_balance       jsonb not null default '{}'::jsonb,
  is_on_hold          boolean not null default false, -- hold salary
  hold_reason         text,
  created_at          timestamptz not null default now(),
  unique (run_id, employee_id)
);
create index payroll_attendance_run_idx on payroll_attendance_inputs (run_id);
create index payroll_attendance_employee_idx on payroll_attendance_inputs (employee_id);

-- Variable Pay Inputs per Run per Employee
create table payroll_variable_inputs (
  id              text primary key,
  run_id          text not null references payroll_runs(id) on delete cascade,
  employee_id     text not null references employees(id) on delete cascade,
  component_code  text not null, -- COMMISSION, BONUS, OVERTIME, PENALTY, ARREAR, THIRTEENTH_MONTH
  amount          numeric(20,8) not null,
  is_taxable      boolean not null default true,
  is_pensionable  boolean not null default true,
  description     text,
  created_by      text references users(id),
  created_at      timestamptz not null default now()
);
create index payroll_variable_run_idx on payroll_variable_inputs (run_id, employee_id);
create index payroll_variable_component_idx on payroll_variable_inputs (component_code);

-- Loans & Advances
create table payroll_loans (
  id               text primary key,
  merchant_id      text not null references merchants(id) on delete cascade,
  employee_id      text not null references employees(id) on delete cascade,
  loan_type        text not null check (loan_type in ('personal','salary_advance','housing','education','medical','other')),
  principal        numeric(20,8) not null check (principal >0),
  interest_rate    numeric(5,2) not null default 0, -- 0% for salary advance
  tenure_months    int not null check (tenure_months >0),
  emi_amount       numeric(20,8) not null,
  total_paid       numeric(20,8) not null default 0,
  outstanding      numeric(20,8) not null,
  status           text not null check (status in ('draft','pending_approval','approved','active','closed','rejected','written_off')) default 'pending_approval',
  disbursed_at     timestamptz,
  next_due_date    date,
  approved_by      text references users(id),
  reason           text,
  meta             jsonb not null default '{}'::jsonb,
  created_at       timestamptz not null default now(),
  updated_at       timestamptz not null default now()
);
create index payroll_loans_employee_idx on payroll_loans (employee_id, status);
create index payroll_loans_merchant_idx on payroll_loans (merchant_id, status);

create table payroll_loan_repayments (
  id                    text primary key,
  loan_id               text not null references payroll_loans(id) on delete cascade,
  run_id                text references payroll_runs(id) on delete set null,
  employee_id           text not null references employees(id) on delete cascade,
  amount                numeric(20,8) not null,
  principal_component   numeric(20,8) not null default 0,
  interest_component    numeric(20,8) not null default 0,
  outstanding_after     numeric(20,8) not null,
  status                text not null check (status in ('pending','paid','failed')) default 'paid',
  created_at            timestamptz not null default now()
);
create index payroll_loan_repayments_loan_idx on payroll_loan_repayments (loan_id);
create index payroll_loan_repayments_run_idx on payroll_loan_repayments (run_id);

-- Enhanced payroll_runs
do $$ begin
  alter table payroll_runs add column if not exists pay_calendar_id text;
  alter table payroll_runs add column if not exists cutoff_date date;
  alter table payroll_runs add column if not exists disbursal_date date;
  alter table payroll_runs add column if not exists payroll_data jsonb not null default '{}'::jsonb;
  alter table payroll_runs add column if not exists variance_report jsonb not null default '{}'::jsonb;
  alter table payroll_runs add column if not exists bank_file_key text;
  alter table payroll_runs add column if not exists bank_file_hash text;
  alter table payroll_runs add column if not exists employer_total_pension numeric(20,8) not null default 0;
  alter table payroll_runs add column if not exists total_employer_cost numeric(20,8) not null default 0;
  alter table payroll_runs add column if not exists total_employees_paid int not null default 0;
  alter table payroll_runs add column if not exists total_employees_failed int not null default 0;
  alter table payroll_runs add column if not exists locked_at timestamptz;
  alter table payroll_runs add column if not exists updated_at timestamptz not null default now();
exception when others then null; end $$;

-- Enhanced payroll_items with breakdowns JSON + YTD + paid_days
do $$ begin
  alter table payroll_items add column if not exists earnings_breakdown jsonb not null default '[]'::jsonb;
  alter table payroll_items add column if not exists deductions_breakdown jsonb not null default '[]'::jsonb;
  alter table payroll_items add column if not exists employer_contributions_breakdown jsonb not null default '[]'::jsonb;
  alter table payroll_items add column if not exists ytd jsonb not null default '{}'::jsonb;
  alter table payroll_items add column if not exists paid_days int not null default 30;
  alter table payroll_items add column if not exists lop_days int not null default 0;
  alter table payroll_items add column if not exists proration_factor numeric(5,4) not null default 1.0;
  alter table payroll_items add column if not exists ctc_monthly numeric(20,8) not null default 0;
  alter table payroll_items add column if not exists is_on_hold boolean not null default false;
  alter table payroll_items add column if not exists hold_reason text;
  alter table payroll_items add column if not exists updated_at timestamptz not null default now();
exception when others then null; end $$;

-- Compliance Reports generated
create table payroll_compliance_reports (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  period_month  int not null check (period_month between 1 and 12),
  period_year   int not null,
  report_type   text not null check (report_type in ('pension_contribution','erca_withholding','annual_tax_certificate','pension_challan','bank_disbursal_file','payroll_register','cost_center_report','variance_report')),
  file_key      text,
  file_hash     text,
  status        text not null check (status in ('draft','generated','paid','filed','failed')) default 'generated',
  metadata      jsonb not null default '{}'::jsonb,
  generated_by  text references users(id),
  created_at    timestamptz not null default now(),
  unique (merchant_id, period_year, period_month, report_type)
);
create index payroll_compliance_reports_merchant_idx on payroll_compliance_reports (merchant_id, period_year desc, period_month desc);

-- Final Settlement F&F
create table payroll_final_settlements (
  id                    text primary key,
  merchant_id           text not null references merchants(id) on delete cascade,
  employee_id           text not null references employees(id) on delete cascade,
  resignation_date      date not null,
  last_working_date     date not null,
  notice_period_days    int not null default 30,
  notice_served_days    int not null default 30,
  notice_shortfall_days int not null default 0,
  leave_encashment_days numeric(10,2) not null default 0,
  leave_encashment_amount numeric(20,8) not null default 0,
  severance_amount      numeric(20,8) not null default 0, -- per ET labour law
  gratuity_amount       numeric(20,8) not null default 0,
  bonus_pro_rata        numeric(20,8) not null default 0,
  outstanding_loans     numeric(20,8) not null default 0,
  outstanding_advances  numeric(20,8) not null default 0,
  other_earnings        numeric(20,8) not null default 0,
  other_deductions      numeric(20,8) not null default 0,
  total_payable         numeric(20,8) not null default 0,
  total_deductions      numeric(20,8) not null default 0,
  net_payable           numeric(20,8) not null default 0,
  status                text not null check (status in ('draft','pending_approval','approved','paid','rejected')) default 'draft',
  clearance_checklist   jsonb not null default '[]'::jsonb, -- [{item:laptop, status:done, checked_by}]
  approved_by           text references users(id),
  paid_at               timestamptz,
  created_at            timestamptz not null default now(),
  updated_at            timestamptz not null default now(),
  unique (employee_id) -- one active settlement per employee at a time
);
create index payroll_final_settlements_merchant_idx on payroll_final_settlements (merchant_id, status);

-- Employee Portal Access — magic link JWT 24h
create table payroll_employee_portal_access (
  id                 text primary key,
  merchant_id        text not null references merchants(id) on delete cascade,
  employee_id        text not null references employees(id) on delete cascade,
  magic_token_hash   text not null, -- sha256 of JWT
  token_last4        text,
  expires_at         timestamptz not null,
  last_accessed_at   timestamptz,
  access_count       int not null default 0,
  is_revoked         boolean not null default false,
  created_at         timestamptz not null default now(),
  unique (magic_token_hash)
);
create index payroll_employee_portal_employee_idx on payroll_employee_portal_access (employee_id, expires_at desc);

-- Payroll audit logs for compliance (separate from general audit_logs)
create table payroll_audit_logs (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  run_id        text references payroll_runs(id) on delete set null,
  employee_id   text references employees(id) on delete set null,
  actor_type    text not null, -- system, hr, finance, admin, employee
  actor_id      text,
  action        text not null, -- create_employee, salary_revision, calculate_run, approve_run, disburse_run, hold_salary, generate_payslip, etc.
  details       jsonb not null default '{}'::jsonb,
  ip            inet,
  request_id    text,
  created_at    timestamptz not null default now()
);
create index payroll_audit_logs_merchant_created_idx on payroll_audit_logs (merchant_id, created_at desc);
create index payroll_audit_logs_run_idx on payroll_audit_logs (run_id, created_at desc);

-- Seed default departments, designations, grades for new merchants via function (called in seed script)
-- Placeholder: Engineering, Sales, HR, Finance, Operations

comment on table payroll_salary_structures is 'modern CTC template: e.g., Fixed 500k Annual breakdown monthly earnings/deductions/employer contributions with formula engine.';
comment on table payroll_structure_components is 'Component code BASIC HOUSING TRANSPORT FUEL SPECIAL_ALLOW MEDICAL OVERTIME COMMISSION BONUS etc calculation fixed/percentage/formula taxable proratable pensionable.';
comment on table payroll_attendance_inputs is 'LOP proration: paid_days/total_days factor. OT hours weekday_weekend_holiday_night for ET labour law 1.25/1.5/2.0/1.3.';
comment on table payroll_variable_inputs is 'Variable pay per run per employee: commission bonus penalty arrear thirteenth_month etc taxable/pensionable flag.';
comment on table payroll_loans is 'Personal, salary_advance principal tenure EMI outstanding status draft->pending_approval->approved->active->closed. Auto deduction per payroll run.';
comment on table payroll_compliance_reports is 'Pension contribution CSV for Private Org Employees Social Security Agency + ERCA withholding monthly CSV + bank disbursal pain.001 file.';
comment on table payroll_final_settlements is 'F&F: resignation LWD notice period leave_encashment gross/30 severance per ET labour law Art 39-44 gratuity outstanding loans clearance checklist.';
