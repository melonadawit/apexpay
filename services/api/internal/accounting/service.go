package accounting

import (
	"context"
	"fmt"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/ledger"
	"apexpay/internal/platform/errors"
	"github.com/shopspring/decimal"
)

// repo is the subset of persistence the service needs, kept as an interface for testability.
type accountingRepo interface {
	EnsureOperatingBook(ctx context.Context, merchantID string) (string, error)
	AccountIDByCode(ctx context.Context, bookID, code string) (string, error)
	AccountNameByCode(ctx context.Context, bookID, code string) (string, error)
	PeriodStatus(ctx context.Context, merchantID, period string) (string, error)
	SetPeriodStatus(ctx context.Context, merchantID, period, status, userID string) error
	ListPeriods(ctx context.Context, merchantID string) ([]FiscalPeriod, error)
	ListManualJournals(ctx context.Context, merchantID string) ([]JournalEntry, error)
}

// Service posts manual journal entries to the merchant's operating ledger and manages
// fiscal period open/close. It enforces double-entry balance, period locks, and resolves
// account codes to real ledger account ids before posting.
type Service struct {
	repo       accountingRepo
	ledgerRepo ledgerPoster
}

type ledgerPoster interface {
	PostJournalTx(ctx context.Context, journal *ledger.Journal, entries []ledger.Entry) error
}

func NewService(repo accountingRepo, ledgerRepo ledgerPoster) *Service {
	return &Service{repo: repo, ledgerRepo: ledgerRepo}
}

// PostJournalEntry validates and posts a balanced manual journal entry, or returns a
// validation/period-closed error without touching the ledger.
func (s *Service) PostJournalEntry(ctx context.Context, merchantID string, req JournalEntryRequest) (*JournalEntry, error) {
	if len(req.Lines) < 2 {
		return nil, errors.Validation("journal entry requires at least 2 lines")
	}

	// Resolve operating book (create chart of accounts if absent).
	bookID, err := s.repo.EnsureOperatingBook(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("resolve operating book: %w", err)
	}

	// Validate lines: direction, amount, account existence.
	entries := make([]ledger.Entry, 0, len(req.Lines))
	lines := make([]JournalLine, 0, len(req.Lines))
	var totalDebit, totalCredit decimal.Decimal
	for _, l := range req.Lines {
		if l.Direction != "debit" && l.Direction != "credit" {
			return nil, errors.Validation(fmt.Sprintf("invalid direction %q for account %s", l.Direction, l.AccountCode))
		}
		amt, err := decimal.NewFromString(l.Amount)
		if err != nil || amt.IsNegative() {
			return nil, errors.Validation(fmt.Sprintf("invalid amount %q for account %s", l.Amount, l.AccountCode))
		}
		acid, err := s.repo.AccountIDByCode(ctx, bookID, l.AccountCode)
		if err != nil {
			return nil, errors.Validation(fmt.Sprintf("account code %q not in chart of accounts", l.AccountCode))
		}
		name, _ := s.repo.AccountNameByCode(ctx, bookID, l.AccountCode)
		if l.Direction == "debit" {
			totalDebit = totalDebit.Add(amt)
		} else {
			totalCredit = totalCredit.Add(amt)
		}
		entries = append(entries, ledger.Entry{
			ID: id.New("lent"), BookID: bookID, AccountID: acid,
			Direction: l.Direction, Amount: amt, Currency: "ETB",
		})
		lines = append(lines, JournalLine{AccountCode: l.AccountCode, AccountName: name, Direction: l.Direction, Amount: amt.String()})
	}

	if !totalDebit.Equal(totalCredit) {
		return nil, errors.Validation(fmt.Sprintf("journal not balanced: debit %s != credit %s", totalDebit, totalCredit))
	}

	// Period lock check.
	period := time.Now().Format("2006-01")
	status, err := s.repo.PeriodStatus(ctx, merchantID, period)
	if err != nil {
		return nil, err
	}
	if status == "closed" {
		return nil, errors.Validation(fmt.Sprintf("period %s is closed; reopen before posting", period))
	}

	journalID := id.NewLedgerJournal()
	journal := &ledger.Journal{
		ID: journalID, BookID: bookID,
		PostingKey: "manual_" + period + "_" + id.New("jk"),
		Memo:       req.Memo, TransferGroup: "",
		ReferenceType: "manual", ReferenceID: req.ReferenceID,
	}
	for i := range entries {
		entries[i].JournalID = journalID
	}

	if err := s.ledgerRepo.PostJournalTx(ctx, journal, entries); err != nil {
		return nil, fmt.Errorf("post journal: %w", err)
	}

	return &JournalEntry{
		ID: journalID, Memo: req.Memo, Period: period,
		PostingKey: journal.PostingKey, Lines: lines,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// ListJournalEntries returns manual journal entries for the merchant.
func (s *Service) ListJournalEntries(ctx context.Context, merchantID string) ([]JournalEntry, error) {
	return s.repo.ListManualJournals(ctx, merchantID)
}

// ListPeriods returns fiscal periods for the merchant.
func (s *Service) ListPeriods(ctx context.Context, merchantID string) ([]FiscalPeriod, error) {
	return s.repo.ListPeriods(ctx, merchantID)
}

// ClosePeriod locks a fiscal period against new postings.
func (s *Service) ClosePeriod(ctx context.Context, merchantID, period, userID string) (*FiscalPeriod, error) {
	if period == "" {
		return nil, errors.Validation("period required (YYYY-MM)")
	}
	if err := s.repo.SetPeriodStatus(ctx, merchantID, period, "closed", userID); err != nil {
		return nil, err
	}
	return &FiscalPeriod{Merchant: merchantID, Period: period, Status: "closed"}, nil
}

// ReopenPeriod unlocks a closed fiscal period (explicit operator action).
func (s *Service) ReopenPeriod(ctx context.Context, merchantID, period string) (*FiscalPeriod, error) {
	if period == "" {
		return nil, errors.Validation("period required (YYYY-MM)")
	}
	if err := s.repo.SetPeriodStatus(ctx, merchantID, period, "open", ""); err != nil {
		return nil, err
	}
	return &FiscalPeriod{Merchant: merchantID, Period: period, Status: "open"}, nil
}
