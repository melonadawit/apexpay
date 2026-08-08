package invoicing

import (
	"context"
	"errors"
	"strconv"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// CreateInvoice inserts an invoice and its line items in one transaction.
func (r *Repository) CreateInvoice(ctx context.Context, merchantID, userID string, inv *Invoice) error {
	inv.ID = id.New("inv")
	inv.Currency = defaultStr(inv.Currency, "ETB")
	if inv.Status == "" {
		inv.Status = "draft"
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `
		INSERT INTO invoices (id, merchant_id, invoice_number, customer_name, customer_email, customer_phone,
			issue_date, due_date, currency, subtotal, tax_amount, withholding_amount, total_amount, amount_paid, status, notes, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7::date,$8::date,$9,$10::numeric,$11::numeric,$12::numeric,$13::numeric,0,$14,$15,$16)`,
		inv.ID, merchantID, inv.InvoiceNumber, inv.CustomerName, inv.CustomerEmail, inv.CustomerPhone,
		inv.IssueDate, inv.DueDate, inv.Currency, inv.Subtotal, inv.TaxAmount, inv.WithholdingAmount, inv.TotalAmount,
		inv.Status, inv.Notes, userID)
	if err != nil {
		return err
	}
	for i, li := range inv.LineItems {
		_, err = tx.Exec(ctx, `
			INSERT INTO invoice_line_items (id, invoice_id, description, quantity, unit_price, line_total, sort_order)
			VALUES ($1,$2,$3,$4::numeric,$5::numeric,$6::numeric,$7)`,
			id.New("line"), inv.ID, li.Description, li.Quantity, li.UnitPrice, li.LineTotal, i+1)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// ListInvoices returns invoices for a merchant, optionally filtered by status.
func (r *Repository) ListInvoices(ctx context.Context, merchantID, status string, limit int) ([]Invoice, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query := `SELECT id, invoice_number, customer_name, COALESCE(customer_email,''), COALESCE(customer_phone,''),
		to_char(issue_date,'YYYY-MM-DD'), to_char(due_date,'YYYY-MM-DD'), currency,
		subtotal::text, tax_amount::text, withholding_amount::text, total_amount::text, amount_paid::text, status,
		COALESCE(hosted_token,''), dunning_stage,
		to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM invoices WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if status != "" {
		query += ` AND status=$2`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC LIMIT $` + itoa(len(args)+1)
	args = append(args, limit)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Invoice{}
	for rows.Next() {
		var v Invoice
		if err := rows.Scan(&v.ID, &v.InvoiceNumber, &v.CustomerName, &v.CustomerEmail, &v.CustomerPhone,
			&v.IssueDate, &v.DueDate, &v.Currency, &v.Subtotal, &v.TaxAmount, &v.WithholdingAmount,
			&v.TotalAmount, &v.AmountPaid, &v.Status, &v.HostedToken, &v.DunningStage, &v.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

// GetInvoice returns a single invoice with its line items.
func (r *Repository) GetInvoice(ctx context.Context, merchantID, invoiceID string) (*Invoice, error) {
	var v Invoice
	err := r.pool.QueryRow(ctx, `
		SELECT id, invoice_number, customer_name, COALESCE(customer_email,''), COALESCE(customer_phone,''),
			to_char(issue_date,'YYYY-MM-DD'), to_char(due_date,'YYYY-MM-DD'), currency,
			subtotal::text, tax_amount::text, withholding_amount::text, total_amount::text, amount_paid::text, status,
			COALESCE(hosted_token,''), dunning_stage,
			to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM invoices WHERE merchant_id=$1 AND id=$2`, merchantID, invoiceID).
		Scan(&v.ID, &v.InvoiceNumber, &v.CustomerName, &v.CustomerEmail, &v.CustomerPhone,
			&v.IssueDate, &v.DueDate, &v.Currency, &v.Subtotal, &v.TaxAmount, &v.WithholdingAmount,
			&v.TotalAmount, &v.AmountPaid, &v.Status, &v.HostedToken, &v.DunningStage, &v.CreatedAt)
	if err != nil {
		return nil, ErrNotFound
	}
	// line items
	rows, err := r.pool.Query(ctx, `SELECT id, description, quantity::text, unit_price::text, line_total::text, sort_order FROM invoice_line_items WHERE invoice_id=$1 ORDER BY sort_order`, invoiceID)
	if err != nil {
		return &v, err
	}
	defer rows.Close()
	for rows.Next() {
		var li LineItem
		if err := rows.Scan(&li.ID, &li.Description, &li.Quantity, &li.UnitPrice, &li.LineTotal, &li.SortOrder); err != nil {
			return &v, err
		}
		v.LineItems = append(v.LineItems, li)
	}
	return &v, rows.Err()
}

// Aging returns AR aging buckets for a merchant.
func (r *Repository) Aging(ctx context.Context, merchantID string) ([]AgingBucket, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
		  CASE
		    WHEN due_date >= current_date THEN 'current'
		    WHEN due_date >= current_date - 30 THEN '30'
		    WHEN due_date >= current_date - 60 THEN '60'
		    WHEN due_date >= current_date - 90 THEN '90'
		    ELSE '90plus'
		  END AS bucket,
		  COUNT(*),
		  COALESCE(SUM(total_amount - amount_paid),0)::text
		FROM invoices
		WHERE merchant_id=$1 AND status IN ('sent','partially_paid','overdue')
		GROUP BY bucket ORDER BY bucket`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AgingBucket{}
	for rows.Next() {
		var b AgingBucket
		if err := rows.Scan(&b.Bucket, &b.Count, &b.Amount); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// MarkOverdue flips sent/partially_paid invoices past due to 'overdue'.
func (r *Repository) MarkOverdue(ctx context.Context, merchantID string) (int, error) {
	ct, err := r.pool.Exec(ctx, `
		UPDATE invoices SET status='overdue', updated_at=now()
		WHERE merchant_id=$1 AND status IN ('sent','partially_paid') AND due_date < current_date`, merchantID)
	if err != nil {
		return 0, err
	}
	return int(ct.RowsAffected()), nil
}

// SetHostedToken attaches a hosted payment link token to an invoice and marks it sent.
func (r *Repository) SetHostedToken(ctx context.Context, merchantID, invoiceID, token string) error {
	_, err := r.pool.Exec(ctx, `UPDATE invoices SET hosted_token=$1, status=CASE WHEN status='draft' THEN 'sent' ELSE status END, updated_at=now() WHERE merchant_id=$2 AND id=$3`,
		token, merchantID, invoiceID)
	return err
}

func defaultStr(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

func itoa(n int) string { return strconv.Itoa(n) }

func tzEAT() *time.Location { return time.FixedZone("EAT", 3*3600) }
