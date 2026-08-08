package tax

// Tax domain types: schedules (VAT/TOT/withholding returns) and GL posting of collected tax.

// TaxType is the kind of indirect tax collected.
type TaxType string

const (
	VAT         TaxType = "vat"         // 15% output VAT
	TOT         TaxType = "tot"         // Turnover Tax
	Withholding TaxType = "withholding" // income withholding tax deducted
)

// ScheduleLine is one period's tax position for a given tax type.
type ScheduleLine struct {
	Period    string `json:"period"`    // YYYY-MM
	TaxType   string `json:"tax_type"`  // vat | tot | withholding
	Collected string `json:"collected"` // gross tax billed on invoices this period
	Paid      string `json:"paid"`      // tax already remitted to the authority
	Due       string `json:"due"`       // outstanding (collected - paid)
	Count     int    `json:"count"`     // number of invoices contributing
}

// Schedule is the full tax return schedule for a merchant.
type Schedule struct {
	MerchantID string         `json:"merchant_id"`
	Lines      []ScheduleLine `json:"lines"`
	TotalDue   string         `json:"total_due"`
}

// PostResult describes a GL posting of collected tax into a liability account.
type PostResult struct {
	JournalID string `json:"journal_id"`
	Period    string `json:"period"`
	Amount    string `json:"amount"`
	Account   string `json:"account"`
}
