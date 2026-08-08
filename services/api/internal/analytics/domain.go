package analytics

// Analytics & cohort domain types.

type Daily struct {
	StatDate        string                 `json:"stat_date"`
	Revenue         string                 `json:"revenue"`
	TPV             string                 `json:"tpv"`
	PaymentCount    int                    `json:"payment_count"`
	SuccessCount    int                    `json:"success_count"`
	FailedCount     int                    `json:"failed_count"`
	RefundAmount    string                 `json:"refund_amount"`
	MethodBreakdown map[string]interface{} `json:"method_breakdown"`
}

type Cohort struct {
	CohortMonth     string `json:"cohort_month"`
	Customers       int    `json:"customers"`
	Month1Retention string `json:"month1_retention"`
	Month2Retention string `json:"month2_retention"`
	Month3Retention string `json:"month3_retention"`
	MRR             string `json:"mrr"`
}
