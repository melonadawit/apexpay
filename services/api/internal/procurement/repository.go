package procurement

import (
	"context"
	"fmt"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// ---- Vendors ----

func (r *Repository) CreateVendor(ctx context.Context, merchantID string, in VendorInput) (*Vendor, error) {
	if in.PaymentTermsDays <= 0 {
		in.PaymentTermsDays = 30
	}
	v := &Vendor{ID: id.New("vend"), Name: in.Name, Email: in.Email, Phone: in.Phone,
		TIN: in.TIN, PaymentTermsDays: in.PaymentTermsDays, Status: "active"}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO vendors (id, merchant_id, name, email, phone, tin, payment_terms_days, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,'active')`,
		v.ID, merchantID, v.Name, v.Email, v.Phone, v.TIN, v.PaymentTermsDays)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (r *Repository) ListVendors(ctx context.Context, merchantID string) ([]Vendor, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, COALESCE(email,''), COALESCE(phone,''), COALESCE(tin,''), payment_terms_days, status
		FROM vendors WHERE merchant_id=$1 ORDER BY name`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Vendor{}
	for rows.Next() {
		var v Vendor
		if err := rows.Scan(&v.ID, &v.Name, &v.Email, &v.Phone, &v.TIN, &v.PaymentTermsDays, &v.Status); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

// ---- Purchase Orders ----

func (r *Repository) CreatePO(ctx context.Context, merchantID, userID string, in POInput) (*PurchaseOrder, error) {
	if len(in.Items) == 0 {
		return nil, fmt.Errorf("po requires items")
	}
	po := &PurchaseOrder{ID: id.New("po"), VendorID: in.VendorID, PONumber: in.PONumber,
		OrderDate: in.OrderDate, ExpectedDelivery: in.ExpectedDelivery, Status: "approved"}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var subtotal, tax decimal.Decimal
	items := make([]POItem, 0, len(in.Items))
	for _, it := range in.Items {
		qty, err1 := decimal.NewFromString(it.Quantity)
		price, err2 := decimal.NewFromString(it.UnitPrice)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid quantity or price")
		}
		line := qty.Mul(price)
		subtotal = subtotal.Add(line)
		items = append(items, POItem{ItemName: it.ItemName, Quantity: qty.String(), UnitPrice: price.String(), LineTotal: line.String()})
	}
	taxRate := decimal.Zero
	if in.TaxRate != "" {
		taxRate, _ = decimal.NewFromString(in.TaxRate)
	}
	tax = subtotal.Mul(taxRate).Round(2)
	total := subtotal.Add(tax)

	po.Subtotal, po.TaxAmount, po.TotalAmount = subtotal.String(), tax.String(), total.String()

	if po.OrderDate == "" {
		po.OrderDate = timeNow()
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO ap_purchase_orders (id, merchant_id, vendor_id, po_number, order_date, expected_delivery, status, subtotal, tax_amount, total_amount, created_by)
		VALUES ($1,$2,$3,$4,$5::date,$6::date,$7,$8,$9,$10,$11)`,
		po.ID, merchantID, in.VendorID, in.PONumber, po.OrderDate, nilDate(in.ExpectedDelivery), po.Status,
		po.Subtotal, po.TaxAmount, po.TotalAmount, userID); err != nil {
		return nil, err
	}
	for i, it := range in.Items {
		lineID := po.ID + "_" + fmt.Sprint(i+1)
		if _, err := tx.Exec(ctx, `
			INSERT INTO ap_purchase_order_items (id, purchase_order_id, item_name, quantity, unit_price, line_total, received_qty)
			VALUES ($1,$2,$3,$4::numeric,$5::numeric,$6::numeric,0)`,
			lineID, po.ID, it.ItemName, it.Quantity, it.UnitPrice, items[i].LineTotal); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	po.Items = items
	return po, nil
}

func (r *Repository) ListPOs(ctx context.Context, merchantID string) ([]PurchaseOrder, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT po.id, po.vendor_id, v.name, po.po_number, to_char(po.order_date,'YYYY-MM-DD'),
		       COALESCE(to_char(po.expected_delivery,'YYYY-MM-DD'),''), po.status,
		       po.subtotal::text, po.tax_amount::text, po.total_amount::text
		FROM ap_purchase_orders po JOIN vendors v ON v.id = po.vendor_id
		WHERE po.merchant_id=$1 ORDER BY po.created_at DESC LIMIT 50`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []PurchaseOrder{}
	for rows.Next() {
		var p PurchaseOrder
		if err := rows.Scan(&p.ID, &p.VendorID, &p.VendorName, &p.PONumber, &p.OrderDate, &p.ExpectedDelivery, &p.Status, &p.Subtotal, &p.TaxAmount, &p.TotalAmount); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Repository) GetPO(ctx context.Context, merchantID, poID string) (*PurchaseOrder, error) {
	var p PurchaseOrder
	err := r.pool.QueryRow(ctx, `
		SELECT po.id, po.vendor_id, v.name, po.po_number, to_char(po.order_date,'YYYY-MM-DD'),
		       COALESCE(to_char(po.expected_delivery,'YYYY-MM-DD'),''), po.status,
		       po.subtotal::text, po.tax_amount::text, po.total_amount::text
		FROM ap_purchase_orders po JOIN vendors v ON v.id = po.vendor_id
		WHERE po.merchant_id=$1 AND po.id=$2`, merchantID, poID).
		Scan(&p.ID, &p.VendorID, &p.VendorName, &p.PONumber, &p.OrderDate, &p.ExpectedDelivery, &p.Status, &p.Subtotal, &p.TaxAmount, &p.TotalAmount)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `
		SELECT item_name, quantity::text, unit_price::text, line_total::text, received_qty::text
		FROM ap_purchase_order_items WHERE purchase_order_id=$1`, poID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var it POItem
		if err := rows.Scan(&it.ItemName, &it.Quantity, &it.UnitPrice, &it.LineTotal, &it.ReceivedQty); err != nil {
			return nil, err
		}
		p.Items = append(p.Items, it)
	}
	return &p, rows.Err()
}

// Receive creates a goods-received receipt against a PO and marks it received.
func (r *Repository) Receive(ctx context.Context, merchantID, userID, poID string) (*Receipt, error) {
	po, err := r.GetPO(ctx, merchantID, poID)
	if err != nil {
		return nil, err
	}
	if po == nil {
		return nil, fmt.Errorf("po not found")
	}
	rec := &Receipt{ID: id.New("rcpt"), VendorID: po.VendorID, PONumber: po.PONumber, ReceiptNumber: "RCPT-" + po.PONumber}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `
		INSERT INTO ap_receipts (id, merchant_id, purchase_order_id, vendor_id, receipt_number, received_at, created_by)
		VALUES ($1,$2,$3,$4,$5,now(),$6)`, rec.ID, merchantID, poID, po.VendorID, rec.ReceiptNumber, userID); err != nil {
		return nil, err
	}
	// Copy PO items into ap_receipt_lines and mark received.
	for _, it := range po.Items {
		if _, err := tx.Exec(ctx, `
			INSERT INTO ap_receipt_lines (id, receipt_id, item_name, quantity, unit_price)
			VALUES ($1,$2,$3,$4::numeric,$5::numeric)`, rec.ID+"_"+it.ItemName, rec.ID, it.ItemName, it.Quantity, it.UnitPrice); err != nil {
			return nil, err
		}
	}
	if _, err := tx.Exec(ctx, `UPDATE ap_purchase_orders SET status='received' WHERE id=$1 AND merchant_id=$2`, poID, merchantID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	rec.VendorName = po.VendorName
	return rec, nil
}

// ---- AP Invoices ----

func (r *Repository) CreateAPInvoice(ctx context.Context, merchantID, userID string, in APInvoiceInput) (*APInvoice, error) {
	sub, err := decimal.NewFromString(in.Subtotal)
	if err != nil {
		return nil, fmt.Errorf("invalid subtotal")
	}
	tax, err := decimal.NewFromString(in.TaxAmount)
	if err != nil {
		return nil, fmt.Errorf("invalid tax")
	}
	total := sub.Add(tax)

	matchStatus := "unmatched"
	status := "pending"
	if in.PurchaseOrderID != "" {
		po, err := r.GetPO(ctx, merchantID, in.PurchaseOrderID)
		if err != nil {
			return nil, err
		}
		if po == nil {
			return nil, fmt.Errorf("po not found")
		}
		// 3-way-ish match: invoice total vs PO total.
		poTotal, _ := decimal.NewFromString(po.TotalAmount)
		if total.Equal(poTotal) {
			matchStatus = "matched"
			status = "matched"
		} else {
			matchStatus = "mismatch"
		}
	}

	inv := &APInvoice{ID: id.New("apinv"), VendorID: in.VendorID, PurchaseOrderID: in.PurchaseOrderID,
		InvoiceNumber: in.InvoiceNumber, InvoiceDate: in.InvoiceDate, DueDate: in.DueDate,
		Subtotal: sub.String(), TaxAmount: tax.String(), TotalAmount: total.String(),
		AmountPaid: "0", Status: status, MatchStatus: matchStatus}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO ap_invoices (id, merchant_id, vendor_id, purchase_order_id, invoice_number, invoice_date, due_date, subtotal, tax_amount, total_amount, amount_paid, status, match_status, created_by)
		VALUES ($1,$2,$3,$4,$5,$6::date,$7::date,$8::numeric,$9::numeric,$10::numeric,0,$11,$12,$13)`,
		inv.ID, merchantID, in.VendorID, in.PurchaseOrderID, in.InvoiceNumber, in.InvoiceDate, in.DueDate,
		inv.Subtotal, inv.TaxAmount, inv.TotalAmount, inv.Status, inv.MatchStatus, userID)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *Repository) ListAPInvoices(ctx context.Context, merchantID string) ([]APInvoice, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT a.id, a.vendor_id, v.name, COALESCE(a.purchase_order_id,''), a.invoice_number,
		       to_char(a.invoice_date,'YYYY-MM-DD'), to_char(a.due_date,'YYYY-MM-DD'),
		       a.subtotal::text, a.tax_amount::text, a.total_amount::text, a.amount_paid::text,
		       a.status, a.match_status
		FROM ap_invoices a JOIN vendors v ON v.id = a.vendor_id
		WHERE a.merchant_id=$1 ORDER BY a.due_date ASC LIMIT 100`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []APInvoice{}
	for rows.Next() {
		var a APInvoice
		if err := rows.Scan(&a.ID, &a.VendorID, &a.VendorName, &a.PurchaseOrderID, &a.InvoiceNumber, &a.InvoiceDate, &a.DueDate, &a.Subtotal, &a.TaxAmount, &a.TotalAmount, &a.AmountPaid, &a.Status, &a.MatchStatus); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Aging computes AP aging buckets from unpaid invoices.
func (r *Repository) Aging(ctx context.Context, merchantID string) ([]AgingBucket, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT CASE
			WHEN due_date >= CURRENT_DATE THEN 'current'
			WHEN due_date >= CURRENT_DATE - 30 THEN '30'
			WHEN due_date >= CURRENT_DATE - 60 THEN '60'
			WHEN due_date >= CURRENT_DATE - 90 THEN '90'
			ELSE '90plus' END AS bucket,
			COUNT(*), COALESCE(SUM(total_amount - amount_paid),0)::text
		FROM ap_invoices
		WHERE merchant_id=$1 AND status NOT IN ('paid','cancelled')
		GROUP BY bucket ORDER BY bucket`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AgingBucket{}
	order := map[string]int{"current": 0, "30": 1, "60": 2, "90": 3, "90plus": 4}
	for rows.Next() {
		var b AgingBucket
		if err := rows.Scan(&b.Bucket, &b.Count, &b.Amount); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	// Sort buckets current -> 90plus for a stable response.
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if order[out[j].Bucket] < order[out[i].Bucket] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out, rows.Err()
}

func nilDate(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func timeNow() string {
	return time.Now().UTC().Format("2006-01-02")
}
