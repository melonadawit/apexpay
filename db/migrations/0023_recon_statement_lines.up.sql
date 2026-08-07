CREATE TABLE IF NOT EXISTS recon_statement_lines (
  id text PRIMARY KEY,
  statement_id text NOT NULL REFERENCES recon_statements(id) ON DELETE CASCADE,
  external_transaction_id text NOT NULL,
  connector_ref text,
  amount numeric(20,8) NOT NULL CHECK (amount > 0),
  currency char(3) NOT NULL,
  occurred_at timestamptz NOT NULL,
  match_status text NOT NULL CHECK (match_status IN ('unmatched','matched','ambiguous')) DEFAULT 'unmatched',
  matched_journal_id text REFERENCES ledger_journals(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (statement_id, external_transaction_id)
);
CREATE INDEX IF NOT EXISTS recon_statement_lines_unmatched_idx ON recon_statement_lines(statement_id,match_status);
CREATE INDEX IF NOT EXISTS recon_statement_lines_connector_ref_idx ON recon_statement_lines(connector_ref);
