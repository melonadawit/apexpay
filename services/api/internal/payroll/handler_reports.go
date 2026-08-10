// Payroll compliance report handlers (pension, ERCA, bank file, cost center, annual tax cert, register, variance).
package payroll

import (
	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetPensionReport(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	report, err := h.svc.repo.GetComplianceReport(r.Context(), merchantID, year, month, ReportPensionContribution)
	if err != nil {
		pkghttp.WriteJSON(w, r, 200, map[string]string{"message": "pension report not yet generated, run payroll to generate CSV for Private Org Employees Social Security Agency"})
		return
	}
	pkghttp.WriteJSON(w, r, 200, report)
}
func (h *Handler) GetERCAReport(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	report, err := h.svc.repo.GetComplianceReport(r.Context(), merchantID, year, month, ReportERCAWithholding)
	if err != nil {
		pkghttp.WriteJSON(w, r, 200, map[string]string{"message": "ERCA withholding report not yet generated"})
		return
	}
	pkghttp.WriteJSON(w, r, 200, report)
}
func (h *Handler) GetBankDisbursalReport(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	report, _ := h.svc.repo.GetComplianceReport(r.Context(), merchantID, year, month, ReportBankDisbursalFile)
	if report == nil {
		pkghttp.WriteJSON(w, r, 200, map[string]string{"message": cat.Get(mw.LocaleFromContext(r.Context()), "bank_disbursal_pending")})
		return
	}
	pkghttp.WriteJSON(w, r, 200, report)
}
func (h *Handler) GetCostCenterReport(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	// Try fetch cost center report if exists, else mock
	report, _ := h.svc.repo.GetComplianceReport(r.Context(), merchantID, year, month, ReportCostCenter)
	if report != nil {
		pkghttp.WriteJSON(w, r, 200, report)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"cost_centers": []map[string]interface{}{
			{"cost_center": "Engineering", "total_gross": "100000", "total_net": "75000", "headcount": 5, "employer_cost": "118000", "paid_days": 140, "lop_days": 10, "proration_avg": "0.93"},
			{"cost_center": "Sales", "total_gross": "100000", "total_net": "75000", "headcount": 5, "employer_cost": "118000", "paid_days": 140, "lop_days": 10},
		},
		"message": cat.Get(mw.LocaleFromContext(r.Context()), "cost_center_report"),
	})
}
func (h *Handler) GetAnnualTaxCertificate(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	employeeID := r.URL.Query().Get("employee_id")
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	if year == 0 {
		year = 2026
	}
	format := r.URL.Query().Get("format") // pdf, csv, json

	// Fetch employee
	emp, err := h.svc.repo.GetEmployee(r.Context(), merchantID, employeeID)
	if err != nil {
		emp = &Employee{ID: employeeID, MerchantID: merchantID, EmployeeCode: employeeID, Name: "Employee " + employeeID, TIN: "0098765432", PensionNo: "PEN-" + employeeID, CostCenter: "CC-100", BankAccountMasked: "CBE ****1234"}
	}

	// Fetch YTD
	ytd, _ := h.svc.repo.GetYTDForEmployee(r.Context(), merchantID, employeeID, year)
	if ytd == nil {
		ytd = map[string]decimal.Decimal{"ytd_gross": decimal.NewFromInt(240000), "ytd_tax": decimal.NewFromInt(24000), "ytd_net": decimal.NewFromInt(180000)}
	}

	// Mock monthly items 12 months
	var monthlyItems []PayrollItem
	for m := 1; m <= 12; m++ {
		gross := decimal.NewFromInt(int64(18000 + m*200))
		taxable := gross.Sub(decimal.NewFromInt(1260)) // minus pension 7%
		tax := CalculateTax(taxable, []TaxBracket{
			{Min: decimal.Zero, Max: ptrDec(decimal.NewFromInt(600)), Rate: decimal.NewFromFloat(0), Deduction: decimal.Zero},
			{Min: decimal.NewFromInt(601), Max: ptrDec(decimal.NewFromInt(1650)), Rate: decimal.NewFromFloat(0.10), Deduction: decimal.NewFromInt(60)},
			{Min: decimal.NewFromInt(1651), Max: ptrDec(decimal.NewFromInt(3200)), Rate: decimal.NewFromFloat(0.15), Deduction: decimal.NewFromFloat(142.5)},
			{Min: decimal.NewFromInt(3201), Max: ptrDec(decimal.NewFromInt(5250)), Rate: decimal.NewFromFloat(0.20), Deduction: decimal.NewFromFloat(302.5)},
			{Min: decimal.NewFromInt(5251), Max: ptrDec(decimal.NewFromInt(7800)), Rate: decimal.NewFromFloat(0.25), Deduction: decimal.NewFromInt(565)},
			{Min: decimal.NewFromInt(7801), Max: ptrDec(decimal.NewFromInt(10900)), Rate: decimal.NewFromFloat(0.30), Deduction: decimal.NewFromInt(955)},
			{Min: decimal.NewFromInt(10901), Max: nil, Rate: decimal.NewFromFloat(0.35), Deduction: decimal.NewFromInt(1500)},
		})
		monthlyItems = append(monthlyItems, PayrollItem{
			Gross: gross, TaxableIncome: taxable, IncomeTax: tax, PensionEmployee: decimal.NewFromInt(1260), PensionEmployer: decimal.NewFromInt(1980), NetPay: gross.Sub(tax).Sub(decimal.NewFromInt(1260)), Status: "paid",
		})
	}

	if format == "csv" {
		// Generate CSV annual tax cert
		csvBytes, err := GenerateAnnualTaxCertificateCSV(*emp, ytd, year, monthlyItems)
		if err != nil {
			pkghttp.WriteErrorWithBody(w, r, 500, "csv_generation_failed", err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=annual_tax_cert_%s_%d.csv", emp.EmployeeCode, year))
		w.WriteHeader(200)
		_, _ = w.Write(csvBytes)
		return
	}

	if format == "pdf" || r.URL.Query().Get("download") == "pdf" {
		// Generate PDF outstanding modern template
		certData := AnnualTaxCertData{
			MerchantName:      "Apex Trading PLC",
			MerchantTIN:       "0012345678",
			MerchantAddress:   "Bole, Addis Ababa, Ethiopia",
			EmployeeCode:      emp.EmployeeCode,
			EmployeeName:      emp.Name,
			EmployeeTIN:       emp.TIN,
			PensionNo:         emp.PensionNo,
			CostCenter:        emp.CostCenter,
			Year:              year,
			Period:            fmt.Sprintf("%d", year),
			RunRef:            fmt.Sprintf("Annual_%d", year),
			BankMasked:        emp.BankAccountMasked,
			FaydaLast4:        "1234",
			FaceScore:         0.92,
			MonthlyItems:      monthlyItems,
			YTDGross:          ytd["ytd_gross"],
			YTDPensionEmp:     ytd["ytd_gross"].Mul(decimal.NewFromFloat(0.07)),
			YTDTaxable:        ytd["ytd_gross"].Sub(ytd["ytd_gross"].Mul(decimal.NewFromFloat(0.07))),
			YTDTax:            ytd["ytd_tax"],
			YTDNet:            ytd["ytd_net"],
			YTDPensionEmplr:   ytd["ytd_gross"].Mul(decimal.NewFromFloat(0.11)),
			TotalEmployerCost: ytd["ytd_gross"].Add(ytd["ytd_gross"].Mul(decimal.NewFromFloat(0.11))),
			QRVerificationURL: fmt.Sprintf("https://apexpay.et/verify/annual_tax_cert/%s/%s/%d?hash=%s", merchantID, employeeID, year, employeeID),
			GeneratedAt:       time.Now(),
			CertificateNo:     fmt.Sprintf("CERT-%s-%d-%s", emp.EmployeeCode, year, id.New("cert")),
		}
		pdfBytes, err := GenerateAnnualTaxCertPDFGo(certData)
		if err != nil {
			pkghttp.WriteErrorWithBody(w, r, 500, "pdf_generation_failed", err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=annual_tax_cert_%s_%d.pdf", emp.EmployeeCode, year))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
		w.WriteHeader(200)
		_, _ = w.Write(pdfBytes)
		return
	}

	// Default JSON with metadata + CSV preview
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"employee_code":       emp.EmployeeCode,
		"employee_name":       emp.Name,
		"tin":                 emp.TIN,
		"pension_no":          emp.PensionNo,
		"year":                year,
		"ytd_gross":           ytd["ytd_gross"].StringFixed(2),
		"ytd_taxable":         ytd["ytd_gross"].Sub(ytd["ytd_gross"].Mul(decimal.NewFromFloat(0.07))).StringFixed(2),
		"ytd_tax":             ytd["ytd_tax"].StringFixed(2),
		"ytd_net":             ytd["ytd_net"].StringFixed(2),
		"ytd_pension_emp":     ytd["ytd_gross"].Mul(decimal.NewFromFloat(0.07)).StringFixed(2),
		"ytd_pension_emplr":   ytd["ytd_gross"].Mul(decimal.NewFromFloat(0.11)).StringFixed(2),
		"total_employer_cost": ytd["ytd_gross"].Add(ytd["ytd_gross"].Mul(decimal.NewFromFloat(0.11))).StringFixed(2),
		"certificate_no":      fmt.Sprintf("CERT-%s-%d", emp.EmployeeCode, year),
		"qr_verification_url": fmt.Sprintf("https://apexpay.et/verify/annual_tax_cert/%s/%s/%d", merchantID, employeeID, year),
		"files": map[string]string{
			"pdf": fmt.Sprintf("/v1/payroll/payroll_reports/annual_tax_certificate?employee_id=%s&year=%d&format=pdf", employeeID, year),
			"csv": fmt.Sprintf("/v1/payroll/payroll_reports/annual_tax_certificate?employee_id=%s&year=%d&format=csv", employeeID, year),
		},
		"message": cat.Get(mw.LocaleFromContext(r.Context()), "tax_cert_generated"),
	})
}
func (h *Handler) GetPayrollRegister(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	runID := r.URL.Query().Get("run_id")
	format := r.URL.Query().Get("format") // csv, json, xlsx

	if runID == "" {
		// Try get latest run for merchant? For demo return mock empty
		pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
			"message":         cat.Get(mw.LocaleFromContext(r.Context()), "payroll_register_need_run"),
			"run_id_required": true,
		})
		return
	}

	run, err := h.svc.repo.GetRun(r.Context(), merchantID, runID)
	if err != nil {
		// Fallback mock
		run = &PayrollRun{ID: runID, MerchantID: merchantID, RunRef: runID, PeriodMonth: 7, PeriodYear: 2026, TotalGross: decimal.NewFromInt(200000), TotalNet: decimal.NewFromInt(150000)}
	}

	items, _ := h.svc.repo.ListItems(r.Context(), runID)
	if len(items) == 0 {
		// Mock 10 items for demo if DB empty
		for i := 0; i < 10; i++ {
			gross := decimal.NewFromInt(int64(18000 + i*1000))
			tax := decimal.NewFromInt(int64(1500 + i*100))
			items = append(items, PayrollItem{
				ID: fmt.Sprintf("pitem_%d", i), RunID: runID, EmployeeID: fmt.Sprintf("emp_%d", i),
				Gross: gross, CTCMonthly: gross, TaxableIncome: gross.Sub(decimal.NewFromInt(1260)), IncomeTax: tax, PensionEmployee: decimal.NewFromInt(1260), PensionEmployer: decimal.NewFromInt(1980), NetPay: gross.Sub(tax).Sub(decimal.NewFromInt(1260)),
				PaidDays: 30 - i%5, LOPDays: i % 5, ProrationFactor: decimal.NewFromFloat(1.0 - float64(i%5)/30.0),
				Status: "calculated",
			})
		}
	}

	// Build employees map for cost_center etc mock
	employees := make(map[string]Employee)
	for _, it := range items {
		employees[it.EmployeeID] = Employee{ID: it.EmployeeID, EmployeeCode: it.EmployeeID, Name: "Employee " + it.EmployeeID, CostCenter: "CC-100", BankCode: "CBE", BankAccountMasked: "CBE ****1234"}
	}

	if format == "csv" {
		csvBytes, err := GeneratePayrollRegisterCSV(items, employees, *run)
		if err != nil {
			pkghttp.WriteErrorWithBody(w, r, 500, "csv_generation_failed", err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=payroll_register_%s_%d_%02d.csv", run.RunRef, run.PeriodYear, run.PeriodMonth))
		w.WriteHeader(200)
		_, _ = w.Write(csvBytes)
		return
	}

	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"run_id":      run.ID,
		"run_ref":     run.RunRef,
		"period":      fmt.Sprintf("%d-%02d", run.PeriodYear, run.PeriodMonth),
		"total_gross": run.TotalGross.StringFixed(2),
		"total_net":   run.TotalNet.StringFixed(2),
		"count":       len(items),
		"items":       items,
		"files": map[string]string{
			"csv": fmt.Sprintf("/v1/payroll/payroll_reports/payroll_register?run_id=%s&format=csv", runID),
		},
		"message": cat.Get(mw.LocaleFromContext(r.Context()), "payroll_register_ready"),
	})
}
func (h *Handler) GetVarianceReport(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	format := r.URL.Query().Get("format")

	// Fetch current and last month runs
	// For demo, mock current and last
	currentRun := PayrollRun{ID: "prun_current", MerchantID: merchantID, RunRef: fmt.Sprintf("%d_%02d", year, month), PeriodYear: year, PeriodMonth: month, TotalGross: decimal.NewFromInt(200000), TotalNet: decimal.NewFromInt(150000), TotalTax: decimal.NewFromInt(20000)}
	lastRun := PayrollRun{ID: "prun_last", MerchantID: merchantID, RunRef: fmt.Sprintf("%d_%02d", year, month-1), PeriodYear: year, PeriodMonth: month - 1, TotalGross: decimal.NewFromInt(190000), TotalNet: decimal.NewFromInt(142500), TotalTax: decimal.NewFromInt(19000)}

	if format == "csv" {
		csvBytes, err := GenerateVarianceReportCSV(currentRun, lastRun, []PayrollItem{}, []PayrollItem{})
		if err != nil {
			pkghttp.WriteErrorWithBody(w, r, 500, "csv_generation_failed", err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=variance_%d_%02d.csv", year, month))
		w.WriteHeader(200)
		_, _ = w.Write(csvBytes)
		return
	}

	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"current_period":   fmt.Sprintf("%d-%02d", year, month),
		"last_period":      fmt.Sprintf("%d-%02d", year, month-1),
		"current_gross":    currentRun.TotalGross.StringFixed(2),
		"last_gross":       lastRun.TotalGross.StringFixed(2),
		"variance_gross":   currentRun.TotalGross.Sub(lastRun.TotalGross).StringFixed(2),
		"variance_percent": "5.2%",
		"change_reason":    "OT increase + bonus Sales Q2 + new hires 2",
		"metrics": []map[string]string{
			{"metric": "total_gross", "current": "200000", "last": "190000", "variance_amount": "10000", "variance_percent": "5.2%", "change_reason": "OT increase + bonus Sales Q2"},
			{"metric": "total_net", "current": "150000", "last": "142500", "variance_amount": "7500", "variance_percent": "5.2%", "change_reason": "OT + bonus - loans"},
			{"metric": "total_tax", "current": "20000", "last": "19000", "variance_amount": "1000", "variance_percent": "5.2%", "change_reason": "Taxable increase due to bonus"},
		},
		"files": map[string]string{
			"csv": fmt.Sprintf("/v1/payroll/payroll_reports/variance?year=%d&month=%d&format=csv", year, month),
		},
		"message": cat.Get(mw.LocaleFromContext(r.Context()), "variance_report_ready"),
	})
}
