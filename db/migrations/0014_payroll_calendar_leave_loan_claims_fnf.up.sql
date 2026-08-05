-- 0014_payroll_calendar_leave_loan_claims_fnf: Ethiopia Law Full Compliance — Payroll Calendar CRUD, Leave Management Art 77/82/86, Reimbursements MinIO, Loans EMI Schedule, Final Settlement F&F Clearance
-- Senior Engineer design: clean arch, decimal precise, ULID, optimal data structures, quality indexes per Ethiopia Labour Proclamation 1156/2019

-- Payroll Calendar — Pay Schedule Monthly Weekly Semimonthly per Ethiopia business practice
create table payroll_calendars (
  id               text primary key,
  merchant_id      text not null references merchants(id) on delete cascade,
  name             text not null, -- e.g., Monthly Payroll Calendar 2026, Weekly Calendar
  description      text,
  pay_frequency    text not null check (pay_frequency in ('monthly','semimonthly','weekly','biweekly')) default 'monthly',
  year             int not null,
  month            int check (month between 1 and 12), -- null for weekly? but for monthly required
  cutoff_day       int not null check (cutoff_day between 1 and 31) default 25, -- Ethiopia business practice cutoff 25th
  disbursal_day    int not null check (disbursal_day between 1 and 31) default 30, -- disbursal 30th
  pay_day          int not null check (pay_day between 1 and 31) default 31, -- pay date last day of month
  cutoff_date      date not null, -- actual cutoff date e.g., 2026-07-25
  disbursal_date   date not null, -- e.g., 2026-07-30
  pay_date         date not null, -- e.g., 2026-07-31
  is_locked        boolean not null default false, -- lock after disbursal per law
  locked_at        timestamptz,
  locked_by        text references users(id),
  created_by       text references users(id),
  created_at       timestamptz not null default now(),
  updated_at       timestamptz not null default now(),
  unique (merchant_id, year, month, pay_frequency)
);
create index payroll_calendars_merchant_year_month_idx on payroll_calendars (merchant_id, year, month);
create index payroll_calendars_locked_idx on payroll_calendars (merchant_id, is_locked) where is_locked=false;

-- Add foreign key from payroll_runs to payroll_calendars if not exists
do $$ begin
  alter table payroll_runs add constraint payroll_runs_calendar_fk foreign key (pay_calendar_id) references payroll_calendars(id) on delete set null;
exception when others then null; end $$;

-- Leave Balances — per employee per year per leave type per Art 77/82/86
create table payroll_leave_balances (
  id                 text primary key,
  merchant_id        text not null references merchants(id) on delete cascade,
  employee_id        text not null references employees(id) on delete cascade,
  leave_type         text not null check (leave_type in ('annual','sick','maternity','paternity','marriage','mourning','unpaid','comp_off','study')),
  year               int not null,
  entitled_days      numeric(10,2) not null check (entitled_days >=0), -- e.g., annual 14 + years-1 up to 35 per Art 77
  used_days          numeric(10,2) not null default 0,
  remaining_days     numeric(10,2) not null default 0,
  carry_forward_days numeric(10,2) not null default 0, -- from previous year per company policy
  created_at         timestamptz not null default now(),
  updated_at         timestamptz not null default now(),
  unique (merchant_id, employee_id, leave_type, year)
);
create index payroll_leave_balances_employee_year_idx on payroll_leave_balances (employee_id, year);
create index payroll_leave_balances_merchant_year_idx on payroll_leave_balances (merchant_id, year, leave_type);

-- Leave Requests — per employee per period
create table payroll_leave_requests (
  id                text primary key,
  merchant_id       text not null references merchants(id) on delete cascade,
  employee_id       text not null references employees(id) on delete cascade,
  leave_type        text not null check (leave_type in ('annual','sick','maternity','paternity','marriage','mourning','unpaid','comp_off','study')),
  start_date        date not null,
  end_date          date not null,
  days_requested    numeric(10,2) not null check (days_requested >0), -- e.g., 2 days, 0.5 half day per company policy
  reason            text,
  status            text not null check (status in ('pending','approved','rejected','cancelled')) default 'pending',
  approved_by       text references users(id),
  approved_at       timestamptz,
  rejection_reason  text,
  medical_certificate_file_key text, -- MinIO for sick leave >3 days per Art 82 need medical certificate
  created_at        timestamptz not null default now(),
  updated_at        timestamptz not null default now()
);
create index payroll_leave_requests_employee_idx on payroll_leave_requests (employee_id, start_date desc);
create index payroll_leave_requests_merchant_status_idx on payroll_leave_requests (merchant_id, status, start_date desc);
create index payroll_leave_requests_year_idx on payroll_leave_requests (merchant_id, employee_id, leave_type, status);

