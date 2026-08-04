package routing

import (
	"time"

	"github.com/shopspring/decimal"
)

type ConnectorID string

const (
	ConnectorMock      ConnectorID = "mock"
	ConnectorTelebirr  ConnectorID = "telebirr"
	ConnectorCBEBirr   ConnectorID = "cbe_birr"
	ConnectorBankIPS   ConnectorID = "bank_ips"
	ConnectorEthSwitch ConnectorID = "ethswitch"
	ConnectorCard      ConnectorID = "card_acquirer"
)

type Strategy string

const (
	StrategySuccessRate Strategy = "success_rate"
	StrategyLatency     Strategy = "latency"
	StrategyCost        Strategy = "cost"
	StrategyRoundRobin  Strategy = "round_robin"
)

type RoutingRule struct {
	ID               string
	MerchantID       *string // nil = global
	Name             string
	MinAmount        *decimal.Decimal
	MaxAmount        *decimal.Decimal
	Currency         string
	PaymentMethod    *string
	PrimaryConnector ConnectorID
	Fallback1        *ConnectorID
	Fallback2        *ConnectorID
	Strategy         Strategy
	Enabled          bool
	Priority         int // lower = higher priority, optimal sorting
	CreatedAt        time.Time
}

type HealthSample struct {
	ID          string
	ConnectorID ConnectorID
	Environment string
	LatencyMS   int
	Success     bool
	ErrorCode   *string
	SampledAt   time.Time
}

// Aggregated health for decision - optimal cached structure
type ConnectorHealth struct {
	ConnectorID   ConnectorID
	SuccessRate5m float64 // 0-1
	AvgLatency5m  int
	Uptime24h     float64
	CircuitState  string // closed, open, half_open
	LastSampleAt  time.Time
}

// Decision result
type RoutingDecision struct {
	RuleID         string
	Primary        ConnectorID
	Fallbacks      []ConnectorID
	Chosen         ConnectorID
	Reason         string
	HealthSnapshot map[ConnectorID]ConnectorHealth
}
