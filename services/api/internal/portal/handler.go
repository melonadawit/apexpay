package portal

import (
	"encoding/json"
	"net/http"
	"time"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ repo *Repository }

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) Routes(r chi.Router) {
	// Admin-facing: issue a portal token for a vendor or customer.
	r.Post("/token", h.Create)
	// Token-gated self-service: callers pass ?token= in the header.
	r.Get("/me", h.Me)
}

// Create issues a portal access token (admin / merchant session required).
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PortalType string `json:"portal_type"` // vendor | customer
		EntityID   string `json:"entity_id"`   // vendor id or customer email
		EntityName string `json:"entity_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if (req.PortalType != "vendor" && req.PortalType != "customer") || req.EntityID == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "portal_type (vendor|customer) and entity_id required")
		return
	}
	acc, err := h.repo.Create(r.Context(), middleware.MerchantID(r.Context()), PortalType(req.PortalType), req.EntityID, req.EntityName, 72*time.Hour)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, acc)
}

// Me returns the caller's self-service dashboard, authenticated by a portal token in the
// X-Portal-Token header.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Portal-Token")
	if token == "" {
		pkghttp.WriteErrorWithBody(w, r, 401, "unauthorized", "portal token required (X-Portal-Token)")
		return
	}
	acc, err := h.repo.Resolve(r.Context(), token)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 401, "unauthorized", "invalid or expired portal token")
		return
	}
	dash := &Dashboard{PortalType: acc.PortalType, EntityName: acc.EntityName, Invoices: []any{}}
	switch acc.PortalType {
	case PortalVendor:
		invs, err := h.repo.VendorInvoices(r.Context(), acc.MerchantID, acc.EntityID)
		if err != nil {
			pkghttp.WriteError(w, r, err)
			return
		}
		for i := range invs {
			dash.Invoices = append(dash.Invoices, invs[i])
		}
	case PortalCustomer:
		invs, err := h.repo.CustomerInvoices(r.Context(), acc.MerchantID, acc.EntityID)
		if err != nil {
			pkghttp.WriteError(w, r, err)
			return
		}
		for i := range invs {
			dash.Invoices = append(dash.Invoices, invs[i])
		}
	}
	pkghttp.WriteJSON(w, r, 200, dash)
}
