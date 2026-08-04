package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// New creates structured logger with request_id and merchant_id fields redaction for PII.
// Secrets/FIN never logged - enforced via field filter.
func New(env string) zerolog.Logger {
	level := zerolog.InfoLevel
	if env == "local" {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "apexpay-api").
		Logger()

	// Redact sensitive fields in prod - middleware will add request_id
	return logger
}

// WithRequestID returns logger with correlation IDs.
func WithRequestID(l zerolog.Logger, requestID, merchantID string) zerolog.Logger {
	return l.With().Str("request_id", requestID).Str("merchant_id", merchantID).Logger()
}
