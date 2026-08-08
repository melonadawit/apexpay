package inventory

// Inventory & Sales (software POS) domain types.

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	SKU         string `json:"sku,omitempty"`
	Price       string `json:"price"`
	CostPrice   string `json:"cost_price"`
	Currency    string `json:"currency"`
	VATCategory string `json:"vat_category"`
	StockQty    string `json:"stock_qty"`
	LowStock    string `json:"low_stock_threshold"`
	Status      string `json:"status"`
}

type OrderItem struct {
	ProductID   string `json:"product_id"`
	Description string `json:"description"`
	Quantity    string `json:"quantity"`
	UnitPrice   string `json:"unit_price"`
	LineTotal   string `json:"line_total"`
}

type Order struct {
	ID            string      `json:"id"`
	OrderNumber   string      `json:"order_number"`
	CustomerName  string      `json:"customer_name,omitempty"`
	CustomerEmail string      `json:"customer_email,omitempty"`
	Status        string      `json:"status"`
	Subtotal      string      `json:"subtotal"`
	TaxAmount     string      `json:"tax_amount"`
	TotalAmount   string      `json:"total_amount"`
	Items         []OrderItem `json:"items"`
	CreatedAt     string      `json:"created_at"`
}

type StockMovement struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Qty       string `json:"qty"`
	Direction string `json:"direction"`
	Reference string `json:"reference,omitempty"`
	Note      string `json:"note,omitempty"`
	CreatedAt string `json:"created_at"`
}
