package invoicing

import (
	"context"

	"github.com/shopspring/decimal"
)

// Service applies invoice math (subtotal, VAT/TOT, withholding) and business rules.
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

// ComputeInvoiceTotals is a pure function that derives line totals, subtotal, tax,
// withholding, and final total from line items + percentages. Unit-testable without a DB.
func ComputeInvoiceTotals(items []LineItem, taxPercent, withholdPercent decimal.Decimal) (subtotal, tax, withholding, total decimal.Decimal, lines []LineItem) {
	subtotal = decimal.Zero
	lines = make([]LineItem, 0, len(items))
	for _, li := range items {
		qty, _ := decimal.NewFromString(li.Quantity)
		price, _ := decimal.NewFromString(li.UnitPrice)
		lineTotal := qty.Mul(price).Round(2)
		lines = append(lines, LineItem{
			Description: li.Description, Quantity: qty.String(), UnitPrice: price.String(), LineTotal: lineTotal.String(),
		})
		subtotal = subtotal.Add(lineTotal)
	}
	tax = subtotal.Mul(taxPercent).Div(decimal.NewFromInt(100)).Round(2)
	withholding = subtotal.Mul(withholdPercent).Div(decimal.NewFromInt(100)).Round(2)
	total = subtotal.Add(tax).Sub(withholding)
	if total.LessThan(decimal.Zero) {
		total = decimal.Zero
	}
	return
}

// BuildInvoice computes line totals, subtotal, tax (VAT/TOT percent), withholding, and
// the final total, then persists the invoice. All money is decimal — no floats.
func (s *Service) BuildInvoice(ctx context.Context, merchantID, userID string, req CreateInvoiceRequest) (*Invoice, error) {
	taxPct, _ := decimal.NewFromString(req.TaxPercent)
	withholdPct, _ := decimal.NewFromString(req.WithholdingPercent)
	subtotal, tax, withholding, total, lines := ComputeInvoiceTotals(req.LineItems, taxPct, withholdPct)

	inv := &Invoice{
		InvoiceNumber: req.InvoiceNumber, CustomerName: req.CustomerName,
		CustomerEmail: req.CustomerEmail, CustomerPhone: req.CustomerPhone,
		IssueDate: req.IssueDate, DueDate: req.DueDate, Currency: req.Currency,
		Subtotal: subtotal.String(), TaxAmount: tax.String(), WithholdingAmount: withholding.String(),
		TotalAmount: total.String(), Status: "draft", LineItems: lines,
	}
	if err := s.repo.CreateInvoice(ctx, merchantID, userID, inv); err != nil {
		return nil, err
	}
	return inv, nil
}

// Issue attaches a hosted payment token to a draft invoice, moving it to 'sent'.
func (s *Service) Issue(ctx context.Context, merchantID, invoiceID, token string) (*Invoice, error) {
	if err := s.repo.SetHostedToken(ctx, merchantID, invoiceID, token); err != nil {
		return nil, err
	}
	return s.repo.GetInvoice(ctx, merchantID, invoiceID)
}
