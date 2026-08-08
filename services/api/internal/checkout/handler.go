package checkout

import (
	"encoding/json"
	"net/http"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/link"
	"apexpay/internal/payment"
	pkghttp "apexpay/internal/platform/http"
	"github.com/go-chi/chi/v5"
)

// Handler is the public (no-API-key) hosted-checkout API. Auth is the payment link's
// public token, which acts as the capability to pay against a specific merchant link.
type Handler struct {
	links    *link.PgRepository
	payments *payment.Service
}

func NewHandler(links *link.PgRepository, payments *payment.Service) *Handler {
	return &Handler{links: links, payments: payments}
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/{token}", h.Get)
	r.Post("/{token}/initialize", h.Initialize)
	r.Get("/{token}/status/{txRef}", h.Status)
	r.Post("/{token}/2fa/{paymentID}", h.Verify2FA)
}

// Get returns the payment link details for the checkout screen.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	pl, err := h.links.GetByToken(r.Context(), token)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "payment link not found")
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"token":       pl.PublicToken,
		"amount":      pl.Amount.String(),
		"currency":    pl.Currency,
		"description": pl.Description,
		"status":      pl.Status,
	})
}

// Initialize creates a payment against the link's merchant and opens a checkout session.
// The public token authorizes this; no merchant API key is exposed to the browser.
func (h *Handler) Initialize(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	pl, err := h.links.GetByToken(r.Context(), token)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "payment link not found")
		return
	}
	if pl.Status != "active" && pl.Status != "open" {
		pkghttp.WriteErrorWithBody(w, r, 409, "conflict", "payment link is not active")
		return
	}

	var req struct {
		Method        string `json:"method"`
		CustomerEmail string `json:"customer_email"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Method == "" {
		req.Method = "telebirr"
	}

	txRef := "txr_" + id.New("chk")[4:24]
	p, err := h.payments.Initialize(r.Context(), payment.InitializeRequest{
		MerchantID:     pl.MerchantID,
		TxRef:          txRef,
		Amount:         pl.Amount,
		Currency:       pl.Currency,
		Method:         req.Method,
		Description:    pl.Description,
		CustomerEmail:  req.CustomerEmail,
		IdempotencyKey: "checkout_" + token + "_" + req.Method,
	})
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}

	// Record the checkout session so the page can resume/verify by token.
	cs := &link.CheckoutSession{
		ID:             id.New("cs"),
		MerchantID:     pl.MerchantID,
		PaymentID:      p.ID,
		PaymentLinkID:  &pl.ID,
		PublicToken:    token,
		Status:         "open",
		SelectedMethod: req.Method,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}
	_ = h.links.CreateCheckoutSession(r.Context(), cs)

	pkghttp.WriteJSON(w, r, 201, map[string]interface{}{
		"id": p.ID, "tx_ref": p.TxRef, "amount": p.Amount.String(), "currency": p.Currency,
		"status": p.Status, "connector_id": p.ConnectorID,
		"requires_2fa": p.Requires2FA, "fee_amount": p.FeeAmount.String(), "net_amount": p.NetAmount.String(),
	})
}

// Status returns the live payment status for the checkout page's poller.
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	txRef := chi.URLParam(r, "txRef")
	pl, err := h.links.GetByToken(r.Context(), token)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "payment link not found")
		return
	}
	p, err := h.payments.Verify(r.Context(), payment.VerifyRequest{MerchantID: pl.MerchantID, TxRef: txRef})
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"id": p.ID, "tx_ref": p.TxRef, "status": p.Status,
		"connector_id": p.ConnectorID, "requires_2fa": p.Requires2FA,
		"two_fa_verified": p.TwoFAVerified, "succeeded_at": p.SucceededAt,
	})
}

// Verify2FA submits the one-time code for payments that require 2FA (>5000 ETB).
func (h *Handler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	paymentID := chi.URLParam(r, "paymentID")
	pl, err := h.links.GetByToken(r.Context(), token)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "payment link not found")
		return
	}
	var req struct {
		OTP string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if err := h.payments.Verify2FA(r.Context(), pl.MerchantID, paymentID, req.OTP); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]bool{"two_fa_verified": true, "can_verify_now": true})
}
