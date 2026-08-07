package payment

import (
	"context"
	"fmt"
	"time"

	"apexpay/internal/ledger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type PgRepository struct {
	pool       *pgxpool.Pool
	ledgerRepo *ledger.PgRepository
}

func NewPgRepository(pool *pgxpool.Pool, ledgerRepo *ledger.PgRepository) *PgRepository {
	return &PgRepository{pool: pool, ledgerRepo: ledgerRepo}
}

func (r *PgRepository) CreatePaymentTx(ctx context.Context, p *Payment, outboxEventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `INSERT INTO payments (id, merchant_id, tx_ref, amount, currency, status, method, description, customer_email, connector_id, connector_ref, routing_rule_id, checkout_url, return_url, callback_url, fee_amount, net_amount, requires_2fa, two_fa_verified, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19, now(), now())`,
		p.ID, p.MerchantID, p.TxRef, p.Amount.String(), p.Currency, p.Status, p.Method, p.Description, p.CustomerEmail, p.ConnectorID, p.ConnectorRef, p.RoutingRuleID, p.CheckoutURL, p.ReturnURL, p.CallbackURL, p.FeeAmount.String(), p.NetAmount.String(), p.Requires2FA, p.TwoFAVerified)
	if err != nil {
		return err
	}

	// Outbox payment.created for webhook worker
	_, err = tx.Exec(ctx, `INSERT INTO outbox_events (id, merchant_id, aggregate_type, aggregate_id, event_type, payload) VALUES ($1,$2,'payment',$3,'payment.created',$4)`,
		outboxEventID, p.MerchantID, p.ID, fmt.Sprintf(`{"payment_id":"%s","tx_ref":"%s","amount":"%s"}`, p.ID, p.TxRef, p.Amount.String()))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PgRepository) GetByTxRef(ctx context.Context, merchantID, txRef string) (*Payment, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, tx_ref, amount::text, currency, status, connector_id, connector_ref, fee_amount::text, net_amount::text, requires_2fa, two_fa_verified, checkout_url FROM payments WHERE merchant_id=$1 AND tx_ref=$2`, merchantID, txRef)
	var p Payment
	var amtStr, feeStr, netStr string
	err := row.Scan(&p.ID, &p.MerchantID, &p.TxRef, &amtStr, &p.Currency, &p.Status, &p.ConnectorID, &p.ConnectorRef, &feeStr, &netStr, &p.Requires2FA, &p.TwoFAVerified, &p.CheckoutURL)
	if err != nil {
		return nil, err
	}
	p.Amount, _ = decimal.NewFromString(amtStr)
	p.FeeAmount, _ = decimal.NewFromString(feeStr)
	p.NetAmount, _ = decimal.NewFromString(netStr)
	return &p, nil
}

func (r *PgRepository) UpdateStatusTx(ctx context.Context, paymentID string, status Status, journal *ledger.Journal, entries []ledger.Entry, succeededAt *time.Time) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Resolve the active merchant operating book and account *IDs* in this same
	// transaction. Services pass account codes; database foreign keys require IDs.
	var bookID, clearingID, payableID, feeID string
	err = tx.QueryRow(ctx, `SELECT lb.id,
		MAX(la.id) FILTER (WHERE la.code=$2),
		MAX(la.id) FILTER (WHERE la.code='liability:merchant_payable'),
		MAX(la.id) FILTER (WHERE la.code='liability:platform_fee_due')
		FROM ledger_books lb JOIN ledger_accounts la ON la.book_id=lb.id
		JOIN payments p ON p.merchant_id=lb.merchant_id
		WHERE p.id=$1 AND lb.book_type='merchant_operating' AND lb.status='open'
		GROUP BY lb.id ORDER BY lb.id LIMIT 1`, paymentID, entries[0].AccountID).Scan(&bookID, &clearingID, &payableID, &feeID)
	if err != nil || clearingID == "" || payableID == "" || feeID == "" { return fmt.Errorf("payment ledger accounts unavailable: %w", err) }
	journal.BookID = bookID
	for i := range entries {
		entries[i].BookID = bookID
		switch entries[i].AccountID {
		case "liability:merchant_payable": entries[i].AccountID = payableID
		case "liability:platform_fee_due": entries[i].AccountID = feeID
		default: entries[i].AccountID = clearingID
		}
	}
	_, err = tx.Exec(ctx, `UPDATE payments SET status=$1, succeeded_at=$2, updated_at=now() WHERE id=$3 AND status IN ('created','pending','processing')`, status, succeededAt, paymentID)
	if err != nil { return err }

	// Ledger post in same Tx per DATABASE transaction boundary
	_, err = tx.Exec(ctx, `INSERT INTO ledger_journals (id, book_id, posting_key, memo, reference_type, reference_id, transfer_group) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (book_id, posting_key) DO NOTHING`,
		journal.ID, journal.BookID, journal.PostingKey, journal.Memo, journal.ReferenceType, journal.ReferenceID, journal.TransferGroup)
	if err != nil {
		return err
	}

	for _, e := range entries {
		_, err = tx.Exec(ctx, `INSERT INTO ledger_entries (id, journal_id, book_id, account_id, direction, amount, currency) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (id) DO NOTHING`,
			e.ID, e.JournalID, e.BookID, e.AccountID, e.Direction, e.Amount.String(), e.Currency)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO ledger_balances (book_id, account_id, amount, updated_at) VALUES ($1,$2,$3, now()) ON CONFLICT (book_id, account_id) DO UPDATE SET amount = ledger_balances.amount + EXCLUDED.amount, updated_at=now()`,
			e.BookID, e.AccountID, e.Amount.String())
		if err != nil {
			return err
		}
	}

	// Outbox payment.succeeded
	outboxID := fmt.Sprintf("outbox_%s_%d", paymentID, time.Now().UnixNano())
	_, err = tx.Exec(ctx, `INSERT INTO outbox_events (id, merchant_id, aggregate_type, aggregate_id, event_type, payload) SELECT $1, merchant_id, 'payment', $2, 'payment.succeeded', $3 FROM payments WHERE id=$2`, outboxID, paymentID, fmt.Sprintf(`{"payment_id":"%s","status":"%s"}`, paymentID, status))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PgRepository) GetIdempotency(ctx context.Context, merchantID, key string) (*Payment, string, error) {
	var paymentID, requestHash string
	err := r.pool.QueryRow(ctx, `SELECT resource_id, request_hash FROM idempotency_keys WHERE merchant_id=$1 AND key=$2 AND resource_type='payment'`, merchantID, key).Scan(&paymentID, &requestHash)
	if err != nil { return nil, "", err }
	p, err := r.getByID(ctx, merchantID, paymentID)
	if err != nil { return nil, "", err }
	return p, requestHash, nil
}

