package payroll

import (
	"context"
	"sort"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/ledger"
	"apexpay/internal/platform/errors"

	"github.com/shopspring/decimal"
)

type Repository interface {
	CreateEmployee(ctx context.Context, e *Employee) error
	ListEmployees(ctx context.Context, merchantID string) ([]Employee, error)
	GetEmployee(ctx context.Context, merchantID, employeeID string) (*Employee, error)

	CreateRun(ctx context.Context, r *PayrollRun) error
	GetRun(ctx context.Context, merchantID, runID string) (*PayrollRun, error)
	UpdateRunStatus(ctx context.Context, runID string, status RunStatus, totals map[string]decimal.Decimal) error

	BulkCreateItems(ctx context.Context, items []PayrollItem) error
	ListItems(ctx context.Context, runID string) ([]PayrollItem, error)

	GetTaxBrackets(ctx context.Context) ([]TaxBracket, error) // ordered by Min asc

	CreateRunBookTx(ctx context.Context, run *PayrollRun, journal *ledger.Journal, entries []ledger.Entry) error
}

type Service struct {
	repo   Repository
	ledger *ledger.Service
}

func NewService(repo Repository, ledgerSvc *ledger.Service) *Service {
	return &Service{repo: repo, ledger: ledgerSvc}
}

// CalculateTax - optimal O(log n) binary search over sorted brackets + decimal precise
func CalculateTax(taxable decimal.Decimal, brackets []TaxBracket) decimal.Decimal {
	if taxable.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero
	}
	// brackets sorted ascending Min
	// Binary search
	idx := sort.Search(len(brackets), func(i int) bool {
		b := brackets[i]
		if b.Max == nil {
			return true // last bracket catches all
		}
		return taxable.LessThan(*b.Max)
	})
	if idx >= len(brackets) {
		idx = len(brackets) - 1
	}
	br := brackets[idx]
	if taxable.LessThan(br.Min) {
		return decimal.Zero
	}
	// ET formula: tax = taxable*rate - deduction
	tax := taxable.Mul(br.Rate).Sub(br.Deduction)
	if tax.LessThan(decimal.Zero) {
		tax = decimal.Zero
	}
	return tax.Round(2)
}

// CalculatePayroll - core algorithm: optimal data structure with decimal, no float
func (s *Service) CalculateRun(ctx context.Context, merchantID, runID string) error {
	r, err := s.repo.GetRun(ctx, merchantID, runID)
	if err != nil {
		return errors.NotFound("payroll run not found")
	}
	if r.Status != StatusDraft && r.Status != StatusCalculating {
		return errors.Validation("run must be draft to calculate")
	}

	_ = s.repo.UpdateRunStatus(ctx, runID, StatusCalculating, nil)

	employees, err := s.repo.ListEmployees(ctx, merchantID)
	if err != nil {
		return err
	}
	brackets, err := s.repo.GetTaxBrackets(ctx)
	if err != nil {
		return err
	}

	items := make([]PayrollItem, 0, len(employees))
	totalGross := decimal.Zero
	totalDeductions := decimal.Zero
	totalNet := decimal.Zero
	totalTax := decimal.Zero
	totalPension := decimal.Zero

	for _, emp := range employees {
		if emp.Status != "active" {
			continue
		}
		// Simplified: gross = base + allowances (OT etc would be input per item request)
		gross := emp.BaseSalary
		// Example OT 5h weekday for demo
		// In real API, OT hours come from request payload per employee

		pensionEmp := gross.Mul(decimal.NewFromFloat(0.07)).Round(2) // 7%
		pensionEmplr := gross.Mul(decimal.NewFromFloat(0.11)).Round(2) // 11%

		// Taxable per ET: gross - pensionEmp (non-taxable portion) - other exemptions
		taxable := gross.Sub(pensionEmp)
		incomeTax := CalculateTax(taxable, brackets)

		deductions := pensionEmp.Add(incomeTax)
		net := gross.Sub(deductions)

		item := PayrollItem{
			ID: id.NewPayrollItem(), RunID: runID, EmployeeID: emp.ID,
			Gross: gross, TaxableIncome: taxable, IncomeTax: incomeTax,
			PensionEmployee: pensionEmp, PensionEmployer: pensionEmplr,
			NetPay: net, Status: "calculated",
		}
		items = append(items, item)
		totalGross = totalGross.Add(gross)
		totalDeductions = totalDeductions.Add(deductions)
		totalNet = totalNet.Add(net)
		totalTax = totalTax.Add(incomeTax)
		totalPension = totalPension.Add(pensionEmp.Add(pensionEmplr))
	}

	if err := s.repo.BulkCreateItems(ctx, items); err != nil {
		_ = s.repo.UpdateRunStatus(ctx, runID, StatusFailed, nil)
		return err
	}

	totals := map[string]decimal.Decimal{
		"total_gross":      totalGross,
		"total_deductions": totalDeductions,
		"total_net":        totalNet,
		"total_tax":        totalTax,
		"total_pension":    totalPension,
	}
	if err := s.repo.UpdateRunStatus(ctx, runID, StatusPendingApproval, totals); err != nil {
		return err
	}

	// Ledger M4: draft posting to payroll_run book
	// Dr expense:salary totalGross Cr liability:payroll_payable totalNet Cr tax payable totalTax Cr pension payable totalPension_emp+emplr
	// Actual journal created on approve/disburse step per DATABASE

	return nil
}

