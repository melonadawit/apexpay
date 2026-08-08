package connector

import (
	"context"
)

// railPaths describes the REST endpoints a rail exposes. Rails map their API into this
// contract; only the paths and (optionally) payload field names differ.
type railPaths struct {
	initialize string
	verify     string
	refund     string
	health     string
}

// baseRail is a DRY HTTP connector. A specific rail (Telebirr, CBE Birr, Amole, EthSwitch,
// card acquirer) is just a baseRail with its own id, paths, and optional health endpoint.
// Each satisfies the Connector interface.
type baseRail struct {
	id     string
	paths  railPaths
	cfg    Config
	client *railClient
}

func newBaseRail(id string, p railPaths, cfg Config) *baseRail {
	return &baseRail{id: id, paths: p, cfg: cfg, client: newRailClient(cfg)}
}

func (r *baseRail) ID() string { return r.id }

func (r *baseRail) Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error) {
	payload := map[string]string{
		"merchant_id": r.cfg.MerchantID,
		"amount":      req.Amount,
		"currency":    req.Currency,
		"tx_ref":      req.TxRef,
		"return_url":  req.ReturnURL,
	}
	res, err := r.client.do(ctx, "POST", r.paths.initialize, payload)
	if err != nil {
		return InitializeResponse{}, err
	}
	return InitializeResponse{
		ConnectorRef: res.ConnectorRef,
		CheckoutURL:  res.CheckoutURL,
		Status:       res.Status,
	}, nil
}

func (r *baseRail) Verify(ctx context.Context, req VerifyRequest) (VerifyResponse, error) {
	payload := map[string]string{"connector_ref": req.ConnectorRef, "tx_ref": req.TxRef}
	res, err := r.client.do(ctx, "POST", r.paths.verify, payload)
	if err != nil {
		return VerifyResponse{}, err
	}
	return VerifyResponse{
		Status:      res.Status,
		Amount:      res.Amount,
		FailureCode: res.FailureCode,
	}, nil
}

func (r *baseRail) Refund(ctx context.Context, req RefundRequest) (RefundResponse, error) {
	payload := map[string]string{"connector_ref": req.ConnectorRef, "amount": req.Amount, "refund_ref": req.RefundRef}
	res, err := r.client.do(ctx, "POST", r.paths.refund, payload)
	if err != nil {
		return RefundResponse{}, err
	}
	return RefundResponse{Status: res.Status, ConnectorRef: res.ConnectorRef}, nil
}

func (r *baseRail) Health(ctx context.Context) (bool, int, error) {
	if r.paths.health == "" {
		return true, 0, nil
	}
	res, err := r.client.do(ctx, "GET", r.paths.health, nil)
	if err != nil {
		return false, 0, err
	}
	return res.Status == "ok", res.LatencyMS, nil
}

// ---------------------------------------------------------------------------
// Factories for each rail. In production these point at the live rail endpoints
// (config is per-merchant from connector_configs). In tests they point at the
// sandbox server. All implement the Connector interface.
// ---------------------------------------------------------------------------

// NewTelebirr builds a Telebirr adapter (Ethio Telecom mobile money).
func NewTelebirr(cfg Config) (Connector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return newBaseRail("telebirr", railPaths{
		initialize: "/v1/payments/initialize",
		verify:     "/v1/payments/verify",
		refund:     "/v1/payments/refund",
		health:     "/v1/health",
	}, cfg), nil
}

// NewCBEBirr builds a CBE Birr adapter (Commercial Bank of Ethiopia mobile money).
func NewCBEBirr(cfg Config) (Connector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return newBaseRail("cbe_birr", railPaths{
		initialize: "/api/pay/initiate",
		verify:     "/api/pay/status",
		refund:     "/api/pay/refund",
		health:     "/api/health",
	}, cfg), nil
}

// NewAmole builds an Amole adapter (Dashen Bank digital wallet).
func NewAmole(cfg Config) (Connector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return newBaseRail("amole", railPaths{
		initialize: "/merchant/checkout/init",
		verify:     "/merchant/checkout/query",
		refund:     "/merchant/checkout/refund",
		health:     "/merchant/health",
	}, cfg), nil
}

// NewEthSwitch builds an EthSwitch interoperable adapter (NBE national switch/QR).
func NewEthSwitch(cfg Config) (Connector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return newBaseRail("ethswitch", railPaths{
		initialize: "/v1/qr/transactions/init",
		verify:     "/v1/qr/transactions/status",
		refund:     "/v1/qr/transactions/refund",
		health:     "/v1/health",
	}, cfg), nil
}

// NewCardAcquirer builds a card acquirer adapter (Visa/Mastercard processor).
func NewCardAcquirer(cfg Config) (Connector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return newBaseRail("card_acquirer", railPaths{
		initialize: "/v1/charges/init",
		verify:     "/v1/charges/verify",
		refund:     "/v1/charges/refund",
		health:     "/v1/health",
	}, cfg), nil
}
