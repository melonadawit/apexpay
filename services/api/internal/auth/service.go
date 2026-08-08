package auth

import (
	"context"
	"time"

	"apexpay/internal/platform/crypto"
)

// Service orchestrates dashboard login, session validation, and logout.
type Service struct {
	repo       *Repository
	sessionTTL time.Duration
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo, sessionTTL: 12 * time.Hour}
}

// Login verifies credentials and issues an opaque session token bound to the user's
// first active merchant. Returns ErrInvalidCredentials on a bad email/password.
func (s *Service) Login(ctx context.Context, email, password, userAgent, ip string) (*LoginResult, error) {
	u, err := s.repo.findUserByEmail(ctx, email)
	if err == ErrNotFound {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if u.Status != "active" {
		return nil, ErrInvalidCredentials
	}
	// Always run a verify to keep timing constant-ish even for unknown users is not
	// possible without a stored hash; this path only reaches here when the user exists.
	if u.PasswordHash == "" || !crypto.VerifyPassword(password, u.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	members, err := s.repo.memberships(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, ErrInvalidCredentials // no active merchant membership
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(s.sessionTTL)
	sessionID, err := s.repo.createSession(ctx, u.ID, members[0].MerchantID, userAgent, ip, expiresAt)
	if err != nil {
		return nil, err
	}
	if err := s.repo.storeTokenHash(ctx, sessionID, hashToken(token)); err != nil {
		return nil, err
	}

	return &LoginResult{
		Token:     token,
		ExpiresAt: expiresAt.UTC().Format(time.RFC3339),
		User:      User{ID: u.ID, Email: u.Email, Name: u.Name, Status: u.Status},
		Merchant:  members[0],
	}, nil
}

// Validate resolves a session from its raw token. Returns ErrNotFound if invalid/expired.
func (s *Service) Validate(ctx context.Context, token string) (*Session, error) {
	if token == "" {
		return nil, ErrNotFound
	}
	sess, err := s.repo.findSessionByTokenHash(ctx, hashToken(token))
	if err != nil {
		return nil, err
	}
	s.repo.touchSession(ctx, sess.ID)
	return sess, nil
}

// Logout revokes the session associated with the raw token.
func (s *Service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.repo.revokeSession(ctx, hashToken(token))
}

// Me returns the current user and all their merchant memberships.
func (s *Service) Me(ctx context.Context, userID string) (*SessionResult, error) {
	var res SessionResult
	// Reuse memberships for the merchant list.
	all, err := s.repo.memberships(ctx, userID)
	if err != nil {
		return nil, err
	}
	res.Merchants = all
	if len(all) > 0 {
		res.Merchant = all[0]
	}
	return &res, nil
}
