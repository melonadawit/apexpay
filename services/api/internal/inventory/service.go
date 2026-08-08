package inventory

import (
	"context"
	"fmt"

	"apexpay/internal/id"
	"apexpay/internal/ledger"
	"github.com/shopspring/decimal"
)

// bookResolver finds the merchant's operating ledger book and resolves account ids by code.
type bookResolver interface {
	EnsureOperatingBook(ctx context.Context, merchantID string) (string, error)
	AccountIDByCode(ctx context.Context, bookID, code string) (string, error)
}

type ledgerPoster interface {
	PostJournalTx(ctx context.Context, journal *ledger.Journal, entries []ledger.Entry) error
}

// Service wraps order creation to also post Cost of Goods Sold to the GL. When a sale
// (order) is created, inventory leaves stock; the corresponding cost is moved from the
// Inventory asset account to Cost of Sales (expense) so the P&L reflects true COGS and the
// balance sheet reflects inventory at cost.
type Service struct {
	repo   *Repository
	book   bookResolver
	ledger ledgerPoster
}

func NewService(repo *Repository, book bookResolver, ledger ledgerPoster) *Service {
	return &Service{repo: repo, book: book, ledger: ledger}
}

// CreateOrder creates the order and posts COGS to the ledger. The COGS post happens after
// the order is committed; if ledger posting fails the order still stands but we return the
// error so callers can surface a partial-posting condition.
func (s *Service) CreateOrder(ctx context.Context, merchantID string, o *Order) error {
	// Cost lookup must happen before the order mutates stock (still accurate at this point).
	cogsByItem := make([]decimal.Decimal, len(o.Items))
	for i, it := range o.Items {
		if it.ProductID == "" {
			continue
		}
		cost, err := s.repo.CostPrice(ctx, merchantID, it.ProductID)
		if err != nil {
			return fmt.Errorf("lookup cost for %s: %w", it.ProductID, err)
		}
		qty, _ := decimal.NewFromString(it.Quantity)
		cogsByItem[i] = cost.Mul(qty).Round(2)
	}

	if err := s.repo.CreateOrder(ctx, merchantID, o); err != nil {
		return err
	}

	// Sum COGS and post to ledger.
	totalCogs := decimal.Zero
	for _, c := range cogsByItem {
		totalCogs = totalCogs.Add(c)
	}
	if totalCogs.GreaterThan(decimal.Zero) {
		if err := s.postCOGS(ctx, merchantID, o.ID, totalCogs); err != nil {
			return fmt.Errorf("post COGS to ledger: %w", err)
		}
	}
	return nil
}

func (s *Service) postCOGS(ctx context.Context, merchantID, orderID string, amount decimal.Decimal) error {
	bookID, err := s.book.EnsureOperatingBook(ctx, merchantID)
	if err != nil {
		return err
	}
	cogsID, err := s.book.AccountIDByCode(ctx, bookID, "expense:cost_of_sales")
	if err != nil {
		return fmt.Errorf("cost of sales account missing: %w", err)
	}
	invID, err := s.book.AccountIDByCode(ctx, bookID, "asset:inventory")
	if err != nil {
		return fmt.Errorf("inventory account missing: %w", err)
	}

	journalID := id.NewLedgerJournal()
	journal := &ledger.Journal{
		ID: journalID, BookID: bookID,
		PostingKey:    "cogs_" + orderID,
		Memo:          "COGS " + orderID,
		ReferenceType: "order", ReferenceID: orderID,
	}
	entries := []ledger.Entry{
		{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: cogsID, Direction: "debit", Amount: amount, Currency: "ETB"},
		{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: invID, Direction: "credit", Amount: amount, Currency: "ETB"},
	}
	return s.ledger.PostJournalTx(ctx, journal, entries)
}
