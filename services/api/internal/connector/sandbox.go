package connector

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// Sandbox simulates the external rail APIs for local/test use. It exposes the same
// endpoint shapes the adapters call, so the whole gateway (adapters + registry +
// payment service) can be exercised end-to-end without real bank credentials.
// In production the adapters point at the real rail endpoints instead.

// NewSandboxHandler returns an http.Handler that answers rail-style requests.
// It accepts any of the five rails' path patterns.
func NewSandboxHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		var status, connectorRef string
		checkout := false

		switch {
		case strings.Contains(path, "initialize") || strings.Contains(path, "initiate") || strings.Contains(path, "init"):
			status, connectorRef, checkout = "pending", "sandbox_ref_"+refSuffix(r), true
		case strings.Contains(path, "verify") || strings.Contains(path, "status") || strings.Contains(path, "query"):
			status, connectorRef = "succeeded", "sandbox_ref_"+refSuffix(r)
		case strings.Contains(path, "refund"):
			status, connectorRef = "succeeded", "sandbox_refund_"+refSuffix(r)
		case strings.Contains(path, "health"):
			writeRail(w, railResponse{Status: "ok", LatencyMS: 40})
			return
		default:
			http.NotFound(w, r)
			return
		}

		_ = time.Now
		resp := railResponse{
			Status:       status,
			ConnectorRef: connectorRef,
			Amount:       r.FormValue("amount"),
			CheckoutURL:  "",
		}
		if checkout {
			resp.CheckoutURL = "https://sandbox.apexpay.et/checkout/" + connectorRef
		}
		writeRail(w, resp)
	})
	return mux
}

func refSuffix(r *http.Request) string {
	// Derive a stable suffix from the tx_ref in the body if present.
	var b struct {
		TxRef        string `json:"tx_ref"`
		ConnectorRef string `json:"connector_ref"`
		RefundRef    string `json:"refund_ref"`
	}
	_ = json.NewDecoder(r.Body).Decode(&b)
	for _, s := range []string{b.TxRef, b.ConnectorRef, b.RefundRef} {
		if s != "" {
			return s
		}
	}
	return "unknown"
}

func writeRail(w http.ResponseWriter, resp railResponse) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
