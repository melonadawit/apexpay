package inventory

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

func TestPostCOGSBalanced(t *testing.T) {
	poster := &fakePoster{}
	svc := &Service{book: &fakeBook{}, ledger: poster}
	if err := svc.postCOGS(context.Background(), "m1", "order1", decimal.NewFromFloat(1250)); err != nil {
		t.Fatalf("postCOGS: %v", err)
	}
	if poster.posted != 1 {
		t.Fatalf("expected 1 journal, got %d", poster.posted)
	}
	if !poster.balanced {
		t.Fatal("expected balanced COGS journal (debit cost_of_sales == credit inventory)")
	}
}
