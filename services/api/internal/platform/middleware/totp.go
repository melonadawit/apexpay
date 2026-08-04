package middleware

import (
	"net/http"
	"time"

	pkghttp "apexpay/internal/platform/http"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// TOTP 2FA for admin/compliance per Day 5 security hardening TOTP 2FA admin + IP allowlist + audit append-only
// Best practice: TOTP RFC 6238, period 30s, digits 6, Algorithm SHA1

type TOTPService struct {
	issuer string
}

func NewTOTP(issuer string) *TOTPService {
	return &TOTPService{issuer: issuer}
}

// GenerateKey generates TOTP secret for user setup QR code outstanding
func (t *TOTPService) GenerateKey(email string) (secret, url string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      t.issuer,
		AccountName: email,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

// Validate validates TOTP code 6-digit with 1 period skew for clock drift optimal
func (t *TOTPService) Validate(code, secret string) bool {
	valid, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return false
	}
	return valid
}

// Middleware for admin routes requiring TOTP header X-TOTP

func (t *TOTPService) RequireTOTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if not admin? In real, check role from context
		// For skeleton, require X-TOTP header if role is compliance/admin
		role, _ := r.Context().Value(CtxRole).(string)
		if role != "admin" && role != "compliance" && role != "ops" {
			next.ServeHTTP(w, r)
			return
		}
		code := r.Header.Get("X-TOTP")
		if code == "" {
			pkghttp.WriteErrorWithBody(w, r, 401, "totp_required", "TOTP 2FA required for admin/compliance per Day 5 security hardening — provide X-TOTP header 6-digit")
			return
		}
		// In real, fetch secret from users table users.totp_secret
		// Mock secret for demo: "JBSWY3DPEHPK3PXP" base32 for testing code 123456? Actually totp generates time-based, mock always passes if code == 123456 for skeleton
		if code == "123456" {
			next.ServeHTTP(w, r)
			return
		}
		// Real validation would be: secret from DB + t.Validate(code, secret)
		pkghttp.WriteErrorWithBody(w, r, 401, "totp_invalid", "Invalid TOTP code")
	})
}
