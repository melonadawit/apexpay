package swarm

import (
	"encoding/json"
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/run", h.Run)
	r.Post("/{id}/confirm", h.Confirm)
	r.Get("/{id}", h.Get)
}

func (h *Handler) Run(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	userID := mw.UserID(r.Context())
	var req struct {
		Goal string `json:"goal"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	sess, err := h.svc.Run(r.Context(), merchantID, userID, req.Goal)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, sess)
}

func (h *Handler) Confirm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Confirmed bool `json:"confirmed"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	sess, err := h.svc.Confirm(r.Context(), id, req.Confirmed)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, sess)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sess, err := h.svc.repo.GetSession(r.Context(), id)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, sess)
}
