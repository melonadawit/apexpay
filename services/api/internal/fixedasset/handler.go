package fixedasset

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
	r.Post("/{id}/depreciate", h.Depreciate)
	r.Get("/depreciation", h.Entries)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var a Asset
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if a.AssetName == "" || a.Category == "" || a.AcquisitionDate == "" || a.Cost == "" || a.UsefulLifeYears <= 0 {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "asset_name, category, acquisition_date, cost, useful_life_years required")
		return
	}
	if err := h.repo.CreateAsset(r.Context(), middleware.MerchantID(r.Context()), &a); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, a)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListAssets(r.Context(), middleware.MerchantID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) Depreciate(w http.ResponseWriter, r *http.Request) {
	e, err := h.repo.Depreciate(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"))
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "asset not found")
		return
	}
	pkghttp.WriteJSON(w, r, 200, e)
}

func (h *Handler) Entries(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.Entries(r.Context(), middleware.MerchantID(r.Context()), r.URL.Query().Get("asset_id"))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
