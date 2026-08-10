package budget

import (
	"context"

	"github.com/shopspring/decimal"
)

type repo interface {
	Upsert(ctx context.Context, merchantID string, in BudgetInput) (*Budget, error)
	List(ctx context.Context, merchantID, period string) ([]Budget, error)
	ActualsForPeriod(ctx context.Context, merchantID, period string) (map[string]decimal.Decimal, error)
}

// Service computes budget-vs-actual variance from budgets and ledger actuals.
type Service struct {
	repo repo
}

func NewService(repo repo) *Service { return &Service{repo: repo} }

// SetBudget upserts a budget.
func (s *Service) SetBudget(ctx context.Context, merchantID string, in BudgetInput) (*Budget, error) {
	return s.repo.Upsert(ctx, merchantID, in)
}

// List returns budgets.
func (s *Service) List(ctx context.Context, merchantID, period string) ([]Budget, error) {
	return s.repo.List(ctx, merchantID, period)
}

// Variance builds a budget-vs-actual report for a period. Categories without a budget are
// omitted; budgets without actuals show actual 0.
func (s *Service) Variance(ctx context.Context, merchantID, period string) (*VarianceReport, error) {
	budgets, err := s.repo.List(ctx, merchantID, period)
	if err != nil {
		return nil, err
	}
	actuals, err := s.repo.ActualsForPeriod(ctx, merchantID, period)
	if err != nil {
		return nil, err
	}

	lines := []VarianceLine{}
	for _, b := range budgets {
		actual := actuals[b.Category]
		budgetAmt, _ := decimal.NewFromString(b.BudgetAmount)
		variance := actual.Sub(budgetAmt)
		variancePct := decimal.Zero
		if !budgetAmt.IsZero() {
			variancePct = variance.Div(budgetAmt).Mul(decimal.NewFromInt(100)).Round(2)
		}
		lines = append(lines, VarianceLine{
			Period: b.Period, Category: b.Category,
			BudgetAmount: budgetAmt.String(), ActualAmount: actual.String(),
			Variance: variance.String(), VariancePct: variancePct.String(),
		})
	}
	return &VarianceReport{MerchantID: merchantID, Period: period, Lines: lines}, nil
}
