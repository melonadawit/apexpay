package payment

import (
	"time"
	"github.com/shopspring/decimal"
)

type Status string
const (
	StatusCreated          Status = "created"
	StatusPending          Status = "pending"
	StatusProcessing       Status = "processing"
	StatusSucceeded        Status = "succeeded"
	StatusFailed           Status = "failed"
	StatusCanceled         Status = "canceled"
	StatusRefunded         Status = "refunded"
	StatusPartiallyRefunded Status = "partially_refunded"
)

type Payment struct {
	ID             string
	MerchantID     string
	TxRef          string
	Amount         decimal.Decimal
	Currency       string
	Status         Status
	Method         string
	Description    string
	CustomerEmail  string
	CustomerName   string
	CustomerPhone  string
	ConnectorID    string
	ConnectorRef   string
	RoutingRuleID  string
	CheckoutURL    string
	ReturnURL      string
	CallbackURL    string
	FeeAmount      decimal.Decimal
	NetAmount      decimal.Decimal
	Requires2FA    bool
	TwoFAVerified  bool
	SucceededAt    *time.Time
	CreatedAt      time.Time
}

type InitializeRequest struct {
	MerchantID   string
	TxRef        string
	Amount       decimal.Decimal
	Currency     string
	Method       string
	Description  string
	CustomerEmail string
	ReturnURL    string
	CallbackURL  string
	IdempotencyKey string
}

type VerifyRequest struct {
	MerchantID string
	TxRef      string
}
