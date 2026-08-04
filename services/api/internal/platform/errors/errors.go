package errors

import "net/http"

// Domain errors -> stable API error codes per SAD §11.
type Code string

const (
	CodeDuplicateTxRef       Code = "duplicate_tx_ref"
	CodeDuplicateRefundRef   Code = "duplicate_refund_ref"
	CodeInvalidFaydaFin      Code = "invalid_fayda_fin"
	CodeFaydaOTPFailed       Code = "fayda_otp_failed"
	CodeFaydaNotVerified     Code = "fayda_not_verified"
	CodeOnboardingNotReady   Code = "onboarding_not_ready"
	CodeDocumentRequired     Code = "document_required"
	CodeBankVerificationFail Code = "bank_verification_failed"
	CodeInsufficientBalance  Code = "insufficient_balance"
	CodeNotFound             Code = "not_found"
	CodeUnauthorized         Code = "unauthorized"
	CodeForbidden            Code = "forbidden"
	CodeValidation           Code = "validation_error"
	CodeConflict             Code = "conflict"
	CodeRateLimited          Code = "rate_limited"
	CodeConnectorDown        Code = "connector_unavailable"
	CodeRefundExceeded       Code = "refund_exceeded"
	CodePayrollNotCalculable Code = "payroll_not_calculable"
)

type AppError struct {
	Code       Code   `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string { return string(e.Code) + ": " + e.Message }

func New(code Code, msg string, status int) *AppError {
	return &AppError{Code: code, Message: msg, HTTPStatus: status}
}

func NotFound(msg string) *AppError       { return New(CodeNotFound, msg, http.StatusNotFound) }
func Validation(msg string) *AppError     { return New(CodeValidation, msg, http.StatusBadRequest) }
func Unauthorized(msg string) *AppError   { return New(CodeUnauthorized, msg, http.StatusUnauthorized) }
func Forbidden(msg string) *AppError      { return New(CodeForbidden, msg, http.StatusForbidden) }
func Conflict(code Code, msg string) *AppError { return New(code, msg, http.StatusConflict) }

// HTTPStatus returns status or 500.
func (e *AppError) Status() int {
	if e.HTTPStatus == 0 {
		return http.StatusInternalServerError
	}
	return e.HTTPStatus
}
