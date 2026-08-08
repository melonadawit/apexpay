package accounting

// Accounting & Bookkeeping domain types.

// Account maps a ledger account to a financial statement classification so reports can be
// derived. Merchant accounts are auto-classified by code prefix; this table lets merchants
// override.
type Account struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Category   string `json:"category"`    // asset, liability, equity, revenue, expense
	NormalSide string `json:"normal_side"` // debit | credit
	IsActive   bool   `json:"is_active"`
}

// TrialBalanceRow is one account's balance in the trial balance.
type TrialBalanceRow struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Debit  string `json:"debit"`
	Credit string `json:"credit"`
}

// StatementLine is a grouped line in a financial statement.
type StatementLine struct {
	Label  string `json:"label"`
	Amount string `json:"amount"`
	Kind   string `json:"kind"` // header, total, normal
}

// FinancialStatement is a P&L or Balance Sheet built from the ledger.
type FinancialStatement struct {
	Title  string          `json:"title"`
	Period string          `json:"period"`
	Lines  []StatementLine `json:"lines"`
}

// CashFlowLine is one line of the cash-flow statement.
type CashFlowLine struct {
	Label  string `json:"label"`
	Amount string `json:"amount"`
	Kind   string `json:"kind"`
}
