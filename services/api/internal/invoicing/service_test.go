package invoicing

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestComputeInvoiceTotals_WithVATAndWithholding(t *testing.T) {
	items := []LineItem{
		{Description: "Service A", Quantity: "2", UnitPrice: "1000"}, // 2000
		{Description: "Service B", Quantity: "1", UnitPrice: "500"},  // 500
	}
	subtotal, tax, withholding, total, lines := ComputeInvoiceTotals(items,
		decimal.NewFromInt(15), // VAT 15%
		decimal.NewFromInt(2),  // withholding 2%
	)
	if !subtotal.Equal(decimal.NewFromInt(2500)) {
		t.Fatalf("subtotal = %s, want 2500", subtotal)
	}
	if !tax.Equal(decimal.NewFromFloat(375)) { // 2500 * 0.15
		t.Fatalf("tax = %s, want 375", tax)
	}
	if !withholding.Equal(decimal.NewFromInt(50)) { // 2500 * 0.02
		t.Fatalf("withholding = %s, want 50", withholding)
	}
	// total = 2500 + 375 - 50 = 2825
	if !total.Equal(decimal.NewFromInt(2825)) {
		t.Fatalf("total = %s, want 2825", total)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 line items, got %d", len(lines))
	}
	if lines[0].LineTotal != "2000" {
		t.Fatalf("line 1 total = %s, want 2000", lines[0].LineTotal)
	}
}

func TestComputeInvoiceTotals_NoTax(t *testing.T) {
	items := []LineItem{{Description: "x", Quantity: "3", UnitPrice: "100"}}
	subtotal, _, _, total, _ := ComputeInvoiceTotals(items, decimal.Zero, decimal.Zero)
	if !subtotal.Equal(decimal.NewFromInt(300)) || !total.Equal(decimal.NewFromInt(300)) {
		t.Fatalf("subtotal=%s total=%s want 300/300", subtotal, total)
	}
}

func TestComputeInvoiceTotals_FractionalQty(t *testing.T) {
	items := []LineItem{{Description: "half day", Quantity: "0.5", UnitPrice: "1000"}}
	subtotal, _, _, total, _ := ComputeInvoiceTotals(items, decimal.Zero, decimal.Zero)
	if !subtotal.Equal(decimal.NewFromInt(500)) {
		t.Fatalf("subtotal = %s, want 500", subtotal)
	}
	_ = total
}
