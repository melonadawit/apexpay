package risk

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// Service orchestrates rule evaluation against live aggregates and persists flags.
type Service struct {
	repo   *Repository
	engine *Engine
}

func NewService(repo *Repository, engine *Engine) *Service {
	return &Service{repo: repo, engine: engine}
}

// Evaluate runs the merchant's rules against a transaction context, populating the
// window aggregates from the DB, and persists any matched findings as flags.
// Returns the evaluation so the caller can decide whether to block the payment.
func (s *Service) Evaluate(ctx context.Context, in EvaluationContext) (Evaluation, error) {
	rules, err := s.repo.ListRules(ctx, in.MerchantID)
	if err != nil {
		return Evaluation{}, err
	}
	// Populate window aggregates from the DB for each rule's window.
	ctx2 := in
	ctx2.AmountInWindow, _ = s.repo.AmountInWindow(ctx, in.MerchantID, 60*time.Minute)
	ctx2.CountInWindow, _ = s.repo.CountInWindow(ctx, in.MerchantID, 60*time.Minute)
	ctx2.FailureRate, _ = s.repo.FailureRateInWindow(ctx, in.MerchantID, 60*time.Minute)

	eval := s.engine.Evaluate(ctx2, rules)

	// Persist matched findings as flags for the operations/review console.
	for _, f := range eval.Matched() {
		_ = s.repo.CreateFlag(ctx, &Flag{
			MerchantID: in.MerchantID,
			EntityType: in.EntityType,
			EntityID:   in.EntityID,
			RuleID:     f.RuleID,
			RuleName:   f.RuleName,
			Severity:   f.Severity,
			Action:     f.Action,
			Reason:     f.Reason,
			Details: map[string]interface{}{
				"amount": in.Amount.String(), "device_id": in.DeviceID, "ip": in.IP, "customer_id": in.CustomerID,
			},
		})
	}
	return eval, nil
}

// DefaultRules returns a sensible set of global rules for a merchant that has none.
func DefaultRules() []Rule {
	return []Rule{
		{RuleType: "threshold_amount", Name: "High-ticket payment", Action: "flag", Severity: "high", Parameters: RuleParams{AmountLimit: decimal.NewFromInt(500000)}, Enabled: true},
		{RuleType: "velocity_amount", Name: "Amount velocity", Action: "flag", Severity: "medium", Parameters: RuleParams{WindowMinutes: 60, AmountLimit: decimal.NewFromInt(2000000)}, Enabled: true},
		{RuleType: "velocity_count", Name: "Transaction count velocity", Action: "review", Severity: "medium", Parameters: RuleParams{WindowMinutes: 60, CountLimit: 50}, Enabled: true},
		{RuleType: "high_failure_rate", Name: "High failure rate", Action: "flag", Severity: "medium", Parameters: RuleParams{WindowMinutes: 60, RatioLimit: decimal.NewFromFloat(0.5)}, Enabled: true},
	}
}
