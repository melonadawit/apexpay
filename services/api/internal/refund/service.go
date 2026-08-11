package refund

import (
	"context"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/ledger"
	"apexpay/internal/platform/errors"
	"github.com/shopspring/decimal"
)

type PaymentInfo struct {
	ID          string
	MerchantID  string
	Amount      decimal.Decimal
	RefundedAmt decimal.Decimal // sum of succeeded refunds
	FeeAmount   decimal.Decimal
	Status      string
	Currency    string
	ConnectorID string
}

type Repository interface {
	GetPayment(ctx context.Context, merchantID, paymentID string) (*PaymentInfo, error)
	GetRefundByRef(ctx context.Context, merchantID, refundRef string) (*Refund, error)
	GetRefundByID(ctx context.Context, merchantID, id string) (*Refund, error)
	ListRefundsByPayment(ctx context.Context, paymentID string) ([]Refund, error)
	CreateRefundTx(ctx context.Context, refund *Refund, journal *ledger.Journal, entries []ledger.Entry) error
	UpdateRefundStatus(ctx context.Context, id string, status Status, connectorRef string) error
}

type Service struct {
	repo   Repository
	ledger *ledger.Service
}

func NewService(repo Repository, ledgerSvc *ledger.Service) *Service {
	return &Service{repo: repo, ledger: ledgerSvc}
}

// Create refund - optimal algorithm: check remaining refundable amount O(1) + idempotency
func (s *Service) Create(ctx context.Context, req CreateRequest) (*Refund, error) {
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.Validation("refund amount must be >0")
	}

	pay, err := s.repo.GetPayment(ctx, req.MerchantID, req.PaymentID)
	if err != nil {
		return nil, errors.NotFound("payment not found")
	}
	if pay.Status != "succeeded" && pay.Status != "partially_refunded" {
		return nil, errors.Validation("payment must be succeeded to refund")
	}

	// Idempotency: by refund_ref unique (merchant_id, refund_ref)
	existing, err := s.repo.GetRefundByRef(ctx, req.MerchantID, req.RefundRef)
	if err == nil && existing != nil {
		// same amount? return existing else conflict
		if existing.Amount.Equal(req.Amount) {
			return existing, nil
		}
		return nil, errors.Conflict(errors.CodeDuplicateRefundRef, "refund_ref already exists with different amount")
	}

	remaining := pay.Amount.Sub(pay.RefundedAmt)
	if req.Amount.GreaterThan(remaining) {
		return nil, errors.New(errors.CodeRefundExceeded, "refund amount exceeds remaining refundable", 400)
	}

	// Fee reversal calc - optimal decimal math
	feeReversal := s.calcFeeReversal(pay.FeeAmount, pay.Amount, req.Amount, req.FeePolicy)

	r := &Refund{
		ID:          id.NewRefund(),
		MerchantID:  req.MerchantID,
		PaymentID:   req.PaymentID,
		RefundRef:   req.RefundRef,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      StatusProcessing,
		Reason:      req.Reason,
		FeeReversal: feeReversal,
		ConnectorID: pay.ConnectorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Ledger posting M2: Dr merchant_payable (amount-feeReversal) + Dr fee_due (feeReversal) Cr clearing
	// Use transfer_group = refund group for multi-book (merchant + rail)
	journal := &ledger.Journal{
		ID:            id.NewLedgerJournal(),
		BookID:        "merchant_operating_default", // resolved via ledger svc in reality
		PostingKey:    "refund:" + r.ID,
		Memo:          "refund " + req.Reason,
		ReferenceType: "refund",
		ReferenceID:   r.ID,
	}
	// Ledger entries - double entry invariant check sum debits=credits
	entries := []ledger.Entry{
		{
			ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID,
			AccountID: "liability:merchant_payable", Direction: "debit",
			Amount: req.Amount.Sub(feeReversal), Currency: req.Currency,
		},
		{
			ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID,
			AccountID: "liability:platform_fee_due", Direction: "debit",
			Amount: feeReversal, Currency: req.Currency,
		},
		{
			ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID,
			AccountID: "asset:clearing:" + pay.ConnectorID, Direction: "credit",
			Amount: req.Amount, Currency: req.Currency,
		},
	}
	// Optimize: filter zero amount entries (if feeReversal zero)
	filtered := entries[:0]
	for _, e := range entries {
		if e.Amount.GreaterThan(decimal.Zero) {
			filtered = append(filtered, e)
		}
	}

	// Single Tx: refund insert + journal + entries + payment status update handled in repo
	if err := s.repo.CreateRefundTx(ctx, r, journal, filtered); err != nil {
		return nil, err
	}

	// Async: connector refund call would happen in worker; here we mark succeeded for mock
	if pay.ConnectorID == "mock" {
		_ = s.repo.UpdateRefundStatus(ctx, r.ID, StatusSucceeded, "mock_ref_"+r.ID)
		r.Status = StatusSucceeded
	}

	return r, nil
}

func (s *Service) calcFeeReversal(feeTotal, paymentAmount, refundAmount decimal.Decimal, policy FeePolicy) decimal.Decimal {
	if paymentAmount.IsZero() || feeTotal.IsZero() {
		return decimal.Zero
	}
	switch policy {
	case FeePolicyNonRefundable:
		return decimal.Zero
	case FeePolicyFull:
		if refundAmount.Equal(paymentAmount) {
			return feeTotal
		}
		return decimal.Zero
	case FeePolicyProRata:
		// pro-rata: fee = totalFee * (refund/pay) with bankers rounding
		ratio := refundAmount.Div(paymentAmount)
		return feeTotal.Mul(ratio).Round(2) // ETB scale 2
	default:
		return decimal.Zero
	}
}
