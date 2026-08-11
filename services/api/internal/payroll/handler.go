// Package payroll handler entrypoint. Per-resource route handlers live in
// handler_{resource}.go (org, salary, employees, loans, runs, reports,
// settlement, portal, audit, calendar, leave, claims). This file defines the
// Handler type, its route wiring, and the shared helpers.

package payroll

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	// Org hierarchy
	r.Post("/departments", h.CreateDepartment)
	r.Get("/departments", h.ListDepartments)
	r.Post("/designations", h.CreateDesignation)
	r.Get("/designations", h.ListDesignations)
	r.Post("/grades", h.CreateGrade)
	r.Post("/branches", h.CreateBranch)
	r.Get("/branches", h.ListBranches)

	// Salary structures — CTC template enterprise-grade
	r.Post("/salary_structures", h.CreateSalaryStructure)
	r.Get("/salary_structures", h.ListSalaryStructures)
	r.Get("/salary_structures/{id}", h.GetSalaryStructure)

	// Employees enhanced
	r.Post("/employees", h.CreateEmployee)
	r.Post("/employees/bulk", h.BulkCreateEmployees) // CSV import 500 employees <2s p99
	r.Get("/employees", h.ListEmployees)
	r.Get("/employees/{id}", h.GetEmployee)
	r.Post("/employees/{id}/revisions", h.CreateSalaryRevision)
	r.Get("/employees/{id}/revisions", h.ListSalaryRevisions)
	r.Get("/employees/{id}/ytd", h.GetYTD)

	// Loans & Advances
	r.Post("/loans", h.CreateLoan)
	r.Get("/employees/{id}/loans", h.ListLoans)

	// Payroll runs comprehensive
	r.Post("/payroll_runs", h.CreateRun)
	r.Get("/payroll_runs", h.ListRuns)                             // additional list endpoint
	r.Post("/payroll_runs/{id}/attendance/bulk", h.BulkAttendance) // attendance + OT + LOP
	r.Post("/payroll_runs/{id}/variable_inputs/bulk", h.BulkVariableInputs)
	r.Post("/payroll_runs/{id}/calculate", h.Calculate)      // V2 formula engine proration
	r.Post("/payroll_runs/{id}/calculate/v2", h.CalculateV2) // explicit V2
	r.Post("/payroll_runs/{id}/approve", h.Approve)
	r.Post("/payroll_runs/{id}/disburse", h.Disburse)
	r.Get("/payroll_runs/{id}/items", h.ListItems)
	r.Get("/payroll_runs/{id}/payslips/{employee_id}/pdf", h.GetPayslipPDF)
	r.Get("/payroll_runs/{id}/payslips/bulk/zip", h.GetPayslipsZip)

	// Compliance reports — outstanding Beyond ApexPay
	r.Get("/payroll_reports/pension", h.GetPensionReport)
	r.Get("/payroll_reports/erca_withholding", h.GetERCAReport)
	r.Get("/payroll_reports/bank_disbursal", h.GetBankDisbursalReport)
	r.Get("/payroll_reports/cost_center", h.GetCostCenterReport)
	r.Get("/payroll_reports/annual_tax_certificate", h.GetAnnualTaxCertificate) // ERCA annual Form16 equivalent
	r.Get("/payroll_reports/payroll_register", h.GetPayrollRegister)
	r.Get("/payroll_reports/variance", h.GetVarianceReport)

	// Final settlement F&F
	r.Post("/final_settlements", h.CreateFinalSettlement)
	r.Get("/final_settlements", h.ListFinalSettlements)

	// Employee portal
	r.Post("/employee_portal/magic_link", h.CreateMagicLink)
	r.Get("/employee_portal/me", h.GetMyPortal)

	// Payroll audit
	r.Get("/payroll_audit_logs", h.ListAuditLogs)

	// Payroll Calendar — Ethiopia Business Practice Cutoff 25th Disbursal 30th Pay Last Day Lock After Disbursal
	r.Post("/calendars", h.CreateCalendar)
	r.Get("/calendars", h.ListCalendars)
	r.Get("/calendars/{id}", h.GetCalendar)
	r.Post("/calendars/{id}/lock", h.LockCalendar)
	r.Post("/calendars/{id}/unlock", h.UnlockCalendar)

	// Leave Management — Art 77 Annual 14+1 up to 35, Art 82 Sick 6 months, Art 86 Maternity 120 days
	r.Post("/leave_balances", h.CreateLeaveBalance)
	r.Get("/leave_balances", h.ListLeaveBalances)
	r.Post("/leave_requests", h.CreateLeaveRequest)
	r.Get("/leave_requests", h.ListLeaveRequests)
	r.Post("/leave_requests/{id}/approve", h.ApproveLeaveRequest)
	r.Post("/leave_requests/{id}/reject", h.RejectLeaveRequest)

	// Claims Enhanced — Receipt Upload MinIO Approval Manager->Finance
	r.Post("/claims", h.CreateClaim)
	r.Get("/claims", h.ListClaims)
	r.Post("/claims/{id}/approve/manager", h.ApproveClaimManager)
	r.Post("/claims/{id}/approve/finance", h.ApproveClaimFinance)

	// Loans EMI Schedule Repayment Tracking UI
	r.Get("/loans/{id}/emi_schedule", h.GetLoanEMISchedule)
	r.Get("/loans/{id}/repayments", h.GetLoanRepayments)
}

// ==================== Org Hierarchy Handlers ====================

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

var timeNow = func() time.Time { return time.Now() }

func nowTime() time.Time { return timeNow() }

func parseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, s)
}

func monthsBetween(from, to time.Time) int {
	if from.After(to) {
		return 0
	}
	years := to.Year() - from.Year()
	months := int(to.Month() - from.Month())
	return years*12 + months
}

func ptrDec(d decimal.Decimal) *decimal.Decimal { return &d }

// createdByPtr returns a pointer to s, or nil when s is empty (no authenticated user).
func createdByPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
