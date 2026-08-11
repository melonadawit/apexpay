package payout

import (
	"apexpay/internal/ledger"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type PgRepository struct {
	pool   *pgxpool.Pool
	ledger *ledger.PgRepository
}

func NewPgRepository(pool *pgxpool.Pool, ledger *ledger.PgRepository) *PgRepository {
	return &PgRepository{pool: pool, ledger: ledger}
}

func (r *PgRepository) CreateBeneficiary(ctx context.Context, b *Beneficiary) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO beneficiaries (id, merchant_id, name, account_no_masked, account_no_hash, bank_code, bank_name, type, verification_status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		b.ID, b.MerchantID, b.Name, b.AccountNoMasked, b.AccountNoHash, b.BankCode, b.BankName, b.Type, b.VerificationStatus)
	return err
}

func (r *PgRepository) GetBeneficiary(ctx context.Context, merchantID, id string) (*Beneficiary, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, name, account_no_masked, bank_code, bank_name FROM beneficiaries WHERE merchant_id=$1 AND id=$2`, merchantID, id)
	var b Beneficiary
	err := row.Scan(&b.ID, &b.MerchantID, &b.Name, &b.AccountNoMasked, &b.BankCode, &b.BankName)
	return &b, err
}

func (r *PgRepository) CreateBatchTx(ctx context.Context, batch *PayoutBatch, journal *ledger.Journal, entries []ledger.Entry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, `INSERT INTO payout_batches (id, merchant_id, book_id, batch_ref, amount, currency, status, total_count) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		batch.ID, batch.MerchantID, batch.BookID, batch.BatchRef, batch.Amount.String(), batch.Currency, batch.Status, len(batch.Payouts))
	if err != nil {
		return err
	}
	// ledger
	_, err = tx.Exec(ctx, `INSERT INTO ledger_journals (id, book_id, posting_key, memo, reference_type, reference_id) VALUES ($1,$2,$3,$4,$5,$6)`, journal.ID, journal.BookID, journal.PostingKey, journal.Memo, journal.ReferenceType, journal.ReferenceID)
	if err != nil {
		return err
	}
	for _, e := range entries {
		_, err = tx.Exec(ctx, `INSERT INTO ledger_entries (id, journal_id, book_id, account_id, direction, amount, currency) VALUES ($1,$2,$3,$4,$5,$6,$7)`, e.ID, e.JournalID, e.BookID, e.AccountID, e.Direction, e.Amount.String(), e.Currency)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PgRepository) CreatePayout(ctx context.Context, p *Payout) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO payouts (id, merchant_id, batch_id, beneficiary_id, payout_ref, amount, currency, status, method) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		p.ID, p.MerchantID, p.BatchID, p.BeneficiaryID, p.PayoutRef, p.Amount.String(), p.Currency, p.Status, p.Method)
	return err
}

func (r *PgRepository) CreateBulkTx(ctx context.Context, batch *PayoutBatch, payouts []Payout, journal *ledger.Journal, entries []ledger.Entry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, `INSERT INTO payout_batches (id, merchant_id, batch_ref, amount, currency, status, total_count) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		batch.ID, batch.MerchantID, batch.BatchRef, batch.Amount.String(), batch.Currency, batch.Status, len(payouts))
	if err != nil {
		return err
	}
	for _, p := range payouts {
		_, err = tx.Exec(ctx, `INSERT INTO payouts (id, merchant_id, batch_id, beneficiary_id, payout_ref, amount, currency, status, method) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			p.ID, p.MerchantID, p.BatchID, p.BeneficiaryID, p.PayoutRef, p.Amount.String(), p.Currency, p.Status, p.Method)
		if err != nil {
			return err
		}
	}
	// ledger M3
	_, err = tx.Exec(ctx, `INSERT INTO ledger_journals (id, book_id, posting_key, memo, reference_type, reference_id) VALUES ($1,$2,$3,$4,$5,$6)`,
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
	}
	return tx.Commit(ctx)
}

func (r *PgRepository) GetBatch(ctx context.Context, merchantID, batchID string) (*PayoutBatch, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, batch_ref, amount::text, currency, status FROM payout_batches WHERE merchant_id=$1 AND id=$2`, merchantID, batchID)
	var b PayoutBatch
	var amt string
	err := row.Scan(&b.ID, &b.MerchantID, &b.BatchRef, &amt, &b.Currency, &b.Status)
	if err != nil {
		return nil, err
	}
	b.Amount, _ = decimal.NewFromString(amt)
	return &b, nil
}

// ListBatches returns the merchant's payout batches, newest first, with a
// payout count per batch.
func (r *PgRepository) ListBatches(ctx context.Context, merchantID string) ([]PayoutBatch, error) {
	rows, err := r.pool.Query(ctx, `SELECT pb.id, pb.merchant_id, COALESCE(pb.book_id,''), pb.batch_ref,
		pb.amount::text, pb.currency, pb.status, COALESCE(pb.approved_by,''), COALESCE(pb.created_at, now()),
		(SELECT COUNT(*) FROM payouts p WHERE p.batch_id=pb.id)
		FROM payout_batches pb WHERE pb.merchant_id=$1 ORDER BY pb.created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []PayoutBatch
	for rows.Next() {
		var b PayoutBatch
		var amt, bookID, approvedBy string
		var count int
		if err := rows.Scan(&b.ID, &b.MerchantID, &bookID, &b.BatchRef, &amt, &b.Currency, &b.Status, &approvedBy, &b.CreatedAt, &count); err != nil {
			return nil, err
		}
		if bookID != "" {
			b.BookID = &bookID
		}
		if approvedBy != "" {
			b.ApprovedBy = &approvedBy
		}
		b.Amount, _ = decimal.NewFromString(amt)
		b.TotalCount = count
		list = append(list, b)
	}
	return list, rows.Err()
}

// ListBeneficiaries returns the merchant's saved beneficiaries.
func (r *PgRepository) ListBeneficiaries(ctx context.Context, merchantID string) ([]Beneficiary, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, name, account_no_masked, bank_code,
		COALESCE(bank_name,''), COALESCE(type,''), verification_status, created_at
		FROM beneficiaries WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Beneficiary
	for rows.Next() {
		var b Beneficiary
		if err := rows.Scan(&b.ID, &b.MerchantID, &b.Name, &b.AccountNoMasked, &b.BankCode, &b.BankName, &b.Type, &b.VerificationStatus, &b.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, rows.Err()
}

func (r *PgRepository) UpdateBatchStatus(ctx context.Context, batchID, status, approvedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payout_batches SET status=$1, approved_by=$2 WHERE id=$3`, status, approvedBy, batchID)
	return err
}

func (r *PgRepository) UpdatePayoutStatus(ctx context.Context, payoutID string, status PayoutStatus, connectorRef string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payouts SET status=$1, connector_ref=$2, updated_at=now() WHERE id=$3`, status, connectorRef, payoutID)
	return err
}

func (r *PgRepository) GetMerchantBalance(ctx context.Context, merchantID string) (decimal.Decimal, error) {
	// Simplified: sum of succeeded payments net - payouts succeeded
	var balStr string
	err := r.pool.QueryRow(ctx, `SELECT COALESCE((SELECT SUM(net_amount) FROM payments WHERE merchant_id=$1 AND status='succeeded'),0) - COALESCE((SELECT SUM(amount) FROM payouts WHERE merchant_id=$1 AND status IN ('queued','processing','succeeded')),0)`, merchantID).Scan(&balStr)
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromString(balStr)
}
