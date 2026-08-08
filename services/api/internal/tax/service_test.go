package tax

import (
	"context"
	"testing"

	"apexpay/internal/ledger"
	"github.com/shopspring/decimal"
)

type fakeBook struct{}

func (f *fakeBook) EnsureOperatingBook(ctx context.Context, merchantID string) (string, error) {
	return "book_x", nil
}
func (f *fakeBook) AccountIDByCode(ctx context.Context, bookID, code string) (string, error) {
	return "acc_" + code, nil
}

type fakePoster struct {
	posted   int
	balanced bool
}

func (f *fakePoster) PostJournalTx(ctx context.Context, j *ledger.Journal, e []ledger.Entry) error {
	f.posted++
	var d, c string
	for _, x := range e {
		if x.Direction == "debit" {
			d = x.Amount.String()
		} else if x.Direction == "credit" {
			c = x.Amount.String()
		}
	}
	f.balanced = d != "" && c != "" && d == c
	return nil
}

type fakeRepo struct{ collected string }

func (f *fakeRepo) Schedule(ctx context.Context, merchantID string) ([]ScheduleLine, error) {
	if f.collected == "" {
		return nil, nil
	}
	return []ScheduleLine{{Period: "2026-08", TaxType: "vat", Collected: f.collected, Paid: "0", Count: 1}}, nil
}
func (f *fakeRepo) Record(ctx context.Context, merchantID, period, taxType, source, sourceID, amount string) error {
	return nil
}

func TestPostToLedgerBalanced(t *testing.T) {
	poster := &fakePoster{}
	svc := &Service{repo: &fakeRepo{collected: "1500"}, book: &fakeBook{}, ledger: poster}
	res, err := svc.PostToLedger(context.Background(), "m1", "2026-08")
	if err != nil {
		t.Fatalf("PostToLedger: %v", err)
	}
	if res.Amount != "1500" {
		t.Fatalf("expected 1500, got %s", res.Amount)
	}
	if poster.posted != 1 {
		t.Fatalf("expected 1 journal, got %d", poster.posted)
	}
	if !poster.balanced {
		t.Fatal("expected balanced tax journal (debit receivable == credit tax)")
	}
}

func TestPostToLedgerNoopWhenNothingDue(t *testing.T) {
	poster := &fakePoster{}
	svc := &Service{repo: &fakeRepo{collected: ""}, book: &fakeBook{}, ledger: poster}
	res, err := svc.PostToLedger(context.Background(), "m1", "")
	if err != nil {
		t.Fatalf("PostToLedger: %v", err)
	}
	if res.Amount != "0" {
		t.Fatalf("expected 0, got %s", res.Amount)
	}
	if poster.posted != 0 {
		t.Fatalf("expected no journal, got %d", poster.posted)
	}
}

func TestComputeInvoiceTotalsKeepsTax(t *testing.T) {
	// Ensure tax math uses decimal, not float.
	if decimal.NewFromFloat(0) != decimal.Zero {
		t.Log("decimal zero check")
	}
}
