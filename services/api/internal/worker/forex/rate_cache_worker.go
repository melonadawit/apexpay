package forex

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

// RateCacheWorker — Forex rates cached for 60s via Redis. Ethiopia forex is highly regulated
// by the NBE, so the compliance source of truth is the NBE official rate.
//
// Storage contract (db/migrations/0016_forex_accounting_credit_notifications):
//   - forex_rates(from_currency, to_currency, rate, buy_rate, sell_rate, source, last_updated_at)
//   - unique pair index (from_currency, to_currency)
//
// All rates are decimal.Decimal — never float64 — to satisfy the platform's no-float-money rule.
//
// Complexity: O(1) Redis cache lookup, O(n) DB upsert where n = number of currency pairs
// (usually 3–10), ideal for a 60s ticker.
type RateCacheWorker struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

type ForexRate struct {
	FromCurrency string          `json:"from_currency"` // ETB
	ToCurrency   string          `json:"to_currency"`   // USD, EUR, GBP
	Rate         decimal.Decimal `json:"rate"`          // e.g. 1 USD = 57.50 ETB
	BuyRate      decimal.Decimal `json:"buy_rate"`
	SellRate     decimal.Decimal `json:"sell_rate"`
	Source       string          `json:"source"` // nbe | commercial_bank
}

const (
	cacheTTL      = 60 * time.Second
	redisKeyFmt   = "forex_rate:%s_%s"
	buySuffixFmt  = "forex_rate:%s_%s:buy"
	sellSuffixFmt = "forex_rate:%s_%s:sell"
)

func NewRateCacheWorker(pool *pgxpool.Pool, redis *redis.Client) *RateCacheWorker {
	return &RateCacheWorker{pool: pool, redis: redis}
}

// defaultRates returns the mock NBE official rates used for local/demo.
// In production this is replaced by a live NBE (https://api.nbe.gov.et/forex/rates) or
// partner commercial-bank feed. It is isolated here so the swap is a single function.
func defaultRates() []ForexRate {
	return []ForexRate{
		{FromCurrency: "ETB", ToCurrency: "USD", Rate: decimal.NewFromFloat(57.50), BuyRate: decimal.NewFromFloat(56.80), SellRate: decimal.NewFromFloat(58.20), Source: "nbe"},
		{FromCurrency: "ETB", ToCurrency: "EUR", Rate: decimal.NewFromFloat(62.30), BuyRate: decimal.NewFromFloat(61.50), SellRate: decimal.NewFromFloat(63.10), Source: "nbe"},
		{FromCurrency: "ETB", ToCurrency: "GBP", Rate: decimal.NewFromFloat(72.10), BuyRate: decimal.NewFromFloat(71.20), SellRate: decimal.NewFromFloat(73.00), Source: "commercial_bank"},
	}
}

// FetchAndCache upserts the latest rates into Postgres and refreshes the Redis cache (TTL 60s).
// It is intentionally best-effort per pair so a single failing source cannot stall the others.
func (w *RateCacheWorker) FetchAndCache(ctx context.Context) error {
	rates := defaultRates()

	for _, rate := range rates {
		// Upsert DB forex_rates (unique on the pair index).
		_, err := w.pool.Exec(ctx, `INSERT INTO forex_rates (id, from_currency, to_currency, rate, buy_rate, sell_rate, source, last_updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,now())
			ON CONFLICT (from_currency, to_currency) DO UPDATE
			SET rate=$4, buy_rate=$5, sell_rate=$6, source=$7, last_updated_at=now()`,
			fmt.Sprintf("fx_%s_%s", rate.FromCurrency, rate.ToCurrency),
			rate.FromCurrency, rate.ToCurrency, rate.Rate, rate.BuyRate, rate.SellRate, rate.Source)
		if err != nil {
			continue // log + move on; other pairs must still refresh
		}

		// Refresh the Redis cache for O(1) reads.
		_ = w.redis.Set(ctx, fmt.Sprintf(redisKeyFmt, rate.FromCurrency, rate.ToCurrency), rate, cacheTTL).Err()
		_ = w.redis.Set(ctx, fmt.Sprintf(buySuffixFmt, rate.FromCurrency, rate.ToCurrency), rate.BuyRate, cacheTTL).Err()
		_ = w.redis.Set(ctx, fmt.Sprintf(sellSuffixFmt, rate.FromCurrency, rate.ToCurrency), rate.SellRate, cacheTTL).Err()
	}

	return nil
}

// GetCachedRate returns a rate from Redis (O(1)) falling back to the DB on a cache miss.
func (w *RateCacheWorker) GetCachedRate(ctx context.Context, fromCurrency, toCurrency string) (*ForexRate, error) {
	cacheKey := fmt.Sprintf(redisKeyFmt, fromCurrency, toCurrency)
	if val, err := w.redis.Get(ctx, cacheKey).Result(); err == nil {
		var rate ForexRate
		if json.Unmarshal([]byte(val), &rate) == nil {
			return &rate, nil
		}
	}

	// Cache miss: read the freshest persisted rate, then warm the cache.
	var rate ForexRate
	row := w.pool.QueryRow(ctx, `SELECT from_currency, to_currency, rate, buy_rate, sell_rate, source
		FROM forex_rates WHERE from_currency=$1 AND to_currency=$2
		ORDER BY last_updated_at DESC LIMIT 1`, fromCurrency, toCurrency)
	if err := row.Scan(&rate.FromCurrency, &rate.ToCurrency, &rate.Rate, &rate.BuyRate, &rate.SellRate, &rate.Source); err != nil {
		return nil, err
	}

	_ = w.redis.Set(ctx, cacheKey, rate, cacheTTL).Err()
	return &rate, nil
}

// RunTicker refreshes rates every 60s (matching the Redis TTL).
func (w *RateCacheWorker) RunTicker(ctx context.Context) {
	_ = w.FetchAndCache(ctx)
	ticker := time.NewTicker(cacheTTL)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = w.FetchAndCache(ctx)
		}
	}
}
