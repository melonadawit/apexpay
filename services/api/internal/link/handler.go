package link

import (
	"encoding/json"
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{token}", h.GetByToken)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	var req struct {
		Amount      string `json:"amount"`
		Currency    string `json:"currency"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	amt, err := decimal.NewFromString(req.Amount)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "amount numeric")
		return
	}
	pl, cs, err := h.svc.Create(r.Context(), CreateRequest{
		MerchantID: merchantID, Amount: amt, Currency: req.Currency, Description: req.Description,
	})
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, map[string]interface{}{
		"id": pl.ID, "amount": pl.Amount.String(), "currency": pl.Currency, "description": pl.Description,
		"status": pl.Status, "public_token": pl.PublicToken,
		"checkout_url":        "https://checkout.apexpay.et/c/" + pl.PublicToken,
		"qr_data_url":         "data:image/png;base64,mock_qr_" + pl.PublicToken,
		"checkout_session_id": cs.ID,
		"expires_at":          pl.ExpiresAt,
		"share": map[string]string{
			"telegram": "https://t.me/share/url?url=https://checkout.apexpay.et/c/" + pl.PublicToken,
			"whatsapp": "https://wa.me/?text=Pay%20ETB%20" + pl.Amount.String() + "%20https://checkout.apexpay.et/c/" + pl.PublicToken,
		},
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	list, err := h.svc.List(r.Context(), merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) GetByToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	pl, err := h.svc.GetByToken(r.Context(), token)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, pl)
}
