package ledger

import (
	"time"

	"github.com/shopspring/decimal"
)

type Journal struct {
	ID            string
	BookID        string
	PostingKey    string // idempotent unique per book
	Memo          string
	TransferGroup string
	ReferenceType string
	ReferenceID   string
	CreatedAt     time.Time
}

type Entry struct {
	ID        string
	JournalID string
	BookID    string
	AccountID string
	Direction string // debit, credit
	Amount    decimal.Decimal
	Currency  string
	Meta      map[string]any
	CreatedAt time.Time
}

// Service handles balance updates - optimal with advisory locks per book
type Service struct {
	repo interface {
		PostJournalTx(journal *Journal, entries []Entry) error
	}
}

func NewService(repo interface {
	PostJournalTx(journal *Journal, entries []Entry) error
}) *Service {
	// Define wrapper to avoid import cycle; we use same repo interface via type assertion in impl
	return &Service{repo: repo}
}

// ValidateBalanced checks double-entry invariant debit==credit - critical per DATABASE
func ValidateBalanced(entries []Entry) bool {
	if len(entries) == 0 {
		return false
	}
	var debit, credit decimal.Decimal
	for _, e := range entries {
		if e.Direction == "debit" {
			debit = debit.Add(e.Amount)
		} else {
			credit = credit.Add(e.Amount)
		}
	}
	return debit.Equal(credit)
}
