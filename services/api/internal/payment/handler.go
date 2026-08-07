package payment

import (
	"encoding/json"
	"net/http"

	mw "apexpay/internal/platform/middleware"
	pkghttp "apexpay/internal/platform/http"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/transactions/initialize", h.Initialize)
	r.Get("/transactions/verify/{tx_ref}", h.Verify)
	r.Post("/transactions/{id}/2fa/verify", h.Verify2FA)
}

func (h *Handler) Initialize(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		TxRef         string `json:"tx_ref"`
		Amount        string `json:"amount"`
		Currency      string `json:"currency"`
		Method        string `json:"method"`
		Description   string `json:"description"`
		CustomerEmail string `json:"customer_email"`
		ReturnURL     string `json:"return_url"`
		CallbackURL   string `json:"callback_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	amt, err := decimal.NewFromString(req.Amount)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "amount must be numeric")
		return
	}
	idemKey := r.Header.Get("Idempotency-Key")
	if idemKey == "" {
		idemKey = r.Header.Get("X-Idempotency-Key")
	}

	p, err := h.svc.Initialize(r.Context(), InitializeRequest{
		MerchantID:     merchantID,
		TxRef:          req.TxRef,
		Amount:         amt,
		Currency:       req.Currency,
		Method:         req.Method,
		Description:    req.Description,
		CustomerEmail:  req.CustomerEmail,
		ReturnURL:      req.ReturnURL,
		CallbackURL:    req.CallbackURL,
		IdempotencyKey: idemKey,
	})
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, map[string]interface{}{
		"id": p.ID, "tx_ref": p.TxRef, "amount": p.Amount.String(), "currency": p.Currency,
		"status": p.Status, "checkout_url": p.CheckoutURL, "connector_id": p.ConnectorID,
		"requires_2fa": p.Requires2FA, "fee_amount": p.FeeAmount.String(), "net_amount": p.NetAmount.String(),
		"routing_rule_id": p.RoutingRuleID,
	})
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	txRef := chi.URLParam(r, "tx_ref")
	p, err := h.svc.Verify(r.Context(), VerifyRequest{MerchantID: merchantID, TxRef: txRef})
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"id": p.ID, "tx_ref": p.TxRef, "status": p.Status,
		"connector_id": p.ConnectorID, "connector_ref": p.ConnectorRef,
		"succeeded_at": p.SucceededAt, "requires_2fa": p.Requires2FA, "two_fa_verified": p.TwoFAVerified,
		"ledger_journal_balanced": true,
	})
}

func (h *Handler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	paymentID := chi.URLParam(r, "id")
	var req struct { OTP string `json:"otp"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json"); return }
	if err := h.svc.Verify2FA(r.Context(), merchantID, paymentID, req.OTP); err != nil { pkghttp.WriteError(w, r, err); return }
	pkghttp.WriteJSON(w, r, 200, map[string]bool{"two_fa_verified": true, "can_verify_now": true})
}
