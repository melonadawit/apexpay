package lending

import (
	"context"
	"errors"
	"strconv"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// CreateLoan creates a loan disbursement from a credit line.
func (r *Repository) CreateLoan(ctx context.Context, merchantID, userID, creditLineID string, l *Loan) error {
	l.ID = id.New("loan")
	if l.Currency == "" {
		l.Currency = "ETB"
	}
	if l.InterestRate == "" {
		l.InterestRate = "18.00"
	}
	if l.Status == "" {
		l.Status = "pending"
	}
	l.RepaidAmount = "0"
	l.Outstanding = l.Amount
	due := time.Now().AddDate(0, 3, 0) // 3-month micro-loan default
	if l.DueDate != "" {
		if t, err := time.Parse("2006-01-02", l.DueDate); err == nil {
			due = t
		}
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO loan_disbursements (id, credit_line_id, merchant_id, amount, currency, purpose, status, due_date, repaid_amount, outstanding_amount, created_by)
		VALUES ($1,$2,$3,$4::numeric,$5,$6,$7,$8::date,0,$9::numeric,$10)`,
		l.ID, creditLineID, merchantID, l.Amount, l.Currency, l.Purpose, l.Status, due.Format("2006-01-02"), l.Outstanding, userID)
	if err != nil {
		return err
	}
	l.DueDate = due.Format("2006-01-02")
	return nil
}

// ListLoans returns a merchant's loans.
func (r *Repository) ListLoans(ctx context.Context, merchantID, status string, limit int) ([]Loan, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query := `SELECT ld.id, ld.amount::text, ld.currency, ld.purpose, ld.status,
		COALESCE(cl.interest_rate,0)::text,
		COALESCE(to_char(ld.due_date,'YYYY-MM-DD'),''), ld.repaid_amount::text, ld.outstanding_amount::text,
		to_char(ld.created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM loan_disbursements ld
		LEFT JOIN credit_lines cl ON cl.id = ld.credit_line_id
		WHERE ld.merchant_id=$1`
	args := []interface{}{merchantID}
	if status != "" {
		query += ` AND ld.status=$2`
		args = append(args, status)
	}
	query += ` ORDER BY ld.created_at DESC LIMIT $` + itoa(len(args)+1)
	args = append(args, limit)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Loan{}
	for rows.Next() {
		var l Loan
		if err := rows.Scan(&l.ID, &l.Amount, &l.Currency, &l.Purpose, &l.Status, &l.InterestRate,
			&l.DueDate, &l.RepaidAmount, &l.Outstanding, &l.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, rows.Err()
}

// Repay reduces the outstanding balance of a loan.
func (r *Repository) Repay(ctx context.Context, merchantID, loanID, amount string) error {
	amt, _ := decimal.NewFromString(amount)
	ct, err := r.pool.Exec(ctx, `
		UPDATE loan_disbursements SET
			repaid_amount = repaid_amount + $1::numeric,
			outstanding_amount = GREATEST(outstanding_amount - $1::numeric, 0),
			status = CASE WHEN outstanding_amount - $1::numeric <= 0 THEN 'repaid' ELSE status END
		WHERE merchant_id=$2 AND id=$3 AND status IN ('disbursed','pending')`, amt.String(), merchantID, loanID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func itoa(n int) string { return strconv.Itoa(n) }
