package budget

import (
	"context"
	"fmt"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// Upsert creates or updates a budget for a (merchant, period, category).
func (r *Repository) Upsert(ctx context.Context, merchantID string, in BudgetInput) (*Budget, error) {
	amt, err := decimal.NewFromString(in.BudgetAmount)
	if err != nil {
		return nil, fmt.Errorf("invalid budget_amount")
	}
	b := &Budget{ID: id.New("budg"), MerchantID: merchantID, Period: in.Period, Category: in.Category, BudgetAmount: amt.String()}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO budgets (id, merchant_id, period, category, budget_amount)
		VALUES ($1,$2,$3,$4,$5::numeric)
		ON CONFLICT (merchant_id, period, category)
		DO UPDATE SET budget_amount = EXCLUDED.budget_amount`,
		b.ID, merchantID, b.Period, b.Category, b.BudgetAmount)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// List returns budgets for a merchant, optionally filtered by period.
func (r *Repository) List(ctx context.Context, merchantID, period string) ([]Budget, error) {
	query := `SELECT id, merchant_id, period, category, budget_amount::text, to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SS')
		FROM budgets WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if period != "" {
		query += ` AND period=$2`
		args = append(args, period)
	}
	query += ` ORDER BY period DESC, category`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Budget{}
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.ID, &b.MerchantID, &b.Period, &b.Category, &b.BudgetAmount, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// ActualsForPeriod computes the actual revenue/expense from the ledger for a period,
// grouped by top-level category (revenue|expense) or a specific category code.
func (r *Repository) ActualsForPeriod(ctx context.Context, merchantID, period string) (map[string]decimal.Decimal, error) {
	from := period + "-01"
	rows, err := r.pool.Query(ctx, `
		SELECT split_part(a.code, ':', 1) AS category,
		       SUM(CASE WHEN e.direction='debit' THEN e.amount ELSE -e.amount END)::text AS net
		FROM ledger_entries e
		JOIN ledger_books b ON b.id = e.book_id
		JOIN ledger_accounts a ON a.id = e.account_id
		JOIN ledger_journals j ON j.id = e.journal_id
		WHERE b.merchant_id=$1 AND j.created_at >= $2::date
		  AND j.created_at < ($2::date + interval '1 month')
		  AND (a.code LIKE 'revenue%' OR a.code LIKE 'expense%')
		GROUP BY category`, merchantID, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]decimal.Decimal{}
	for rows.Next() {
		var cat, s string
		if err := rows.Scan(&cat, &s); err != nil {
			return nil, err
		}
		d, _ := decimal.NewFromString(s)
		out[cat] = d
	}
	return out, rows.Err()
}
