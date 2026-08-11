// List queries for payroll runs and final settlements. These were stubbed in the
// handlers; this file backs them with real DB reads so the merchant UI shows
// live data instead of empty arrays.
package payroll

import (
	"context"
	"encoding/json"
	"github.com/shopspring/decimal"
)

// ListRuns returns the merchant's payroll runs, newest period first.
func (r *PgRepository) ListRuns(ctx context.Context, merchantID string) ([]PayrollRun, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, run_ref, period_month, period_year, type, status,
		total_gross::text, total_deductions::text, total_net::text, total_tax::text, total_pension::text,
		COALESCE(employer_total_pension::text,'0'), COALESCE(total_employer_cost::text,'0'), total_count,
		COALESCE(total_employees_paid,0), COALESCE(total_employees_failed,0), book_id, payroll_data, variance_report,
		COALESCE(created_by,'') FROM payroll_runs WHERE merchant_id=$1 ORDER BY period_year DESC, period_month DESC, created_at DESC`,
		merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []PayrollRun
	for rows.Next() {
		pr := PayrollRun{}
		var gross, ded, net, tax, pension, employerPension, employerCost string
		var payrollDataStr, varianceStr string
		var totalCount, paid, failed int
		var createdBy string
		if err := rows.Scan(&pr.ID, &pr.MerchantID, &pr.RunRef, &pr.PeriodMonth, &pr.PeriodYear, &pr.Type, &pr.Status,
			&gross, &ded, &net, &tax, &pension, &employerPension, &employerCost, &totalCount, &paid, &failed,
			&pr.BookID, &payrollDataStr, &varianceStr, &createdBy); err != nil {
			return nil, err
		}
		pr.TotalGross, _ = decimal.NewFromString(gross)
		pr.TotalDeductions, _ = decimal.NewFromString(ded)
		pr.TotalNet, _ = decimal.NewFromString(net)
		pr.TotalTax, _ = decimal.NewFromString(tax)
		pr.TotalPension, _ = decimal.NewFromString(pension)
		pr.EmployerTotalPension, _ = decimal.NewFromString(employerPension)
		pr.TotalEmployerCost, _ = decimal.NewFromString(employerCost)
		pr.TotalCount = totalCount
		pr.TotalEmployeesPaid = paid
		pr.TotalEmployeesFailed = failed
		if createdBy != "" {
			pr.CreatedBy = &createdBy
		}
		_ = json.Unmarshal([]byte(payrollDataStr), &pr.PayrollData)
		_ = json.Unmarshal([]byte(varianceStr), &pr.VarianceReport)
		list = append(list, pr)
	}
	return list, rows.Err()
}

// ListFinalSettlements returns the merchant's final-settlement (F&F) requests,
// newest first.
func (r *PgRepository) ListFinalSettlements(ctx context.Context, merchantID string) ([]FinalSettlement, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, employee_id, resignation_date, last_working_date,
		notice_period_days, notice_served_days, notice_shortfall_days,
		COALESCE(leave_encashment_days::text,'0'), COALESCE(leave_encashment_amount::text,'0'),
		COALESCE(severance_amount::text,'0'), COALESCE(gratuity_amount::text,'0'), COALESCE(bonus_pro_rata::text,'0'),
		COALESCE(outstanding_loans::text,'0'), COALESCE(outstanding_advances::text,'0'),
		COALESCE(other_earnings::text,'0'), COALESCE(other_deductions::text,'0'),
		COALESCE(total_payable::text,'0'), COALESCE(total_deductions::text,'0'), COALESCE(net_payable::text,'0'),
		status, COALESCE(clearance_checklist,'[]'::jsonb), COALESCE(clearance_items_detailed,'[]'::jsonb),
		COALESCE(assets_returned,'[]'::jsonb), COALESCE(exit_interview,'{}'::jsonb), created_at
		FROM payroll_final_settlements WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []FinalSettlement
	for rows.Next() {
		fs := FinalSettlement{}
		var encDays, encAmt, sev, grat, bonus, loans, adv, otherEarn, otherDed, totalPay, totalDed, netPay string
		var checklistJSON, detailedJSON, assetsJSON, exitJSON string
		if err := rows.Scan(&fs.ID, &fs.MerchantID, &fs.EmployeeID, &fs.ResignationDate, &fs.LastWorkingDate,
			&fs.NoticePeriodDays, &fs.NoticeServedDays, &fs.NoticeShortfallDays,
			&encDays, &encAmt, &sev, &grat, &bonus, &loans, &adv, &otherEarn, &otherDed,
			&totalPay, &totalDed, &netPay, &fs.Status,
			&checklistJSON, &detailedJSON, &assetsJSON, &exitJSON, &fs.CreatedAt); err != nil {
			return nil, err
		}
		fs.LeaveEncashmentDays, _ = decimal.NewFromString(encDays)
		fs.LeaveEncashmentAmount, _ = decimal.NewFromString(encAmt)
		fs.SeveranceAmount, _ = decimal.NewFromString(sev)
		fs.GratuityAmount, _ = decimal.NewFromString(grat)
		fs.BonusProRata, _ = decimal.NewFromString(bonus)
		fs.OutstandingLoans, _ = decimal.NewFromString(loans)
		fs.OutstandingAdvances, _ = decimal.NewFromString(adv)
		fs.OtherEarnings, _ = decimal.NewFromString(otherEarn)
		fs.OtherDeductions, _ = decimal.NewFromString(otherDed)
		fs.TotalPayable, _ = decimal.NewFromString(totalPay)
		fs.TotalDeductions, _ = decimal.NewFromString(totalDed)
		fs.NetPayable, _ = decimal.NewFromString(netPay)
		_ = json.Unmarshal([]byte(checklistJSON), &fs.ClearanceChecklist)
		// structured clearance fields are kept in these fields; unmarshal best-effort
		_ = json.Unmarshal([]byte(detailedJSON), &fs.ClearanceDetailed)
		_ = json.Unmarshal([]byte(assetsJSON), &fs.AssetsReturned)
		_ = json.Unmarshal([]byte(exitJSON), &fs.ExitInterview)
		list = append(list, fs)
	}
	return list, rows.Err()
}
