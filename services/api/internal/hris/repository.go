package hris

import (
	"context"
	"errors"
	"fmt"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// ---- Teams ----

func (r *Repository) CreateTeam(ctx context.Context, merchantID string, t *Team) error {
	t.ID = id.New("team")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO hris_teams (id, merchant_id, name, department_id, manager_id, description)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		t.ID, merchantID, t.Name, nilStr(t.DepartmentID), nilStr(t.ManagerID), t.Description)
	return err
}

func (r *Repository) ListTeams(ctx context.Context, merchantID string) ([]Team, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, department_id, manager_id, COALESCE(description,''),
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM hris_teams WHERE merchant_id=$1 ORDER BY name`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Team{}
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.Name, &t.DepartmentID, &t.ManagerID, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

// ---- Contracts ----

func (r *Repository) CreateContract(ctx context.Context, merchantID string, c *Contract) error {
	c.ID = id.New("contr")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO hris_contracts (id, merchant_id, employee_id, contract_type, start_date, end_date, salary_amount, salary_currency, probation_months, notice_days, status, doc_key, signed_at)
		VALUES ($1,$2,$3,$4,$5::date,NULLIF($6,'')::date,$7::numeric,$8,$9,$10,$11,$12,$13)`,
		c.ID, merchantID, c.EmployeeID, c.ContractType, c.StartDate, c.EndDate, c.SalaryAmount, c.SalaryCurrency,
		c.ProbationMonths, c.NoticeDays, c.Status, nilStr(&c.DocKey), nilTime(c.SignedAt))
	return err
}

func (r *Repository) ListContracts(ctx context.Context, merchantID, employeeID string) ([]Contract, error) {
	query := `SELECT id, employee_id, contract_type, to_char(start_date,'YYYY-MM-DD'), COALESCE(to_char(end_date,'YYYY-MM-DD'),''),
		COALESCE(salary_amount,0)::text, salary_currency, probation_months, notice_days, status, COALESCE(doc_key,''), COALESCE(to_char(signed_at,'YYYY-MM-DD'),'')
		FROM hris_contracts WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if employeeID != "" {
		query += ` AND employee_id=$2`
		args = append(args, employeeID)
	}
	query += ` ORDER BY start_date DESC`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Contract{}
	for rows.Next() {
		var c Contract
		if err := rows.Scan(&c.ID, &c.EmployeeID, &c.ContractType, &c.StartDate, &c.EndDate, &c.SalaryAmount,
			&c.SalaryCurrency, &c.ProbationMonths, &c.NoticeDays, &c.Status, &c.DocKey, &c.SignedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// ---- Shifts ----

func (r *Repository) CreateShift(ctx context.Context, merchantID string, s *Shift) error {
	s.ID = id.New("shift")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO hris_shifts (id, merchant_id, name, start_time, end_time, break_min)
		VALUES ($1,$2,$3,$4::time,$5::time,$6)`,
		s.ID, merchantID, s.Name, s.StartTime, s.EndTime, s.BreakMin)
	return err
}

func (r *Repository) ListShifts(ctx context.Context, merchantID string) ([]Shift, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, to_char(start_time,'HH24:MI'), to_char(end_time,'HH24:MI'), break_min
		FROM hris_shifts WHERE merchant_id=$1 ORDER BY start_time`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Shift{}
	for rows.Next() {
		var s Shift
		if err := rows.Scan(&s.ID, &s.Name, &s.StartTime, &s.EndTime, &s.BreakMin); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// ---- Attendance clock ----

func (r *Repository) PunchIn(ctx context.Context, merchantID, employeeID string, shiftID *string, at time.Time, source string) (*AttendanceClock, error) {
	// Idempotent on (merchant_id, employee_id, clock_date): re-punch updates the punch-in.
	date := at.In(tzEAT()).Format("2006-01-02")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO hris_attendance_clock (id, merchant_id, employee_id, shift_id, clock_date, punch_in, status, source)
		VALUES ($1,$2,$3,$4,$5::date,$6,'present',$7)
		ON CONFLICT (merchant_id, employee_id, clock_date) DO UPDATE SET punch_in=EXCLUDED.punch_in, shift_id=EXCLUDED.shift_id, updated_at=now()`,
		id.New("clock"), merchantID, employeeID, nilStr(shiftID), date, at, source)
	if err != nil {
		return nil, err
	}
	return &AttendanceClock{EmployeeID: employeeID, ClockDate: date, PunchIn: at.Format(time.RFC3339)}, nil
}

