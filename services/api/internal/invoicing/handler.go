package invoicing

import (
	"encoding/json"
	"net/http"

	"apexpay/internal/id"
	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/aging", h.Aging)
	r.Get("/{id}", h.Get)
	r.Post("/{id}/issue", h.Issue)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if req.InvoiceNumber == "" || req.CustomerName == "" || req.IssueDate == "" || req.DueDate == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invoice_number, customer_name, issue_date, due_date required")
		return
	}
	inv, err := h.svc.BuildInvoice(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), req)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, inv)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.repo.ListInvoices(r.Context(), middleware.MerchantID(r.Context()), r.URL.Query().Get("status"), 50)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	inv, err := h.svc.repo.GetInvoice(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"))
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "invoice not found")
		return
	}
	pkghttp.WriteJSON(w, r, 200, inv)
}

func (h *Handler) Aging(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.repo.Aging(r.Context(), middleware.MerchantID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}

// Issue marks a draft invoice as sent with a hosted payment token.
func (h *Handler) Issue(w http.ResponseWriter, r *http.Request) {
	token := "inv_" + id.New("tok")[5:]
	inv, err := h.svc.Issue(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"), token)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "invoice not found")
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"invoice":    inv,
		"hosted_url": "https://checkout.apexpay.et/pay/" + token,
	})
}
