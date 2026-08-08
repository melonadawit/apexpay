package loyalty

// Loyalty & cashback domain.

type Tier struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	MinSpend        string `json:"min_spend"`
	CashbackPercent string `json:"cashback_percent"`
}

type Account struct {
	ID            string `json:"id"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerPhone string `json:"customer_phone,omitempty"`
	Points        string `json:"points"`
	TierID        string `json:"tier_id,omitempty"`
	TierName      string `json:"tier_name,omitempty"`
	TotalSpend    string `json:"total_spend"`
}

type CashbackTx struct {
	ID        string `json:"id"`
	PaymentID string `json:"payment_id,omitempty"`
	Amount    string `json:"amount"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}
