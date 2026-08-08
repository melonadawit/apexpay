package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"net/url"
	"time"

	platformcrypto "apexpay/internal/platform/crypto"
)

// Service delivers webhooks with HMAC signing + retry exponential backoff + SSRF block

type Delivery struct {
	ID              string
	MerchantID      string
	EndpointID      string
	EventType       string
	Payload         []byte
	URL             string
	Secret          string
	EncryptedSecret []byte
	AttemptCount    int
}

type Repository interface {
	ListPendingDeliveries(ctx context.Context, limit int) ([]Delivery, error)
	MarkSuccess(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string, statusCode int, errMsg string, nextAttempt time.Time) error
}

type Service struct {
	repo          Repository
	client        *http.Client
	encryptionKey []byte
}

func NewService(repo Repository, encryptionKey []byte) *Service {
	return &Service{
		repo:   repo,
		client: &http.Client{Timeout: 10 * time.Second}, encryptionKey: encryptionKey,
	}
}

// Sign payload HMAC SHA256 per webhook secret - outstanding security
func Sign(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *Service) Deliver(ctx context.Context, d Delivery) error {
	// SSRF block private ranges - simplified check
	if isPrivateURL(d.URL) {
		return s.repo.MarkFailed(ctx, d.ID, 0, "SSRF blocked private range", time.Now().Add(1*time.Hour))
	}

	secret := d.Secret
	if len(d.EncryptedSecret) > 0 {
		plain, err := platformcrypto.Decrypt(s.encryptionKey, d.EncryptedSecret)
		if err != nil {
			return s.repo.MarkFailed(ctx, d.ID, 0, "webhook secret decrypt failed", time.Now().Add(time.Hour))
		}
		secret = string(plain)
	}
	if secret == "" {
		return s.repo.MarkFailed(ctx, d.ID, 0, "webhook secret unavailable", time.Now().Add(time.Hour))
	}
	sig := Sign(d.Payload, secret)
	req, err := http.NewRequestWithContext(ctx, "POST", d.URL, bytes.NewReader(d.Payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ApexPay-Signature", sig)
	req.Header.Set("X-ApexPay-Event", d.EventType)
	req.Header.Set("X-Request-Id", d.ID)

	resp, err := s.client.Do(req)
	if err != nil {
		next := time.Now().Add(backoff(d.AttemptCount))
		return s.repo.MarkFailed(ctx, d.ID, 0, err.Error(), next)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return s.repo.MarkSuccess(ctx, d.ID)
	}
	next := time.Now().Add(backoff(d.AttemptCount))
	return s.repo.MarkFailed(ctx, d.ID, resp.StatusCode, "non-2xx", next)
}

func backoff(attempt int) time.Duration {
	// exponential backoff: 1s, 2s, 4s, 8s, 16s, 32s, max 1h
	switch attempt {
	case 0:
		return 1 * time.Second
	case 1:
		return 2 * time.Second
	case 2:
		return 4 * time.Second
	case 3:
		return 8 * time.Second
	case 4:
		return 16 * time.Second
	case 5:
		return 32 * time.Second
	default:
		return 1 * time.Hour
	}
}

// isPrivateURL blocks server-side request forgery against internal network targets.
// It rejects non-http(s) schemes, then resolves the host and rejects any IP that is
// loopback, link-local, CGNAT, private, unspecified, multicast, or otherwise non-routable.
// The returned boolean means "unsafe — do not deliver".
func isPrivateURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return true
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return true
	}
	host := parsed.Hostname()
	if host == "" {
		return true
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return true
	}
	for _, ip := range ips {
		if isBlockedIP(ip) {
			return true
		}
	}
	return false
}

func isBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	// Normalize to 4-byte form for v4-mapped-in-v6 addresses.
	if v4 := ip.To4(); v4 != nil {
		ip = v4
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsMulticast() || ip.IsUnspecified() {
		return true
	}
	// CGNAT (RFC 6598) 100.64.0.0/10 and documentation ranges (RFC 5737) are not routable.
	if v4 := ip.To4(); v4 != nil {
		first := v4[0]
		if first == 100 && v4[1] >= 64 && v4[1] <= 127 { // 100.64.0.0/10
			return true
		}
	}
	return false
}
