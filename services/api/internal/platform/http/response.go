package http

import (
	"encoding/json"
	"net/http"

	appErrors "apexpay/internal/platform/errors"
	"github.com/go-chi/chi/v5/middleware"
)

// Standard API response with request_id for audit per SAD §11 correlation

type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorBody  `json:"error,omitempty"`
	RequestID string      `json:"request_id"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := Response{
		Success:   status < 400,
		Data:      data,
		RequestID: middleware.GetReqID(r.Context()),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	// Map AppError to stable code per SAD
	if appErr, ok := err.(*appErrors.AppError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.Status())
		resp := Response{
			Success:   false,
			Error:     &ErrorBody{Code: string(appErr.Code), Message: appErr.Message},
			RequestID: middleware.GetReqID(r.Context()),
		}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	// generic 500
	WriteJSON(w, r, http.StatusInternalServerError, map[string]string{"message": "internal error"})
}

func WriteErrorWithBody(w http.ResponseWriter, r *http.Request, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := Response{
		Success:   false,
		Error:     &ErrorBody{Code: code, Message: msg},
		RequestID: middleware.GetReqID(r.Context()),
	}
	_ = json.NewEncoder(w).Encode(resp)
}
