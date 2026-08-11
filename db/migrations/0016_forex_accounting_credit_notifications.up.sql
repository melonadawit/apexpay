-- 0016_forex_accounting_credit_notifications: P0 remaining to surpass ApexPay — Forex FDI Transfers, Accounting Integrations Two-way Sync Tally Zoho QuickBooks, Instant Loans Digital Lending Collateral-free Credit Lines, Notifications Bulk Payouts Approval Refresh Button, Dedicated RM Priority Support
-- Ethiopia law compliance: Forex highly regulated by NBE, Accounting integrations for ERCA compliance, Credit scoring based on TPV payroll data

-- Forex Requests + Rates + Transactions — 2.5% Forex Markup Flat 1% Cashback per Corporate Cards
create table forex_rates (
  id                  text primary key,
  from_currency       char(3) not null, -- ETB
  to_currency         char(3) not null, -- USD, EUR, GBP, etc.
  rate                numeric(20,8) not null, -- e.g., 1 USD = 57.5 ETB
  buy_rate            numeric(20,8) not null,
  sell_rate           numeric(20,8) not null,
  source              text not null default 'nbe', -- nbe, commercial_bank, black_market? For compliance, use NBE official rate
  last_updated_at     timestamptz not null default now(),
  created_at          timestamptz not null default now(),
  unique (from_currency, to_currency)
);
create index forex_rates_pair_idx on forex_rates (from_currency, to_currency, last_updated_at desc);

create table forex_requests (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  from_currency       char(3) not null default 'ETB',
  to_currency         char(3) not null, -- USD, EUR, GBP, etc.
  from_amount         numeric(20,8) not null check (from_amount >0),
  to_amount           numeric(20,8) not null check (to_amount >0),
  forex_rate_id       text references forex_rates(id),
  rate_used           numeric(20,8) not null,
  forex_fee_percent   numeric(5,2) not null default 2.50, -- 2.5% forex markup per Corporate Cards
  forex_fee_amount    numeric(20,8) not null default 0,
  purpose             text not null, -- import_payment, service_payment, fdi, etc. per NBE
  purpose_description text,
  status              text not null check (status in ('draft','pending_nbe_approval','pending_bank_approval','approved','rejected','processing','completed','failed','cancelled')) default 'draft',
  nbe_approval_required boolean not null default true, -- Forex highly regulated by NBE per Ethiopia law
  nbe_approval_status text check (nbe_approval_status in ('pending','approved','rejected')) default 'pending',
  nbe_reference       text,
  bank_reference      text,
  created_by          text references users(id),
  approved_by         text references users(id),
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now()
);
create index forex_requests_merchant_idx on forex_requests (merchant_id, status, created_at desc);
create index forex_requests_nbe_approval_idx on forex_requests (nbe_approval_required, nbe_approval_status) where nbe_approval_required=true;

create table forex_transactions (
  id                  text primary key,
  request_id          text not null references forex_requests(id) on delete cascade,
  merchant_id         text not null references merchants(id) on delete cascade,
  from_amount         numeric(20,8) not null,
  to_amount           numeric(20,8) not null,
  from_currency       char(3) not null,
  to_currency         char(3) not null,
  rate_used           numeric(20,8) not null,
  status              text not null check (status in ('pending','processing','completed','failed','cancelled')) default 'pending',
  ledger_book_id      text references ledger_books(id),
  created_at          timestamptz not null default now()
);
create index forex_transactions_request_idx on forex_transactions (request_id);
create index forex_transactions_merchant_idx on forex_transactions (merchant_id, status);

