package link

import (
	"context"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/platform/errors"
	"github.com/shopspring/decimal"
)

type Repository interface {
	CreateLink(ctx context.Context, pl *PaymentLink) error
	GetByToken(ctx context.Context, token string) (*PaymentLink, error)
	ListByMerchant(ctx context.Context, merchantID string) ([]PaymentLink, error)
	MarkPaid(ctx context.Context, linkID, paymentID string) error
	CreateCheckoutSession(ctx context.Context, cs *CheckoutSession) error
}

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, req CreateRequest) (*PaymentLink, *CheckoutSession, error) {
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, nil, errors.Validation("amount must be >0")
	}
	publicToken := id.New("tok") // outstanding unique for checkout URL
	linkID := id.New("pl")
	pl := &PaymentLink{
		ID: linkID, MerchantID: req.MerchantID, Amount: req.Amount, Currency: req.Currency,
		Description: req.Description, Status: StatusActive, PublicToken: publicToken,
		ExpiresAt: req.ExpiresAt, CreatedAt: time.Now(),
	}
	if err := s.repo.CreateLink(ctx, pl); err != nil {
		return nil, nil, err
	}

	// Create checkout session for this link (will be linked to payment later on initialize)
	cs := &CheckoutSession{
		ID: id.New("cs"), MerchantID: req.MerchantID, PaymentID: linkID, // placeholder payment ID = link ID initially, real payment created on init
		PaymentLinkID: &linkID, PublicToken: publicToken, Status: "open",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	// In real flow, checkout session created when payment initialized, not here. For skeleton we skip creation error check.
	_ = s.repo.CreateCheckoutSession(ctx, cs)

	return pl, cs, nil
}

func (s *Service) List(ctx context.Context, merchantID string) ([]PaymentLink, error) {
	return s.repo.ListByMerchant(ctx, merchantID)
}

func (s *Service) GetByToken(ctx context.Context, token string) (*PaymentLink, error) {
	pl, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return nil, errors.NotFound("payment link not found")
	}
	if pl.Status != StatusActive {
		return nil, errors.Validation("link not active, already paid/expired")
	}
	return pl, nil
}
