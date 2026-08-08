package fxreval

// FX revaluation domain types.

// RevaluationLine is one foreign-currency account's revaluation for a period.
type RevaluationLine struct {
	AccountID     string `json:"account_id"`
	AccountNumber string `json:"account_number,omitempty"`
	Currency      string `json:"currency"`
	AmountFX      string `json:"amount_fx"`
	Rate          string `json:"rate"`         // ETB per unit of foreign currency
	AmountETB     string `json:"amount_etb"`   // revalued ETB equivalent
	FXGainLoss    string `json:"fx_gain_loss"` // signed
}

// Revaluation is the result of revaluing all foreign-currency accounts for a period.
type Revaluation struct {
	Period    string            `json:"period"`
	Lines     []RevaluationLine `json:"lines"`
	NetFX     string            `json:"net_fx"` // net unrealized FX gain/loss (ETB)
	JournalID string            `json:"journal_id,omitempty"`
}
