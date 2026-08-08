package fixedasset

import (
	"context"
	"fmt"

	"apexpay/internal/id"
	"apexpay/internal/ledger"
	"github.com/shopspring/decimal"
)

// bookResolver finds the merchant's operating ledger book and resolves account ids by code.
// Kept as a small interface so the service is testable without a DB.
type bookResolver interface {
	EnsureOperatingBook(ctx context.Context, merchantID string) (string, error)
	AccountIDByCode(ctx context.Context, bookID, code string) (string, error)
}

type ledgerPoster interface {
	PostJournalTx(ctx context.Context, journal *ledger.Journal, entries []ledger.Entry) error
}

// Service records a depreciation entry and posts the corresponding double-entry journal
// (debit Depreciation Expense, credit Accumulated Depreciation) to the operating ledger so
// accumulated depreciation flows into the balance sheet. Both happen so depreciation never
// changes the ledger without a matching entry (and vice-versa).
type Service struct {
	repo   *Repository
	book   bookResolver
	ledger ledgerPoster
}

func NewService(repo *Repository, book bookResolver, ledger ledgerPoster) *Service {
	return &Service{repo: repo, book: book, ledger: ledger}
}

// Depreciate computes and records one month's depreciation, then posts it to the ledger.
func (s *Service) Depreciate(ctx context.Context, merchantID, assetID string) (*DepreciationEntry, error) {
	entry, err := s.repo.Depreciate(ctx, merchantID, assetID)
	if err != nil {
		return nil, err
	}
	amount, err := decimal.NewFromString(entry.Amount)
	if err != nil {
		return nil, err
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		// Nothing to post (asset fully depreciated / zero annual).
		return entry, nil
	}
	if err := s.postToLedger(ctx, merchantID, assetID, entry.Period, amount); err != nil {
		return nil, fmt.Errorf("post depreciation to ledger: %w", err)
	}
	return entry, nil
}

func (s *Service) postToLedger(ctx context.Context, merchantID, assetID, period string, amount decimal.Decimal) error {
	bookID, err := s.book.EnsureOperatingBook(ctx, merchantID)
	if err != nil {
		return err
	}
	depExpID, err := s.book.AccountIDByCode(ctx, bookID, "expense:depreciation")
	if err != nil {
		return fmt.Errorf("depreciation expense account missing: %w", err)
	}
	accumDepID, err := s.book.AccountIDByCode(ctx, bookID, "asset:accumulated_depreciation")
	if err != nil {
		return fmt.Errorf("accumulated depreciation account missing: %w", err)
	}

	journalID := id.NewLedgerJournal()
	journal := &ledger.Journal{
		ID: journalID, BookID: bookID,
		PostingKey:    "dep_" + period + "_" + assetID,
		Memo:          "Depreciation " + period,
		ReferenceType: "depreciation", ReferenceID: assetID,
	}
	entries := []ledger.Entry{
		{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: depExpID, Direction: "debit", Amount: amount, Currency: "ETB"},
		{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: accumDepID, Direction: "credit", Amount: amount, Currency: "ETB"},
	}
	return s.ledger.PostJournalTx(ctx, journal, entries)
}
