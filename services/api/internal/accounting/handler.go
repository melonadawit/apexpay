package accounting

import (
	"encoding/json"
	"net/http"
	"time"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repository
	svc  *Service
}

func NewHandler(repo *Repository, svc *Service) *Handler { return &Handler{repo: repo, svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Get("/accounts", h.ChartOfAccounts)
	r.Get("/trial-balance", h.TrialBalance)
	r.Get("/profit-loss", h.ProfitLoss)
	r.Get("/balance-sheet", h.BalanceSheet)
	r.Get("/cash-flow", h.CashFlow)

	// Real GL: manual journal entries + fiscal period close.
	r.Get("/journal-entries", h.ListJournalEntries)
	r.Post("/journal-entries", h.PostJournalEntry)
	r.Get("/periods", h.ListPeriods)
	r.Post("/periods/close", h.ClosePeriod)
	r.Post("/periods/reopen", h.ReopenPeriod)
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

// ---- Real GL: journal entries + fiscal periods ----

func (h *Handler) ListJournalEntries(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.ListJournalEntries(r.Context(), middleware.MerchantID(r.Context()))
	write(w, r, out, err)
}

func (h *Handler) PostJournalEntry(w http.ResponseWriter, r *http.Request) {
	var req JournalEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, http.StatusBadRequest, "validation_error", "invalid json")
		return
	}
	out, err := h.svc.PostJournalEntry(r.Context(), middleware.MerchantID(r.Context()), req)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, out)
}

func (h *Handler) ListPeriods(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.ListPeriods(r.Context(), middleware.MerchantID(r.Context()))
	write(w, r, out, err)
}

func (h *Handler) ClosePeriod(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Period string `json:"period"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	out, err := h.svc.ClosePeriod(r.Context(), middleware.MerchantID(r.Context()), req.Period, middleware.UserID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}

func (h *Handler) ReopenPeriod(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Period string `json:"period"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	out, err := h.svc.ReopenPeriod(r.Context(), middleware.MerchantID(r.Context()), req.Period)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}
