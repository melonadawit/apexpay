package loyalty

import (
	"context"
	"errors"
	"strconv"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// CreateTier adds a loyalty tier.
func (r *Repository) CreateTier(ctx context.Context, merchantID string, t *Tier) error {
	t.ID = id.New("tier")
	_, err := r.pool.Exec(ctx, `INSERT INTO loyalty_tiers (id, merchant_id, name, min_spend, cashback_percent)
		VALUES ($1,$2,$3,$4::numeric,$5::numeric)`,
		t.ID, merchantID, t.Name, t.MinSpend, t.CashbackPercent)
	return err
}

func (r *Repository) ListTiers(ctx context.Context, merchantID string) ([]Tier, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, min_spend::text, cashback_percent::text FROM loyalty_tiers WHERE merchant_id=$1 ORDER BY min_spend`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Tier{}
	for rows.Next() {
		var t Tier
		if err := rows.Scan(&t.ID, &t.Name, &t.MinSpend, &t.CashbackPercent); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

// GetOrCreateAccount upserts a loyalty account by customer email.
func (r *Repository) GetOrCreateAccount(ctx context.Context, merchantID, email string) (*Account, error) {
	if email == "" {
		return nil, errors.New("customer_email required")
	}
	a := &Account{ID: id.New("loyal"), CustomerEmail: email}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO loyalty_accounts (id, merchant_id, customer_email) VALUES ($1,$2,$3)
		ON CONFLICT (merchant_id, customer_email) DO NOTHING`, a.ID, merchantID, email)
	if err != nil {
		return nil, err
	}
	// Read it back with tier name.
	err = r.pool.QueryRow(ctx, `
		SELECT a.id, a.customer_email, COALESCE(a.customer_phone,''), a.points::text,
			COALESCE(a.tier_id,''), COALESCE(t.name,''), a.total_spend::text
		FROM loyalty_accounts a LEFT JOIN loyalty_tiers t ON t.id = a.tier_id
		WHERE a.merchant_id=$1 AND a.customer_email=$2`, merchantID, email).
		Scan(&a.ID, &a.CustomerEmail, &a.CustomerPhone, &a.Points, &a.TierID, &a.TierName, &a.TotalSpend)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// EarnCashback records earned cashback and adds to the account.
func (r *Repository) EarnCashback(ctx context.Context, merchantID, accountID, paymentID, amount string) (*CashbackTx, error) {
	tx := &CashbackTx{ID: id.New("cb"), PaymentID: paymentID, Amount: amount, Type: "earned", CreatedAt: nowStr()}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO cashback_transactions (id, merchant_id, payment_id, account_id, amount, type)
		VALUES ($1,$2,$3,$4,$5::numeric,'earned')`, tx.ID, merchantID, paymentID, accountID, amount)
	if err != nil {
		return nil, err
	}
	_, err = r.pool.Exec(ctx, `UPDATE loyalty_accounts SET points = points + $2::numeric, total_spend = total_spend + $2::numeric WHERE id=$1`,
		accountID, amount)
	return tx, err
}

// Accounts returns loyalty accounts for a merchant.
func (r *Repository) Accounts(ctx context.Context, merchantID string, limit int) ([]Account, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT a.id, COALESCE(a.customer_email,''), COALESCE(a.customer_phone,''), a.points::text,
			COALESCE(a.tier_id,''), COALESCE(t.name,''), a.total_spend::text
		FROM loyalty_accounts a LEFT JOIN loyalty_tiers t ON t.id = a.tier_id
		WHERE a.merchant_id=$1 ORDER BY a.total_spend DESC LIMIT $2`, merchantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Account{}
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.CustomerEmail, &a.CustomerPhone, &a.Points, &a.TierID, &a.TierName, &a.TotalSpend); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

// Transactions returns cashback history.
func (r *Repository) Transactions(ctx context.Context, merchantID string, limit int) ([]CashbackTx, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE(payment_id,''), amount::text, type,
			to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM cashback_transactions WHERE merchant_id=$1 ORDER BY created_at DESC LIMIT $2`, merchantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []CashbackTx{}
	for rows.Next() {
		var c CashbackTx
		if err := rows.Scan(&c.ID, &c.PaymentID, &c.Amount, &c.Type, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func itoa(n int) string { return strconv.Itoa(n) }

func nowStr() string        { return time.Now().In(tzEAT()).Format(time.RFC3339) }
func tzEAT() *time.Location { return time.FixedZone("EAT", 3*3600) }
