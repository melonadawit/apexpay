package fayda

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MockVerifier for local dev - returns deterministic results matching id.gov.et spec
// Use FaydaMode=mock in config.

type MockVerifier struct{}

func NewMockVerifier() *MockVerifier { return &MockVerifier{} }

func (m *MockVerifier) RequestOTP(ctx context.Context, fin, fan string) (string, error) {
	// Simulate network latency 150-300ms
	time.Sleep(200 * time.Millisecond)
	// Always succeed, return fake txID
	identifier := fin
	if identifier == "" {
		identifier = fan
	}
	return fmt.Sprintf("mock_tx_%s_%d", identifier, time.Now().Unix()), nil
}

func (m *MockVerifier) VerifyOTP(ctx context.Context, faydaTxID, otp string) (bool, float64, error) {
	time.Sleep(300 * time.Millisecond)
	// Mock OTP: 123456 always succeeds, 000000 fails
	if otp == "123456" {
		return true, 0.92, nil
	}
	if otp == "000000" {
		return false, 0.2, fmt.Errorf("invalid OTP")
	}
	// Any other 6-digit passes with 0.88 if contains even sum? Simplified always true for demo
	if len(otp) == 6 {
		return true, 0.88, nil
	}
	return false, 0.0, fmt.Errorf("invalid OTP format")
}

func (m *MockVerifier) VerifyOfflineQR(ctx context.Context, qrData string) (bool, QRData, error) {
	// FaydaEncode QR format: "FIN_LAST4|NAME|DOB|SIG"
	parts := strings.Split(qrData, "|")
	if len(parts) < 3 {
		return false, QRData{}, fmt.Errorf("invalid QR format")
	}
	return true, QRData{
		FINLast4:  parts[0],
		Name:      parts[1],
		DOB:       parts[2],
		Signature: "mock_sig_valid",
	}, nil
}

func (m *MockVerifier) VerifyOIDC(ctx context.Context, code string) (FaydaVerification, error) {
	// OIDC eSignet flow - mock returns verified
	return FaydaVerification{
		Status: StatusVerified,
	}, nil
}

// LiveVerifier placeholder for real id.gov.et integration
type LiveVerifier struct {
	PartnerCode string
	PartnerKey  string
	BaseURL     string // https://id.gov.et/api
}

func NewLiveVerifier(partnerCode, partnerKey, baseURL string) *LiveVerifier {
	return &LiveVerifier{PartnerCode: partnerCode, PartnerKey: partnerKey, BaseURL: baseURL}
}

func (l *LiveVerifier) RequestOTP(ctx context.Context, fin, fan string) (string, error) {
	// Real implementation would call:
	// POST {BaseURL}/auth/otp with PartnerAPIKey header
	// body {id: FIN/FAN, otpChannel: SMS}
	// For skeleton, we keep same mock but log warning
	// TODO: implement with resty, HMAC, PartnerAPIKey approval flow per https://id.gov.et/api
	return fmt.Sprintf("live_tx_%d", time.Now().UnixNano()), nil
}

func (l *LiveVerifier) VerifyOTP(ctx context.Context, faydaTxID, otp string) (bool, float64, error) {
	// Real: POST {BaseURL}/auth/verify
	return true, 0.95, nil
}

func (l *LiveVerifier) VerifyOfflineQR(ctx context.Context, qrData string) (bool, QRData, error) {
	// Offline QR verification uses NIDP public key to verify signature - implement with crypto
	return true, QRData{}, nil
}

func (l *LiveVerifier) VerifyOIDC(ctx context.Context, code string) (FaydaVerification, error) {
	// OIDC eSignet: exchange code for id_token, verify JWT, extract FIN
	return FaydaVerification{}, nil
}
