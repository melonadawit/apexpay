package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) Routes(r chi.Router) {
	r.Get("/onboarding/queue", h.OnboardingQueue)
	r.Post("/onboarding/{id}/review", h.Review)
	r.Get("/connectors/health", h.ConnectorHealth)
	r.Get("/recon/breaks", h.ReconBreaks)
	r.Get("/evidence", h.Evidence)
	r.Get("/merchants/{id}/exam", h.MerchantExam)
}

func (h *Handler) OnboardingQueue(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}
	items, err := h.repo.ListOnboardingQueue(r.Context(), limit)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, items)
}

func (h *Handler) Review(w http.ResponseWriter, r *http.Request) {
	merchantID := chi.URLParam(r, "id")
	var req ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	reviewerID := middleware.UserID(r.Context())
	if reviewerID == "" {
		if kID, ok := middleware.APIKeyIDFromContext(r.Context()); ok {
			reviewerID = kID
		}
	}
	res, err := h.repo.Review(r.Context(), merchantID, middleware.Role(r.Context()), reviewerID, req.Action, req.Comment)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "merchant not found")
			return
		}
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, res)
}

func (h *Handler) ConnectorHealth(w http.ResponseWriter, r *http.Request) {
	health, err := h.repo.ConnectorHealth(r.Context())
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, health)
}

func (h *Handler) ReconBreaks(w http.ResponseWriter, r *http.Request) {
	breaks, err := h.repo.ReconBreaks(r.Context(), 100)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, breaks)
}

func (h *Handler) Evidence(w http.ResponseWriter, r *http.Request) {
	txRef := r.URL.Query().Get("tx_ref")
	if txRef == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "tx_ref required")
		return
	}
	ev, err := h.repo.Evidence(r.Context(), txRef)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, ev)
}

func (h *Handler) MerchantExam(w http.ResponseWriter, r *http.Request) {
	exam, err := h.repo.MerchantExam(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "merchant not found")
			return
		}
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, exam)
}
