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

// JournalLine is one debit/credit leg of a manual journal entry.
type JournalLine struct {
	AccountCode string `json:"account_code"`
	AccountName string `json:"account_name,omitempty"`
	Direction   string `json:"direction"` // debit | credit
	Amount      string `json:"amount"`
}

// JournalEntryRequest is the input to create a manual journal entry.
type JournalEntryRequest struct {
	Memo          string        `json:"memo"`
	ReferenceType string        `json:"reference_type,omitempty"`
	ReferenceID   string        `json:"reference_id,omitempty"`
	Lines         []JournalLine `json:"lines"`
}

// JournalEntry is a persisted manual journal entry with its legs.
type JournalEntry struct {
	ID         string        `json:"id"`
	Memo       string        `json:"memo"`
	Period     string        `json:"period"` // YYYY-MM
	PostingKey string        `json:"posting_key,omitempty"`
	Lines      []JournalLine `json:"lines"`
	CreatedAt  string        `json:"created_at"`
}

// FiscalPeriod tracks open/closed status of an accounting period (month) so postings can
// be locked once books are closed.
type FiscalPeriod struct {
	ID       string `json:"id"`
	Merchant string `json:"merchant_id"`
	Period   string `json:"period"` // YYYY-MM
	Status   string `json:"status"` // open | closed
	ClosedAt string `json:"closed_at,omitempty"`
	ClosedBy string `json:"closed_by,omitempty"`
}
