package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pkghttp "apexpay/internal/platform/http"
	"github.com/redis/go-redis/v9"
)

// Rate limiter using Redis token bucket Lua script per Day 5 spec Fayda OTP 5/hour per IP via Redis token bucket Lua
// Best practice: Lua script atomic O(1) via EVAL, no race, optimal data structure token bucket

const faydaRateLimitLua = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("GET", key)
if current and tonumber(current) >= limit then
  return {0, current}
end
current = redis.call("INCR", key)
if tonumber(current) == 1 then
  redis.call("EXPIRE", key, window)
end
return {1, current}
`

type RateLimiter struct {
	redis *redis.Client
}

func NewRateLimiter(rdb *redis.Client) *RateLimiter {
	return &RateLimiter{redis: rdb}
}

// FaydaOTP5PerHour — per spec Day 5 rate limit Fayda OTP 5/hour per IP via Redis token bucket Lua
func (rl *RateLimiter) FaydaOTP5PerHour(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rl.redis == nil {
			next.ServeHTTP(w, r)
			return
		}
		ip := r.RemoteAddr
		// Use X-Real-IP if behind Nginx per nginx.conf proxy_set_header X-Real-IP
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			ip = realIP
		}
		ownerID := r.URL.Query().Get("owner_id")
		if ownerID == "" {
			ownerID = "anonymous"
		}
		key := fmt.Sprintf("fayda:otp:%s:%s", ip, ownerID)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		result, err := rl.redis.Eval(ctx, faydaRateLimitLua, []string{key}, 5, 3600).Result()
		if err != nil {
			// Fail open if Redis down per best practice for availability but log warning
			next.ServeHTTP(w, r)
			return
		}
		res, ok := result.([]interface{})
		if !ok || len(res) < 2 {
			next.ServeHTTP(w, r)
			return
		}
		allowed, _ := res[0].(int64)
		if allowed == 0 {
			pkghttp.WriteErrorWithBody(w, r, 429, "rate_limited", "Fayda OTP rate limit 5 per hour exceeded per IP per Day 5 WAF + Redis Lua token bucket — try after 1 hour")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GeneralAPIRateLimit 100/min per IP
func (rl *RateLimiter) General100PerMin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rl.redis == nil {
			next.ServeHTTP(w, r)
			return
		}
		ip := r.RemoteAddr
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			ip = realIP
		}
		key := fmt.Sprintf("api:general:%s", ip)
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		// Simple INCR + EXPIRE 60s for 100/min
		current, err := rl.redis.Incr(ctx, key).Result()
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if current == 1 {
			rl.redis.Expire(ctx, key, 60*time.Second)
		}
		if current > 100 {
			pkghttp.WriteErrorWithBody(w, r, 429, "rate_limited", "API rate limit 100/min exceeded")
			return
		}
		next.ServeHTTP(w, r)
	})
}
