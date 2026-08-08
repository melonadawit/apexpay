//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"apexpay/internal/admin"
	"apexpay/internal/id"
	"apexpay/internal/onboarding"
)

// Exercises the DB-backed admin package: onboarding queue, maker-checker dual
// approval for high-risk merchants, operating-book provisioning on final approval,
// the compliance exam, and reconciliation evidence reads.

func setupAdminRepo(t *testing.T) (*admin.Repository, *pgxpool.Pool) {
	pool := setupPool(t)
	repo := admin.NewRepository(pool)
	return repo, pool
}

// seedHighRiskMerchant creates a merchant + KYC whose risk score is high enough to
// trigger maker-checker (risk >= 70), returning its IDs.
func seedHighRiskMerchant(t *testing.T, pool *pgxpool.Pool) (merchantID, kycID string) {
	ctx := context.Background()
	merchantID = id.NewMerchant()
	_, err := pool.Exec(ctx, `INSERT INTO merchants (id, legal_name, display_name, email, status, onboarding_status, risk_score, risk_tier)
		VALUES ($1,$2,$3,$4,'pending','submitted',85,'high')`,
		merchantID, "High Risk PLC", "High Risk", fmt.Sprintf("highrisk_%s@example.et", merchantID))
	require.NoError(t, err)

	kycID = id.NewKYCProfile()
	kyc := &onboarding.KYCProfile{
		ID: kycID, MerchantID: merchantID, Version: 1,
		LegalName: "High Risk PLC", BusinessType: onboarding.BusinessTypePLC,
		RegistrationNumber: "MT/AA/999", TINNumber: "0023456789", Industry: onboarding.IndustryEcommerce,
		Description: "high risk", Region: "Addis Ababa", City: "Addis Ababa",
		AddressFull: "Bole", ContactPersonName: "A", ContactPersonRole: "owner",
		ContactEmail: "a@example.et", ContactPhone: "0911111111",
		ExpectedMonthlyTPV: decimal.NewFromInt(2000000), OnboardingStatus: onboarding.StatusSubmitted,
	}
	err = onboarding.NewPgRepository(pool).CreateKYCProfile(ctx, kyc)
	require.NoError(t, err)
	return merchantID, kycID
}

func TestAdminReviewHighRiskMakerChecker(t *testing.T) {
	repo, pool := setupAdminRepo(t)
	defer pool.Close()
	ctx := context.Background()

	merchantID, _ := seedHighRiskMerchant(t, pool)

	// First reviewer approves → high-risk, must NOT reach terminal approval yet.
	res, err := repo.Review(ctx, merchantID, "ops", "user_reviewer_a", "approve", "")
	require.NoError(t, err)
	assert.Equal(t, "pending_approval", res.Status, "high-risk first approval must be parked for second approver")

	// Different reviewer approves → final approval + operating book provisioned.
	res2, err := repo.Review(ctx, merchantID, "compliance", "user_reviewer_b", "approve", "")
	require.NoError(t, err)
	assert.Equal(t, "approved", res2.Status)

	// Operating book must now exist.
	var bookCount int
	err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM ledger_books WHERE merchant_id=$1 AND book_type='merchant_operating'`, merchantID).Scan(&bookCount)
	require.NoError(t, err)
	assert.Equal(t, 1, bookCount, "operating book provisioned on final approval")

	// Terminal state prevents further review.
	_, err = repo.Review(ctx, merchantID, "ops", "user_reviewer_c", "approve", "")
	assert.Error(t, err, "already in terminal state")

	// The queue should no longer list this merchant (now approved).
	queue, err := repo.ListOnboardingQueue(ctx, 50)
	require.NoError(t, err)
	for _, q := range queue {
		assert.NotEqual(t, merchantID, q.MerchantID, "approved merchant should leave the review queue")
	}
}

func TestAdminReviewLowRiskSingleApproval(t *testing.T) {
	repo, pool := setupAdminRepo(t)
	defer pool.Close()
	ctx := context.Background()

	// A low-risk merchant (risk < 70, low TPV) approves in a single pass.
	merchantID := id.NewMerchant()
	_, err := pool.Exec(ctx, `INSERT INTO merchants (id, legal_name, display_name, email, status, onboarding_status, risk_score, risk_tier)
		VALUES ($1,'Low Risk','Low Risk',$2,'pending','submitted',20,'low')`,
		merchantID, fmt.Sprintf("lowrisk_%s@example.et", merchantID))
	require.NoError(t, err)
	kycID := id.NewKYCProfile()
	err = onboarding.NewPgRepository(pool).CreateKYCProfile(ctx, &onboarding.KYCProfile{
		ID: kycID, MerchantID: merchantID, Version: 1, LegalName: "Low Risk",
		BusinessType: onboarding.BusinessTypePLC, RegistrationNumber: "MT/AA/1", TINNumber: "0023456789",
		Industry: onboarding.IndustryEcommerce, Description: "low", Region: "AA", City: "AA",
		AddressFull: "x", ContactPersonName: "A", ContactPersonRole: "owner", ContactEmail: "a@example.et",
		ContactPhone: "0911111111", ExpectedMonthlyTPV: decimal.NewFromInt(100000), OnboardingStatus: onboarding.StatusSubmitted,
	})
	require.NoError(t, err)

	res, err := repo.Review(ctx, merchantID, "ops", "user_low", "approve", "")
	require.NoError(t, err)
	assert.Equal(t, "approved", res.Status, "low-risk merchant approves on first pass")
}

func TestAdminEvidenceAndRecon(t *testing.T) {
	repo, pool := setupAdminRepo(t)
	defer pool.Close()
	ctx := context.Background()

	merchantID, kycID := seedHighRiskMerchant(t, pool)

	// Insert a doc + an open reconciliation case for evidence/recon endpoints.
	_, err := pool.Exec(ctx, `INSERT INTO merchant_documents (id, merchant_id, kyc_profile_id, doc_type, file_key, file_hash, mime_type, file_size_bytes, status)
		VALUES ($1,$2,$3,'company_registration','x.pdf','hash_x','application/pdf',1024,'verified')`,
		id.NewDocument(), merchantID, kycID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO payment_reconciliation_cases (merchant_id, idempotency_key, tx_ref, status)
		VALUES ($1,'idem_ev','txr_ev','open') ON CONFLICT (merchant_id, idempotency_key) DO NOTHING`, merchantID)
	require.NoError(t, err)

	breaks, err := repo.ReconBreaks(ctx, 50)
	require.NoError(t, err)
	found := false
	for _, b := range breaks {
		if b.TxRef == "txr_ev" {
			found = true
			assert.Equal(t, "open", b.Status)
		}
	}
	assert.True(t, found, "open reconciliation case listed")

	exam, err := repo.MerchantExam(ctx, merchantID)
	require.NoError(t, err)
	assert.Equal(t, "High Risk PLC", exam.LegalName)
	assert.Len(t, exam.Documents, 1)
	assert.Equal(t, "verified", exam.Documents[0].Status)
}
