package payroll

import (
	"encoding/json"
	"net/http"

	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/employees", h.CreateEmployee)
	r.Get("/employees", h.ListEmployees)
	r.Post("/payroll_runs", h.CreateRun)
	r.Post("/payroll_runs/{id}/calculate", h.Calculate)
	r.Post("/payroll_runs/{id}/approve", h.Approve)
	r.Post("/payroll_runs/{id}/disburse", h.Disburse)
	r.Get("/payroll_runs/{id}/items", h.ListItems)
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	merchantID, _ := r.Context().Value("merchant_id").(string)
	var req struct {
		EmployeeCode, Name, BaseSalary, BankCode, BankAccount string `json:"employee_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	base, _ := decimal.NewFromString(req.BaseSalary)
	emp := &Employee{ID: id.NewEmployee(), MerchantID: merchantID, EmployeeCode: req.EmployeeCode, Name: req.Name, BaseSalary: base, BankCode: req.BankCode, BankAccountMasked: req.BankAccount, Status: "active"}
	// Set employment date now
	// ...
	if err := h.svc.repo.CreateEmployee(r.Context(), emp); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, emp)
}

func (h *Handler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	merchantID, _ := r.Context().Value("merchant_id").(string)
	list, err := h.svc.repo.ListEmployees(r.Context(), merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) CreateRun(w http.ResponseWriter, r *http.Request) {
	merchantID, _ := r.Context().Value("merchant_id").(string)
	var req struct {
		RunRef                  string `json:"run_ref"`
		PeriodMonth, PeriodYear int    `json:"period_month"`
		Type                    string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	run := &PayrollRun{ID: id.NewPayrollRun(), MerchantID: merchantID, RunRef: req.RunRef, PeriodMonth: req.PeriodMonth, PeriodYear: req.PeriodYear, Type: RunType(req.Type), Status: StatusDraft}
	if err := h.svc.repo.CreateRun(r.Context(), run); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, run)
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	merchantID, _ := r.Context().Value("merchant_id").(string)
	runID := chi.URLParam(r, "id")
	if err := h.svc.CalculateRun(r.Context(), merchantID, runID); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"run_id": runID, "status": string(StatusPendingApproval), "message": "calculated totals gross/net/tax/pension binary search O(log n)"})
}

func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	merchantID, _ := r.Context().Value("merchant_id").(string)
	runID := chi.URLParam(r, "id")
	userID, _ := r.Context().Value("user_id").(string)
	if err := h.svc.ApproveRun(r.Context(), merchantID, runID, userID); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"run_id": runID, "status": string(StatusApproved)})
}

func (h *Handler) Disburse(w http.ResponseWriter, r *http.Request) {
	merchantID, _ := r.Context().Value("merchant_id").(string)
	runID := chi.URLParam(r, "id")
	if err := h.svc.DisburseRun(r.Context(), merchantID, runID); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"run_id": runID, "status": string(StatusProcessing), "message": "ledger M4 Dr salary Cr payroll_payable Cr tax Cr pension + payout batch created"})
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	runID := chi.URLParam(r, "id")
	list, err := h.svc.repo.ListItems(r.Context(), runID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
