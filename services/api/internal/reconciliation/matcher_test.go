package reconciliation

import (
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestMatchLinesPrefersReferenceAndRejectsAmbiguity(t *testing.T) {
	now := time.Date(2026, 8, 7, 10, 0, 0, 0, time.UTC)
	amount := decimal.NewFromInt(100)
	lines := []StatementLine{{ID: "s1", ConnectorRef: "ref-1", Amount: amount, Currency: "ETB", OccurredAt: now}, {ID: "s2", Amount: amount, Currency: "ETB", OccurredAt: now}}
	journals := []LedgerCandidate{{JournalID: "j1", ConnectorRef: "ref-1", Amount: amount, Currency: "ETB", OccurredAt: now}, {JournalID: "j2", Amount: amount, Currency: "ETB", OccurredAt: now.Add(time.Hour)}}
	matches, breaks := MatchLines(lines, journals)
	if len(matches) != 1 || matches[0].Method != "connector_ref" {
		t.Fatal("reference match must win")
	}
	if len(breaks) != 1 {
		t.Fatal("ambiguous fallback must be a break")
	}
}
