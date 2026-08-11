// Package developer exposes a real, DB-backed developer portal API: API-key
// management (list / create / revoke). Webhooks live under the webhook module.
// Keys created here are immediately usable because they follow the auth
// middleware's scheme: prefix = first 12 chars of the token, secret_hash =
// sha256(full token) stored at rest.
package developer

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler { return &Handler{pool: pool} }

func (h *Handler) Routes(r chi.Router) {
	r.Get("/api_keys", h.ListKeys)
	r.Post("/api_keys", h.CreateKey)
	r.Post("/api_keys/{id}/revoke", h.RevokeKey)
}

// KeyRow is the safe (non-secret) representation of an API key.
type KeyRow struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	KeyType   string   `json:"key_type"`
	KeyPrefix string   `json:"key_prefix"`
	Env       string   `json:"environment"`
	Status    string   `json:"status"`
	Scopes    []string `json:"scopes"`
	LastUsed  *string  `json:"last_used_at"`
	CreatedAt string   `json:"created_at"`
}

func (h *Handler) ListKeys(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	rows, err := h.pool.Query(r.Context(), `SELECT id, name, key_type, key_prefix, environment, status,
		scopes::text, COALESCE(last_used_at::text,''), created_at::text FROM api_keys WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	defer rows.Close()
	var list []KeyRow
	for rows.Next() {
		var k KeyRow
		var scopes, lastUsed string
		if err := rows.Scan(&k.ID, &k.Name, &k.KeyType, &k.KeyPrefix, &k.Env, &k.Status,
			&scopes, &lastUsed, &k.CreatedAt); err != nil {
			pkghttp.WriteError(w, r, err)
			return
		}
		_ = json.Unmarshal([]byte(scopes), &k.Scopes)
		if lastUsed != "" {
			k.LastUsed = &lastUsed
		}
		list = append(list, k)
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

// createKey generates a random secret, stores sha256 at rest, and returns the
// plaintext secret exactly once (the dashboard shows it then never again).
func createKey(ctx context.Context, pool *pgxpool.Pool, merchantID string, name, env, scopes string) (string, error) {
	buf := make([]byte, 18)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	raw := hex.EncodeToString(buf) // 36 hex chars
	token := "sk_" + env + "_" + raw
	prefix := token[:12]
	sum := sha256.Sum256([]byte(token))
	hash := hex.EncodeToString(sum[:])
	id := id.New("key")

	if scopes == "" {
		scopes = `["payments:read","payments:write","refunds:write","payouts:write"]`
	}
	_, err := pool.Exec(ctx, `INSERT INTO api_keys (id, merchant_id, name, key_type, key_prefix, secret_hash, environment, status, scopes)
		VALUES ($1,$2,$3,'secret',$4,$5,$6,'active',$7::jsonb)`,
		id, merchantID, name, prefix, hash, env, scopes)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (h *Handler) CreateKey(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		Name   string `json:"name"`
		Env    string `json:"environment"`
		Scopes string `json:"scopes"` // optional JSON array string
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if req.Name == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "name required")
		return
	}
	if req.Env == "" {
		req.Env = "test"
	}
	if req.Env != "test" && req.Env != "live" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "environment must be test or live")
		return
	}
	token, err := createKey(r.Context(), h.pool, merchantID, req.Name, req.Env, req.Scopes)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	// Secret returned exactly once.
	pkghttp.WriteJSON(w, r, 201, map[string]string{"secret": token, "message": "secret shown once — store it now"})
}

func (h *Handler) RevokeKey(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	keyID := chi.URLParam(r, "id")
	res, err := h.pool.Exec(r.Context(), `UPDATE api_keys SET status='revoked', revoked_at=now() WHERE id=$1 AND merchant_id=$2`, keyID, merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	if res.RowsAffected() == 0 {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "api key not found")
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"id": keyID, "status": "revoked"})
}
