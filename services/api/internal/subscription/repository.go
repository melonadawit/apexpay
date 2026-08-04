package subscription

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct{ pool *pgxpool.Pool }

func NewPgRepository(pool *pgxpool.Pool) *PgRepository { return &PgRepository{pool: pool} }

func (r *PgRepository) CreatePlan(ctx context.Context, p *Plan) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO subscription_plans (id, merchant_id, name, description, amount, currency, interval_type, interval_count, trial_days, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		p.ID, p.MerchantID, p.Name, p.Description, p.Amount.String(), p.Currency, p.IntervalType, p.IntervalCount, p.TrialDays, p.Status)
	return err
}

func (r *PgRepository) GetPlan(ctx context.Context, merchantID, planID string) (*Plan, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, name, amount::text, currency, interval_type, interval_count, trial_days, status FROM subscription_plans WHERE merchant_id=$1 AND id=$2`, merchantID, planID)
	var pl Plan
	var amt string
	err := row.Scan(&pl.ID, &pl.MerchantID, &pl.Name, &amt, &pl.Currency, &pl.IntervalType, &pl.IntervalCount, &pl.TrialDays, &pl.Status)
	return &pl, err
}

func (r *PgRepository) CreateCustomer(ctx context.Context, c *Customer) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO customers (id, merchant_id, email, phone, name) VALUES ($1,$2,$3,$4,$5)`, c.ID, c.MerchantID, c.Email, c.Phone, c.Name)
	return err
}

func (r *PgRepository) CreateSubscription(ctx context.Context, s *Subscription) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO subscriptions (id, merchant_id, customer_id, plan_id, status, current_period_start, current_period_end, trial_end) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		s.ID, s.MerchantID, s.CustomerID, s.PlanID, s.Status, s.CurrentPeriodStart, s.CurrentPeriodEnd, s.TrialEnd)
	return err
}

func (r *PgRepository) CreateInvoice(ctx context.Context, inv *Invoice) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO subscription_invoices (id, merchant_id, subscription_id, amount, currency, status, attempt_count, due_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		inv.ID, inv.MerchantID, inv.SubscriptionID, inv.Amount.String(), inv.Currency, inv.Status, inv.AttemptCount, inv.DueAt)
	return err
}

func (r *PgRepository) ListSubscriptions(ctx context.Context, merchantID string, status *SubscriptionStatus) ([]Subscription, error) {
	query := `SELECT id, merchant_id, customer_id, plan_id, status, current_period_start, current_period_end FROM subscriptions WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if status != nil {
		query += ` AND status=$2`
		args = append(args, *status)
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Subscription
	for rows.Next() {
		var s Subscription
		if err := rows.Scan(&s.ID, &s.MerchantID, &s.CustomerID, &s.PlanID, &s.Status, &s.CurrentPeriodStart, &s.CurrentPeriodEnd); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *PgRepository) UpdateSubscriptionStatus(ctx context.Context, id string, status SubscriptionStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE subscriptions SET status=$1 WHERE id=$2`, status, id)
	return err
}
