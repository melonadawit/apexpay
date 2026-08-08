package bankverification

import (
	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strings"
	"time"
)

type Handler struct{ pool *pgxpool.Pool }

func NewHandler(pool *pgxpool.Pool) *Handler { return &Handler{pool: pool} }
func (h *Handler) Routes(r chi.Router) {
	r.Post("/bank_accounts/verify", h.Create)
	r.Get("/bank_accounts/verifications/{id}", h.Get)
}
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	merchant := mw.MerchantID(r.Context())
	var q struct {
		BankCode      string `json:"bank_code"`
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		Method        string `json:"method"`
		ConnectorID   string `json:"connector_id"`
	}
	if json.NewDecoder(r.Body).Decode(&q) != nil || q.BankCode == "" || len(q.AccountNumber) < 4 {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "bank code and valid account number required")
		return
	}
	if q.Method == "" {
		q.Method = "manual"
	}
	if q.Method != "penny_test" && q.Method != "micro_deposit" && q.Method != "bank_letter" && q.Method != "manual" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "unsupported verification method")
		return
	}
	sum := sha256.Sum256([]byte(q.AccountNumber))
	vid := id.New("bav")
	masked := "****" + q.AccountNumber[len(q.AccountNumber)-4:]
	_, err := h.pool.Exec(r.Context(), `INSERT INTO bank_account_verifications (id,merchant_id,bank_code,account_number_masked,account_number_hash,account_name,verification_method,amount,connector_id,status,expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7,1.00,$8,'pending',$9)`, vid, merchant, strings.ToUpper(q.BankCode), masked, hex.EncodeToString(sum[:]), q.AccountName, q.Method, q.ConnectorID, time.Now().Add(24*time.Hour))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, map[string]string{"id": vid, "status": "pending", "account_number_masked": masked, "message": "verification queued; no bank transfer is initiated until a partner connector is approved"})
}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	merchant := mw.MerchantID(r.Context())
	var status, masked, method string
	var verified *time.Time
	err := h.pool.QueryRow(r.Context(), `SELECT status,account_number_masked,verification_method,verified_at FROM bank_account_verifications WHERE id=$1 AND merchant_id=$2`, chi.URLParam(r, "id"), merchant).Scan(&status, &masked, &method, &verified)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "verification not found")
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{"status": status, "account_number_masked": masked, "method": method, "verified_at": verified})
}
