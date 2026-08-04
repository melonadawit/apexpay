package middleware

import (
	"net"
	"net/http"
	"strings"

	pkghttp "apexpay/internal/platform/http"
)

// IPAllowlist per Day 5 IP allowlist admin + audit append-only + WAF
// Best practice: CIDR matching O(n) where n small (office IPs), optimal trie would be for large n but office list small

type IPAllowlist struct {
	allowedNets []*net.IPNet
}

func NewIPAllowlist(cidrs []string) (*IPAllowlist, error) {
	var nets []*net.IPNet
	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			// Try single IP
			ip := net.ParseIP(cidr)
			if ip == nil {
				continue
			}
			// /32 for single IP
			_, ipnet, _ = net.ParseCIDR(cidr + "/32")
			if ipnet == nil {
				continue
			}
		}
		nets = append(nets, ipnet)
	}
	return &IPAllowlist{allowedNets: nets}, nil
}

func (a *IPAllowlist) Contains(ipStr string) bool {
	// Strip port if present
	if strings.Contains(ipStr, ":") {
		host, _, err := net.SplitHostPort(ipStr)
		if err == nil {
			ipStr = host
		}
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	for _, net := range a.allowedNets {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}

func (a *IPAllowlist) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get real IP from X-Real-IP if behind Nginx per nginx.conf proxy_set_header X-Real-IP
		ip := r.RemoteAddr
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			ip = realIP
		}
		if !a.Contains(ip) {
			pkghttp.WriteErrorWithBody(w, r, 403, "forbidden", "IP not allowed per Day 5 IP allowlist admin — office IPs only")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Default office allowlist example per nginx.conf geo $admin_allowed
func DefaultAdminAllowlist() *IPAllowlist {
	// Example: Ethio Telecom range 196.188.0.0/16 + office
	cidrs := []string{"196.188.0.0/16", "127.0.0.1/32", "10.0.0.0/8"}
	allowlist, _ := NewIPAllowlist(cidrs)
	return allowlist
}
