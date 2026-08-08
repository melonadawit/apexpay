package procurement

// Procurement & Accounts Payable domain types.

type Vendor struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email,omitempty"`
	Phone            string `json:"phone,omitempty"`
	TIN              string `json:"tin,omitempty"`
	PaymentTermsDays int    `json:"payment_terms_days"`
	Status           string `json:"status"`
}

type VendorInput struct {
	Name             string `json:"name"`
	Email            string `json:"email,omitempty"`
	Phone            string `json:"phone,omitempty"`
	TIN              string `json:"tin,omitempty"`
	PaymentTermsDays int    `json:"payment_terms_days"`
}

type POItem struct {
	ItemName    string `json:"item_name"`
	Quantity    string `json:"quantity"`
	UnitPrice   string `json:"unit_price"`
	LineTotal   string `json:"line_total"`
	ReceivedQty string `json:"received_qty,omitempty"`
}

type PurchaseOrder struct {
	ID               string   `json:"id"`
	VendorID         string   `json:"vendor_id"`
	VendorName       string   `json:"vendor_name,omitempty"`
	PONumber         string   `json:"po_number"`
	OrderDate        string   `json:"order_date"`
	ExpectedDelivery string   `json:"expected_delivery,omitempty"`
	Status           string   `json:"status"`
	Subtotal         string   `json:"subtotal"`
	TaxAmount        string   `json:"tax_amount"`
	TotalAmount      string   `json:"total_amount"`
	Items            []POItem `json:"items"`
	CreatedAt        string   `json:"created_at"`
}

type POInput struct {
	VendorID         string        `json:"vendor_id"`
	PONumber         string        `json:"po_number"`
	OrderDate        string        `json:"order_date"`
	ExpectedDelivery string        `json:"expected_delivery,omitempty"`
	Items            []POItemInput `json:"items"`
	TaxRate          string        `json:"tax_rate,omitempty"` // e.g. "0.15"
}

type POItemInput struct {
	ItemName  string `json:"item_name"`
	Quantity  string `json:"quantity"`
	UnitPrice string `json:"unit_price"`
}

type Receipt struct {
	ID            string `json:"id"`
	PONumber      string `json:"po_number,omitempty"`
	VendorID      string `json:"vendor_id"`
	VendorName    string `json:"vendor_name,omitempty"`
	ReceiptNumber string `json:"receipt_number"`
	Note          string `json:"note,omitempty"`
	ReceivedAt    string `json:"received_at"`
}

type APInvoice struct {
	ID              string `json:"id"`
	VendorID        string `json:"vendor_id"`
	VendorName      string `json:"vendor_name,omitempty"`
	PurchaseOrderID string `json:"purchase_order_id,omitempty"`
	InvoiceNumber   string `json:"invoice_number"`
	InvoiceDate     string `json:"invoice_date"`
	DueDate         string `json:"due_date"`
	Subtotal        string `json:"subtotal"`
	TaxAmount       string `json:"tax_amount"`
	TotalAmount     string `json:"total_amount"`
	AmountPaid      string `json:"amount_paid"`
	Status          string `json:"status"`
	MatchStatus     string `json:"match_status"`
}

type APInvoiceInput struct {
	VendorID        string `json:"vendor_id"`
	PurchaseOrderID string `json:"purchase_order_id,omitempty"`
	InvoiceNumber   string `json:"invoice_number"`
	InvoiceDate     string `json:"invoice_date"`
	DueDate         string `json:"due_date"`
	Subtotal        string `json:"subtotal"`
	TaxAmount       string `json:"tax_amount"`
}

// AgingBucket is an AP aging bucket (current, 30, 60, 90, 90plus).
type AgingBucket struct {
	Bucket string `json:"bucket"`
	Count  int    `json:"count"`
	Amount string `json:"amount"`
}
