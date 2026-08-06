-- Down migration 0016 P0 remaining to surpass RazorpayX — Forex, Accounting Integrations, Credit Lines, Notifications, RM, Support Tickets
drop table if exists support_tickets;
drop table if exists relationship_managers;
drop table if exists notifications;
drop table if exists loan_disbursements;
drop table if exists credit_lines;
drop table if exists accounting_sync_logs;
drop table if exists accounting_integrations;
drop table if exists forex_transactions;
drop table if exists forex_requests;
drop table if exists forex_rates;
