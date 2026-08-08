package fixedasset

import (
	"context"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// CreateAsset inserts an asset with its initial net book value = cost.
func (r *Repository) CreateAsset(ctx context.Context, merchantID string, a *Asset) error {
	a.ID = id.New("fa")
	cost, _ := decimal.NewFromString(a.Cost)
	if a.SalvageValue == "" {
		a.SalvageValue = "0"
	}
	if a.DepreciationMethod == "" {
		a.DepreciationMethod = "straight_line"
	}
	if a.Status == "" {
		a.Status = "active"
	}
	a.AccumulatedDepreciation = "0"
	a.NetBookValue = cost.String()
	_, err := r.pool.Exec(ctx, `
		INSERT INTO fixed_assets (id, merchant_id, asset_name, category, acquisition_date, cost, salvage_value,
			useful_life_years, depreciation_method, depreciation_rate, accumulated_depreciation, net_book_value, status, notes)
		VALUES ($1,$2,$3,$4,$5::date,$6::numeric,$7::numeric,$8,$9,$10::numeric,0,$11::numeric,$12,$13)`,
		a.ID, merchantID, a.AssetName, a.Category, a.AcquisitionDate, a.Cost, a.SalvageValue,
		a.UsefulLifeYears, a.DepreciationMethod, a.DepreciationRate, cost.String(), a.Status, "")
	return err
}

// ListAssets returns a merchant's assets.
func (r *Repository) ListAssets(ctx context.Context, merchantID string) ([]Asset, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, asset_name, category, to_char(acquisition_date,'YYYY-MM-DD'), cost::text, salvage_value::text,
			useful_life_years, depreciation_method, COALESCE(depreciation_rate,0)::text,
			accumulated_depreciation::text, net_book_value::text, status
		FROM fixed_assets WHERE merchant_id=$1 ORDER BY acquisition_date DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Asset{}
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.AssetName, &a.Category, &a.AcquisitionDate, &a.Cost, &a.SalvageValue,
			&a.UsefulLifeYears, &a.DepreciationMethod, &a.DepreciationRate, &a.AccumulatedDepreciation,
			&a.NetBookValue, &a.Status); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

// Depreciate computes and records a depreciation entry for an asset for the current period.
func (r *Repository) Depreciate(ctx context.Context, merchantID, assetID string) (*DepreciationEntry, error) {
	// Load asset.
	var a Asset
	var cost, salvage, rate string
	var years int
	err := r.pool.QueryRow(ctx, `SELECT id, cost::text, salvage_value::text, useful_life_years, depreciation_method, COALESCE(depreciation_rate,0)::text FROM fixed_assets WHERE merchant_id=$1 AND id=$2`,
		merchantID, assetID).Scan(&a.ID, &cost, &salvage, &years, &a.DepreciationMethod, &rate)
	if err != nil {
		return nil, err
	}
	annual := AnnualDepreciation(DepreciationParams{
		Cost: mustDec(cost), SalvageValue: mustDec(salvage), UsefulLifeYears: years,
		Method: a.DepreciationMethod, Rate: mustDec(rate),
	})
	amount := MonthlyDepreciation(annual)
	if amount.LessThanOrEqual(decimal.Zero) {
		amount = decimal.Zero
	}

	period := time.Now().In(tzEAT()).Format("2006-01")
	entry := &DepreciationEntry{ID: id.New("dep"), AssetID: assetID, Period: period, Amount: amount.String()}

	// Insert entry (idempotent on asset+period) and update accumulated NBV.
	err = r.pool.QueryRow(ctx, `
		INSERT INTO depreciation_entries (id, asset_id, merchant_id, period, amount, book_value_after)
		SELECT $1, $2, $3, $4, $5::numeric, GREATEST(net_book_value - $5::numeric, salvage_value)
		FROM fixed_assets WHERE id=$2
		ON CONFLICT (asset_id, period) DO NOTHING
		RETURNING book_value_after::text`, entry.ID, assetID, merchantID, period, amount.String()).Scan(&entry.BookValueAfter)
	if err != nil {
		// conflict -> already recorded this period; return existing
		_ = r.pool.QueryRow(ctx, `SELECT amount::text, book_value_after::text FROM depreciation_entries WHERE asset_id=$1 AND period=$2`,
			assetID, period).Scan(&entry.Amount, &entry.BookValueAfter)
		return entry, nil
	}
	_, _ = r.pool.Exec(ctx, `UPDATE fixed_assets SET accumulated_depreciation = accumulated_depreciation + $2, net_book_value = $3, updated_at=now() WHERE id=$1`,
		assetID, amount.String(), entry.BookValueAfter)
	return entry, nil
}

// Entries returns depreciation history for an asset or merchant.
func (r *Repository) Entries(ctx context.Context, merchantID, assetID string) ([]DepreciationEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, asset_id, period, amount::text, book_value_after::text
		FROM depreciation_entries WHERE merchant_id=$1 AND ($2='' OR asset_id=$2)
		ORDER BY period DESC LIMIT 100`, merchantID, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []DepreciationEntry{}
	for rows.Next() {
		var e DepreciationEntry
		if err := rows.Scan(&e.ID, &e.AssetID, &e.Period, &e.Amount, &e.BookValueAfter); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func mustDec(s string) decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return d
}

func tzEAT() *time.Location { return time.FixedZone("EAT", 3*3600) }
