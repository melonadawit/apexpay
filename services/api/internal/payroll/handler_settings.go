package payroll

import (
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
)

// GetTaxBracketsHandler returns the active Ethiopia payroll tax brackets for
// the payroll settings screen.
func (h *Handler) GetTaxBracketsHandler(w http.ResponseWriter, r *http.Request) {
	// Brackets are national (not merchant-scoped), but require a valid session.
	_ = mw.MerchantID(r.Context())
	brackets, err := h.svc.repo.GetTaxBrackets(r.Context())
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, brackets)
}
