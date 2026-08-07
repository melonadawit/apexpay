package webhook

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"

	"apexpay/internal/id"
	platformcrypto "apexpay/internal/platform/crypto"
	mw "apexpay/internal/platform/middleware"
	pkghttp "apexpay/internal/platform/http"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ repo *PgRepository; encryptionKey []byte }

func NewHandler(repo *PgRepository, encryptionKey []byte) *Handler { return &Handler{repo: repo, encryptionKey: encryptionKey} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/endpoints", h.CreateEndpoint)
	r.Get("/endpoints", h.ListEndpoints)
	r.Get("/deliveries", h.ListDeliveries)
	r.Post("/deliveries/{id}/resend", h.Resend)
}

func (h *Handler) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	merchantID, _ := mw.MerchantID(r.Context())
	var req struct {
		URL string `json:"url"`
		Secret string `json:"secret"`
		Events      []string `json:"events"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	// SSRF block private ranges per security hardening
	if req.URL == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "url required")
		return
	}
	
	// Basic SSRF protection (simplified for demonstration)
	importURL, err := url.Parse(req.URL)
	if err != nil || importURL.Scheme != "https" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "must be valid https url")
		return
	}
	
	ips, err := net.LookupIP(importURL.Hostname())
	if err == nil {
		for _, ip := range ips {
			if ip.IsPrivate() || ip.IsLoopback() {
				pkghttp.WriteErrorWithBody(w, r, 403, "security_error", "SSRF protection: internal IPs blocked")
				return
			}
		}
	}
	if len(req.Secret) < 16 { pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "webhook secret must be at least 16 characters"); return }
	endpointID := id.New("we")
	encryptedSecret, err := platformcrypto.Encrypt(h.encryptionKey, []byte(req.Secret))
	if err != nil { pkghttp.WriteError(w, r, err); return }
	_, err = h.repo.pool.Exec(r.Context(), `INSERT INTO webhook_endpoints (id, merchant_id, url, secret_hash, secret_prefix, secret_encrypted, status, events) VALUES ($1,$2,$3,$4,$5,$6,'active',$7)`,
		endpointID, merchantID, req.URL, "hash_"+req.Secret, req.Secret[:8], encryptedSecret, req.Events)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, map[string]string{"id": endpointID, "url": req.URL, "status": "active"})
}

func (h *Handler) ListEndpoints(w http.ResponseWriter, r *http.Request) {
	merchantID, _ := mw.MerchantID(r.Context())
	rows, err := h.repo.pool.Query(r.Context(), `SELECT id, url, status FROM webhook_endpoints WHERE merchant_id=$1`, merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	defer rows.Close()
	var list []map[string]string
	for rows.Next() {
		var id, url, status string
		_ = rows.Scan(&id, &url, &status)
		list = append(list, map[string]string{"id": id, "url": url, "status": status})
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	merchantID, _ := mw.MerchantID(r.Context())
	rows, err := h.repo.pool.Query(r.Context(), `SELECT id, event_type, status, attempt_count, last_status_code FROM webhook_deliveries WHERE merchant_id=$1 ORDER BY created_at DESC LIMIT 50`, merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var id, eventType, status string
		var attempt, lastCode int
		_ = rows.Scan(&id, &eventType, &status, &attempt, &lastCode)
		list = append(list, map[string]interface{}{"id": id, "event_type": eventType, "status": status, "attempt_count": attempt, "last_status_code": lastCode})
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) Resend(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.repo.pool.Exec(r.Context(), `UPDATE webhook_deliveries SET status='pending', next_attempt_at=now(), updated_at=now() WHERE id=$1`, id)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"id": id, "status": "pending", "message": "resend queued"})
}
