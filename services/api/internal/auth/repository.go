package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"apexpay/internal/id"
)

var (
	// ErrNotFound is returned when a user/session/merchant does not exist.
	ErrNotFound = errors.New("not found")
	// ErrInvalidCredentials is returned on a bad email/password.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Repository persists users, merchant memberships, and sessions.
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

type userRow struct {
	ID           string
	Email        string
	Name         string
	Status       string
	PasswordHash string
}

// findUserByEmail returns the user row (including password hash for verification).
func (r *Repository) findUserByEmail(ctx context.Context, email string) (*userRow, error) {
	var u userRow
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, name, status, COALESCE(password_hash,'') FROM users WHERE lower(email)=lower($1)`,
		email).Scan(&u.ID, &u.Email, &u.Name, &u.Status, &u.PasswordHash)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// memberships returns all merchants a user belongs to, with roles.
func (r *Repository) memberships(ctx context.Context, userID string) ([]MerchantContext, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT mm.merchant_id, m.legal_name, mm.role
		FROM merchant_members mm JOIN merchants m ON m.id = mm.merchant_id
		WHERE mm.user_id = $1 AND m.status = 'active'
		ORDER BY m.legal_name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []MerchantContext{}
	for rows.Next() {
		var mc MerchantContext
		if err := rows.Scan(&mc.MerchantID, &mc.LegalName, &mc.Role); err != nil {
			return nil, err
		}
		list = append(list, mc)
	}
	return list, rows.Err()
}

// createSession inserts a session row for the given user+merchant and returns its ID.
func (r *Repository) createSession(ctx context.Context, userID, merchantID, userAgent, ip string, expiresAt time.Time) (string, error) {
	sessionID := id.New("sess")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO auth_sessions (id, user_id, merchant_id, token_hash, user_agent, ip, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6::inet,$7)`,
		sessionID, userID, merchantID, "", userAgent, ip, expiresAt)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

// storeTokenHash records the token digest for an existing session row.
func (r *Repository) storeTokenHash(ctx context.Context, sessionID, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `UPDATE auth_sessions SET token_hash=$1 WHERE id=$2`, tokenHash, sessionID)
	return err
}

// findSessionByTokenHash returns the live session for a token digest.
func (r *Repository) findSessionByTokenHash(ctx context.Context, tokenHash string) (*Session, error) {
	var s Session
	var role string
	err := r.pool.QueryRow(ctx, `
		SELECT s.id, s.user_id, s.merchant_id, mm.role
		FROM auth_sessions s
		JOIN merchant_members mm ON mm.merchant_id = s.merchant_id AND mm.user_id = s.user_id
		WHERE s.token_hash = $1 AND s.revoked_at IS NULL AND s.expires_at > now()`,
		tokenHash).Scan(&s.ID, &s.UserID, &s.MerchantID, &role)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	s.Role = role
	s.TokenHash = tokenHash
	return &s, nil
}

// touchSession updates last_active_at for a validated session.
func (r *Repository) touchSession(ctx context.Context, sessionID string) {
	_, _ = r.pool.Exec(ctx, `UPDATE auth_sessions SET last_active_at=now() WHERE id=$1`, sessionID)
}

// revokeSession invalidates a session by token digest.
func (r *Repository) revokeSession(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `UPDATE auth_sessions SET revoked_at=now() WHERE token_hash=$1 AND revoked_at IS NULL`, tokenHash)
	return err
}

// generateToken returns a cryptographically random opaque session token.
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashToken returns the SHA-256 digest of a session token (only this is persisted).
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
