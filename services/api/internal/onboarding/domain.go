package onboarding

import (
	"time"

	"github.com/shopspring/decimal"
)

// Enums per DATABASE v1.1.0
type BusinessType string

const (
	BusinessTypeSoleProp BusinessType = "sole_proprietorship"
	BusinessTypePLC      BusinessType = "plc"
	BusinessTypeShare    BusinessType = "share_company"
	BusinessTypePartner  BusinessType = "partnership"
	BusinessTypeCoop     BusinessType = "cooperative"
)

type OnboardingStatus string

const (
	StatusDraft           OnboardingStatus = "draft"
	StatusSubmitted       OnboardingStatus = "submitted"
	StatusInReview        OnboardingStatus = "in_review"
	StatusFaydaPending    OnboardingStatus = "fayda_pending"
	StatusComplianceCheck OnboardingStatus = "compliance_check"
	StatusNeedsMoreInfo   OnboardingStatus = "needs_more_info"
	StatusApproved        OnboardingStatus = "approved"
	StatusRejected        OnboardingStatus = "rejected"
	StatusActive          OnboardingStatus = "active"
)

type KYCLevel string

const (
	Level1 KYCLevel = "level1" // referral, <=5000 bal per NBE tiered
	Level2 KYCLevel = "level2" // standard
	Level3 KYCLevel = "level3" // full + physical
)

type Industry string

const (
	IndustryEcommerce Industry = "e-commerce"
	IndustryEducation Industry = "education"
	IndustryHealth    Industry = "health"
	IndustryLogistics Industry = "logistics"
	IndustryFood      Industry = "food"
	IndustryTech      Industry = "tech"
	IndustryOther     Industry = "other"
	// Restricted
	IndustryGambling Industry = "gambling"
	IndustryCrypto   Industry = "crypto"
	IndustryAdult    Industry = "adult"
)

var RestrictedIndustries = map[Industry]bool{
	IndustryGambling: true,
	IndustryCrypto:   true,
	IndustryAdult:    true,
}

// KYCProfile is merchant_kyc_profiles
type KYCProfile struct {
	ID                 string
	MerchantID         string
	Version            int
	LegalName          string
	TradeName          string
	BusinessType       BusinessType
	RegistrationNumber string
	TINNumber          string
	VATNumber          string
	BusinessLicenseNo  string
	Industry           Industry
	Description        string
	WebsiteURL         string
	ExpectedMonthlyTPV decimal.Decimal
	AvgTicketAmount    decimal.Decimal
	Region             string
	City               string
	SubCity            string
	Woreda             string
	AddressFull        string
	ContactPersonName  string
	ContactPersonRole  string
	ContactEmail       string
	ContactPhone       string
	HasRefundPolicy    bool
	HasPrivacyPolicy   bool
	HasTerms           bool
	OnboardingStatus   OnboardingStatus
	KYCLevel           KYCLevel
	SubmittedAt        *time.Time
	CreatedAt          time.Time
}

// BeneficialOwner
type OwnerRole string

const (
	RoleOwner         OwnerRole = "owner"
	RoleShareholder   OwnerRole = "shareholder"
	RoleDirector      OwnerRole = "director"
	RoleAuthorizedRep OwnerRole = "authorized_rep"
	RoleUBO           OwnerRole = "ubo"
)

type BeneficialOwner struct {
	ID                    string
	MerchantID            string
	KYCProfileID          string
	FullName              string
	FullNameAM            string
	Role                  OwnerRole
	OwnershipPercentage   decimal.Decimal
	Nationality           string
	IDType                string // fayda, passport...
	FinHash               string
	FinLast4              string
	FAN                   string
	FaydaVerified         bool
	DateOfBirth           *time.Time
	Phone                 string
	Email                 string
	IsPEP                 bool
	IsAuthorizedSignatory bool
	VerificationStatus    string // pending, fayda_verified, verified
	CreatedAt             time.Time
}

