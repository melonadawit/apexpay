package payroll

import (
	"time"

	"github.com/shopspring/decimal"
)

// ==================== Enums ====================

type EmploymentType string

const (
	EmploymentPermanent EmploymentType = "permanent"
	EmploymentContract  EmploymentType = "contract"
	EmploymentPartTime  EmploymentType = "part_time"
	EmploymentIntern    EmploymentType = "intern"
)

type ConfirmationStatus string

const (
	ConfirmationProbation  ConfirmationStatus = "probation"
	ConfirmationConfirmed  ConfirmationStatus = "confirmed"
	ConfirmationNotice     ConfirmationStatus = "notice"
	ConfirmationTerminated ConfirmationStatus = "terminated"
)

type RunType string

const (
	RunRegular    RunType = "regular"
	RunOffCycle   RunType = "off_cycle"
	RunBonus      RunType = "bonus"
	RunAdjustment RunType = "adjustment"
	RunFnF        RunType = "full_and_final"
)

type RunStatus string

const (
	StatusDraft           RunStatus = "draft"
	StatusCalculating     RunStatus = "calculating"
	StatusPendingApproval RunStatus = "pending_approval"
	StatusApproved        RunStatus = "approved"
	StatusProcessing      RunStatus = "processing"
	StatusCompleted       RunStatus = "completed"
	StatusFailed          RunStatus = "failed"
	StatusVoided          RunStatus = "voided"
)

// OT Rates per ET labour law
// Weekday 1.25x, Weekend 1.5x, Holiday 2x, Night 1.3x
type OTType string

const (
	OTWeekday OTType = "weekday"
	OTWeekend OTType = "weekday_weekend"
	OTHoliday OTType = "holiday"
	OTNight   OTType = "night"
)

var OTRates = map[OTType]decimal.Decimal{
	OTWeekday: decimal.NewFromFloat(1.25),
	OTWeekend: decimal.NewFromFloat(1.5),
	OTHoliday: decimal.NewFromFloat(2.0),
	OTNight:   decimal.NewFromFloat(1.3),
}

// Component types for salary structure — enterprise-grade
type ComponentType string

const (
	ComponentEarning              ComponentType = "earning"
	ComponentDeduction            ComponentType = "deduction"
	ComponentEmployerContribution ComponentType = "employer_contribution"
	ComponentReimbursement        ComponentType = "reimbursement"
)

type CalculationType string

const (
	CalcFixed             CalculationType = "fixed"
	CalcPercentageOfBasic CalculationType = "percentage_of_basic"
	CalcPercentageOfCTC   CalculationType = "percentage_of_ctc"
	CalcPercentageOfGross CalculationType = "percentage_of_gross"
	CalcFormula           CalculationType = "formula"
)

type LoanType string

const (
	LoanPersonal      LoanType = "personal"
	LoanSalaryAdvance LoanType = "salary_advance"
	LoanHousing       LoanType = "housing"
	LoanEducation     LoanType = "education"
	LoanMedical       LoanType = "medical"
	LoanOther         LoanType = "other"
)

type LoanStatus string

const (
	LoanDraft           LoanStatus = "draft"
	LoanPendingApproval LoanStatus = "pending_approval"
	LoanApproved        LoanStatus = "approved"
	LoanActive          LoanStatus = "active"
	LoanClosed          LoanStatus = "closed"
	LoanRejected        LoanStatus = "rejected"
	LoanWrittenOff      LoanStatus = "written_off"
)

type ClaimType string

const (
	ClaimExpense ClaimType = "expense"
	ClaimMedical ClaimType = "medical"
	ClaimTravel  ClaimType = "travel"
	ClaimOther   ClaimType = "other"
)

type ReportType string

const (
	ReportPensionContribution ReportType = "pension_contribution"
	ReportERCAWithholding     ReportType = "erca_withholding"
	ReportAnnualTaxCert       ReportType = "annual_tax_certificate"
	ReportPensionChallan      ReportType = "pension_challan"
	ReportBankDisbursalFile   ReportType = "bank_disbursal_file"
	ReportPayrollRegister     ReportType = "payroll_register"
	ReportCostCenter          ReportType = "cost_center_report"
	ReportVariance            ReportType = "variance_report"
)

