package accounting

import (
	"context"
	"fmt"
	"testing"
	"time"

	"apexpay/internal/ledger"
)

// fakeRepo satisfies accountingRepo without a DB.
type fakeRepo struct {
	status map[string]string
	postN  int
}

func (f *fakeRepo) EnsureOperatingBook(ctx context.Context, merchantID string) (string, error) {
	return "book_x", nil
}
func (f *fakeRepo) AccountIDByCode(ctx context.Context, bookID, code string) (string, error) {
	return "acc_" + code, nil
}
func (f *fakeRepo) AccountNameByCode(ctx context.Context, bookID, code string) (string, error) {
	return code, nil
}
func (f *fakeRepo) PeriodStatus(ctx context.Context, merchantID, period string) (string, error) {
	if f.status == nil {
		return "open", nil
	}
	if s, ok := f.status[period]; ok {
		return s, nil
	}
	return "open", nil
}
func (f *fakeRepo) SetPeriodStatus(ctx context.Context, merchantID, period, status, userID string) error {
	if f.status == nil {
		f.status = map[string]string{}
	}
	f.status[period] = status
	return nil
}
func (f *fakeRepo) ListPeriods(ctx context.Context, merchantID string) ([]FiscalPeriod, error) {
	return nil, nil
}
func (f *fakeRepo) ListManualJournals(ctx context.Context, merchantID string) ([]JournalEntry, error) {
	return nil, nil
}

type fakePoster struct{ err error }

func (f *fakePoster) PostJournalTx(ctx context.Context, j *ledger.Journal, e []ledger.Entry) error {
	return f.err
}

func newFakeSvc() (*Service, *fakeRepo) {
	r := &fakeRepo{}
	return NewService(r, &fakePoster{}), r
}

func TestPostJournalEntry_Balanced(t *testing.T) {
	svc, _ := newFakeSvc()
	req := JournalEntryRequest{
		Memo: "test",
		Lines: []JournalLine{
			{AccountCode: "asset:bank", Direction: "debit", Amount: "1000.00"},
			{AccountCode: "revenue:product", Direction: "credit", Amount: "1000.00"},
		},
	}
	je, err := svc.PostJournalEntry(context.Background(), "m1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if je.ID == "" || len(je.Lines) != 2 {
		t.Fatalf("bad result: %+v", je)
	}
}

func TestPostJournalEntry_Unbalanced(t *testing.T) {
	svc, _ := newFakeSvc()
	req := JournalEntryRequest{
		Memo: "bad",
		Lines: []JournalLine{
			{AccountCode: "asset:bank", Direction: "debit", Amount: "1000.00"},
			{AccountCode: "revenue:product", Direction: "credit", Amount: "900.00"},
		},
	}
	if _, err := svc.PostJournalEntry(context.Background(), "m1", req); err == nil {
		t.Fatal("expected error for unbalanced entry")
	}
}

func TestPostJournalEntry_SingleLine(t *testing.T) {
	svc, _ := newFakeSvc()
	req := JournalEntryRequest{Memo: "one", Lines: []JournalLine{{AccountCode: "a", Direction: "debit", Amount: "1"}}}
	if _, err := svc.PostJournalEntry(context.Background(), "m1", req); err == nil {
		t.Fatal("expected error for <2 lines")
	}
}

func TestPostJournalEntry_ClosedPeriod(t *testing.T) {
	svc, r := newFakeSvc()
	closed := time.Now().Format("2006-01")
	r.status = map[string]string{closed: "closed"}
	req := JournalEntryRequest{
		Memo: "closed test",
		Lines: []JournalLine{
			{AccountCode: "asset:bank", Direction: "debit", Amount: "100"},
			{AccountCode: "revenue:product", Direction: "credit", Amount: "100"},
		},
	}
	if _, err := svc.PostJournalEntry(context.Background(), "m1", req); err == nil {
		t.Fatal("expected period-closed error")
	}
}

func TestPostJournalEntry_InvalidDirection(t *testing.T) {
	svc, _ := newFakeSvc()
	req := JournalEntryRequest{
		Memo: "bad dir",
		Lines: []JournalLine{
			{AccountCode: "asset:bank", Direction: "sideways", Amount: "100"},
			{AccountCode: "revenue:product", Direction: "credit", Amount: "100"},
		},
	}
	if _, err := svc.PostJournalEntry(context.Background(), "m1", req); err == nil {
		t.Fatal("expected invalid-direction error")
	}
}

// TestPostJournalEntry_MissingAccount covers the error path when an account code is not in
// the chart of accounts (AccountIDByCode fails).
func TestPostJournalEntry_MissingAccount(t *testing.T) {
	repo := &missingAccountRepo{}
	svc := NewService(repo, &fakePoster{})
	req := JournalEntryRequest{
		Memo: "missing acct",
		Lines: []JournalLine{
			{AccountCode: "asset:unknown", Direction: "debit", Amount: "100"},
			{AccountCode: "revenue:product", Direction: "credit", Amount: "100"},
		},
	}
	if _, err := svc.PostJournalEntry(context.Background(), "m1", req); err == nil {
		t.Fatal("expected missing-account error")
	}
}

// TestPeriodCloseThenReopen verifies the close -> reopen lifecycle returns proper status.
func TestPeriodCloseThenReopen(t *testing.T) {
	svc, r := newFakeSvc()
	closed, err := svc.ClosePeriod(context.Background(), "m1", "2026-08", "u1")
	if err != nil || closed.Status != "closed" {
		t.Fatalf("ClosePeriod: %v %+v", err, closed)
	}
	if r.status["2026-08"] != "closed" {
		t.Fatal("expected period marked closed in repo")
	}
	reopened, err := svc.ReopenPeriod(context.Background(), "m1", "2026-08")
	if err != nil || reopened.Status != "open" {
		t.Fatalf("ReopenPeriod: %v %+v", err, reopened)
	}
	if r.status["2026-08"] != "open" {
		t.Fatal("expected period marked open in repo")
	}
}

// TestClosePeriodRequiresValue guards the empty-period error path.
func TestClosePeriodRequiresValue(t *testing.T) {
	svc, _ := newFakeSvc()
	if _, err := svc.ClosePeriod(context.Background(), "m1", "", "u1"); err == nil {
		t.Fatal("expected error for empty period")
	}
}

// missingAccountRepo fails AccountIDByCode for all codes.
type missingAccountRepo struct{ fakeRepo }

func (f *missingAccountRepo) AccountIDByCode(ctx context.Context, bookID, code string) (string, error) {
	return "", fmt.Errorf("account %s not in chart of accounts", code)
}
