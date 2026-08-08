package lending

// Lending / micro-loan domain.

type Loan struct {
	ID           string `json:"id"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Purpose      string `json:"purpose"`
	Status       string `json:"status"`
	InterestRate string `json:"interest_rate"`
	DueDate      string `json:"due_date"`
	RepaidAmount string `json:"repaid_amount"`
	Outstanding  string `json:"outstanding_amount"`
	CreatedAt    string `json:"created_at"`
}