func (r *PgRepository) getByID(ctx context.Context, merchantID, paymentID string) (*Payment, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, tx_ref, amount::text, currency, status, connector_id, connector_ref, fee_amount::text, net_amount::text, requires_2fa, two_fa_verified, checkout_url FROM payments WHERE merchant_id=$1 AND id=$2`, merchantID, paymentID)
	var p Payment
	var amount, fee, net string
	if err := row.Scan(&p.ID, &p.MerchantID, &p.TxRef, &amount, &p.Currency, &p.Status, &p.ConnectorID, &p.ConnectorRef, &fee, &net, &p.Requires2FA, &p.TwoFAVerified, &p.CheckoutURL); err != nil { return nil, err }
	p.Amount, _ = decimal.NewFromString(amount); p.FeeAmount, _ = decimal.NewFromString(fee); p.NetAmount, _ = decimal.NewFromString(net)
	return &p, nil
}

func (r *PgRepository) SaveIdempotency(ctx context.Context, merchantID, key, requestHash string, payment *Payment) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO idempotency_keys (merchant_id, key, request_hash, response_code, response_body, resource_type, resource_id) VALUES ($1,$2,$3,200,$4,'payment',$5) ON CONFLICT (merchant_id, key) DO NOTHING`,
		merchantID, key, requestHash, fmt.Sprintf(`{"id":"%s","checkout_url":"%s"}`, payment.ID, payment.CheckoutURL), payment.ID)
	return err
}

func (r *PgRepository) Mark2FAVerified(ctx context.Context, merchantID, paymentID string) error {
	command, err := r.pool.Exec(ctx, `UPDATE payments
		SET two_fa_verified=true, updated_at=now()
		WHERE id=$1 AND merchant_id=$2 AND requires_2fa=true AND status IN ('created','pending','processing')`, paymentID, merchantID)
	if err != nil { return err }
	if command.RowsAffected() != 1 { return fmt.Errorf("payment not found or not awaiting 2FA") }
	return nil
}
