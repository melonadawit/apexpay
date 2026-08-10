-- Budgeting / FP&A: budgets per period+category, and budget-vs-actual variance.
create table if not exists budgets (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  period        text not null,           -- YYYY-MM
  category      text not null,           -- revenue | expense | by cost center
  budget_amount numeric(20,8) not null default 0,
  created_by    text,
  created_at    timestamptz not null default now(),
  unique (merchant_id, period, category)
);
create index if not exists budgets_merchant_idx on budgets (merchant_id, period);
