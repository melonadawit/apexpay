package admin

// Admin-domain DTOs. These are the read models returned to the compliance/ops/admin
// dashboards. All money-adjacent values are string-rendered decimals to keep the wire
// format JSON-safe and avoid float drift.

// OnboardingQueueItem is one row in the reviewer work queue.
type OnboardingQueueItem struct {
	MerchantID       string `json:"merchant_id"`
	LegalName        string `json:"legal_name"`
	Email            string `json:"email"`
	OnboardingStatus string `json:"onboarding_status"`
	RiskScore        int    `json:"risk_score"`
	RiskTier         string `json:"risk_tier"`
	FaydaVerified    bool   `json:"fayda_verified"`
	SubmittedAt      string `json:"submitted_at,omitempty"`
	CreatedAt        string `json:"created_at"`
}

// ReviewRequest is the payload of POST /admin/onboarding/{id}/review.
type ReviewRequest struct {
	Action  string `json:"action"` // approve | reject | request_info
	Comment string `json:"comment"`
}

// ReviewResult describes the outcome of a review action.
type ReviewResult struct {
	MerchantID string `json:"merchant_id"`
	Status     string `json:"status"`
	Action     string `json:"action"`
	Message    string `json:"message"`
}

// ConnectorHealth is one row in the live connector-health dashboard.
type ConnectorHealth struct {
	Connector    string  `json:"connector"`
	AvgLatencyMS int     `json:"avg_latency_5m"`
	SuccessRate  float64 `json:"success_rate_5m"`
	Samples      int     `json:"samples"`
	Circuit      string  `json:"circuit"`
}

// ReconBreak is an open reconciliation case needing a human decision.
type ReconBreak struct {
	MerchantID     string `json:"merchant_id"`
	IdempotencyKey string `json:"idempotency_key"`
	TxRef          string `json:"tx_ref,omitempty"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
}

// Evidence is the compliance evidence bundle for a transaction.
type Evidence struct {
	TxRef             string            `json:"tx_ref"`
	LedgerJournals    []LedgerJournal   `json:"ledger_journals"`
	FaydaVerified     *bool             `json:"fayda_verified,omitempty"`
	DocHashes         []string          `json:"docs_hashes"`
	OnboardingReviews []ReviewRecord    `json:"onboarding_reviews"`
	AuditLogs         []AuditLog        `json:"audit_logs"`
	WebhookDeliveries []WebhookDelivery `json:"webhook_deliveries"`
}

// LedgerJournal is a journal with its posting key and memo.
type LedgerJournal struct {
	ID         string `json:"id"`
	PostingKey string `json:"posting_key"`
	Memo       string `json:"memo,omitempty"`
	RefType    string `json:"reference_type,omitempty"`
	RefID      string `json:"reference_id,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// ReviewRecord is a recorded onboarding review transition.
type ReviewRecord struct {
	ID           string `json:"id"`
	ReviewerType string `json:"reviewer_type"`
	FromStatus   string `json:"from_status"`
	ToStatus     string `json:"to_status"`
	Action       string `json:"action"`
	Comments     string `json:"comments,omitempty"`
	CreatedAt    string `json:"created_at"`
}

// AuditLog is one append-only audit entry.
type AuditLog struct {
	ID           string `json:"id"`
	ActorType    string `json:"actor_type"`
	ActorID      string `json:"actor_id,omitempty"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type,omitempty"`
	ResourceID   string `json:"resource_id,omitempty"`
	CreatedAt    string `json:"created_at"`
}

// WebhookDelivery is one webhook delivery attempt record.
type WebhookDelivery struct {
	ID           string `json:"id"`
	EventType    string `json:"event_type"`
	Status       string `json:"status"`
	AttemptCount int    `json:"attempt_count"`
	LastCode     int    `json:"last_status_code,omitempty"`
	LastError    string `json:"last_error,omitempty"`
	CreatedAt    string `json:"created_at"`
}

// MerchantExam is the full onboarding compliance file for a merchant.
type MerchantExam struct {
	MerchantID        string            `json:"merchant_id"`
	LegalName         string            `json:"legal_name"`
	OnboardingStatus  string            `json:"onboarding_status"`
	RiskScore         int               `json:"risk_score"`
	RiskTier          string            `json:"risk_tier"`
	KYCProfiles       []KYCProfile      `json:"kyc_profiles"`
	Owners            []Owner           `json:"owners"`
	Documents         []Document        `json:"documents"`
	ComplianceChecks  []ComplianceCheck `json:"compliance_checks"`
	OnboardingReviews []ReviewRecord    `json:"onboarding_reviews"`
	Banks             []BankAccount     `json:"banks"`
	LedgerBooks       []LedgerBook      `json:"ledger_books"`
}

// KYCProfile summarizes one KYC submission version.
type KYCProfile struct {
	ID               string `json:"id"`
	Version          int    `json:"version"`
	LegalName        string `json:"legal_name"`
	TINNumber        string `json:"tin_number"`
	BusinessType     string `json:"business_type"`
	OnboardingStatus string `json:"onboarding_status"`
	SubmittedAt      string `json:"submitted_at,omitempty"`
	ReviewedAt       string `json:"reviewed_at,omitempty"`
}

// Owner summarizes a beneficial owner (PII-safe: hashes/last4 only).
type Owner struct {
	ID                 string `json:"id"`
	FullName           string `json:"full_name"`
	Role               string `json:"role"`
	FaydaVerified      bool   `json:"fayda_verified"`
	IDType             string `json:"id_type"`
	IDNumberLast4      string `json:"id_number_last4,omitempty"`
	VerificationStatus string `json:"verification_status"`
}

// Document summarizes an uploaded KYC document.
type Document struct {
	ID        string `json:"id"`
	DocType   string `json:"doc_type"`
	Status    string `json:"status"`
	FileHash  string `json:"file_hash"`
	CreatedAt string `json:"created_at"`
}

// ComplianceCheck summarizes a single compliance check result.
type ComplianceCheck struct {
	CheckType string `json:"check_type"`
	Status    string `json:"status"`
	Score     int    `json:"score"`
	Provider  string `json:"provider,omitempty"`
	UpdatedAt string `json:"updated_at"`
}

// BankAccount summarizes a merchant bank account (masked).
type BankAccount struct {
	ID                  string `json:"id"`
	BankCode            string `json:"bank_code"`
	AccountNumberMasked string `json:"account_number_masked"`
	AccountName         string `json:"account_name"`
	IsVerified          bool   `json:"is_verified"`
}

// LedgerBook summarizes a merchant ledger book.
type LedgerBook struct {
	ID       string `json:"id"`
	BookType string `json:"book_type"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}
