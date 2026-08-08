package dispute

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
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Post("/{id}/evidence", h.Evidence)
	r.Post("/{id}/decide", h.Decide)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var d Dispute
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if d.PaymentID == "" || d.ReasonCode == "" || d.Amount == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "payment_id, reason_code, amount required")
		return
	}
	if err := h.repo.Create(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), &d); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, d)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.List(r.Context(), middleware.MerchantID(r.Context()), r.URL.Query().Get("status"), 50)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}

func (h *Handler) Evidence(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Evidence []EvidenceItem `json:"evidence"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if err := h.repo.SubmitEvidence(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"), req.Evidence); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"status": "evidence_submitted"})
}

func (h *Handler) Decide(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Decision   string `json:"decision"`
		Resolution string `json:"resolution"`
		Fee        string `json:"fee"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if err := h.repo.Decide(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"), req.Decision, req.Resolution, req.Fee); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", err.Error())
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"status": req.Decision})
}
