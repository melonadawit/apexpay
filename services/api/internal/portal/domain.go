package portal

// Self-service portal domain types (vendor + customer).

// PortalType identifies which party is using a portal.
type PortalType string

const (
	PortalVendor   PortalType = "vendor"
	PortalCustomer PortalType = "customer"
)

// Access is a token-gated portal session.
type Access struct {
	ID         string     `json:"id"`
	MerchantID string     `json:"merchant_id"`
	PortalType PortalType `json:"portal_type"`
	EntityID   string     `json:"entity_id"`
	EntityName string     `json:"entity_name"`
	Token      string     `json:"token,omitempty"` // only returned at creation
	ExpiresAt  string     `json:"expires_at"`
}

// VendorInvoice is what a vendor sees in their portal: their AP invoices + payment status.
type VendorInvoice struct {
	InvoiceNumber string `json:"invoice_number"`
	InvoiceDate   string `json:"invoice_date"`
	DueDate       string `json:"due_date"`
	Subtotal      string `json:"subtotal"`
	TaxAmount     string `json:"tax_amount"`
	TotalAmount   string `json:"total_amount"`
	AmountPaid    string `json:"amount_paid"`
	Status        string `json:"status"`
}

// CustomerInvoice is what a customer sees in their portal: their invoices + payment status.
type CustomerInvoice struct {
	InvoiceNumber string `json:"invoice_number"`
	CustomerName  string `json:"customer_name"`
	IssueDate     string `json:"issue_date"`
	DueDate       string `json:"due_date"`
	Subtotal      string `json:"subtotal"`
	TaxAmount     string `json:"tax_amount"`
	TotalAmount   string `json:"total_amount"`
	AmountPaid    string `json:"amount_paid"`
	Status        string `json:"status"`
}

// Dashboard is the portal home for a party.
type Dashboard struct {
	PortalType PortalType `json:"portal_type"`
	EntityName string     `json:"entity_name"`
	Invoices   []any      `json:"invoices"`
}
