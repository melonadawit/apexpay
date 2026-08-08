package compliance

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
	r.Get("/status", h.Status)
	r.Get("/checks", h.Checks)
	r.Post("/checks", h.AddCheck)
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	s, err := h.repo.GetStatus(r.Context(), middleware.MerchantID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, s)
}

func (h *Handler) Checks(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListChecks(r.Context(), middleware.MerchantID(r.Context()), 30)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) AddCheck(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CheckType string `json:"check_type"`
		Status    string `json:"status"`
		Detail    string `json:"detail"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CheckType == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "check_type required")
		return
	}
	if err := h.repo.LogCheck(r.Context(), middleware.MerchantID(r.Context()), req.CheckType, req.Status, req.Detail); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, map[string]string{"status": "logged"})
}
