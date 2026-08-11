package payout

import (
	"time"

	"github.com/shopspring/decimal"
)

type PayoutStatus string

const (
	StatusCreated         PayoutStatus = "created"
	StatusPendingApproval PayoutStatus = "pending_approval"
	StatusQueued          PayoutStatus = "queued"
	StatusProcessing      PayoutStatus = "processing"
	StatusSucceeded       PayoutStatus = "succeeded"
	StatusFailed          PayoutStatus = "failed"
	StatusReturned        PayoutStatus = "returned"
)

type Beneficiary struct {
	ID                 string    `json:"id"`
	MerchantID         string    `json:"merchant_id"`
	Name               string    `json:"name"`
	AccountNoMasked    string    `json:"account_no_masked"`
	AccountNoHash      string    `json:"account_no_hash"`
	BankCode           string    `json:"bank_code"`
	BankName           string    `json:"bank_name"`
	Type               string    `json:"type"` // individual, business
	VerificationStatus string    `json:"verification_status"`
	CreatedAt          time.Time `json:"created_at"`
}

type PayoutBatch struct {
	ID         string          `json:"id"`
	MerchantID string          `json:"merchant_id"`
	BookID     *string         `json:"book_id"`
	BatchRef   string          `json:"batch_ref"`
	Amount     decimal.Decimal `json:"amount"`
	Currency   string          `json:"currency"`
	Status     string          `json:"status"`
	ApprovedBy *string         `json:"approved_by"`
	CreatedAt  time.Time       `json:"created_at"`
	TotalCount int             `json:"total_count"`
	Payouts    []Payout        `json:"payouts,omitempty"`
}

type Payout struct {
	ID            string          `json:"id"`
	MerchantID    string          `json:"merchant_id"`
	BatchID       *string         `json:"batch_id"`
	BeneficiaryID string          `json:"beneficiary_id"`
	PayoutRef     string          `json:"payout_ref"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	Status        PayoutStatus    `json:"status"`
	Method        string          `json:"method"`
	ConnectorID   string          `json:"connector_id"`
	ConnectorRef  string          `json:"connector_ref"`
	FailureCode   string          `json:"failure_code"`
	CreatedAt     time.Time       `json:"created_at"`
}

type CreateBulkRequest struct {
	MerchantID string
	BatchRef   string
	Currency   string
	Items      []struct {
		BeneficiaryID string
		Amount        decimal.Decimal
		PayoutRef     string
	}
}