-- Accounting Integrations Two-way Sync Tally Zoho QuickBooks CA Access Controls
create table accounting_integrations (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  provider            text not null check (provider in ('tally','zoho','quickbooks','xero','sage','other')), -- tally, zoho, quickbooks per ApexPay
  status              text not null check (status in ('connected','disconnected','error','pending')) default 'pending',
  credentials_encrypted text, -- encrypted via AES-GCM CONNECTOR_ENCRYPTION_KEY
  last_sync_at        timestamptz,
  last_sync_status    text check (last_sync_status in ('success','failed','partial')) default 'success',
  last_sync_error     text,
  created_by          text references users(id),
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now(),
  unique (merchant_id, provider)
);
create index accounting_integrations_merchant_idx on accounting_integrations (merchant_id, status);

create table accounting_sync_logs (
  id                  text primary key,
  integration_id      text not null references accounting_integrations(id) on delete cascade,
  merchant_id         text not null references merchants(id) on delete cascade,
  sync_type           text not null check (sync_type in ('payments','payouts','invoices','bills','journal_entries','contacts','full_sync')),
  status              text not null check (status in ('success','failed','partial','pending')) default 'pending',
  payload             jsonb not null default '{}'::jsonb,
  response            jsonb not null default '{}'::jsonb,
  records_synced      int not null default 0,
  error_message       text,
  created_at          timestamptz not null default now()
);
create index accounting_sync_logs_integration_idx on accounting_sync_logs (integration_id, created_at desc);
create index accounting_sync_logs_merchant_idx on accounting_sync_logs (merchant_id, sync_type, created_at desc);

-- Instant Loans Digital Lending Collateral-free Credit Lines — Capital Line of Credit
create table credit_lines (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  credit_limit        numeric(20,8) not null check (credit_limit >0), -- e.g., up to 2Cr ETB equivalent (20L-2Cr INR in India)
  available_credit    numeric(20,8) not null,
  utilized_credit     numeric(20,8) not null default 0,
  interest_rate       numeric(5,2) not null default 18.00, -- 18% per annum? Configurable
  status              text not null check (status in ('draft','pending_approval','approved','active','suspended','closed','rejected')) default 'draft',
  credit_score        int check (credit_score between 300 and 900), -- credit scoring based on TPV payroll data etc.
  approved_by         text references users(id),
  approved_at         timestamptz,
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now()
);
create index credit_lines_merchant_idx on credit_lines (merchant_id, status);

create table loan_disbursements (
  id                  text primary key,
  credit_line_id      text not null references credit_lines(id) on delete cascade,
  merchant_id         text not null references merchants(id) on delete cascade,
  amount              numeric(20,8) not null check (amount >0),
  currency            char(3) not null default 'ETB',
  purpose             text not null, -- working_capital, inventory, payroll, etc.
  status              text not null check (status in ('pending','approved','disbursed','repaid','defaulted','cancelled')) default 'pending',
  disbursed_at        timestamptz,
  due_date            date,
  repaid_amount       numeric(20,8) not null default 0,
  outstanding_amount  numeric(20,8) not null,
  ledger_book_id      text references ledger_books(id),
  created_by          text references users(id),
  created_at          timestamptz not null default now()
);
create index loan_disbursements_credit_line_idx on loan_disbursements (credit_line_id, status);
create index loan_disbursements_merchant_idx on loan_disbursements (merchant_id, status);

-- Notifications Bulk Payouts Approval Refresh Button Latest Balance Instantly + Pending Payouts Notification
create table notifications (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  user_id             text references users(id) on delete cascade, -- null for all users in merchant
  type                text not null check (type in ('bulk_payouts_approval','pending_payout','payout_failed','payroll_run_pending_approval','payroll_run_completed','tax_payment_due','compliance_alert','bank_file_generated','pension_csv_generated','erca_csv_generated','loan_emi_due','leave_request_pending','claim_pending','escrow_held','escrow_released','current_account_opened','corporate_card_transaction','forex_rate_alert','accounting_sync_failed','other')), -- bulk_payouts_approval per ApexPay
  title               text not null,
  message             text not null,
  data                jsonb not null default '{}'::jsonb, -- {payout_batch_id, payroll_run_id, amount, etc.}
  is_read             boolean not null default false,
  read_at             timestamptz,
  action_url          text, -- e.g., /payout_batches/{id} or /payroll/{id}
  created_at          timestamptz not null default now()
);
create index notifications_merchant_user_idx on notifications (merchant_id, user_id, is_read, created_at desc);
create index notifications_type_idx on notifications (merchant_id, type, is_read);
create index notifications_unread_idx on notifications (merchant_id, user_id) where is_read=false;

