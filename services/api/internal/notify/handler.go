package notify

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
	r.Get("/", h.List)
	r.Put("/", h.Upsert)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	var p Preference
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.EventType == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "event_type required")
		return
	}
	if err := h.repo.Upsert(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), &p); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, p)
}
