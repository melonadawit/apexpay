package link

import (
	"github.com/shopspring/decimal"
	"time"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusPaid      Status = "paid"
	StatusExpired   Status = "expired"
	StatusCancelled Status = "cancelled"
)

type PaymentLink struct {
	ID          string
	MerchantID  string
	PaymentID   *string
	Amount      decimal.Decimal
	Currency    string
	Description string
	Status      Status
	PublicToken string // unique for checkout URL
	ExpiresAt   *time.Time
	CreatedAt   time.Time
}

type CreateRequest struct {
	MerchantID  string
	Amount      decimal.Decimal
	Currency    string
	Description string
	ExpiresAt   *time.Time
}

type CheckoutSession struct {
	ID             string
	MerchantID     string
	PaymentID      string
	PaymentLinkID  *string
	PublicToken    string
	Status         string // open, completed, expired
	SelectedMethod string
	ExpiresAt      time.Time
}
