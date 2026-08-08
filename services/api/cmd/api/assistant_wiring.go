package main

import (
	"context"
	"encoding/json"
	"strings"

	"apexpay/internal/accounting"
	"apexpay/internal/assistant"
	"apexpay/internal/inventory"
	"apexpay/internal/invoicing"
	"apexpay/internal/lending"
	"apexpay/internal/payment"
	"apexpay/internal/payroll"
	"apexpay/internal/treasury"
	"github.com/jackc/pgx/v5/pgxpool"
)

// buildAssistantReaders adapts the concrete domain repositories into the narrow read-only
// Readers interface the assistant depends on. Adapting here keeps the assistant package free
// of direct imports to every domain package and forces a single, explicit conversion seam.
func buildAssistantReaders(pool *pgxpool.Pool) assistant.Readers {
	return assistant.Readers{
		Payments:   &paymentsReader{repo: payment.NewPgRepository(pool, nil)},
		Invoices:   &invoicesReader{repo: invoicing.NewRepository(pool)},
		Inventory:  &inventoryReader{repo: inventory.NewRepository(pool)},
		Treasury:   &treasuryReader{repo: treasury.NewRepository(pool)},
		Loans:      &loansReader{repo: lending.NewRepository(pool)},
		Accounting: &accountingReader{repo: accounting.NewRepository(pool)},
		Employee:   &employeeReader{repo: payroll.NewPgRepository(pool, nil)},
	}
}

// ---- Adapters (concrete typed repos -> assistant.Readers map-based interface) ----

type paymentsReader struct{ repo *payment.PgRepository }

func (a *paymentsReader) DashboardSummary(ctx context.Context, merchantID string) (map[string]any, error) {
	s, err := a.repo.DashboardSummary(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	return structToMap(s), nil
}

func (a *paymentsReader) ListByMerchant(ctx context.Context, merchantID string, limit int) ([]any, error) {
	list, err := a.repo.ListByMerchant(ctx, merchantID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]any, 0, len(list))
	for _, p := range list {
		out = append(out, structToMap(p))
	}
	return out, nil
}

type invoicesReader struct{ repo *invoicing.Repository }

func (a *invoicesReader) ListInvoices(ctx context.Context, merchantID, status string, limit int) ([]map[string]any, error) {
	list, err := a.repo.ListInvoices(ctx, merchantID, status, limit)
	if err != nil {
		return nil, err
	}
	return structsToMaps(list), nil
}

func (a *invoicesReader) Aging(ctx context.Context, merchantID string) ([]map[string]any, error) {
	list, err := a.repo.Aging(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	return structsToMaps(list), nil
}

type inventoryReader struct{ repo *inventory.Repository }

func (a *inventoryReader) ListProducts(ctx context.Context, merchantID string) ([]map[string]any, error) {
	list, err := a.repo.ListProducts(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	return structsToMaps(list), nil
}

type treasuryReader struct{ repo *treasury.Repository }

func (a *treasuryReader) CashPosition(ctx context.Context, merchantID string) (map[string]any, error) {
	p, err := a.repo.CashPosition(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	return structToMap(p), nil
}

func (a *treasuryReader) LatestForecast(ctx context.Context, merchantID string) (map[string]any, error) {
	f, err := a.repo.LatestForecast(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	return structToMap(f), nil
}

type loansReader struct{ repo *lending.Repository }

func (a *loansReader) ListLoans(ctx context.Context, merchantID, status string, limit int) ([]map[string]any, error) {
	list, err := a.repo.ListLoans(ctx, merchantID, status, limit)
	if err != nil {
		return nil, err
	}
	return structsToMaps(list), nil
}

type accountingReader struct{ repo *accounting.Repository }

func (a *accountingReader) ProfitLoss(ctx context.Context, merchantID, from, to string) (map[string]any, error) {
	st, err := a.repo.ProfitLoss(ctx, merchantID, from, to)
	if err != nil {
		return nil, err
	}
	return structToMap(st), nil
}

func (a *accountingReader) BalanceSheet(ctx context.Context, merchantID, asOf string) (map[string]any, error) {
	st, err := a.repo.BalanceSheet(ctx, merchantID, asOf)
	if err != nil {
		return nil, err
	}
	return structToMap(st), nil
}

type employeeReader struct{ repo *payroll.PgRepository }

func (a *employeeReader) GetYTDForEmployee(ctx context.Context, merchantID, employeeID string, year int) (map[string]any, error) {
	m, err := a.repo.GetYTDForEmployee(ctx, merchantID, employeeID, year)
	if err != nil {
		return nil, err
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v.String()
	}
	return out, nil
}

func (a *employeeReader) LeaveBalancesByEmployee(ctx context.Context, merchantID, employeeID string, year int) ([]map[string]any, error) {
	list, err := a.repo.ListLeaveBalancesByEmployee(ctx, merchantID, employeeID, year)
	if err != nil {
		return nil, err
	}
	return lowerKeys(structsToMaps(list)), nil
}

func (a *employeeReader) ClaimsByEmployee(ctx context.Context, merchantID, employeeID string) ([]map[string]any, error) {
	list, err := a.repo.ListClaimsByEmployee(ctx, merchantID, employeeID, nil)
	if err != nil {
		return nil, err
	}
	return lowerKeys(structsToMaps(list)), nil
}

// lowerKeys lowercases the top-level keys of each map. Domain structs without JSON tags
// marshal with capitalized field names; normalizing here keeps the assistant's formatting
// helpers (which read lowercase keys) correct regardless of source struct tags.
func lowerKeys(in []map[string]any) []map[string]any {
	for i := range in {
		out := make(map[string]any, len(in[i]))
		for k, v := range in[i] {
			out[toSnakeCase(k)] = v
		}
		in[i] = out
	}
	return in
}

// toSnakeCase converts Go/CamelCase field names (e.g. "LeaveType", "RemainingDays")
// into snake_case ("leave_type", "remaining_days") so the assistant's formatting helpers,
// which read lowercase snake_case keys, work regardless of source struct JSON tags.
func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r + ('a' - 'A'))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ---- JSON-based conversion helpers ----

func structToMap(v any) map[string]any {
	b, err := json.Marshal(v)
	if err != nil {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]any{}
	}
	return m
}

func structsToMaps(v any) []map[string]any {
	b, err := json.Marshal(v)
	if err != nil {
		return []map[string]any{}
	}
	var out []map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		return []map[string]any{}
	}
	return out
}
