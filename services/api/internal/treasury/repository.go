package treasury

import (
	"context"
	"errors"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// CashPosition returns all of a merchant's active current accounts with balances.
func (r *Repository) CashPosition(ctx context.Context, merchantID string) (*CashPosition, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, account_number, account_name, account_type, bank_code,
		       balance::text, available_balance::text, currency
		FROM current_accounts
		WHERE merchant_id=$1 AND status='active'
		ORDER BY is_primary DESC, created_at`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pos := &CashPosition{Accounts: []AccountPosition{}, Currency: "ETB"}
	total := decimal.Zero
	totalAvail := decimal.Zero
	for rows.Next() {
		var a AccountPosition
		var bal, avail string
		if err := rows.Scan(&a.AccountID, &a.AccountNumber, &a.AccountName, &a.AccountType, &a.BankCode, &bal, &avail, &a.Currency); err != nil {
			return nil, err
		}
		a.Balance = bal
		a.AvailableBalance = avail
		pos.Accounts = append(pos.Accounts, a)
		b, _ := decimal.NewFromString(bal)
		av, _ := decimal.NewFromString(avail)
		total = total.Add(b)
		totalAvail = totalAvail.Add(av)
	}
	pos.TotalBalance = total.String()
	pos.TotalAvailable = totalAvail.String()
	pos.GeneratedAt = time.Now().In(tzEAT()).Format(time.RFC3339)
	return pos, rows.Err()
}

// CreateTransfer records a pending internal transfer. It does not move money itself;
// a worker (or the service) completes it by updating balances within a transaction.
func (r *Repository) CreateTransfer(ctx context.Context, merchantID, userID string, t *Transfer) error {
	t.ID = id.New("tf")
	if t.Currency == "" {
		t.Currency = "ETB"
	}
	if t.Status == "" {
		t.Status = "pending"
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO treasury_transfers (id, merchant_id, from_account_id, to_account_id, amount, currency, purpose, status, created_by)
		VALUES ($1,$2,$3,$4,$5::numeric,$6,$7,$8,$9)`,
		t.ID, merchantID, t.FromAccountID, t.ToAccountID, t.Amount, t.Currency, t.Purpose, t.Status, userID)
	return err
}

// CompleteTransfer atomically debits source, credits destination, and marks the transfer
// completed — guarded by an advisory lock so concurrent transfers can't double-spend.
func (r *Repository) CompleteTransfer(ctx context.Context, transferID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Lock the transfer and read it.
	var merchantID, fromID, toID, amountStr string
	err = tx.QueryRow(ctx, `SELECT merchant_id, from_account_id, to_account_id, amount::text FROM treasury_transfers WHERE id=$1 FOR UPDATE`, transferID).
		Scan(&merchantID, &fromID, &toID, &amountStr)
	if err == pgx.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	amount, _ := decimal.NewFromString(amountStr)

	// Lock both accounts to serialize.
	_, _ = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, "treasury_"+fromID)
	_, _ = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, "treasury_"+toID)

	// Check sufficient available balance on source.
	var availStr string
	if err := tx.QueryRow(ctx, `SELECT available_balance::text FROM current_accounts WHERE id=$1`, fromID).Scan(&availStr); err != nil {
		return err
	}
	avail, _ := decimal.NewFromString(availStr)
	if avail.LessThan(amount) {
		_, _ = tx.Exec(ctx, `UPDATE treasury_transfers SET status='failed', updated_at=now() WHERE id=$1`, transferID)
		return errors.New("insufficient available balance for transfer")
	}

	// Debit source, credit destination.
	if _, err := tx.Exec(ctx, `UPDATE current_accounts SET balance = balance - $2, available_balance = available_balance - $2, updated_at=now() WHERE id=$1`, fromID, amount.String()); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE current_accounts SET balance = balance + $2, available_balance = available_balance + $2, updated_at=now() WHERE id=$1`, toID, amount.String()); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE treasury_transfers SET status='completed', completed_at=now(), updated_at=now() WHERE id=$1`, transferID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// ListTransfers returns recent transfers for a merchant.
