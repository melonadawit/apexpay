-- Reserve an idempotency key before contacting an external payment connector.
-- Existing rows represent completed historical requests.
ALTER TABLE idempotency_keys
  ADD COLUMN IF NOT EXISTS state text NOT NULL DEFAULT 'completed'
  CHECK (state IN ('in_progress', 'completed', 'failed'));
CREATE INDEX IF NOT EXISTS idempotency_keys_stale_idx
  ON idempotency_keys (created_at) WHERE state = 'in_progress';
