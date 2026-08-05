-- 0015_current_accounts_escrow_corporate_cards: P0 Business Banking Core — Current Accounts Real Partner Bank (CBE/Awash/Dashen) + Cheque Book + Debit Card + Lite Interim Account + Multiple Accounts Balance Snapshot + Escrow Automated Marketplace + Corporate Cards Virtual + Physical + Payout Links QR Scan & Pay
-- Senior Engineer design: clean arch, decimal precise, ULID, optimal data structures, quality indexes per Ethiopia business practice + NBE ONPS/02/2020 09/2023 10/2025 + RazorpayX parity

-- Current Accounts Real — Issued by partner banks CBE/Awash/Dashen + cheque book + debit card + unlimited deposits/withdrawals + lite interim + multiple accounts
create table current_accounts (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  account_number      text not null, -- virtual account number e.g., ETB-CBE-1234567890
  account_name        text not null,
  account_type        text not null check (account_type in ('current','saving','virtual','escrow','reserve')) default 'current',
  currency            char(3) not null default 'ETB',
  bank_code           text not null, -- CBE, AWASH, DASHEN, ABYSSINIA, etc.
  partner_bank_name   text not null, -- Commercial Bank of Ethiopia, Awash Bank, Dashen Bank, etc. — Ethiopia equivalent of ICICI/Axis/RBL/YES in RazorpayX India
  status              text not null check (status in ('draft','pending_kyc','pending_approval','active','suspended','closed','frozen')) default 'draft',
  balance             numeric(20,8) not null default 0,
  available_balance   numeric(20,8) not null default 0,
  overdraft_limit     numeric(20,8) not null default 0,
  is_primary          boolean not null default false, -- primary settlement account
  is_lite             boolean not null default false, -- lite interim account until current account active per RazorpayX Lite concept
  is_virtual          boolean not null default false, -- virtual account for collections smart collect
  cheque_book_issued  boolean not null default false,
  debit_card_issued   boolean not null default false,
  debit_card_type     text check (debit_card_type in ('virtual','physical','both')),
  created_by          text references users(id),
  approved_by         text references users(id),
  approved_at         timestamptz,
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now(),
  unique (merchant_id, account_number),
  unique (account_number)
);
create index current_accounts_merchant_idx on current_accounts (merchant_id, status);
create index current_accounts_merchant_primary_idx on current_accounts (merchant_id) where is_primary=true;
create index current_accounts_bank_code_idx on current_accounts (bank_code, status);
create index current_accounts_balance_idx on current_accounts (balance) where status='active';

-- Current Account Opening Requests — Online <24h paperless per RazorpayX, with NBE approval flow
create table current_account_opening_requests (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  kyc_profile_id      text references merchant_kyc_profiles(id) on delete set null,
  partner_bank_code   text not null, -- CBE, AWASH, DASHEN etc.
  account_type        text not null check (account_type in ('current','saving','virtual')) default 'current',
  requested_account_name text not null,
  status              text not null check (status in ('draft','submitted','in_review','kyc_pending','approved','rejected','needs_more_info')) default 'draft',
  risk_score          int not null default 0,
  risk_tier           text check (risk_tier in ('low','medium','high')) default 'medium',
  submitted_at        timestamptz,
  reviewed_at         timestamptz,
  reviewer_id         text references users(id),
  rejection_reason    text,
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now()
);
create index current_account_opening_requests_merchant_idx on current_account_opening_requests (merchant_id, status);
create index current_account_opening_requests_status_idx on current_account_opening_requests (status, created_at desc);

-- Cheque Books — Issuance Tracking
create table cheque_books (
  id                  text primary key,
  current_account_id  text not null references current_accounts(id) on delete cascade,
  merchant_id         text not null references merchants(id) on delete cascade,
  cheque_book_number  text not null,
  start_cheque_number int not null,
  end_cheque_number   int not null,
  total_cheques       int not null check (total_cheques >0),
  used_cheques        int not null default 0,
  status              text not null check (status in ('ordered','issued','active','used_up','blocked','cancelled')) default 'ordered',
  issued_at           timestamptz,
  issued_by           text references users(id),
  created_at          timestamptz not null default now()
);
create index cheque_books_account_idx on cheque_books (current_account_id, status);
create unique index cheque_books_number_uidx on cheque_books (cheque_book_number);

