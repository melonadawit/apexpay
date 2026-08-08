package fixedasset

import (
	"github.com/shopspring/decimal"
)

// AnnualDepreciation computes the annual depreciation amount for an asset.
// Pure function — unit-testable.
func AnnualDepreciation(p DepreciationParams) decimal.Decimal {
	depreciable := p.Cost.Sub(p.SalvageValue)
	if depreciable.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero
	}
	if p.UsefulLifeYears <= 0 {
		return decimal.Zero
	}
	switch p.Method {
	case "declining_balance":
		rate := p.Rate
		if rate.IsZero() {
			rate = decimal.NewFromFloat(0.2) // default 20%
		}
		// First-year depreciation = cost * rate (book value declines each year).
		return p.Cost.Mul(rate).Round(2)
	default: // straight_line
		return depreciable.Div(decimal.NewFromInt(int64(p.UsefulLifeYears))).Round(2)
	}
}

// MonthlyDepreciation splits the annual amount into 12 even months.
func MonthlyDepreciation(annual decimal.Decimal) decimal.Decimal {
	return annual.Div(decimal.NewFromInt(12)).Round(2)
}

// StraightLineNBV computes the net book value after n full years of straight-line depreciation.
func StraightLineNBV(p DepreciationParams, years int) decimal.Decimal {
	annual := AnnualDepreciation(p)
	if years > p.UsefulLifeYears {
		years = p.UsefulLifeYears
	}
	dep := annual.Mul(decimal.NewFromInt(int64(years)))
	nbv := p.Cost.Sub(dep)
	if nbv.LessThan(p.SalvageValue) {
		nbv = p.SalvageValue
	}
	return nbv
}
