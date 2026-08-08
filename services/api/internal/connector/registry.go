package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Resolver resolves the connector to use for a merchant+connector pair, so the payment
// service can pick the right rail per request (connectors are per-merchant).
type Resolver interface {
	Get(ctx context.Context, merchantID, connectorID string) Connector
	ForMerchant(ctx context.Context, merchantID string) (map[string]Connector, error)
}

// Factory builds a Connector from its id and decrypted config.
type Factory func(id string, cfg Config) (Connector, error)

// defaultFactory maps connector ids to their constructors.
var defaultFactory Factory = func(id string, cfg Config) (Connector, error) {
	switch id {
	case "telebirr":
		return NewTelebirr(cfg)
	case "cbe_birr":
		return NewCBEBirr(cfg)
	case "amole":
		return NewAmole(cfg)
	case "ethswitch":
		return NewEthSwitch(cfg)
	case "card_acquirer":
		return NewCardAcquirer(cfg)
	case "mock":
		return NewMock(), nil
	default:
		return nil, fmt.Errorf("connector factory: unknown connector %q", id)
	}
}

// Registry builds connectors from the connector_configs table. A row may be global
// (merchant_id NULL) or merchant-specific; the registry prefers merchant-specific and
// falls back to global. Secrets stored in the config column are decrypted via decryptFn.
type Registry struct {
	pool        *pgxpool.Pool
	environment string
	factory     Factory
	decryptFn   func([]byte) ([]byte, error)

	mu    sync.Mutex
	cache map[string]map[string]Connector
	ttl   time.Duration
}

func NewRegistry(pool *pgxpool.Pool, environment string, decryptFn func([]byte) ([]byte, error)) *Registry {
	return &Registry{
		pool: pool, environment: environment, factory: defaultFactory, decryptFn: decryptFn,
		cache: map[string]map[string]Connector{}, ttl: 60 * time.Second,
	}
}

// Get returns the connector for a merchant+connector, falling back to mock. Never returns nil.
func (r *Registry) Get(ctx context.Context, merchantID, connectorID string) Connector {
	conns, err := r.ForMerchant(ctx, merchantID)
	if err != nil {
		return NewMock()
	}
	if c, ok := conns[connectorID]; ok {
		return c
	}
	if c, ok := conns["mock"]; ok {
		return c
	}
	return NewMock()
}

// ForMerchant returns the merchant's connector map, cached briefly to avoid a DB hit per
// payment. Merchant-specific config overrides global. Always includes a mock fallback.
func (r *Registry) ForMerchant(ctx context.Context, merchantID string) (map[string]Connector, error) {
	key := merchantID
	r.mu.Lock()
	if c, ok := r.cache[key]; ok {
		r.mu.Unlock()
		return c, nil
	}
	r.mu.Unlock()

	built, err := r.build(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	r.cache[key] = built
	r.mu.Unlock()

	// Simple expiry: schedule a one-shot clear for this merchant.
	time.AfterFunc(r.ttl, func() {
		r.mu.Lock()
		delete(r.cache, key)
		r.mu.Unlock()
	})
	return built, nil
}

// build performs the DB read and connector construction (uncached).
func (r *Registry) build(ctx context.Context, merchantID string) (map[string]Connector, error) {
	out := map[string]Connector{"mock": NewMock()}

	rows, err := r.pool.Query(ctx, `
		SELECT connector_id, COALESCE(merchant_id,''), config, enabled
		FROM connector_configs
		WHERE environment=$1 AND (merchant_id IS NULL OR merchant_id=$2)
		ORDER BY merchant_id NULLS FIRST`, r.environment, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, merchant string
		var cfgJSON []byte
		var enabled bool
		if err := rows.Scan(&id, &merchant, &cfgJSON, &enabled); err != nil {
			return nil, err
		}
		if !enabled {
			continue
		}
		decrypted, err := r.decryptFn(cfgJSON)
		if err != nil {
			continue
		}
		var cfg Config
		if err := json.Unmarshal(decrypted, &cfg); err != nil {
			continue
		}
		conn, err := r.factory(id, cfg)
		if err != nil {
			continue
		}
		out[id] = conn
	}
	return out, rows.Err()
}

// Build returns a connector_id -> Connector map for a merchant, merging enabled global and
// merchant-specific configs (merchant-specific wins). Falls back to the mock connector so
// the API is always functional even before any rail is configured.
func (r *Registry) Build(ctx context.Context, merchantID string) (map[string]Connector, error) {
	out := map[string]Connector{"mock": NewMock()}

	rows, err := r.pool.Query(ctx, `
		SELECT connector_id, COALESCE(merchant_id,''), config, enabled
		FROM connector_configs
		WHERE environment=$1 AND (merchant_id IS NULL OR merchant_id=$2)
		ORDER BY merchant_id NULLS FIRST`, r.environment, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Order guarantees globals (merchant_id NULL) come first, so a merchant-specific row
	// for the same connector id naturally overrides the global one as we iterate.
	for rows.Next() {
		var id, merchant string
		var cfgJSON []byte
		var enabled bool
		if err := rows.Scan(&id, &merchant, &cfgJSON, &enabled); err != nil {
			return nil, err
		}
		if !enabled {
			continue
		}
		decrypted, err := r.decryptFn(cfgJSON)
		if err != nil {
			continue // keep other connectors working
		}
		var cfg Config
		if err := json.Unmarshal(decrypted, &cfg); err != nil {
			continue
		}
		conn, err := r.factory(id, cfg)
		if err != nil {
			continue
		}
		out[id] = conn
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
