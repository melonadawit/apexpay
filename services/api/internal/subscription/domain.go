package subscription

import (
	"time"

	"github.com/shopspring/decimal"
)

type Plan struct {
	ID            string          `json:"id"`
	MerchantID    string          `json:"merchant_id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	IntervalType  string          `json:"interval_type"` // day, week, month, year
	IntervalCount int             `json:"interval_count"`
	TrialDays     int             `json:"trial_days"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
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
	ID                 string             `json:"id"`
	MerchantID         string             `json:"merchant_id"`
	CustomerID         string             `json:"customer_id"`
	PlanID             string             `json:"plan_id"`
	Status             SubscriptionStatus `json:"status"`
	CurrentPeriodStart time.Time          `json:"current_period_start"`
	CurrentPeriodEnd   time.Time          `json:"current_period_end"`
	TrialEnd           *time.Time         `json:"trial_end"`
	CancelAt           *time.Time         `json:"cancel_at"`
	CreatedAt          time.Time          `json:"created_at"`
}

type InvoiceStatus string

const (
	InvoiceDraft         InvoiceStatus = "draft"
	InvoiceOpen          InvoiceStatus = "open"
	InvoicePaid          InvoiceStatus = "paid"
	InvoiceUncollectible InvoiceStatus = "uncollectible"
)

type Invoice struct {
	ID             string          `json:"id"`
	MerchantID     string          `json:"merchant_id"`
	SubscriptionID string          `json:"subscription_id"`
	PaymentID      *string         `json:"payment_id"`
	Amount         decimal.Decimal `json:"amount"`
	Currency       string          `json:"currency"`
	Status         InvoiceStatus   `json:"status"`
	AttemptCount   int             `json:"attempt_count"`
	DueAt          time.Time       `json:"due_at"`
	CreatedAt      time.Time       `json:"created_at"`
}

// SubscriptionDetail enriches a subscription with its plan, customer and invoices.
type SubscriptionDetail struct {
	Subscription *Subscription `json:"subscription"`
	Plan         *Plan         `json:"plan"`
	Customer     *Customer     `json:"customer"`
	Invoices     []Invoice     `json:"invoices"`
}

type Customer struct {
	ID         string `json:"id"`
	MerchantID string `json:"merchant_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Name       string `json:"name"`
	FinHash    string `json:"fin_hash"` // optional Fayda link
}
