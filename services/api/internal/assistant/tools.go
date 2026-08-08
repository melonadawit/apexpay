package assistant

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ---- Minimal interfaces (loose coupling). main.go passes the concrete repos. ----

type PaymentsReader interface {
	DashboardSummary(ctx context.Context, merchantID string) (map[string]any, error)
	ListByMerchant(ctx context.Context, merchantID string, limit int) ([]any, error)
}

type InvoicesReader interface {
	ListInvoices(ctx context.Context, merchantID, status string, limit int) ([]map[string]any, error)
	Aging(ctx context.Context, merchantID string) ([]map[string]any, error)
}

type InventoryReader interface {
	ListProducts(ctx context.Context, merchantID string) ([]map[string]any, error)
}

type TreasuryReader interface {
	CashPosition(ctx context.Context, merchantID string) (map[string]any, error)
	LatestForecast(ctx context.Context, merchantID string) (map[string]any, error)
}

type LoansReader interface {
	ListLoans(ctx context.Context, merchantID, status string, limit int) ([]map[string]any, error)
}

type AccountingReader interface {
	ProfitLoss(ctx context.Context, merchantID, from, to string) (map[string]any, error)
	BalanceSheet(ctx context.Context, merchantID, asOf string) (map[string]any, error)
}

// EmployeeReader exposes an employee's own self-service data only.
type EmployeeReader interface {
	GetYTDForEmployee(ctx context.Context, merchantID, employeeID string, year int) (map[string]any, error)
	LeaveBalancesByEmployee(ctx context.Context, merchantID, employeeID string, year int) ([]map[string]any, error)
	ClaimsByEmployee(ctx context.Context, merchantID, employeeID string) ([]map[string]any, error)
}

// Readers bundles all data sources the assistant can read. This is the single seam through
// which the assistant reaches business data — kept deliberately narrow and read-only.
type Readers struct {
	Payments   PaymentsReader
	Invoices   InvoicesReader
	Inventory  InventoryReader
	Treasury   TreasuryReader
	Loans      LoansReader
	Accounting AccountingReader
	Employee   EmployeeReader
}

// BuildTools constructs the read-only tool registry wired to the given readers.
func BuildTools(r Readers) []Tool {
	return []Tool{
		{Name: "summary", Description: "Business summary", Actors: []ActorType{ActorMerchant},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return merchantSummary(ctx, r, s) }},
		{Name: "payments", Description: "Recent payments", Actors: []ActorType{ActorMerchant},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return payments(ctx, r, s) }},
		{Name: "invoices", Description: "Invoices + aging", Actors: []ActorType{ActorMerchant},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return invoices(ctx, r, s) }},
		{Name: "inventory", Description: "Inventory + stock levels", Actors: []ActorType{ActorMerchant},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return inventory(ctx, r, s) }},
		{Name: "treasury", Description: "Cash position + forecast", Actors: []ActorType{ActorMerchant},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return treasury(ctx, r, s) }},
		{Name: "loans", Description: "Loans outstanding", Actors: []ActorType{ActorMerchant},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return loans(ctx, r, s) }},
		{Name: "profit_loss", Description: "Profit & loss", Actors: []ActorType{ActorMerchant},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return profitLoss(ctx, r, s) }},
		{Name: "balance_sheet", Description: "Balance sheet", Actors: []ActorType{ActorMerchant},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return balanceSheet(ctx, r, s) }},
		{Name: "my_pay", Description: "My YTD pay", Actors: []ActorType{ActorEmployee},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return myPay(ctx, r, s) }},
		{Name: "leave_balance", Description: "My leave balance", Actors: []ActorType{ActorEmployee},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return myLeave(ctx, r, s) }},
		{Name: "my_claims", Description: "My expense claims", Actors: []ActorType{ActorEmployee},
			Run: func(ctx context.Context, s Scope) (ToolResult, error) { return myClaims(ctx, r, s) }},
	}
}

// ---- Merchant tools ----

