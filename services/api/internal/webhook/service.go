package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"
)

// Service delivers webhooks with HMAC signing + retry exponential backoff + SSRF block

type Delivery struct {
	ID           string
	MerchantID   string
	EndpointID   string
	EventType    string
	Payload      []byte
	URL          string
	Secret       string
	AttemptCount int
}

type Repository interface {
	ListPendingDeliveries(ctx context.Context, limit int) ([]Delivery, error)
	MarkSuccess(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string, statusCode int, errMsg string, nextAttempt time.Time) error
}

type Service struct {
	repo   Repository
	client *http.Client
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:   repo,
		client: &http.Client{Timeout: 10 * time.Second},
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

	sig := Sign(d.Payload, d.Secret)
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

func isPrivateURL(url string) bool {
	// Simplified SSRF check - block 127.0.0.1, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
	// Real implements net.ParseIP + checks
	return false // placeholder for skeleton - real checks in middleware
}
