package accounting

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// ---- Journal-entry + fiscal-period persistence (kept separate from report queries). ----

// EnsureOperatingBook returns the merchant's operating ledger book id, creating it (and its
// standard chart of accounts) if it does not exist. The operating book is where manual
// journal entries post.
func (r *Repository) EnsureOperatingBook(ctx context.Context, merchantID string) (string, error) {
	var bookID string
	err := r.pool.QueryRow(ctx, `SELECT id FROM ledger_books WHERE merchant_id=$1 AND book_type='merchant_operating'`, merchantID).Scan(&bookID)
	if err == nil {
		return bookID, nil
	}
	if err != pgx.ErrNoRows {
		return "", err
	}

	// Create the book within a transaction along with a standard chart of accounts.
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	bookID = "book_" + shortID(merchantID) + "_operating"
	if _, err := tx.Exec(ctx, `
		INSERT INTO ledger_books (id, merchant_id, book_type, name, currency, status)
		VALUES ($1,$2,'merchant_operating','Operating ledger','ETB','open')`, bookID, merchantID); err != nil {
		return "", err
	}
	// Standard chart of accounts: asset, liability, equity, revenue, expense.
	accounts := []struct {
		code, name, side string
	}{
		{"asset:bank", "Cash & Bank", "debit"},
		{"asset:receivable", "Accounts Receivable", "debit"},
		{"asset:inventory", "Inventory", "debit"},
		{"asset:fixed", "Fixed Assets", "debit"},
		{"asset:accumulated_depreciation", "Accumulated Depreciation", "credit"},
		{"liability:payable", "Accounts Payable", "credit"},
		{"liability:tax", "Tax Payable", "credit"},
		{"liability:tax", "Tax Payable", "credit"},
		{"equity:owner", "Owner's Equity", "credit"},
		{"equity:retained", "Retained Earnings", "credit"},
		{"revenue:product", "Product Revenue", "credit"},
		{"revenue:service", "Service Revenue", "credit"},
		{"expense:cost_of_sales", "Cost of Sales", "debit"},
		{"expense:operating", "Operating Expenses", "debit"},
		{"expense:admin", "Admin & Overhead", "debit"},
		{"expense:depreciation", "Depreciation", "debit"},
	}
	for i, a := range accounts {
		acid := fmt.Sprintf("%s_%d", bookID, i+1)
		if _, err := tx.Exec(ctx, `
			INSERT INTO ledger_accounts (id, book_id, code, name, normal_balance)
			VALUES ($1,$2,$3,$4,$5)`, acid, bookID, a.code, a.name, a.side); err != nil {
			return "", err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return bookID, nil
}

// AccountIDByCode resolves a ledger account id in a book by its code.
func (r *Repository) AccountIDByCode(ctx context.Context, bookID, code string) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `SELECT id FROM ledger_accounts WHERE book_id=$1 AND code=$2`, bookID, code).Scan(&id)
	return id, err
}

// AccountNameByCode returns the account name for a code in a book ("" if missing).
func (r *Repository) AccountNameByCode(ctx context.Context, bookID, code string) (string, error) {
	var name string
	err := r.pool.QueryRow(ctx, `SELECT name FROM ledger_accounts WHERE book_id=$1 AND code=$2`, bookID, code).Scan(&name)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	return name, err
}

// PeriodStatus returns the status of a fiscal period for a merchant ("open" if none recorded).
func (r *Repository) PeriodStatus(ctx context.Context, merchantID, period string) (string, error) {
	var status string
	err := r.pool.QueryRow(ctx, `SELECT status FROM fiscal_periods WHERE merchant_id=$1 AND period=$2`, merchantID, period).Scan(&status)
	if err == pgx.ErrNoRows {
		return "open", nil
	}
	if err != nil {
		return "", err
	}
	return status, nil
}

// SetPeriodStatus opens or closes a fiscal period. Closing is append-only; reopening is an
// explicit operator action recorded on the row.
func (r *Repository) SetPeriodStatus(ctx context.Context, merchantID, period, status, userID string) error {
	if status == "closed" {
		_, err := r.pool.Exec(ctx, `
			INSERT INTO fiscal_periods (id, merchant_id, period, status, closed_at, closed_by)
			VALUES ($1,$2,$3,'closed',now(),$4)
			ON CONFLICT (merchant_id, period) DO UPDATE SET status='closed', closed_at=now(), closed_by=$4`,
			"fp_"+shortID(merchantID)+"_"+period, merchantID, period, userID)
		return err
	}
	// reopen
	_, err := r.pool.Exec(ctx, `
		UPDATE fiscal_periods SET status='open', closed_at=NULL WHERE merchant_id=$1 AND period=$2`, merchantID, period)
	return err
}

// ListPeriods returns the merchant's fiscal periods newest first.
func (r *Repository) ListPeriods(ctx context.Context, merchantID string) ([]FiscalPeriod, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, merchant_id, period, status, COALESCE(to_char(closed_at,'YYYY-MM-DD"T"HH24:MI:SS'),''), COALESCE(closed_by,'')
		FROM fiscal_periods WHERE merchant_id=$1 ORDER BY period DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []FiscalPeriod{}
	for rows.Next() {
		var p FiscalPeriod
		if err := rows.Scan(&p.ID, &p.Merchant, &p.Period, &p.Status, &p.ClosedAt, &p.ClosedBy); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ListManualJournals returns manual journal entries (reference_type='manual').
func (r *Repository) ListManualJournals(ctx context.Context, merchantID string) ([]JournalEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT j.id, j.memo, to_char(j.created_at,'YYYY-MM'), j.posting_key, j.reference_type,
		       to_char(j.created_at,'YYYY-MM-DD"T"HH24:MI:SS')
		FROM ledger_journals j
		JOIN ledger_books b ON b.id = j.book_id
		WHERE b.merchant_id=$1 AND j.reference_type='manual'
		ORDER BY j.created_at DESC LIMIT 50`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []JournalEntry{}
	for rows.Next() {
		var je JournalEntry
		var refType string
		if err := rows.Scan(&je.ID, &je.Memo, &je.Period, &je.PostingKey, &refType, &je.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, je)
	}
	return out, rows.Err()
}

func shortID(s string) string {
	if len(s) <= 8 {
		return s
	}
	return s[len(s)-8:]
}
