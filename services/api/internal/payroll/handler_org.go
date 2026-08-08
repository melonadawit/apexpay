// Payroll organization hierarchy handlers (departments, designations, grades, branches).
package payroll

import (
	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"encoding/json"
	"github.com/shopspring/decimal"
	"net/http"
)

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		Name        string `json:"name"`
		NameAM      string `json:"name_am"`
		Code        string `json:"code"`
		CostCenter  string `json:"cost_center"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	dept := &Department{
		ID: id.New("dept"), MerchantID: merchantID,
		Name: req.Name, NameAM: req.NameAM, Code: req.Code, CostCenter: req.CostCenter, Description: req.Description,
	}
	if err := h.svc.repo.CreateDepartment(r.Context(), dept); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, dept)
}
func (h *Handler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	list, err := h.svc.repo.ListDepartments(r.Context(), merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
func (h *Handler) CreateDesignation(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		Title       string `json:"title"`
		TitleAM     string `json:"title_am"`
		Level       int    `json:"level"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	desg := &Designation{ID: id.New("desg"), MerchantID: merchantID, Title: req.Title, TitleAM: req.TitleAM, Level: req.Level, Description: req.Description}
	if err := h.svc.repo.CreateDesignation(r.Context(), desg); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, desg)
}
func (h *Handler) ListDesignations(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	list, _ := h.svc.repo.ListDesignations(r.Context(), merchantID)
	pkghttp.WriteJSON(w, r, 200, list)
}
func (h *Handler) CreateGrade(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		Name        string `json:"name"`
		NameAM      string `json:"name_am"`
		MinSalary   string `json:"min_salary"`
		MaxSalary   string `json:"max_salary"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	min, _ := decimal.NewFromString(req.MinSalary)
	max, _ := decimal.NewFromString(req.MaxSalary)
	grade := &Grade{ID: id.New("grade"), MerchantID: merchantID, Name: req.Name, NameAM: req.NameAM, MinSalary: min, MaxSalary: max, Description: req.Description}
	if err := h.svc.repo.CreateGrade(r.Context(), grade); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, grade)
}
func (h *Handler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		Name    string `json:"name"`
		NameAM  string `json:"name_am"`
		Region  string `json:"region"`
		City    string `json:"city"`
		SubCity string `json:"sub_city"`
		Address string `json:"address"`
		IsHead  bool   `json:"is_head"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	branch := &Branch{ID: id.New("branch"), MerchantID: merchantID, Name: req.Name, NameAM: req.NameAM, Region: req.Region, City: req.City, SubCity: req.SubCity, Address: req.Address, IsHead: req.IsHead}
	if err := h.svc.repo.CreateBranch(r.Context(), branch); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, branch)
}
func (h *Handler) ListBranches(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	list, _ := h.svc.repo.ListBranches(r.Context(), merchantID)
	pkghttp.WriteJSON(w, r, 200, list)
}
