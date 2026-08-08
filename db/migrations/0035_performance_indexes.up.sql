-- 0035: performance indexes for the hot query paths (dashboard, reports, checkout).

-- Payments: common filters are (merchant_id, status, created_at) and (merchant_id, method).
create index if not exists perf_payments_merchant_status_created
  on payments (merchant_id, status, created_at desc);

-- Ledger: trial balance / P&L / balance sheet group by account over merchant books.
create index if not exists perf_ledger_entries_book
  on ledger_entries (book_id, account_id, journal_id);

-- Ledger journals: report windowing by reference + created_at.
create index if not exists perf_ledger_journals_created
  on ledger_journals (reference_type, reference_id, created_at);

-- Connector health: the router reads last-5m per connector.
create index if not exists perf_connector_health_recent
  on connector_health_samples (connector_id, sampled_at desc);

-- Notifications: delivery worker scans unread recent.
create index if not exists perf_notifications_unread
  on notifications (is_read, created_at desc);

-- Order items lookups.
create index if not exists perf_order_items_order on order_items (order_id);

-- Payroll items: disbursal + report reads by run.
create index if not exists perf_payroll_items_run on payroll_items (run_id);
