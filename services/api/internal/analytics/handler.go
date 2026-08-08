package analytics

import (
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
	r.Get("/revenue", h.Revenue)
	r.Get("/methods", h.Methods)
	r.Get("/cohorts", h.Cohorts)
}

func (h *Handler) Revenue(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.Revenue(r.Context(), middleware.MerchantID(r.Context()), 30)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) Methods(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.MethodBreakdown(r.Context(), middleware.MerchantID(r.Context()), 30)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) Cohorts(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.Cohorts(r.Context(), middleware.MerchantID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