func (r *Repository) PunchOut(ctx context.Context, merchantID, employeeID string, at time.Time, source string) (*AttendanceClock, error) {
	date := at.In(tzEAT()).Format("2006-01-02")
	// Compute hours from existing punch_in.
	var hours float64
	row := r.pool.QueryRow(ctx, `SELECT punch_in FROM hris_attendance_clock WHERE merchant_id=$1 AND employee_id=$2 AND clock_date=$3::date`,
		merchantID, employeeID, date)
	var pin time.Time
	if err := row.Scan(&pin); err != nil {
		return nil, ErrNotFound
	}
	if at.After(pin) {
		hours = at.Sub(pin).Minutes() / 60
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE hris_attendance_clock SET punch_out=$1, hours=$2, updated_at=now()
		WHERE merchant_id=$3 AND employee_id=$4 AND clock_date=$5::date`,
		at, fmt.Sprintf("%.2f", hours), merchantID, employeeID, date)
	if err != nil {
		return nil, err
	}
	return &AttendanceClock{EmployeeID: employeeID, ClockDate: date, PunchOut: at.Format(time.RFC3339), Hours: fmt.Sprintf("%.2f", hours)}, nil
}

func (r *Repository) ListAttendance(ctx context.Context, merchantID, employeeID string, from, to string) ([]AttendanceClock, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, shift_id, to_char(clock_date,'YYYY-MM-DD'), COALESCE(to_char(punch_in,'YYYY-MM-DD"T"HH24:MI:SS'),''),
		       COALESCE(to_char(punch_out,'YYYY-MM-DD"T"HH24:MI:SS'),''), hours::text, status, source, COALESCE(note,'')
		FROM hris_attendance_clock
		WHERE merchant_id=$1 AND clock_date BETWEEN $2::date AND $3::date
		ORDER BY clock_date DESC`, merchantID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []AttendanceClock{}
	for rows.Next() {
		var a AttendanceClock
		if err := rows.Scan(&a.ID, &a.EmployeeID, &a.ShiftID, &a.ClockDate, &a.PunchIn, &a.PunchOut, &a.Hours, &a.Status, &a.Source, &a.Note); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

// ---- Onboarding checklist ----

func (r *Repository) CreateOnboardingTask(ctx context.Context, merchantID string, t *OnboardingTask) error {
	t.ID = id.New("obt")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO hris_onboarding_checklists (id, merchant_id, employee_id, task, category, due_in_days, status, assigned_to)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		t.ID, merchantID, t.EmployeeID, t.Task, t.Category, t.DueInDays, t.Status, nilStr(&t.AssignedTo))
	return err
}

func (r *Repository) ListOnboardingTasks(ctx context.Context, merchantID, employeeID string) ([]OnboardingTask, error) {
	query := `SELECT id, employee_id, task, category, due_in_days, status, COALESCE(assigned_to,'') FROM hris_onboarding_checklists WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if employeeID != "" {
		query += ` AND employee_id=$2`
		args = append(args, employeeID)
	}
	query += ` ORDER BY created_at`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []OnboardingTask{}
	for rows.Next() {
		var t OnboardingTask
		if err := rows.Scan(&t.ID, &t.EmployeeID, &t.Task, &t.Category, &t.DueInDays, &t.Status, &t.AssignedTo); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

// ---- Performance reviews ----

func (r *Repository) CreateReview(ctx context.Context, merchantID string, rev *PerformanceReview) error {
	rev.ID = id.New("perf")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO hris_performance_reviews (id, merchant_id, employee_id, reviewer_id, period, rating, goals, comments, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		rev.ID, merchantID, rev.EmployeeID, nilStr(&rev.ReviewerID), rev.Period, rev.Rating, rev.Goals, rev.Comments, rev.Status)
	return err
}

func (r *Repository) ListReviews(ctx context.Context, merchantID, employeeID string) ([]PerformanceReview, error) {
	query := `SELECT id, employee_id, COALESCE(reviewer_id,''), period, rating, COALESCE(goals,''), COALESCE(comments,''), status FROM hris_performance_reviews WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if employeeID != "" {
		query += ` AND employee_id=$2`
		args = append(args, employeeID)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []PerformanceReview{}
	for rows.Next() {
		var r PerformanceReview
		if err := rows.Scan(&r.ID, &r.EmployeeID, &r.ReviewerID, &r.Period, &r.Rating, &r.Goals, &r.Comments, &r.Status); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, rows.Err()
}

// ---- helpers ----

func nilStr(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}

func nilTime(s string) interface{} {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return t
}

func tzEAT() *time.Location { return time.FixedZone("EAT", 3*3600) }
