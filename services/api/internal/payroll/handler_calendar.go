// Payroll calendar handlers (Ethiopia cutoff 25th / disbursal 30th).
package payroll

import (
	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) CreateCalendar(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		PayFrequency  string `json:"pay_frequency"` // monthly/semimonthly/weekly/biweekly
		Year          int    `json:"year"`
		Month         *int   `json:"month"`
		CutoffDay     int    `json:"cutoff_day"`     // Ethiopia business practice cutoff 25th
		DisbursalDay  int    `json:"disbursal_day"`  // disbursal 30th
		PayDay        int    `json:"pay_day"`        // pay date last day
		CutoffDate    string `json:"cutoff_date"`    // 2026-07-25
		DisbursalDate string `json:"disbursal_date"` // 2026-07-30
		PayDate       string `json:"pay_date"`       // 2026-07-31
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	cutoffDate, _ := parseDate(req.CutoffDate)
	disbursalDate, _ := parseDate(req.DisbursalDate)
	payDate, _ := parseDate(req.PayDate)
	if cutoffDate.IsZero() {
		// Auto-calc from cutoff_day year month
		if req.Month != nil {
			cutoffDate = time.Date(req.Year, time.Month(*req.Month), req.CutoffDay, 0, 0, 0, 0, time.UTC)
		}
	}
	if disbursalDate.IsZero() && req.Month != nil {
		disbursalDate = time.Date(req.Year, time.Month(*req.Month), req.DisbursalDay, 0, 0, 0, 0, time.UTC)
	}
	if payDate.IsZero() && req.Month != nil {
		// Last day of month
		firstOfNextMonth := time.Date(req.Year, time.Month(*req.Month)+1, 1, 0, 0, 0, 0, time.UTC)
		payDate = firstOfNextMonth.AddDate(0, 0, -1)
		if req.PayDay != 0 && req.PayDay <= 31 {
			// If pay_day specified, use it, but if last day requested and month has 31, use 31 else last day
			if req.PayDay == 31 {
				// keep last day
			} else {
				payDate = time.Date(req.Year, time.Month(*req.Month), req.PayDay, 0, 0, 0, 0, time.UTC)
			}
		}
	}
	cal := &PayrollCalendar{
		ID:            id.New("pcal"),
		MerchantID:    merchantID,
		Name:          req.Name,
		Description:   req.Description,
		PayFrequency:  PayFrequency(req.PayFrequency),
		Year:          req.Year,
		Month:         req.Month,
		CutoffDay:     req.CutoffDay,
		DisbursalDay:  req.DisbursalDay,
		PayDay:        req.PayDay,
		CutoffDate:    cutoffDate,
		DisbursalDate: disbursalDate,
		PayDate:       payDate,
		IsLocked:      false,
	}
	if cal.PayFrequency == "" {
		cal.PayFrequency = PayFrequencyMonthly
	}
	if cal.CutoffDay == 0 {
		cal.CutoffDay = 25 // Ethiopia business practice cutoff 25th
	}
	if cal.DisbursalDay == 0 {
		cal.DisbursalDay = 30 // disbursal 30th
	}
	if cal.PayDay == 0 {
		cal.PayDay = 31 // pay date last day
	}
	if err := h.svc.repo.CreateCalendar(r.Context(), cal); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, cal)
}
func (h *Handler) ListCalendars(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	if year == 0 {
		year = 2026
	}
	list, err := h.svc.repo.ListCalendars(r.Context(), merchantID, year)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
func (h *Handler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	calID := chi.URLParam(r, "id")
	cal, err := h.svc.repo.GetCalendar(r.Context(), merchantID, calID)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "calendar not found")
		return
	}
	pkghttp.WriteJSON(w, r, 200, cal)
}
func (h *Handler) LockCalendar(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	calID := chi.URLParam(r, "id")
	userID := mw.UserID(r.Context())
	if err := h.svc.repo.LockCalendar(r.Context(), merchantID, calID, userID); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"id": calID, "is_locked": "true", "locked_by": userID, "message": "Locked after disbursal per Ethiopia business practice • Prevents re-run amendment unless unlocked by admin with audit log payroll_audit_logs actor admin action unlock_calendar details locked_by IP inet request_id immutable"})
}
func (h *Handler) UnlockCalendar(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	calID := chi.URLParam(r, "id")
	if err := h.svc.repo.UnlockCalendar(r.Context(), merchantID, calID); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"id": calID, "is_locked": "false", "message": "Unlocked by admin with audit log payroll_audit_logs actor admin action unlock_calendar"})
}
