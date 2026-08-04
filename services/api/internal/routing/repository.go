package routing

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"time"
)

type PgRepository struct{ pool *pgxpool.Pool }

func NewPgRepository(pool *pgxpool.Pool) *PgRepository { return &PgRepository{pool: pool} }

func (r *PgRepository) ListRules(ctx context.Context, merchantID *string) ([]RoutingRule, error) {
	query := `SELECT id, merchant_id, name, min_amount::text, max_amount::text, currency, payment_method, primary_connector, fallback1, fallback2, strategy, enabled, priority FROM routing_rules WHERE enabled=true AND (merchant_id = $1 OR ($1 IS NULL AND merchant_id IS NULL)) ORDER BY priority ASC`
	rows, err := r.pool.Query(ctx, query, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []RoutingRule
	for rows.Next() {
		var rr RoutingRule
		var minStr, maxStr *string
		var pm *string
		var f1, f2 *string
		var merchID *string
		if err := rows.Scan(&rr.ID, &merchID, &rr.Name, &minStr, &maxStr, &rr.Currency, &pm, &rr.PrimaryConnector, &f1, &f2, &rr.Strategy, &rr.Enabled, &rr.Priority); err != nil {
			return nil, err
		}
		if minStr != nil {
			dec, _ := decimal.NewFromString(*minStr)
			rr.MinAmount = &dec
		}
		if maxStr != nil {
			dec, _ := decimal.NewFromString(*maxStr)
			rr.MaxAmount = &dec
		}
		if merchID != nil {
			rr.MerchantID = merchID
		}
		if pm != nil {
			rr.PaymentMethod = pm
		}
		if f1 != nil {
			c := ConnectorID(*f1)
			rr.Fallback1 = &c
		}
		if f2 != nil {
			c := ConnectorID(*f2)
			rr.Fallback2 = &c
		}
		rules = append(rules, rr)
	}
	return rules, nil
}

func (r *PgRepository) ListHealthSamples(ctx context.Context, connectorID ConnectorID, since time.Time) ([]HealthSample, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, connector_id, environment, latency_ms, success, error_code, sampled_at FROM connector_health_samples WHERE connector_id=$1 AND sampled_at >= $2 ORDER BY sampled_at DESC LIMIT 100`, connectorID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var samples []HealthSample
	for rows.Next() {
		var hs HealthSample
		if err := rows.Scan(&hs.ID, &hs.ConnectorID, &hs.Environment, &hs.LatencyMS, &hs.Success, &hs.ErrorCode, &hs.SampledAt); err != nil {
			return nil, err
		}
		samples = append(samples, hs)
	}
	return samples, nil
}

func (r *PgRepository) SaveHealthSample(ctx context.Context, s HealthSample) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO connector_health_samples (id, connector_id, environment, latency_ms, success, error_code, sampled_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		s.ID, s.ConnectorID, s.Environment, s.LatencyMS, s.Success, s.ErrorCode, s.SampledAt)
	return err
}
