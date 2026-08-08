-- 0027: Workforce OS (HRIS) — complements payroll with the HR lifecycle that feeds it.
-- Contracts, org teams, attendance clocking, onboarding checklists, performance reviews.

create table hris_teams (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  name          text not null,
  department_id text references departments(id) on delete set null,
  manager_id    text references employees(id) on delete set null,
  description   text,
  created_at    timestamptz not null default now()
);
create index hris_teams_merchant_idx on hris_teams (merchant_id);

create table hris_contracts (
  id               text primary key,
  merchant_id      text not null references merchants(id) on delete cascade,
  employee_id      text not null references employees(id) on delete cascade,
  contract_type    text not null check (contract_type in ('permanent','fixed_term','probation','consultancy','internship','daily_labor')),
  start_date       date not null,
  end_date         date,
  salary_amount    numeric(20,8),
  salary_currency  char(3) not null default 'ETB',
  probation_months int not null default 3,
  notice_days      int not null default 30,
  status           text not null default 'active' check (status in ('draft','active','expired','terminated')),
  doc_key          text,
  doc_hash         text,
  signed_at        timestamptz,
  created_by       text references users(id),
  created_at       timestamptz not null default now(),
  updated_at       timestamptz not null default now()
);
create index hris_contracts_employee_idx on hris_contracts (employee_id);
create index hris_contracts_expiry_idx on hris_contracts (end_date) where end_date is not null;

create table hris_shifts (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  name        text not null,
  start_time  time not null,
  end_time    time not null,
  break_min   int not null default 0,
  created_at  timestamptz not null default now()
);
create index hris_shifts_merchant_idx on hris_shifts (merchant_id);

create table hris_attendance_clock (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  employee_id text not null references employees(id) on delete cascade,
  shift_id    text references hris_shifts(id) on delete set null,
  clock_date  date not null,
  punch_in    timestamptz,
  punch_out   timestamptz,
  hours       numeric(6,2) not null default 0,
  status      text not null default 'present' check (status in ('present','late','absent','on_leave','half_day')),
  source      text not null default 'manual' check (source in ('manual','clock','admin')),
  note        text,
  created_at  timestamptz not null default now(),
  unique (merchant_id, employee_id, clock_date)
);
create index hris_clock_merchant_date_idx on hris_attendance_clock (merchant_id, clock_date);

create table hris_onboarding_checklists (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  employee_id text not null references employees(id) on delete cascade,
  task        text not null,
  category    text not null default 'hr' check (category in ('hr','it','finance','security','compliance')),
  due_in_days int not null default 7,
  status      text not null default 'pending' check (status in ('pending','in_progress','done')),
  assigned_to text references users(id) on delete set null,
  completed_at timestamptz,
  created_at  timestamptz not null default now()
);
create index hris_onboarding_employee_idx on hris_onboarding_checklists (employee_id);

create table hris_performance_reviews (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  employee_id text not null references employees(id) on delete cascade,
  reviewer_id text references users(id) on delete set null,
  period      text not null, -- e.g. "2026-Q3" or "2026-Annual"
  rating      int check (rating between 1 and 5),
  goals       text,
  comments    text,
  status      text not null default 'draft' check (status in ('draft','submitted','completed')),
  reviewed_at timestamptz,
  created_at  timestamptz not null default now()
);
create index hris_perf_employee_idx on hris_performance_reviews (employee_id);
