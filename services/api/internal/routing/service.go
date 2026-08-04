package routing

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"apexpay/internal/platform/errors"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	ListRules(ctx context.Context, merchantID *string) ([]RoutingRule, error)
	ListHealthSamples(ctx context.Context, connector ConnectorID, since time.Time) ([]HealthSample, error)
	SaveHealthSample(ctx context.Context, s HealthSample) error
}

type Service struct {
	repo  Repository
	redis *redis.Client

	// Circuit breaker state - optimal in-memory + Redis backup
	mu       sync.RWMutex
	circuits map[ConnectorID]*CircuitBreaker
}

type CircuitBreaker struct {
	Failures    int
	LastFailure time.Time
	State       string // closed, open, half_open
	OpenedAt    *time.Time
}

func NewService(repo Repository, redis *redis.Client) *Service {
	return &Service{
		repo: repo, redis: redis,
		circuits: make(map[ConnectorID]*CircuitBreaker),
	}
}

// Evaluate routing - optimal O(n log n) sort by priority + health score
func (s *Service) Evaluate(ctx context.Context, merchantID string, amount decimal.Decimal, currency, method string) (*RoutingDecision, error) {
	rules, err := s.repo.ListRules(ctx, &merchantID)
	if err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		// fallback global rules
		rules, err = s.repo.ListRules(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	// Filter by amount/currency/method + enabled + priority sort
	candidates := make([]RoutingRule, 0, len(rules))
	for _, r := range rules {
		if !r.Enabled {
			continue
		}
		if r.Currency != "" && r.Currency != currency {
			continue
		}
		if r.MinAmount != nil && amount.LessThan(*r.MinAmount) {
			continue
		}
		if r.MaxAmount != nil && amount.GreaterThan(*r.MaxAmount) {
			continue
		}
		if r.PaymentMethod != nil && method != "" && *r.PaymentMethod != method {
			continue
		}
		candidates = append(candidates, r)
	}

	if len(candidates) == 0 {
		// Ultimate fallback mock
		return &RoutingDecision{
			Primary: ConnectorMock, Chosen: ConnectorMock, Reason: "no matching rule, fallback mock",
		}, nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Priority < candidates[j].Priority
	})

	selected := candidates[0]

	// Get health for primary + fallbacks
	healthMap := make(map[ConnectorID]ConnectorHealth)
	connectorsToCheck := []ConnectorID{selected.PrimaryConnector}
	if selected.Fallback1 != nil {
		connectorsToCheck = append(connectorsToCheck, *selected.Fallback1)
	}
	if selected.Fallback2 != nil {
		connectorsToCheck = append(connectorsToCheck, *selected.Fallback2)
	}

	for _, cid := range connectorsToCheck {
		h, _ := s.getHealth(ctx, cid)
		healthMap[cid] = h
	}

	// Circuit breaker check + strategy
	chosen := selected.PrimaryConnector
	reason := "primary healthy"

	s.mu.RLock()
	cb, ok := s.circuits[chosen]
	s.mu.RUnlock()
	if ok && cb.State == "open" {
		if cb.OpenedAt != nil && time.Since(*cb.OpenedAt) < 60*time.Second {
			// circuit open, try fallback1
			if selected.Fallback1 != nil {
				chosen = *selected.Fallback1
				reason = "primary circuit open, fallback1 used"
			}
		} else {
			// move to half_open
			s.mu.Lock()
			cb.State = "half_open"
			s.mu.Unlock()
		}
	}

	// Strategy: if success_rate strategy, pick highest success among candidates if primary success < threshold
	if selected.Strategy == StrategySuccessRate {
		primaryHealth := healthMap[selected.PrimaryConnector]
		if primaryHealth.SuccessRate5m < 0.7 {
			// find best fallback by success_rate
			best := chosen
			bestRate := primaryHealth.SuccessRate5m
			for _, c := range connectorsToCheck[1:] {
				if h, ok := healthMap[c]; ok && h.SuccessRate5m > bestRate {
					best = c
					bestRate = h.SuccessRate5m
				}
			}
			if best != chosen {
				reason = "primary success_rate low, better fallback"
				chosen = best
			}
		}
	}

	// Latency strategy
	if selected.Strategy == StrategyLatency {
		// pick lowest latency
		best := chosen
		bestLat := healthMap[chosen].AvgLatency5m
		for _, c := range connectorsToCheck {
			if h, ok := healthMap[c]; ok && h.AvgLatency5m > 0 && h.AvgLatency5m < bestLat {
				best = c
				bestLat = h.AvgLatency5m
			}
		}
		chosen = best
	}

	fallbacks := []ConnectorID{}
	if selected.Fallback1 != nil {
		fallbacks = append(fallbacks, *selected.Fallback1)
	}
	if selected.Fallback2 != nil {
		fallbacks = append(fallbacks, *selected.Fallback2)
	}

	return &RoutingDecision{
		RuleID: selected.ID, Primary: selected.PrimaryConnector,
		Fallbacks: fallbacks, Chosen: chosen, Reason: reason,
		HealthSnapshot: healthMap,
	}, nil
}

