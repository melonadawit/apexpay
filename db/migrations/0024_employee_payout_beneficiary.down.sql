DROP INDEX IF EXISTS employees_payout_beneficiary_idx;
ALTER TABLE employees DROP COLUMN IF EXISTS payout_beneficiary_id;
