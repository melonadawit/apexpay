package budget

import (
	"encoding/json"
	"net/http"
	"time"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Get("/budgets", h.List)
	r.Post("/budgets", h.Set)
	r.Get("/variance", h.Variance)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.List(r.Context(), middleware.MerchantID(r.Context()), r.URL.Query().Get("period"))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}

func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	var in BudgetInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Period == "" || in.Category == "" || in.BudgetAmount == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "period, category, budget_amount required")
		return
	}
	out, err := h.svc.SetBudget(r.Context(), middleware.MerchantID(r.Context()), in)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, out)
}

func (h *Handler) Variance(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = timeNowPeriod()
	}
	out, err := h.svc.Variance(r.Context(), middleware.MerchantID(r.Context()), period)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}

func timeNowPeriod() string {
	// YYYY-MM of the current month (UTC).
	return time.Now().UTC().Format("2006-01")
}
