package fixedasset

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestStraightLineDepreciation(t *testing.T) {
	// cost 120000, salvage 0, 10 years -> 12000/yr, 1000/mo
	p := DepreciationParams{Cost: decimal.NewFromInt(120000), SalvageValue: decimal.Zero, UsefulLifeYears: 10, Method: "straight_line"}
	annual := AnnualDepreciation(p)
	if !annual.Equal(decimal.NewFromInt(12000)) {
		t.Fatalf("annual = %s, want 12000", annual)
	}
	if !MonthlyDepreciation(annual).Equal(decimal.NewFromInt(1000)) {
		t.Fatalf("monthly = %s, want 1000", MonthlyDepreciation(annual))
	}
}

func TestStraightLineWithSalvage(t *testing.T) {
	// cost 100000, salvage 20000, 8 years -> 10000/yr
	p := DepreciationParams{Cost: decimal.NewFromInt(100000), SalvageValue: decimal.NewFromInt(20000), UsefulLifeYears: 8, Method: "straight_line"}
	annual := AnnualDepreciation(p)
	if !annual.Equal(decimal.NewFromInt(10000)) {
		t.Fatalf("annual = %s, want 10000", annual)
	}
	// NBV never below salvage.
	nbv := StraightLineNBV(p, 20)
	if !nbv.Equal(decimal.NewFromInt(20000)) {
		t.Fatalf("nbv = %s, want 20000 (salvage floor)", nbv)
	}
}

func TestDecliningBalance(t *testing.T) {
	// cost 100000, rate 20% -> first-year 20000
	p := DepreciationParams{Cost: decimal.NewFromInt(100000), UsefulLifeYears: 5, Method: "declining_balance", Rate: decimal.NewFromFloat(0.2)}
	annual := AnnualDepreciation(p)
	if !annual.Equal(decimal.NewFromInt(20000)) {
		t.Fatalf("annual = %s, want 20000", annual)
	}
}

func TestFullyDepreciatedAsset(t *testing.T) {
	// cost == salvage -> no depreciation
	p := DepreciationParams{Cost: decimal.NewFromInt(5000), SalvageValue: decimal.NewFromInt(5000), UsefulLifeYears: 5}
	if !AnnualDepreciation(p).Equal(decimal.Zero) {
		t.Fatal("fully depreciated asset should have zero annual depreciation")
	}
}
