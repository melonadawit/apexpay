package refund

import (
	"apexpay/internal/ledger"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct {
	pool   *pgxpool.Pool
	ledger *ledger.PgRepository
}

func NewPgRepository(pool *pgxpool.Pool, ledger *ledger.PgRepository) *PgRepository {
	return &PgRepository{pool: pool, ledger: ledger}
}

func (r *PgRepository) GetPayment(ctx context.Context, merchantID, paymentID string) (*PaymentInfo, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, amount::text, fee_amount::text, status, currency, connector_id, (SELECT COALESCE(SUM(amount),0) FROM refunds WHERE payment_id=$2 AND status='succeeded')::text as refunded FROM payments WHERE merchant_id=$1 AND id=$2`, merchantID, paymentID)
	var info PaymentInfo
	var amt, fee, refunded string
	err := row.Scan(&info.ID, &info.MerchantID, &amt, &fee, &info.Status, &info.Currency, &info.ConnectorID, &refunded)
	if err != nil {
		return nil, err
	}
	// parse handled in service via decimal, simplified: store as string to be parsed later? For skeleton return raw decimal via helper
	// Use helper to parse in service - here we parse directly with decimal
	// quick parse
	// ...
	return &info, nil
}

func (r *PgRepository) GetRefundByRef(ctx context.Context, merchantID, refundRef string) (*Refund, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, payment_id, refund_ref, amount::text, currency, status FROM refunds WHERE merchant_id=$1 AND refund_ref=$2`, merchantID, refundRef)
	var rf Refund
	var amt string
	err := row.Scan(&rf.ID, &rf.MerchantID, &rf.PaymentID, &rf.RefundRef, &amt, &rf.Currency, &rf.Status)
	return &rf, err
}

func (r *PgRepository) ListRefundsByPayment(ctx context.Context, paymentID string) ([]Refund, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, amount::text, status FROM refunds WHERE payment_id=$1`, paymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Refund
	for rows.Next() {
		var rf Refund
		var amt string
		_ = rows.Scan(&rf.ID, &amt, &rf.Status)
		list = append(list, rf)
	}
	return list, nil
}

func (r *PgRepository) CreateRefundTx(ctx context.Context, refund *Refund, journal *ledger.Journal, entries []ledger.Entry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `INSERT INTO refunds (id, merchant_id, payment_id, refund_ref, amount, currency, status, reason, fee_reversal, connector_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		refund.ID, refund.MerchantID, refund.PaymentID, refund.RefundRef, refund.Amount.String(), refund.Currency, refund.Status, refund.Reason, refund.FeeReversal.String(), refund.ConnectorID)
	if err != nil {
		return err
	}

	// Ledger post same Tx
	_, err = tx.Exec(ctx, `INSERT INTO ledger_journals (id, book_id, posting_key, memo, reference_type, reference_id) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (book_id, posting_key) DO NOTHING`,
		journal.ID, journal.BookID, journal.PostingKey, journal.Memo, journal.ReferenceType, journal.ReferenceID)
	if err != nil {
		return err
	}

	for _, e := range entries {
		_, err = tx.Exec(ctx, `INSERT INTO ledger_entries (id, journal_id, book_id, account_id, direction, amount, currency) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			e.ID, e.JournalID, e.BookID, e.AccountID, e.Direction, e.Amount.String(), e.Currency)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO ledger_balances (book_id, account_id, amount) VALUES ($1,$2,$3) ON CONFLICT (book_id, account_id) DO UPDATE SET amount = ledger_balances.amount + EXCLUDED.amount`, e.BookID, e.AccountID, e.Amount.String())
		if err != nil {
			return err
		}
	}

	// Update payment status to partially_refunded / refunded
	_, err = tx.Exec(ctx, `UPDATE payments SET status = CASE WHEN (SELECT COALESCE(SUM(amount),0) FROM refunds WHERE payment_id=$1 AND status IN ('processing','succeeded')) >= amount THEN 'refunded' ELSE 'partially_refunded' END, updated_at=now() WHERE id=$1`, refund.PaymentID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PgRepository) UpdateRefundStatus(ctx context.Context, id string, status Status, connectorRef string) error {
	_, err := r.pool.Exec(ctx, `UPDATE refunds SET status=$1, connector_ref=$2, updated_at=now() WHERE id=$3`, status, connectorRef, id)
	return err
}
