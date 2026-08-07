DROP INDEX IF EXISTS idempotency_keys_connector_started_idx;
ALTER TABLE idempotency_keys DROP CONSTRAINT IF EXISTS idempotency_keys_state_check;
ALTER TABLE idempotency_keys ADD CONSTRAINT idempotency_keys_state_check
  CHECK (state IN ('in_progress', 'completed', 'failed'));
