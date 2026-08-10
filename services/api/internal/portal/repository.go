package portal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidToken = errors.New("invalid or expired portal token")

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// Create issues a new portal token for a vendor/customer and stores its hash.
func (r *Repository) Create(ctx context.Context, merchantID string, portalType PortalType, entityID, entityName string, ttl time.Duration) (*Access, error) {
	token := id.New("ptok") + id.New("r")
	hash := hashToken(token)
	last4 := tokenLast4(token)
	a := &Access{
		ID: id.New("pac"), MerchantID: merchantID, PortalType: portalType,
		EntityID: entityID, EntityName: entityName, Token: token,
		ExpiresAt: time.Now().Add(ttl).UTC().Format(time.RFC3339),
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO portal_access (id, merchant_id, portal_type, entity_id, entity_name, token_hash, token_last4, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		a.ID, merchantID, string(portalType), entityID, entityName, hash, last4, time.Now().Add(ttl))
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Resolve validates a raw portal token and returns the access + merchant id.
func (r *Repository) Resolve(ctx context.Context, rawToken string) (*Access, error) {
	var a Access
	var merchantID, ptype, entityID, entityName string
	var expires time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT id, merchant_id, portal_type, entity_id, COALESCE(entity_name,''), expires_at
		FROM portal_access WHERE token_hash=$1 AND is_revoked=false`, hashToken(rawToken)).
		Scan(&a.ID, &merchantID, &ptype, &entityID, &entityName, &expires)
	if err == pgx.ErrNoRows {
		return nil, ErrInvalidToken
	}
	if err != nil {
		return nil, err
	}
	if time.Now().After(expires) {
		return nil, ErrInvalidToken
	}
	a.MerchantID = merchantID
	a.PortalType = PortalType(ptype)
	a.EntityID = entityID
	a.EntityName = entityName
	a.ExpiresAt = expires.UTC().Format(time.RFC3339)
	// Bump access counter.
	_, _ = r.pool.Exec(ctx, `UPDATE portal_access SET last_accessed_at=now(), access_count=access_count+1 WHERE id=$1`, a.ID)
	return &a, nil
}

// VendorInvoices returns the AP invoices for a vendor (scoped to merchant + vendor id).
func (r *Repository) VendorInvoices(ctx context.Context, merchantID, vendorID string) ([]VendorInvoice, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT invoice_number, to_char(invoice_date,'YYYY-MM-DD'), to_char(due_date,'YYYY-MM-DD'),
		       subtotal::text, tax_amount::text, total_amount::text, amount_paid::text, status
		FROM ap_invoices WHERE merchant_id=$1 AND vendor_id=$2 ORDER BY due_date`, merchantID, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []VendorInvoice{}
	for rows.Next() {
		var v VendorInvoice
		if err := rows.Scan(&v.InvoiceNumber, &v.InvoiceDate, &v.DueDate, &v.Subtotal, &v.TaxAmount, &v.TotalAmount, &v.AmountPaid, &v.Status); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

// CustomerInvoices returns the invoices for a customer email (scoped to merchant + email).
func (r *Repository) CustomerInvoices(ctx context.Context, merchantID, email string) ([]CustomerInvoice, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT invoice_number, customer_name, to_char(issue_date,'YYYY-MM-DD'), to_char(due_date,'YYYY-MM-DD'),
		       subtotal::text, tax_amount::text, total_amount::text, amount_paid::text, status
		FROM invoices WHERE merchant_id=$1 AND lower(customer_email)=lower($2) ORDER BY due_date`, merchantID, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []CustomerInvoice{}
	for rows.Next() {
		var c CustomerInvoice
		if err := rows.Scan(&c.InvoiceNumber, &c.CustomerName, &c.IssueDate, &c.DueDate, &c.Subtotal, &c.TaxAmount, &c.TotalAmount, &c.AmountPaid, &c.Status); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func hashToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return hex.EncodeToString(sum[:])
}

func tokenLast4(t string) string {
	if len(t) < 4 {
		return t
	}
	return t[len(t)-4:]
}
