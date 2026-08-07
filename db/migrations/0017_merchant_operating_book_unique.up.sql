-- One active operating ledger per merchant prevents ambiguous payment posting.
CREATE UNIQUE INDEX IF NOT EXISTS ledger_books_merchant_type_uidx
  ON ledger_books (merchant_id, book_type);
