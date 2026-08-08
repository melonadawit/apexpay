package reconciliation

import (
	"strings"
	"testing"
)

func TestParseCSVRejectsDuplicateAndParsesNormalizedLine(t *testing.T) {
	csv := "transaction_id,connector_ref,amount,currency,occurred_at\ns1,ref1,100.00,etb,2026-08-07T10:00:00Z\n"
	lines, e := ParseCSV(strings.NewReader(csv))
	if e != nil || len(lines) != 1 || lines[0].Currency != "ETB" {
		t.Fatal("valid normalized csv must parse")
	}
	_, e = ParseCSV(strings.NewReader("transaction_id,connector_ref,amount,currency,occurred_at\ns1,,1,ETB,2026-08-07T10:00:00Z\ns1,,1,ETB,2026-08-07T10:01:00Z\n"))
	if e == nil {
		t.Fatal("duplicate external transaction IDs must be rejected")
	}
}
