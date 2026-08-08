package dispute

// Disputes & chargebacks domain.

type Dispute struct {
	ID         string         `json:"id"`
	PaymentID  string         `json:"payment_id,omitempty"`
	Amount     string         `json:"amount"`
	Currency   string         `json:"currency"`
	ReasonCode string         `json:"reason_code"`
	Status     string         `json:"status"`
	Evidence   []EvidenceItem `json:"evidence"`
	Resolution string         `json:"resolution,omitempty"`
	Fee        string         `json:"fee"`
	CreatedAt  string         `json:"created_at"`
}

type EvidenceItem struct {
	FileKey     string `json:"file_key"`
	Description string `json:"description,omitempty"`
}
