-- 0008_payroll: workforce money OS with ET tax brackets + pension 7%/11% + OT rates
-- Each payroll_run gets own ledger book per DATABASE

create table employees (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  employee_code       text not null,
  name                text not null,
  name_am             text,
  email               text,
  phone               text,
  tin                 text,
  fayda_fin_hash      text,
  pension_no          text,
  bank_account_hash   text,
  bank_account_masked text,
  bank_code           text,
  base_salary         numeric(20,8) not null check (base_salary >=0),
  employment_date     date not null,
  employment_type     text not null check (employment_type in ('permanent','contract','part_time','intern')) default 'permanent',
  cost_center         text,
  status              text not null check (status in ('active','inactive','terminated','on_leave')) default 'active',
  metadata            jsonb not null default '{}'::jsonb,
  created_at          timestamptz not null default now(),
  unique (merchant_id, employee_code)
);
create index employees_merchant_idx on employees (merchant_id, status);
create index employees_merchant_cost_center_idx on employees (merchant_id, cost_center);
create index employees_fayda_hash_idx on employees (fayda_fin_hash) where fayda_fin_hash is not null;

create table payroll_tax_brackets (
  id               text primary key,
  min_amount       numeric(20,8) not null,
  max_amount       numeric(20,8), -- null = infinity
  rate             numeric(5,4) not null check (rate >=0 and rate <=1),
  deduction        numeric(20,8) not null,
  effective_from   date not null default '2024-01-01',
  effective_to     date,
  created_at       timestamptz not null default now()
);
create index tax_brackets_effective_idx on payroll_tax_brackets (effective_from, effective_to);
-- Seed ET 2024 brackets placeholder
insert into payroll_tax_brackets (id, min_amount, max_amount, rate, deduction, effective_from) values
('brack_600', 0, 600, 0.0000, 0, '2024-01-01'),
('brack_1650', 601, 1650, 0.1000, 60, '2024-01-01'),
('brack_3200', 1651, 3200, 0.1500, 142.50, '2024-01-01'),
('brack_5250', 3201, 5250, 0.2000, 302.50, '2024-01-01'),
('brack_7800', 5251, 7800, 0.2500, 565, '2024-01-01'),
('brack_10900', 7801, 10900, 0.3000, 955, '2024-01-01'),
('brack_inf', 10901, null, 0.3500, 1500, '2024-01-01')
on conflict (id) do nothing;

create table payroll_runs (
  id               text primary key,
  merchant_id      text not null references merchants(id) on delete cascade,
  book_id          text references ledger_books(id),
  run_ref          text not null,
  period_month     int not null check (period_month between 1 and 12),
  period_year      int not null check (period_year >=2020),
  type             text not null check (type in ('regular','off_cycle','bonus','adjustment')) default 'regular',
  status           text not null check (status in ('draft','calculating','pending_approval','approved','processing','completed','failed','voided')) default 'draft',
  total_gross      numeric(20,8) not null default 0,
  total_deductions numeric(20,8) not null default 0,
  total_net        numeric(20,8) not null default 0,
  total_tax        numeric(20,8) not null default 0,
  total_pension    numeric(20,8) not null default 0,
  total_count      int not null default 0,
  approved_by      text references users(id),
  created_at       timestamptz not null default now(),
  unique (merchant_id, run_ref)
);
create index payroll_runs_merchant_idx on payroll_runs (merchant_id, period_year desc, period_month desc);
create index payroll_runs_status_idx on payroll_runs (status);

create table payroll_items (
  id                text primary key,
  run_id            text not null references payroll_runs(id) on delete cascade,
  employee_id       text not null references employees(id) on delete cascade,
  gross             numeric(20,8) not null check (gross >=0),
  ot_hours          numeric(10,2) not null default 0,
  ot_amount         numeric(20,8) not null default 0,
  commission        numeric(20,8) not null default 0,
  bonus             numeric(20,8) not null default 0,
  other_allowances  numeric(20,8) not null default 0,
  taxable_income    numeric(20,8) not null,
  income_tax        numeric(20,8) not null check (income_tax >=0),
  pension_employee  numeric(20,8) not null check (pension_employee >=0),
  pension_employer  numeric(20,8) not null check (pension_employer >=0),
  other_deductions  numeric(20,8) not null default 0,
  net_pay           numeric(20,8) not null check (net_pay >=0),
  status            text not null check (status in ('pending','calculated','approved','paid','failed')) default 'pending',
  failure_reason    text,
  created_at        timestamptz not null default now()
);
create index payroll_items_run_idx on payroll_items (run_id);
create index payroll_items_employee_idx on payroll_items (employee_id);
create index payroll_items_run_employee_unique on payroll_items (run_id, employee_id);

create table payroll_claims (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  employee_id     text not null references employees(id) on delete cascade,
  run_id          text references payroll_runs(id) on delete set null,
  claim_type      text not null check (claim_type in ('expense','medical','travel','other')),
  amount          numeric(20,8) not null check (amount >=0),
  description     text,
  receipt_file_key text, -- MinIO
  status          text not null check (status in ('pending','approved','rejected','paid')) default 'pending',
  created_at      timestamptz not null default now()
);
create index payroll_claims_employee_idx on payroll_claims (employee_id);
create index payroll_claims_run_idx on payroll_claims (run_id);

comment on table payroll_runs is 'Ledger Model M4: Dr expense:salary total_gross Cr liability:payroll_payable total_net Cr liability:et_income_tax_payable total_tax Cr liability:pension_payable total_pension. Per run book_id. Then disburse via payout batch: Dr payroll_payable Cr asset:clearing:bank.';
