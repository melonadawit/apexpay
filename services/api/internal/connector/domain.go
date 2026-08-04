package connector

import "context"

type InitializeRequest struct {
	MerchantID string
	Amount     string
	Currency   string
	TxRef      string
	ReturnURL  string
	Meta       map[string]any
}

type InitializeResponse struct {
	ConnectorRef string
	CheckoutURL  string
	Status       string
}

type VerifyRequest struct {
	ConnectorRef string
	TxRef        string
}

type VerifyResponse struct {
	Status      string // succeeded, failed, pending
	Amount      string
	FailureCode string
}

type RefundRequest struct {
	ConnectorRef string
	Amount       string
	RefundRef    string
}

type RefundResponse struct {
	Status       string
	ConnectorRef string
}

// Connector interface - optimal Strategy pattern per SAD ADR-004
type Connector interface {
	ID() string
	Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error)
	Verify(ctx context.Context, req VerifyRequest) (VerifyResponse, error)
	Refund(ctx context.Context, req RefundRequest) (RefundResponse, error)
	Health(ctx context.Context) (bool, int, error) // ok, latency_ms, err
}
