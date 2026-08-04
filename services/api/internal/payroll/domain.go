package payroll

import (
	"time"

	"github.com/shopspring/decimal"
)

type EmploymentType string

const (
	EmploymentPermanent EmploymentType = "permanent"
	EmploymentContract  EmploymentType = "contract"
	EmploymentPartTime  EmploymentType = "part_time"
)

type Employee struct {
	ID                string
	MerchantID        string
	EmployeeCode      string
	Name              string
	NameAM            string
	Email             string
	Phone             string
	TIN               string
	FinHash           string // Fayda
	PensionNo         string
	BankAccountMasked string
	BankAccountHash   string
	BankCode          string
	BaseSalary        decimal.Decimal
	EmploymentDate    time.Time
	EmploymentType    EmploymentType
	CostCenter        string
	Status            string
	CreatedAt         time.Time
}

type RunType string

const (
	RunRegular    RunType = "regular"
	RunOffCycle   RunType = "off_cycle"
	RunBonus      RunType = "bonus"
	RunAdjustment RunType = "adjustment"
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
)

type PayrollRun struct {
	ID              string
	MerchantID      string
	BookID          *string // each run gets own ledger book (per DATABASE)
	RunRef          string
	PeriodMonth     int
	PeriodYear      int
	Type            RunType
	Status          RunStatus
	TotalGross      decimal.Decimal
	TotalDeductions decimal.Decimal
	TotalNet        decimal.Decimal
	TotalTax        decimal.Decimal
	TotalPension    decimal.Decimal
	ApprovedBy      *string
	CreatedAt       time.Time
}

type PayrollItem struct {
	ID              string
	RunID           string
	EmployeeID      string
	Gross           decimal.Decimal
	OTHours         decimal.Decimal
	OTAmount        decimal.Decimal
	Commission      decimal.Decimal
	Bonus           decimal.Decimal
	OtherAllowances decimal.Decimal
	TaxableIncome   decimal.Decimal
	IncomeTax       decimal.Decimal
	PensionEmployee decimal.Decimal // 7%
	PensionEmployer decimal.Decimal // 11%
	OtherDeductions decimal.Decimal
	NetPay          decimal.Decimal
	Status          string
}

// Tax bracket table - DB driven per DATABASE v1.1.0 but struct here for optimal binary search
type TaxBracket struct {
	Min           decimal.Decimal  // inclusive
	Max           *decimal.Decimal // exclusive, nil = infinity
	Rate          decimal.Decimal  // 0.1 = 10%
	Deduction     decimal.Decimal  // fixed deduction per ET law
	EffectiveFrom time.Time
}

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
