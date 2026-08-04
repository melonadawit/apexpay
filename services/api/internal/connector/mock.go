package connector

import (
	"context"
	"fmt"
	"time"
)

type MockConnector struct{}

func NewMock() *MockConnector { return &MockConnector{} }

func (m *MockConnector) ID() string { return "mock" }

func (m *MockConnector) Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error) {
	time.Sleep(50 * time.Millisecond)
	return InitializeResponse{
		ConnectorRef: "mock_ref_" + req.TxRef,
		CheckoutURL:  fmt.Sprintf("https://checkout.apexpay.et/mock/%s", req.TxRef),
		Status: "pending",
	}, nil
}

func (m *MockConnector) Verify(ctx context.Context, req VerifyRequest) (VerifyResponse, error) {
	time.Sleep(30 * time.Millisecond)
	return VerifyResponse{Status: "succeeded", Amount: "100.00"}, nil
}

func (m *MockConnector) Refund(ctx context.Context, req RefundRequest) (RefundResponse, error) {
	return RefundResponse{Status: "succeeded", ConnectorRef: "mock_refund_" + req.RefundRef}, nil
}

func (m *MockConnector) Health(ctx context.Context) (bool, int, error) { return true, 45, nil }
