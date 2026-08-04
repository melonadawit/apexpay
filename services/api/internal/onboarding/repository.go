package onboarding

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"apexpay/internal/id"
	"time"
)

type PgRepository struct{ pool *pgxpool.Pool }

func NewPgRepository(pool *pgxpool.Pool) *PgRepository { return &PgRepository{pool: pool} }

func (r *PgRepository) CreateKYCProfile(ctx context.Context, p *KYCProfile) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO merchant_kyc_profiles (id, merchant_id, version, legal_name, trade_name, business_type, registration_number, tin_number, vat_number, business_license_no, industry_category, business_description, website_url, expected_monthly_tpv, avg_ticket_amount, region, city, sub_city, woreda, office_address_full, contact_person_name, contact_person_role, contact_email, contact_phone, onboarding_status, kyc_level, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26, now(), now())`,
		p.ID, p.MerchantID, p.Version, p.LegalName, p.TradeName, p.BusinessType, p.RegistrationNumber, p.TINNumber, p.VATNumber, p.BusinessLicenseNo, p.Industry, p.Description, p.WebsiteURL, p.ExpectedMonthlyTPV.String(), p.AvgTicketAmount.String(), p.Region, p.City, p.SubCity, p.Woreda, p.AddressFull, p.ContactPersonName, p.ContactPersonRole, p.ContactEmail, p.ContactPhone, p.OnboardingStatus, p.KYCLevel,
	)
	return err
}

func (r *PgRepository) GetKYCProfile(ctx context.Context, merchantID, profileID string) (*KYCProfile, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, version, legal_name, trade_name, business_type, registration_number, tin_number, industry_category, business_description, region, city, office_address_full, contact_person_name, onboarding_status, kyc_level, created_at FROM merchant_kyc_profiles WHERE merchant_id=$1 AND id=$2`, merchantID, profileID)
	var p KYCProfile
	err := row.Scan(&p.ID, &p.MerchantID, &p.Version, &p.LegalName, &p.TradeName, &p.BusinessType, &p.RegistrationNumber, &p.TINNumber, &p.Industry, &p.Description, &p.Region, &p.City, &p.AddressFull, &p.ContactPersonName, &p.OnboardingStatus, &p.KYCLevel, &p.CreatedAt)
	return &p, err
}

func (r *PgRepository) GetLatestKYCProfile(ctx context.Context, merchantID string) (*KYCProfile, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, version, legal_name, business_type, onboarding_status, kyc_level FROM merchant_kyc_profiles WHERE merchant_id=$1 ORDER BY version DESC LIMIT 1`, merchantID)
	var p KYCProfile
	err := row.Scan(&p.ID, &p.MerchantID, &p.Version, &p.LegalName, &p.BusinessType, &p.OnboardingStatus, &p.KYCLevel)
	return &p, err
}

func (r *PgRepository) UpdateKYCStatus(ctx context.Context, merchantID, profileID string, status OnboardingStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE merchant_kyc_profiles SET onboarding_status=$1, updated_at=now() WHERE merchant_id=$2 AND id=$3`, status, merchantID, profileID)
	return err
}

func (r *PgRepository) CreateOwner(ctx context.Context, o *BeneficialOwner) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO merchant_beneficial_owners (id, merchant_id, kyc_profile_id, full_name, full_name_am, role, ownership_percentage, nationality, id_type, phone, email, is_pep, is_authorized_signatory, verification_status, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14, now())`,
		o.ID, o.MerchantID, o.KYCProfileID, o.FullName, o.FullNameAM, o.Role, o.OwnershipPercentage.String(), o.Nationality, o.IDType, o.Phone, o.Email, o.IsPEP, o.IsAuthorizedSignatory, o.VerificationStatus)
	return err
}

func (r *PgRepository) ListOwners(ctx context.Context, merchantID, kycProfileID string) ([]BeneficialOwner, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, kyc_profile_id, full_name, role, ownership_percentage::text, phone, is_authorized_signatory, fayda_verified, verification_status FROM merchant_beneficial_owners WHERE merchant_id=$1 AND kyc_profile_id=$2`, merchantID, kycProfileID)
	if err != nil { return nil, err }
	defer rows.Close()
	var owners []BeneficialOwner
	for rows.Next() {
		var o BeneficialOwner
		var pctStr string
		if err := rows.Scan(&o.ID, &o.MerchantID, &o.KYCProfileID, &o.FullName, &o.Role, &pctStr, &o.Phone, &o.IsAuthorizedSignatory, &o.FaydaVerified, &o.VerificationStatus); err != nil { return nil, err }
		owners = append(owners, o)
	}
	return owners, nil
}

