-- Operations decisions are durable, auditable, and keyed to the original idempotency request.
ALTER TABLE idempotency_keys DROP CONSTRAINT IF EXISTS idempotency_keys_state_check;
ALTER TABLE idempotency_keys ADD CONSTRAINT idempotency_keys_state_check
  CHECK (state IN ('in_progress','connector_started','completed','failed','manual_review','retry_authorized'));
CREATE TABLE IF NOT EXISTS payment_reconciliation_cases (
  merchant_id text NOT NULL REFERENCES merchants(id),
  idempotency_key text NOT NULL,
  tx_ref text,
  status text NOT NULL CHECK (status IN ('open','confirmed_paid','confirmed_not_paid','requires_connector_investigation')) DEFAULT 'open',
  reviewer_id text,
  reviewer_note text,
  decided_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (merchant_id, idempotency_key)
);
CREATE INDEX IF NOT EXISTS payment_reconciliation_cases_status_idx ON payment_reconciliation_cases (status, created_at);