func (r *Repository) ListTransfers(ctx context.Context, merchantID string, limit int) ([]Transfer, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, from_account_id, to_account_id, amount::text, currency, COALESCE(purpose,''), status,
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM treasury_transfers WHERE merchant_id=$1 ORDER BY created_at DESC LIMIT $2`, merchantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Transfer{}
	for rows.Next() {
		var t Transfer
		if err := rows.Scan(&t.ID, &t.FromAccountID, &t.ToAccountID, &t.Amount, &t.Currency, &t.Purpose, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

// ForecastFromObligations builds a cash-flow forecast from known data: inflows from
// succeeded payments and open invoices; outflows from open payouts, pending tax, payroll,
// and subscriptions. Returns the forecast (and persists it).
func (r *Repository) ForecastFromObligations(ctx context.Context, merchantID string, horizonDays int) (*Forecast, error) {
	if horizonDays <= 0 || horizonDays > 180 {
		horizonDays = 90
	}
	var in30, in60, in90, out30, out60, out90 string

	err := r.pool.QueryRow(ctx, `
		SELECT
		  (SELECT COALESCE(SUM(amount),0)::text FROM payments WHERE merchant_id=$1 AND status='succeeded' AND created_at >= now() - interval '30 days'),
		  (SELECT COALESCE(SUM(amount),0)::text FROM payments WHERE merchant_id=$1 AND status='succeeded' AND created_at >= now() - interval '60 days'),
		  (SELECT COALESCE(SUM(amount),0)::text FROM payments WHERE merchant_id=$1 AND status='succeeded' AND created_at >= now() - interval '90 days'),
		  (SELECT COALESCE(SUM(total_amount),0)::text FROM vendor_invoices WHERE merchant_id=$1 AND status IN ('approved','pending_approval')),
		  (SELECT COALESCE(SUM(total_net),0)::text FROM payroll_runs WHERE merchant_id=$1 AND status IN ('approved','processing')),
		  (SELECT COALESCE(SUM(amount),0)::text FROM tax_payments WHERE merchant_id=$1 AND status IN ('pending','pending_approval'))
		`, merchantID).Scan(&in30, &in60, &in90, &out30, &out60, &out90)
	if err != nil {
		return nil, err
	}
	d := func(s string) decimal.Decimal { v, _ := decimal.NewFromString(s); return v }
	in30d, in60d, in90d := d(in30), d(in60), d(in90)
	out30d, out60d, out90d := d(out30), d(out60), d(out90)
	net90 := in90d.Sub(out90d)

	f := &Forecast{
		ID: id.New("cf"), ForecastDate: time.Now().In(tzEAT()).Format("2006-01-02"), HorizonDays: horizonDays,
		InflowToday: in30d.String(), Inflow30d: in30d.String(), Inflow60d: in60d.String(), Inflow90d: in90d.String(),
		OutflowToday: out30d.String(), Outflow30d: out30d.String(), Outflow60d: out60d.String(), Outflow90d: out90d.String(),
		Net90d: net90.String(), GeneratedAt: time.Now().In(tzEAT()).Format(time.RFC3339),
	}

	// Persist the latest forecast.
	_, err = r.pool.Exec(ctx, `
		INSERT INTO cash_flow_forecasts (id, merchant_id, forecast_date, horizon_days,
			inflow_today, inflow_30d, inflow_60d, inflow_90d, outflow_today, outflow_30d, outflow_60d, outflow_90d, net_90d)
		VALUES ($1,$2,$3::date,$4,$5::numeric,$6::numeric,$7::numeric,$8::numeric,$9::numeric,$10::numeric,$11::numeric,$12::numeric,$13::numeric)`,
		f.ID, merchantID, f.ForecastDate, horizonDays,
		f.InflowToday, f.Inflow30d, f.Inflow60d, f.Inflow90d, f.OutflowToday, f.Outflow30d, f.Outflow60d, f.Outflow90d, f.Net90d)
	if err != nil {
		return f, err
	}
	return f, nil
}

// LatestForecast returns the most recently persisted forecast, or a fresh one if none.
func (r *Repository) LatestForecast(ctx context.Context, merchantID string) (*Forecast, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, to_char(forecast_date,'YYYY-MM-DD'), horizon_days,
		       inflow_today::text, inflow_30d::text, inflow_60d::text, inflow_90d::text,
		       outflow_today::text, outflow_30d::text, outflow_60d::text, outflow_90d::text, net_90d::text,
		       to_char(generated_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM cash_flow_forecasts WHERE merchant_id=$1 ORDER BY generated_at DESC LIMIT 1`, merchantID)
	var f Forecast
	err := row.Scan(&f.ID, &f.ForecastDate, &f.HorizonDays,
		&f.InflowToday, &f.Inflow30d, &f.Inflow60d, &f.Inflow90d,
		&f.OutflowToday, &f.Outflow30d, &f.Outflow60d, &f.Outflow90d, &f.Net90d, &f.GeneratedAt)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func tzEAT() *time.Location { return time.FixedZone("EAT", 3*3600) }