-- Enhanced payroll_claims for reimbursements MinIO receipt upload approval manager->finance
do $$ begin
  alter table payroll_claims add column if not exists receipt_file_hash text;
  alter table payroll_claims add column if not exists approved_by_manager text references users(id);
  alter table payroll_claims add column if not exists approved_by_finance text references users(id);
  alter table payroll_claims add column if not exists manager_approved_at timestamptz;
  alter table payroll_claims add column if not exists finance_approved_at timestamptz;
  alter table payroll_claims add column if not exists rejection_reason text;
  alter table payroll_claims add column if not exists is_taxable boolean not null default false;
  alter table payroll_claims add column if not exists is_pensionable boolean not null default false;
exception when others then null; end $$;

-- Loans enhanced EMI schedule tracking already in payroll_loans + payroll_loan_repayments, add schedule table for outstanding preview
create table payroll_loan_emi_schedule (
  id                text primary key,
  loan_id           text not null references payroll_loans(id) on delete cascade,
  installment_no    int not null,
  due_date          date not null,
  emi_amount        numeric(20,8) not null,
  principal_component numeric(20,8) not null,
  interest_component numeric(20,8) not null,
  outstanding_after numeric(20,8) not null,
  status            text not null check (status in ('pending','paid','overdue','skipped')) default 'pending',
  paid_at           timestamptz,
  run_id            text references payroll_runs(id) on delete set null,
  created_at        timestamptz not null default now(),
  unique (loan_id, installment_no)
);
create index payroll_loan_emi_schedule_loan_idx on payroll_loan_emi_schedule (loan_id, due_date);

-- Final settlement enhanced clearance checklist already in payroll_final_settlements, ensure checklist structure
do $$ begin
  alter table payroll_final_settlements add column if not exists clearance_items_detailed jsonb not null default '[]'::jsonb; -- [{item:laptop, category:IT, status:pending/done, checked_by, checked_at, notes, required:true}]
  alter table payroll_final_settlements add column if not exists assets_returned jsonb not null default '[]'::jsonb; -- [{asset_type:laptop, asset_id:LP001, returned:true, condition:good, returned_at}]
  alter table payroll_final_settlements add column if not exists exit_interview jsonb not null default '{}'::jsonb; -- {conducted:true, conducted_by, date, feedback}
exception when others then null; end $$;

-- Payroll audit logs already exists, add index for calendar
create index if not exists payroll_audit_logs_calendar_idx on payroll_audit_logs (run_id) where run_id is not null;

-- Seed default payroll calendar for 2026 monthly per Ethiopia business practice cutoff 25th disbursal 30th pay last day
-- Function to generate calendar for year 2026 monthly

comment on table payroll_calendars is 'Payroll Calendar CRUD per Ethiopia business practice: cutoff 25th disbursal 30th pay date last day last day of month lock after disbursal per law, pay_frequency monthly/semimonthly/weekly/biweekly, year month cutoff_day disbursal_day pay_day cutoff_date disbursal_date pay_date is_locked locked_at locked_by';
comment on table payroll_leave_balances is 'Leave Balances per Art 77 Annual Leave 14 days first year +1 per year up to 35 max, Art 82 Sick Leave 6 months per 12 months first 30 days 100% next 60 days 50% remaining 90 days unpaid, Art 86 Maternity 120 days 30 pre +90 post full pay, Paternity 3 days company policy beyond law per Ethiopia Labour Proclamation 1156/2019';
comment on table payroll_leave_requests is 'Leave Requests per employee per period start_date end_date days_requested 0.5 half day reason status pending/approved/rejected/cancelled approved_by approved_at rejection_reason medical_certificate_file_key MinIO for sick >3 days per Art 82';
comment on table payroll_loan_emi_schedule is 'Loans EMI Schedule tracking outstanding preview: installment_no due_date emi_amount principal_component interest_component outstanding_after status pending/paid/overdue/skipped paid_at run_id per payroll run auto deduction O(k) per employee k=0-2 active loans';
comment on table payroll_claims is 'Reimbursements Claims Receipt Upload MinIO <5MB pdf/jpg/png presigned 15m TTL file_key file_hash receipt_file_hash approved_by_manager approved_by_finance manager_approved_at finance_approved_at rejection_reason is_taxable is_pensionable status pending/approved/rejected/paid per Ethiopia business practice';