-- Debit Cards — Virtual + Physical Issuance Tracking
create table debit_cards (
  id                  text primary key,
  current_account_id  text not null references current_accounts(id) on delete cascade,
  merchant_id         text not null references merchants(id) on delete cascade,
  card_number_masked  text not null, -- ****1234
  card_number_hash    text not null, -- sha256 hash for lookup
  card_type           text not null check (card_type in ('virtual','physical','both')) default 'virtual',
  card_network        text not null check (card_network in ('visa','mastercard','verve','ethswitch')) default 'visa',
  status              text not null check (status in ('ordered','active','blocked','expired','cancelled')) default 'ordered',
  daily_limit         numeric(20,8) not null default 50000,
  monthly_limit       numeric(20,8) not null default 500000,
  cardholder_name     text not null,
  expiry_month        int check (expiry_month between 1 and 12),
  expiry_year         int,
  cvv_hash            text,
  is_contactless      boolean not null default true,
  created_by          text references users(id),
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now()
);
create index debit_cards_account_idx on debit_cards (current_account_id, status);
create index debit_cards_merchant_idx on debit_cards (merchant_id, status);

-- Escrow Accounts Automated — For Marketplaces P2P Hold & Release Funds Under Defined Conditions
create table escrow_accounts (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade, -- marketplace operator
  agreement_id        text not null, -- references escrow_agreements
  account_number      text not null unique,
  account_name        text not null,
  amount              numeric(20,8) not null check (amount >0),
  currency            char(3) not null default 'ETB',
  status              text not null check (status in ('draft','held','released','returned','disputed','expired')) default 'draft',
  held_at             timestamptz,
  release_at          timestamptz,
  return_at           timestamptz,
  expires_at          timestamptz,
  buyer_merchant_id   text references merchants(id),
  seller_merchant_id  text references merchants(id),
  order_id            text,
  order_amount        numeric(20,8),
  platform_fee        numeric(20,8) not null default 0, -- e.g., 10% 100 ETB
  seller_amount       numeric(20,8) not null default 0, -- e.g., 90% 900 ETB
  withholding_tax     numeric(20,8) not null default 0, -- e.g., 2% 20 ETB
  ledger_book_id      text references ledger_books(id),
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now()
);
create index escrow_accounts_merchant_idx on escrow_accounts (merchant_id, status);
create index escrow_accounts_buyer_seller_idx on escrow_accounts (buyer_merchant_id, seller_merchant_id);
create index escrow_accounts_status_expires_idx on escrow_accounts (status, expires_at) where status='held';

create table escrow_agreements (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  agreement_number    text not null unique,
  title               text not null,
  description         text,
  buyer_merchant_id   text references merchants(id),
  seller_merchant_id  text references merchants(id),
  amount              numeric(20,8) not null,
  currency            char(3) not null default 'ETB',
  platform_fee_percent numeric(5,2) not null default 10.00, -- 10% platform fee
  withholding_tax_percent numeric(5,2) not null default 2.00, -- 2% withholding tax for services per Ethiopia Income Tax Proclamation
  conditions          jsonb not null default '[]'::jsonb, -- [{type: delivery_confirmed, days: 7}, {type: inspection_period, days: 3}]
  auto_release        boolean not null default true,
  auto_release_after_days int not null default 7,
  status              text not null check (status in ('draft','active','completed','disputed','cancelled')) default 'draft',
  created_by          text references users(id),
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now()
);
create index escrow_agreements_merchant_idx on escrow_agreements (merchant_id, status);