-- Dedicated Relationship Manager + Priority Support + SLA
create table relationship_managers (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  rm_user_id          text not null references users(id) on delete cascade, -- relationship manager user
  assigned_at         timestamptz not null default now(),
  assigned_by         text references users(id),
  status              text not null check (status in ('active','inactive','reassigned')) default 'active',
  created_at          timestamptz not null default now(),
  unique (merchant_id) -- one RM per merchant for simplicity
);

create table support_tickets (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  user_id             text references users(id) on delete set null,
  subject             text not null,
  description         text not null,
  priority            text not null check (priority in ('low','medium','high','urgent')) default 'medium',
  status              text not null check (status in ('open','in_progress','pending_customer','resolved','closed')) default 'open',
  assigned_to         text references users(id) on delete set null, -- RM or support agent
  sla_due_at          timestamptz,
  resolved_at         timestamptz,
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now()
);
create index support_tickets_merchant_idx on support_tickets (merchant_id, status, priority, created_at desc);
create index support_tickets_assigned_idx on support_tickets (assigned_to, status);

comment on table forex_rates is 'Forex Rates cached 60s via Redis per Ethiopia business practice highly regulated by NBE per Ethiopia law: from_currency ETB to_currency USD EUR GBP etc rate buy_rate sell_rate source nbe commercial_bank black_market? For compliance use NBE official rate last_updated_at';
comment on table forex_requests is 'Forex Requests + Rates + Transactions — 2.5% Forex Markup Flat 1% Cashback per Corporate Cards: from_currency ETB to_currency USD EUR GBP etc from_amount to_amount forex_rate_id rate_used forex_fee_percent 2.50 forex_fee_amount purpose import_payment service_payment fdi etc per NBE purpose_description status draft/pending_nbe_approval/pending_bank_approval/approved/rejected/processing/completed/failed/cancelled nbe_approval_required true Forex highly regulated by NBE per Ethiopia law nbe_approval_status pending/approved/rejected nbe_reference bank_reference created_by approved_by';
comment on table accounting_integrations is 'Accounting Integrations Two-way Sync Tally Zoho QuickBooks CA Access Controls per ApexPay: provider tally/zoho/quickbooks/xero/sage/other status connected/disconnected/error/pending credentials_encrypted via AES-GCM CONNECTOR_ENCRYPTION_KEY last_sync_at last_sync_status success/failed/partial last_sync_error created_by';
comment on table credit_lines is 'Instant Loans Digital Lending Collateral-free Credit Lines — Capital Line of Credit: credit_limit up to 2Cr ETB equivalent 20L-2Cr INR in available for Ethiopia_credit utilized_credit interest_rate 18% per annum configurable status draft/pending_approval/approved/active/suspended/closed/rejected credit_score 300-900 credit scoring based on TPV payroll data etc approved_by approved_at';
comment on table notifications is 'Notifications Bulk Payouts Approval Refresh Button Latest Balance Instantly + Pending Payouts Notification per ApexPay: type bulk_payouts_approval pending_payout payout_failed payroll_run_pending_approval payroll_run_completed tax_payment_due compliance_alert bank_file_generated pension_csv_generated erca_csv_generated loan_emi_due leave_request_pending claim_pending escrow_held escrow_released current_account_opened corporate_card_transaction forex_rate_alert accounting_sync_failed other title message data jsonb payout_batch_id payroll_run_id amount etc is_read read_at action_url /payout_batches/{id} or /payroll/{id}';
