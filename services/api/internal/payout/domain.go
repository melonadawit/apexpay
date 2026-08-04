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
	ID                 string
	MerchantID         string
	Name               string
	AccountNoMasked    string
	AccountNoHash      string
	BankCode           string
	BankName           string
	Type               string // individual, business
	VerificationStatus string
	CreatedAt          time.Time
}

type PayoutBatch struct {
	ID         string
	MerchantID string
	BookID     *string
	BatchRef   string
	Amount     decimal.Decimal
	Currency   string
	Status     string
	ApprovedBy *string
	CreatedAt  time.Time
	Payouts    []Payout
}

type Payout struct {
	ID            string
	MerchantID    string
	BatchID       *string
	BeneficiaryID string
	PayoutRef     string
	Amount        decimal.Decimal
	Currency      string
	Status        PayoutStatus
	Method        string
	ConnectorID   string
	ConnectorRef  string
	FailureCode   string
	CreatedAt     time.Time
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
