package refund

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	pkghttp "apexpay/internal/platform/http"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Get("/payment/{paymentId}", h.ListByPayment)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	merchantID, _ := r.Context().Value("merchant_id").(string)
	var req struct {
		PaymentID  string `json:"payment_id"`
		RefundRef  string `json:"refund_ref"`
		Amount     string `json:"amount"`
		Currency   string `json:"currency"`
		Reason     string `json:"reason"`
		FeePolicy  string `json:"fee_policy"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	amt, err := decimal.NewFromString(req.Amount)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "amount numeric")
		return
	}
	idemKey := r.Header.Get("Idempotency-Key")
	rf, err := h.svc.Create(r.Context(), CreateRequest{
		MerchantID: merchantID, PaymentID: req.PaymentID, RefundRef: req.RefundRef,
		Amount: amt, Currency: req.Currency, Reason: req.Reason, FeePolicy: FeePolicy(req.FeePolicy),
		IdempotencyKey: idemKey,
	})
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, map[string]interface{}{
		"id": rf.ID, "refund_ref": rf.RefundRef, "amount": rf.Amount.String(), "fee_reversal": rf.FeeReversal.String(),
		"status": rf.Status, "payment_id": rf.PaymentID, "ledger_model": "M2 Dr payable R-FR + Dr fee_due FR Cr clearing R",
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// Simplified: would call repo GetByID
	pkghttp.WriteJSON(w, r, 200, map[string]string{"id": id, "status": "succeeded"})
}

func (h *Handler) ListByPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "paymentId")
	list, err := h.svc.repo.ListRefundsByPayment(r.Context(), paymentID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
