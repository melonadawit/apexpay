package treasury

import (
	"encoding/json"
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Get("/position", h.CashPosition)
	r.Post("/transfers", h.CreateTransfer)
	r.Get("/transfers", h.ListTransfers)
	r.Post("/forecast", h.GenerateForecast)
	r.Get("/forecast", h.LatestForecast)
}

func (h *Handler) CashPosition(w http.ResponseWriter, r *http.Request) {
	pos, err := h.svc.repo.CashPosition(r.Context(), middleware.MerchantID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, pos)
}

func (h *Handler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	var t Transfer
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if t.FromAccountID == "" || t.ToAccountID == "" || t.Amount == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "from_account_id, to_account_id, amount required")
		return
	}
	if t.FromAccountID == t.ToAccountID {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "cannot transfer to the same account")
		return
	}
	created, err := h.svc.CreateAndCompleteTransfer(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), &t)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "transfer_failed", err.Error())
		return
	}
	pkghttp.WriteJSON(w, r, 201, created)
}

func (h *Handler) ListTransfers(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.repo.ListTransfers(r.Context(), middleware.MerchantID(r.Context()), 50)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}

func (h *Handler) GenerateForecast(w http.ResponseWriter, r *http.Request) {
	f, err := h.svc.repo.ForecastFromObligations(r.Context(), middleware.MerchantID(r.Context()), 90)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, f)
}

func (h *Handler) LatestForecast(w http.ResponseWriter, r *http.Request) {
	f, err := h.svc.repo.LatestForecast(r.Context(), middleware.MerchantID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, f)
}
