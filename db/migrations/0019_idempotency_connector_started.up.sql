ALTER TABLE idempotency_keys DROP CONSTRAINT IF EXISTS idempotency_keys_state_check;
ALTER TABLE idempotency_keys
  ADD CONSTRAINT idempotency_keys_state_check
  CHECK (state IN ('in_progress', 'connector_started', 'completed', 'failed'));
CREATE INDEX IF NOT EXISTS idempotency_keys_connector_started_idx
  ON idempotency_keys (created_at) WHERE state = 'connector_started';
