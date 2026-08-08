package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// railResponse is the canonical JSON response shape that a rail (sandbox or live) returns
// for initialize/verify/refund. It is intentionally versioned so future live rails can map
// into it without changing the Connector interface.
type railResponse struct {
	Status       string `json:"status"` // succeeded | failed | pending
	ConnectorRef string `json:"connector_ref"`
	CheckoutURL  string `json:"checkout_url"`
	Amount       string `json:"amount"`
	FailureCode  string `json:"failure_code"`
	LatencyMS    int    `json:"latency_ms"`
}

// railClient wraps an http.Client bound to one connector's config and provides
// signed JSON requests to the rail's REST API.
type railClient struct {
	cfg  Config
	http *http.Client
}

func newRailClient(cfg Config) *railClient {
	return &railClient{cfg: cfg, http: cfg.newHTTPClient()}
}

// do sends a signed JSON request and decodes the rail response.
func (c *railClient) do(ctx context.Context, method, path string, payload any) (*railResponse, error) {
	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.cfg.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", c.cfg.APIKey)
	req.Header.Set("X-ApexPay-Ts", time.Now().UTC().Format(time.RFC3339))
	// HMAC-SHA256 over the payload signs the request; some rails use base64.
	if len(body) > 0 {
		req.Header.Set("X-Signature", c.cfg.Sign(body))
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rail request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read rail response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("rail http %d: %s", resp.StatusCode, truncate(raw))
	}

	var out railResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode rail response: %w", err)
	}
	return &out, nil
}

func truncate(b []byte) string {
	const max = 300
	if len(b) > max {
		return string(b[:max]) + "..."
	}
	return string(b)
}
