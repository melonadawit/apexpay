package risk

import (
	"encoding/json"
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Get("/rules", h.ListRules)
	r.Post("/rules", h.CreateRule)
	r.Get("/flags", h.ListFlags)
	r.Post("/evaluate", h.Evaluate)
}

func (h *Handler) ListRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.svc.repo.ListRules(r.Context(), middleware.MerchantID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, rules)
}

func (h *Handler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var rule Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if rule.Name == "" || rule.RuleType == "" || rule.Action == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "name, rule_type, action required")
		return
	}
	if rule.Severity == "" {
		rule.Severity = "medium"
	}
	rule.Enabled = true
	if err := h.svc.repo.CreateRule(r.Context(), middleware.MerchantID(r.Context()), &rule); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, rule)
}

func (h *Handler) ListFlags(w http.ResponseWriter, r *http.Request) {
	flags, err := h.svc.repo.ListFlags(r.Context(), middleware.MerchantID(r.Context()), r.URL.Query().Get("status"), 100)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, flags)
}

// Evaluate runs the risk engine against a candidate transaction and returns matched findings.
func (h *Handler) Evaluate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EntityType string `json:"entity_type"`
		EntityID   string `json:"entity_id"`
		Amount     string `json:"amount"`
		DeviceID   string `json:"device_id"`
		IP         string `json:"ip"`
		CustomerID string `json:"customer_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "amount must be numeric")
		return
	}
	etype := req.EntityType
	if etype == "" {
		etype = "payment"
	}
	eval, err := h.svc.Evaluate(r.Context(), EvaluationContext{
		MerchantID: middleware.MerchantID(r.Context()),
		EntityType: etype, EntityID: req.EntityID, Amount: amount,
		DeviceID: req.DeviceID, IP: req.IP, CustomerID: req.CustomerID,
	})
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"findings": eval.Matched(),
		"block":    eval.HasBlock(),
		"approved": !eval.HasBlock(),
	})
}
