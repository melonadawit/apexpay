package auth

// Auth domain types for the merchant dashboard session layer.

// User is a dashboard user (owners/admin/finance/etc.), PII-safe in responses.
type User struct {
	ID                 string `json:"id"`
	Email              string `json:"email"`
	Name               string `json:"name"`
	Status             string `json:"status"`
	LanguagePreference string `json:"language_preference"`
}

// MerchantContext is the active merchant + role for a session.
type MerchantContext struct {
	MerchantID string `json:"merchant_id"`
	LegalName  string `json:"legal_name"`
	Role       string `json:"role"`
}

// Session is a validated dashboard session (used internally).
type Session struct {
	ID         string
	UserID     string
	MerchantID string
	Role       string
	TokenHash  string
}

// LoginResult is returned on a successful login.
type LoginResult struct {
	Token     string          `json:"token"`
	ExpiresAt string          `json:"expires_at"`
	User      User            `json:"user"`
	Merchant  MerchantContext `json:"merchant"`
}

// SessionResult is returned by GET /auth/me.
type SessionResult struct {
	User      User              `json:"user"`
	Merchant  MerchantContext   `json:"merchant"`
	Merchants []MerchantContext `json:"merchants"`
}
