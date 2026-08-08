package treasury

// Treasury & Cash Management domain types.

// Transfer moves funds between a merchant's own current accounts.
type Transfer struct {
	ID            string `json:"id"`
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Purpose       string `json:"purpose,omitempty"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// AccountPosition is one account in the merchant's cash position.
type AccountPosition struct {
	AccountID        string `json:"account_id"`
	AccountNumber    string `json:"account_number"`
	AccountName      string `json:"account_name"`
	AccountType      string `json:"account_type"`
	BankCode         string `json:"bank_code"`
	Balance          string `json:"balance"`
	AvailableBalance string `json:"available_balance"`
	Currency         string `json:"currency"`
}

// CashPosition is the aggregate cash position across a merchant's accounts.
type CashPosition struct {
	Accounts       []AccountPosition `json:"accounts"`
	TotalBalance   string            `json:"total_balance"`
	TotalAvailable string            `json:"total_available"`
	Currency       string            `json:"currency"`
	GeneratedAt    string            `json:"generated_at"`
}

// Forecast is a cash-flow forecast over a horizon, with inflow/outflow buckets.
type Forecast struct {
	ID           string `json:"id"`
	ForecastDate string `json:"forecast_date"`
	HorizonDays  int    `json:"horizon_days"`
	InflowToday  string `json:"inflow_today"`
	Inflow30d    string `json:"inflow_30d"`
	Inflow60d    string `json:"inflow_60d"`
	Inflow90d    string `json:"inflow_90d"`
	OutflowToday string `json:"outflow_today"`
	Outflow30d   string `json:"outflow_30d"`
	Outflow60d   string `json:"outflow_60d"`
	Outflow90d   string `json:"outflow_90d"`
	Net90d       string `json:"net_90d"`
	GeneratedAt  string `json:"generated_at"`
}
