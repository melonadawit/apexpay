package fxreval

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// ForeignAccounts returns active non-ETB current accounts with their balance and the
// current forex rate (ETB per unit) for their currency. Accounts without a rate are skipped.
func (r *Repository) ForeignAccounts(ctx context.Context, merchantID string) ([]RevaluationLine, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT ca.id, ca.account_number, ca.currency, ca.balance::text, fr.rate::text
		FROM current_accounts ca
		JOIN forex_rates fr ON fr.to_currency = ca.currency AND fr.from_currency = 'ETB'
		WHERE ca.merchant_id=$1 AND ca.status='active' AND ca.currency <> 'ETB'`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []RevaluationLine{}
	for rows.Next() {
		var l RevaluationLine
		if err := rows.Scan(&l.AccountID, &l.AccountNumber, &l.Currency, &l.AmountFX, &l.Rate); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

// PriorETB returns the most recent recorded ETB value for an account (0 if none).
func (r *Repository) PriorETB(ctx context.Context, merchantID, accountID string) (decimal.Decimal, error) {
	var s string
	err := r.pool.QueryRow(ctx, `
		SELECT amount_etb::text FROM fx_revaluations
		WHERE merchant_id=$1 AND account_id=$2
		ORDER BY created_at DESC LIMIT 1`, merchantID, accountID).Scan(&s)
	if err == pgx.ErrNoRows {
		return decimal.Zero, nil
	}
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromString(s)
}

// Record persists a revaluation line for the period (idempotent per account+period).
func (r *Repository) Record(ctx context.Context, merchantID string, line RevaluationLine, period string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO fx_revaluations (id, merchant_id, period, account_id, currency, amount_fx, rate, amount_etb, fx_gain)
		VALUES ('fxr_'||$4||'_'||$3, $1, $2, $4, $5, $6::numeric, $7::numeric, $8::numeric, $9::numeric)
		ON CONFLICT (id) DO NOTHING`,
		merchantID, period, line.AccountID, line.AccountID, line.Currency,
		line.AmountFX, line.Rate, line.AmountETB, line.FXGainLoss)
	if err != nil {
		return fmt.Errorf("record fx revaluation: %w", err)
	}
	return nil
}
