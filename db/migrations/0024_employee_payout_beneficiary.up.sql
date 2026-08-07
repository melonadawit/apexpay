ALTER TABLE employees ADD COLUMN IF NOT EXISTS payout_beneficiary_id text REFERENCES beneficiaries(id);
CREATE INDEX IF NOT EXISTS employees_payout_beneficiary_idx ON employees (payout_beneficiary_id) WHERE payout_beneficiary_id IS NOT NULL;
