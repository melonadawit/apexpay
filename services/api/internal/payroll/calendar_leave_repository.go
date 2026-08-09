package payroll

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/ledger"
	"github.com/shopspring/decimal"
)

// ==================== Payroll Calendar CRUD per Ethiopia Business Practice ====================

func (r *PgRepository) CreateCalendar(ctx context.Context, cal *PayrollCalendar) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO payroll_calendars (id, merchant_id, name, description, pay_frequency, year, month, cutoff_day, disbursal_day, pay_day, cutoff_date, disbursal_date, pay_date, is_locked, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		cal.ID, cal.MerchantID, cal.Name, cal.Description, cal.PayFrequency, cal.Year, cal.Month, cal.CutoffDay, cal.DisbursalDay, cal.PayDay, cal.CutoffDate, cal.DisbursalDate, cal.PayDate, cal.IsLocked, cal.CreatedBy)
	return err
}

func (r *PgRepository) ListCalendars(ctx context.Context, merchantID string, year int) ([]PayrollCalendar, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, name, COALESCE(description,''), pay_frequency, year, month, cutoff_day, disbursal_day, pay_day, cutoff_date, disbursal_date, pay_date, is_locked, locked_at, locked_by, created_at FROM payroll_calendars WHERE merchant_id=$1 AND year=$2 ORDER BY month ASC`, merchantID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []PayrollCalendar
	for rows.Next() {
		var c PayrollCalendar
		var month *int
		if err := rows.Scan(&c.ID, &c.MerchantID, &c.Name, &c.Description, &c.PayFrequency, &c.Year, &month, &c.CutoffDay, &c.DisbursalDay, &c.PayDay, &c.CutoffDate, &c.DisbursalDate, &c.PayDate, &c.IsLocked, &c.LockedAt, &c.LockedBy, &c.CreatedAt); err != nil {
			return nil, err
		}
		if month != nil {
			c.Month = month
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *PgRepository) GetCalendar(ctx context.Context, merchantID, calendarID string) (*PayrollCalendar, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, name, COALESCE(description,''), pay_frequency, year, month, cutoff_day, disbursal_day, pay_day, cutoff_date, disbursal_date, pay_date, is_locked, locked_at, locked_by, created_at FROM payroll_calendars WHERE merchant_id=$1 AND id=$2`, merchantID, calendarID)
	var c PayrollCalendar
	var month *int
	err := row.Scan(&c.ID, &c.MerchantID, &c.Name, &c.Description, &c.PayFrequency, &c.Year, &month, &c.CutoffDay, &c.DisbursalDay, &c.PayDay, &c.CutoffDate, &c.DisbursalDate, &c.PayDate, &c.IsLocked, &c.LockedAt, &c.LockedBy, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	if month != nil {
		c.Month = month
	}
	return &c, nil
}

func (r *PgRepository) LockCalendar(ctx context.Context, merchantID, calendarID, lockedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_calendars SET is_locked=true, locked_at=now(), locked_by=$3, updated_at=now() WHERE merchant_id=$1 AND id=$2`, merchantID, calendarID, lockedBy)
	return err
}

func (r *PgRepository) UnlockCalendar(ctx context.Context, merchantID, calendarID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_calendars SET is_locked=false, locked_at=NULL, locked_by=NULL, updated_at=now() WHERE merchant_id=$1 AND id=$2`, merchantID, calendarID)
	return err
}

// ==================== Leave Balances per Ethiopia Law Art 77/82/86 ====================

func (r *PgRepository) CreateLeaveBalance(ctx context.Context, balance *LeaveBalance) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO payroll_leave_balances (id, merchant_id, employee_id, leave_type, year, entitled_days, used_days, remaining_days, carry_forward_days) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (merchant_id, employee_id, leave_type, year) DO UPDATE SET entitled_days=$6, used_days=$7, remaining_days=$8, carry_forward_days=$9, updated_at=now()`,
		balance.ID, balance.MerchantID, balance.EmployeeID, balance.LeaveType, balance.Year, balance.EntitledDays.String(), balance.UsedDays.String(), balance.RemainingDays.String(), balance.CarryForwardDays.String())
	return err
}

func (r *PgRepository) GetLeaveBalance(ctx context.Context, merchantID, employeeID string, leaveType LeaveType, year int) (*LeaveBalance, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, employee_id, leave_type, year, entitled_days::text, used_days::text, remaining_days::text, carry_forward_days::text FROM payroll_leave_balances WHERE merchant_id=$1 AND employee_id=$2 AND leave_type=$3 AND year=$4`, merchantID, employeeID, leaveType, year)
	var b LeaveBalance
	var entitled, used, remaining, carry string
	err := row.Scan(&b.ID, &b.MerchantID, &b.EmployeeID, &b.LeaveType, &b.Year, &entitled, &used, &remaining, &carry)
	if err != nil {
		return nil, err
	}
	b.EntitledDays, _ = decimal.NewFromString(entitled)
	b.UsedDays, _ = decimal.NewFromString(used)
	b.RemainingDays, _ = decimal.NewFromString(remaining)
	b.CarryForwardDays, _ = decimal.NewFromString(carry)
	return &b, nil
}

