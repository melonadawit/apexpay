package budget

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
)

type fakeRepo struct {
	budgets []Budget
	actuals map[string]decimal.Decimal
}

func (f *fakeRepo) Upsert(ctx context.Context, merchantID string, in BudgetInput) (*Budget, error) {
	b := &Budget{ID: "budg_1", MerchantID: merchantID, Period: in.Period, Category: in.Category, BudgetAmount: in.BudgetAmount}
	f.budgets = append(f.budgets, *b)
	return b, nil
}
func (f *fakeRepo) List(ctx context.Context, merchantID, period string) ([]Budget, error) {
	return f.budgets, nil
}
func (f *fakeRepo) ActualsForPeriod(ctx context.Context, merchantID, period string) (map[string]decimal.Decimal, error) {
	return f.actuals, nil
}

func TestVarianceComputation(t *testing.T) {
	repo := &fakeRepo{
		budgets: []Budget{
			{Period: "2026-08", Category: "expense", BudgetAmount: "10000"},
			{Period: "2026-08", Category: "revenue", BudgetAmount: "50000"},
		},
		actuals: map[string]decimal.Decimal{
			"expense": decimal.RequireFromString("12000"),
			"revenue": decimal.RequireFromString("60000"),
		},
	}
	svc := NewService(repo)
	report, err := svc.Variance(context.Background(), "m1", "2026-08")
	if err != nil {
		t.Fatalf("Variance: %v", err)
	}
	if len(report.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(report.Lines))
	}
	// expense: actual 12000 - budget 10000 = +2000 variance, +20%
	find := func(cat string) *VarianceLine {
		for i := range report.Lines {
			if report.Lines[i].Category == cat {
				return &report.Lines[i]
			}
		}
		return nil
	}
	exp := find("expense")
	if exp == nil || exp.Variance != "2000" || exp.VariancePct != "20" {
		t.Fatalf("expense variance wrong: %+v", exp)
	}
	rev := find("revenue")
	if rev == nil || rev.Variance != "10000" || rev.VariancePct != "20" {
		t.Fatalf("revenue variance wrong: %+v", rev)
	}
}

func TestSetBudget(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo)
	b, err := svc.SetBudget(context.Background(), "m1", BudgetInput{Period: "2026-09", Category: "expense", BudgetAmount: "1500"})
	if err != nil {
		t.Fatalf("SetBudget: %v", err)
	}
	if b.BudgetAmount != "1500" {
		t.Fatalf("expected 1500, got %s", b.BudgetAmount)
	}
}
