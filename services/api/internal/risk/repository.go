package risk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// ListRules returns enabled rules relevant to a merchant (global + merchant-specific),
// merchant-specific overrides global by rule_type.
func (r *Repository) ListRules(ctx context.Context, merchantID string) ([]Rule, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE(merchant_id,''), name, COALESCE(description,''), rule_type,
		       parameters::text, action, severity, enabled
		FROM risk_rules
		WHERE enabled=true AND (merchant_id IS NULL OR merchant_id=$1)
		ORDER BY merchant_id NULLS FIRST`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Rule
	for rows.Next() {
		var rule Rule
		var m, params, desc string
		if err := rows.Scan(&rule.ID, &m, &rule.Name, &desc, &rule.RuleType, &params, &rule.Action, &rule.Severity, &rule.Enabled); err != nil {
			return nil, err
		}
		if m != "" {
			rule.MerchantID = &m
		}
		rule.Description = desc
		_ = json.Unmarshal([]byte(params), &rule.Parameters)
		list = append(list, rule)
	}
	return list, rows.Err()
}

// CreateRule inserts a rule and returns it.
func (r *Repository) CreateRule(ctx context.Context, merchantID string, rule *Rule) error {
	rule.ID = id.New("rule")
	params, _ := json.Marshal(rule.Parameters)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO risk_rules (id, merchant_id, name, description, rule_type, parameters, action, severity, enabled)
		VALUES ($1,$2,$3,$4,$5,$6::jsonb,$7,$8,$9)`,
		rule.ID, nilStr(merchantID), rule.Name, rule.Description, rule.RuleType, params, rule.Action, rule.Severity, rule.Enabled)
	return err
}

// --- Window aggregates for evaluation ---

// AmountInWindow returns the sum of succeeded+failed payment amounts in the window.
func (r *Repository) AmountInWindow(ctx context.Context, merchantID string, window time.Duration) (decimal.Decimal, error) {
	var s string
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount),0)::text FROM payments
		WHERE merchant_id=$1 AND created_at >= now() - ($2 || ' minutes')::interval
		  AND status IN ('succeeded','failed')`, merchantID, int(window.Minutes())).Scan(&s)
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromString(s)
}

// CountInWindow returns the number of payments in the window.
func (r *Repository) CountInWindow(ctx context.Context, merchantID string, window time.Duration) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM payments
		WHERE merchant_id=$1 AND created_at >= now() - ($2 || ' minutes')::interval
		  AND status IN ('succeeded','failed')`, merchantID, int(window.Minutes())).Scan(&n)
	return n, err
}

// FailureRateInWindow returns the failure ratio in the window (0..1).
func (r *Repository) FailureRateInWindow(ctx context.Context, merchantID string, window time.Duration) (decimal.Decimal, error) {
	var f float64
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(COUNT(*) FILTER (WHERE status='failed')::float / NULLIF(COUNT(*),0),0)
		FROM payments WHERE merchant_id=$1 AND created_at >= now() - ($2 || ' minutes')::interval
		  AND status IN ('succeeded','failed')`, merchantID, int(window.Minutes())).Scan(&f)
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromFloat(f), nil
}

// --- Flags ---

type Flag struct {
	ID         string
	MerchantID string
	EntityType string
	EntityID   string
	RuleID     string
	RuleName   string
	Severity   string
	Action     string
	Reason     string
	Details    map[string]interface{}
	Status     string
}

// CreateFlag inserts a risk flag.
func (r *Repository) CreateFlag(ctx context.Context, f *Flag) error {
	f.ID = id.New("rflag")
	details, _ := json.Marshal(f.Details)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO risk_flags (id, merchant_id, entity_type, entity_id, rule_id, rule_name, severity, action, reason, details, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10::jsonb,'open')`,
		f.ID, f.MerchantID, f.EntityType, f.EntityID, f.RuleID, f.RuleName, f.Severity, f.Action, f.Reason, details)
	return err
}

// ListFlags returns open flags for a merchant, newest first.
func (r *Repository) ListFlags(ctx context.Context, merchantID string, status string, limit int) ([]Flag, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query := `SELECT id, merchant_id, entity_type, entity_id, COALESCE(rule_id,''), COALESCE(rule_name,''), severity, action, COALESCE(reason,''), details::text, status FROM risk_flags WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if status != "" {
		query += ` AND status=$2`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Flag
	for rows.Next() {
		var f Flag
		var details string
		if err := rows.Scan(&f.ID, &f.MerchantID, &f.EntityType, &f.EntityID, &f.RuleID, &f.RuleName,
			&f.Severity, &f.Action, &f.Reason, &details, &f.Status); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(details), &f.Details)
		list = append(list, f)
	}
	return list, rows.Err()
}

func nilStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
