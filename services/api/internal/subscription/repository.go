package subscription

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
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

// GetSubscriptionDetail returns a single subscription joined with its plan,
// customer and invoices for the detail page.
func (r *PgRepository) GetSubscriptionDetail(ctx context.Context, merchantID, subID string) (*SubscriptionDetail, error) {
	var s Subscription
	var trialEnd, cancelAt *time.Time
	err := r.pool.QueryRow(ctx, `SELECT id, merchant_id, customer_id, plan_id, status, current_period_start, current_period_end, trial_end, cancel_at FROM subscriptions WHERE merchant_id=$1 AND id=$2`, merchantID, subID).
		Scan(&s.ID, &s.MerchantID, &s.CustomerID, &s.PlanID, &s.Status, &s.CurrentPeriodStart, &s.CurrentPeriodEnd, &trialEnd, &cancelAt)
	if err != nil {
		return nil, err
	}
	s.TrialEnd = trialEnd
	s.CancelAt = cancelAt

	detail := &SubscriptionDetail{Subscription: &s}

	var plan Plan
	var amount string
	err = r.pool.QueryRow(ctx, `SELECT id, merchant_id, name, COALESCE(description,''), amount::text, currency, interval_type, interval_count, trial_days, status FROM subscription_plans WHERE merchant_id=$1 AND id=$2`, merchantID, s.PlanID).
		Scan(&plan.ID, &plan.MerchantID, &plan.Name, &plan.Description, &amount, &plan.Currency, &plan.IntervalType, &plan.IntervalCount, &plan.TrialDays, &plan.Status)
	if err == nil {
		plan.Amount, _ = decimal.NewFromString(amount)
		detail.Plan = &plan
	}

	var cust Customer
	err = r.pool.QueryRow(ctx, `SELECT id, merchant_id, COALESCE(email,''), COALESCE(phone,''), COALESCE(name,''), COALESCE(fayda_fin_hash,'') FROM customers WHERE merchant_id=$1 AND id=$2`, merchantID, s.CustomerID).
		Scan(&cust.ID, &cust.MerchantID, &cust.Email, &cust.Phone, &cust.Name, &cust.FinHash)
	if err == nil {
		detail.Customer = &cust
	}

	invRows, err := r.pool.Query(ctx, `SELECT id, merchant_id, subscription_id, payment_id, amount::text, currency, status, attempt_count, due_at, created_at FROM subscription_invoices WHERE merchant_id=$1 AND subscription_id=$2 ORDER BY created_at DESC`, merchantID, subID)
	if err == nil {
		defer invRows.Close()
		for invRows.Next() {
			var inv Invoice
			var invAmount string
			if err := invRows.Scan(&inv.ID, &inv.MerchantID, &inv.SubscriptionID, &inv.PaymentID, &invAmount, &inv.Currency, &inv.Status, &inv.AttemptCount, &inv.DueAt, &inv.CreatedAt); err == nil {
				inv.Amount, _ = decimal.NewFromString(invAmount)
				detail.Invoices = append(detail.Invoices, inv)
			}
		}
	}

	return detail, nil
}
