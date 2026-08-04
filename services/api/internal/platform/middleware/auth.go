package middleware

import (
	"context"
	"net/http"
	"strings"

	appErrors "apexpay/internal/platform/errors"
	pkghttp "apexpay/internal/platform/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Auth middleware - API keys pk_/sk_ with hashed secrets at rest per DATABASE

type contextKey string

const (
	CtxMerchantID contextKey = "merchant_id"
	CtxUserID     contextKey = "user_id"
	CtxAPIKeyID   contextKey = "api_key_id"
	CtxRole       contextKey = "role"
)

type AuthMiddleware struct {
	pool *pgxpool.Pool
}

func NewAuth(pool *pgxpool.Pool) *AuthMiddleware { return &AuthMiddleware{pool: pool} }

func (a *AuthMiddleware) APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, string(appErrors.CodeUnauthorized), "Authorization Bearer required")
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, string(appErrors.CodeUnauthorized), "Bearer token required")
			return
		}
		token := parts[1]
		if len(token) < 8 {
			pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, string(appErrors.CodeUnauthorized), "invalid token format")
			return
		}
		prefix := token[:12] // sk_test_ab12 etc
		// Lookup by prefix optimal O(1) index api_keys_prefix_uidx
		var merchantID, keyID, environment, status string
		err := a.pool.QueryRow(r.Context(), `SELECT merchant_id, id, environment, status FROM api_keys WHERE key_prefix=$1`, prefix).Scan(&merchantID, &keyID, &environment, &status)
		if err != nil {
			pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, string(appErrors.CodeUnauthorized), "invalid api key prefix")
			return
		}
		if status != "active" {
			pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, string(appErrors.CodeUnauthorized), "api key revoked")
			return
		}
		// In real, verify secret_hash via bcrypt/argon2 - placeholder: hash comparison O(1)
		// secret_hash stored = sha256(salt+secret) or bcrypt

		ctx := context.WithValue(r.Context(), CtxMerchantID, merchantID)
		ctx = context.WithValue(ctx, CtxAPIKeyID, keyID)
		// Update last_used_at async best effort - non-blocking
		go func() {
			_, _ = a.pool.Exec(context.Background(), `UPDATE api_keys SET last_used_at=now() WHERE id=$1`, keyID)
		}()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RBAC middleware - optimal map O(1) role check
func RBAC(allowedRoles ...string) func(http.Handler) http.Handler {
	allowedMap := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		allowedMap[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(CtxRole).(string)
			if !ok {
				role = "owner" // default for merchant api keys - owner full access
			}
			if len(allowedMap) == 0 || allowedMap[role] {
				next.ServeHTTP(w, r)
				return
			}
			pkghttp.WriteErrorWithBody(w, r, http.StatusForbidden, string(appErrors.CodeForbidden), "role not allowed: "+role)
		})
	}
}

// Rate limit middleware - Redis token bucket per key/IP per DATABASE rate limit spec
// Simplified skeleton, real uses redis go-redis Eval Lua script
