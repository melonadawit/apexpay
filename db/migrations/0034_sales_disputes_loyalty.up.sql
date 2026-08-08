-- 0034: Inventory & Sales (software POS), Disputes & Chargebacks, Loyalty & Cashback.

-- ---- Inventory & Sales ----
create table products (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  name          text not null,
  description   text,
  sku           text,
  price         numeric(20,8) not null check (price >= 0),
  cost_price    numeric(20,8) not null default 0,
  currency      char(3) not null default 'ETB',
  vat_category  text not null default 'standard' check (vat_category in ('standard','zero','exempt')),
  stock_qty     numeric(12,4) not null default 0,
  low_stock_threshold numeric(12,4) not null default 5,
  status        text not null default 'active' check (status in ('active','inactive')),
  created_at    timestamptz not null default now(),
  updated_at    timestamptz not null default now()
);
create index products_merchant_idx on products (merchant_id, status);

create table orders (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  order_number  text not null,
  customer_name text,
  customer_email text,
  status        text not null default 'draft' check (status in ('draft','awaiting_payment','paid','fulfilled','cancelled','refunded')),
  payment_id    text references payments(id) on delete set null,
  invoice_id    text references invoices(id) on delete set null,
  subtotal      numeric(20,8) not null default 0,
  tax_amount    numeric(20,8) not null default 0,
  total_amount  numeric(20,8) not null default 0,
  created_at    timestamptz not null default now(),
  updated_at    timestamptz not null default now(),
  unique (merchant_id, order_number)
);
create index orders_merchant_status_idx on orders (merchant_id, status, created_at);

create table order_items (
  id          text primary key,
  order_id    text not null references orders(id) on delete cascade,
  product_id  text references products(id) on delete set null,
  description text not null,
  quantity    numeric(12,4) not null,
  unit_price  numeric(20,8) not null,
  line_total  numeric(20,8) not null
);
create index order_items_order_idx on order_items (order_id);

-- stock_movements records purchases (in) and sales (out) for inventory tracking.
create table stock_movements (
  id          text primary key,
  merchant_id text not null references merchants(id) on delete cascade,
  product_id  text not null references products(id) on delete cascade,
  qty         numeric(12,4) not null,
  direction   text not null check (direction in ('in','out')),
  reference   text, -- purchase_order_id / order_id
  note        text,
  created_at  timestamptz not null default now()
);
create index stock_movements_product_idx on stock_movements (product_id, created_at);

-- ---- Disputes & Chargebacks ----
create table disputes (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  payment_id      text references payments(id) on delete set null,
  amount          numeric(20,8) not null,
  currency        char(3) not null default 'ETB',
  reason_code     text not null, -- fraud, service_not_received, duplicate, refund_requested, other
  status          text not null default 'open' check (status in ('open','evidence_submitted','won','lost','reversed','closed')),
  evidence        jsonb not null default '[]'::jsonb, -- [{file_key, description}]
  resolution      text,
  fee             numeric(20,8) not null default 0,
  decided_at      timestamptz,
  created_by      text references users(id),
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now()
);
create index disputes_merchant_status_idx on disputes (merchant_id, status);

-- ---- Loyalty & Cashback ----
create table loyalty_tiers (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  name            text not null,
  min_spend       numeric(20,8) not null default 0,
  cashback_percent numeric(5,2) not null default 0,
  created_at      timestamptz not null default now()
);
create index loyalty_tiers_merchant_idx on loyalty_tiers (merchant_id);

create table loyalty_accounts (
  id              text primary key,
  merchant_id     text not null references merchants(id) on delete cascade,
  customer_email  text,
  customer_phone  text,
  points          numeric(14,2) not null default 0,
  tier_id         text references loyalty_tiers(id) on delete set null,
  total_spend     numeric(20,8) not null default 0,
  created_at      timestamptz not null default now(),
  unique (merchant_id, customer_email)
);
create index loyalty_accounts_merchant_idx on loyalty_accounts (merchant_id);

create table cashback_transactions (
  id            text primary key,
  merchant_id   text not null references merchants(id) on delete cascade,
  payment_id    text references payments(id) on delete set null,
  account_id    text references loyalty_accounts(id) on delete cascade,
  amount        numeric(20,8) not null,
  type          text not null check (type in ('earned','redeemed','expired')),
  created_at    timestamptz not null default now()
);
create index cashback_merchant_idx on cashback_transactions (merchant_id, created_at);
