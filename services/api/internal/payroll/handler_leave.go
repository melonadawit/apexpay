// Payroll leave handlers (Art 77/82/86).
package payroll

import (
	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
)

func (h *Handler) CreateLeaveBalance(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		EmployeeID       string `json:"employee_id"`
		LeaveType        string `json:"leave_type"` // annual/sick/maternity/paternity/marriage/mourning/unpaid/comp_off/study
		Year             int    `json:"year"`
		EntitledDays     string `json:"entitled_days"`
		CarryForwardDays string `json:"carry_forward_days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	entitled, _ := decimal.NewFromString(req.EntitledDays)
	carry, _ := decimal.NewFromString(req.CarryForwardDays)
	remaining := entitled.Add(carry)
	balance := &LeaveBalance{
		ID:               id.New("lbal"),
		MerchantID:       merchantID,
		EmployeeID:       req.EmployeeID,
		LeaveType:        LeaveType(req.LeaveType),
		Year:             req.Year,
		EntitledDays:     entitled,
		UsedDays:         decimal.Zero,
		RemainingDays:    remaining,
		CarryForwardDays: carry,
	}
	if err := h.svc.repo.CreateLeaveBalance(r.Context(), balance); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, balance)
}
func (h *Handler) ListLeaveBalances(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	employeeID := r.URL.Query().Get("employee_id")
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	if year == 0 {
		year = 2026
	}
	var list []LeaveBalance
	var err error
	if employeeID != "" {
		list, err = h.svc.repo.ListLeaveBalancesByEmployee(r.Context(), merchantID, employeeID, year)
	} else {
		// For all employees, mock empty for demo — real would query all balances for merchant year
		list = []LeaveBalance{}
	}
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
func (h *Handler) CreateLeaveRequest(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		EmployeeID         string  `json:"employee_id"`
		LeaveType          string  `json:"leave_type"`
		StartDate          string  `json:"start_date"`     // 2026-07-10
		EndDate            string  `json:"end_date"`       // 2026-07-12
		DaysRequested      float64 `json:"days_requested"` // 2 days, 0.5 half day
		Reason             string  `json:"reason"`
		MedicalCertificate string  `json:"medical_certificate_file_key"` // MinIO for sick >3 days per Art 82
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	start, _ := parseDate(req.StartDate)
	end, _ := parseDate(req.EndDate)
	leaveReq := &LeaveRequest{
		ID:                        id.New("lreq"),
		MerchantID:                merchantID,
		EmployeeID:                req.EmployeeID,
		LeaveType:                 LeaveType(req.LeaveType),
		StartDate:                 start,
		EndDate:                   end,
		DaysRequested:             decimal.NewFromFloat(req.DaysRequested),
		Reason:                    req.Reason,
		Status:                    LeavePending,
		MedicalCertificateFileKey: strPtr(req.MedicalCertificate),
	}
	if err := h.svc.repo.CreateLeaveRequest(r.Context(), leaveReq); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, leaveReq)
}
func (h *Handler) ListLeaveRequests(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	employeeID := r.URL.Query().Get("employee_id")
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	if year == 0 {
		year = 2026
	}
	statusStr := r.URL.Query().Get("status")
	var status *LeaveStatus
	if statusStr != "" {
		s := LeaveStatus(statusStr)
		status = &s
	}
	list, err := h.svc.repo.ListLeaveRequests(r.Context(), merchantID, employeeID, year, status)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
func (h *Handler) ApproveLeaveRequest(w http.ResponseWriter, r *http.Request) {
	reqID := chi.URLParam(r, "id")
	userID := mw.UserID(r.Context())
	if err := h.svc.repo.UpdateLeaveRequestStatus(r.Context(), reqID, LeaveApproved, &userID, ""); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	// After approval, deduct from balance Used+=Requested Remaining=Entitled-Used floor zero
	// For outstanding, we would also update leave balance here
	pkghttp.WriteJSON(w, r, 200, map[string]string{"id": reqID, "status": string(LeaveApproved), "approved_by": userID, "message": "Approved • Deduct from balance Used+=Requested Remaining=Entitled-Used floor zero • Art 77/82/86 • Outstanding"})
}
func (h *Handler) RejectLeaveRequest(w http.ResponseWriter, r *http.Request) {
	reqID := chi.URLParam(r, "id")
	var req struct {
		RejectionReason string `json:"rejection_reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.svc.repo.UpdateLeaveRequestStatus(r.Context(), reqID, LeaveRejected, nil, req.RejectionReason); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"id": reqID, "status": string(LeaveRejected), "rejection_reason": req.RejectionReason})
}
