package payment

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"apexpay/internal/connector"
	"apexpay/internal/id"
	"apexpay/internal/ledger"
	"apexpay/internal/platform/errors"
	"apexpay/internal/routing"
	"github.com/shopspring/decimal"
)

type Repository interface {
	CreatePaymentTx(ctx context.Context, p *Payment, outboxEventID, idempotencyKey string) error
	GetByTxRef(ctx context.Context, merchantID, txRef string) (*Payment, error)
	UpdateStatusTx(ctx context.Context, paymentID string, status Status, journal *ledger.Journal, entries []ledger.Entry, succeededAt *time.Time) error
	ReserveIdempotency(ctx context.Context, merchantID, key, requestHash string) (*Payment, error)
	MarkConnectorStarted(ctx context.Context, merchantID, key, txRef string) error
	FailIdempotency(ctx context.Context, merchantID, key string) error
	Mark2FAVerified(ctx context.Context, merchantID, paymentID string) error
}

type Service struct {
	repo     Repository
	ledger   *ledger.Service
	router   *routing.Service
	registry map[string]connector.Connector // connector_id -> Connector optimal O(1)
	mdrRate  decimal.Decimal                // 2.9%
	allowDemoOTP bool // strictly local-only until a real challenge provider is wired
}

func NewService(repo Repository, ledgerSvc *ledger.Service, router *routing.Service, registry map[string]connector.Connector, mdrRate decimal.Decimal, allowDemoOTP bool) *Service {
	return &Service{repo: repo, ledger: ledgerSvc, router: router, registry: registry, mdrRate: mdrRate, allowDemoOTP: allowDemoOTP}
}

func (s *Service) Initialize(ctx context.Context, req InitializeRequest) (*Payment, error) {
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.Validation("amount must be >0")
	}
	// Idempotency keys are scoped to a merchant and bound to a canonical request.
	// Reusing a key with different business inputs is a conflict, never a new charge.
	requestHash := canonicalRequestHash(req)
	if req.IdempotencyKey != "" {
		existing, err := s.repo.ReserveIdempotency(ctx, req.MerchantID, req.IdempotencyKey, requestHash)
		if err != nil {
			if err == ErrIdempotencyConflict || err == ErrIdempotencyInProgress { return nil, errors.Conflict(errors.CodeConflict, err.Error()) }
			return nil, err
		}
		if existing != nil { return existing, nil }
	}

	// Duplicate tx_ref check handled by DB unique (merchant_id, tx_ref)
	// Routing evaluation
	decision, err := s.router.Evaluate(ctx, req.MerchantID, req.Amount, req.Currency, req.Method)
	if err != nil {
		decision = &routing.RoutingDecision{Chosen: routing.ConnectorMock, Primary: routing.ConnectorMock, Reason: "router error fallback mock"}
	}

	connID := string(decision.Chosen)
	conn, ok := s.registry[connID]
	if !ok {
		conn = s.registry["mock"]
		connID = "mock"
	}

	// Fee calc: fee = amount * mdrRate rounded ETB scale 2
	fee := req.Amount.Mul(s.mdrRate).Round(2)
	net := req.Amount.Sub(fee)

	// 2FA check per NBE ONPS/10/2025 >5000 ETB
	requires2FA := false
	if req.Currency == "ETB" && req.Amount.GreaterThan(decimal.NewFromInt(5000)) {
		requires2FA = true
	}

	// Persist this boundary before the external call. A crash after this point is
	// operationally ambiguous and must be reconciled, never blindly retried.
	if req.IdempotencyKey != "" {
		if err := s.repo.MarkConnectorStarted(ctx, req.MerchantID, req.IdempotencyKey, req.TxRef); err != nil { return nil, err }
	}
	// Connector Initialize
	initResp, err := conn.(connector.Connector).Initialize(ctx, connector.InitializeRequest{
		MerchantID: req.MerchantID, Amount: req.Amount.String(), Currency: req.Currency, TxRef: req.TxRef, ReturnURL: req.ReturnURL,
	})
	if err != nil {
		// circuit breaker record failure
		s.router.RecordFailure(routing.ConnectorID(connID))
		if req.IdempotencyKey != "" { _ = s.repo.FailIdempotency(ctx, req.MerchantID, req.IdempotencyKey) }
		return nil, errors.New(errors.CodeConnectorDown, fmt.Sprintf("connector %s failed: %v", connID, err), 502)
	}
	s.router.RecordSuccess(routing.ConnectorID(connID))

	now := time.Now()
	p := &Payment{
		ID: id.NewPayment(), MerchantID: req.MerchantID, TxRef: req.TxRef,
		Amount: req.Amount, Currency: req.Currency, Status: StatusPending,
		Method: req.Method, Description: req.Description, CustomerEmail: req.CustomerEmail,
		ConnectorID: connID, ConnectorRef: initResp.ConnectorRef,
		RoutingRuleID: decision.RuleID, CheckoutURL: initResp.CheckoutURL,
		ReturnURL: req.ReturnURL, CallbackURL: req.CallbackURL,
		FeeAmount: fee, NetAmount: net, Requires2FA: requires2FA, CreatedAt: now,
	}

	// Create with outbox event payment.created
	outboxID := id.NewOutbox()
	if err := s.repo.CreatePaymentTx(ctx, p, outboxID, req.IdempotencyKey); err != nil {
		// handle duplicate tx_ref conflict
		if isDuplicate(err) {
			return nil, errors.Conflict(errors.CodeDuplicateTxRef, "duplicate tx_ref")
		}
		if req.IdempotencyKey != "" { _ = s.repo.FailIdempotency(ctx, req.MerchantID, req.IdempotencyKey) }
		return nil, err
	}


	return p, nil
}