func (r *PgRepository) UpdateLeaveBalance(ctx context.Context, balance *LeaveBalance) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_leave_balances SET entitled_days=$1, used_days=$2, remaining_days=$3, carry_forward_days=$4, updated_at=now() WHERE id=$5`,
		balance.EntitledDays.String(), balance.UsedDays.String(), balance.RemainingDays.String(), balance.CarryForwardDays.String(), balance.ID)
	return err
}

func (r *PgRepository) ListLeaveBalancesByEmployee(ctx context.Context, merchantID, employeeID string, year int) ([]LeaveBalance, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, employee_id, leave_type, year, entitled_days::text, used_days::text, remaining_days::text, carry_forward_days::text FROM payroll_leave_balances WHERE merchant_id=$1 AND employee_id=$2 AND year=$3`, merchantID, employeeID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []LeaveBalance
	for rows.Next() {
		var b LeaveBalance
		var entitled, used, remaining, carry string
		if err := rows.Scan(&b.ID, &b.MerchantID, &b.EmployeeID, &b.LeaveType, &b.Year, &entitled, &used, &remaining, &carry); err != nil {
			return nil, err
		}
		b.EntitledDays, _ = decimal.NewFromString(entitled)
		b.UsedDays, _ = decimal.NewFromString(used)
		b.RemainingDays, _ = decimal.NewFromString(remaining)
		b.CarryForwardDays, _ = decimal.NewFromString(carry)
		list = append(list, b)
	}
	return list, nil
}

// ==================== Leave Requests ====================

func (r *PgRepository) CreateLeaveRequest(ctx context.Context, req *LeaveRequest) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO payroll_leave_requests (id, merchant_id, employee_id, leave_type, start_date, end_date, days_requested, reason, status, medical_certificate_file_key) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		req.ID, req.MerchantID, req.EmployeeID, req.LeaveType, req.StartDate, req.EndDate, req.DaysRequested.String(), req.Reason, req.Status, req.MedicalCertificateFileKey)
	return err
}

func (r *PgRepository) GetLeaveRequest(ctx context.Context, merchantID, requestID string) (*LeaveRequest, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, employee_id, leave_type, start_date, end_date, days_requested::text, COALESCE(reason,''), status, approved_by, approved_at, COALESCE(rejection_reason,''), medical_certificate_file_key FROM payroll_leave_requests WHERE merchant_id=$1 AND id=$2`, merchantID, requestID)
	var req LeaveRequest
	var days string
	err := row.Scan(&req.ID, &req.MerchantID, &req.EmployeeID, &req.LeaveType, &req.StartDate, &req.EndDate, &days, &req.Reason, &req.Status, &req.ApprovedBy, &req.ApprovedAt, &req.RejectionReason, &req.MedicalCertificateFileKey)
	if err != nil {
		return nil, err
	}
	req.DaysRequested, _ = decimal.NewFromString(days)
	return &req, nil
}

func (r *PgRepository) ListLeaveRequests(ctx context.Context, merchantID, employeeID string, year int, status *LeaveStatus) ([]LeaveRequest, error) {
	query := `SELECT id, merchant_id, employee_id, leave_type, start_date, end_date, days_requested::text, COALESCE(reason,''), status FROM payroll_leave_requests WHERE merchant_id=$1 AND EXTRACT(YEAR FROM start_date)=$2`
	args := []interface{}{merchantID, year}
	if employeeID != "" {
		query += ` AND employee_id=$3`
		args = append(args, employeeID)
	}
	if status != nil {
		query += ` AND status=$4`
		args = append(args, *status)
	}
	query += ` ORDER BY start_date DESC`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []LeaveRequest
	for rows.Next() {
		var req LeaveRequest
		var days string
		if err := rows.Scan(&req.ID, &req.MerchantID, &req.EmployeeID, &req.LeaveType, &req.StartDate, &req.EndDate, &days, &req.Reason, &req.Status); err != nil {
			return nil, err
		}
		req.DaysRequested, _ = decimal.NewFromString(days)
		list = append(list, req)
	}
	return list, nil
}

func (r *PgRepository) UpdateLeaveRequestStatus(ctx context.Context, requestID string, status LeaveStatus, approvedBy *string, rejectionReason string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_leave_requests SET status=$1, approved_by=$2, approved_at=now(), rejection_reason=$3, updated_at=now() WHERE id=$4`,
		status, approvedBy, rejectionReason, requestID)
	return err
}

