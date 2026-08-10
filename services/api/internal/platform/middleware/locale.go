package middleware

import (
	"context"
	"net/http"

	"apexpay/internal/i18n"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Locale middleware resolves the caller's language preference into the request context.
// Preference order:
//  1. Explicit X-Lang header (client-driven, e.g. portals) — validated.
//  2. The authenticated user's saved language_preference (per-user).
//  3. Accept-Language header (best-effort).
//  4. Default (English).
type LocaleMiddleware struct {
	pool *pgxpool.Pool
}

func NewLocale(pool *pgxpool.Pool) *LocaleMiddleware { return &LocaleMiddleware{pool: pool} }

const CtxLocale contextKey = "locale"

// LocaleFromContext returns the resolved locale, defaulting to English.
func LocaleFromContext(ctx context.Context) i18n.Locale {
	if v, ok := ctx.Value(CtxLocale).(i18n.Locale); ok && v != "" {
		return v
	}
	return i18n.DefaultLocale
}

// WithLocale attaches a locale to the context.
func WithLocale(ctx context.Context, l i18n.Locale) context.Context {
	return context.WithValue(ctx, CtxLocale, l)
}

// Handler wraps the next handler with locale resolution.
func (m *LocaleMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := m.resolve(r)
		next.ServeHTTP(w, r.WithContext(WithLocale(r.Context(), locale)))
	})
}

func (m *LocaleMiddleware) resolve(r *http.Request) i18n.Locale {
	// 1. Explicit header.
	if v := r.Header.Get("X-Lang"); v != "" && i18n.IsValid(v) {
		return i18n.Normalize(v)
	}
	// 2. Per-user preference (only meaningful when a session user is present).
	if userID := UserID(r.Context()); userID != "" && m.pool != nil {
		var lang string
		if err := m.pool.QueryRow(r.Context(), `SELECT COALESCE(language_preference,'en') FROM users WHERE id=$1`, userID).Scan(&lang); err == nil && i18n.IsValid(lang) {
			return i18n.Normalize(lang)
		}
	}
	// 3. Accept-Language best-effort.
	if v := r.Header.Get("Accept-Language"); v != "" {
		tag := v
		if i := indexByte(v, ','); i > 0 {
			tag = v[:i]
		}
		if i18n.IsValid(tag) {
			return i18n.Normalize(tag)
		}
	}
	// 4. Default.
	return i18n.DefaultLocale
}

func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
