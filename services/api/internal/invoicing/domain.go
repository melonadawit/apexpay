package invoicing

// Invoicing & Receivables domain types.

type LineItem struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
	Quantity    string `json:"quantity"`
	UnitPrice   string `json:"unit_price"`
	LineTotal   string `json:"line_total"`
	SortOrder   int    `json:"sort_order,omitempty"`
}

type Invoice struct {
	ID                string     `json:"id"`
	InvoiceNumber     string     `json:"invoice_number"`
	CustomerName      string     `json:"customer_name"`
	CustomerEmail     string     `json:"customer_email,omitempty"`
	CustomerPhone     string     `json:"customer_phone,omitempty"`
	IssueDate         string     `json:"issue_date"`
	DueDate           string     `json:"due_date"`
	Currency          string     `json:"currency"`
	Subtotal          string     `json:"subtotal"`
	TaxAmount         string     `json:"tax_amount"`
	WithholdingAmount string     `json:"withholding_amount"`
	TotalAmount       string     `json:"total_amount"`
	AmountPaid        string     `json:"amount_paid"`
	Status            string     `json:"status"`
	HostedToken       string     `json:"hosted_token,omitempty"`
	DunningStage      int        `json:"dunning_stage"`
	Notes             string     `json:"notes,omitempty"`
	LineItems         []LineItem `json:"line_items"`
	CreatedAt         string     `json:"created_at"`
}

// AgingBucket counts/sums invoices by age bucket for AR aging.
type AgingBucket struct {
	Bucket string `json:"bucket"` // current, 30, 60, 90, 90plus
	Count  int    `json:"count"`
	Amount string `json:"amount"`
}

type CreateInvoiceRequest struct {
	InvoiceNumber      string     `json:"invoice_number"`
	CustomerName       string     `json:"customer_name"`
	CustomerEmail      string     `json:"customer_email"`
	CustomerPhone      string     `json:"customer_phone"`
	IssueDate          string     `json:"issue_date"`
	DueDate            string     `json:"due_date"`
	Currency           string     `json:"currency"`
	TaxPercent         string     `json:"tax_percent"`         // VAT 15% / TOT
	WithholdingPercent string     `json:"withholding_percent"` // withholding 2%
	LineItems          []LineItem `json:"line_items"`
	Notes              string     `json:"notes"`
}
