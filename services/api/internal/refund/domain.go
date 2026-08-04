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
	ID            string
	MerchantID    string
	PaymentID     string
	RefundRef     string // merchant supplied unique
	Amount        decimal.Decimal
	Currency      string
	Status        Status
	Reason        string
	FeeReversal   decimal.Decimal // policy: how much of platform fee reversed
	ConnectorID   string
	ConnectorRef  string
	FailureCode   string
	FailureMsg    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Policy: Fee reversal strategy per SAD / DATABASE comment.
type FeePolicy string
const (
	FeePolicyNonRefundable FeePolicy = "non_refundable" // platform keeps fee
	FeePolicyProRata       FeePolicy = "pro_rata"       // reverse pro-rata
	FeePolicyFull          FeePolicy = "full"           // full fee reversal on full refund
)

type CreateRequest struct {
	MerchantID  string
	PaymentID   string
	RefundRef   string
	Amount      decimal.Decimal
	Currency    string
	Reason      string
	FeePolicy   FeePolicy
	IdempotencyKey string
}
