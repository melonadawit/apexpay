package fxreval

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

func TestPostToLedgerBalancedForGain(t *testing.T) {
	poster := &fakePoster{}
	svc := &Service{book: &fakeBook{}, ledger: poster}
	// Net gain: debit bank, credit fx_gain.
	if err := svc.postToLedger(context.Background(), "m1", "2026-08", decimal.NewFromFloat(75000)); err != nil {
		t.Fatalf("postToLedger: %v", err)
	}
	if poster.posted != 1 {
		t.Fatalf("expected 1 journal, got %d", poster.posted)
	}
	if !poster.balanced {
		t.Fatal("expected balanced FX journal")
	}
}

func TestPostToLedgerBalancedForLoss(t *testing.T) {
	poster := &fakePoster{}
	svc := &Service{book: &fakeBook{}, ledger: poster}
	// Net loss: debit fx_loss, credit bank (magnitude positive).
	if err := svc.postToLedger(context.Background(), "m1", "2026-08", decimal.NewFromFloat(-5000)); err != nil {
		t.Fatalf("postToLedger: %v", err)
	}
	if poster.posted != 1 {
		t.Fatalf("expected 1 journal, got %d", poster.posted)
	}
	if !poster.balanced {
		t.Fatal("expected balanced FX journal for loss")
	}
}
