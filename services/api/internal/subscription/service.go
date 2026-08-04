package subscription

import (
	"context"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/platform/errors"

	"github.com/shopspring/decimal"
)

type Repository interface {
	CreatePlan(ctx context.Context, p *Plan) error
	GetPlan(ctx context.Context, merchantID, planID string) (*Plan, error)
	CreateCustomer(ctx context.Context, c *Customer) error
	CreateSubscription(ctx context.Context, s *Subscription) error
	CreateInvoice(ctx context.Context, inv *Invoice) error
	ListSubscriptions(ctx context.Context, merchantID string, status *SubscriptionStatus) ([]Subscription, error)
	UpdateSubscriptionStatus(ctx context.Context, id string, status SubscriptionStatus) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) CreatePlan(ctx context.Context, p *Plan) (*Plan, error) {
	if p.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.Validation("plan amount >0")
	}
	if p.IntervalCount <= 0 {
		p.IntervalCount = 1
	}
	p.ID = id.NewSubPlan()
	p.CreatedAt = time.Now()
	p.Status = "active"
	if err := s.repo.CreatePlan(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) CreateSubscription(ctx context.Context, merchantID, customerID, planID string) (*Subscription, error) {
	plan, err := s.repo.GetPlan(ctx, merchantID, planID)
	if err != nil {
		return nil, errors.NotFound("plan not found")
	}

	now := time.Now()
	sub := &Subscription{
		ID: merchantID + planID, // temp, will be overwritten
	}
	sub.ID = id.NewSubscription()
	sub.MerchantID = merchantID
	sub.CustomerID = customerID
	sub.PlanID = planID
	sub.CreatedAt = now

	// Trial handling
	if plan.TrialDays > 0 {
		sub.Status = StatusTrialing
		trialEnd := now.AddDate(0, 0, plan.TrialDays)
		sub.TrialEnd = &trialEnd
		sub.CurrentPeriodStart = now
		sub.CurrentPeriodEnd = trialEnd
	} else {
		sub.Status = StatusActive
		sub.CurrentPeriodStart = now
		sub.CurrentPeriodEnd = addInterval(now, plan.IntervalType, plan.IntervalCount)
	}

	if err := s.repo.CreateSubscription(ctx, sub); err != nil {
		return nil, err
	}

	// Create first invoice (draft if trialing)
	invStatus := InvoiceOpen
	if sub.Status == StatusTrialing {
		invStatus = InvoiceDraft
	}
	inv := &Invoice{
		ID: id.New("sinv"), MerchantID: merchantID, SubscriptionID: sub.ID,
		Amount: plan.Amount, Currency: plan.Currency, Status: invStatus,
		DueAt: sub.CurrentPeriodEnd, CreatedAt: now,
	}
	_ = s.repo.CreateInvoice(ctx, inv)

	return sub, nil
}

func addInterval(t time.Time, typ string, count int) time.Time {
	switch typ {
	case "day":
		return t.AddDate(0, 0, count)
	case "week":
		return t.AddDate(0, 0, 7*count)
	case "month":
		return t.AddDate(0, count, 0)
	case "year":
		return t.AddDate(count, 0, 0)
	default:
		return t.AddDate(0, count, 0)
	}
}

// Dunning schedule: 1d, 3d, 5d - optimal next attempt calculation
func NextDunningAttempt(attemptCount int, lastAttempt time.Time) time.Time {
	switch attemptCount {
	case 0:
		return lastAttempt.Add(24 * time.Hour)
	case 1:
		return lastAttempt.Add(72 * time.Hour)
	case 2:
		return lastAttempt.Add(120 * time.Hour)
	default:
		return lastAttempt // max attempts reached
	}
}
