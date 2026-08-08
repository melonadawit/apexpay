-- Procurement & Accounts Payable: vendors, purchase orders, receipts, AP invoices.
-- Note: table names are ap_-prefixed for POs/receipts to avoid colliding with the
-- escrow purchase_orders table already created by migration 0015.
create table if not exists vendors (
  id                 text primary key,
  merchant_id        text not null references merchants(id) on delete cascade,
  name               text not null,
  email              text,
  phone              text,
  tin                text,
  payment_terms_days int  not null default 30,
  status             text not null default 'active' check (status in ('active','inactive')),
  created_at         timestamptz not null default now()
);
create index if not exists vendors_merchant_idx on vendors (merchant_id);

create table if not exists ap_purchase_orders (
  id               text primary key,
  merchant_id      text not null references merchants(id) on delete cascade,
  vendor_id        text not null references vendors(id),
  po_number        text not null,
  order_date       date not null,
  expected_delivery date,
  status           text not null default 'draft' check (status in ('draft','approved','received','closed','cancelled')),
  subtotal         numeric(20,8) not null default 0,
  tax_amount       numeric(20,8) not null default 0,
  total_amount     numeric(20,8) not null default 0,
  created_by       text,
  created_at       timestamptz not null default now(),
  unique (merchant_id, po_number)
);
create index if not exists ap_purchase_orders_merchant_idx on ap_purchase_orders (merchant_id, created_at desc);

create table if not exists ap_purchase_order_items (
  id              text primary key,
  purchase_order_id text not null references ap_purchase_orders(id) on delete cascade,
  item_name       text not null,
  quantity        numeric(20,8) not null,
  unit_price      numeric(20,8) not null,
  line_total      numeric(20,8) not null,
  received_qty    numeric(20,8) not null default 0
);

create table if not exists ap_receipts (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  purchase_order_id text references ap_purchase_orders(id),
  vendor_id       text not null references vendors(id),
  receipt_number  text not null,
  received_at     timestamptz not null default now(),
  note            text,
  created_by      text,
  created_at      timestamptz not null default now()
);
create index if not exists ap_receipts_merchant_idx on ap_receipts (merchant_id);

create table if not exists ap_receipt_lines (
  id          text primary key,
  receipt_id  text not null references ap_receipts(id) on delete cascade,
  item_name   text not null,
  quantity    numeric(20,8) not null,
  unit_price  numeric(20,8) not null
);

create table if not exists ap_invoices (
  id                text primary key,
  merchant_id       text not null references merchants(id) on delete cascade,
  vendor_id         text not null references vendors(id),
  purchase_order_id text references ap_purchase_orders(id),
  invoice_number    text not null,
  invoice_date      date not null,
  due_date          date not null,
  subtotal          numeric(20,8) not null default 0,
  tax_amount        numeric(20,8) not null default 0,
  total_amount      numeric(20,8) not null default 0,
  amount_paid       numeric(20,8) not null default 0,
  status            text not null default 'pending' check (status in ('pending','matched','partially_paid','paid','overdue')),
  match_status      text not null default 'unmatched' check (match_status in ('unmatched','matched','mismatch')),
  created_by        text,
  created_at        timestamptz not null default now(),
  unique (merchant_id, invoice_number)
);
create index if not exists ap_invoices_merchant_idx on ap_invoices (merchant_id, due_date);
