package assistant

import (
	"encoding/json"
	"net/http"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc  *Service
	repo *Repository
}

func NewHandler(svc *Service, repo *Repository) *Handler { return &Handler{svc: svc, repo: repo} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/chat", h.Chat)
	r.Get("/threads/{id}", h.Thread)
}

// Chat is the session-authenticated assistant entrypoint. It resolves the actor from the
// user: if the user maps to an employee row they get the employee scope, otherwise merchant.
func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	merchantID := middleware.MerchantID(r.Context())
	userID := middleware.UserID(r.Context())
	if userID == "" {
		pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, "unauthorized", "session required for assistant")
		return
	}
	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, http.StatusBadRequest, "validation_error", "invalid json")
		return
	}

	// Resolve actor: employee if the user has an active employee record in this merchant.
	actor := ActorMerchant
	empID, isEmp, err := h.repo.EmployeeIDForUser(r.Context(), merchantID, userID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	if isEmp {
		actor = ActorEmployee
	}

	scope := Scope{MerchantID: merchantID, UserID: userID, Actor: actor}
	if actor == ActorEmployee {
		scope.EmployeeID = empID
	}

	reply, err := h.svc.Chat(r.Context(), scope, req.Message)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, http.StatusOK, reply)
}

// Thread returns a conversation and its messages, gated to the owner.
func (h *Handler) Thread(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	thread, err := h.repo.GetThread(r.Context(), chi.URLParam(r, "id"), userID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	if thread == nil {
		pkghttp.WriteErrorWithBody(w, r, http.StatusNotFound, "not_found", "thread not found")
		return
	}
	msgs, err := h.repo.ListMessages(r.Context(), thread.ID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, http.StatusOK, map[string]any{"thread": thread, "messages": msgs})
}
