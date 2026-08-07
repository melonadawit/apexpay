DROP INDEX IF EXISTS idempotency_keys_stale_idx;
ALTER TABLE idempotency_keys DROP COLUMN IF EXISTS state;