-- Corporate Cards — Collateral-free Credit Cards Virtual + Physical Dynamic Limit up to 2Cr ETB equivalent, Spending Controls, Real-time Expense Tracking, Multi-level Approvals
create table corporate_cards (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  current_account_id  text references current_accounts(id) on delete set null,
  card_number_masked  text not null, -- ****1234
  card_number_hash    text not null, -- sha256 hash for lookup
  card_type           text not null check (card_type in ('virtual','physical','both')) default 'virtual',
  card_network        text not null check (card_network in ('visa','mastercard','verve','ethswitch')) default 'visa',
  cardholder_name     text not null,
  cardholder_email    text,
  status              text not null check (status in ('ordered','active','blocked','expired','cancelled','suspended')) default 'ordered',
  credit_limit        numeric(20,8) not null default 2000000, -- up to 2Cr ETB equivalent (20L-2Cr INR in RazorpayX India)
  available_credit    numeric(20,8) not null default 2000000,
  daily_limit         numeric(20,8) not null default 50000,
  monthly_limit       numeric(20,8) not null default 500000,
  category_restrictions jsonb not null default '[]'::jsonb, -- ["SaaS", "Cloud", "Marketing"] etc.
  spending_controls   jsonb not null default '{}'::jsonb, -- {daily_limit: 50000, monthly_limit: 500000, allowed_categories: ["SaaS", "Cloud"], blocked_merchants: []}
  cashback_percent    numeric(5,2) not null default 1.00, -- flat 1% cashback per RazorpayX
  forex_markup_percent numeric(5,2) not null default 2.50, -- 2.5% forex markup
  interest_free_days  int not null default 45, -- up to 45-50 day interest-free period
  is_addon            boolean not null default false, -- addon card unlimited
  parent_card_id      text references corporate_cards(id),
  created_by          text references users(id),
  approved_by         text references users(id),
  approved_at         timestamptz,
  expiry_month        int check (expiry_month between 1 and 12),
  expiry_year         int,
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now()
);
create index corporate_cards_merchant_idx on corporate_cards (merchant_id, status);
create index corporate_cards_current_account_idx on corporate_cards (current_account_id);
create index corporate_cards_cardholder_idx on corporate_cards (cardholder_email);

create table corporate_card_transactions (
  id                  text primary key,
  card_id             text not null references corporate_cards(id) on delete cascade,
  merchant_id         text not null references merchants(id) on delete cascade,
  amount              numeric(20,8) not null,
  currency            char(3) not null default 'ETB',
  merchant_name       text not null, -- e.g., AWS, Google Cloud, Facebook Ads
  merchant_category   text not null, -- SaaS, Cloud, Marketing, etc.
  status              text not null check (status in ('pending','approved','declined','reversed')) default 'pending',
  decline_reason      text,
  cashback_amount     numeric(20,8) not null default 0, -- flat 1% cashback
  forex_fee           numeric(20,8) not null default 0, -- 2.5% forex markup if international
  created_at          timestamptz not null default now()
);
create index corporate_card_transactions_card_idx on corporate_card_transactions (card_id, created_at desc);
create index corporate_card_transactions_merchant_idx on corporate_card_transactions (merchant_id, created_at desc);

-- Payout Links Enhanced — QR + Scan & Pay + SMS/Email/WhatsApp + Recipient Enters Account Details + OTP Claim + Escrow Book
create table payout_links_enhanced (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  amount              numeric(20,8) not null check (amount >0),
  currency            char(3) not null default 'ETB',
  public_token        text not null unique, -- for QR + public link
  qr_code_data        text, -- QR code data for EthSwitch interoperable QR
  recipient_name      text,
  recipient_phone     text,
  recipient_email     text,
  purpose             text, -- refund, cashback, reward, vendor payment
  status              text not null check (status in ('active','claimed','expired','cancelled','failed')) default 'active',
  expires_at          timestamptz not null,
  claimed_at          timestamptz,
  beneficiary_id      text references beneficiaries(id) on delete set null, -- once claimed, beneficiary created
  escrow_book_id      text references ledger_books(id), -- escrow book until claimed
  ledger_book_id      text references ledger_books(id),
  created_by          text references users(id),
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now()
);
create index payout_links_enhanced_merchant_idx on payout_links_enhanced (merchant_id, status);
create index payout_links_enhanced_token_uidx on payout_links_enhanced (public_token);
create index payout_links_enhanced_expires_idx on payout_links_enhanced (status, expires_at) where status='active';

