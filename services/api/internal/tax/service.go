package tax

import (
	"context"
	"fmt"

	"apexpay/internal/id"
	"apexpay/internal/ledger"
	"github.com/shopspring/decimal"
)

// bookResolver finds the operating ledger book and resolves account ids by code.
type bookResolver interface {
	EnsureOperatingBook(ctx context.Context, merchantID string) (string, error)
	AccountIDByCode(ctx context.Context, bookID, code string) (string, error)
}

type ledgerPoster interface {
	PostJournalTx(ctx context.Context, journal *ledger.Journal, entries []ledger.Entry) error
}

// Service builds the tax schedule and posts collected VAT/TOT into the GL tax-liability
// account. When tax is collected on an invoice, the merchant owes it to the tax authority,
// so we credit a liability:tax account (with an offsetting receivable/cash debit) so the
// balance sheet reflects the tax obligation.
// taxRepo is the persistence subset the service needs (interface for testability).
type taxRepo interface {
	Schedule(ctx context.Context, merchantID string) ([]ScheduleLine, error)
	Record(ctx context.Context, merchantID, period, taxType, source, sourceID, amount string) error
}

type Service struct {
	repo   taxRepo
	book   bookResolver
	ledger ledgerPoster
}

func NewService(repo taxRepo, book bookResolver, ledger ledgerPoster) *Service {
	return &Service{repo: repo, book: book, ledger: ledger}
}

// Schedule returns the merchant's tax schedule with the net amount due.
func (s *Service) Schedule(ctx context.Context, merchantID string) (*Schedule, error) {
	lines, err := s.repo.Schedule(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	totalDue := decimal.Zero
	for i := range lines {
		collected, _ := decimal.NewFromString(lines[i].Collected)
		paid, _ := decimal.NewFromString(lines[i].Paid)
		due := collected.Sub(paid)
		lines[i].Due = due.String()
		totalDue = totalDue.Add(due)
	}
	return &Schedule{MerchantID: merchantID, Lines: lines, TotalDue: totalDue.String()}, nil
}

// RecordCollected adds a collected-tax line to the register for a period.
func (s *Service) RecordCollected(ctx context.Context, merchantID, period, taxType, source, sourceID, amount string) error {
	return s.repo.Record(ctx, merchantID, period, taxType, source, sourceID, amount)
}

// PostToLedger posts the outstanding VAT/TOT/withholding due for a period into the GL
// tax-liability account (debit accounts receivable, credit liability:tax). Returns the
// journal id. No-op if nothing is due.
func (s *Service) PostToLedger(ctx context.Context, merchantID, period string) (*PostResult, error) {
	lines, err := s.repo.Schedule(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	total := decimal.Zero
	for _, l := range lines {
		if period != "" && l.Period != period {
			continue
		}
		collected, _ := decimal.NewFromString(l.Collected)
		paid, _ := decimal.NewFromString(l.Paid)
		total = total.Add(collected.Sub(paid))
	}
	if total.LessThanOrEqual(decimal.Zero) {
		return &PostResult{Period: period, Amount: "0"}, nil
	}

	bookID, err := s.book.EnsureOperatingBook(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	taxID, err := s.book.AccountIDByCode(ctx, bookID, "liability:tax")
	if err != nil {
		return nil, fmt.Errorf("tax liability account missing: %w", err)
	}
	recvID, err := s.book.AccountIDByCode(ctx, bookID, "asset:receivable")
	if err != nil {
		return nil, fmt.Errorf("accounts receivable account missing: %w", err)
	}

	journalID := id.NewLedgerJournal()
	pk := "tax_" + period
	if period == "" {
		pk = "tax_all"
	}
	journal := &ledger.Journal{
		ID: journalID, BookID: bookID, PostingKey: pk,
		Memo: "Sales tax liability " + period, ReferenceType: "tax", ReferenceID: period,
	}
	entries := []ledger.Entry{
		{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: recvID, Direction: "debit", Amount: total, Currency: "ETB"},
		{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: taxID, Direction: "credit", Amount: total, Currency: "ETB"},
	}
	if err := s.ledger.PostJournalTx(ctx, journal, entries); err != nil {
		return nil, err
	}
	return &PostResult{JournalID: journalID, Period: period, Amount: total.String(), Account: "liability:tax"}, nil
}
