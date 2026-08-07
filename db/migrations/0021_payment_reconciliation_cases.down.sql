DROP TABLE IF EXISTS payment_reconciliation_cases;
ALTER TABLE idempotency_keys DROP CONSTRAINT IF EXISTS idempotency_keys_state_check;
ALTER TABLE idempotency_keys ADD CONSTRAINT idempotency_keys_state_check
  CHECK (state IN ('in_progress','connector_started','completed','failed','manual_review'));
