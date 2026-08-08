package banking

import (
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
	r.Get("/current_accounts", h.CurrentAccounts)
	r.Get("/credit_lines", h.CreditLines)
	r.Get("/forex/rates", h.ForexRates)
	r.Get("/forex/requests", h.ForexRequests)
	r.Get("/virtual_accounts", h.VirtualAccounts)
	r.Get("/notifications", h.Notifications)
	r.Get("/corporate_cards", h.CorporateCards)
	r.Get("/escrow", h.EscrowAccounts)
	r.Get("/support_tickets", h.SupportTickets)
	r.Get("/relationship_managers", h.RelationshipManagers)
	r.Get("/bank_verifications", h.BankVerifications)
}

func (h *Handler) CurrentAccounts(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.CurrentAccounts(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) CreditLines(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.CreditLines(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) ForexRates(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ForexRates(r.Context())
	writeList(w, r, out, err)
}

func (h *Handler) ForexRequests(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ForexRequests(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) VirtualAccounts(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.VirtualAccounts(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.Notifications(r.Context(), middleware.MerchantID(r.Context()), 50)
	writeList(w, r, out, err)
}

func (h *Handler) CorporateCards(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.CorporateCards(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) EscrowAccounts(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.EscrowAccounts(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) SupportTickets(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.SupportTickets(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) RelationshipManagers(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.RelationshipManagers(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) BankVerifications(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.BankVerifications(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

// writeList serializes a slice, ensuring empty (not null) JSON arrays.
func writeList(w http.ResponseWriter, r *http.Request, list any, err error) {
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
