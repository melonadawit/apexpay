package auth

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"

	"apexpay/internal/i18n"
	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes(r chi.Router) {
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)
	r.Get("/me", h.Me)
	r.Get("/language", h.GetLanguage)
	r.Put("/language", h.SetLanguage)
}

// GetLanguage returns the caller's language preference.
func (h *Handler) GetLanguage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, "unauthorized", "no session")
		return
	}
	u, err := h.svc.repo.findUserByID(r.Context(), userID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{
		"language_preference": u.LanguagePreference,
	})
}

// SetLanguage updates the caller's language preference to 'en' or 'am'.
func (h *Handler) SetLanguage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		pkghttp.WriteErrorWithBody(w, r, http.StatusUnauthorized, "unauthorized", "no session")
		return
	}
	var req struct {
		LanguagePreference string `json:"language_preference"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, http.StatusBadRequest, "validation_error", "invalid json")
		return
	}
	if !i18n.IsValid(req.LanguagePreference) {
		pkghttp.WriteErrorWithBody(w, r, http.StatusBadRequest, "validation_error", i18n.New().Get(i18n.DefaultLocale, "language_preference_required"))
		return
	}
	u, err := h.svc.SetLanguagePreference(r.Context(), userID, req.LanguagePreference)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, u)
}

func bearerToken(r *http.Request) string {
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		if i := strings.IndexByte(ip, ','); i > 0 {
			ip = ip[:i]
		}
		return strings.TrimSpace(ip)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// Login accepts {email, password} and returns an opaque session token.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if req.Email == "" || req.Password == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "email and password required")
		return
	}
	res, err := h.svc.Login(r.Context(), req.Email, req.Password, r.UserAgent(), clientIP(r))
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			pkghttp.WriteErrorWithBody(w, r, 401, "unauthorized", "invalid email or password")
			return
		}
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, res)
}

// Logout revokes the current session token.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Logout(r.Context(), bearerToken(r)); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"status": "logged_out"})
}

// Me returns the current user + merchant context. Requires the session middleware.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	if userID == "" {
		pkghttp.WriteErrorWithBody(w, r, 401, "unauthorized", "no session")
		return
	}
	res, err := h.svc.Me(r.Context(), userID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, res)
}
