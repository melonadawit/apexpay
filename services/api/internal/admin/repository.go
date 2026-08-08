package admin

import (
	"context"
	"errors"
	"fmt"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a resource does not exist.
var ErrNotFound = errors.New("not found")

// Repository provides admin/compliance read models and review operations against Postgres.
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// ListOnboardingQueue returns merchants awaiting a compliance decision, oldest first.
func (r *Repository) ListOnboardingQueue(ctx context.Context, limit int) ([]OnboardingQueueItem, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT m.id, m.legal_name, m.email, kyc.onboarding_status,
		       m.risk_score, COALESCE(m.risk_tier,'low'), m.fayda_verified,
		       COALESCE(to_char(kyc.submitted_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS'),''),
		       to_char(kyc.created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM merchants m
		JOIN merchant_kyc_profiles kyc ON kyc.merchant_id = m.id
		WHERE kyc.onboarding_status IN ('submitted','in_review','fayda_pending','compliance_check','needs_more_info')
		ORDER BY kyc.created_at ASC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []OnboardingQueueItem{}
	for rows.Next() {
		var it OnboardingQueueItem
		if err := rows.Scan(&it.MerchantID, &it.LegalName, &it.Email, &it.OnboardingStatus,
			&it.RiskScore, &it.RiskTier, &it.FaydaVerified, &it.SubmittedAt, &it.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

// Review applies a reviewer action (approve/reject/request_info) to a merchant's KYC.
//
// Maker-checker: an approve on a high-risk merchant (risk_score >= 70) or one with high
// expected monthly TPV (> 1,000,000 ETB) requires a second, distinct approver. The first
// approve parks the record in 'pending_approval'; the second transitions it to 'approved'
// and provisions the operating book. Lower-risk approvals proceed on the first pass.
func (r *Repository) Review(ctx context.Context, merchantID, reviewerType, reviewerID, action, comment string) (ReviewResult, error) {
	switch action {
	case "approve", "reject", "request_info":
	default:
		return ReviewResult{}, fmt.Errorf("invalid action %q", action)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return ReviewResult{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Lock the merchant + current KYC so concurrent reviews serialize.
	var currentKYCStatus, curRiskTier string
	var riskScore int
	var expectedTPV float64
	var kycID string
	err = tx.QueryRow(ctx, `
		SELECT kyc.id, kyc.onboarding_status, m.risk_score, COALESCE(m.risk_tier,'low'),
		       COALESCE(kyc.expected_monthly_tpv,0)
		FROM merchant_kyc_profiles kyc
		JOIN merchants m ON m.id = kyc.merchant_id
		WHERE kyc.merchant_id = $1 AND kyc.version = (SELECT MAX(version) FROM merchant_kyc_profiles WHERE merchant_id = $1)
		FOR UPDATE OF kyc, m`, merchantID).Scan(&kycID, &currentKYCStatus, &riskScore, &curRiskTier, &expectedTPV)
	if err == pgx.ErrNoRows {
		return ReviewResult{}, ErrNotFound
	}
	if err != nil {
		return ReviewResult{}, err
	}

	if currentKYCStatus == "approved" || currentKYCStatus == "rejected" {
		return ReviewResult{}, fmt.Errorf("merchant already in terminal state %q", currentKYCStatus)
	}

	highRisk := riskScore >= 70 || expectedTPV > 1_000_000

	var toStatus, resultStatus string
	switch action {
	case "approve":
		if highRisk {
			// Count prior approvals by a different reviewer.
			var prior int
			if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM onboarding_reviews
				WHERE merchant_id=$1 AND action='approve' AND reviewer_id IS DISTINCT FROM $2`,
				merchantID, reviewerID).Scan(&prior); err != nil {
				return ReviewResult{}, err
			}
			if prior == 0 {
				toStatus = "pending_approval"
				resultStatus = "pending_approval"
			} else {
				toStatus = "approved"
				resultStatus = "approved"
			}
		} else {
			toStatus = "approved"
			resultStatus = "approved"
		}
	case "reject":
		toStatus, resultStatus = "rejected", "rejected"
	case "request_info":
		toStatus, resultStatus = "needs_more_info", "needs_more_info"
	}

	reviewID := id.New("rev")
	reviewerTypeVal := "ops"
	switch reviewerType {
	case "admin", "compliance":
		reviewerTypeVal = reviewerType
	}

	if _, err := tx.Exec(ctx, `INSERT INTO onboarding_reviews
		(id, merchant_id, kyc_profile_id, reviewer_id, reviewer_type, from_status, to_status, action, comments, internal_notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NULL)`,
		reviewID, merchantID, kycID, nilStr(reviewerID), reviewerTypeVal, currentKYCStatus, toStatus, action, comment); err != nil {
		return ReviewResult{}, err
	}

	if _, err := tx.Exec(ctx, `UPDATE merchant_kyc_profiles SET onboarding_status=$1, reviewed_at=now(), updated_at=now() WHERE id=$2`,
		toStatus, kycID); err != nil {
		return ReviewResult{}, err
	}

	// Mirror onto merchant-level onboarding_status (values overlap for terminal states).
	if _, err := tx.Exec(ctx, `UPDATE merchants SET onboarding_status=$1, status=CASE WHEN $2 THEN 'active' ELSE status END, updated_at=now() WHERE id=$3`,
		toStatus, toStatus == "approved", merchantID); err != nil {
		return ReviewResult{}, err
	}

	// On final approval, provision the merchant operating book + core accounts + outbox event.
	if resultStatus == "approved" {
		if err := r.provisionOperatingBook(ctx, tx, merchantID); err != nil {
			return ReviewResult{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return ReviewResult{}, err
	}

	msg := "review recorded"
	if resultStatus == "pending_approval" {
		msg = "high-risk merchant: first approval recorded, second approver required"
	} else if resultStatus == "approved" {
		msg = "merchant active + operating book provisioned"
	}
	return ReviewResult{MerchantID: merchantID, Status: resultStatus, Action: action, Message: msg}, nil
}

// provisionOperatingBook creates the merchant operating ledger book, its core accounts,
// and emits a merchant.activated outbox event, all inside the caller's transaction.
func (r *Repository) provisionOperatingBook(ctx context.Context, tx pgx.Tx, merchantID string) error {
	bookID := id.New("lbk")
	if _, err := tx.Exec(ctx, `INSERT INTO ledger_books (id, merchant_id, book_type, name, currency, status)
		VALUES ($1,$2,'merchant_operating','Merchant Operating Book','ETB','open')`,
		bookID, merchantID); err != nil {
		return err
	}

	// Core operating accounts. Codes are stable for downstream accounting exports.
	accounts := []struct{ code, name, balance string }{
		{"1000", "Merchant Settlement Clearing", "credit"},
		{"1100", "Merchant Revenue", "credit"},
		{"1200", "Platform Fee Receivable", "debit"},
		{"1300", "Withholding Tax Payable", "credit"},
		{"1400", "Escrow Hold", "credit"},
		{"1500", "Suspense / Reconciliation", "debit"},
	}
	for _, a := range accounts {
		if _, err := tx.Exec(ctx, `INSERT INTO ledger_accounts (id, book_id, code, name, normal_balance)
			VALUES ($1,$2,$3,$4,$5)`, id.New("lacct"), bookID, a.code, a.name, a.balance); err != nil {
			return err
		}
	}

	// Emit merchant.activated for the outbox publisher.
	outboxID := id.New("outbox")
	if _, err := tx.Exec(ctx, `INSERT INTO outbox_events (id, merchant_id, event_type, payload)
		VALUES ($1,$2,'merchant.activated', jsonb_build_object('merchant_id',$2::text,'activated_at', now()))`,
		outboxID, merchantID); err != nil {
		return err
	}
	return nil
}

// ConnectorHealth returns the last-5-minute health per connector.
func (r *Repository) ConnectorHealth(ctx context.Context) ([]ConnectorHealth, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT connector_id,
		       COALESCE(AVG(latency_ms)::int, 0),
		       COALESCE(COUNT(*) FILTER (WHERE success)::float / NULLIF(COUNT(*),0)::float, 0),
		       COUNT(*)
		FROM connector_health_samples
		WHERE sampled_at >= now() - interval '5 minutes'
		GROUP BY connector_id
		ORDER BY connector_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	health := []ConnectorHealth{}
	for rows.Next() {
		var h ConnectorHealth
		if err := rows.Scan(&h.Connector, &h.AvgLatencyMS, &h.SuccessRate, &h.Samples); err != nil {
			return nil, err
		}
		// A fully-open circuit is derived from the success window; heuristic for the UI.
		h.Circuit = "closed"
		if h.SuccessRate < 0.7 && h.Samples >= 5 {
			h.Circuit = "open"
		} else if h.SuccessRate < 0.9 && h.Samples >= 3 {
			h.Circuit = "degraded"
		}
		health = append(health, h)
	}
	return health, rows.Err()
}

// ReconBreaks returns open reconciliation cases awaiting a decision.
func (r *Repository) ReconBreaks(ctx context.Context, limit int) ([]ReconBreak, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := r.pool.Query(ctx, `
		SELECT merchant_id, idempotency_key, COALESCE(tx_ref,''), status,
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM payment_reconciliation_cases
		WHERE status = 'open'
		ORDER BY created_at ASC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	breaks := []ReconBreak{}
	for rows.Next() {
		var b ReconBreak
		if err := rows.Scan(&b.MerchantID, &b.IdempotencyKey, &b.TxRef, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		breaks = append(breaks, b)
	}
	return breaks, rows.Err()
}

// Evidence builds the compliance evidence bundle for a transaction reference.
func (r *Repository) Evidence(ctx context.Context, txRef string) (Evidence, error) {
	ev := Evidence{TxRef: txRef, DocHashes: []string{}, LedgerJournals: []LedgerJournal{},
		OnboardingReviews: []ReviewRecord{}, AuditLogs: []AuditLog{}, WebhookDeliveries: []WebhookDelivery{}}

	// Ledger journals that reference this transaction.
	jrows, err := r.pool.Query(ctx, `SELECT id, posting_key, COALESCE(memo,''), COALESCE(reference_type,''), COALESCE(reference_id,''),
		to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM ledger_journals WHERE reference_id=$1 ORDER BY created_at ASC`, txRef)
	if err != nil {
		return ev, err
	}
	for jrows.Next() {
		var j LedgerJournal
		if err := jrows.Scan(&j.ID, &j.PostingKey, &j.Memo, &j.RefType, &j.RefID, &j.CreatedAt); err != nil {
			jrows.Close()
			return ev, err
		}
		ev.LedgerJournals = append(ev.LedgerJournals, j)
	}
	jrows.Close()
	if err := jrows.Err(); err != nil {
		return ev, err
	}

	// Audit log trail.
	arows, err := r.pool.Query(ctx, `SELECT id, actor_type, COALESCE(actor_id,''), action, COALESCE(resource_type,''), COALESCE(resource_id,''),
		to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM audit_logs WHERE resource_id=$1 OR actor_id=$1 ORDER BY created_at ASC`, txRef)
	if err != nil {
		return ev, err
	}
	for arows.Next() {
		var a AuditLog
		if err := arows.Scan(&a.ID, &a.ActorType, &a.ActorID, &a.Action, &a.ResourceType, &a.ResourceID, &a.CreatedAt); err != nil {
			arows.Close()
			return ev, err
		}
		ev.AuditLogs = append(ev.AuditLogs, a)
	}
	arows.Close()
	if err := arows.Err(); err != nil {
		return ev, err
	}

	// Webhook deliveries for the transaction.
	wrows, err := r.pool.Query(ctx, `SELECT id, event_type, status, attempt_count, COALESCE(last_status_code,0), COALESCE(last_error,''),
		to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM webhook_deliveries WHERE outbox_event_id IN (
			SELECT id FROM outbox_events WHERE payload->>'tx_ref'=$1 OR payload->>'payment_id'=$1 OR payload->>'refund_id'=$1)
		ORDER BY created_at ASC`, txRef)
	if err != nil {
		return ev, err
	}
	for wrows.Next() {
		var w WebhookDelivery
		if err := wrows.Scan(&w.ID, &w.EventType, &w.Status, &w.AttemptCount, &w.LastCode, &w.LastError, &w.CreatedAt); err != nil {
			wrows.Close()
			return ev, err
		}
		ev.WebhookDeliveries = append(ev.WebhookDeliveries, w)
	}
	wrows.Close()
	if err := wrows.Err(); err != nil {
		return ev, err
	}

	return ev, nil
}

// MerchantExam assembles the full compliance file for a merchant.
func (r *Repository) MerchantExam(ctx context.Context, merchantID string) (MerchantExam, error) {
	exam := MerchantExam{
		MerchantID: merchantID, KYCProfiles: []KYCProfile{}, Owners: []Owner{},
		Documents: []Document{}, ComplianceChecks: []ComplianceCheck{},
		OnboardingReviews: []ReviewRecord{}, Banks: []BankAccount{}, LedgerBooks: []LedgerBook{},
	}

	err := r.pool.QueryRow(ctx, `SELECT COALESCE(legal_name,''), COALESCE(onboarding_status,''), risk_score, COALESCE(risk_tier,'low')
		FROM merchants WHERE id=$1`, merchantID).Scan(&exam.LegalName, &exam.OnboardingStatus, &exam.RiskScore, &exam.RiskTier)
	if err == pgx.ErrNoRows {
		return exam, ErrNotFound
	}
	if err != nil {
		return exam, err
	}

	// KYC profiles.
	krows, err := r.pool.Query(ctx, `SELECT id, version, legal_name, tin_number, business_type, onboarding_status,
		COALESCE(to_char(submitted_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS'),''),
		COALESCE(to_char(reviewed_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS'),'')
		FROM merchant_kyc_profiles WHERE merchant_id=$1 ORDER BY version ASC`, merchantID)
	if err != nil {
		return exam, err
	}
	for krows.Next() {
		var k KYCProfile
		if err := krows.Scan(&k.ID, &k.Version, &k.LegalName, &k.TINNumber, &k.BusinessType, &k.OnboardingStatus, &k.SubmittedAt, &k.ReviewedAt); err != nil {
			krows.Close()
			return exam, err
		}
		exam.KYCProfiles = append(exam.KYCProfiles, k)
	}
	krows.Close()

	// Beneficial owners (PII-safe).
	orows, err := r.pool.Query(ctx, `SELECT id, full_name, role, fayda_verified, id_type, COALESCE(id_number_last4,''), COALESCE(verification_status,'pending')
		FROM merchant_beneficial_owners WHERE merchant_id=$1 ORDER BY created_at ASC`, merchantID)
	if err != nil {
		return exam, err
	}
	for orows.Next() {
		var o Owner
		if err := orows.Scan(&o.ID, &o.FullName, &o.Role, &o.FaydaVerified, &o.IDType, &o.IDNumberLast4, &o.VerificationStatus); err != nil {
			orows.Close()
			return exam, err
		}
		exam.Owners = append(exam.Owners, o)
	}
	orows.Close()

	// Documents.
	drows, err := r.pool.Query(ctx, `SELECT id, doc_type, status, file_hash,
		to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM merchant_documents WHERE merchant_id=$1 ORDER BY created_at ASC`, merchantID)
	if err != nil {
		return exam, err
	}
	for drows.Next() {
		var d Document
		if err := drows.Scan(&d.ID, &d.DocType, &d.Status, &d.FileHash, &d.CreatedAt); err != nil {
			drows.Close()
			return exam, err
		}
		exam.Documents = append(exam.Documents, d)
	}
	drows.Close()

	// Compliance checks.
	crows, err := r.pool.Query(ctx, `SELECT check_type, status, COALESCE(score,0), COALESCE(provider,''),
		to_char(updated_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM compliance_checks WHERE merchant_id=$1 ORDER BY created_at ASC`, merchantID)
	if err != nil {
		return exam, err
	}
	for crows.Next() {
		var c ComplianceCheck
		if err := crows.Scan(&c.CheckType, &c.Status, &c.Score, &c.Provider, &c.UpdatedAt); err != nil {
			crows.Close()
			return exam, err
		}
		exam.ComplianceChecks = append(exam.ComplianceChecks, c)
	}
	crows.Close()

	// Onboarding review trail.
	rrows, err := r.pool.Query(ctx, `SELECT id, reviewer_type, from_status, to_status, action, COALESCE(comments,''),
		to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM onboarding_reviews WHERE merchant_id=$1 ORDER BY created_at ASC`, merchantID)
	if err != nil {
		return exam, err
	}
	for rrows.Next() {
		var rc ReviewRecord
		if err := rrows.Scan(&rc.ID, &rc.ReviewerType, &rc.FromStatus, &rc.ToStatus, &rc.Action, &rc.Comments, &rc.CreatedAt); err != nil {
			rrows.Close()
			return exam, err
		}
		exam.OnboardingReviews = append(exam.OnboardingReviews, rc)
	}
	rrows.Close()

	// Bank accounts (masked).
	brows, err := r.pool.Query(ctx, `SELECT id, COALESCE(bank_code,''), account_number_masked, account_name,
		(verification_status = 'verified') FROM bank_accounts WHERE merchant_id=$1 ORDER BY created_at ASC`, merchantID)
	if err != nil {
		return exam, err
	}
	for brows.Next() {
		var ba BankAccount
		var acct, name string
		var verified bool
		if err := brows.Scan(&ba.ID, &ba.BankCode, &acct, &name, &verified); err != nil {
			brows.Close()
			return exam, err
		}
		ba.AccountNumberMasked = maskAccount(acct)
		ba.AccountName = name
		ba.IsVerified = verified
		exam.Banks = append(exam.Banks, ba)
	}
	brows.Close()

	// Ledger books.
	lrows, err := r.pool.Query(ctx, `SELECT id, book_type, name, currency, status
		FROM ledger_books WHERE merchant_id=$1 ORDER BY created_at ASC`, merchantID)
	if err != nil {
		return exam, err
	}
	for lrows.Next() {
		var lb LedgerBook
		if err := lrows.Scan(&lb.ID, &lb.BookType, &lb.Name, &lb.Currency, &lb.Status); err != nil {
			lrows.Close()
			return exam, err
		}
		exam.LedgerBooks = append(exam.LedgerBooks, lb)
	}
	lrows.Close()

	return exam, nil
}

// maskAccount returns the last-4 of an account number, keeping PII out of dashboards.
func maskAccount(acct string) string {
	if len(acct) <= 4 {
		return "****"
	}
	return "****" + acct[len(acct)-4:]
}

func nilStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
