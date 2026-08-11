package payment

import (
	"github.com/shopspring/decimal"
	"time"
)

type Status string

const (
	StatusCreated           Status = "created"
	StatusPending           Status = "pending"
	StatusProcessing        Status = "processing"
	StatusSucceeded         Status = "succeeded"
	StatusFailed            Status = "failed"
	StatusCanceled          Status = "canceled"
	StatusRefunded          Status = "refunded"
	StatusPartiallyRefunded Status = "partially_refunded"
)

type Payment struct {
	ID            string          `json:"id"`
	MerchantID    string          `json:"merchant_id"`
	TxRef         string          `json:"tx_ref"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	Status        Status          `json:"status"`
	Method        string          `json:"method"`
	Description   string          `json:"description"`
	CustomerEmail string          `json:"customer_email"`
	CustomerName  string          `json:"customer_name"`
	CustomerPhone string          `json:"customer_phone"`
	ConnectorID   string          `json:"connector_id"`
	ConnectorRef  string          `json:"connector_ref"`
	RoutingRuleID string          `json:"routing_rule_id"`
	CheckoutURL   string          `json:"checkout_url"`
	ReturnURL     string          `json:"return_url"`
	CallbackURL   string          `json:"callback_url"`
	FeeAmount     decimal.Decimal `json:"fee_amount"`
	NetAmount     decimal.Decimal `json:"net_amount"`
	Requires2FA   bool            `json:"requires_2fa"`
	TwoFAVerified bool            `json:"two_fa_verified"`
	SucceededAt   *time.Time      `json:"succeeded_at"`
	CreatedAt     time.Time       `json:"created_at"`
}

type InitializeRequest struct {
	MerchantID     string
	TxRef          string
	Amount         decimal.Decimal
	Currency       string
	Method         string
	Description    string
	CustomerEmail  string
	ReturnURL      string
	CallbackURL    string
	IdempotencyKey string
}

type VerifyRequest struct {
	MerchantID string
	TxRef      string
}

// PaymentDetail is the full NBE exam-console view for one transaction: the
// payment row plus its lifecycle ledger journals and entries.
type PaymentDetail struct {
	Payment  *Payment         `json:"payment"`
	Journals []PaymentJournal `json:"journals"`
}

type PaymentJournal struct {
	ID         string         `json:"id"`
	BookID     string         `json:"book_id"`
	PostingKey string         `json:"posting_key"`
	Memo       string         `json:"memo"`
	CreatedAt  time.Time      `json:"created_at"`
	Entries    []PaymentEntry `json:"entries"`
}

type PaymentEntry struct {
	Direction   string `json:"direction"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	AccountCode string `json:"account_code"`
	AccountName string `json:"account_name"`
}
