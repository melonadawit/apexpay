package accounting

import (
	"net/http"
	"time"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) Routes(r chi.Router) {
	r.Get("/accounts", h.ChartOfAccounts)
	r.Get("/trial-balance", h.TrialBalance)
	r.Get("/profit-loss", h.ProfitLoss)
	r.Get("/balance-sheet", h.BalanceSheet)
	r.Get("/cash-flow", h.CashFlow)
}

func (h *Handler) ChartOfAccounts(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ChartOfAccounts(r.Context(), middleware.MerchantID(r.Context()))
	write(w, r, out, err)
}

func (h *Handler) TrialBalance(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.TrialBalance(r.Context(), middleware.MerchantID(r.Context()))
	write(w, r, out, err)
}

func (h *Handler) ProfitLoss(w http.ResponseWriter, r *http.Request) {
	from, to := period(r)
	out, err := h.repo.ProfitLoss(r.Context(), middleware.MerchantID(r.Context()), from, to)
	write(w, r, out, err)
}

func (h *Handler) BalanceSheet(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.BalanceSheet(r.Context(), middleware.MerchantID(r.Context()), time.Now().Format("2006-01-02"))
	write(w, r, out, err)
}

func (h *Handler) CashFlow(w http.ResponseWriter, r *http.Request) {
	from, to := period(r)
	out, err := h.repo.CashFlow(r.Context(), middleware.MerchantID(r.Context()), from, to)
	write(w, r, out, err)
}

func period(r *http.Request) (string, string) {
	from := r.URL.Query().Get("from")
	if from == "" {
		from = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	to := r.URL.Query().Get("to")
	if to == "" {
		to = time.Now().Format("2006-01-02")
	}
	return from, to
}

func write(w http.ResponseWriter, r *http.Request, out any, err error) {
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}
