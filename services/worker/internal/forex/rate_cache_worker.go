package forex

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// RateCacheWorker — Forex Rates cached 60s via Redis per Ethiopia business practice highly regulated by NBE per Ethiopia law
// Per spec: forex_rates from_currency ETB to_currency USD EUR GBP etc rate buy_rate sell_rate source nbe commercial_bank black_market last_updated_at unique from_currency to_currency index pair
// Cached 60s via Redis per Ethiopia business practice highly regulated by NBE, for compliance use NBE official rate
// Outstanding: O(1) Redis cache lookup + O(n) NBE API fetch + upsert DB + cache TTL 60s

type RateCacheWorker struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

type ForexRate struct {
	FromCurrency string  `json:"from_currency"` // ETB
	ToCurrency   string  `json:"to_currency"`   // USD, EUR, GBP
	Rate         float64 `json:"rate"`          // e.g., 1 USD = 57.5 ETB
	BuyRate      float64 `json:"buy_rate"`
	SellRate     float64 `json:"sell_rate"`
	Source       string  `json:"source"` // nbe, commercial_bank
}

func NewRateCacheWorker(pool *pgxpool.Pool, redis *redis.Client) *RateCacheWorker {
	return &RateCacheWorker{pool: pool, redis: redis}
}

// FetchAndCache — fetches from NBE official rate API (mock for demo) and caches in Redis TTL 60s + upserts DB forex_rates
// O(n) where n = number of currency pairs (usually 3-10), optimal for 60s ticker
func (w *RateCacheWorker) FetchAndCache(ctx context.Context) error {
	// Mock NBE official rates — in prod, call https://api.nbe.gov.et/forex/rates or commercial bank API CBE/Awash/Dashen
	rates := []ForexRate{
		{FromCurrency: "ETB", ToCurrency: "USD", Rate: 57.50, BuyRate: 56.80, SellRate: 58.20, Source: "nbe"},
		{FromCurrency: "ETB", ToCurrency: "EUR", Rate: 62.30, BuyRate: 61.50, SellRate: 63.10, Source: "nbe"},
		{FromCurrency: "ETB", ToCurrency: "GBP", Rate: 72.10, BuyRate: 71.20, SellRate: 73.00, Source: "commercial_bank"},
	}

	for _, rate := range rates {
		// Upsert DB forex_rates
		_, err := w.pool.Exec(ctx, `INSERT INTO forex_rates (id, from_currency, to_currency, rate, buy_rate, sell_rate, source, last_updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,now()) ON CONFLICT (from_currency, to_currency) DO UPDATE SET rate=$4, buy_rate=$5, sell_rate=$6, source=$7, last_updated_at=now()`,
			fmt.Sprintf("fx_%s_%s", rate.FromCurrency, rate.ToCurrency), rate.FromCurrency, rate.ToCurrency, rate.Rate, rate.BuyRate, rate.SellRate, rate.Source)
		if err != nil {
			continue // log error but continue for other rates O(n)
		}

		// Cache in Redis TTL 60s O(1) lookup
		cacheKey := fmt.Sprintf("forex_rate:%s_%s", rate.FromCurrency, rate.ToCurrency)
		rateJSON, _ := json.Marshal(rate)
		_ = w.redis.Set(ctx, cacheKey, rateJSON, 60*time.Second).Err()

		// Also cache buy/sell separately for quick lookup
		_ = w.redis.Set(ctx, fmt.Sprintf("forex_rate:%s_%s:buy", rate.FromCurrency, rate.ToCurrency), rate.BuyRate, 60*time.Second).Err()
		_ = w.redis.Set(ctx, fmt.Sprintf("forex_rate:%s_%s:sell", rate.FromCurrency, rate.ToCurrency), rate.SellRate, 60*time.Second).Err()
	}

	return nil
}

// GetCachedRate — O(1) Redis cache lookup per Ethiopia business practice highly regulated by NBE, cached 60s via Redis
func (w *RateCacheWorker) GetCachedRate(ctx context.Context, fromCurrency, toCurrency string) (*ForexRate, error) {
	cacheKey := fmt.Sprintf("forex_rate:%s_%s", fromCurrency, toCurrency)
	val, err := w.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var rate ForexRate
		if json.Unmarshal([]byte(val), &rate) == nil {
			return &rate, nil
		}
	}

	// Cache miss: fetch from DB
	row := w.pool.QueryRow(ctx, `SELECT from_currency, to_currency, rate::text, buy_rate::text, sell_rate::text, source FROM forex_rates WHERE from_currency=$1 AND to_currency=$2 ORDER BY last_updated_at DESC LIMIT 1`, fromCurrency, toCurrency)
	var rate ForexRate
	var rateStr, buyStr, sellStr string
	if err := row.Scan(&rate.FromCurrency, &rate.ToCurrency, &rateStr, &buyStr, &sellStr, &rate.Source); err != nil {
		return nil, err
	}
	// Parse float
	fmt.Sscanf(rateStr, "%f", &rate.Rate)
	fmt.Sscanf(buyStr, "%f", &rate.BuyRate)
	fmt.Sscanf(sellStr, "%f", &rate.SellRate)

	// Re-cache for next request TTL 60s
	rateJSON, _ := json.Marshal(rate)
	_ = w.redis.Set(ctx, cacheKey, rateJSON, 60*time.Second).Err()

	return &rate, nil
}

// RunTicker — runs every 60s per Ethiopia business practice cached 60s via Redis per Ethiopia law highly regulated by NBE
func (w *RateCacheWorker) RunTicker(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// Initial fetch
	_ = w.FetchAndCache(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = w.FetchAndCache(ctx)
		}
	}
}

var _ = time.Now
