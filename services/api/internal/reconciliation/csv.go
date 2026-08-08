package reconciliation

import (
	"encoding/csv"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"strings"
	"time"
)

// ParseCSV accepts a deliberately small normalized settlement format:
// transaction_id,connector_ref,amount,currency,occurred_at (RFC3339).
// Unknown/malformed rows fail closed; financial imports are never best-effort.
func ParseCSV(r io.Reader) ([]StatementLine, error) {
	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("statement must contain a header and at least one row")
	}
	expected := []string{"transaction_id", "connector_ref", "amount", "currency", "occurred_at"}
	if len(rows[0]) != len(expected) {
		return nil, fmt.Errorf("invalid statement header")
	}
	for i := range expected {
		if strings.TrimSpace(strings.ToLower(rows[0][i])) != expected[i] {
			return nil, fmt.Errorf("invalid statement header")
		}
	}
	out := make([]StatementLine, 0, len(rows)-1)
	seen := map[string]bool{}
	for rowNo, row := range rows[1:] {
		if len(row) != 5 {
			return nil, fmt.Errorf("row %d: expected 5 columns", rowNo+2)
		}
		id := strings.TrimSpace(row[0])
		if id == "" || seen[id] {
			return nil, fmt.Errorf("row %d: missing or duplicate transaction_id", rowNo+2)
		}
		amount, e := decimal.NewFromString(strings.TrimSpace(row[2]))
		if e != nil || !amount.GreaterThan(decimal.Zero) {
			return nil, fmt.Errorf("row %d: invalid amount", rowNo+2)
		}
		at, e := time.Parse(time.RFC3339, strings.TrimSpace(row[4]))
		if e != nil {
			return nil, fmt.Errorf("row %d: invalid occurred_at", rowNo+2)
		}
		seen[id] = true
		out = append(out, StatementLine{ID: id, ConnectorRef: strings.TrimSpace(row[1]), Amount: amount, Currency: strings.ToUpper(strings.TrimSpace(row[3])), OccurredAt: at})
	}
	return out, nil
}
