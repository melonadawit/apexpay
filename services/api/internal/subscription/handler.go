package subscription

import (
	"encoding/json"
	"net/http"

	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/customers", h.CreateCustomer)
	r.Post("/subscription_plans", h.CreatePlan)
	r.Post("/subscriptions", h.CreateSubscription)
	r.Get("/subscriptions", h.List)
	r.Post("/subscriptions/{id}/cancel", h.Cancel)
}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
		Name  string `json:"name"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	cust := &Customer{ID: id.NewCustomer(), MerchantID: merchantID, Email: req.Email, Phone: req.Phone, Name: req.Name}
	if err := h.svc.repo.CreateCustomer(r.Context(), cust); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, cust)
}

func (h *Handler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		Amount        string `json:"amount"`
		Currency      string `json:"currency"`
		IntervalType  string `json:"interval_type"`
		IntervalCount int    `json:"interval_count"`
		TrialDays     int    `json:"trial_days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	amt, _ := decimal.NewFromString(req.Amount)
	plan := &Plan{MerchantID: merchantID, Name: req.Name, Description: req.Description, Amount: amt, Currency: req.Currency, IntervalType: req.IntervalType, IntervalCount: req.IntervalCount, TrialDays: req.TrialDays}
	resp, err := h.svc.CreatePlan(r.Context(), plan)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, resp)
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		CustomerID string `json:"customer_id"`
		PlanID     string `json:"plan_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	sub, err := h.svc.CreateSubscription(r.Context(), merchantID, req.CustomerID, req.PlanID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, sub)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	list, err := h.svc.repo.ListSubscriptions(r.Context(), merchantID, nil)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_ = h.svc.repo.UpdateSubscriptionStatus(r.Context(), id, StatusCanceled)
	pkghttp.WriteJSON(w, r, 200, map[string]string{"id": id, "status": string(StatusCanceled)})
}
