package banking

import (
	"encoding/json"
	"errors"
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

	// Action-oriented modules (read + write).
	r.Get("/vendor_invoices", h.VendorInvoices)
	r.Post("/vendor_invoices", h.CreateVendorInvoice)
	r.Get("/petty_cash_budgets", h.PettyCashBudgets)
	r.Post("/petty_cash_budgets", h.CreatePettyCashBudget)
	r.Get("/petty_cash_expenses", h.PettyCashExpenses)
	r.Post("/petty_cash_expenses", h.CreatePettyCashExpense)
	r.Get("/tax_payments", h.TaxPayments)
	r.Post("/tax_payments", h.CreateTaxPayment)
	r.Get("/payout_links", h.PayoutLinks)
	r.Post("/payout_links", h.CreatePayoutLink)
	r.Get("/accounting_integrations", h.AccountingIntegrations)
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

// ---- Action modules: read + write ----

func (h *Handler) VendorInvoices(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.VendorInvoices(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) CreateVendorInvoice(w http.ResponseWriter, r *http.Request) {
	var in VendorInvoice
	if err := decodeJSON(r, &in); err != nil {
		writeErr(w, r, err)
		return
	}
	if err := h.repo.CreateVendorInvoice(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), &in); err != nil {
		writeErr(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, in)
}

func (h *Handler) PettyCashBudgets(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.PettyCashBudgets(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) CreatePettyCashBudget(w http.ResponseWriter, r *http.Request) {
	var in PettyCashBudget
	if err := decodeJSON(r, &in); err != nil {
		writeErr(w, r, err)
		return
	}
	if err := h.repo.CreatePettyCashBudget(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), &in); err != nil {
		writeErr(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, in)
}

func (h *Handler) PettyCashExpenses(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.PettyCashExpenses(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) CreatePettyCashExpense(w http.ResponseWriter, r *http.Request) {
	var in PettyCashExpense
	if err := decodeJSON(r, &in); err != nil {
		writeErr(w, r, err)
		return
	}
	if err := h.repo.CreatePettyCashExpense(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), &in); err != nil {
		writeErr(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, in)
}

func (h *Handler) TaxPayments(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.TaxPayments(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) CreateTaxPayment(w http.ResponseWriter, r *http.Request) {
	var in TaxPayment
	if err := decodeJSON(r, &in); err != nil {
		writeErr(w, r, err)
		return
	}
	if err := h.repo.CreateTaxPayment(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), &in); err != nil {
		writeErr(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, in)
}

func (h *Handler) PayoutLinks(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.PayoutLinks(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) CreatePayoutLink(w http.ResponseWriter, r *http.Request) {
	var in PayoutLink
	if err := decodeJSON(r, &in); err != nil {
		writeErr(w, r, err)
		return
	}
	if err := h.repo.CreatePayoutLink(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), &in); err != nil {
		writeErr(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, in)
}

// decodeJSON decodes a request body into v, returning a user-friendly error.
func decodeJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New("invalid json body")
	}
	return nil
}

// writeErr serializes an error via the platform helper.
func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	pkghttp.WriteError(w, r, err)
}

func (h *Handler) AccountingIntegrations(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.AccountingIntegrations(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}