-- Vendor Invoices + Purchase Orders + Petty Cash — End-to-End Accounts Payable Automation OCR
create table vendor_invoices (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  vendor_id           text, -- references beneficiaries or vendor_contacts
  invoice_number      text not null,
  invoice_date        date not null,
  due_date            date,
  amount              numeric(20,8) not null check (amount >0),
  currency            char(3) not null default 'ETB',
  tax_amount          numeric(20,8) not null default 0, -- VAT 15% TOT 2%/10%
  withholding_tax_amount numeric(20,8) not null default 0, -- 2% for services per Ethiopia Income Tax Proclamation
  total_amount        numeric(20,8) not null, -- amount + tax - withholding
  status              text not null check (status in ('draft','pending_approval','approved','paid','rejected','cancelled')) default 'draft',
  ocr_raw             jsonb not null default '{}'::jsonb, -- {extracted_text, confidence, vendor_name, tin, invoice_number, amount, tax, withholding, etc.}
  file_key            text, -- MinIO file key for invoice PDF/image
  file_hash           text,
  created_by          text references users(id),
  approved_by         text references users(id),
  approved_at         timestamptz,
  paid_at             timestamptz,
  payout_id           text references payouts(id) on delete set null,
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now(),
  unique (merchant_id, invoice_number)
);
create index vendor_invoices_merchant_idx on vendor_invoices (merchant_id, status, due_date);

create table purchase_orders (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  vendor_id           text,
  po_number           text not null,
  amount              numeric(20,8) not null,
  currency            char(3) not null default 'ETB',
  status              text not null check (status in ('draft','sent','approved','received','cancelled','closed')) default 'draft',
  created_by          text references users(id),
  approved_by         text references users(id),
  created_at          timestamptz not null default now(),
  unique (merchant_id, po_number)
);

create table petty_cash_budgets (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  budget_name         text not null,
  amount              numeric(20,8) not null check (amount >0),
  assigned_to         text references users(id) on delete set null,
  status              text not null check (status in ('active','closed','exhausted')) default 'active',
  spent_amount        numeric(20,8) not null default 0,
  remaining_amount    numeric(20,8) not null default 0,
  created_by          text references users(id),
  created_at          timestamptz not null default now()
);
create index petty_cash_budgets_merchant_idx on petty_cash_budgets (merchant_id, status);

create table petty_cash_expenses (
  id                  text primary key,
  budget_id           text not null references petty_cash_budgets(id) on delete cascade,
  merchant_id         text not null references merchants(id) on delete cascade,
  amount              numeric(20,8) not null check (amount >0),
  description         text not null,
  receipt_file_key    text,
  receipt_file_hash   text,
  status              text not null check (status in ('pending','approved','rejected','paid')) default 'pending',
  approved_by         text references users(id),
  created_by          text references users(id),
  created_at          timestamptz not null default now()
);
create index petty_cash_expenses_budget_idx on petty_cash_expenses (budget_id, status);

-- Tax Payments Automated Pre-filled Forms Challans Inbox Accountant Collaboration VAT 15% TOT Withholding 2%
create table tax_payments (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  tax_type            text not null check (tax_type in ('vat','tot','withholding','paye','pension','corporate_tax','excise','other')), -- VAT 15% TOT 2%/10% Withholding 2% PAYE (income tax) Pension 7%/11%
  amount              numeric(20,8) not null check (amount >0),
  currency            char(3) not null default 'ETB',
  period_month        int check (period_month between 1 and 12),
  period_year         int,
  due_date            date,
  status              text not null check (status in ('draft','pending_approval','pending','paid','failed','cancelled')) default 'draft',
  challan_file_key    text, -- MinIO file key for challan PDF
  challan_file_hash   text,
  payment_reference   text, -- bank payment reference / UTR
  paid_at             timestamptz,
  created_by          text references users(id),
  approved_by         text references users(id),
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now()
);
create index tax_payments_merchant_idx on tax_payments (merchant_id, tax_type, status, due_date);

