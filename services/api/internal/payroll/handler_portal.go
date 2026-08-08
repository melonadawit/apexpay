// Payroll employee-portal handlers (magic link).
package payroll

import (
	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) CreateMagicLink(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeID string `json:"employee_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	// Magic link JWT 24h HMAC signed employee_id+merchant_id+expiry
	// Placeholder: return URL
	magicURL := fmt.Sprintf("https://employee.apexpay.et/portal?token=%s_%s_magic_24h", req.EmployeeID, id.New("tok"))
	pkghttp.WriteJSON(w, r, 201, map[string]string{"magic_link": magicURL, "expires_in": "24h", "message": "employee portal magic link JWT 24h + WhatsApp integration"})
}
func (h *Handler) GetMyPortal(w http.ResponseWriter, r *http.Request) {
	pkghttp.WriteJSON(w, r, 200, map[string]string{"message": "employee self-service portal: payslips YTD, claims, loans, docs"})
}
