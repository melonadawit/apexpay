package routing

import (
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Get("/", h.Ranked)
	r.Get("/methods", h.Ranked)
	r.Get("/rules", h.ListRules)
}

func (h *Handler) Ranked(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	amountStr := r.URL.Query().Get("amount")
	currency := r.URL.Query().Get("currency")
	if currency == "" {
		currency = "ETB"
	}
	amt, _ := decimal.NewFromString(amountStr)
	if amt.IsZero() {
		amt = decimal.NewFromInt(1000)
	}

	ranked, err := h.svc.RankedMethods(r.Context(), merchantID, amt, currency)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, ranked)
}

func (h *Handler) ListRules(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	mID := &merchantID
	rules, err := h.svc.repo.ListRules(r.Context(), mID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, rules)
}
