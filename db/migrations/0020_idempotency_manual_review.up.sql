-- External-call ambiguity is a financial operations case, not a retry case.
ALTER TABLE idempotency_keys DROP CONSTRAINT IF EXISTS idempotency_keys_state_check;
ALTER TABLE idempotency_keys
  ADD CONSTRAINT idempotency_keys_state_check
  CHECK (state IN ('in_progress', 'connector_started', 'completed', 'failed', 'manual_review'));
CREATE INDEX IF NOT EXISTS idempotency_keys_manual_review_idx
  ON idempotency_keys (created_at) WHERE state = 'manual_review';
