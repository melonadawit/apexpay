package connector

import (
	"context"
	"net/http/httptest"
	"testing"
)

// TestRailAdapters exercises each rail adapter end-to-end against the sandbox server.
func TestRailAdapters(t *testing.T) {
	srv := httptest.NewServer(NewSandboxHandler())
	defer srv.Close()

	cfg := Config{
		BaseURL:    srv.URL,
		APIKey:     "test-key",
		Secret:     "test-secret",
		MerchantID: "mer_test",
		TimeoutMS:  2000,
	}

	factories := map[string]func() (Connector, error){
		"telebirr":      func() (Connector, error) { return NewTelebirr(cfg) },
		"cbe_birr":      func() (Connector, error) { return NewCBEBirr(cfg) },
		"amole":         func() (Connector, error) { return NewAmole(cfg) },
		"ethswitch":     func() (Connector, error) { return NewEthSwitch(cfg) },
		"card_acquirer": func() (Connector, error) { return NewCardAcquirer(cfg) },
	}

	for id, makeConn := range factories {
		conn, err := makeConn()
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if conn.ID() != id {
			t.Fatalf("%s: id = %q", id, conn.ID())
		}

		// Initialize
		init, err := conn.Initialize(context.Background(), InitializeRequest{
			MerchantID: "mer_test", Amount: "100.00", Currency: "ETB", TxRef: "txr_" + id, ReturnURL: "https://m.et/ret",
		})
		if err != nil {
			t.Fatalf("%s initialize: %v", id, err)
		}
		if init.ConnectorRef == "" {
			t.Fatalf("%s: empty connector ref", id)
		}
		if init.Status != "pending" {
			t.Fatalf("%s: unexpected status %q", id, init.Status)
		}

		// Verify
		ver, err := conn.Verify(context.Background(), VerifyRequest{ConnectorRef: init.ConnectorRef, TxRef: "txr_" + id})
		if err != nil {
			t.Fatalf("%s verify: %v", id, err)
		}
		if ver.Status != "succeeded" {
			t.Fatalf("%s: verify status %q", id, ver.Status)
		}

		// Refund
		ref, err := conn.Refund(context.Background(), RefundRequest{ConnectorRef: init.ConnectorRef, Amount: "100.00", RefundRef: "ref_" + id})
		if err != nil {
			t.Fatalf("%s refund: %v", id, err)
		}
		if ref.Status != "succeeded" {
			t.Fatalf("%s: refund status %q", id, ref.Status)
		}

		// Health
		ok, latency, err := conn.Health(context.Background())
		if err != nil {
			t.Fatalf("%s health: %v", id, err)
		}
		if !ok {
			t.Fatalf("%s: not healthy", id)
		}
		_ = latency
		t.Logf("%s: init=%s verify=%s refund=%s", id, init.Status, ver.Status, ref.Status)
	}
}

// TestConnectorConfigValidation ensures factories reject incomplete configs.
func TestConnectorConfigValidation(t *testing.T) {
	if _, err := NewTelebirr(Config{}); err == nil {
		t.Fatal("expected error for empty config")
	}
	if _, err := NewCBEBirr(Config{BaseURL: "http://x"}); err == nil {
		t.Fatal("expected error for missing api_key/secret")
	}
}

// TestMockConnector covers the fallback connector.
func TestMockConnector(t *testing.T) {
	m := NewMock()
	if m.ID() != "mock" {
		t.Fatalf("mock id = %q", m.ID())
	}
	init, err := m.Initialize(context.Background(), InitializeRequest{Amount: "10", TxRef: "t"})
	if err != nil {
		t.Fatal(err)
	}
	if init.ConnectorRef == "" {
		t.Fatal("mock should return a ref")
	}
	ok, _, err := m.Health(context.Background())
	if err != nil || !ok {
		t.Fatalf("mock health: ok=%v err=%v", ok, err)
	}
}
