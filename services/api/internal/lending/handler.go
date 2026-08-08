package lending

import (
	"encoding/json"
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Post("/{id}/repay", h.Repay)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CreditLineID string `json:"credit_line_id"`
		Amount       string `json:"amount"`
		Purpose      string `json:"purpose"`
		DueDate      string `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if req.CreditLineID == "" || req.Amount == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "credit_line_id, amount required")
		return
	}
	l := &Loan{Amount: req.Amount, Purpose: req.Purpose, DueDate: req.DueDate}
	if err := h.repo.CreateLoan(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), req.CreditLineID, l); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, l)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListLoans(r.Context(), middleware.MerchantID(r.Context()), r.URL.Query().Get("status"), 50)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}

func (h *Handler) Repay(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount string `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "amount required")
		return
	}
	if err := h.repo.Repay(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"), req.Amount); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"status": "repaid"})
}
