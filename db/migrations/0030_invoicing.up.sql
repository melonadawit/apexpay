-- 0030: Invoicing & Receivables Automation.
-- Invoices with line items, VAT/TOT/withholding, hosted payment, dunning, AR aging.

create table invoices (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  invoice_number      text not null,
  customer_name       text not null,
  customer_email      text,
  customer_phone      text,
  issue_date          date not null,
  due_date            date not null,
  currency            char(3) not null default 'ETB',
  subtotal            numeric(20,8) not null default 0,
  tax_amount          numeric(20,8) not null default 0,  -- VAT 15% / TOT
  withholding_amount  numeric(20,8) not null default 0,  -- withholding 2% for services
  total_amount        numeric(20,8) not null default 0,
  amount_paid         numeric(20,8) not null default 0,
  status              text not null default 'draft' check (status in ('draft','sent','partially_paid','paid','overdue','cancelled')),
  hosted_token        text, -- hosted payment link token
  payment_link_id     text references payment_links(id) on delete set null,
  dunning_stage       int not null default 0,  -- 0 none, 1..3 reminders
  last_dunning_at     timestamptz,
  notes               text,
  created_by          text references users(id),
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now(),
  unique (merchant_id, invoice_number)
);
create index invoices_merchant_status_idx on invoices (merchant_id, status, due_date);
create index invoices_customer_idx on invoices (merchant_id, customer_email);

create table invoice_line_items (
  id          text primary key,
  invoice_id  text not null references invoices(id) on delete cascade,
  description text not null,
  quantity    numeric(12,4) not null default 1,
  unit_price  numeric(20,8) not null default 0,
  line_total  numeric(20,8) not null default 0,
  sort_order  int not null default 0
);
create index invoice_lines_invoice_idx on invoice_line_items (invoice_id);

create table dunning_logs (
  id          text primary key,
  invoice_id  text not null references invoices(id) on delete cascade,
  merchant_id text not null references merchants(id),
  stage       int not null,
  channel     text not null default 'email' check (channel in ('email','sms','push')),
  sent_at     timestamptz not null default now(),
  response    text
);
create index dunning_logs_invoice_idx on dunning_logs (invoice_id);
