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

var (
	// Sentinel errors for idempotency-state outcomes surfaced to the service layer.
	// A key reused with identical business inputs returns the existing payment; reused
	// with different inputs or while the connector call is still in flight is a conflict.
	ErrIdempotencyConflict   = fmt.Errorf("idempotency key conflict: request differs from original or already resolved")
	ErrIdempotencyInProgress = fmt.Errorf("idempotency request still in progress")
)

func NewPgRepository(pool *pgxpool.Pool, ledgerRepo *ledger.PgRepository) *PgRepository {
	return &PgRepository{pool: pool, ledgerRepo: ledgerRepo}
}

func (r *PgRepository) CreatePaymentTx(ctx context.Context, p *Payment, outboxEventID, idempotencyKey string) error {
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

	if idempotencyKey != "" {
		command, err := tx.Exec(ctx, `UPDATE idempotency_keys SET state='completed', resource_id=$3, response_code=201,
			response_body=jsonb_build_object('id',$3::text) WHERE merchant_id=$1 AND key=$2 AND state='connector_started'`, p.MerchantID, idempotencyKey, p.ID)
		if err != nil {
			return err
		}
		if command.RowsAffected() != 1 {
			return fmt.Errorf("idempotency reservation was not available for completion")
		}
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

// ListByMerchant returns the merchant's most recent payments, newest first.
func (r *PgRepository) ListByMerchant(ctx context.Context, merchantID string, limit int) ([]*Payment, error) {
	if limit <= 0 || limit > 200 {
		limit = 25
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, merchant_id, tx_ref, amount::text, currency, status, COALESCE(method,''), COALESCE(description,''),
		       COALESCE(customer_email,''), connector_id, COALESCE(fee_amount,0)::text, COALESCE(net_amount,0)::text,
		       requires_2fa, two_fa_verified, succeeded_at, created_at
		FROM payments WHERE merchant_id=$1 ORDER BY created_at DESC LIMIT $2`, merchantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Payment
	for rows.Next() {
		var p Payment
		var amtStr, feeStr, netStr string
		var succeededAt *time.Time
		if err := rows.Scan(&p.ID, &p.MerchantID, &p.TxRef, &amtStr, &p.Currency, &p.Status,
			&p.Method, &p.Description, &p.CustomerEmail, &p.ConnectorID, &feeStr, &netStr,
			&p.Requires2FA, &p.TwoFAVerified, &succeededAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.Amount, _ = decimal.NewFromString(amtStr)
		p.FeeAmount, _ = decimal.NewFromString(feeStr)
		p.NetAmount, _ = decimal.NewFromString(netStr)
		p.SucceededAt = succeededAt
		list = append(list, &p)
	}
	return list, rows.Err()
}

// Summary returns dashboard aggregates for a merchant over the trailing window.
type Summary struct {
	TPVToday      decimal.Decimal `json:"tpv_today"`
	TPV7Days      decimal.Decimal `json:"tpv_7_days"`
	SuccessCount  int64           `json:"success_count_7_days"`
	TotalCount    int64           `json:"total_count_7_days"`
	SuccessRate   float64         `json:"success_rate_7_days"`
	ActiveLinks   int64           `json:"active_links"`
	RefundedCount int64           `json:"refunded_count_7_days"`
}

// DashboardSummary computes lightweight aggregates for the merchant dashboard. All money is
// summed as numeric in Postgres and returned as decimal strings (no float money).
func (r *PgRepository) DashboardSummary(ctx context.Context, merchantID string) (*Summary, error) {
	var s Summary
	var tpvToday, tpv7 string
	err := r.pool.QueryRow(ctx, `
		SELECT
		  COALESCE(SUM(amount) FILTER (WHERE status IN ('succeeded') AND created_at >= date_trunc('day', now())), 0)::text,
		  COALESCE(SUM(amount) FILTER (WHERE status IN ('succeeded') AND created_at >= now() - interval '7 days'), 0)::text,
		  COUNT(*) FILTER (WHERE status='succeeded' AND created_at >= now() - interval '7 days'),
		  COUNT(*) FILTER (WHERE created_at >= now() - interval '7 days'),
		  COALESCE(COUNT(*) FILTER (WHERE status='succeeded' AND created_at >= now() - interval '7 days'))::float
		   / NULLIF(COUNT(*) FILTER (WHERE created_at >= now() - interval '7 days'),0),
		  (SELECT COUNT(*) FROM payment_links WHERE merchant_id=$1 AND status IN ('active')),
		  COUNT(*) FILTER (WHERE status IN ('refunded','partially_refunded') AND created_at >= now() - interval '7 days')
		FROM payments WHERE merchant_id=$1`, merchantID).
		Scan(&tpvToday, &tpv7, &s.SuccessCount, &s.TotalCount, &s.SuccessRate, &s.ActiveLinks, &s.RefundedCount)
	if err != nil {
		return nil, err
	}
	s.TPVToday, _ = decimal.NewFromString(tpvToday)
	s.TPV7Days, _ = decimal.NewFromString(tpv7)
	return &s, nil
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
	if err != nil || clearingID == "" || payableID == "" || feeID == "" {
		return fmt.Errorf("payment ledger accounts unavailable: %w", err)
	}
	journal.BookID = bookID
	for i := range entries {
		entries[i].BookID = bookID
		switch entries[i].AccountID {
		case "liability:merchant_payable":
			entries[i].AccountID = payableID
		case "liability:platform_fee_due":
			entries[i].AccountID = feeID
		default:
			entries[i].AccountID = clearingID
		}
	}
	_, err = tx.Exec(ctx, `UPDATE payments SET status=$1, succeeded_at=$2, updated_at=now() WHERE id=$3 AND status IN ('created','pending','processing')`, status, succeededAt, paymentID)
	if err != nil {
		return err
	}

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

func (r *PgRepository) ReserveIdempotency(ctx context.Context, merchantID, key, requestHash string) (*Payment, error) {
	command, err := r.pool.Exec(ctx, `INSERT INTO idempotency_keys (merchant_id, key, request_hash, response_code, response_body, resource_type, state)
		VALUES ($1,$2,$3,0,'{}'::jsonb,'payment','in_progress') ON CONFLICT (merchant_id, key) DO NOTHING`, merchantID, key, requestHash)
	if err != nil {
		return nil, err
	}
	if command.RowsAffected() == 1 {
		return nil, nil
	}

	var storedHash, state, paymentID string
	err = r.pool.QueryRow(ctx, `SELECT request_hash, state, COALESCE(resource_id,'') FROM idempotency_keys WHERE merchant_id=$1 AND key=$2 AND resource_type='payment'`, merchantID, key).Scan(&storedHash, &state, &paymentID)
	if err != nil {
		return nil, err
	}
	if storedHash != requestHash {
		return nil, ErrIdempotencyConflict
	}
	if state == "completed" && paymentID != "" {
		return r.getByID(ctx, merchantID, paymentID)
	}
	if state == "retry_authorized" {
		command, err := r.pool.Exec(ctx, `UPDATE idempotency_keys SET state='in_progress',response_code=0,response_body='{}'::jsonb WHERE merchant_id=$1 AND key=$2 AND state='retry_authorized'`, merchantID, key)
		if err != nil {
			return nil, err
		}
		if command.RowsAffected() == 1 {
			return nil, nil
		}
	}
	return nil, ErrIdempotencyInProgress
}

func (r *PgRepository) MarkConnectorStarted(ctx context.Context, merchantID, key, txRef string) error {
	command, err := r.pool.Exec(ctx, `UPDATE idempotency_keys
		SET state='connector_started', response_body=jsonb_build_object('tx_ref',$3::text)
		WHERE merchant_id=$1 AND key=$2 AND state='in_progress'`, merchantID, key, txRef)
	if err != nil {
		return err
	}
	if command.RowsAffected() != 1 {
		return fmt.Errorf("idempotency reservation was not available before connector call")
	}
	return nil
}

func (r *PgRepository) FailIdempotency(ctx context.Context, merchantID, key string) error {
	_, err := r.pool.Exec(ctx, `UPDATE idempotency_keys SET state='failed', response_code=502,
		response_body='{"error":"initialization_failed"}'::jsonb WHERE merchant_id=$1 AND key=$2 AND state IN ('in_progress','connector_started')`, merchantID, key)
	return err
}

// GetPaymentDetail returns a payment plus its lifecycle ledger journals and
// entries for the NBE exam console. Real DB-backed, not mocked.
func (r *PgRepository) GetPaymentDetail(ctx context.Context, merchantID, paymentID string) (*PaymentDetail, error) {
	p, err := r.getByID(ctx, merchantID, paymentID)
	if err != nil {
		return nil, err
	}
	journals, err := r.ledgerRepo.ListJournalsByRef(ctx, "payment", paymentID)
	if err != nil {
		return nil, err
	}
	detail := &PaymentDetail{Payment: p, Journals: []PaymentJournal{}}
	for _, j := range journals {
		pj := PaymentJournal{
			ID:         j.ID,
			BookID:     j.BookID,
			PostingKey: j.PostingKey,
			Memo:       j.Memo,
			CreatedAt:  j.CreatedAt,
		}
		rows, err := r.pool.Query(ctx, `SELECT e.direction, e.amount::text, e.currency, COALESCE(a.code,''), COALESCE(a.name,'')
			FROM ledger_entries e LEFT JOIN ledger_accounts a ON a.id=e.account_id
			WHERE e.journal_id=$1 ORDER BY e.created_at, e.id`, j.ID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var pe PaymentEntry
			if err := rows.Scan(&pe.Direction, &pe.Amount, &pe.Currency, &pe.AccountCode, &pe.AccountName); err != nil {
				rows.Close()
				return nil, err
			}
			pj.Entries = append(pj.Entries, pe)
		}
		rows.Close()
		detail.Journals = append(detail.Journals, pj)
	}
	return detail, nil
}

func (r *PgRepository) getByID(ctx context.Context, merchantID, paymentID string) (*Payment, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, tx_ref, amount::text, currency, status, connector_id, connector_ref, fee_amount::text, net_amount::text, requires_2fa, two_fa_verified, checkout_url FROM payments WHERE merchant_id=$1 AND id=$2`, merchantID, paymentID)
	var p Payment
	var amount, fee, net string
	if err := row.Scan(&p.ID, &p.MerchantID, &p.TxRef, &amount, &p.Currency, &p.Status, &p.ConnectorID, &p.ConnectorRef, &fee, &net, &p.Requires2FA, &p.TwoFAVerified, &p.CheckoutURL); err != nil {
		return nil, err
	}
	p.Amount, _ = decimal.NewFromString(amount)
	p.FeeAmount, _ = decimal.NewFromString(fee)
	p.NetAmount, _ = decimal.NewFromString(net)
	return &p, nil
}

func (r *PgRepository) Mark2FAVerified(ctx context.Context, merchantID, paymentID string) error {
	command, err := r.pool.Exec(ctx, `UPDATE payments
		SET two_fa_verified=true, updated_at=now()
		WHERE id=$1 AND merchant_id=$2 AND requires_2fa=true AND status IN ('created','pending','processing')`, paymentID, merchantID)
	if err != nil {
		return err
	}
	if command.RowsAffected() != 1 {
		return fmt.Errorf("payment not found or not awaiting 2FA")
	}
	return nil
}