func merchantSummary(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	sum, err := r.Payments.DashboardSummary(ctx, s.MerchantID)
	if err != nil {
		return ToolResult{}, err
	}
	line := formatSummary(sum)
	return ToolResult{Line: line, Data: sum}, nil
}

func payments(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	list, err := r.Payments.ListByMerchant(ctx, s.MerchantID, 5)
	if err != nil {
		return ToolResult{}, err
	}
	if len(list) == 0 {
		return ToolResult{Line: "No recent payments yet.", Data: map[string]any{"count": 0}}, nil
	}
	return ToolResult{Line: fmt.Sprintf("You have %d recent payments.", len(list)), Data: map[string]any{"count": len(list), "payments": list}}, nil
}

func invoices(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	aging, err := r.Invoices.Aging(ctx, s.MerchantID)
	if err != nil {
		return ToolResult{}, err
	}
	var totalOverdue string
	overdueCount := 0
	var buckets []string
	for _, b := range aging {
		if bucket, ok := b["bucket"].(string); ok {
			buckets = append(buckets, bucket)
			if bucket == "overdue" || bucket == "90plus" {
				overdueCount += toInt(b["count"])
				totalOverdue = toStr(b["amount"])
			}
		}
	}
	if overdueCount == 0 {
		return ToolResult{Line: "No overdue invoices.", Data: map[string]any{"aging": aging}}, nil
	}
	return ToolResult{Line: fmt.Sprintf("%d invoices are overdue totalling %s ETB.", overdueCount, totalOverdue), Data: map[string]any{"aging": aging}}, nil
}

func inventory(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	prods, err := r.Inventory.ListProducts(ctx, s.MerchantID)
	if err != nil {
		return ToolResult{}, err
	}
	low := 0
	for _, p := range prods {
		if strings.EqualFold(toStr(p["status"]), "low_stock") {
			low++
		}
	}
	return ToolResult{Line: fmt.Sprintf("You have %d products in inventory, %d below reorder threshold.", len(prods), low),
		Data: map[string]any{"products": len(prods), "low_stock": low}}, nil
}

func treasury(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	pos, err := r.Treasury.CashPosition(ctx, s.MerchantID)
	if err != nil {
		return ToolResult{}, err
	}
	bal := toStr(pos["total_balance"])
	fc, _ := r.Treasury.LatestForecast(ctx, s.MerchantID)
	line := fmt.Sprintf("Cash position is %s ETB.", bal)
	if fc != nil {
		net := toStr(fc["net_90d"])
		if net != "" && net != "0" {
			line += fmt.Sprintf(" Net 90-day cash flow is forecast at %s ETB.", net)
		}
	}
	return ToolResult{Line: line, Data: pos}, nil
}

func loans(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	list, err := r.Loans.ListLoans(ctx, s.MerchantID, "", 10)
	if err != nil {
		return ToolResult{}, err
	}
	if len(list) == 0 {
		return ToolResult{Line: "No loans on record.", Data: map[string]any{"count": 0}}, nil
	}
	outstanding := 0.0
	for _, l := range list {
		outstanding += toFloat(l["outstanding_amount"])
	}
	return ToolResult{Line: fmt.Sprintf("You have %d loans with %s ETB outstanding.", len(list), formatFloat(outstanding)),
		Data: map[string]any{"count": len(list), "outstanding": formatFloat(outstanding)}}, nil
}

func profitLoss(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	from, to := defaultPeriod()
	st, err := r.Accounting.ProfitLoss(ctx, s.MerchantID, from, to)
	if err != nil {
		return ToolResult{}, err
	}
	return ToolResult{Line: formatStatement(st), Data: st}, nil
}

func balanceSheet(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	st, err := r.Accounting.BalanceSheet(ctx, s.MerchantID, "")
	if err != nil {
		return ToolResult{}, err
	}
	return ToolResult{Line: formatStatement(st), Data: st}, nil
}