func (s *Service) ApproveRun(ctx context.Context, merchantID, runID, approverID string) error {
	r, err := s.repo.GetRun(ctx, merchantID, runID)
	if err != nil {
		return err
	}
	if r.Status != StatusPendingApproval {
		return errors.Validation("run not pending approval")
	}
	// Dual approval check if >100k ETB
	if r.TotalNet.GreaterThan(decimal.NewFromInt(100000)) {
		// Need second approver logic - simplified: if approvedBy already set and different
		if r.ApprovedBy != nil && *r.ApprovedBy != approverID {
			// second approval completes
		} else if r.ApprovedBy == nil {
			// first approval - stay pending but record approver
			// In real impl, would track approvals table
			return s.repo.UpdateRunStatus(ctx, runID, StatusPendingApproval, nil) // placeholder
		}
	}

	return s.repo.UpdateRunStatus(ctx, runID, StatusApproved, nil)
}

func (s *Service) DisburseRun(ctx context.Context, merchantID, runID string) error {
	r, err := s.repo.GetRun(ctx, merchantID, runID)
	if err != nil {
		return err
	}
	if r.Status != StatusApproved {
		return errors.Validation("run must be approved to disburse")
	}

	// Create ledger book per run if not exists + post M4
	// Then create payout batch for employees banks

	journal := &ledger.Journal{
		ID: id.NewLedgerJournal(),
		BookID: func() string {
			if r.BookID != nil {
				return *r.BookID
			}
			return id.NewLedgerBook()
		}(),
		PostingKey: "payroll_run:" + r.ID,
		Memo:       "payroll run " + r.RunRef,
		ReferenceType: "payroll_run",
		ReferenceID: r.ID,
	}

	entries := []ledger.Entry{
		{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "expense:salary", Direction: "debit", Amount: r.TotalGross, Currency: "ETB"},
		{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "liability:payroll_payable", Direction: "credit", Amount: r.TotalNet, Currency: "ETB"},
		{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "liability:et_income_tax_payable", Direction: "credit", Amount: r.TotalTax, Currency: "ETB"},
		{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "liability:pension_payable", Direction: "credit", Amount: r.TotalPension, Currency: "ETB"},
	}

	if !ledger.ValidateBalanced(entries) {
		return errors.Validation("ledger entries not balanced for payroll run")
	}

	if err := s.repo.CreateRunBookTx(ctx, r, journal, entries); err != nil {
		return err
	}

	return s.repo.UpdateRunStatus(ctx, runID, StatusProcessing, nil)
}
