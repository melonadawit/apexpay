package middleware

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"strings"

	appErrors "apexpay/internal/platform/errors"
	pkghttp "apexpay/internal/platform/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Private key types prevent collisions with values set by other middleware.
type contextKey string

const (
	CtxMerchantID contextKey = "merchant_id"
	CtxUserID     contextKey = "user_id"
	CtxAPIKeyID   contextKey = "api_key_id"
	CtxRole       contextKey = "role"
)

// MerchantIDFromContext is the only supported way for handlers to read the authenticated tenant.
func MerchantIDFromContext(ctx context.Context) (string, bool) {
	merchantID, ok := ctx.Value(CtxMerchantID).(string)
	return merchantID, ok && merchantID != ""
}
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(CtxUserID).(string)
	return userID, ok && userID != ""
}
func APIKeyIDFromContext(ctx context.Context) (string, bool) {
	keyID, ok := ctx.Value(CtxAPIKeyID).(string)
	return keyID, ok && keyID != ""
}

type AuthMiddleware struct{ pool *pgxpool.Pool }
func NewAuth(pool *pgxpool.Pool) *AuthMiddleware { return &AuthMiddleware{pool: pool} }

// APIKeyAuth verifies the whole bearer secret, not merely its public prefix. secret_hash
// must contain a lowercase SHA-256 hex digest of the complete generated API key. New key
// issuance must generate high-entropy values and persist only this digest.
func (a *AuthMiddleware) APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || len(parts[1]) < 16 {
			pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, string(appErrors.CodeUnauthorized), "Bearer API key required")
			return
		}
		token := parts[1]
		prefixLen := 12
		if len(token) < prefixLen { // defensive; retained separately to make prefix policy explicit
			pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, string(appErrors.CodeUnauthorized), "invalid API key")
			return
		}

		var merchantID, keyID, status, storedHash string
		err := a.pool.QueryRow(r.Context(), `SELECT merchant_id, id, status, COALESCE(secret_hash, '') FROM api_keys WHERE key_prefix=$1`, token[:prefixLen]).Scan(&merchantID, &keyID, &status, &storedHash)
		if err != nil || status != "active" || storedHash == "" {
			pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, string(appErrors.CodeUnauthorized), "invalid API key")
			return
		}
		digest := sha256.Sum256([]byte(token))
		if subtle.ConstantTimeCompare([]byte(strings.ToLower(storedHash)), []byte(fmtSHA256(digest))) != 1 {
			pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, string(appErrors.CodeUnauthorized), "invalid API key")
			return
		}

		ctx := context.WithValue(r.Context(), CtxMerchantID, merchantID)
		ctx = context.WithValue(ctx, CtxAPIKeyID, keyID)
		// Compatibility bridge while all handlers are migrated to the typed accessors.
		// Do not add new reads using string keys.
		ctx = context.WithValue(ctx, "merchant_id", merchantID)
		go func() { _, _ = a.pool.Exec(context.Background(), `UPDATE api_keys SET last_used_at=now() WHERE id=$1`, keyID) }()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func fmtSHA256(sum [sha256.Size]byte) string {
	const hexdigits = "0123456789abcdef"
	buf := make([]byte, sha256.Size*2)
	for i, b := range sum { buf[i*2], buf[i*2+1] = hexdigits[b>>4], hexdigits[b&0x0f] }
	return string(buf)
}

func RBAC(allowedRoles ...string) func(http.Handler) http.Handler {
	allowedMap := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles { allowedMap[r] = true }
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(CtxRole).(string)
			if !ok || role == "" { pkghttp.WriteErrorWithBody(w, r, http.StatusForbidden, string(appErrors.CodeForbidden), "role required"); return }
			if allowedMap[role] { next.ServeHTTP(w, r); return }
			pkghttp.WriteErrorWithBody(w, r, http.StatusForbidden, string(appErrors.CodeForbidden), "role not allowed")
		})
	}
}
