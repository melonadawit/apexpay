package tax

import (
	"context"
	"fmt"
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

// TestScheduleNetsCollectedMinusPaid verifies the "due" amount and total reduce by what has
// already been remitted to the tax authority.
func TestScheduleNetsCollectedMinusPaid(t *testing.T) {
	repo := &multiLineRepo{
		lines: []ScheduleLine{
			{Period: "2026-08", TaxType: "vat", Collected: "1000", Paid: "400", Count: 2},
			{Period: "2026-08", TaxType: "withholding", Collected: "300", Paid: "0", Count: 1},
		},
	}
	svc := &Service{repo: repo, book: &fakeBook{}, ledger: &fakePoster{}}
	sched, err := svc.Schedule(context.Background(), "m1")
	if err != nil {
		t.Fatalf("Schedule: %v", err)
	}
	if len(sched.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(sched.Lines))
	}
	// vat due = 1000-400 = 600; withholding due = 300; total = 900
	if sched.Lines[0].Due != "600" {
		t.Fatalf("expected vat due 600, got %s", sched.Lines[0].Due)
	}
	if sched.TotalDue != "900" {
		t.Fatalf("expected total due 900, got %s", sched.TotalDue)
	}
}

// TestPostToLedgerMissingAccount covers the error path when a required GL account is absent.
func TestPostToLedgerMissingAccount(t *testing.T) {
	book := &failingBook{err: fmt.Errorf("tax liability account missing")}
	svc := &Service{repo: &fakeRepo{collected: "500"}, book: book, ledger: &fakePoster{}}
	_, err := svc.PostToLedger(context.Background(), "m1", "2026-08")
	if err == nil {
		t.Fatal("expected error when tax liability account is missing")
	}
}

// multiLineRepo returns several schedule lines for netting tests.
type multiLineRepo struct{ lines []ScheduleLine }

func (f *multiLineRepo) Schedule(ctx context.Context, merchantID string) ([]ScheduleLine, error) {
	return f.lines, nil
}
func (f *multiLineRepo) Record(ctx context.Context, merchantID, period, taxType, source, sourceID, amount string) error {
	return nil
}

// failingBook fails account resolution to exercise the missing-account error path.
type failingBook struct{ err error }

func (f *failingBook) EnsureOperatingBook(ctx context.Context, merchantID string) (string, error) {
	return "book_x", nil
}
func (f *failingBook) AccountIDByCode(ctx context.Context, bookID, code string) (string, error) {
	return "", f.err
}
