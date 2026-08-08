package team

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// validRoles are the roles a merchant can assign to members.
var validRoles = map[string]bool{
	"owner": true, "admin": true, "developer": true, "finance": true,
	"support": true, "ops": true, "compliance": true, "viewer": true,
}

func ValidRole(role string) bool { return validRoles[role] }

// InviteMember creates (or reuses) a user and adds a merchant membership. Returns the member.
func (r *Repository) InviteMember(ctx context.Context, merchantID, actorID string, req InviteRequest) (*Member, error) {
	if req.Email == "" || req.Name == "" || !ValidRole(req.Role) {
		return nil, errors.New("email, name and a valid role are required")
	}
	perms, _ := json.Marshal(req.Permissions)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Reuse an existing user by email, or create one.
	var userID, existingName string
	err = tx.QueryRow(ctx, `SELECT id, name FROM users WHERE lower(email)=lower($1)`, req.Email).Scan(&userID, &existingName)
	if err != nil {
		// not found -> create
		userID = id.New("user")
		if _, err := tx.Exec(ctx, `INSERT INTO users (id, email, name, status) VALUES ($1,$2,$3,'pending_verification')`,
			userID, req.Email, req.Name); err != nil {
			return nil, err
		}
	} else if req.Name == "" {
		req.Name = existingName
	}

	if _, err := tx.Exec(ctx, `INSERT INTO merchant_members (merchant_id, user_id, role, permissions)
		VALUES ($1,$2,$3,$4::jsonb)
		ON CONFLICT (merchant_id, user_id) DO UPDATE SET role=EXCLUDED.role, permissions=EXCLUDED.permissions`,
		merchantID, userID, req.Role, string(perms)); err != nil {
		return nil, err
	}

	return &Member{UserID: userID, Email: req.Email, Name: req.Name, Role: req.Role, Permissions: req.Permissions}, tx.Commit(ctx)
}

// ListMembers returns all members of a merchant with their user info.
func (r *Repository) ListMembers(ctx context.Context, merchantID string) ([]Member, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT u.id, u.email, u.name, mm.role, COALESCE(mm.permissions,'[]')::text,
		       to_char(mm.created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM merchant_members mm JOIN users u ON u.id = mm.user_id
		WHERE mm.merchant_id=$1 ORDER BY u.name`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Member{}
	for rows.Next() {
		var m Member
		var perms string
		if err := rows.Scan(&m.UserID, &m.Email, &m.Name, &m.Role, &perms, &m.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(perms), &m.Permissions)
		list = append(list, m)
	}
	return list, rows.Err()
}

// UpdateRole changes a member's role (owners cannot be demoted by non-owners is enforced
// in the handler/service; here it just updates).
func (r *Repository) UpdateRole(ctx context.Context, merchantID, userID, role string) error {
	if !ValidRole(role) {
		return errors.New("invalid role")
	}
	ct, err := r.pool.Exec(ctx, `UPDATE merchant_members SET role=$1 WHERE merchant_id=$2 AND user_id=$3`, role, merchantID, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// RemoveMember removes a user from the merchant.
func (r *Repository) RemoveMember(ctx context.Context, merchantID, userID string) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM merchant_members WHERE merchant_id=$1 AND user_id=$2`, merchantID, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ---- Approvals inbox ----

// CreateApproval records a new approval request (defaults to 1 required approver; set
// requiredApprovals=2 for maker-checker).
func (r *Repository) CreateApproval(ctx context.Context, merchantID, requestedBy string, a *Approval) error {
	a.ID = id.New("apr")
	if a.Currency == "" {
		a.Currency = "ETB"
	}
	if a.Status == "" {
		a.Status = "pending"
	}
	if a.RequiredApprovals < 1 {
		a.RequiredApprovals = 1
	}
	votes, _ := json.Marshal([]ApprovalVote{})
	_, err := r.pool.Exec(ctx, `
		INSERT INTO approval_requests (id, merchant_id, resource_type, resource_id, action, summary, amount, currency, status, requested_by, required_approvals, approvals)
		VALUES ($1,$2,$3,$4,$5,$6,$7::numeric,$8,$9,$10,$11,$12::jsonb)`,
		a.ID, merchantID, a.ResourceType, a.ResourceID, a.Action, a.Summary, a.Amount, a.Currency,
		a.Status, requestedBy, a.RequiredApprovals, string(votes))
	return err
}

// ListApprovals returns approval requests for a merchant, optionally by status.
func (r *Repository) ListApprovals(ctx context.Context, merchantID, status string, limit int) ([]Approval, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query := `SELECT id, resource_type, resource_id, action, COALESCE(summary,''), amount::text, currency, status, required_approvals, approvals::text,
		to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM approval_requests WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if status != "" {
		query += ` AND status=$2`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC LIMIT $` + itoa(len(args)+1)
	args = append(args, limit)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Approval{}
	for rows.Next() {
		var a Approval
		var votes string
		if err := rows.Scan(&a.ID, &a.ResourceType, &a.ResourceID, &a.Action, &a.Summary, &a.Amount,
			&a.Currency, &a.Status, &a.RequiredApprovals, &votes, &a.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(votes), &a.Approvals)
		list = append(list, a)
	}
	return list, rows.Err()
}

// DecideApproval records one approver's decision. Returns the updated approval and a bool
// indicating whether it reached final approval.
func (r *Repository) DecideApproval(ctx context.Context, merchantID, approvalID, userID, userName, userRole, decision string) (*Approval, bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var a Approval
	var votesJSON string
	var required, voteCount int
	err = tx.QueryRow(ctx, `SELECT id, resource_type, resource_id, action, COALESCE(summary,''), amount::text, currency, status, required_approvals, approvals::text,
		(SELECT COUNT(*) FROM jsonb_array_elements(approvals::jsonb)) FROM approval_requests WHERE merchant_id=$1 AND id=$2 FOR UPDATE`,
		merchantID, approvalID).Scan(&a.ID, &a.ResourceType, &a.ResourceID, &a.Action, &a.Summary, &a.Amount,
		&a.Currency, &a.Status, &a.RequiredApprovals, &votesJSON, &voteCount)
	if err != nil {
		return nil, false, ErrNotFound
	}
	if a.Status != "pending" {
		return nil, false, errors.New("approval already decided")
	}
	required = a.RequiredApprovals
	_ = json.Unmarshal([]byte(votesJSON), &a.Approvals)

	// Record the vote (append).
	a.Approvals = append(a.Approvals, ApprovalVote{
		UserID: userID, Name: userName, Role: userRole, Decision: decision, DecidedAt: time.Now().UTC().Format(time.RFC3339),
	})
	newVotes, _ := json.Marshal(a.Approvals)

	voteCount++
	finalDecision := a.Status
	if decision == "reject" {
		finalDecision = "rejected"
	} else if voteCount >= required {
		finalDecision = "approved"
	}
	_, err = tx.Exec(ctx, `UPDATE approval_requests SET approvals=$1::jsonb, status=$2, decided_at=now(), updated_at=now() WHERE id=$3`,
		string(newVotes), finalDecision, approvalID)
	if err != nil {
		return nil, false, err
	}
	a.Status = finalDecision
	_ = tx.Commit(ctx)
	return &a, finalDecision == "approved", nil
}

func itoa(n int) string { return strconv.Itoa(n) }
