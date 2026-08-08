package payroll

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// Ethiopia-law payroll math tests: pension, OT, leave, severance, tax boundaries.

func TestPensionET_Rates(t *testing.T) {
	emp, emplr, total := CalculatePensionET(decimal.NewFromInt(10000), 0, 0)
	require.True(t, emp.Equal(decimal.NewFromFloat(700)), "employee 7%")
	require.True(t, emplr.Equal(decimal.NewFromFloat(1100)), "employer 11%")
	require.True(t, total.Equal(decimal.NewFromFloat(1800)), "total 18%")
}

func TestTaxableIncomeET_ClampsAtZero(t *testing.T) {
	// gross 5000, pension 350, exempt 600 -> 4050
	got := TaxableIncomeET(decimal.NewFromInt(5000), decimal.NewFromInt(350), decimal.NewFromInt(600))
	require.True(t, got.Equal(decimal.NewFromInt(4050)))
	// gross smaller than pension -> 0
	got2 := TaxableIncomeET(decimal.NewFromInt(100), decimal.NewFromInt(350), decimal.NewFromInt(0))
	require.True(t, got2.Equal(decimal.Zero))
}

func TestTaxBrackets_Boundaries(t *testing.T) {
	// The ERCA 2024 progressive brackets (bracket boundaries must not double-tax).
	brackets := []TaxBracket{
		{Min: decimal.Zero, Max: etPtrDec(decimal.NewFromInt(600)), Rate: decimal.NewFromFloat(0), Deduction: decimal.Zero},
		{Min: decimal.NewFromInt(601), Max: etPtrDec(decimal.NewFromInt(1650)), Rate: decimal.NewFromFloat(0.10), Deduction: decimal.NewFromInt(60)},
		{Min: decimal.NewFromInt(1651), Max: etPtrDec(decimal.NewFromInt(3200)), Rate: decimal.NewFromFloat(0.15), Deduction: decimal.NewFromFloat(142.5)},
		{Min: decimal.NewFromInt(3201), Max: etPtrDec(decimal.NewFromInt(5250)), Rate: decimal.NewFromFloat(0.20), Deduction: decimal.NewFromFloat(302.5)},
		{Min: decimal.NewFromInt(5251), Max: etPtrDec(decimal.NewFromInt(7800)), Rate: decimal.NewFromFloat(0.25), Deduction: decimal.NewFromInt(565)},
		{Min: decimal.NewFromInt(7801), Max: etPtrDec(decimal.NewFromInt(10900)), Rate: decimal.NewFromFloat(0.30), Deduction: decimal.NewFromInt(955)},
		{Min: decimal.NewFromInt(10901), Max: nil, Rate: decimal.NewFromFloat(0.35), Deduction: decimal.NewFromInt(1500)},
	}
	cases := []struct {
		taxable float64
		want    float64
	}{
		{600, 0},       // 0% bracket
		{601, 0.1},     // 601*0.10 - 60 = 0.1
		{1650, 105},    // 1650*0.10 - 60 = 105 (upper boundary stays in 10% bracket)
		{1651, 105.15}, // 1651*0.15 - 142.5
		{3200, 337.5},  // 3200*0.15 - 142.5
		{5250, 747.5},  // 5250*0.20 - 302.5
		{7800, 1385},   // 7800*0.25 - 565
		{10900, 2315},  // 10900*0.30 - 955
		{20000, 5500},  // 20000*0.35 - 1500
	}
	for _, c := range cases {
		tax := CalculateTax(decimal.NewFromFloat(c.taxable), brackets)
		require.True(t, tax.Equal(decimal.NewFromFloat(c.want)),
			"taxable %v should be %.2f, got %s", c.taxable, c.want, tax.String())
	}
}

func TestOTRates_ET(t *testing.T) {
	base := decimal.NewFromInt(20800) // hourly = 100 ETB (20800/208)
	// 2h weekday (1.25x) + 2h weekend (1.5x) + 1h holiday (2x) + 2h night (1.3x)
	got := CalculateOTAmountET(base,
		decimal.NewFromInt(2), decimal.NewFromInt(2), decimal.NewFromInt(1), decimal.NewFromInt(2))
	// 2*100*1.25 + 2*100*1.5 + 1*100*2 + 2*100*1.3 = 250 + 300 + 200 + 260 = 1010
	require.True(t, got.Equal(decimal.NewFromInt(1010)))
}

func TestHourlyRateET(t *testing.T) {
	require.True(t, HourlyRateET(decimal.NewFromInt(20800)).Equal(decimal.NewFromInt(100)))
}

func TestAnnualLeaveEntitlement_Art77(t *testing.T) {
	require.Equal(t, 14, AnnualLeaveEntitlementET(1))
	require.Equal(t, 15, AnnualLeaveEntitlementET(2))
	require.Equal(t, 34, AnnualLeaveEntitlementET(21)) // 14 + (21-1) = 34
	require.Equal(t, 35, AnnualLeaveEntitlementET(22)) // 14 + 21 = 35
	require.Equal(t, 35, AnnualLeaveEntitlementET(40), "capped at 35")
}

func TestSickLeaveEntitlement_Art82(t *testing.T) {
	e := SickLeaveEntitlementET()
	require.Equal(t, 30, e.FirstMonth100Pct)
	require.Equal(t, 60, e.Next2Months50Pct)
	require.Equal(t, 90, e.Remaining3Months0Pct)
	require.Equal(t, 180, e.Total6Months)
}

func TestSeverancePay(t *testing.T) {
	// base 20000 * 5 years * 1.0 = 100000
	got := SeverancePayET(decimal.NewFromInt(20000), 5, 1.0)
	require.True(t, got.Equal(decimal.NewFromInt(100000)))
	// double factor (illegal termination) -> 200000
	got2 := SeverancePayET(decimal.NewFromInt(20000), 5, 2.0)
	require.True(t, got2.Equal(decimal.NewFromInt(200000)))
}

func TestLeaveEncashment(t *testing.T) {
	// monthly gross 30000, 3 unused days -> 30000/30*3 = 3000
	got := LeaveEncashmentET(decimal.NewFromInt(30000), decimal.NewFromInt(3))
	require.True(t, got.Equal(decimal.NewFromInt(3000)))
}

func TestValidateTIN(t *testing.T) {
	require.NoError(t, ValidateTIN("0023456789"))
	require.Error(t, ValidateTIN("12345"), "must be 10 digits")
	require.Error(t, ValidateTIN("00abc234567"), "must be numeric")
}
