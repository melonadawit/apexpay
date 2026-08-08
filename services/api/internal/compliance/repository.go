package compliance

import (
	"context"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// GetStatus returns the merchant's compliance console status, deriving it from KYC,
// documents, and obligations if no console row exists yet.
func (r *Repository) GetStatus(ctx context.Context, merchantID string) (*Status, error) {
	// Derive from merchants + kyc profiles; if a console row exists it refines it.
	var s Status
	s.MerchantID = merchantID
	err := r.pool.QueryRow(ctx, `
		SELECT m.onboarding_status, COALESCE(k.license_expiry,'')::text, m.fayda_verified, COALESCE(m.risk_tier,'low'),
		       COALESCE(to_char(c.kyc_expiry_date,'YYYY-MM-DD'),''),
		       COALESCE(to_char(c.next_erca_due,'YYYY-MM-DD'),''),
		       COALESCE(to_char(c.next_pension_due,'YYYY-MM-DD'),''),
		       COALESCE(to_char(c.annual_tax_filing_due,'YYYY-MM-DD'),''),
		       COALESCE(to_char(c.aml_due,'YYYY-MM-DD'),''),
		       COALESCE(c.overall_status,'attention'), COALESCE(c.notes,'')
		FROM merchants m
		LEFT JOIN merchant_kyc_profiles k ON k.merchant_id = m.id AND k.version = (SELECT MAX(version) FROM merchant_kyc_profiles WHERE merchant_id=m.id)
		LEFT JOIN compliance_console c ON c.merchant_id = m.id
		WHERE m.id=$1`, merchantID).
		Scan(&s.OnboardingStatus, &s.LicenseExpiry, &s.FaydaVerified, &s.RiskTier,
			&s.KYCExpiryDate, &s.NextERCADue, &s.NextPensionDue, &s.AnnualTaxFilingDue, &s.AMLDue,
			&s.OverallStatus, &s.Notes)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// UpsertStatus saves/refreshes the compliance console row.
func (r *Repository) UpsertStatus(ctx context.Context, merchantID string, s *Status) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO compliance_console (merchant_id, onboarding_status, kyc_expiry_date, license_expiry, fayda_verified, risk_tier,
			next_erca_due, next_pension_due, annual_tax_filing_due, aml_due, overall_status, notes, updated_at)
		VALUES ($1,$2,$3::date,$4::date,$5,$6,$7::date,$8::date,$9::date,$10::date,$11,$12,now())
		ON CONFLICT (merchant_id) DO UPDATE SET
			onboarding_status=EXCLUDED.onboarding_status, kyc_expiry_date=EXCLUDED.kyc_expiry_date,
			license_expiry=EXCLUDED.license_expiry, fayda_verified=EXCLUDED.fayda_verified, risk_tier=EXCLUDED.risk_tier,
			next_erca_due=EXCLUDED.next_erca_due, next_pension_due=EXCLUDED.next_pension_due,
			annual_tax_filing_due=EXCLUDED.annual_tax_filing_due, aml_due=EXCLUDED.aml_due,
			overall_status=EXCLUDED.overall_status, notes=EXCLUDED.notes, updated_at=now()`,
		merchantID, s.OnboardingStatus, s.KYCExpiryDate, s.LicenseExpiry, s.FaydaVerified, s.RiskTier,
		s.NextERCADue, s.NextPensionDue, s.AnnualTaxFilingDue, s.AMLDue, s.OverallStatus, s.Notes)
	return err
}

// ListChecks returns recent compliance check results.
func (r *Repository) ListChecks(ctx context.Context, merchantID string, limit int) ([]CheckLog, error) {
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, check_type, status, COALESCE(detail,''),
		       to_char(checked_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM compliance_checks_log WHERE merchant_id=$1 ORDER BY checked_at DESC LIMIT $2`, merchantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []CheckLog{}
	for rows.Next() {
		var c CheckLog
		if err := rows.Scan(&c.ID, &c.CheckType, &c.Status, &c.Detail, &c.CheckedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// LogCheck records a compliance check result.
func (r *Repository) LogCheck(ctx context.Context, merchantID, checkType, status, detail string) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO compliance_checks_log (id, merchant_id, check_type, status, detail)
		VALUES ($1,$2,$3,$4,$5)`, id.New("cchk"), merchantID, checkType, status, detail)
	return err
}

// ComputeOverall derives the overall status (good/attention/overdue) from due dates.
func ComputeOverall(now time.Time, kyc, license, erca, pension, aml string) string {
	overdue := 0
	for _, d := range []string{kyc, license, erca, pension, aml} {
		if d == "" {
			continue
		}
		if t, err := time.Parse("2006-01-02", d); err == nil && t.Before(now) {
			overdue++
		}
	}
	if overdue >= 2 {
		return "overdue"
	}
	if overdue >= 1 {
		return "attention"
	}
	return "good"
}
