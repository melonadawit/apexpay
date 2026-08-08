package reconciliation

import (
	"time"

	"github.com/shopspring/decimal"
)

// StatementLine is the normalized, rail-independent representation imported
// from a bank/connector CSV. Parsers must preserve the original reference.
type StatementLine struct {
	ID, ConnectorRef, Currency string
	Amount                     decimal.Decimal
	OccurredAt                 time.Time
}

type LedgerCandidate struct {
	JournalID, ConnectorRef, Currency string
	Amount                            decimal.Decimal
	OccurredAt                        time.Time
}

type Match struct {
	StatementID, JournalID string
	Method                 string
}
type Break struct {
	StatementID string
	Reason      string
}

// MatchLines first uses exact connector references (O(n+m)), then a bounded
// amount/currency/day bucket fallback. Fallback accepts only one candidate
// within 24h, preventing an ambiguous financial auto-match.
func MatchLines(lines []StatementLine, journals []LedgerCandidate) ([]Match, []Break) {
	byRef := make(map[string]LedgerCandidate, len(journals))
	byBucket := make(map[string][]LedgerCandidate, len(journals))
	for _, journal := range journals {
		if journal.ConnectorRef != "" {
			byRef[journal.ConnectorRef] = journal
		}
		key := bucketKey(journal.Amount, journal.Currency, journal.OccurredAt)
		byBucket[key] = append(byBucket[key], journal)
	}
	matches := make([]Match, 0, len(lines))
	breaks := make([]Break, 0)
	for _, line := range lines {
		if line.ConnectorRef != "" {
			if journal, ok := byRef[line.ConnectorRef]; ok && journal.Currency == line.Currency && journal.Amount.Equal(line.Amount) {
				matches = append(matches, Match{line.ID, journal.JournalID, "connector_ref"})
				continue
			}
		}
		candidates := byBucket[bucketKey(line.Amount, line.Currency, line.OccurredAt)]
		var chosen *LedgerCandidate
		for i := range candidates {
			if candidates[i].Amount.Equal(line.Amount) && absDuration(candidates[i].OccurredAt.Sub(line.OccurredAt)) <= 24*time.Hour {
				if chosen != nil {
					chosen = nil
					break
				}
				chosen = &candidates[i]
			}
		}
		if chosen == nil {
			breaks = append(breaks, Break{line.ID, "no unique matching journal"})
			continue
		}
		matches = append(matches, Match{line.ID, chosen.JournalID, "amount_currency_24h"})
	}
	return matches, breaks
}
func bucketKey(amount decimal.Decimal, currency string, at time.Time) string {
	return amount.String() + "|" + currency + "|" + at.UTC().Format("2006-01-02")
}
func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
