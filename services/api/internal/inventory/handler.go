package inventory

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
	// Products
	r.Post("/products", h.CreateProduct)
	r.Get("/products", h.ListProducts)
	r.Post("/products/{id}/stock", h.AddStock)
	// Orders (software POS)
	r.Post("/orders", h.CreateOrder)
	r.Get("/orders", h.ListOrders)
	r.Post("/orders/{id}/mark-paid", h.MarkPaid)
	// Stock movements
	r.Get("/stock-movements", h.StockMovements)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if p.Name == "" || p.Price == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "name and price required")
		return
	}
	if err := h.repo.CreateProduct(r.Context(), middleware.MerchantID(r.Context()), &p); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, p)
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListProducts(r.Context(), middleware.MerchantID(r.Context()))
	write(w, r, out, err)
}

func (h *Handler) AddStock(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Qty  string `json:"qty"`
		Note string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Qty == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "qty required")
		return
	}
	m, err := h.repo.AddStock(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"), req.Qty, req.Note)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, m)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if len(o.Items) == 0 {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "order requires items")
		return
	}
	if err := h.repo.CreateOrder(r.Context(), middleware.MerchantID(r.Context()), &o); err != nil {
		if err == ErrInsufficientStock {
			pkghttp.WriteErrorWithBody(w, r, 409, "insufficient_stock", err.Error())
			return
		}
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, o)
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListOrders(r.Context(), middleware.MerchantID(r.Context()), r.URL.Query().Get("status"), 50)
	write(w, r, out, err)
}

func (h *Handler) MarkPaid(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PaymentID string `json:"payment_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.repo.MarkPaid(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"), req.PaymentID); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"status": "paid"})
}

func (h *Handler) StockMovements(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.StockMovements(r.Context(), middleware.MerchantID(r.Context()), 50)
	write(w, r, out, err)
}

func write(w http.ResponseWriter, r *http.Request, out any, err error) {
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}
