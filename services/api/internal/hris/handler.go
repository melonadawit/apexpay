package hris

import (
	"encoding/json"
	"net/http"
	"time"

	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) Routes(r chi.Router) {
	// Teams / org
	r.Post("/teams", h.CreateTeam)
	r.Get("/teams", h.ListTeams)
	// Contracts
	r.Post("/contracts", h.CreateContract)
	r.Get("/contracts", h.ListContracts)
	r.Get("/employees/{id}/contracts", h.ListEmployeeContracts)
	// Shifts
	r.Post("/shifts", h.CreateShift)
	r.Get("/shifts", h.ListShifts)
	// Attendance clock
	r.Post("/attendance/clock-in", h.PunchIn)
	r.Post("/attendance/clock-out", h.PunchOut)
	r.Get("/attendance", h.ListAttendance)
	r.Get("/employees/{id}/attendance", h.ListEmployeeAttendance)
	// Onboarding checklist
	r.Post("/onboarding-tasks", h.CreateOnboardingTask)
	r.Get("/onboarding-tasks", h.ListOnboardingTasks)
	r.Get("/employees/{id}/onboarding-tasks", h.ListEmployeeOnboardingTasks)
	// Performance reviews
	r.Post("/performance-reviews", h.CreateReview)
	r.Get("/performance-reviews", h.ListReviews)
	r.Get("/employees/{id}/performance-reviews", h.ListEmployeeReviews)
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var t Team
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if t.Name == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "name required")
		return
	}
	if err := h.repo.CreateTeam(r.Context(), middleware.MerchantID(r.Context()), &t); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, t)
}

func (h *Handler) ListTeams(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListTeams(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) CreateContract(w http.ResponseWriter, r *http.Request) {
	var c Contract
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if c.EmployeeID == "" || c.ContractType == "" || c.StartDate == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "employee_id, contract_type, start_date required")
		return
	}
	if c.SalaryCurrency == "" {
		c.SalaryCurrency = "ETB"
	}
	if c.ProbationMonths == 0 {
		c.ProbationMonths = 3
	}
	if c.NoticeDays == 0 {
		c.NoticeDays = 30
	}
	if c.Status == "" {
		c.Status = "draft"
	}
	if err := h.repo.CreateContract(r.Context(), middleware.MerchantID(r.Context()), &c); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, c)
}

func (h *Handler) ListContracts(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListContracts(r.Context(), middleware.MerchantID(r.Context()), "")
	writeList(w, r, out, err)
}

func (h *Handler) ListEmployeeContracts(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListContracts(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"))
	writeList(w, r, out, err)
}

func (h *Handler) CreateShift(w http.ResponseWriter, r *http.Request) {
	var s Shift
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if s.Name == "" || s.StartTime == "" || s.EndTime == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "name, start_time, end_time required")
		return
	}
	if err := h.repo.CreateShift(r.Context(), middleware.MerchantID(r.Context()), &s); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, s)
}

func (h *Handler) ListShifts(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListShifts(r.Context(), middleware.MerchantID(r.Context()))
	writeList(w, r, out, err)
}

func (h *Handler) PunchIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeID string `json:"employee_id"`
		ShiftID    string `json:"shift_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.EmployeeID == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "employee_id required")
		return
	}
	var shiftPtr *string
	if req.ShiftID != "" {
		shiftPtr = &req.ShiftID
	}
	out, err := h.repo.PunchIn(r.Context(), middleware.MerchantID(r.Context()), req.EmployeeID, shiftPtr, time.Now(), "manual")
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, out)
}

func (h *Handler) PunchOut(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeID string `json:"employee_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.EmployeeID == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "employee_id required")
		return
	}
	out, err := h.repo.PunchOut(r.Context(), middleware.MerchantID(r.Context()), req.EmployeeID, time.Now(), "manual")
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, out)
}

func (h *Handler) ListAttendance(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListAttendance(r.Context(), middleware.MerchantID(r.Context()), "", qFrom(r), qTo(r))
	writeList(w, r, out, err)
}

func (h *Handler) ListEmployeeAttendance(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListAttendance(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"), qFrom(r), qTo(r))
	writeList(w, r, out, err)
}

func (h *Handler) CreateOnboardingTask(w http.ResponseWriter, r *http.Request) {
	var t OnboardingTask
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if t.EmployeeID == "" || t.Task == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "employee_id, task required")
		return
	}
	if t.DueInDays == 0 {
		t.DueInDays = 7
	}
	if err := h.repo.CreateOnboardingTask(r.Context(), middleware.MerchantID(r.Context()), &t); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, t)
}

func (h *Handler) ListOnboardingTasks(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListOnboardingTasks(r.Context(), middleware.MerchantID(r.Context()), "")
	writeList(w, r, out, err)
}

func (h *Handler) ListEmployeeOnboardingTasks(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListOnboardingTasks(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"))
	writeList(w, r, out, err)
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var rev PerformanceReview
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if rev.EmployeeID == "" || rev.Period == "" {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "employee_id, period required")
		return
	}
	if rev.Status == "" {
		rev.Status = "draft"
	}
	if err := h.repo.CreateReview(r.Context(), middleware.MerchantID(r.Context()), &rev); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, rev)
}

func (h *Handler) ListReviews(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListReviews(r.Context(), middleware.MerchantID(r.Context()), "")
	writeList(w, r, out, err)
}

func (h *Handler) ListEmployeeReviews(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListReviews(r.Context(), middleware.MerchantID(r.Context()), chi.URLParam(r, "id"))
	writeList(w, r, out, err)
}

func qFrom(r *http.Request) string {
	f := r.URL.Query().Get("from")
	if f == "" {
		return "2020-01-01"
	}
	return f
}

func qTo(r *http.Request) string {
	t := r.URL.Query().Get("to")
	if t == "" {
		return "2099-12-31"
	}
	return t
}

func writeList(w http.ResponseWriter, r *http.Request, list any, err error) {
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, list)
}
