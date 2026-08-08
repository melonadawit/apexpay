package fixedasset

import "github.com/shopspring/decimal"

// Fixed-asset domain types.

type Asset struct {
	ID                      string `json:"id"`
	AssetName               string `json:"asset_name"`
	Category                string `json:"category"`
	AcquisitionDate         string `json:"acquisition_date"`
	Cost                    string `json:"cost"`
	SalvageValue            string `json:"salvage_value"`
	UsefulLifeYears         int    `json:"useful_life_years"`
	DepreciationMethod      string `json:"depreciation_method"`
	DepreciationRate        string `json:"depreciation_rate,omitempty"`
	AccumulatedDepreciation string `json:"accumulated_depreciation"`
	NetBookValue            string `json:"net_book_value"`
	Status                  string `json:"status"`
}

// DepreciationEntry is one period's depreciation amount.
type DepreciationEntry struct {
	ID             string `json:"id"`
	AssetID        string `json:"asset_id"`
	Period         string `json:"period"`
	Amount         string `json:"amount"`
	BookValueAfter string `json:"book_value_after"`
}

// DepreciationParams for pure math.
type DepreciationParams struct {
	Cost            decimal.Decimal
	SalvageValue    decimal.Decimal
	UsefulLifeYears int
	Method          string          // straight_line | declining_balance
	Rate            decimal.Decimal // annual % for declining balance
}
