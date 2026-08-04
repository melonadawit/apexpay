package subscription

import (
	"time"

	"github.com/shopspring/decimal"
)

type Plan struct {
	ID            string
	MerchantID    string
	Name          string
	Description   string
	Amount        decimal.Decimal
	Currency      string
	IntervalType  string // day, week, month, year
	IntervalCount int
	TrialDays     int
	Status        string
	CreatedAt     time.Time
}

type SubscriptionStatus string
const (
	StatusIncomplete SubscriptionStatus = "incomplete"
	StatusTrialing   SubscriptionStatus = "trialing"
	StatusActive     SubscriptionStatus = "active"
	StatusPastDue    SubscriptionStatus = "past_due"
	StatusCanceled   SubscriptionStatus = "canceled"
	StatusPaused     SubscriptionStatus = "paused"
)

type Subscription struct {
	ID                  string
	MerchantID          string
	CustomerID          string
	PlanID              string
	Status              SubscriptionStatus
	CurrentPeriodStart  time.Time
	CurrentPeriodEnd    time.Time
	TrialEnd            *time.Time
	CancelAt            *time.Time
	CreatedAt           time.Time
}

type InvoiceStatus string
const (
	InvoiceDraft         InvoiceStatus = "draft"
	InvoiceOpen          InvoiceStatus = "open"
	InvoicePaid          InvoiceStatus = "paid"
	InvoiceUncollectible InvoiceStatus = "uncollectible"
)

type Invoice struct {
	ID             string
	MerchantID     string
	SubscriptionID string
	PaymentID      *string
	Amount         decimal.Decimal
	Currency       string
	Status         InvoiceStatus
	AttemptCount   int
	DueAt          time.Time
	CreatedAt      time.Time
}

type Customer struct {
	ID         string
	MerchantID string
	Email      string
	Phone      string
	Name       string
	FinHash    string // optional Fayda link
}
