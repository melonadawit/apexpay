package payout

import (
	"encoding/json"
	"net/http"

	"apexpay/internal/id"
	"apexpay/internal/platform/crypto"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/beneficiaries", h.CreateBeneficiary)
	r.Get("/beneficiaries", h.ListBeneficiaries)
	r.Post("/", h.CreatePayout)
	r.Post("/bulk", h.CreateBulk)
	r.Post("/batches/{id}/approve", h.ApproveBatch)
	r.Get("/batches/{id}", h.GetBatch)
	r.Get("/batches", h.ListBatches)
}

// ListBatches returns the merchant's payout batches (newest first).
func (h *Handler) ListBatches(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.repo.ListBatches(r.Context(), mw.MerchantID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

// ListBeneficiaries returns the merchant's saved beneficiaries.
func (h *Handler) ListBeneficiaries(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.repo.ListBeneficiaries(r.Context(), mw.MerchantID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) CreateBeneficiary(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct{ Name, AccountNo, BankCode, BankName, Type string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	// Name fuzzy match Levenshtein <3 check would be here vs legal_name - simplified
	hash := crypto.HashFIN("salt", req.AccountNo)
	masked := crypto.MaskAccount(req.AccountNo)
	b := &Beneficiary{ID: id.NewBeneficiary(), MerchantID: merchantID, Name: req.Name, AccountNoMasked: masked, AccountNoHash: hash, BankCode: req.BankCode, BankName: req.BankName, Type: req.Type, VerificationStatus: "pending"}
	if err := h.svc.repo.CreateBeneficiary(r.Context(), b); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, b)
}

func (h *Handler) CreatePayout(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct{ BeneficiaryID, PayoutRef, Amount, Currency, Method string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	amt, _ := decimal.NewFromString(req.Amount)
	p := &Payout{MerchantID: merchantID, BeneficiaryID: req.BeneficiaryID, PayoutRef: req.PayoutRef, Amount: amt, Currency: req.Currency, Method: req.Method}
	resp, err := h.svc.CreateSingle(r.Context(), p)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, resp)
}

func (h *Handler) CreateBulk(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req CreateBulkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	req.MerchantID = merchantID
	batch, err := h.svc.CreateBulk(r.Context(), req)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, batch)
}

func (h *Handler) ApproveBatch(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	batchID := chi.URLParam(r, "id")
	userID := mw.UserID(r.Context())
	if err := h.svc.ApproveBatch(r.Context(), merchantID, batchID, userID); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"batch_id": batchID, "status": "approved"})
}

func (h *Handler) GetBatch(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	batchID := chi.URLParam(r, "id")
	b, err := h.svc.repo.GetBatch(r.Context(), merchantID, batchID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, b)
}
