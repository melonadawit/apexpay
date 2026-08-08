package risk

import "github.com/shopspring/decimal"

// Rule is a transaction-monitoring rule. Parameters are type-specific and decoded from
// the parameters jsonb.
type Rule struct {
	ID          string
	MerchantID  *string
	Name        string
	Description string
	RuleType    string // velocity_amount, velocity_count, threshold_amount, threshold_count, new_device, new_ip, high_ticket, high_failure_rate, manual
	Parameters  RuleParams
	Action      string // flag | review | block
	Severity    string // low | medium | high | critical
	Enabled     bool
}

// RuleParams holds decoded rule thresholds.
type RuleParams struct {
	WindowMinutes int             `json:"window_minutes"`
	AmountLimit   decimal.Decimal `json:"amount_limit"`
	CountLimit    int             `json:"count_limit"`
	RatioLimit    decimal.Decimal `json:"ratio_limit"` // e.g. failure-rate threshold 0..1
	Severity      string          `json:"severity"`
}

// EvaluationContext is what the engine uses to evaluate a transaction against rules.
type EvaluationContext struct {
	MerchantID string
	EntityType string // payment | payout | refund
	EntityID   string
	Amount     decimal.Decimal
	// Device/IP fields for new-device / new-ip rules.
	DeviceID   string
	IP         string
	CustomerID string
	// Aggregates provided by the repository.
	AmountInWindow decimal.Decimal // sum of amount in the window
	CountInWindow  int             // count of transactions in the window
	FailureRate    decimal.Decimal // 0..1 within the window
}

// Finding is the result of evaluating one rule against the context.
type Finding struct {
	RuleID   string
	RuleName string
	RuleType string
	Severity string
	Action   string
	Reason   string
	Matched  bool
}

// Evaluation is the aggregate result of running all rules.
type Evaluation struct {
	Findings []Finding
	// HighestSeverity returns the strongest severity among matched findings.
}

func (e Evaluation) Matched() []Finding {
	var out []Finding
	for _, f := range e.Findings {
		if f.Matched {
			out = append(out, f)
		}
	}
	return out
}

// HasBlock returns true if any matched finding is an action=block.
func (e Evaluation) HasBlock() bool {
	for _, f := range e.Findings {
		if f.Matched && f.Action == "block" {
			return true
		}
	}
	return false
}