create table tax_accountants (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  user_id             text not null references users(id) on delete cascade,
  role                text not null check (role in ('compliance','accountant','auditor','viewer')) default 'accountant',
  permissions         jsonb not null default '[]'::jsonb,
  created_at          timestamptz not null default now(),
  unique (merchant_id, user_id)
);

-- Bank Account Verification Penny Testing 1 ETB
create table bank_account_verifications (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  bank_code           text not null,
  account_number_masked text not null,
  account_number_hash text not null,
  account_name        text,
  verification_method text not null check (verification_method in ('penny_test','micro_deposit','bank_letter','manual')) default 'penny_test',
  amount              numeric(20,8) not null default 1.00, -- 1 ETB penny test
  connector_id        text not null, -- bank_ips, telebirr, etc.
  status              text not null check (status in ('pending','processing','verified','failed','expired')) default 'pending',
  verification_response jsonb not null default '{}'::jsonb, -- {beneficiary_name, account_name_match_score, bank_details, etc.}
  beneficiary_name_returned text,
  match_score         numeric(5,2), -- fuzzy Levenshtein <3 match score
  verified_at         timestamptz,
  expires_at          timestamptz,
  created_at          timestamptz not null default now()
);
create index bank_account_verifications_merchant_idx on bank_account_verifications (merchant_id, status);
create index bank_account_verifications_hash_idx on bank_account_verifications (account_number_hash);

-- Collections / Smart Collect / Virtual Accounts
create table virtual_accounts (
  id                  text primary key,
  merchant_id         text not null references merchants(id) on delete cascade,
  virtual_account_number text not null unique, -- e.g., VA-CBE-1234567890
  customer_id         text, -- references customers or beneficiaries
  purpose             text not null, -- e.g., customer collections, vendor payments, payroll
  status              text not null check (status in ('active','inactive','closed')) default 'active',
  bank_code           text not null,
  created_at          timestamptz not null default now()
);
create index virtual_accounts_merchant_idx on virtual_accounts (merchant_id, status);

create table virtual_account_transactions (
  id                  text primary key,
  virtual_account_id  text not null references virtual_accounts(id) on delete cascade,
  merchant_id         text not null references merchants(id) on delete cascade,
  amount              numeric(20,8) not null,
  currency            char(3) not null default 'ETB',
  utr                 text, -- Unique Transaction Reference
  sender_name         text,
  sender_account      text,
  status              text not null check (status in ('pending','matched','unmatched','reconciled','failed')) default 'pending',
  matched_invoice_id  text references vendor_invoices(id) on delete set null,
  matched_at          timestamptz,
  created_at          timestamptz not null default now()
);
create index virtual_account_transactions_virtual_account_idx on virtual_account_transactions (virtual_account_id, status, created_at desc);
create index virtual_account_transactions_merchant_idx on virtual_account_transactions (merchant_id, status);

