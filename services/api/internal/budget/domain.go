package budget

// Budgeting / FP&A domain types.

// Budget is a planned amount for a period and category.
type Budget struct {
	ID           string `json:"id"`
	MerchantID   string `json:"merchant_id"`
	Period       string `json:"period"`   // YYYY-MM
	Category     string `json:"category"` // revenue | expense | cost center code
	BudgetAmount string `json:"budget_amount"`
	CreatedAt    string `json:"created_at"`
}

type BudgetInput struct {
	Period       string `json:"period"`
	Category     string `json:"category"`
	BudgetAmount string `json:"budget_amount"`
}

// VarianceLine is budget vs actual for one category in a period.
type VarianceLine struct {
	Period       string `json:"period"`
	Category     string `json:"category"`
	BudgetAmount string `json:"budget_amount"`
	ActualAmount string `json:"actual_amount"`
	Variance     string `json:"variance"`         // actual - budget
	VariancePct  string `json:"variance_percent"` // variance / budget * 100
}

// VarianceReport is the full budget-vs-actual report for a period.
type VarianceReport struct {
	MerchantID string         `json:"merchant_id"`
	Period     string         `json:"period"`
	Lines      []VarianceLine `json:"lines"`
}