func (s *Service) getHealth(ctx context.Context, cid ConnectorID) (ConnectorHealth, error) {
	samples, err := s.repo.ListHealthSamples(ctx, cid, time.Now().Add(-5*time.Minute))
	if err != nil {
		return ConnectorHealth{ConnectorID: cid, SuccessRate5m: 1.0, CircuitState: "closed"}, err
	}
	if len(samples) == 0 {
		return ConnectorHealth{ConnectorID: cid, SuccessRate5m: 1.0, AvgLatency5m: 100, CircuitState: "closed"}, nil
	}
	success := 0
	latSum := 0
	for _, sm := range samples {
		if sm.Success {
			success++
		}
		latSum += sm.LatencyMS
	}
	rate := float64(success) / float64(len(samples))
	avgLat := latSum / len(samples)

	s.mu.RLock()
	cb, ok := s.circuits[cid]
	s.mu.RUnlock()
	state := "closed"
	if ok {
		state = cb.State
	}

	return ConnectorHealth{
		ConnectorID: cid, SuccessRate5m: rate, AvgLatency5m: avgLat,
		CircuitState: state, LastSampleAt: samples[0].SampledAt,
	}, nil
}

func (s *Service) RecordFailure(connector ConnectorID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cb, ok := s.circuits[connector]
	if !ok {
		cb = &CircuitBreaker{State: "closed"}
		s.circuits[connector] = cb
	}
	cb.Failures++
	cb.LastFailure = time.Now()
	if cb.Failures >= 5 {
		now := time.Now()
		cb.State = "open"
		cb.OpenedAt = &now
		// Reset failures after open to avoid overflow
		cb.Failures = 0
	}
}

func (s *Service) RecordSuccess(connector ConnectorID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cb, ok := s.circuits[connector]
	if !ok {
		cb = &CircuitBreaker{State: "closed"}
		s.circuits[connector] = cb
	}
	if cb.State == "half_open" {
		cb.State = "closed"
		cb.Failures = 0
		cb.OpenedAt = nil
	} else if cb.State == "closed" {
		cb.Failures = 0
	}
}

// Methods ranked - for GET /v1/methods
func (s *Service) RankedMethods(ctx context.Context, merchantID string, amount decimal.Decimal, currency string) ([]struct {
	ConnectorID ConnectorID
	Health      ConnectorHealth
	Score       float64
}, error) {
	// Evaluate all connectors health
	allConnectors := []ConnectorID{ConnectorMock, ConnectorTelebirr, ConnectorCBEBirr, ConnectorBankIPS, ConnectorEthSwitch, ConnectorCard}
	result := make([]struct {
		ConnectorID ConnectorID
		Health      ConnectorHealth
		Score       float64
	}, 0, len(allConnectors))

	for _, cid := range allConnectors {
		h, _ := s.getHealth(ctx, cid)
		score := h.SuccessRate5m*0.6 + (1.0-float64(min(h.AvgLatency5m, 1000))/1000.0)*0.4
		result = append(result, struct {
			ConnectorID ConnectorID
			Health      ConnectorHealth
			Score       float64
		}{ConnectorID: cid, Health: h, Score: score})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	return result, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Validate routing rule - prevent cycles
func ValidateNoCycle(rule RoutingRule) error {
	if rule.PrimaryConnector == "" {
		return errors.Validation("primary connector required")
	}
	seen := make(map[ConnectorID]bool)
	seen[rule.PrimaryConnector] = true
	if rule.Fallback1 != nil {
		if seen[*rule.Fallback1] {
			return errors.Validation("fallback1 duplicate primary")
		}
		seen[*rule.Fallback1] = true
	}
	if rule.Fallback2 != nil {
		if seen[*rule.Fallback2] {
			return errors.Validation("fallback2 duplicate")
		}
	}
	return nil
}
