package fayda

import (
	"encoding/json"
	"net/http"

	"apexpay/internal/i18n"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

var cat = i18n.New()

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/verify/init", h.Init)
	r.Post("/verify/confirm", h.ConfirmOTP)
	r.Post("/verify/qr", h.VerifyQR)
	r.Get("/owner/{ownerId}", h.ListByOwner)
}

func (h *Handler) Init(w http.ResponseWriter, r *http.Request) {
	var req InitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	// merchant_id from context or body
	if req.MerchantID == "" {
		if merchantID := mw.MerchantID(r.Context()); merchantID != "" {
			req.MerchantID = merchantID
		}
	}
	req.ConsentIP = r.RemoteAddr
	resp, err := h.svc.Init(r.Context(), req)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	// Return only last4 per privacy, not plain FIN
	pkghttp.WriteJSON(w, r, 201, map[string]interface{}{
		"request_id": resp.RequestID, "fin_last4": resp.FinLast4, "status": resp.Status, "otp_sent": true, "fayda_transaction_id": resp.FaydaTransactionID,
		"message": cat.Get(mw.LocaleFromContext(r.Context()), "fayda_otp_sent"),
	})
}

func (h *Handler) ConfirmOTP(w http.ResponseWriter, r *http.Request) {
	var req ConfirmOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	resp, err := h.svc.ConfirmOTP(r.Context(), req)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"request_id": resp.RequestID, "status": resp.Status, "otp_verified": resp.OTPVerified,
		"demographics_match": resp.DemographicsMatch, "face_match": resp.FaceMatch, "face_score": resp.FaceMatchScore,
		"fin_last4": resp.FinLast4,
	})
}

func (h *Handler) VerifyQR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RequestID string `json:"request_id"`
		QRData    string `json:"qr_data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	valid, err := h.svc.VerifyOfflineQR(r.Context(), req.QRData, req.RequestID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]bool{"offline_verified": valid})
}

func (h *Handler) ListByOwner(w http.ResponseWriter, r *http.Request) {
	ownerID := chi.URLParam(r, "ownerId")
	list, err := h.svc.repo.GetByOwner(r.Context(), ownerID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
