package team

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
	// Members
	r.Post("/members", h.Invite)
	r.Get("/members", h.List)
	r.Put("/members/{userID}/role", h.UpdateRole)
	r.Delete("/members/{userID}", h.Remove)
	// Approvals
	r.Post("/approvals", h.CreateApproval)
	r.Get("/approvals", h.ListApprovals)
	r.Post("/approvals/{id}/decide", h.DecideApproval)
}

func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
	var req InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	m, err := h.repo.InviteMember(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), req)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", err.Error())
		return
	}
	pkghttp.WriteJSON(w, r, 201, m)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListMembers(r.Context(), middleware.MerchantID(r.Context()))
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if err := h.repo.UpdateRole(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "userID"), req.Role); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", err.Error())
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"status": "updated"})
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.RemoveMember(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "userID")); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "member not found")
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"status": "removed"})
}

func (h *Handler) CreateApproval(w http.ResponseWriter, r *http.Request) {
	var a Approval
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if a.ResourceType == "" || a.ResourceID == "" || a.Action == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "resource_type, resource_id, action required")
		return
	}
	if err := h.repo.CreateApproval(r.Context(), middleware.MerchantID(r.Context()), middleware.UserID(r.Context()), &a); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, a)
}

func (h *Handler) ListApprovals(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListApprovals(r.Context(), middleware.MerchantID(r.Context()), r.URL.Query().Get("status"), 50)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}

func (h *Handler) DecideApproval(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Decision string `json:"decision"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if req.Decision != "approve" && req.Decision != "reject" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "decision must be approve or reject")
		return
	}
	userID := middleware.UserID(r.Context())
	role := middleware.Role(r.Context())
	approval, final, err := h.repo.DecideApproval(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"), userID, userID, role, req.Decision)
	if err != nil {
		pkghttp.WriteErrorWithBody(w, r, 409, "conflict", err.Error())
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"approval": approval, "final": final, "status": approval.Status,
	})
}
