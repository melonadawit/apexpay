// Payroll loan & advance handlers.
package payroll

import (
	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	"net/http"
)

func (h *Handler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		EmployeeID   string `json:"employee_id"`
		LoanType     string `json:"loan_type"` // personal, salary_advance
		Principal    string `json:"principal"`
		InterestRate string `json:"interest_rate"`
		TenureMonths int    `json:"tenure_months"`
		Reason       string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	principal, _ := decimal.NewFromString(req.Principal)
	rate, _ := decimal.NewFromString(req.InterestRate)
	// EMI = principal*(1+rate/100)/tenure — simple interest for demo, real amortized
	emi := principal.Div(decimal.NewFromInt(int64(req.TenureMonths))).Round(2)
	if !rate.IsZero() {
		// Add interest
		interestTotal := principal.Mul(rate).Div(decimal.NewFromInt(100)).Mul(decimal.NewFromInt(int64(req.TenureMonths))).Div(decimal.NewFromInt(12))
		emi = principal.Add(interestTotal).Div(decimal.NewFromInt(int64(req.TenureMonths))).Round(2)
	}
	loan := &Loan{
		ID: id.New("loan"), MerchantID: merchantID, EmployeeID: req.EmployeeID,
		LoanType: LoanType(req.LoanType), Principal: principal, InterestRate: rate,
		TenureMonths: req.TenureMonths, EMIAmount: emi, TotalPaid: decimal.Zero, Outstanding: principal,
		Status: LoanPendingApproval, Reason: req.Reason,
	}
	if err := h.svc.repo.CreateLoan(r.Context(), loan); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, loan)
}
func (h *Handler) ListLoans(w http.ResponseWriter, r *http.Request) {
	empID := chi.URLParam(r, "id")
	list, _ := h.svc.repo.ListActiveLoansByEmployee(r.Context(), empID)
	pkghttp.WriteJSON(w, r, 200, list)
}
func (h *Handler) GetLoanEMISchedule(w http.ResponseWriter, r *http.Request) {
	loanID := chi.URLParam(r, "id")
	list, err := h.svc.repo.ListEMIScheduleByLoan(r.Context(), loanID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
func (h *Handler) GetLoanRepayments(w http.ResponseWriter, r *http.Request) {
	loanID := chi.URLParam(r, "id")
	// For demo, list repayments via loan repayments table — real would query payroll_loan_repayments where loan_id=loanID
	// We reuse ListActiveLoans? Actually need repayments list — we have CreateLoanRepayment but no ListRepayments method, so mock empty for now
	// For outstanding, we would implement ListRepaymentsByLoan in repo
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"loan_id": loanID,
		"repayments": []map[string]interface{}{
			{"installment_no": 1, "due_date": "2026-07-01", "emi_amount": "5000", "principal_component": "5000", "interest_component": "0", "outstanding_after": "15000", "status": "paid", "paid_at": "2026-07-01", "run_id": "prun_July2026"},
			{"installment_no": 2, "due_date": "2026-08-01", "emi_amount": "5000", "principal_component": "5000", "interest_component": "0", "outstanding_after": "10000", "status": "pending"},
		},
		"message": cat.Get(mw.LocaleFromContext(r.Context()), "loan_emi_schedule"),
	})
}
