//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"apexpay/internal/auth"
	"apexpay/internal/id"
	"apexpay/internal/platform/crypto"
)

// Exercises the dashboard session-auth layer end-to-end against a real DB:
// password hashing round-trip, login, session validation, logout, revocation.

func setupAuth(t *testing.T) (*auth.Service, *pgxpool.Pool) {
	pool := setupPool(t)
	svc := auth.NewService(auth.NewRepository(pool))
	return svc, pool
}

func TestAuthLoginValidateLogout(t *testing.T) {
	svc, pool := setupAuth(t)
	defer pool.Close()
	ctx := context.Background()

	// Seed a merchant + user + membership.
	merchantID := id.NewMerchant()
	_, err := pool.Exec(ctx, `INSERT INTO merchants (id, legal_name, display_name, email, status, onboarding_status)
		VALUES ($1,'Auth Test PLC','Auth Test',$2,'active','approved')`,
		merchantID, fmt.Sprintf("auth_%s@example.et", merchantID))
	require.NoError(t, err)

	userID := id.New("user")
	hash, err := crypto.HashPassword("s3cret-Pass!")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO users (id, email, name, status, password_hash, email_verified)
		VALUES ($1,$2,'Auth User','active',$3,true)`, userID, "authuser@example.et", hash)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO merchant_members (merchant_id, user_id, role)
		VALUES ($1,$2,'owner')`, merchantID, userID)
	require.NoError(t, err)

	// Wrong password rejected.
	_, err = svc.Login(ctx, "authuser@example.et", "wrong", "", "")
	assert.ErrorIs(t, err, auth.ErrInvalidCredentials)

	// Correct login returns a token.
	res, err := svc.Login(ctx, "authuser@example.et", "s3cret-Pass!", "go-test", "127.0.0.1")
	require.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, merchantID, res.Merchant.MerchantID)
	assert.Equal(t, "owner", res.Merchant.Role)

	// Session validates.
	sess, err := svc.Validate(ctx, res.Token)
	require.NoError(t, err)
	assert.Equal(t, userID, sess.UserID)
	assert.Equal(t, merchantID, sess.MerchantID)

	// Logout revokes → subsequent validate fails.
	require.NoError(t, svc.Logout(ctx, res.Token))
	_, err = svc.Validate(ctx, res.Token)
	assert.ErrorIs(t, err, auth.ErrNotFound)
}