func (r *PgRepository) UpdateOwnerFaydaVerified(ctx context.Context, ownerID string, finHash, finLast4 string, verified bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE merchant_beneficial_owners SET fayda_fin_hash=$1, id_number_last4=$2, fayda_verified=$3, verification_status='fayda_verified', updated_at=now() WHERE id=$4`, finHash, finLast4, verified, ownerID)
	return err
}

func (r *PgRepository) CreateBankAccount(ctx context.Context, b *BankAccount) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO bank_accounts (id, merchant_id, account_name, account_number_masked, account_number_hash, bank_code, bank_name, is_settlement_default, verification_status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		b.ID, b.MerchantID, b.AccountName, b.AccountNumberMasked, b.AccountNumberHash, b.BankCode, b.BankName, b.IsSettlementDefault, b.VerificationStatus)
	return err
}

func (r *PgRepository) ListBankAccounts(ctx context.Context, merchantID string) ([]BankAccount, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, account_name, account_number_masked, account_number_hash, bank_code, bank_name, is_settlement_default, verification_status FROM bank_accounts WHERE merchant_id=$1`, merchantID)
	if err != nil { return nil, err }
	defer rows.Close()
	var accs []BankAccount
	for rows.Next() {
		var b BankAccount
		if err := rows.Scan(&b.ID, &b.MerchantID, &b.AccountName, &b.AccountNumberMasked, &b.AccountNumberHash, &b.BankCode, &b.BankName, &b.IsSettlementDefault, &b.VerificationStatus); err != nil { return nil, err }
		accs = append(accs, b)
	}
	return accs, nil
}

func (r *PgRepository) CreateDocument(ctx context.Context, d *Document) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO merchant_documents (id, merchant_id, kyc_profile_id, owner_id, doc_type, file_key, file_hash, mime_type, file_size_bytes, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		d.ID, d.MerchantID, d.KYCProfileID, d.OwnerID, d.Type, d.FileKey, d.FileHash, d.MimeType, d.SizeBytes, d.Status)
	return err
}

func (r *PgRepository) ListDocuments(ctx context.Context, merchantID, kycProfileID string) ([]Document, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, kyc_profile_id, doc_type, file_key, file_hash, mime_type, file_size_bytes, status FROM merchant_documents WHERE merchant_id=$1 AND kyc_profile_id=$2`, merchantID, kycProfileID)
	if err != nil { return nil, err }
	defer rows.Close()
	var docs []Document
	for rows.Next() {
		var d Document
		if err := rows.Scan(&d.ID, &d.MerchantID, &d.KYCProfileID, &d.Type, &d.FileKey, &d.FileHash, &d.MimeType, &d.SizeBytes, &d.Status); err != nil { return nil, err }
		docs = append(docs, d)
	}
	return docs, nil
}

func (r *PgRepository) GetDocument(ctx context.Context, docID string) (*Document, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, kyc_profile_id, doc_type, file_key, file_hash, mime_type, file_size_bytes, status FROM merchant_documents WHERE id=$1`, docID)
	var d Document
	err := row.Scan(&d.ID, &d.MerchantID, &d.KYCProfileID, &d.Type, &d.FileKey, &d.FileHash, &d.MimeType, &d.SizeBytes, &d.Status)
	return &d, err
}

func (r *PgRepository) CreateComplianceCheck(ctx context.Context, c *ComplianceCheck) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO compliance_checks (id, merchant_id, kyc_profile_id, check_type, status, score, provider, details) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		c.ID, c.MerchantID, c.KYCProfileID, c.Type, c.Status, c.Score, c.Provider, c.Details)
	return err
}

func (r *PgRepository) ListComplianceChecks(ctx context.Context, merchantID, kycProfileID string) ([]ComplianceCheck, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, kyc_profile_id, check_type, status, score FROM compliance_checks WHERE merchant_id=$1 AND kyc_profile_id=$2`, merchantID, kycProfileID)
	if err != nil { return nil, err }
	defer rows.Close()
	var checks []ComplianceCheck
	for rows.Next() {
		var c ComplianceCheck
		if err := rows.Scan(&c.ID, &c.MerchantID, &c.KYCProfileID, &c.Type, &c.Status, &c.Score); err != nil { return nil, err }
		checks = append(checks, c)
	}
	return checks, nil
}

func (r *PgRepository) CreateReview(ctx context.Context, rev *OnboardingReview) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO onboarding_reviews (id, merchant_id, kyc_profile_id, reviewer_id, reviewer_type, from_status, to_status, action, comments, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9, now())`,
		rev.ID, rev.MerchantID, rev.KYCProfileID, rev.ReviewerID, rev.ReviewerType, rev.FromStatus, rev.ToStatus, rev.Action, rev.Comments)
	return err
}