// ==================== Organizational Structure ====================

type Department struct {
	ID          string
	MerchantID  string
	Name        string
	NameAM      string
	Code        string
	CostCenter  string
	Description string
	CreatedAt   time.Time
}

type Designation struct {
	ID          string
	MerchantID  string
	Title       string
	TitleAM     string
	Level       int
	Description string
	CreatedAt   time.Time
}

type Grade struct {
	ID          string
	MerchantID  string
	Name        string
	NameAM      string
	MinSalary   decimal.Decimal
	MaxSalary   decimal.Decimal
	Description string
	CreatedAt   time.Time
}

type Branch struct {
	ID         string
	MerchantID string
	Name       string
	NameAM     string
	Region     string
	City       string
	SubCity    string
	Address    string
	IsHead     bool
	CreatedAt  time.Time
}

// ==================== Salary Structure — CTC Template ====================

type SalaryStructure struct {
	ID            string
	MerchantID    string
	Name          string
	NameAM        string
	Description   string
	CTCAnnual     decimal.Decimal
	CTCMonthly    decimal.Decimal
	Currency      string
	EffectiveFrom time.Time
	EffectiveTo   *time.Time
	Status        string // active, archived, draft
	IsDefault     bool
	CreatedBy     *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Components    []StructureComponent // populated via join
}

type StructureComponent struct {
	ID              string
	StructureID     string
	ComponentType   ComponentType
	Code            string // BASIC, HOUSING, TRANSPORT, FUEL, SPECIAL_ALLOW, MEDICAL, etc.
	Name            string
	NameAM          string
	CalculationType CalculationType
	Amount          decimal.Decimal // for fixed
	Percentage      decimal.Decimal // 40.00 = 40%
	Formula         string          // e.g., "CTC_MONTHLY * 0.4" or "BASIC * 0.1"
	IsTaxable       bool
	IsPartOfGross   bool
	IsProratable    bool
	IsPensionable   bool
	IsOptional      bool
	TaxExemptLimit  decimal.Decimal
	OrderNo         int
	Meta            map[string]interface{}
	CreatedAt       time.Time
}

// Component breakdown for payslip
type EarningsBreakdown struct {
	Code         string          `json:"code"`
	Name         string          `json:"name"`
	NameAM       string          `json:"name_am,omitempty"`
	Amount       decimal.Decimal `json:"amount"`
	IsTaxable    bool            `json:"is_taxable"`
	IsProratable bool            `json:"is_proratable,omitempty"`
	// TaxExemptLimit is the ETB amount of this component excluded from taxable income
	// (e.g. MEDICAL/TRANSPORT allowances). Applied up to the component amount.
	TaxExemptLimit decimal.Decimal `json:"tax_exempt_limit,omitempty"`
}

type DeductionsBreakdown struct {
	Code     string          `json:"code"`
	Name     string          `json:"name"`
	Amount   decimal.Decimal `json:"amount"`
	IsPreTax bool            `json:"is_pre_tax,omitempty"`
}

type EmployerContributionsBreakdown struct {
	Code   string          `json:"code"`
	Name   string          `json:"name"`
	Amount decimal.Decimal `json:"amount"`
	Rate   decimal.Decimal `json:"rate,omitempty"`
}

// ==================== Employee Enhanced ====================

