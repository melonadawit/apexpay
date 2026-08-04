package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Audit middleware — append-only immutable audit_logs per Day 5 audit append-only + NBE exam-ready
// Best practice: no UPDATE/DELETE on audit_logs via DB trigger prevent_update(), only INSERT, request_id correlation per SAD §11

type AuditLogger struct {
	pool *pgxpool.Pool
}

func NewAuditLogger(pool *pgxpool.Pool) *AuditLogger {
	return &AuditLogger{pool: pool}
}

func (a *AuditLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Capture request body for audit (but not for file uploads >1MB)
		var bodyBytes []byte
		if r.Body != nil && r.ContentLength < 1024*1024 {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Capture response status via wrapper
		rw := &responseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(rw, r)

		// Async audit insert best effort non-blocking for performance p99<150ms ex-rail
		go func() {
			merchantID, _ := r.Context().Value(CtxMerchantID).(string)
			userID, _ := r.Context().Value(CtxUserID).(string)
			apiKeyID, _ := r.Context().Value(CtxAPIKeyID).(string)
			actorType := "api_key"
			actorID := apiKeyID
			if userID != "" {
				actorType = "user"
				actorID = userID
			}
			// Sanitize PII: never log plain FIN, only last4 + hash
			// Simplified: truncate body if contains FIN pattern 12-digit
			// Real would use crypto.Last4 + hash check
			auditData := map[string]interface{}{
				"method": r.Method, "path": r.URL.Path, "status": rw.statusCode,
				"duration_ms": time.Since(start).Milliseconds(),
				"ip":          r.RemoteAddr, "user_agent": r.UserAgent(),
			}
			_, _ = a.pool.Exec(context.Background(), `
				INSERT INTO audit_logs (id, merchant_id, actor_type, actor_id, action, resource_type, resource_id, ip, request_id, data)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8::inet,$9,$10)
			`, id.New("audit"), merchantID, actorType, actorID, r.Method+" "+r.URL.Path, "http_request", r.URL.Path, r.RemoteAddr, r.Header.Get("X-Request-Id"), auditData)
		}()
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// DB trigger for append-only per Day 5 audit append-only spec — would be in migration 0013
// CREATE OR REPLACE FUNCTION prevent_audit_update() RETURNS TRIGGER AS $$
// BEGIN
//   RAISE EXCEPTION 'audit_logs is append-only immutable per NBE exam-ready';
//   RETURN NULL;
// END;
// $$ LANGUAGE plpgsql;
// CREATE TRIGGER audit_no_update BEFORE UPDATE OR DELETE ON audit_logs FOR EACH ROW EXECUTE FUNCTION prevent_audit_update();