// defaultPeriod returns a trailing 30-day window (YYYY-MM-DD) as the accounting default,
// mirroring the dashboard handlers so date-cast queries never see empty strings.
func defaultPeriod() (string, string) {
	now := time.Now()
	from := now.AddDate(0, 0, -30).Format("2006-01-02")
	to := now.Format("2006-01-02")
	return from, to
}

// ---- Employee tools (self-service, scoped to the caller's employee row) ----

func myPay(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	ytd, err := r.Employee.GetYTDForEmployee(ctx, s.MerchantID, s.EmployeeID, time.Now().Year())
	if err != nil {
		return ToolResult{}, err
	}
	return ToolResult{Line: fmt.Sprintf("Your YTD gross is %s, net is %s, tax %s ETB.",
		toStr(ytd["ytd_gross"]), toStr(ytd["ytd_net"]), toStr(ytd["ytd_tax"])), Data: ytd}, nil
}

func myLeave(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	list, err := r.Employee.LeaveBalancesByEmployee(ctx, s.MerchantID, s.EmployeeID, time.Now().Year())
	if err != nil {
		return ToolResult{}, err
	}
	if len(list) == 0 {
		return ToolResult{Line: "No leave balance on record for this year.", Data: map[string]any{"count": 0}}, nil
	}
	annual := ""
	for _, b := range list {
		if toStr(b["leave_type"]) == "annual" {
			annual = toStr(b["remaining_days"])
		}
	}
	if annual == "" {
		return ToolResult{Line: fmt.Sprintf("You have %d leave types on record this year.", len(list)), Data: map[string]any{"count": len(list), "leave_types": list}}, nil
	}
	return ToolResult{Line: fmt.Sprintf("You have %s annual leave days remaining this year.", annual), Data: map[string]any{"annual_remaining": annual, "leave_types": list}}, nil
}

func myClaims(ctx context.Context, r Readers, s Scope) (ToolResult, error) {
	list, err := r.Employee.ClaimsByEmployee(ctx, s.MerchantID, s.EmployeeID)
	if err != nil {
		return ToolResult{}, err
	}
	if len(list) == 0 {
		return ToolResult{Line: "You have no expense claims on record.", Data: map[string]any{"count": 0}}, nil
	}
	pending := 0
	total := 0.0
	for _, c := range list {
		total += toFloat(c["amount"])
		if strings.EqualFold(toStr(c["status"]), "pending") {
			pending++
		}
	}
	return ToolResult{Line: fmt.Sprintf("You have %d expense claims (%d pending) totalling %s ETB.",
		len(list), pending, formatFloat(total)), Data: map[string]any{"count": len(list), "pending": pending, "total": formatFloat(total)}}, nil
}

// ---- formatting helpers ----

func formatStatement(st map[string]any) string {
	title := toStr(st["title"])
	if title == "" {
		title = "Statement"
	}
	// Pull the first "total" line if present for a concise headline.
	if lines, ok := st["lines"].([]any); ok && len(lines) > 0 {
		last := lines[len(lines)-1]
		if lm, ok := last.(map[string]any); ok && toStr(lm["kind"]) == "total" {
			return fmt.Sprintf("%s: %s %s.", title, toStr(lm["label"]), toStr(lm["amount"]))
		}
	}
	return title + " is available."
}

func formatSummary(sum map[string]any) string {
	tpv := toStr(sum["tpv_today"])
	if tpv == "" {
		tpv = toStr(sum["tpv"])
	}
	count := toStr(sum["count"])
	if count == "" {
		count = toStr(sum["payment_count"])
	}
	return fmt.Sprintf("Today's TPV is %s ETB across %s transactions.", tpv, count)
}

func toStr(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	if f, ok := v.(float64); ok {
		return formatFloat(f)
	}
	return fmt.Sprintf("%v", v)
}

func toInt(v any) int {
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case string:
		i, _ := strconv.Atoi(n)
		return i
	}
	return 0
}

func toFloat(v any) float64 {
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case string:
		f, _ := strconv.ParseFloat(n, 64)
		return f
	}
	return 0
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 0, 64)
}
