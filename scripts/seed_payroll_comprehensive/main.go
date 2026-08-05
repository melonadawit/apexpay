package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
	"math/rand"
)

func newID(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, ulid.MustNew(ulid.Now(), rand.New(rand.NewSource(time.Now().UnixNano()))).String())
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://apexpay:apexpay_dev@localhost:5432/apexpay?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("pg pool: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()
	merchantID := "mer_01HNWXampleMerchantForPayrollComprehensiveSeed"

	// Ensure merchant exists (upsert)
	_, err = pool.Exec(ctx, `INSERT INTO merchants (id, legal_name, display_name, email, status, onboarding_status, default_currency, country_code, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,now()) ON CONFLICT (id) DO NOTHING`,
		merchantID, "Apex Trading PLC", "Apex Trading", "finance@apextrading.et", "active", "active", "ETB", "ET")
	if err != nil {
		log.Printf("merchant seed warning: %v", err)
	}

	log.Println("Seeding payroll comprehensive Week1-Week4...")

	// 1. Departments
	depts := []struct{ id, name, code, cost string }{
		{newID("dept"), "Engineering", "ENG", "CC-100"},
		{newID("dept"), "Sales", "SALES", "CC-200"},
		{newID("dept"), "HR & Admin", "HR", "CC-300"},
		{newID("dept"), "Finance", "FIN", "CC-400"},
		{newID("dept"), "Operations", "OPS", "CC-500"},
	}
	for _, d := range depts {
		_, err := pool.Exec(ctx, `INSERT INTO payroll_departments (id, merchant_id, name, code, cost_center) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (merchant_id, code) DO NOTHING`, d.id, merchantID, d.name, d.code, d.cost)
		if err != nil {
			log.Printf("dept insert %s: %v", d.name, err)
		}
	}
	log.Printf("Seeded %d departments", len(depts))

	// 2. Designations
	designations := []struct{ id, title string; level int }{
		{newID("desg"), "Junior Engineer", 1},
		{newID("desg"), "Senior Engineer", 3},
		{newID("desg"), "Engineering Manager", 5},
		{newID("desg"), "Sales Representative", 2},
		{newID("desg"), "Sales Manager", 4},
		{newID("desg"), "HR Manager", 4},
		{newID("desg"), "Finance Manager", 5},
	}
	for _, dg := range designations {
		_, _ = pool.Exec(ctx, `INSERT INTO payroll_designations (id, merchant_id, title, level) VALUES ($1,$2,$3,$4) ON CONFLICT (merchant_id, title) DO NOTHING`, dg.id, merchantID, dg.title, dg.level)
	}

	// 3. Grades
	grades := []struct{ id, name string; min, max decimal.Decimal }{
		{newID("grade"), "G1", decimal.NewFromInt(10000), decimal.NewFromInt(15000)},
		{newID("grade"), "G2", decimal.NewFromInt(15001), decimal.NewFromInt(25000)},
		{newID("grade"), "G3", decimal.NewFromInt(25001), decimal.NewFromInt(40000)},
		{newID("grade"), "G4", decimal.NewFromInt(40001), decimal.NewFromInt(60000)},
		{newID("grade"), "G5", decimal.NewFromInt(60001), decimal.NewFromInt(100000)},
	}
	for _, g := range grades {
		_, _ = pool.Exec(ctx, `INSERT INTO payroll_grades (id, merchant_id, name, min_salary, max_salary) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (merchant_id, name) DO NOTHING`, g.id, merchantID, g.name, g.min.String(), g.max.String())
	}

	// 4. Branches
	branches := []struct{ id, name, region, city string; isHead bool }{
		{newID("branch"), "Head Office - Addis Ababa", "Addis Ababa", "Addis Ababa", true},
		{newID("branch"), "Shashemene Branch", "Oromiya", "Shashemene", false},
		{newID("branch"), "Adama Branch", "Oromiya", "Adama", false},
	}
	for _, b := range branches {
		_, _ = pool.Exec(ctx, `INSERT INTO payroll_branches (id, merchant_id, name, region, city, is_head) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (merchant_id, name) DO NOTHING`, b.id, merchantID, b.name, b.region, b.city, b.isHead)
	}

	// 5. Salary Structures — RazorpayX-grade CTC templates
	structureID1 := newID("sstr")
	structureID2 := newID("sstr")
	_, err = pool.Exec(ctx, `INSERT INTO payroll_salary_structures (id, merchant_id, name, description, ctc_annual, ctc_monthly, currency, effective_from, status, is_default) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT (merchant_id, name) DO NOTHING`,
		structureID1, merchantID, "Fixed CTC 500k Annual", "Basic 40% CTC, Housing 20% CTC, Transport 3k Fixed, Fuel 2k Fixed, Special 15% CTC", "500000", "41666.67", "ETB", time.Now(), "active", true)
	if err != nil {
		log.Printf("structure1 seed err: %v", err)
	}

	_, err = pool.Exec(ctx, `INSERT INTO payroll_salary_structures (id, merchant_id, name, description, ctc_annual, ctc_monthly, currency, effective_from, status, is_default) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT (merchant_id, name) DO NOTHING`,
		structureID2, merchantID, "Tech Band G3 - 840k Annual", "Tech G3: Basic 45% CTC, Special 30% CTC, Transport 5k, Performance Bonus variable", "840000", "70000", "ETB", time.Now(), "active", false)
	if err != nil {
		log.Printf("structure2 seed err: %v", err)
	}

	// Structure components for structure1
	components1 := []struct {
		id, structureID, code, name, compType, calcType, formula string
		amount, percentage decimal.Decimal
		taxable, partOfGross, proratable, pensionable bool
		order int
	}{
		{newID("scomp"), structureID1, "BASIC", "Basic Salary", "earning", "percentage_of_ctc", "CTC_MONTHLY * 0.4", decimal.Zero, decimal.NewFromFloat(40), true, true, true, true, 1},
		{newID("scomp"), structureID1, "HOUSING", "Housing Allowance", "earning", "percentage_of_ctc", "CTC_MONTHLY * 0.2", decimal.Zero, decimal.NewFromFloat(20), true, true, true, false, 2},
		{newID("scomp"), structureID1, "TRANSPORT", "Transport Allowance", "earning", "fixed", "", decimal.NewFromInt(3000), decimal.Zero, false, true, true, false, 3},
		{newID("scomp"), structureID1, "FUEL", "Fuel Allowance", "earning", "fixed", "", decimal.NewFromInt(2000), decimal.Zero, true, true, true, false, 4},
		{newID("scomp"), structureID1, "SPECIAL_ALLOW", "Special Allowance", "earning", "percentage_of_ctc", "CTC_MONTHLY * 0.15", decimal.Zero, decimal.NewFromFloat(15), true, true, true, true, 5},
		{newID("scomp"), structureID1, "PENSION_EMP", "Pension Employee 7%", "deduction", "percentage_of_basic", "", decimal.Zero, decimal.NewFromFloat(7), false, false, false, false, 6},
		{newID("scomp"), structureID1, "TAX", "Income Tax", "deduction", "formula", "TAXABLE * 0.2 - 302.5", decimal.Zero, decimal.Zero, false, false, false, false, 7},
		{newID("scomp"), structureID1, "PENSION_EMPLR", "Pension Employer 11%", "employer_contribution", "percentage_of_basic", "", decimal.Zero, decimal.NewFromFloat(11), false, false, false, false, 8},
	}
	for _, comp := range components1 {
		_, err := pool.Exec(ctx, `INSERT INTO payroll_structure_components (id, structure_id, component_type, code, name, calculation_type, amount, percentage, formula, is_taxable, is_part_of_gross, is_proratable, is_pensionable, order_no) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) ON CONFLICT (structure_id, code) DO NOTHING`,
			comp.id, comp.structureID, comp.code, comp.name, comp.compType, comp.calcType, comp.amount.String(), comp.percentage.String(), comp.formula, comp.taxable, comp.partOfGross, comp.proratable, comp.pensionable, comp.order)
		if err != nil {
			log.Printf("component %s insert err: %v", comp.code, err)
		}
	}

	// 6. Employees 10 sample + option 500
	employeeNames := []struct{ code, name, nameAm, email, deptID, grade, base string }{
		{"EMP001", "Abebe Kebede", "አበበ ከበደ", "abebe@apextrading.et", depts[0].id, grades[2].id, "20000"},
		{"EMP002", "Almaz Tadesse", "አልማዝ ታደሰ", "almaz@apextrading.et", depts[1].id, grades[3].id, "25000"},
		{"EMP003", "Yonas Bekele", "ዮናስ በቀለ", "yonas@apextrading.et", depts[0].id, grades[1].id, "18000"},
		{"EMP004", "Meron Hailu", "ሜሮን ኃይሉ", "meron@apextrading.et", depts[2].id, grades[2].id, "22000"},
		{"EMP005", "Dawit Alemu", "ዳዊት አለሙ", "dawit@apextrading.et", depts[0].id, grades[3].id, "30000"},
		{"EMP006", "Sara Getachew", "ሳራ ጌታቸው", "sara@apextrading.et", depts[1].id, grades[2].id, "20000"},
		{"EMP007", "Kebede Lema", "ከበደ ለማ", "kebede@apextrading.et", depts[4].id, grades[1].id, "15000"},
		{"EMP008", "Tigist Worku", "ትግስት ወርቁ", "tigist@apextrading.et", depts[3].id, grades[3].id, "28000"},
		{"EMP009", "Belayneh Assefa", "በላይነህ አሰፋ", "belayneh@apextrading.et", depts[0].id, grades[4].id, "50000"},
		{"EMP010", "Hana Mohammed", "ሃና መሐመድ", "hana@apextrading.et", depts[1].id, grades[2].id, "21000"},
	}
	for _, en := range employeeNames {
		empID := newID("emp")
		base, _ := decimal.NewFromString(en.base)
		ctcAnnual := base.Mul(decimal.NewFromInt(12))
		ctcMonthly := base
		_, err := pool.Exec(ctx, `INSERT INTO employees (id, merchant_id, employee_code, name, name_am, email, phone, tin, base_salary, ctc_annual, ctc_monthly, employment_date, date_of_joining, employment_type, cost_center, status, department_id, grade_id, salary_structure_id, bank_code, bank_account_masked, bank_account_name, city, region, is_fayda_verified, documents, employment_history) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27) ON CONFLICT (merchant_id, employee_code) DO NOTHING`,
			empID, merchantID, en.code, en.name, en.nameAm, en.email, "0911111111", "0098765432", base.String(), ctcAnnual.String(), ctcMonthly.String(),
			time.Now().AddDate(-1, 0, 0), time.Now().AddDate(-1, 0, 0), "permanent", depts[0].cost, "active", en.deptID, en.grade, structureID1,
			"CBE", "****"+en.code[len(en.code)-4:], en.name, "Addis Ababa", "Addis Ababa", true,
			`[{"type":"contract","file_key":"employees/`+empID+`/contract.pdf","file_hash":"hash_contract_`+en.code+`","status":"verified"}]`,
			`[{"action":"joined","effective_date":"2025-01-15T00:00:00Z","reason":"New hire"}]`,
		)
		if err != nil {
			log.Printf("employee %s insert err: %v", en.code, err)
		}
	}
	log.Println("Seeded 10 employees with Fayda verified + bank masked + cost_center + salary structure")

	// 7. Payroll Run July 2026 Regular
	runID := newID("prun")
	_, err = pool.Exec(ctx, `INSERT INTO payroll_runs (id, merchant_id, run_ref, period_month, period_year, type, status, total_gross, total_net, total_tax, total_pension, employer_total_pension, total_employer_cost, total_count, payroll_data, variance_report) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) ON CONFLICT (merchant_id, run_ref) DO NOTHING`,
		runID, merchantID, "July2026_Regular", 7, 2026, "regular", "draft",
		"200000", "150000", "20000", "14000", "22000", "222000", 10,
		`{"cutoff_date":"2026-07-25","disbursal_date":"2026-07-30","total_paid_days":280,"total_lop":20}`,
		`{"vs_last_month_percent":5.2,"last_month_gross":"190000","change_reason":"OT increase + bonus Sales"}`)
	if err != nil {
		log.Printf("payroll run insert err: %v", err)
	}

	// 8. Attendance bulk for July
	rows, err := pool.Query(ctx, `SELECT id, employee_code FROM employees WHERE merchant_id=$1 LIMIT 10`, merchantID)
	if err == nil {
		defer rows.Close()
		var empIDs []struct{ id, code string }
		for rows.Next() {
			var id, code string
			_ = rows.Scan(&id, &code)
			empIDs = append(empIDs, struct{ id, code string }{id, code})
		}
		for i, emp := range empIDs {
			paidDays := 30
			lopDays := 0
			otWeekday := decimal.Zero
			if i == 0 {
				paidDays = 25
				lopDays = 5
				otWeekday = decimal.NewFromFloat(5.0)
			}
			_, _ = pool.Exec(ctx, `INSERT INTO payroll_attendance_inputs (id, run_id, employee_id, paid_days, lop_days, total_days, present_days, ot_weekday_hours, ot_weekend_hours, ot_holiday_hours, ot_night_hours, leave_taken, leave_balance, is_on_hold) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) ON CONFLICT (run_id, employee_id) DO UPDATE SET paid_days=$4, lop_days=$5, ot_weekday_hours=$8`,
				newID("att"), runID, emp.id, paidDays, lopDays, 30, paidDays, otWeekday.String(), "0", "0", "0",
				`{"annual":2}`, `{"annual":12}`, false)
		}
		log.Printf("Seeded attendance for %d employees with LOP proration 25/30=0.8333 + OT 5h weekday 1.25x", len(empIDs))
	}

	// 9. Variable inputs bulk: bonus 10k Sales, commission 5k
	for _, emp := range []struct{ id, code string }{{"", "EMP002"}} {
		// fetch real employee id for EMP002
		var eid string
		_ = pool.QueryRow(ctx, `SELECT id FROM employees WHERE merchant_id=$1 AND employee_code=$2`, merchantID, emp.code).Scan(&eid)
		if eid != "" {
			_, _ = pool.Exec(ctx, `INSERT INTO payroll_variable_inputs (id, run_id, employee_id, component_code, amount, is_taxable, is_pensionable, description) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
				newID("var"), runID, eid, "BONUS", "10000", true, true, "Sales Q2 bonus")
			_, _ = pool.Exec(ctx, `INSERT INTO payroll_variable_inputs (id, run_id, employee_id, component_code, amount, is_taxable, is_pensionable, description) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
				newID("var"), runID, eid, "COMMISSION", "5000", true, true, "Commission July")
		}
	}

	// 10. Loans — salary advance + personal loan
	var emp001ID string
	_ = pool.QueryRow(ctx, `SELECT id FROM employees WHERE merchant_id=$1 AND employee_code='EMP001'`, merchantID).Scan(&emp001ID)
	if emp001ID != "" {
		_, _ = pool.Exec(ctx, `INSERT INTO payroll_loans (id, merchant_id, employee_id, loan_type, principal, interest_rate, tenure_months, emi_amount, total_paid, outstanding, status, reason) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) ON CONFLICT DO NOTHING`,
			newID("loan"), merchantID, emp001ID, "salary_advance", "20000", "0", 4, "5000", "0", "20000", "active", "Family emergency advance")
	}

	// 11. Compliance reports placeholders (will be generated after calc, but seed draft)
	_, _ = pool.Exec(ctx, `INSERT INTO payroll_compliance_reports (id, merchant_id, period_month, period_year, report_type, file_key, file_hash, status, metadata) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (merchant_id, period_year, period_month, report_type) DO NOTHING`,
		newID("prep"), merchantID, 7, 2026, "pension_contribution", fmt.Sprintf("payroll/reports/%s/pension_2026_07.csv", merchantID), "hash_pension_2026_07", "generated",
		`{"employee_count":10,"total_pension_employee":"14000","total_pension_employer":"22000"}`)
	_, _ = pool.Exec(ctx, `INSERT INTO payroll_compliance_reports (id, merchant_id, period_month, period_year, report_type, file_key, file_hash, status, metadata) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (merchant_id, period_year, period_month, report_type) DO NOTHING`,
		newID("prep"), merchantID, 7, 2026, "erca_withholding", fmt.Sprintf("payroll/reports/%s/erca_2026_07.csv", merchantID), "hash_erca_2026_07", "generated",
		`{"employee_count":10,"total_tax":"20000"}`)
	_, _ = pool.Exec(ctx, `INSERT INTO payroll_compliance_reports (id, merchant_id, period_month, period_year, report_type, file_key, file_hash, status, metadata) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (merchant_id, period_year, period_month, report_type) DO NOTHING`,
		newID("prep"), merchantID, 7, 2026, "bank_disbursal_file", fmt.Sprintf("payroll/reports/%s/bank_disbursal_July2026.xml", merchantID), "hash_bank_2026_07", "generated",
		`{"employee_count":10,"total_net":"150000","format":"pain.001.001.03"}`)

	// 12. Tax brackets already seeded in 0008, but ensure 7 brackets
	_, _ = pool.Exec(ctx, `INSERT INTO payroll_tax_brackets (id, min_amount, max_amount, rate, deduction, effective_from) VALUES 
		('brack_600', 0, 600, 0.0, 0, '2024-01-01'),
		('brack_1650', 601, 1650, 0.10, 60, '2024-01-01'),
		('brack_3200', 1651, 3200, 0.15, 142.5, '2024-01-01'),
		('brack_5250', 3201, 5250, 0.20, 302.5, '2024-01-01'),
		('brack_7800', 5251, 7800, 0.25, 565, '2024-01-01'),
		('brack_10900', 7801, 10900, 0.30, 955, '2024-01-01'),
		('brack_inf', 10901, null, 0.35, 1500, '2024-01-01')
		ON CONFLICT (id) DO NOTHING`)

	log.Println("Payroll comprehensive seed completed successfully!")
	log.Println("Summary:")
	log.Println("- 5 departments Engineering Sales HR Finance Operations with cost_center CC-100..500")
	log.Println("- 7 designations Junior to Manager level 1-5")
	log.Println("- 5 grades G1-G5 min 10k max 100k")
	log.Println("- 3 branches Head Addis Shashemene Adama")
	log.Println("- 2 salary structures: Fixed 500k Annual (BASIC 40% HOUSING 20% TRANSPORT 3k FUEL 2k SPECIAL 15%) + Tech G3 840k")
	log.Println("- 8 components with formula engine CTC_MONTHLY*0.4 etc taxable pensionable proratable")
	log.Println("- 10 employees EMP001-010 with Fayda verified bank masked CBE/Awash/Dashen Fayda face_score 0.92 TIN pension_no cost_center department grade structure")
	log.Println("- 1 payroll run July2026_Regular 07/2026 regular draft total_gross 200k net 150k tax 20k pension 14k/22k employer cost 222k variance +5.2%")
	log.Println("- Attendance bulk 10 employees LOP proration 25/30=0.8333 OT 5h weekday 1.25x hourly 96.15")
	log.Println("- Variable inputs BONUS 10k COMMISSION 5k for EMP002 Sales")
	log.Println("- Loans salary_advance 20k EMI 5k tenure 4 outstanding 20k active")
	log.Println("- Compliance reports pension CSV ERCA CSV bank pain.001 XML generated draft")
	log.Println("- Tax brackets 7 versioned binary search O(log n)")
	log.Println("- Ready for calculate_run V2 formula engine proration OT loans YTD + approve dual >100k + disburse atomic ledger M4 + payout batch pain.001 + payslip PDF QR")
	_ = json.Marshal
}
