package payroll

import (
	"context"
	"encoding/json"
	"time"

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
	_, err := r.pool.Exec(ctx, `UPDATE payroll_claims SET approved_by_manager=$1, manager_approved_at=now(), status='approved_by_manager', updated_at=now() WHERE id=$2`, managerID, claimID)
	return err
}

func (r *PgRepository) ApproveClaimFinance(ctx context.Context, claimID, financeID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payroll_claims SET approved_by_finance=$1, finance_approved_at=now(), status='approved', updated_at=now() WHERE id=$2`, financeID, claimID)
	return err
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