func canonicalRequestHash(req InitializeRequest) string {
	// Delimit and normalize all fields that change the business operation.
	// Do not include presentation-only whitespace or the idempotency key itself.
	canonical := strings.Join([]string{req.TxRef, req.Amount.String(), strings.ToUpper(req.Currency), strings.ToLower(req.Method), req.CustomerEmail, req.ReturnURL, req.CallbackURL}, "\x1f")
	digest := sha256.Sum256([]byte(canonical))
	return hex.EncodeToString(digest[:])
}

// Verify2FA persists authorization before a payment can be verified/captured. A real
// challenge provider is mandatory outside local development; the legacy demo OTP is
// deliberately unavailable in staging/production.
func (s *Service) Verify2FA(ctx context.Context, merchantID, paymentID, otp string) error {
	if !s.allowDemoOTP || otp != "123456" {
		return errors.New(errors.CodeValidation, "invalid or unavailable 2FA challenge", 400)
	}
	return s.repo.Mark2FAVerified(ctx, merchantID, paymentID)
}

func (s *Service) Verify(ctx context.Context, req VerifyRequest) (*Payment, error) {
	p, err := s.repo.GetByTxRef(ctx, req.MerchantID, req.TxRef)
	if err != nil {
		return nil, errors.NotFound("payment not found")
	}
	if p.Status == StatusSucceeded {
		return p, nil // idempotent no-op per MVP B6
	}
	// If requires 2FA and not verified, return pending
	if p.Requires2FA && !p.TwoFAVerified {
		return p, nil
	}

	// Call connector Verify
	conn, ok := s.registry[p.ConnectorID]
	if !ok {
		conn = s.registry["mock"]
	}
	verifyResp, err := conn.(connector.Connector).Verify(ctx, connector.VerifyRequest{ConnectorRef: p.ConnectorRef, TxRef: p.TxRef})
	if err != nil {
		return p, nil // keep pending, worker will retry
	}
	if verifyResp.Status != "succeeded" {
		return p, nil
	}

	// Ledger M1 posting atomically with status update per DATABASE transaction boundary
	now := time.Now()
	journal := &ledger.Journal{
		ID: id.NewLedgerJournal(), BookID: "merchant_operating_default", // resolved via ledger svc in prod
		PostingKey: fmt.Sprintf("payment_success:%s", p.ID),
		Memo:       "payment success", ReferenceType: "payment", ReferenceID: p.ID,
		TransferGroup: fmt.Sprintf("pay_%s", p.ID),
	}
	entries := []ledger.Entry{
		{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "asset:clearing:" + p.ConnectorID, Direction: "debit", Amount: p.Amount, Currency: p.Currency},
		{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "liability:merchant_payable", Direction: "credit", Amount: p.NetAmount, Currency: p.Currency},
		{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "liability:platform_fee_due", Direction: "credit", Amount: p.FeeAmount, Currency: p.Currency},
	}
	// Filter zero fee if any
	filtered := entries[:0]
	for _, e := range entries {
		if e.Amount.GreaterThan(decimal.Zero) {
			filtered = append(filtered, e)
		}
	}

	if err := s.repo.UpdateStatusTx(ctx, p.ID, StatusSucceeded, journal, filtered, &now); err != nil {
		return nil, err
	}
	p.Status = StatusSucceeded
	p.SucceededAt = &now
	return p, nil
}

func isDuplicate(err error) bool {
	// check pg unique violation code 23505
	return err != nil && (contains(err.Error(), "duplicate") || contains(err.Error(), "unique"))
}
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	})()
}
