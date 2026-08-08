package procurement

import (
	"encoding/json"
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ repo *Repository }

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) Routes(r chi.Router) {
	r.Get("/vendors", h.ListVendors)
	r.Post("/vendors", h.CreateVendor)
	r.Get("/purchase-orders", h.ListPOs)
	r.Post("/purchase-orders", h.CreatePO)
	r.Get("/purchase-orders/{id}", h.GetPO)
	r.Post("/purchase-orders/{id}/receive", h.Receive)
	r.Get("/invoices", h.ListAPInvoices)
	r.Post("/invoices", h.CreateAPInvoice)
	r.Get("/aging", h.Aging)
}

func (h *Handler) ListVendors(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListVendors(r.Context(), middleware.MerchantID(r.Context()))
	write(w, r, out, err)
}

func (h *Handler) CreateVendor(w http.ResponseWriter, r *http.Request) {
	var in VendorInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Name == "" {
		pkghttp.WriteErrorWithBody(w, r, http.StatusBadRequest, "validation_error", "name required")
		return
	}
	out, err := h.repo.CreateVendor(r.Context(), middleware.MerchantID(r.Context()), in)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, http.StatusCreated, out)
}

func (h *Handler) ListPOs(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListPOs(r.Context(), middleware.MerchantID(r.Context()))
	write(w, r, out, err)
}

func (h *Handler) GetPO(w http.ResponseWriter, r *http.Request) {
	po, err := h.repo.GetPO(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	if po == nil {
		pkghttp.WriteErrorWithBody(w, r, http.StatusNotFound, "not_found", "po not found")
		return
	}
	pkghttp.WriteJSON(w, r, http.StatusOK, po)
}

func (h *Handler) CreatePO(w http.ResponseWriter, r *http.Request) {
	var in POInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.VendorID == "" || in.PONumber == "" {
		pkghttp.WriteErrorWithBody(w, r, http.StatusBadRequest, "validation_error", "vendor_id, po_number required")
		return
	}
	out, err := h.repo.CreatePO(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), in)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, http.StatusCreated, out)
}

func (h *Handler) Receive(w http.ResponseWriter, r *http.Request) {
	rec, err := h.repo.Receive(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), chi.URLParam(r, "id"))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, http.StatusCreated, rec)
}

func (h *Handler) ListAPInvoices(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListAPInvoices(r.Context(), middleware.MerchantID(r.Context()))
	write(w, r, out, err)
}

func (h *Handler) CreateAPInvoice(w http.ResponseWriter, r *http.Request) {
	var in APInvoiceInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.VendorID == "" || in.InvoiceNumber == "" {
		pkghttp.WriteErrorWithBody(w, r, http.StatusBadRequest, "validation_error", "vendor_id, invoice_number required")
		return
	}
	out, err := h.repo.CreateAPInvoice(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), in)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, http.StatusCreated, out)
}

func (h *Handler) Aging(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.Aging(r.Context(), middleware.MerchantID(r.Context()))
	write(w, r, out, err)
}

func write(w http.ResponseWriter, r *http.Request, out any, err error) {
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, http.StatusOK, out)
}
