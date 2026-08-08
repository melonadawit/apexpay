package fxreval

import (
	"encoding/json"
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/", h.Revalue)
}

func (h *Handler) Revalue(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Period string `json:"period"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	out, err := h.svc.Revalue(r.Context(), middleware.MerchantID(r.Context()), req.Period)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}
