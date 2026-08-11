package refund

import (
	"time"

	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusCreated    Status = "created"
	StatusProcessing Status = "processing"
	StatusSucceeded  Status = "succeeded"
	StatusFailed     Status = "failed"
)

type Refund struct {
	ID           string          `json:"id"`
	MerchantID   string          `json:"merchant_id"`
	PaymentID    string          `json:"payment_id"`
	RefundRef    string          `json:"refund_ref"`
	Amount       decimal.Decimal `json:"amount"`
	Currency     string          `json:"currency"`
	Status       Status          `json:"status"`
	Reason       string          `json:"reason"`
	FeeReversal  decimal.Decimal `json:"fee_reversal"`
	ConnectorID  string          `json:"connector_id"`
	ConnectorRef string          `json:"connector_ref"`
	FailureCode  string          `json:"failure_code"`
	FailureMsg   string          `json:"failure_message"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// Policy: Fee reversal strategy per SAD / DATABASE comment.
type FeePolicy string

const (
	FeePolicyNonRefundable FeePolicy = "non_refundable" // platform keeps fee
	FeePolicyProRata       FeePolicy = "pro_rata"       // reverse pro-rata
	FeePolicyFull          FeePolicy = "full"           // full fee reversal on full refund
)

type CreateRequest struct {
	MerchantID     string
	PaymentID      string
	RefundRef      string
	Amount         decimal.Decimal
	Currency       string
	Reason         string
	FeePolicy      FeePolicy
	IdempotencyKey string
}
