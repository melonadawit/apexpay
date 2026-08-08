// Payroll salary-structure handlers.
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

func (h *Handler) CreateSalaryStructure(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		Name        string `json:"name"`
		NameAM      string `json:"name_am"`
		Description string `json:"description"`
		CTCAnnual   string `json:"ctc_annual"`
		Currency    string `json:"currency"`
		Components  []struct {
			Code            string `json:"code"`
			Name            string `json:"name"`
			NameAM          string `json:"name_am"`
			ComponentType   string `json:"component_type"`   // earning, deduction, employer_contribution, reimbursement
			CalculationType string `json:"calculation_type"` // fixed, percentage_of_basic, percentage_of_ctc, formula
			Amount          string `json:"amount"`
			Percentage      string `json:"percentage"`
			Formula         string `json:"formula"`
			IsTaxable       bool   `json:"is_taxable"`
			IsPartOfGross   bool   `json:"is_part_of_gross"`
			IsProratable    bool   `json:"is_proratable"`
			IsPensionable   bool   `json:"is_pensionable"`
			OrderNo         int    `json:"order_no"`
		} `json:"components"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	ctc, _ := decimal.NewFromString(req.CTCAnnual)
	structure := &SalaryStructure{
		ID: id.New("sstr"), MerchantID: merchantID,
		Name: req.Name, NameAM: req.NameAM, Description: req.Description,
		CTCAnnual: ctc, Currency: req.Currency, Status: "active",
	}
	for _, c := range req.Components {
		amt, _ := decimal.NewFromString(c.Amount)
		perc, _ := decimal.NewFromString(c.Percentage)
		comp := StructureComponent{
			ID: id.New("scomp"), StructureID: structure.ID,
			Code: c.Code, Name: c.Name, NameAM: c.NameAM,
			ComponentType: ComponentType(c.ComponentType), CalculationType: CalculationType(c.CalculationType),
			Amount: amt, Percentage: perc, Formula: c.Formula,
			IsTaxable: c.IsTaxable, IsPartOfGross: c.IsPartOfGross, IsProratable: c.IsProratable, IsPensionable: c.IsPensionable,
			OrderNo: c.OrderNo,
		}
		structure.Components = append(structure.Components, comp)
	}
	if err := h.svc.CreateSalaryStructure(r.Context(), merchantID, structure); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, structure)
}
func (h *Handler) ListSalaryStructures(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	list, err := h.svc.repo.ListSalaryStructures(r.Context(), merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
func (h *Handler) GetSalaryStructure(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	idParam := chi.URLParam(r, "id")
	s, err := h.svc.repo.GetSalaryStructure(r.Context(), merchantID, idParam)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "structure not found")
		return
	}
	pkghttp.WriteJSON(w, r, 200, s)
}
