package tax

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// Schedule computes the tax schedule from the tax_register, grouped by period and type.
// It derives collected (from invoices) and paid (remitted), and nets the balance due.
func (r *Repository) Schedule(ctx context.Context, merchantID string) ([]ScheduleLine, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT period, tax_type,
		       COALESCE(SUM(amount),0)::text,
		       COALESCE(SUM(paid),0)::text,
		       COUNT(*)
		FROM tax_register
		WHERE merchant_id=$1
		GROUP BY period, tax_type
		ORDER BY period DESC, tax_type`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []ScheduleLine{}
	for rows.Next() {
		var l ScheduleLine
		if err := rows.Scan(&l.Period, &l.TaxType, &l.Collected, &l.Paid, &l.Count); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

// Record inserts a collected-tax entry (idempotent per source) into the register.
func (r *Repository) Record(ctx context.Context, merchantID, period, taxType, source, sourceID, amount string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO tax_register (id, merchant_id, period, tax_type, source, source_id, amount)
		VALUES ('tax_'||$6||'_'||$4, $1, $2, $3, $4, $5, $6::numeric)
		ON CONFLICT (merchant_id, tax_type, source, source_id) DO NOTHING`,
		merchantID, period, taxType, source, sourceID, amount)
	return err
}
