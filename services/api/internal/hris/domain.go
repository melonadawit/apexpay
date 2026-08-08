package hris

// Workforce OS (HRIS) domain types. Complements payroll with the HR lifecycle.

type Team struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	DepartmentID *string `json:"department_id,omitempty"`
	ManagerID    *string `json:"manager_id,omitempty"`
	Description  string  `json:"description,omitempty"`
	CreatedAt    string  `json:"created_at"`
}

type Contract struct {
	ID              string `json:"id"`
	EmployeeID      string `json:"employee_id"`
	ContractType    string `json:"contract_type"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date,omitempty"`
	SalaryAmount    string `json:"salary_amount"`
	SalaryCurrency  string `json:"salary_currency"`
	ProbationMonths int    `json:"probation_months"`
	NoticeDays      int    `json:"notice_days"`
	Status          string `json:"status"`
	DocKey          string `json:"doc_key,omitempty"`
	SignedAt        string `json:"signed_at,omitempty"`
}

type Shift struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	BreakMin  int    `json:"break_min"`
}

type AttendanceClock struct {
	ID         string  `json:"id"`
	EmployeeID string  `json:"employee_id"`
	ShiftID    *string `json:"shift_id,omitempty"`
	ClockDate  string  `json:"clock_date"`
	PunchIn    string  `json:"punch_in,omitempty"`
	PunchOut   string  `json:"punch_out,omitempty"`
	Hours      string  `json:"hours"`
	Status     string  `json:"status"`
	Source     string  `json:"source"`
	Note       string  `json:"note,omitempty"`
}

type OnboardingTask struct {
	ID         string `json:"id"`
	EmployeeID string `json:"employee_id"`
	Task       string `json:"task"`
	Category   string `json:"category"`
	DueInDays  int    `json:"due_in_days"`
	Status     string `json:"status"`
	AssignedTo string `json:"assigned_to,omitempty"`
}

type PerformanceReview struct {
	ID         string `json:"id"`
	EmployeeID string `json:"employee_id"`
	ReviewerID string `json:"reviewer_id,omitempty"`
	Period     string `json:"period"`
	Rating     *int   `json:"rating,omitempty"`
	Goals      string `json:"goals,omitempty"`
	Comments   string `json:"comments,omitempty"`
	Status     string `json:"status"`
}
