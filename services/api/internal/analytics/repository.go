package analytics

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// Revenue returns revenue by day for the trailing N days (computed live from payments).
func (r *Repository) Revenue(ctx context.Context, merchantID string, days int) ([]Daily, error) {
	if days <= 0 || days > 365 {
		days = 30
	}
	rows, err := r.pool.Query(ctx, `
		SELECT to_char(date_trunc('day', created_at),'YYYY-MM-DD'),
		       COALESCE(SUM(net_amount) FILTER (WHERE status='succeeded'),0)::text,
		       COALESCE(SUM(amount) FILTER (WHERE status='succeeded'),0)::text,
		       COUNT(*) FILTER (WHERE status IN ('succeeded','failed')),
		       COUNT(*) FILTER (WHERE status='succeeded'),
		       COUNT(*) FILTER (WHERE status='failed'),
		       COALESCE(SUM(amount) FILTER (WHERE status IN ('refunded','partially_refunded')),0)::text
		FROM payments
		WHERE merchant_id=$1 AND created_at >= current_date - ($2 * interval '1 day')
		GROUP BY date_trunc('day', created_at) ORDER BY 1 DESC`, merchantID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Daily{}
	for rows.Next() {
		var d Daily
		var methods string
		if err := rows.Scan(&d.StatDate, &d.Revenue, &d.TPV, &d.PaymentCount, &d.SuccessCount, &d.FailedCount, &d.RefundAmount); err != nil {
			return nil, err
		}
		d.MethodBreakdown = map[string]interface{}{}
		_ = methods
		list = append(list, d)
	}
	return list, rows.Err()
}

// MethodBreakdown returns success by payment method (for the success-by-method chart).
func (r *Repository) MethodBreakdown(ctx context.Context, merchantID string, days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	rows, err := r.pool.Query(ctx, `
		SELECT COALESCE(method,'unknown'),
		       COUNT(*), COUNT(*) FILTER (WHERE status='succeeded'), COALESCE(SUM(amount),0)::text
		FROM payments
		WHERE merchant_id=$1 AND created_at >= current_date - ($2 * interval '1 day')
		GROUP BY method ORDER BY 2 DESC`, merchantID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []map[string]interface{}{}
	for rows.Next() {
		var method string
		var count, succ int
		var rev string
		if err := rows.Scan(&method, &count, &succ, &rev); err != nil {
			return nil, err
		}
		list = append(list, map[string]interface{}{"method": method, "count": count, "success": succ, "revenue": rev})
	}
	return list, rows.Err()
}

// Cohorts returns subscription cohort retention (from stored subscription_cohorts or a
// computed view from subscriptions). Falls back to stored rows.
func (r *Repository) Cohorts(ctx context.Context, merchantID string) ([]Cohort, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT to_char(cohort_month,'YYYY-MM'), customers, month1_retention::text, month2_retention::text,
		       month3_retention::text, mrr::text
		FROM subscription_cohorts WHERE merchant_id=$1 ORDER BY cohort_month DESC LIMIT 12`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Cohort{}
	for rows.Next() {
		var c Cohort
		if err := rows.Scan(&c.CohortMonth, &c.Customers, &c.Month1Retention, &c.Month2Retention, &c.Month3Retention, &c.MRR); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}