// BankAccount
type BankAccount struct {
	ID                  string
	MerchantID          string
	AccountName         string
	AccountNumberMasked string
	AccountNumberHash   string
	BankCode            string
	BankName            string
	IsSettlementDefault bool
	VerificationStatus  string
}

// Document types per NBE + PayAtlas
type DocType string

const (
	DocCompanyReg        DocType = "company_registration"
	DocTINCertificate    DocType = "tin_certificate"
	DocBusinessLicense   DocType = "business_license"
	DocVAT               DocType = "vat_certificate"
	DocMemorandum        DocType = "memorandum_articles"
	DocBoardResolution   DocType = "board_resolution"
	DocShareholderList   DocType = "shareholder_list"
	DocUBOIDFront        DocType = "ubo_id_front"
	DocUBOIDBack         DocType = "ubo_id_back"
	DocFaydaFront        DocType = "fayda_card_front"
	DocFaydaBack         DocType = "fayda_card_back"
	DocProofOfAddress    DocType = "proof_of_address"
	DocBankLetter        DocType = "bank_letter"
	DocWebsiteScreenshot DocType = "website_screenshot"
	DocRefundPolicy      DocType = "refund_policy_doc"
	DocOther             DocType = "other"
)

type Document struct {
	ID           string
	MerchantID   string
	KYCProfileID string
	OwnerID      *string
	Type         DocType
	FileKey      string // MinIO
	FileHash     string
	MimeType     string
	SizeBytes    int
	Status       string // pending, uploaded, ocr_done, verified, rejected
	OCRRaw       map[string]any
	ExpiresAt    *time.Time
	CreatedAt    time.Time
}

// ComplianceCheck
type CheckType string

const (
	CheckTIN        CheckType = "tin_validation"
	CheckLicense    CheckType = "business_license_validation"
	CheckBank       CheckType = "bank_account_validation"
	CheckAML        CheckType = "aml_screening"
	CheckPEP        CheckType = "pep_check"
	CheckRestricted CheckType = "restricted_industry"
	CheckWebsite    CheckType = "website_policy_check"
	CheckFayda      CheckType = "fayda_verification"
	CheckDocAuth    CheckType = "document_authenticity"
	CheckRisk       CheckType = "risk_scoring"
)

type ComplianceCheck struct {
	ID           string
	MerchantID   string
	KYCProfileID string
	Type         CheckType
	Status       string // pending, passed, failed, needs_review
	Score        int
	Provider     string
	Details      map[string]any
	CreatedAt    time.Time
}

// Onboarding Review - maker-checker
type ReviewAction string

const (
	ActionSubmit      ReviewAction = "submit"
	ActionApprove     ReviewAction = "approve"
	ActionReject      ReviewAction = "reject"
	ActionRequestInfo ReviewAction = "request_info"
)

type OnboardingReview struct {
	ID           string
	MerchantID   string
	KYCProfileID string
	ReviewerID   *string
	ReviewerType string // system, ops, compliance, admin
	FromStatus   OnboardingStatus
	ToStatus     OnboardingStatus
	Action       ReviewAction
	Comments     string
	CreatedAt    time.Time
}

// Required docs calculator per business type & KYC level - optimal data structure
func RequiredDocs(bType BusinessType, level KYCLevel) []DocType {
	base := []DocType{DocTINCertificate, DocBusinessLicense, DocProofOfAddress, DocBankLetter}
	switch bType {
	case BusinessTypePLC, BusinessTypeShare:
		base = append(base, DocCompanyReg, DocMemorandum, DocBoardResolution, DocShareholderList)
	}
	if level == Level2 || level == Level3 {
		base = append(base, DocFaydaFront, DocFaydaBack)
	}
	// Website policy docs are mandatory per PayAtlas for PSP
	base = append(base, DocWebsiteScreenshot, DocRefundPolicy)
	return unique(base)
}

func unique(in []DocType) []DocType {
	seen := make(map[DocType]bool, len(in))
	out := make([]DocType, 0, len(in))
	for _, v := range in {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
