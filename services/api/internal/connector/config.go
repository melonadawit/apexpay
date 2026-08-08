package connector

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// Config is the per-connector configuration decoded from the `config` jsonb column of
// connector_configs. Secrets are encrypted at rest (AES-GCM) by the platform; the app
// decrypts them before constructing a connector and never logs them.
type Config struct {
	BaseURL    string            `json:"base_url"`    // sandbox or live rail endpoint
	APIKey     string            `json:"api_key"`     // rail-provided app/API key
	Secret     string            `json:"secret"`      // shared signing secret
	MerchantID string            `json:"merchant_id"` // rail-side merchant/app id
	TimeoutMS  int               `json:"timeout_ms"`
	Extra      map[string]string `json:"extra,omitempty"`
}

// Validate ensures a config has what a connector needs to operate.
func (c Config) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("connector config: base_url required")
	}
	if c.APIKey == "" {
		return fmt.Errorf("connector config: api_key required")
	}
	if c.Secret == "" {
		return fmt.Errorf("connector config: secret required")
	}
	return nil
}

// defaultTimeout returns the configured request timeout or a sane default.
func (c Config) defaultTimeout() time.Duration {
	if c.TimeoutMS > 0 {
		return time.Duration(c.TimeoutMS) * time.Millisecond
	}
	return 10 * time.Second
}

// newHTTPClient returns an *http.Client with the connector's timeout.
func (c Config) newHTTPClient() *http.Client {
	return &http.Client{Timeout: c.defaultTimeout()}
}

// Sign computes an HMAC-SHA256 signature of the payload for the rail.
func (c Config) Sign(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(c.Secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// SignBase64 returns an HMAC-SHA256 signature base64-encoded (used by some rails).
func (c Config) SignBase64(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(c.Secret))
	mac.Write(payload)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// sha256Hex returns the SHA-256 digest of data (used for signature payloads).
func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