type Employee struct {
	ID                  string
	MerchantID          string
	EmployeeCode        string
	Name                string
	NameAM              string
	Email               string
	Phone               string
	TIN                 string
	FinHash             string // Fayda hashed
	PensionNo           string
	BankAccountMasked   string
	BankAccountHash     string
	BankCode            string
	BankAccountName     string
	BaseSalary          decimal.Decimal
	CTCAnnual           decimal.Decimal
	CTCMonthly          decimal.Decimal
	EmploymentDate      time.Time
	DateOfJoining       time.Time
	EmploymentType      EmploymentType
	ConfirmationStatus  ConfirmationStatus
	DepartmentID        *string
	DesignationID       *string
	GradeID             *string
	BranchID            *string
	ReportingManagerID  *string
	SalaryStructureID   *string
	ProbationEndDate    *time.Time
	CostCenter          string
	Status              string // active, inactive, terminated, on_leave, on_hold
	EmploymentStatusExt string // active, on_hold, notice_period, terminated, retired
	Nationality         string // ET
	Gender              string
	Address             string
	City                string
	Region              string
	IsFaydaVerified     bool
	FaydaVerifiedAt     *time.Time
	Documents           []EmployeeDocument // JSON
	EmploymentHistory   []EmploymentHistoryItem
	CreatedAt           time.Time
	UpdatedAt           time.Time

	// Populated joins
	Department  *Department
	Designation *Designation
	Grade       *Grade
	Branch      *Branch
	Structure   *SalaryStructure
}

type EmployeeDocument struct {
	Type     string `json:"type"` // contract, tin_certificate, fayda_front, bank_letter
	FileKey  string `json:"file_key"`
	FileHash string `json:"file_hash"`
	Status   string `json:"status"` // pending, verified, rejected
	MimeType string `json:"mime_type,omitempty"`
	Size     int    `json:"size,omitempty"`
}

type EmploymentHistoryItem struct {
	Action        string    `json:"action"` // joined, promoted, salary_revision, transferred
	From          string    `json:"from,omitempty"`
	To            string    `json:"to,omitempty"`
	EffectiveDate time.Time `json:"effective_date"`
	Reason        string    `json:"reason,omitempty"`
	ApprovedBy    string    `json:"approved_by,omitempty"`
}

// ==================== Salary Revision + Arrears ====================

type SalaryRevision struct {
	ID             string
	MerchantID     string
	EmployeeID     string
	OldBase        decimal.Decimal
	NewBase        decimal.Decimal
	OldCTC         decimal.Decimal
	NewCTC         decimal.Decimal
	OldStructureID *string
	NewStructureID *string
	EffectiveFrom  time.Time
	Reason         string
	ApprovedBy     *string
	Status         string // pending, approved, rejected
	ArrearAmount   decimal.Decimal
	ArrearMonths   int
	CreatedAt      time.Time
}

// ==================== Attendance & Variable Inputs ====================

type AttendanceInput struct {
	ID             string
	RunID          string
	EmployeeID     string
	PaidDays       int
	LOPDays        int
	TotalDays      int
	PresentDays    int
	OTWeekdayHours decimal.Decimal
	OTWeekendHours decimal.Decimal
	OTHolidayHours decimal.Decimal
	OTNightHours   decimal.Decimal
	LeaveTaken     map[string]int // annual, sick, maternity
	LeaveBalance   map[string]int
	IsOnHold       bool
	HoldReason     string
	CreatedAt      time.Time
}

type VariableInput struct {
	ID            string
	RunID         string
	EmployeeID    string
	ComponentCode string // COMMISSION, BONUS, PENALTY, ARREAR, THIRTEENTH_MONTH, OVERTIME
	Amount        decimal.Decimal
	IsTaxable     bool
	IsPensionable bool
	Description   string
	CreatedBy     *string
	CreatedAt     time.Time
}

// ==================== Loans & Advances ====================

