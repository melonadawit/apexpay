DROP INDEX IF EXISTS idempotency_keys_manual_review_idx;
ALTER TABLE idempotency_keys DROP CONSTRAINT IF EXISTS idempotency_keys_state_check;
ALTER TABLE idempotency_keys ADD CONSTRAINT idempotency_keys_state_check
  CHECK (state IN ('in_progress', 'connector_started', 'completed', 'failed'));
