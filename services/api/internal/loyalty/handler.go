package loyalty

import (
	"encoding/json"
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/tiers", h.CreateTier)
	r.Get("/tiers", h.ListTiers)
	r.Get("/accounts", h.Accounts)
	r.Post("/accounts/{id}/earn", h.Earn)
	r.Get("/transactions", h.Transactions)
}

func (h *Handler) CreateTier(w http.ResponseWriter, r *http.Request) {
	var t Tier
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if t.Name == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "name required")
		return
	}
	if err := h.repo.CreateTier(r.Context(), middleware.MerchantID(r.Context()), &t); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, t)
}

func (h *Handler) ListTiers(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListTiers(r.Context(), middleware.MerchantID(r.Context()))
	write(w, r, out, err)
}

func (h *Handler) Accounts(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.Accounts(r.Context(), middleware.MerchantID(r.Context()), 50)
	write(w, r, out, err)
}

func (h *Handler) Earn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PaymentID string `json:"payment_id"`
		Amount    string `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "amount required")
		return
	}
	tx, err := h.repo.EarnCashback(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"), req.PaymentID, req.Amount)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, tx)
}

func (h *Handler) Transactions(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.Transactions(r.Context(), middleware.MerchantID(r.Context()), 50)
	write(w, r, out, err)
}

func write(w http.ResponseWriter, r *http.Request, out any, err error) {
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}
