package assistant

import (
	"context"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists assistant threads and messages, and resolves the actor scope.
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// CreateThread inserts a new conversation and returns it.
func (r *Repository) CreateThread(ctx context.Context, t *Thread) error {
	t.ID = id.New("athr")
	t.CreatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		INSERT INTO assistant_threads (id, user_id, merchant_id, actor, title)
		VALUES ($1,$2,$3,$4,$5)`, t.ID, t.UserID, t.MerchantID, string(t.Actor), t.Title)
	return err
}

// AppendMessage records one turn. Messages are append-only.
func (r *Repository) AppendMessage(ctx context.Context, m *Message) error {
	m.ID = id.New("amsg")
	m.CreatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
		INSERT INTO assistant_messages (id, thread_id, role, content, intent, tools_used, data)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		m.ID, m.ThreadID, m.Role, m.Content, m.Intent, joinTools(m.ToolsUsed), m.Data)
	return err
}

// ListMessages returns a thread's messages oldest-first.
func (r *Repository) ListMessages(ctx context.Context, threadID string) ([]Message, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, thread_id, role, content, COALESCE(intent,''), COALESCE(tools_used,'{}'), COALESCE(data,''), created_at
		FROM assistant_messages WHERE thread_id=$1 ORDER BY created_at ASC`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Message{}
	for rows.Next() {
		var m Message
		var tools string
		if err := rows.Scan(&m.ID, &m.ThreadID, &m.Role, &m.Content, &m.Intent, &tools, &m.Data, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.ToolsUsed = splitTools(tools)
		out = append(out, m)
	}
	return out, rows.Err()
}

// GetThread returns a thread only if it belongs to the given user (ownership gate).
func (r *Repository) GetThread(ctx context.Context, threadID, userID string) (*Thread, error) {
	var t Thread
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, merchant_id, actor, title, created_at
		FROM assistant_threads WHERE id=$1 AND user_id=$2`, threadID, userID).
		Scan(&t.ID, &t.UserID, &t.MerchantID, &t.Actor, &t.Title, &t.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t.Actor = ActorType(t.Actor)
	return &t, nil
}

// EmployeeIDForUser resolves the employee row for a dashboard user (matched by email to the
// users table), returning ("", false) when the user is not an employee (i.e. merchant actor).
func (r *Repository) EmployeeIDForUser(ctx context.Context, merchantID, userID string) (string, bool, error) {
	var empID string
	err := r.pool.QueryRow(ctx, `
		SELECT e.id FROM employees e
		JOIN users u ON u.email = e.email
		WHERE e.merchant_id=$1 AND u.id=$2 AND e.status='active'
		LIMIT 1`, merchantID, userID).Scan(&empID)
	if err == pgx.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return empID, true, nil
}

func joinTools(tools []string) string {
	if len(tools) == 0 {
		return "[]"
	}
	return "[\"" + joinStrings(tools, "\",\"") + "\"]"
}

func splitTools(s string) []string {
	// Stored as a JSON string array; parse simply.
	if len(s) < 2 {
		return []string{}
	}
	body := s[1 : len(s)-1]
	if body == "" {
		return []string{}
	}
	return splitString(body, "\",\"")
}

func joinStrings(a []string, sep string) string {
	out := ""
	for i, s := range a {
		if i > 0 {
			out += sep
		}
		out += s
	}
	return out
}

func splitString(s, sep string) []string {
	var out []string
	cur := ""
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			out = append(out, cur)
			cur = ""
			i += len(sep) - 1
		} else {
			cur += string(s[i])
		}
	}
	out = append(out, cur)
	return out
}
