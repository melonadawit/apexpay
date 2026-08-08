package twofa

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// TwoFA implements real time-based one-time-password challenges (TOTP, RFC 6238) using the
// pquerna/otp library — a production-ready replacement for the local-only demo OTP in the
// payment flow. Supports both enrolling a secret (for authenticator apps) and verifying a
// user-provided code.

// Provider is the interface for 2FA challenge providers.
type Provider interface {
	// GenerateSecret issues a new TOTP secret (base32) for enrollment.
	GenerateSecret(account string) (secret, otpauthURL string, err error)
	// VerifyCode checks a 6-digit code against a secret with a small time-step window.
	VerifyCode(secret, code string) bool
}

type TOTPProvider struct{}

func New() *TOTPProvider { return &TOTPProvider{} }

// GenerateSecret creates a random TOTP secret and its otpauth URL for authenticator apps.
func (p *TOTPProvider) GenerateSecret(account string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ApexPay",
		AccountName: account,
		Period:      30,
		SecretSize:  20,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

// VerifyCode verifies a code within the default window (±1 step).
func (p *TOTPProvider) VerifyCode(secret, code string) bool {
	if code == "" {
		return false
	}
	valid, err := totp.ValidateCustom(code, secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	return err == nil && valid
}

// RandomSecret returns a raw random base32 secret (used when no enrollment is required).
func RandomSecret() string {
	b := make([]byte, 20)
	_, _ = rand.Read(b)
	return strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")
}