type Loan struct {
	ID           string
	MerchantID   string
	EmployeeID   string
	LoanType     LoanType
	Principal    decimal.Decimal
	InterestRate decimal.Decimal
	TenureMonths int
	EMIAmount    decimal.Decimal
	TotalPaid    decimal.Decimal
	Outstanding  decimal.Decimal
	Status       LoanStatus
	DisbursedAt  *time.Time
	NextDueDate  *time.Time
	ApprovedBy   *string
	Reason       string
	Meta         map[string]interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type LoanRepayment struct {
	ID                 string
	LoanID             string
	RunID              *string
	EmployeeID         string
	Amount             decimal.Decimal
	PrincipalComponent decimal.Decimal
	InterestComponent  decimal.Decimal
	OutstandingAfter   decimal.Decimal
	Status             string
	CreatedAt          time.Time
}

// ==================== Payroll Run Enhanced ====================

type PayrollRun struct {
	ID                   string                 `json:"id"`
	MerchantID           string                 `json:"merchant_id"`
	BookID               *string                `json:"book_id"` // each run gets own ledger book per DATABASE
	RunRef               string                 `json:"run_ref"`
	PeriodMonth          int                    `json:"period_month"`
	PeriodYear           int                    `json:"period_year"`
	Type                 RunType                `json:"type"`
	Status               RunStatus              `json:"status"`
	TotalGross           decimal.Decimal        `json:"total_gross"`
	TotalDeductions      decimal.Decimal        `json:"total_deductions"`
	TotalNet             decimal.Decimal        `json:"total_net"`
	TotalTax             decimal.Decimal        `json:"total_tax"`
	TotalPension         decimal.Decimal        `json:"total_pension"`          // employee
	EmployerTotalPension decimal.Decimal        `json:"employer_total_pension"` // employer 11%
	TotalEmployerCost    decimal.Decimal        `json:"total_employer_cost"`    // gross + employer pension
	TotalEmployeesPaid   int                    `json:"total_employees_paid"`
	TotalEmployeesFailed int                    `json:"total_employees_failed"`
	TotalCount           int                    `json:"total_count"`
	PayCalendarID        *string                `json:"pay_calendar_id"`
	CutoffDate           *time.Time             `json:"cutoff_date"`
	DisbursalDate        *time.Time             `json:"disbursal_date"`
	PayrollData          map[string]interface{} `json:"payroll_data"`    // {total_paid_days, total_lop}
	VarianceReport       map[string]interface{} `json:"variance_report"` // vs last month %
	BankFileKey          *string                `json:"bank_file_key"`
	BankFileHash         *string                `json:"bank_file_hash"`
	LockedAt             *time.Time             `json:"locked_at"`
	ApprovedBy           *string                `json:"approved_by"`
	CreatedBy            *string                `json:"created_by"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

// ==================== Payroll Item Enhanced with Breakdowns ====================

type PayrollItem struct {
	ID                             string                           `json:"id"`
	RunID                          string                           `json:"run_id"`
	EmployeeID                     string                           `json:"employee_id"`
	Gross                          decimal.Decimal                  `json:"gross"`
	CTCMonthly                     decimal.Decimal                  `json:"ctc_monthly"`
	OTHours                        decimal.Decimal                  `json:"ot_hours"`
	OTAmount                       decimal.Decimal                  `json:"ot_amount"`
	Commission                     decimal.Decimal                  `json:"commission"`
	Bonus                          decimal.Decimal                  `json:"bonus"`
	OtherAllowances                decimal.Decimal                  `json:"other_allowances"`
	TaxableIncome                  decimal.Decimal                  `json:"taxable_income"`
	IncomeTax                      decimal.Decimal                  `json:"income_tax"`
	PensionEmployee                decimal.Decimal                  `json:"pension_employee"` // 7%
	PensionEmployer                decimal.Decimal                  `json:"pension_employer"` // 11%
	OtherDeductions                decimal.Decimal                  `json:"other_deductions"`
	NetPay                         decimal.Decimal                  `json:"net_pay"`
	Status                         string                           `json:"status"`
	FailureReason                  string                           `json:"failure_reason"`
	EarningsBreakdown              []EarningsBreakdown              `json:"earnings_breakdown"`
	DeductionsBreakdown            []DeductionsBreakdown            `json:"deductions_breakdown"`
	EmployerContributionsBreakdown []EmployerContributionsBreakdown `json:"employer_contributions_breakdown"`
	YTD                            map[string]decimal.Decimal       `json:"ytd"` // ytd_gross, ytd_tax, ytd_net
	PaidDays                       int                              `json:"paid_days"`
	LOPDays                        int                              `json:"lop_days"`
	ProrationFactor                decimal.Decimal                  `json:"proration_factor"`
	IsOnHold                       bool                             `json:"is_on_hold"`
	HoldReason                     string                           `json:"hold_reason"`
	CreatedAt                      time.Time                        `json:"created_at"`
	UpdatedAt                      time.Time                        `json:"updated_at"`
}

// ==================== Tax Bracket ====================

type TaxBracket struct {
	Min           decimal.Decimal  // inclusive
	Max           *decimal.Decimal // exclusive, nil = infinity
	Rate          decimal.Decimal  // 0.1 = 10%
	Deduction     decimal.Decimal  // fixed deduction per ET law
	EffectiveFrom time.Time
}

// ==================== Compliance Reports ====================

type ComplianceReport struct {
	ID          string                 `json:"id"`
	MerchantID  string                 `json:"merchant_id"`
	PeriodMonth int                    `json:"period_month"`
	PeriodYear  int                    `json:"period_year"`
	ReportType  ReportType             `json:"report_type"`
	FileKey     *string                `json:"file_key"`
	FileHash    *string                `json:"file_hash"`
	Status      string                 `json:"status"` // draft, generated, paid, filed, failed
	Metadata    map[string]interface{} `json:"metadata"`
	GeneratedBy *string                `json:"generated_by"`
	CreatedAt   time.Time              `json:"created_at"`
}

// ==================== Final Settlement F&F ====================

type FinalSettlement struct {
	ID                    string                   `json:"id"`
	MerchantID            string                   `json:"merchant_id"`
	EmployeeID            string                   `json:"employee_id"`
	ResignationDate       time.Time                `json:"resignation_date"`
	LastWorkingDate       time.Time                `json:"last_working_date"`
	NoticePeriodDays      int                      `json:"notice_period_days"`
	NoticeServedDays      int                      `json:"notice_served_days"`
	NoticeShortfallDays   int                      `json:"notice_shortfall_days"`
	LeaveEncashmentDays   decimal.Decimal          `json:"leave_encashment_days"`
	LeaveEncashmentAmount decimal.Decimal          `json:"leave_encashment_amount"`
	SeveranceAmount       decimal.Decimal          `json:"severance_amount"` // per ET labour law Art 39-44
	GratuityAmount        decimal.Decimal          `json:"gratuity_amount"`
	BonusProRata          decimal.Decimal          `json:"bonus_pro_rata"`
	OutstandingLoans      decimal.Decimal          `json:"outstanding_loans"`
	OutstandingAdvances   decimal.Decimal          `json:"outstanding_advances"`
	OtherEarnings         decimal.Decimal          `json:"other_earnings"`
	OtherDeductions       decimal.Decimal          `json:"other_deductions"`
	TotalPayable          decimal.Decimal          `json:"total_payable"`
	TotalDeductions       decimal.Decimal          `json:"total_deductions"`
	NetPayable            decimal.Decimal          `json:"net_payable"`
	Status                string                   `json:"status"` // draft, pending_approval, approved, paid, rejected
	ClearanceChecklist    []ClearanceItem          `json:"clearance_checklist"`
	ClearanceDetailed     []map[string]interface{} `json:"clearance_items_detailed"`
	AssetsReturned        []map[string]interface{} `json:"assets_returned"`
	ExitInterview         map[string]interface{}   `json:"exit_interview"`
	ApprovedBy            *string                  `json:"approved_by"`
	PaidAt                *time.Time               `json:"paid_at"`
	CreatedAt             time.Time                `json:"created_at"`
	UpdatedAt             time.Time                `json:"updated_at"`
}

type ClearanceItem struct {
	Item      string     `json:"item"`   // laptop, id_card, etc.
	Status    string     `json:"status"` // pending, done
	CheckedBy string     `json:"checked_by,omitempty"`
	CheckedAt *time.Time `json:"checked_at,omitempty"`
}

// ==================== Employee Portal Access ====================

type EmployeePortalAccess struct {
	ID             string
	MerchantID     string
	EmployeeID     string
	MagicTokenHash string
	TokenLast4     string
	ExpiresAt      time.Time
	LastAccessedAt *time.Time
	AccessCount    int
	IsRevoked      bool
	CreatedAt      time.Time
}

// ==================== Audit Logs ====================

type AuditLog struct {
	ID         string
	MerchantID string
	RunID      *string
	EmployeeID *string
	ActorType  string // system, hr, finance, admin, employee
	ActorID    *string
	Action     string // create_employee, salary_revision, calculate_run, etc.
	Details    map[string]interface{}
	IP         string
	RequestID  string
	CreatedAt  time.Time
}

// ==================== Payroll Calendar — Ethiopia Business Practice Cutoff 25th Disbursal 30th Pay Last Day Lock After Disbursal ====================

type PayFrequency string

const (
	PayFrequencyMonthly     PayFrequency = "monthly"
	PayFrequencySemimonthly PayFrequency = "semimonthly"
	PayFrequencyWeekly      PayFrequency = "weekly"
	PayFrequencyBiweekly    PayFrequency = "biweekly"
)

type PayrollCalendar struct {
	ID            string       `json:"id"`
	MerchantID    string       `json:"merchant_id"`
	Name          string       `json:"name"` // e.g., Monthly Payroll Calendar 2026
	Description   string       `json:"description"`
	PayFrequency  PayFrequency `json:"pay_frequency"`
	Year          int          `json:"year"`
	Month         *int         `json:"month"`         // null for weekly
	CutoffDay     int          `json:"cutoff_day"`    // Ethiopia business practice cutoff 25th
	DisbursalDay  int          `json:"disbursal_day"` // disbursal 30th
	PayDay        int          `json:"pay_day"`       // pay date last day of month
	CutoffDate    time.Time    `json:"cutoff_date"`   // actual cutoff date e.g., 2026-07-25
	DisbursalDate time.Time    `json:"disbursal_date"`
	PayDate       time.Time    `json:"pay_date"`
	IsLocked      bool         `json:"is_locked"`
	LockedAt      *time.Time   `json:"locked_at"`
	LockedBy      *string      `json:"locked_by"`
	CreatedBy     *string      `json:"created_by"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// ==================== Leave Management — Ethiopia Labour Proclamation 1156/2019 Art 77/82/86 ====================

type LeaveBalance struct {
	ID               string
	MerchantID       string
	EmployeeID       string
	LeaveType        LeaveType
	Year             int
	EntitledDays     decimal.Decimal // e.g., annual 14 + years-1 up to 35 per Art 77
	UsedDays         decimal.Decimal
	RemainingDays    decimal.Decimal
	CarryForwardDays decimal.Decimal
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type LeaveRequest struct {
	ID                        string
	MerchantID                string
	EmployeeID                string
	LeaveType                 LeaveType
	StartDate                 time.Time
	EndDate                   time.Time
	DaysRequested             decimal.Decimal // e.g., 2 days, 0.5 half day
	Reason                    string
	Status                    LeaveStatus
	ApprovedBy                *string
	ApprovedAt                *time.Time
	RejectionReason           string
	MedicalCertificateFileKey *string // MinIO for sick >3 days per Art 82
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

// ==================== Loan EMI Schedule ====================

type LoanEMISchedule struct {
	ID                 string
	LoanID             string
	InstallmentNo      int
	DueDate            time.Time
	EMIAmount          decimal.Decimal
	PrincipalComponent decimal.Decimal
	InterestComponent  decimal.Decimal
	OutstandingAfter   decimal.Decimal
	Status             string // pending, paid, overdue, skipped
	PaidAt             *time.Time
	RunID              *string
	CreatedAt          time.Time
}

// ==================== Claims Enhanced — Reimbursements MinIO ====================

type ClaimEnhanced struct {
	ID                string
	MerchantID        string
	EmployeeID        string
	RunID             *string
	ClaimType         ClaimType
	Amount            decimal.Decimal
	Description       string
	ReceiptFileKey    *string // MinIO presigned 15m TTL <5MB
	ReceiptFileHash   *string
	Status            string // pending, approved, rejected, paid
	ApprovedByManager *string
	ApprovedByFinance *string
	ManagerApprovedAt *time.Time
	FinanceApprovedAt *time.Time
	RejectionReason   string
	IsTaxable         bool
	IsPensionable     bool
	CreatedAt         time.Time
}
