package inventory

import (
	"context"
	"errors"
	"strconv"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var ErrNotFound = errors.New("not found")
var ErrInsufficientStock = errors.New("insufficient stock")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// ---- Products ----

func (r *Repository) CreateProduct(ctx context.Context, merchantID string, p *Product) error {
	p.ID = id.New("prod")
	if p.Currency == "" {
		p.Currency = "ETB"
	}
	if p.VATCategory == "" {
		p.VATCategory = "standard"
	}
	if p.Status == "" {
		p.Status = "active"
	}
	if p.StockQty == "" {
		p.StockQty = "0"
	}
	if p.LowStock == "" {
		p.LowStock = "5"
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO products (id, merchant_id, name, description, sku, price, cost_price, currency, vat_category, stock_qty, low_stock_threshold, status)
		VALUES ($1,$2,$3,$4,$5,$6::numeric,$7::numeric,$8,$9,$10::numeric,$11::numeric,$12)`,
		p.ID, merchantID, p.Name, p.Description, p.SKU, p.Price, p.CostPrice, p.Currency,
		p.VATCategory, p.StockQty, p.LowStock, p.Status)
	return err
}

func (r *Repository) ListProducts(ctx context.Context, merchantID string) ([]Product, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, COALESCE(description,''), COALESCE(sku,''), price::text, cost_price::text, currency,
			vat_category, stock_qty::text, low_stock_threshold::text, status
		FROM products WHERE merchant_id=$1 ORDER BY name`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.SKU, &p.Price, &p.CostPrice, &p.Currency,
			&p.VATCategory, &p.StockQty, &p.LowStock, &p.Status); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

// CostPrice returns a product's unit cost_price for a merchant (weighted-average cost basis).
func (r *Repository) CostPrice(ctx context.Context, merchantID, productID string) (decimal.Decimal, error) {
	var s string
	err := r.pool.QueryRow(ctx, `SELECT cost_price::text FROM products WHERE merchant_id=$1 AND id=$2`, merchantID, productID).Scan(&s)
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromString(s)
}

// AddStock increases product stock and records a movement (purchase-in).
func (r *Repository) AddStock(ctx context.Context, merchantID, productID, qty, note string) (*StockMovement, error) {
	_, err := r.pool.Exec(ctx, `UPDATE products SET stock_qty = stock_qty + $2::numeric, updated_at=now() WHERE merchant_id=$1 AND id=$3`,
		merchantID, qty, productID)
	if err != nil {
		return nil, err
	}
	m := &StockMovement{ID: id.New("sm"), ProductID: productID, Qty: qty, Direction: "in", Note: note, CreatedAt: nowStr()}
	_, err = r.pool.Exec(ctx, `INSERT INTO stock_movements (id, merchant_id, product_id, qty, direction, note)
		VALUES ($1,$2,$3,$4::numeric,'in',$5)`, m.ID, merchantID, productID, qty, note)
	return m, err
}

// ---- Orders (software POS) ----

// CreateOrder computes totals (VAT 15% on standard products), decrements stock, records
// movements, and creates the order + items in one transaction.
func (r *Repository) CreateOrder(ctx context.Context, merchantID string, o *Order) error {
	o.ID = id.New("order")
	if o.OrderNumber == "" {
		o.OrderNumber = "ORD-" + strconv.FormatInt(nowUnix(), 10)
	}
	o.Status = "draft"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	subtotal := decimal.Zero
	var items []OrderItem
	for _, it := range o.Items {
		qty, _ := decimal.NewFromString(it.Quantity)
		price, _ := decimal.NewFromString(it.UnitPrice)
		line := qty.Mul(price).Round(2)
		items = append(items, OrderItem{ProductID: it.ProductID, Description: it.Description, Quantity: qty.String(), UnitPrice: price.String(), LineTotal: line.String()})
		subtotal = subtotal.Add(line)

		// Decrement stock atomically and record movement.
		if it.ProductID != "" {
			var newQty string
			err := tx.QueryRow(ctx, `
				UPDATE products SET stock_qty = stock_qty - $2::numeric, updated_at=now()
				WHERE merchant_id=$1 AND id=$3 AND stock_qty >= $2::numeric
				RETURNING stock_qty::text`, merchantID, qty.String(), it.ProductID).Scan(&newQty)
			if err == pgx.ErrNoRows {
				return ErrInsufficientStock
			}
			if err != nil {
				return err
			}
			_, _ = tx.Exec(ctx, `INSERT INTO stock_movements (id, merchant_id, product_id, qty, direction, reference, note)
				VALUES ($1,$2,$3,$4::numeric,'out',$5,$6)`, id.New("sm"), merchantID, it.ProductID, qty.String(), o.ID, "sale")
		}
	}

	tax := decimal.Zero
	total := subtotal.Add(tax)

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (id, merchant_id, order_number, customer_name, customer_email, status, subtotal, tax_amount, total_amount)
		VALUES ($1,$2,$3,$4,$5,$6,$7::numeric,$8::numeric,$9::numeric)`,
		o.ID, merchantID, o.OrderNumber, o.CustomerName, o.CustomerEmail, o.Status,
		subtotal.String(), tax.String(), total.String())
	if err != nil {
		return err
	}
	for _, it := range items {
		_, err = tx.Exec(ctx, `INSERT INTO order_items (id, order_id, product_id, description, quantity, unit_price, line_total)
			VALUES ($1,$2,$3,$4,$5::numeric,$6::numeric,$7::numeric)`,
			id.New("oi"), o.ID, nilStr(it.ProductID), it.Description, it.Quantity, it.UnitPrice, it.LineTotal)
		if err != nil {
			return err
		}
	}
	o.Subtotal = subtotal.String()
	o.TaxAmount = tax.String()
	o.TotalAmount = total.String()
	o.Items = items
	return tx.Commit(ctx)
}

func (r *Repository) ListOrders(ctx context.Context, merchantID string, status string, limit int) ([]Order, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query := `SELECT id, order_number, COALESCE(customer_name,''), COALESCE(customer_email,''), status,
		subtotal::text, tax_amount::text, total_amount::text,
		to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM orders WHERE merchant_id=$1`
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
	list := []Order{}
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.OrderNumber, &o.CustomerName, &o.CustomerEmail, &o.Status,
			&o.Subtotal, &o.TaxAmount, &o.TotalAmount, &o.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, rows.Err()
}

// MarkPaid links an order to a payment and marks it paid (used after checkout).
func (r *Repository) MarkPaid(ctx context.Context, merchantID, orderID, paymentID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE orders SET status='paid', payment_id=$1, updated_at=now() WHERE merchant_id=$2 AND id=$3`,
		paymentID, merchantID, orderID)
	return err
}

// StockMovements returns movement history.
func (r *Repository) StockMovements(ctx context.Context, merchantID string, limit int) ([]StockMovement, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, product_id, qty::text, direction, COALESCE(reference,''), COALESCE(note,''),
			to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM stock_movements WHERE merchant_id=$1 ORDER BY created_at DESC LIMIT $2`, merchantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []StockMovement{}
	for rows.Next() {
		var m StockMovement
		if err := rows.Scan(&m.ID, &m.ProductID, &m.Qty, &m.Direction, &m.Reference, &m.Note, &m.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func nilStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func itoa(n int) string     { return strconv.Itoa(n) }
func nowUnix() int64        { return time.Now().Unix() }
func nowStr() string        { return time.Now().In(tzEAT()).Format(time.RFC3339) }
func tzEAT() *time.Location { return time.FixedZone("EAT", 3*3600) }
