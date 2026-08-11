// Payroll final-settlement (F&F) handlers.
package payroll

import (
	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"encoding/json"
	"github.com/shopspring/decimal"
	"net/http"
)

func (h *Handler) CreateFinalSettlement(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		EmployeeID          string  `json:"employee_id"`
		ResignationDate     string  `json:"resignation_date"`
		LastWorkingDate     string  `json:"last_working_date"`
		NoticePeriodDays    int     `json:"notice_period_days"`
		NoticeServedDays    int     `json:"notice_served_days"`
		LeaveEncashmentDays float64 `json:"leave_encashment_days"`
		SeveranceAmount     string  `json:"severance_amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	resig, _ := parseDate(req.ResignationDate)
	lwd, _ := parseDate(req.LastWorkingDate)
	sev, _ := decimal.NewFromString(req.SeveranceAmount)
	fs := &FinalSettlement{
		ID: id.New("fnf"), MerchantID: merchantID, EmployeeID: req.EmployeeID,
		ResignationDate: resig, LastWorkingDate: lwd,
		NoticePeriodDays: req.NoticePeriodDays, NoticeServedDays: req.NoticeServedDays,
		NoticeShortfallDays: req.NoticePeriodDays - req.NoticeServedDays,
		LeaveEncashmentDays: decimal.NewFromFloat(req.LeaveEncashmentDays),
		SeveranceAmount:     sev, Status: "draft",
	}
	if err := h.svc.CreateFinalSettlement(r.Context(), fs); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, fs)
}
func (h *Handler) ListFinalSettlements(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	list, err := h.svc.repo.ListFinalSettlements(r.Context(), merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
