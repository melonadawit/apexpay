package accounting

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Repository struct {
	pool  *pgxpool.Pool
	cache *Cache
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// SetCache attaches a Redis read-through cache for report queries.
func (r *Repository) SetCache(c *Cache) { r.cache = c }

// AccountCategories maps a ledger account code prefix to a statement category.
func categoryForCode(code string) string {
	c := ""
	for _, r := range code {
		if r == ':' || r == '-' {
			break
		}
		c += string(r)
	}
	switch c {
	case "asset", "liability", "equity", "revenue", "expense":
		return c
	default:
		// Heuristic by first char of code.
		if len(code) > 0 {
			switch code[0] {
			case '1':
				return "asset"
			case '2':
				return "liability"
			case '3':
				return "equity"
			case '4':
				return "revenue"
			case '5':
				return "expense"
			}
		}
		return "asset"
	}
}

// ChartOfAccounts lists all ledger accounts for a merchant's books with classification.
func (r *Repository) ChartOfAccounts(ctx context.Context, merchantID string) ([]Account, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT a.code, a.name, a.normal_balance
		FROM ledger_accounts a
		JOIN ledger_books b ON b.id = a.book_id
		WHERE b.merchant_id=$1
		ORDER BY a.code`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Account{}
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.Code, &a.Name, &a.NormalSide); err != nil {
			return nil, err
		}
		a.Category = categoryForCode(a.Code)
		a.IsActive = true
		list = append(list, a)
	}
	return list, rows.Err()
}

// TrialBalance returns each account's net debit/credit balance across all the merchant's books.
func (r *Repository) TrialBalance(ctx context.Context, merchantID string) ([]TrialBalanceRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT a.code, a.name,
		       COALESCE(SUM(e.amount) FILTER (WHERE e.direction='debit'),0)::text,
		       COALESCE(SUM(e.amount) FILTER (WHERE e.direction='credit'),0)::text
		FROM ledger_entries e
		JOIN ledger_books b ON b.id = e.book_id
		JOIN ledger_accounts a ON a.id = e.account_id
		WHERE b.merchant_id=$1
		GROUP BY a.code, a.name
		ORDER BY a.code`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []TrialBalanceRow{}
	for rows.Next() {
		var r TrialBalanceRow
		if err := rows.Scan(&r.Code, &r.Name, &r.Debit, &r.Credit); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, rows.Err()
}

// balanceOf returns the net balance (debit-credit) for a set of account codes.
func (r *Repository) balanceOf(ctx context.Context, merchantID string, codes []string) (decimal.Decimal, error) {
	sql := `
		SELECT COALESCE(SUM(CASE WHEN e.direction='debit' THEN e.amount ELSE -e.amount END),0)::text
		FROM ledger_entries e
		JOIN ledger_books b ON b.id = e.book_id
		JOIN ledger_accounts a ON a.id = e.account_id
		WHERE b.merchant_id=$1 AND (` + codesOR(codes) + `)`
	rows, err := r.pool.Query(ctx, sql, merchantID)
	if err != nil {
		return decimal.Zero, err
	}
	defer rows.Close()
	var s string
	if rows.Next() {
		if err := rows.Scan(&s); err != nil {
			return decimal.Zero, err
		}
	}
	return decimal.NewFromString(s)
}

func codesOR(codes []string) string {
	out := ""
	for i, c := range codes {
		if i > 0 {
			out += " OR "
		}
		out += "a.code LIKE '" + c + "%'"
	}
	return out
}

// ProfitLoss builds a P&L for the period. Revenue minus expenses. Cached briefly.
func (r *Repository) ProfitLoss(ctx context.Context, merchantID, from, to string) (*FinancialStatement, error) {
	key := "acct:pnl:" + merchantID + ":" + from + ":" + to
	val, err := r.cacheGet(ctx, key, func() (any, error) { return r.profitLossUncached(ctx, merchantID, from, to) })
	if err != nil {
		return nil, err
	}
	return val.(*FinancialStatement), nil
}

func (r *Repository) profitLossUncached(ctx context.Context, merchantID, from, to string) (*FinancialStatement, error) {
	rev, err := r.balanceOfPeriod(ctx, merchantID, []string{"revenue", "4"}, from, to)
	if err != nil {
		return nil, err
	}
	exp, err := r.balanceOfPeriod(ctx, merchantID, []string{"expense", "5"}, from, to)
	if err != nil {
		return nil, err
	}
	net := rev.Sub(exp)
	lines := []StatementLine{
		{Label: "Revenue", Amount: rev.String(), Kind: "normal"},
		{Label: "Expenses", Amount: exp.String(), Kind: "normal"},
		{Label: "Net Profit / (Loss)", Amount: net.String(), Kind: "total"},
	}
	return &FinancialStatement{Title: "Profit & Loss", Period: from + " → " + to, Lines: lines}, nil
}

// BalanceSheet builds assets vs liabilities+equity.
func (r *Repository) BalanceSheet(ctx context.Context, merchantID, asOf string) (*FinancialStatement, error) {
	assets, err := r.balanceOf(ctx, merchantID, []string{"asset", "1"})
	if err != nil {
		return nil, err
	}
	liab, err := r.balanceOf(ctx, merchantID, []string{"liability", "2"})
	if err != nil {
		return nil, err
	}
	equity, err := r.balanceOf(ctx, merchantID, []string{"equity", "3"})
	if err != nil {
		return nil, err
	}
	lines := []StatementLine{
		{Label: "Assets", Amount: assets.String(), Kind: "header"},
		{Label: "Total Assets", Amount: assets.String(), Kind: "total"},
		{Label: "Liabilities", Amount: liab.String(), Kind: "header"},
		{Label: "Equity", Amount: equity.String(), Kind: "header"},
		{Label: "Total Liabilities + Equity", Amount: liab.Add(equity).String(), Kind: "total"},
	}
	return &FinancialStatement{Title: "Balance Sheet", Period: "as of " + asOf, Lines: lines}, nil
}

// CashFlow summarizes net cash movement by category.
func (r *Repository) CashFlow(ctx context.Context, merchantID, from, to string) ([]CashFlowLine, error) {
	// Operating = revenue - expense on cash accounts; investing/financing are inferred from
	// account categories. For a payment gateway, operating cash flow is the merchant's net
	// cleared settlements.
	operating, err := r.balanceOfPeriod(ctx, merchantID, []string{"asset:clearing", "revenue", "expense", "liability:merchant_payable"}, from, to)
	if err != nil {
		return nil, err
	}
	return []CashFlowLine{
		{Label: "Net Operating Cash Flow", Amount: operating.String(), Kind: "normal"},
	}, nil
}

// balanceOfPeriod sums net balance within a date window.
func (r *Repository) balanceOfPeriod(ctx context.Context, merchantID string, codes []string, from, to string) (decimal.Decimal, error) {
	sql := `
		SELECT COALESCE(SUM(CASE WHEN e.direction='debit' THEN e.amount ELSE -e.amount END),0)::text
		FROM ledger_entries e
		JOIN ledger_books b ON b.id = e.book_id
		JOIN ledger_accounts a ON a.id = e.account_id
		JOIN ledger_journals j ON j.id = e.journal_id
		WHERE b.merchant_id=$1 AND j.created_at BETWEEN $2::date AND $3::date AND (` + codesOR(codes) + `)`
	rows, err := r.pool.Query(ctx, sql, merchantID, from, to)
	if err != nil {
		return decimal.Zero, err
	}
	defer rows.Close()
	var s string
	if rows.Next() {
		if err := rows.Scan(&s); err != nil {
			return decimal.Zero, err
		}
	}
	return decimal.NewFromString(s)
}

// cacheGet runs compute through the optional Redis cache.
func (r *Repository) cacheGet(ctx context.Context, key string, compute func() (any, error)) (any, error) {
	if r.cache == nil {
		return compute()
	}
	return r.cache.GetOrCompute(ctx, key, compute)
}
