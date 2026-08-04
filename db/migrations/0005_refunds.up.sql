-- 0005_refunds: FULL refund support with fee reversal policy + ledger M2
-- Refund posting: Dr merchant_payable (R-FR) + Dr fee_due FR Cr clearing

create table refunds (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  payment_id      text not null references payments(id) on delete cascade,
  refund_ref      text not null,
  amount          numeric(20,8) not null check (amount > 0),
  currency        char(3) not null,
  status          text not null check (status in ('created','processing','succeeded','failed')),
  reason          text,
  fee_reversal    numeric(20,8) not null default 0 check (fee_reversal >=0),
  fee_policy      text not null default 'non_refundable' check (fee_policy in ('non_refundable','pro_rata','full')),
  connector_id    text not null,
  connector_ref   text,
  failure_code    text,
  failure_message text,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now(),
  unique (merchant_id, refund_ref)
);
create index refunds_payment_idx on refunds (payment_id);
create index refunds_merchant_created_idx on refunds (merchant_id, created_at desc);
create index refunds_status_idx on refunds (status);

-- Add payment refunded amount tracking via view or computed, but we keep sum query
comment on table refunds is 'Ledger Model M2: Dr liability:merchant_payable amount-feeReversal + Dr liability:platform_fee_due feeReversal Cr asset:clearing:connector amount. Filter zero entries if feeReversal zero.';
