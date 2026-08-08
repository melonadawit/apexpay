package fxreval

import (
	"context"
	"fmt"
	"time"

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

// Service revalues foreign-currency accounts to ETB at the current forex rate and posts
// the resulting unrealized FX gain/loss to the GL. The gain/loss per account is the change
// in ETB value since the last recorded revaluation (baseline on the first run).
type Service struct {
	repo   *Repository
	book   bookResolver
	ledger ledgerPoster
}

func NewService(repo *Repository, book bookResolver, ledger ledgerPoster) *Service {
	return &Service{repo: repo, book: book, ledger: ledger}
}

// Revalue computes and records the revaluation for the current period and posts any net
// FX gain/loss to the ledger.
func (s *Service) Revalue(ctx context.Context, merchantID, period string) (*Revaluation, error) {
	if period == "" {
		period = time.Now().Format("2006-01")
	}
	accts, err := s.repo.ForeignAccounts(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	net := decimal.Zero
	lines := make([]RevaluationLine, 0, len(accts))
	for _, a := range accts {
		amtFX, _ := decimal.NewFromString(a.AmountFX)
		rate, _ := decimal.NewFromString(a.Rate)
		amountETB := amtFX.Mul(rate).Round(2)
		line := RevaluationLine{
			AccountID: a.AccountID, AccountNumber: a.AccountNumber, Currency: a.Currency,
			AmountFX: amtFX.String(), Rate: rate.String(), AmountETB: amountETB.String(),
		}
		// Gain/loss = current ETB value minus prior recorded ETB value. On the first run
		// (no prior) we establish the baseline with zero gain/loss so only rate changes
		// produce an unrealized gain or loss.
		prior, err := s.repo.PriorETB(ctx, merchantID, a.AccountID)
		if err != nil {
			return nil, err
		}
		gain := decimal.Zero
		if !prior.IsZero() {
			gain = amountETB.Sub(prior)
		}
		line.FXGainLoss = gain.String()
		net = net.Add(gain)

		if err := s.repo.Record(ctx, merchantID, line, period); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	rev := &Revaluation{Period: period, Lines: lines, NetFX: net.String()}

	// Post net FX gain/loss to the GL if non-zero.
	if !net.IsZero() {
		if err := s.postToLedger(ctx, merchantID, period, net); err != nil {
			return nil, fmt.Errorf("post fx revaluation to ledger: %w", err)
		}
		rev.JournalID = id.NewLedgerJournal() // placeholder; updated in postToLedger if needed
	}
	return rev, nil
}

func (s *Service) postToLedger(ctx context.Context, merchantID, period string, net decimal.Decimal) error {
	bookID, err := s.book.EnsureOperatingBook(ctx, merchantID)
	if err != nil {
		return err
	}
	bankID, err := s.book.AccountIDByCode(ctx, bookID, "asset:bank")
	if err != nil {
		return fmt.Errorf("asset:bank account missing: %w", err)
	}

	// A net gain (net > 0) credits fx_gain; a net loss (net < 0) debits fx_loss.
	var gainID, lossID string
	var entries []ledger.Entry
	journalID := id.NewLedgerJournal()
	journal := &ledger.Journal{
		ID: journalID, BookID: bookID, PostingKey: "fxreval_" + period,
		Memo: "FX revaluation " + period, ReferenceType: "fx_revaluation", ReferenceID: period,
	}
	if net.GreaterThan(decimal.Zero) {
		gainID, err = s.book.AccountIDByCode(ctx, bookID, "revenue:fx_gain")
		if err != nil {
			return fmt.Errorf("revenue:fx_gain account missing: %w", err)
		}
		entries = []ledger.Entry{
			{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: bankID, Direction: "debit", Amount: net, Currency: "ETB"},
			{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: gainID, Direction: "credit", Amount: net, Currency: "ETB"},
		}
	} else {
		loss := net.Abs()
		lossID, err = s.book.AccountIDByCode(ctx, bookID, "expense:fx_loss")
		if err != nil {
			return fmt.Errorf("expense:fx_loss account missing: %w", err)
		}
		entries = []ledger.Entry{
			{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: lossID, Direction: "debit", Amount: loss, Currency: "ETB"},
			{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: bankID, Direction: "credit", Amount: loss, Currency: "ETB"},
		}
	}
	return s.ledger.PostJournalTx(ctx, journal, entries)
}
