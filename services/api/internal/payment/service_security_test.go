package payment

import (
	"context"
	"testing"
	"time"

	"apexpay/internal/ledger"
	"github.com/shopspring/decimal"
)

type securityTestRepo struct {
	verifiedMerchant string
	verifiedPayment  string
}

func (r *securityTestRepo) CreatePaymentTx(context.Context, *Payment, string, string) error {
	return nil
}
func (r *securityTestRepo) GetByTxRef(context.Context, string, string) (*Payment, error) {
	return nil, nil
}
func (r *securityTestRepo) UpdateStatusTx(context.Context, string, Status, *ledger.Journal, []ledger.Entry, *time.Time) error {
	return nil
}
func (r *securityTestRepo) ReserveIdempotency(context.Context, string, string, string) (*Payment, error) {
	return nil, nil
}
func (r *securityTestRepo) MarkConnectorStarted(context.Context, string, string, string) error {
	return nil
}
func (r *securityTestRepo) FailIdempotency(context.Context, string, string) error { return nil }
func (r *securityTestRepo) Mark2FAVerified(_ context.Context, merchantID, paymentID string) error {
	r.verifiedMerchant, r.verifiedPayment = merchantID, paymentID
	return nil
}

func (r *securityTestRepo) ListByMerchant(context.Context, string, int) ([]*Payment, error) {
	return []*Payment{}, nil
}

func (r *securityTestRepo) GetPaymentDetail(context.Context, string, string) (*PaymentDetail, error) {
	return nil, nil
}
func (r *securityTestRepo) DashboardSummary(context.Context, string) (*Summary, error) {
	return &Summary{}, nil
}

func TestCanonicalRequestHashNormalizesOnlyIntendedInputs(t *testing.T) {
	a := InitializeRequest{TxRef: "tx-1", Amount: decimal.RequireFromString("100.00"), Currency: "etb", Method: "CARD", CustomerEmail: "a@example.et", ReturnURL: "https://m.et/return", CallbackURL: "https://m.et/callback"}
	b := a
	b.Currency, b.Method = "ETB", "card"
	if canonicalRequestHash(a) != canonicalRequestHash(b) {
		t.Fatal("normalization should yield the same idempotency hash")
	}
	b.Amount = decimal.RequireFromString("100.01")
	if canonicalRequestHash(a) == canonicalRequestHash(b) {
		t.Fatal("a changed amount must change the idempotency hash")
	}
}

func TestVerify2FADemoOTPIsLocalOnlyAndPersists(t *testing.T) {
	repo := &securityTestRepo{}
	local := &Service{repo: repo, allowDemoOTP: true}
	if err := local.Verify2FA(context.Background(), "merchant-1", "payment-1", "123456"); err != nil {
		t.Fatalf("local demo challenge should persist: %v", err)
	}
	if repo.verifiedMerchant != "merchant-1" || repo.verifiedPayment != "payment-1" {
		t.Fatal("2FA verification was not persisted against the supplied tenant/payment")
	}
	production := &Service{repo: repo, allowDemoOTP: false}
	if err := production.Verify2FA(context.Background(), "merchant-1", "payment-1", "123456"); err == nil {
		t.Fatal("demo OTP must not work outside local mode")
	}
}
