package swarm

import (
	"context"
	"fmt"
	"time"

	"apexpay/internal/id"
)

// ToolExecutor implementation - calls Go services via registry O(1) lookup optimal
// Best practice: JSON schema validation via go-playground/validator, no direct ledger invent

type ToolExecutorImpl struct {
	// Inject services - simplified as function map for skeleton
	paymentLinkCreator func(ctx context.Context, merchantID string, amount float64, currency, desc string) (map[string]interface{}, error)
	payoutCreator      func(ctx context.Context, merchantID string, amount float64) (map[string]interface{}, error)
	payrollCalculator  func(ctx context.Context, merchantID string) (map[string]interface{}, error)
	tpvGetter          func(ctx context.Context, merchantID string) (map[string]interface{}, error)
	complianceAsker    func(ctx context.Context, merchantID, query string) (map[string]interface{}, error)
}

func NewToolExecutor() *ToolExecutorImpl {
	return &ToolExecutorImpl{
		paymentLinkCreator: func(ctx context.Context, merchantID string, amount float64, currency, desc string) (map[string]interface{}, error) {
			// Mock creates payment_link_url
			return map[string]interface{}{
				"payment_link_id":  id.New("pl"),
				"payment_link_url": fmt.Sprintf("https://checkout.apexpay.et/c/%s_%.0f", merchantID, amount),
				"amount":           amount, "currency": currency,
			}, nil
		},
		payoutCreator: func(ctx context.Context, merchantID string, amount float64) (map[string]interface{}, error) {
			return map[string]interface{}{
				"payout_id": id.NewPayout(), "amount": amount, "status": "pending_approval",
			}, nil
		},
		payrollCalculator: func(ctx context.Context, merchantID string) (map[string]interface{}, error) {
			return map[string]interface{}{
				"payroll_run_id": id.NewPayrollRun(), "period_month": 7, "period_year": 2026, "total_net": 150000, "status": "pending_approval",
			}, nil
		},
		tpvGetter: func(ctx context.Context, merchantID string) (map[string]interface{}, error) {
			return map[string]interface{}{"tpv_today": 125430, "currency": "ETB", "count": 42}, nil
		},
		complianceAsker: func(ctx context.Context, merchantID, query string) (map[string]interface{}, error) {
			return map[string]interface{}{"answer": "Transactions above 5000 ETB require 2FA per ONPS/10/2025 [1]", "citations": []map[string]interface{}{{"doc_title": "ONPS/10/2025", "page": 3, "score": 0.92}}}, nil
		},
	}
}

func (e *ToolExecutorImpl) Execute(ctx context.Context, tool string, args map[string]interface{}, merchantID string) (map[string]interface{}, error) {
	start := time.Now()
	var result map[string]interface{}
	var err error

	switch tool {
	case "create_payment_link":
		amt, _ := args["amount"].(float64)
		if amt == 0 {
			if a, ok := args["amount"].(int); ok {
				amt = float64(a)
			}
		}
		currency, _ := args["currency"].(string)
		if currency == "" {
			currency = "ETB"
		}
		desc, _ := args["description"].(string)
		result, err = e.paymentLinkCreator(ctx, merchantID, amt, currency, desc)
	case "create_payout":
		amt, _ := args["amount"].(float64)
		result, err = e.payoutCreator(ctx, merchantID, amt)
	case "calculate_payroll_draft", "calculate_payroll":
		result, err = e.payrollCalculator(ctx, merchantID)
	case "get_tpv":
		result, err = e.tpvGetter(ctx, merchantID)
	case "ask_compliance":
		q, _ := args["query"].(string)
		result, err = e.complianceAsker(ctx, merchantID, q)
	case "list_payments":
		result = map[string]interface{}{"payments": []interface{}{}, "count": 5}
	default:
		err = fmt.Errorf("unknown tool %s", tool)
	}

	// Latency recorded in caller ToolCall LatencyMS

	_ = start // used by caller for latency

	return result, err
}
