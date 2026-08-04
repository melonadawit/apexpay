package link

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct{ pool *pgxpool.Pool }

func NewPgRepository(pool *pgxpool.Pool) *PgRepository { return &PgRepository{pool: pool} }

func (r *PgRepository) CreateLink(ctx context.Context, pl *PaymentLink) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO payment_links (id, merchant_id, amount, currency, description, status, public_token, expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		pl.ID, pl.MerchantID, pl.Amount.String(), pl.Currency, pl.Description, pl.Status, pl.PublicToken, pl.ExpiresAt)
	return err
}

func (r *PgRepository) GetByToken(ctx context.Context, token string) (*PaymentLink, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, amount::text, currency, description, status, public_token FROM payment_links WHERE public_token=$1`, token)
	var pl PaymentLink
	var amt string
	err := row.Scan(&pl.ID, &pl.MerchantID, &amt, &pl.Currency, &pl.Description, &pl.Status, &pl.PublicToken)
	return &pl, err
}

func (r *PgRepository) ListByMerchant(ctx context.Context, merchantID string) ([]PaymentLink, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, amount::text, currency, description, status, public_token, created_at FROM payment_links WHERE merchant_id=$1 ORDER BY created_at DESC LIMIT 100`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []PaymentLink
	for rows.Next() {
		var pl PaymentLink
		var amt string
		if err := rows.Scan(&pl.ID, &amt, &pl.Currency, &pl.Description, &pl.Status, &pl.PublicToken, &pl.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, pl)
	}
	return list, nil
}

func (r *PgRepository) MarkPaid(ctx context.Context, linkID, paymentID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payment_links SET status='paid', payment_id=$1, updated_at=now() WHERE id=$2`, paymentID, linkID)
	return err
}

func (r *PgRepository) CreateCheckoutSession(ctx context.Context, cs *CheckoutSession) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO checkout_sessions (id, merchant_id, payment_id, payment_link_id, public_token, status, expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		cs.ID, cs.MerchantID, cs.PaymentID, cs.PaymentLinkID, cs.PublicToken, cs.Status, cs.ExpiresAt)
	return err
}