comment on table current_accounts is 'Current Accounts Real Partner Bank (CBE/Awash/Dashen) + Cheque Book + Debit Card + Unlimited Deposits/Withdrawals + Lite Interim Account + Multiple Accounts Balance Snapshot per RazorpayX Current Account free online <24h paperless no min balance issued by partner banks RBL/ICICI/Axis/YES equivalent for Ethiopia CBE/Awash/Dashen';
comment on table cheque_books is 'Cheque Books Issuance Tracking per current account: start_cheque_number end_cheque_number total_cheques used_cheques status ordered/issued/active/used_up/blocked/cancelled issued_at issued_by';
comment on table debit_cards is 'Debit Cards Virtual + Physical Issuance Tracking: card_number_masked ****1234 card_number_hash sha256 hash last4 card_type virtual/physical/both card_network visa/mastercard/verve/ethswitch status ordered/active/blocked/expired/cancelled daily_limit monthly_limit cardholder_name expiry_month year cvv_hash is_contactless';
comment on table escrow_accounts is 'Escrow Accounts Automated Marketplace P2P Hold & Release Funds Under Defined Conditions reduces legal overhead: buyer_merchant_id seller_merchant_id order_id order_amount platform_fee 10% seller_amount 90% withholding_tax 2% ledger_book_id per agreement book_type escrow auto_release_after_days 7';
comment on table escrow_agreements is 'Escrow Agreements parties conditions auto_release: buyer merchant seller merchant amount currency platform_fee_percent 10% withholding_tax_percent 2% conditions JSON [{type: delivery_confirmed, days: 7}, {type: inspection_period, days: 3}] auto_release true auto_release_after_days 7 status draft/active/completed/disputed/cancelled';
comment on table corporate_cards is 'Corporate Cards Collateral-free Credit Cards Virtual + Physical Dynamic Limit up to 2Cr ETB equivalent up to 45-50 day interest-free period 2.5% forex markup flat 1% cashback 30% off SaaS custom spending controls daily_limit monthly_limit category_restrictions spending_controls jsonb cashback_percent 1% forex_markup_percent 2.5% interest_free_days 45 is_addon parent_card_id unlimited add-on cards';
comment on table corporate_card_transactions is 'Corporate Card Transactions real-time expense tracking monitor all transactions instantly for improved visibility and timely reconciliation: amount merchant_name merchant_category SaaS Cloud Marketing etc status pending/approved/declined decline_reason cashback_amount 1% flat forex_fee 2.5%';
comment on table payout_links_enhanced is 'Payout Links Enhanced QR + Scan & Pay + SMS/Email/WhatsApp No Bank Details Needed Recipient Enters Account Details Bank Account or UPI ID Instant Payout QR Based Payouts View Attached Vendor Invoices: amount currency public_token unique QR preview share Telegram/WhatsApp purpose refund/cashback/reward/vendor payment status active/claimed/expired/cancelled recipient_name phone email expires_at claimed_at beneficiary_id once claimed escrow book until claimed ledger per agreement book_type escrow';
comment on table vendor_invoices is 'Vendor Invoices OCR-enabled Invoice Capture Multi-layer Approval Workflows Automated TDS Calculation and Filing to NSDL Integrated Payouts: invoice_number invoice_date due_date amount currency tax_amount VAT 15% TOT 2%/10% withholding_tax_amount 2% for services per Ethiopia Income Tax Proclamation total_amount status draft/pending_approval/approved/paid/rejected/cancelled ocr_raw JSON extracted_text confidence vendor_name tin invoice_number amount tax withholding file_key MinIO file_key file_hash';
comment on table tax_payments is 'Tax Payments Automated Pre-filled Forms Challans Inbox Accountant Collaboration VAT 15% TOT 2%/10% Withholding 2% PAYE Pension 7%/11%: tax_type vat/tot/withholding/paye/pension/corporate_tax/excise/other amount currency period_month year due_date status draft/pending_approval/pending/paid/failed/cancelled challan_file_key file_hash payment_reference paid_at created_by approved_by';
comment on table bank_account_verifications is 'Bank Account Verification Penny Testing Fund Account Validation 1 ETB Deposit Single Rupee Returns Validated Bank Details + Beneficiary Name per RazorpayX Penny Testing: bank_code account_number_masked hash account_name verification_method penny_test/micro_deposit/bank_letter/manual amount 1 ETB connector_id bank_ips telebirr status pending/processing/verified/failed/expired verification_response JSON beneficiary_name_returned match_score fuzzy Levenshtein <3 verified_at expires_at';
comment on table virtual_accounts is 'Collections Smart Collect Virtual Accounts Automatically Reconcile Incoming NEFT RTGS IMPS UPI Payments Using Virtual Accounts & UPI-IDs: virtual_account_number customer_id purpose status active/inactive/closed bank_code';
