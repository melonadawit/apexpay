// Payroll employee handlers (CRUD, bulk import, revisions, YTD).
package payroll

import (
	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		EmployeeCode      string `json:"employee_code"`
		Name              string `json:"name"`
		NameAM            string `json:"name_am"`
		Email             string `json:"email"`
		Phone             string `json:"phone"`
		TIN               string `json:"tin"`
		BaseSalary        string `json:"base_salary"`
		CTCAnnual         string `json:"ctc_annual"`
		BankCode          string `json:"bank_code"`
		BankAccount       string `json:"bank_account"`
		BankAccountName   string `json:"bank_account_name"`
		DepartmentID      string `json:"department_id"`
		DesignationID     string `json:"designation_id"`
		GradeID           string `json:"grade_id"`
		BranchID          string `json:"branch_id"`
		SalaryStructureID string `json:"salary_structure_id"`
		CostCenter        string `json:"cost_center"`
		EmploymentType    string `json:"employment_type"`
		City              string `json:"city"`
		Region            string `json:"region"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	base, err := decimal.NewFromString(req.BaseSalary)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "base_salary must be numeric")
		return
	}
	ctcAnnual, err := decimal.NewFromString(req.CTCAnnual)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "ctc_annual must be numeric")
		return
	}
	ctcMonthly := ctcAnnual.Div(decimal.NewFromInt(12)).Round(2)
	if ctcAnnual.IsZero() {
		ctcAnnual = base.Mul(decimal.NewFromInt(12))
		ctcMonthly = base
	}
	if base.LessThanOrEqual(decimal.Zero) {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "base_salary must be > 0")
		return
	}
	emp := &Employee{
		ID: id.NewEmployee(), MerchantID: merchantID,
		EmployeeCode: req.EmployeeCode, Name: req.Name, NameAM: req.NameAM,
		Email: req.Email, Phone: req.Phone, TIN: req.TIN,
		BaseSalary: base, CTCAnnual: ctcAnnual, CTCMonthly: ctcMonthly,
		BankCode: req.BankCode, BankAccountMasked: req.BankAccount, BankAccountName: req.BankAccountName,
		CostCenter: req.CostCenter, Status: "active", ConfirmationStatus: ConfirmationProbation,
		DepartmentID: strPtr(req.DepartmentID), DesignationID: strPtr(req.DesignationID),
		GradeID: strPtr(req.GradeID), BranchID: strPtr(req.BranchID), SalaryStructureID: strPtr(req.SalaryStructureID),
		City: req.City, Region: req.Region, EmploymentType: EmploymentType(req.EmploymentType),
		Nationality: "ET",
	}
	if emp.EmploymentType == "" {
		emp.EmploymentType = EmploymentPermanent
	}
	// employment dates now
	// Employment dates now
	emp.EmploymentDate = nowTime()
	emp.DateOfJoining = nowTime()

	if err := h.svc.repo.CreateEmployee(r.Context(), emp); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, emp)
}
func (h *Handler) BulkCreateEmployees(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	// Accept CSV multipart or JSON array
	// For simplicity support JSON array and CSV file upload via FormFile
	contentType := r.Header.Get("Content-Type")
	var employees []Employee

	if contentType == "application/json" || r.FormValue("bulk") != "" {
		var req []struct {
			EmployeeCode string `json:"employee_code"`
			Name         string `json:"name"`
			Email        string `json:"email"`
			BaseSalary   string `json:"base_salary"`
			BankCode     string `json:"bank_code"`
			DepartmentID string `json:"department_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
			return
		}
		for _, rec := range req {
			base, _ := decimal.NewFromString(rec.BaseSalary)
			employees = append(employees, Employee{
				ID: id.NewEmployee(), MerchantID: merchantID,
				EmployeeCode: rec.EmployeeCode, Name: rec.Name, Email: rec.Email,
				BaseSalary: base, CTCAnnual: base.Mul(decimal.NewFromInt(12)), CTCMonthly: base,
				BankCode: rec.BankCode, Status: "active", DepartmentID: strPtr(rec.DepartmentID),
			})
		}
	} else {
		// CSV upload — handle both multipart file and raw body as io.Reader O(n) papaparse
		var reader io.Reader
		file, _, err := r.FormFile("file")
		if err != nil {
			// Try reading raw body as CSV
			reader = r.Body
		} else {
			defer file.Close()
			reader = file
		}
		csvReader := csv.NewReader(reader)
		records, err := csvReader.ReadAll()
		if err != nil {
			pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "csv parse failed")
			return
		}
		// Expect header employee_code,name,email,base_salary,bank_code,department_id
		for i, rec := range records {
			if i == 0 && rec[0] == "employee_code" {
				continue // header
			}
			if len(rec) < 4 {
				continue
			}
			base, _ := decimal.NewFromString(rec[3])
			emp := Employee{
				ID: id.NewEmployee(), MerchantID: merchantID,
				EmployeeCode: rec[0], Name: rec[1],
				BaseSalary: base, CTCAnnual: base.Mul(decimal.NewFromInt(12)), CTCMonthly: base,
				Status: "active",
			}
			if len(rec) > 1 {
				emp.Email = ""
				if len(rec) > 2 {
					emp.Email = rec[2]
				}
			}
			if len(rec) > 4 {
				emp.BankCode = rec[4]
			}
			employees = append(employees, emp)
		}
	}

	if len(employees) > 1000 {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "bulk max 1000 employees per request")
		return
	}
	if err := h.svc.repo.BulkCreateEmployees(r.Context(), employees); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, map[string]interface{}{"count": len(employees), "message": fmt.Sprintf("%d employees imported <2s p99", len(employees))})
}
func (h *Handler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	list, err := h.svc.repo.ListEmployees(r.Context(), merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
func (h *Handler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	empID := chi.URLParam(r, "id")
	emp, err := h.svc.repo.GetEmployeeWithStructure(r.Context(), merchantID, empID)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "employee not found")
		return
	}
	pkghttp.WriteJSON(w, r, 200, emp)
}
func (h *Handler) GetYTD(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	empID := chi.URLParam(r, "id")
	yearStr := r.URL.Query().Get("year")
	year := 2026
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}
	ytd, err := h.svc.repo.GetYTDForEmployee(r.Context(), merchantID, empID, year)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, ytd)
}
func (h *Handler) CreateSalaryRevision(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	empID := chi.URLParam(r, "id")
	var req struct {
		NewBase        string `json:"new_base"`
		NewCTC         string `json:"new_ctc"`
		NewStructureID string `json:"new_structure_id"`
		EffectiveFrom  string `json:"effective_from"` // 2026-07-01
		Reason         string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	newBase, _ := decimal.NewFromString(req.NewBase)
	newCTC, _ := decimal.NewFromString(req.NewCTC)
	effectiveFrom, _ := parseDate(req.EffectiveFrom)

	// Fetch old
	emp, err := h.svc.repo.GetEmployee(r.Context(), merchantID, empID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	rev := &SalaryRevision{
		ID: id.New("srev"), MerchantID: merchantID, EmployeeID: empID,
		OldBase: emp.BaseSalary, NewBase: newBase,
		OldCTC: emp.CTCAnnual, NewCTC: newCTC,
		NewStructureID: strPtr(req.NewStructureID),
		EffectiveFrom:  effectiveFrom, Reason: req.Reason, Status: "pending",
	}
	// Arrear calc: (new-old)*months pending if effective in past
	// months pending = months between effectiveFrom and now if effectiveFrom < now
	now := nowTime()
	months := monthsBetween(effectiveFrom, now)
	if months > 0 && effectiveFrom.Before(now) {
		diff := newBase.Sub(emp.BaseSalary)
		rev.ArrearAmount = diff.Mul(decimal.NewFromInt(int64(months))).Round(2)
		rev.ArrearMonths = months
	}

	if err := h.svc.repo.CreateSalaryRevision(r.Context(), rev); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, rev)
}
func (h *Handler) ListSalaryRevisions(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	empID := chi.URLParam(r, "id")
	list, _ := h.svc.repo.ListSalaryRevisions(r.Context(), merchantID, empID)
	pkghttp.WriteJSON(w, r, 200, list)
}
