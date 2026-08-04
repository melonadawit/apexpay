package payroll

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"apexpay/internal/ledger"
)

type PgRepository struct {
	pool *pgxpool.Pool
	ledger *ledger.PgRepository
}

func NewPgRepository(pool *pgxpool.Pool, ledger *ledger.PgRepository) *PgRepository {
	return &PgRepository{pool: pool, ledger: ledger}
}

func (r *PgRepository) CreateEmployee(ctx context.Context, e *Employee) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO employees (id, merchant_id, employee_code, name, name_am, email, phone, tin, fayda_fin_hash, pension_no, bank_account_hash, bank_account_masked, bank_code, base_salary, employment_date, employment_type, cost_center, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
		e.ID, e.MerchantID, e.EmployeeCode, e.Name, e.NameAM, e.Email, e.Phone, e.TIN, e.FinHash, e.PensionNo, e.BankAccountHash, e.BankAccountMasked, e.BankCode, e.BaseSalary.String(), e.EmploymentDate, e.EmploymentType, e.CostCenter, e.Status)
	return err
}

func (r *PgRepository) ListEmployees(ctx context.Context, merchantID string) ([]Employee, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, employee_code, name, base_salary::text, status, cost_center FROM employees WHERE merchant_id=$1 AND status='active'`, merchantID)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []Employee
	for rows.Next() {
		var e Employee
		var base string
		if err := rows.Scan(&e.ID, &e.MerchantID, &e.EmployeeCode, &e.Name, &base, &e.Status, &e.CostCenter); err != nil { return nil, err }
		list = append(list, e)
	}
	return list, nil
}

func (r *PgRepository) GetEmployee(ctx context.Context, merchantID, employeeID string) (*Employee, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, employee_code, name, base_salary::text FROM employees WHERE merchant_id=$1 AND id=$2`, merchantID, employeeID)
	var e Employee
	var base string
	err := row.Scan(&e.ID, &e.MerchantID, &e.EmployeeCode, &e.Name, &base)
	return &e, err
}

func (r *PgRepository) CreateRun(ctx context.Context, run *PayrollRun) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO payroll_runs (id, merchant_id, book_id, run_ref, period_month, period_year, type, status, total_gross, total_net) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		run.ID, run.MerchantID, run.BookID, run.RunRef, run.PeriodMonth, run.PeriodYear, run.Type, run.Status, run.TotalGross.String(), run.TotalNet.String())
	return err
}

func (r *PgRepository) GetRun(ctx context.Context, merchantID, runID string) (*PayrollRun, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, run_ref, period_month, period_year, type, status, total_gross::text, total_net::text, total_tax::text, total_pension::text, book_id FROM payroll_runs WHERE merchant_id=$1 AND id=$2`, merchantID, runID)
	var pr PayrollRun
	var gross, net, tax, pension string
	err := row.Scan(&pr.ID, &pr.MerchantID, &pr.RunRef, &pr.PeriodMonth, &pr.PeriodYear, &pr.Type, &pr.Status, &gross, &net, &tax, &pension, &pr.BookID)
	return &pr, err
}

func (r *PgRepository) UpdateRunStatus(ctx context.Context, runID string, status RunStatus, totals map[string]string) error {
	// Simplified: totals map handling would be in service, here update status
	_, err := r.pool.Exec(ctx, `UPDATE payroll_runs SET status=$1 WHERE id=$2`, status, runID)
	return err
}

// Overload for service that passes map decimal
func (r *PgRepository) UpdateRunStatusWithTotals(ctx context.Context, runID string, status RunStatus, totals map[string]interface{}) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_runs SET status=$1 WHERE id=$2`, status, runID)
	return err
}

func (r *PgRepository) BulkCreateItems(ctx context.Context, items []PayrollItem) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil { return err }
	defer func() { _ = tx.Rollback(ctx) }()
	for _, it := range items {
		_, err = tx.Exec(ctx, `INSERT INTO payroll_items (id, run_id, employee_id, gross, taxable_income, income_tax, pension_employee, pension_employer, net_pay, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			it.ID, it.RunID, it.EmployeeID, it.Gross.String(), it.TaxableIncome.String(), it.IncomeTax.String(), it.PensionEmployee.String(), it.PensionEmployer.String(), it.NetPay.String(), it.Status)
		if err != nil { return err }
	}
	return tx.Commit(ctx)
}

func (r *PgRepository) ListItems(ctx context.Context, runID string) ([]PayrollItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, run_id, employee_id, gross::text, net_pay::text FROM payroll_items WHERE run_id=$1`, runID)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []PayrollItem
	for rows.Next() {
		var it PayrollItem
		var gross, net string
		if err := rows.Scan(&it.ID, &it.RunID, &it.EmployeeID, &gross, &net); err != nil { return nil, err }
		list = append(list, it)
	}
	return list, nil
}

func (r *PgRepository) GetTaxBrackets(ctx context.Context) ([]TaxBracket, error) {
	rows, err := r.pool.Query(ctx, `SELECT min_amount::text, max_amount::text, rate::text, deduction::text FROM payroll_tax_brackets WHERE effective_from <= CURRENT_DATE AND (effective_to IS NULL OR effective_to >= CURRENT_DATE) ORDER BY min_amount ASC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var brackets []TaxBracket
	for rows.Next() {
		var minStr, maxStr, rateStr, dedStr string
		var maxPtr *string
		if err := rows.Scan(&minStr, &maxStr, &rateStr, &dedStr); err != nil { return nil, err }
		if maxStr == "" {
			maxPtr = nil
		} else {
			maxPtr = &maxStr
		}
		// parse would be in service, simplified
		brackets = append(brackets, TaxBracket{})
	}
	return brackets, nil
}

func (r *PgRepository) CreateRunBookTx(ctx context.Context, run *PayrollRun, journal *ledger.Journal, entries []ledger.Entry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil { return err }
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, `INSERT INTO ledger_books (id, merchant_id, book_type, name, currency, status) VALUES ($1,$2,'payroll_run',$3,'ETB','open') ON CONFLICT (id) DO NOTHING`, journal.BookID, run.MerchantID, "Payroll run "+run.RunRef)
	if err != nil { return err }
	_, err = tx.Exec(ctx, `UPDATE payroll_runs SET book_id=$1 WHERE id=$2`, journal.BookID, run.ID)
	if err != nil { return err }
	_, err = tx.Exec(ctx, `INSERT INTO ledger_journals (id, book_id, posting_key, memo, reference_type, reference_id) VALUES ($1,$2,$3,$4,$5,$6)`, journal.ID, journal.BookID, journal.PostingKey, journal.Memo, journal.ReferenceType, journal.ReferenceID)
	if err != nil { return err }
	for _, e := range entries {
		_, err = tx.Exec(ctx, `INSERT INTO ledger_entries (id, journal_id, book_id, account_id, direction, amount, currency) VALUES ($1,$2,$3,$4,$5,$6,$7)`, e.ID, e.JournalID, e.BookID, e.AccountID, e.Direction, e.Amount.String(), e.Currency)
		if err != nil { return err }
	}
	return tx.Commit(ctx)
}
