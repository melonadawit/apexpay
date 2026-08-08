// Payroll claims handlers (manager -> finance approval).
package payroll

import (
	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	"net/http"
)

func (h *Handler) CreateClaim(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		EmployeeID     string `json:"employee_id"`
		ClaimType      string `json:"claim_type"` // expense/medical/travel/other
		Amount         string `json:"amount"`
		Description    string `json:"description"`
		ReceiptFileKey string `json:"receipt_file_key"` // MinIO presigned 15m TTL <5MB pdf/jpg/png
		IsTaxable      bool   `json:"is_taxable"`
		IsPensionable  bool   `json:"is_pensionable"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	amount, _ := decimal.NewFromString(req.Amount)
	claim := &ClaimEnhanced{
		ID:             id.New("claim"),
		MerchantID:     merchantID,
		EmployeeID:     req.EmployeeID,
		ClaimType:      ClaimType(req.ClaimType),
		Amount:         amount,
		Description:    req.Description,
		ReceiptFileKey: strPtr(req.ReceiptFileKey),
		Status:         "pending",
		IsTaxable:      req.IsTaxable,
		IsPensionable:  req.IsPensionable,
	}
	if err := h.svc.repo.CreateClaimEnhanced(r.Context(), claim); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, claim)
}
func (h *Handler) ListClaims(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	employeeID := r.URL.Query().Get("employee_id")
	status := r.URL.Query().Get("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}
	list, err := h.svc.repo.ListClaimsByEmployee(r.Context(), merchantID, employeeID, statusPtr)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
func (h *Handler) ApproveClaimManager(w http.ResponseWriter, r *http.Request) {
	claimID := chi.URLParam(r, "id")
	userID := mw.UserID(r.Context())
	if err := h.svc.repo.ApproveClaimManager(r.Context(), claimID, userID); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"id": claimID, "status": "approved_by_manager", "approved_by_manager": userID, "message": "Manager approved • Next finance approval • Status approved_by_manager • Receipt MinIO presigned 15m • Hash integrity • Encrypted SSE-S3 • 7y retention NBE"})
}
func (h *Handler) ApproveClaimFinance(w http.ResponseWriter, r *http.Request) {
	claimID := chi.URLParam(r, "id")
	userID := mw.UserID(r.Context())
	if err := h.svc.repo.ApproveClaimFinance(r.Context(), claimID, userID); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"id": claimID, "status": "approved", "approved_by_finance": userID, "message": "Finance approved • Status approved • Paid via next payroll run • Reimbursement non-taxable added after tax • Payroll item other_allowances reimbursement non-taxable • Outstanding"})
}