// ==================== Loan EMI Schedule ====================

func (r *PgRepository) CreateLoanEMIScheduleBulk(ctx context.Context, schedules []LoanEMISchedule) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	for _, s := range schedules {
		_, err := tx.Exec(ctx, `INSERT INTO payroll_loan_emi_schedule (id, loan_id, installment_no, due_date, emi_amount, principal_component, interest_component, outstanding_after, status, run_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT (loan_id, installment_no) DO NOTHING`,
			s.ID, s.LoanID, s.InstallmentNo, s.DueDate, s.EMIAmount.String(), s.PrincipalComponent.String(), s.InterestComponent.String(), s.OutstandingAfter.String(), s.Status, s.RunID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PgRepository) ListEMIScheduleByLoan(ctx context.Context, loanID string) ([]LoanEMISchedule, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, loan_id, installment_no, due_date, emi_amount::text, principal_component::text, interest_component::text, outstanding_after::text, status, paid_at, run_id FROM payroll_loan_emi_schedule WHERE loan_id=$1 ORDER BY installment_no ASC`, loanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []LoanEMISchedule
	for rows.Next() {
		var s LoanEMISchedule
		var emi, principal, interest, outstanding string
		if err := rows.Scan(&s.ID, &s.LoanID, &s.InstallmentNo, &s.DueDate, &emi, &principal, &interest, &outstanding, &s.Status, &s.PaidAt, &s.RunID); err != nil {
			return nil, err
		}
		s.EMIAmount, _ = decimal.NewFromString(emi)
		s.PrincipalComponent, _ = decimal.NewFromString(principal)
		s.InterestComponent, _ = decimal.NewFromString(interest)
		s.OutstandingAfter, _ = decimal.NewFromString(outstanding)
		list = append(list, s)
	}
	return list, nil
}

// ==================== Claims Enhanced — Receipt Upload MinIO ====================

func (r *PgRepository) CreateClaimEnhanced(ctx context.Context, claim *ClaimEnhanced) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO payroll_claims (id, merchant_id, employee_id, claim_type, amount, description, receipt_file_key, receipt_file_hash, status, is_taxable, is_pensionable) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		claim.ID, claim.MerchantID, claim.EmployeeID, claim.ClaimType, claim.Amount.String(), claim.Description, claim.ReceiptFileKey, claim.ReceiptFileHash, claim.Status, claim.IsTaxable, claim.IsPensionable)
	return err
}

func (r *PgRepository) ListClaimsByEmployee(ctx context.Context, merchantID, employeeID string, status *string) ([]ClaimEnhanced, error) {
	query := `SELECT id, merchant_id, employee_id, claim_type, amount::text, COALESCE(description,''), receipt_file_key, receipt_file_hash, status, is_taxable, is_pensionable FROM payroll_claims WHERE merchant_id=$1 AND employee_id=$2`
	args := []interface{}{merchantID, employeeID}
	if status != nil {
		query += ` AND status=$3`
		args = append(args, *status)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ClaimEnhanced
	for rows.Next() {
		var c ClaimEnhanced
		var amount string
		if err := rows.Scan(&c.ID, &c.MerchantID, &c.EmployeeID, &c.ClaimType, &amount, &c.Description, &c.ReceiptFileKey, &c.ReceiptFileHash, &c.Status, &c.IsTaxable, &c.IsPensionable); err != nil {
			return nil, err
		}
		c.Amount, _ = decimal.NewFromString(amount)
		list = append(list, c)
	}
	return list, nil
}

func (r *PgRepository) ApproveClaimManager(ctx context.Context, claimID, managerID string) error {
	// Use NULL for the approver when no dashboard user is present (API-key actors have no
	// session user and the column has an FK to users(id)).
	_, err := r.pool.Exec(ctx, `UPDATE payroll_claims SET approved_by_manager=$1, manager_approved_at=now(), status='approved_by_manager', updated_at=now() WHERE id=$2`, nilStrApprov(managerID), claimID)
	return err
}

func (r *PgRepository) ApproveClaimFinance(ctx context.Context, claimID, financeID string) error {
	// Read the claim's merchant + amount so we can post the reimbursement to the GL.
	var merchantID, amountStr string
	if err := r.pool.QueryRow(ctx, `SELECT merchant_id, amount::text FROM payroll_claims WHERE id=$1`, claimID).Scan(&merchantID, &amountStr); err != nil {
		return err
	}

	// Mark the claim approved, then post the expense to the operating ledger in the same
	// logical step. Debit expense:operating, credit liability:payable (reimbursement owed).
	if _, err := r.pool.Exec(ctx, `UPDATE payroll_claims SET approved_by_finance=$1, finance_approved_at=now(), status='approved', updated_at=now() WHERE id=$2`, nilStrApprov(financeID), claimID); err != nil {
		return err
	}
	if r.ledger != nil {
		amount, _ := decimal.NewFromString(amountStr)
		if amount.GreaterThan(decimal.Zero) {
			if err := r.postClaimReimbursement(ctx, merchantID, claimID, amount); err != nil {
				return err
			}
		}
	}
	return nil
}

// postClaimReimbursement posts an approved expense claim to the merchant's operating ledger
// (debit expense:operating, credit liability:payable). Idempotent per claim via the posting
// key so re-approvals never double-post.
func (r *PgRepository) postClaimReimbursement(ctx context.Context, merchantID, claimID string, amount decimal.Decimal) error {
	var bookID, expenseID, payableID string
	err := r.pool.QueryRow(ctx, `
		SELECT lb.id,
			MAX(la.id) FILTER (WHERE la.code='expense:operating'),
			MAX(la.id) FILTER (WHERE la.code='liability:payable')
		FROM ledger_books lb JOIN ledger_accounts la ON la.book_id=lb.id
		WHERE lb.merchant_id=$1 AND lb.book_type='merchant_operating' AND lb.status='open'
		GROUP BY lb.id ORDER BY lb.id LIMIT 1`, merchantID).Scan(&bookID, &expenseID, &payableID)
	if err != nil || expenseID == "" || payableID == "" {
		return fmt.Errorf("claim ledger accounts unavailable: %w", err)
	}
	journalID := id.NewLedgerJournal()
	journal := &ledger.Journal{
		ID: journalID, BookID: bookID, PostingKey: "claim_" + claimID,
		Memo: "Expense claim reimbursement " + claimID, ReferenceType: "expense_claim", ReferenceID: claimID,
	}
	entries := []ledger.Entry{
		{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: expenseID, Direction: "debit", Amount: amount, Currency: "ETB"},
		{ID: id.New("lent"), JournalID: journalID, BookID: bookID, AccountID: payableID, Direction: "credit", Amount: amount, Currency: "ETB"},
	}
	return r.ledger.PostJournalTx(ctx, journal, entries)
}

// ==================== Escrow Accounts Automated Marketplace P2P Hold & Release ====================

func (r *PgRepository) CreateEscrowAgreement(ctx context.Context, agreement *EscrowAgreement) error {
	conditions, _ := json.Marshal(agreement.Conditions)
	_, err := r.pool.Exec(ctx, `INSERT INTO escrow_agreements (id, merchant_id, agreement_number, title, description, buyer_merchant_id, seller_merchant_id, amount, currency, platform_fee_percent, withholding_tax_percent, conditions, auto_release, auto_release_after_days, status, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		agreement.ID, agreement.MerchantID, agreement.AgreementNumber, agreement.Title, agreement.Description, agreement.BuyerMerchantID, agreement.SellerMerchantID, agreement.Amount.String(), agreement.Currency, agreement.PlatformFeePercent.String(), agreement.WithholdingTaxPercent.String(), string(conditions), agreement.AutoRelease, agreement.AutoReleaseAfterDays, agreement.Status, agreement.CreatedBy)
	return err
}

func (r *PgRepository) GetEscrowAgreement(ctx context.Context, merchantID, agreementID string) (*EscrowAgreement, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, agreement_number, title, COALESCE(description,''), buyer_merchant_id, seller_merchant_id, amount::text, currency, platform_fee_percent::text, withholding_tax_percent::text, conditions, auto_release, auto_release_after_days, status FROM escrow_agreements WHERE merchant_id=$1 AND id=$2`, merchantID, agreementID)
	var agr EscrowAgreement
	var amount, feePct, taxPct, conditionsStr string
	err := row.Scan(&agr.ID, &agr.MerchantID, &agr.AgreementNumber, &agr.Title, &agr.Description, &agr.BuyerMerchantID, &agr.SellerMerchantID, &amount, &agr.Currency, &feePct, &taxPct, &conditionsStr, &agr.AutoRelease, &agr.AutoReleaseAfterDays, &agr.Status)
	if err != nil {
		return nil, err
	}
	agr.Amount, _ = decimal.NewFromString(amount)
	agr.PlatformFeePercent, _ = decimal.NewFromString(feePct)
	agr.WithholdingTaxPercent, _ = decimal.NewFromString(taxPct)
	_ = json.Unmarshal([]byte(conditionsStr), &agr.Conditions)
	return &agr, nil
}

func (r *PgRepository) CreateEscrowAccountTx(ctx context.Context, escrow *EscrowAccount, journal *ledger.Journal, entries []ledger.Entry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, `INSERT INTO escrow_accounts (id, merchant_id, agreement_id, account_number, account_name, amount, currency, status, held_at, buyer_merchant_id, seller_merchant_id, order_id, order_amount, platform_fee, seller_amount, withholding_tax, ledger_book_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
		escrow.ID, escrow.MerchantID, escrow.AgreementID, escrow.AccountNumber, escrow.AccountName, escrow.Amount.String(), escrow.Currency, escrow.Status, escrow.HeldAt, escrow.BuyerMerchantID, escrow.SellerMerchantID, escrow.OrderID, escrow.OrderAmount.String(), escrow.PlatformFee.String(), escrow.SellerAmount.String(), escrow.WithholdingTax.String(), escrow.LedgerBookID)
	if err != nil {
		return err
	}
	// Ledger
	_, err = tx.Exec(ctx, `INSERT INTO ledger_books (id, merchant_id, book_type, name, currency, status) VALUES ($1,$2,'escrow',$3,'ETB','open') ON CONFLICT (id) DO NOTHING`, journal.BookID, escrow.MerchantID, "Escrow "+escrow.AccountNumber)
	if err != nil {
		return err
	}
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

func (r *PgRepository) GetEscrowAccount(ctx context.Context, merchantID, escrowID string) (*EscrowAccount, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, agreement_id, account_number, account_name, amount::text, currency, status, held_at, buyer_merchant_id, seller_merchant_id, order_id, order_amount::text, platform_fee::text, seller_amount::text, withholding_tax::text, ledger_book_id FROM escrow_accounts WHERE merchant_id=$1 AND id=$2`, merchantID, escrowID)
	var esc EscrowAccount
	var amount, orderAmount, fee, sellerAmt, withholding string
	err := row.Scan(&esc.ID, &esc.MerchantID, &esc.AgreementID, &esc.AccountNumber, &esc.AccountName, &amount, &esc.Currency, &esc.Status, &esc.HeldAt, &esc.BuyerMerchantID, &esc.SellerMerchantID, &esc.OrderID, &orderAmount, &fee, &sellerAmt, &withholding, &esc.LedgerBookID)
	if err != nil {
		return nil, err
	}
	esc.Amount, _ = decimal.NewFromString(amount)
	esc.OrderAmount, _ = decimal.NewFromString(orderAmount)
	esc.PlatformFee, _ = decimal.NewFromString(fee)
	esc.SellerAmount, _ = decimal.NewFromString(sellerAmt)
	esc.WithholdingTax, _ = decimal.NewFromString(withholding)
	return &esc, nil
}

func (r *PgRepository) ReleaseEscrowTx(ctx context.Context, escrowID string, journal *ledger.Journal, entries []ledger.Entry, releaserID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, `UPDATE escrow_accounts SET status='released', release_at=now(), updated_at=now() WHERE id=$1`, escrowID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO ledger_journals (id, book_id, posting_key, memo, reference_type, reference_id) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (book_id, posting_key) DO NOTHING`, journal.ID, journal.BookID, journal.PostingKey, journal.Memo, journal.ReferenceType, journal.ReferenceID)
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

func (r *PgRepository) ReturnEscrowTx(ctx context.Context, escrowID string, journal *ledger.Journal, entries []ledger.Entry, returnerID, reason string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, `UPDATE escrow_accounts SET status='returned', return_at=now(), updated_at=now() WHERE id=$1`, escrowID)
	if err != nil {
		return err
	}
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

func (r *PgRepository) ListExpiredEscrowsForAutoRelease(ctx context.Context) ([]EscrowAccount, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, account_number, amount::text FROM escrow_accounts WHERE status='held' AND expires_at <= now() AND (SELECT auto_release FROM escrow_agreements WHERE escrow_agreements.id = escrow_accounts.agreement_id) = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []EscrowAccount
	for rows.Next() {
		var esc EscrowAccount
		var amount string
		if err := rows.Scan(&esc.ID, &esc.MerchantID, &esc.AccountNumber, &amount); err != nil {
			return nil, err
		}
		esc.Amount, _ = decimal.NewFromString(amount)
		list = append(list, esc)
	}
	return list, nil
}

// ==================== Payout Links Enhanced QR + Scan & Pay ====================

func (r *PgRepository) CreateEnhancedPayoutLink(ctx context.Context, link *EnhancedPayoutLink) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO payout_links_enhanced (id, merchant_id, amount, currency, public_token, qr_code_data, recipient_name, recipient_phone, recipient_email, purpose, status, expires_at, escrow_book_id, ledger_book_id, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		link.ID, link.MerchantID, link.Amount.String(), link.Currency, link.PublicToken, link.QRCodeData, link.RecipientName, link.RecipientPhone, link.RecipientEmail, link.Purpose, link.Status, link.ExpiresAt, link.EscrowBookID, link.LedgerBookID, link.CreatedBy)
	return err
}

func (r *PgRepository) GetEnhancedPayoutLinkByToken(ctx context.Context, publicToken string) (*EnhancedPayoutLink, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, amount::text, currency, public_token, COALESCE(qr_code_data,''), COALESCE(recipient_name,''), COALESCE(recipient_phone,''), COALESCE(recipient_email,''), COALESCE(purpose,''), status, expires_at, claimed_at, beneficiary_id, escrow_book_id, ledger_book_id FROM payout_links_enhanced WHERE public_token=$1`, publicToken)
	var link EnhancedPayoutLink
	var amount string
	err := row.Scan(&link.ID, &link.MerchantID, &amount, &link.Currency, &link.PublicToken, &link.QRCodeData, &link.RecipientName, &link.RecipientPhone, &link.RecipientEmail, &link.Purpose, &link.Status, &link.ExpiresAt, &link.ClaimedAt, &link.BeneficiaryID, &link.EscrowBookID, &link.LedgerBookID)
	if err != nil {
		return nil, err
	}
	link.Amount, _ = decimal.NewFromString(amount)
	return &link, nil
}

func (r *PgRepository) ClaimEnhancedPayoutLinkTx(ctx context.Context, linkID, beneficiaryID string, journal *ledger.Journal, entries []ledger.Entry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, `UPDATE payout_links_enhanced SET status='claimed', claimed_at=now(), beneficiary_id=$2, updated_at=now() WHERE id=$1`, linkID, beneficiaryID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO ledger_journals (id, book_id, posting_key, memo, reference_type, reference_id) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (book_id, posting_key) DO NOTHING`, journal.ID, journal.BookID, journal.PostingKey, journal.Memo, journal.ReferenceType, journal.ReferenceID)
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

// ==================== Helpers for JSON ====================

func toJSONLeave(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

var _ = time.Now
var _ = toJSONLeave

// ApprovedLeaveDaysForPeriod returns the total approved leave days per employee whose
// leave falls within the given month/year. Used by CalculateRun to ensure approved leave
// counts as paid days (not LOP) automatically.
func (r *PgRepository) ApprovedLeaveDaysForPeriod(ctx context.Context, merchantID string, month, year int) (map[string]decimal.Decimal, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT employee_id, COALESCE(SUM(days_requested),0)::text
		FROM payroll_leave_requests
		WHERE merchant_id=$1 AND status='approved'
		  AND EXTRACT(MONTH FROM start_date)=$2 AND EXTRACT(YEAR FROM start_date)=$3
		GROUP BY employee_id`, merchantID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]decimal.Decimal{}
	for rows.Next() {
		var empID, days string
		if err := rows.Scan(&empID, &days); err != nil {
			return nil, err
		}
		d, _ := decimal.NewFromString(days)
		out[empID] = d
	}
	return out, rows.Err()
}

// nilStrApprov returns NULL for an empty approver id so FK-constrained approver columns
// accept API-key actors who have no dashboard user.
func nilStrApprov(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