func (r *PgRepository) ListReviews(ctx context.Context, merchantID string) ([]OnboardingReview, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, kyc_profile_id, reviewer_type, from_status, to_status, action, comments, created_at FROM onboarding_reviews WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil { return nil, err }
	defer rows.Close()
	var revs []OnboardingReview
	for rows.Next() {
		var rv OnboardingReview
		if err := rows.Scan(&rv.ID, &rv.MerchantID, &rv.KYCProfileID, &rv.ReviewerType, &rv.FromStatus, &rv.ToStatus, &rv.Action, &rv.Comments, &rv.CreatedAt); err != nil { return nil, err }
		revs = append(revs, rv)
	}
	return revs, nil
}

func (r *PgRepository) ApproveMerchantTx(ctx context.Context, merchantID, kycProfileID, reviewerID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil { return err }
	defer func() { _ = tx.Rollback(ctx) }()

	// Update KYC approved
	_, err = tx.Exec(ctx, `UPDATE merchant_kyc_profiles SET onboarding_status='approved', reviewed_at=now(), updated_at=now() WHERE merchant_id=$1 AND id=$2`, merchantID, kycProfileID)
	if err != nil { return err }

	// Update merchant active + operating book creation
	_, err = tx.Exec(ctx, `UPDATE merchants SET status='active', onboarding_status='active', fayda_verified=true, updated_at=now() WHERE id=$1`, merchantID)
	if err != nil { return err }

	// Create operating book - idempotent
	bookID := id.NewLedgerBook()
	_, err = tx.Exec(ctx, `INSERT INTO ledger_books (id, merchant_id, book_type, name, currency, status) VALUES ($1,$2,'merchant_operating',$3,'ETB','open') ON CONFLICT (id) DO NOTHING`, bookID, merchantID, fmt.Sprintf("Operating book %s", merchantID))
	if err != nil { return err }

	// Seed standard accounts - optimal batch
	accounts := [][]string{
		{id.New("la"), bookID, "asset:clearing:mock", "Clearing Mock", "debit"},
		{id.New("la"), bookID, "asset:clearing:bank", "Clearing Bank", "debit"},
		{id.New("la"), bookID, "liability:merchant_payable", "Merchant Payable", "credit"},
		{id.New("la"), bookID, "liability:platform_fee_due", "Platform Fee Due", "credit"},
		{id.New("la"), bookID, "liability:payroll_payable", "Payroll Payable", "credit"},
		{id.New("la"), bookID, "expense:salary", "Salary Expense", "debit"},
	}
	for _, acc := range accounts {
		_, err = tx.Exec(ctx, `INSERT INTO ledger_accounts (id, book_id, code, name, normal_balance) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (book_id, code) DO NOTHING`, acc[0], acc[1], acc[2], acc[3], acc[4])
		if err != nil { return err }
	}

	// Audit review
	_, err = tx.Exec(ctx, `INSERT INTO onboarding_reviews (id, merchant_id, kyc_profile_id, reviewer_id, reviewer_type, from_status, to_status, action, comments, created_at) VALUES ($1,$2,$3,$4,'compliance','compliance_check','active','approve','approved merchant active + operating book created', now())`, id.New("orev"), merchantID, kycProfileID, reviewerID)
	if err != nil { return err }

	// Outbox merchant.activated
	_, err = tx.Exec(ctx, `INSERT INTO outbox_events (id, merchant_id, aggregate_type, aggregate_id, event_type, payload) VALUES ($1,$2,'merchant',$3,'merchant.activated',$4)`, id.NewOutbox(), merchantID, merchantID, fmt.Sprintf(`{"merchant_id":"%s","kyc_profile_id":"%s"}`, merchantID, kycProfileID))
	if err != nil { return err }

	// Update merchants.kyc_profile_id
	_, err = tx.Exec(ctx, `UPDATE merchants SET kyc_profile_id=$1 WHERE id=$2`, kycProfileID, merchantID)
	if err != nil { return err }

	return tx.Commit(ctx)
}

// Helper for audit
func init() { _ = time.Now() }
