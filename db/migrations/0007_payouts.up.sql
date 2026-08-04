-- 0007_payouts: disbursement single/bulk/payout links with escrow book + maker-checker

create table beneficiaries (
  id                text primary key,
  merchant_id       text not null references merchants(id) on delete cascade,
  name              text not null,
  account_no_masked text not null,
  account_no_hash   text not null,
  bank_code         text not null,
  bank_name         text not null,
  type              text not null check (type in ('individual','business')),
  verification_status text not null default 'pending' check (verification_status in ('pending','verified','failed')),
  verification_method text check (verification_method in ('name_match','micro_deposit','manual')),
  created_at        timestamptz not null default now()
);
create index beneficiaries_merchant_idx on beneficiaries (merchant_id);
create index beneficiaries_bank_idx on beneficiaries (bank_code);
create index beneficiaries_account_hash_idx on beneficiaries (account_no_hash);

create table payout_batches (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  book_id         text references ledger_books(id),
  batch_ref       text not null,
  amount          numeric(20,8) not null check (amount >0),
  currency        char(3) not null,
  status          text not null check (status in ('draft','pending_approval','approved','processing','completed','failed')),
  total_count     int not null default 0,
  success_count   int not null default 0,
  failed_count    int not null default 0,
  approved_by     text references users(id),
  created_at      timestamptz not null default now(),
  unique (merchant_id, batch_ref)
);
create index payout_batches_merchant_idx on payout_batches (merchant_id, created_at desc);
create index payout_batches_status_idx on payout_batches (status);

create table payouts (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  batch_id        text references payout_batches(id) on delete cascade,
  beneficiary_id  text references beneficiaries(id) on delete set null,
  payout_ref      text not null,
  amount          numeric(20,8) not null check (amount >0),
  currency        char(3) not null,
  status          text not null check (status in ('created','pending_approval','queued','processing','succeeded','failed','returned')),
  method          text not null check (method in ('bank','mobile_money','payout_link')),
  connector_id    text,
  connector_ref   text,
  failure_code    text,
  failure_message text,
  claimed_at      timestamptz, -- for payout links escrow
  expires_at      timestamptz, -- for payout links
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now(),
  unique (merchant_id, payout_ref)
);
create index payouts_merchant_created_idx on payouts (merchant_id, created_at desc);
create index payouts_batch_idx on payouts (batch_id);
create index payouts_beneficiary_idx on payouts (beneficiary_id);
create index payouts_status_idx on payouts (status);

comment on table payouts is 'Ledger Model M3: Dr liability:merchant_payable amount Cr asset:clearing:bank amount. For bulk, single journal per batch. Payout links use escrow book until claimed.';
