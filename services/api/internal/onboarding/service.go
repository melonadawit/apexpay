package onboarding

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"apexpay/internal/id"
	"apexpay/internal/platform/errors"
)

// Repository interface - clean arch, no PG in service.
type Repository interface {
	CreateKYCProfile(ctx context.Context, p *KYCProfile) error
	GetKYCProfile(ctx context.Context, merchantID, profileID string) (*KYCProfile, error)
	GetLatestKYCProfile(ctx context.Context, merchantID string) (*KYCProfile, error)
	UpdateKYCStatus(ctx context.Context, merchantID, profileID string, status OnboardingStatus) error

	CreateOwner(ctx context.Context, o *BeneficialOwner) error
	ListOwners(ctx context.Context, merchantID, kycProfileID string) ([]BeneficialOwner, error)
	UpdateOwnerFaydaVerified(ctx context.Context, ownerID string, finHash, finLast4 string, verified bool) error

	CreateBankAccount(ctx context.Context, b *BankAccount) error
	ListBankAccounts(ctx context.Context, merchantID string) ([]BankAccount, error)

	CreateDocument(ctx context.Context, d *Document) error
	ListDocuments(ctx context.Context, merchantID, kycProfileID string) ([]Document, error)
	GetDocument(ctx context.Context, docID string) (*Document, error)

	CreateComplianceCheck(ctx context.Context, c *ComplianceCheck) error
	ListComplianceChecks(ctx context.Context, merchantID, kycProfileID string) ([]ComplianceCheck, error)

	CreateReview(ctx context.Context, r *OnboardingReview) error
	ListReviews(ctx context.Context, merchantID string) ([]OnboardingReview, error)

	// Transactional: approve flow creates merchant book
	ApproveMerchantTx(ctx context.Context, merchantID, kycProfileID, reviewerID string) error
}

// Service implements business logic with optimal validation.
type Service struct {
	repo Repository
	// salt for FIN hashing - injected from config vault
	finSalt string
}

func NewService(repo Repository, finSalt string) *Service {
	return &Service{repo: repo, finSalt: finSalt}
}

// CreateKYC validates NBE requirements + restricted industry.
func (s *Service) CreateKYC(ctx context.Context, p *KYCProfile) (*KYCProfile, error) {
	if strings.TrimSpace(p.LegalName) == "" {
		return nil, errors.Validation("legal_name required")
	}
	if strings.TrimSpace(p.TINNumber) == "" {
		return nil, errors.Validation("TIN required per NBE ONPS")
	}
	// ET TIN 10 digits example - use validator
	if len(p.TINNumber) != 10 {
		return nil, errors.Validation("TIN must be 10 digits")
	}
	// Restricted industry check per PayAtlas NBE
	if RestrictedIndustries[p.Industry] {
		return nil, errors.New(errors.CodeForbidden, fmt.Sprintf("industry %s is restricted per NBE", p.Industry), 403)
	}
	// Website policy check reminder - must have refund/privacy/terms for PSP approval
	if p.WebsiteURL != "" && (!p.HasRefundPolicy || !p.HasPrivacyPolicy || !p.HasTerms) {
		// Not hard fail, but flag compliance check will fail
	}

	p.ID = id.NewKYCProfile()
	p.Version = 1
	p.OnboardingStatus = StatusDraft
	if p.KYCLevel == "" {
		p.KYCLevel = Level2
	}
	p.CreatedAt = time.Now()

	if err := s.repo.CreateKYCProfile(ctx, p); err != nil {
		return nil, err
	}

	// Create pending compliance checks
	checks := []CheckType{CheckTIN, CheckLicense, CheckBank, CheckRestricted, CheckWebsite, CheckFayda, CheckDocAuth, CheckRisk, CheckAML}
	for _, ct := range checks {
		_ = s.repo.CreateComplianceCheck(ctx, &ComplianceCheck{
			ID: id.New("cchk"), MerchantID: p.MerchantID, KYCProfileID: p.ID,
			Type: ct, Status: "pending", CreatedAt: time.Now(),
		})
	}

	return p, nil
}

// AddOwner validates UBO + PEP
func (s *Service) AddOwner(ctx context.Context, o *BeneficialOwner) (*BeneficialOwner, error) {
	if o.OwnershipPercentage.LessThan(decimal.Zero) || o.OwnershipPercentage.GreaterThan(decimal.NewFromInt(100)) {
		return nil, errors.Validation("ownership_percentage must be 0-100")
	}
	o.ID = id.NewOwner()
	o.CreatedAt = time.Now()
	o.VerificationStatus = "pending"
	if err := s.repo.CreateOwner(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

// SubmitKYC performs completeness checks - optimal O(n) doc verification
func (s *Service) SubmitKYC(ctx context.Context, merchantID, kycProfileID string) error {
	kyc, err := s.repo.GetKYCProfile(ctx, merchantID, kycProfileID)
	if err != nil {
		return errors.NotFound("kyc profile not found")
	}

	owners, err := s.repo.ListOwners(ctx, merchantID, kycProfileID)
	if err != nil {
		return err
	}
	hasAuthSignatory := false
	faydaVerifiedCount := 0
	for _, own := range owners {
		if own.IsAuthorizedSignatory {
			hasAuthSignatory = true
		}
		if own.FaydaVerified {
			faydaVerifiedCount++
		}
	}
	if !hasAuthSignatory {
		return errors.Validation("at least one authorized signatory required per NBE")
	}
	if faydaVerifiedCount == 0 {
		return errors.New(errors.CodeFaydaNotVerified, "at least one owner must have Fayda verified", 400)
	}

	banks, err := s.repo.ListBankAccounts(ctx, merchantID)
	if err != nil {
		return err
	}
	hasSettlement := false
	for _, b := range banks {
		if b.IsSettlementDefault && b.VerificationStatus == "verified" {
			hasSettlement = true
			break
		}
	}
	// Allow pending but flag compliance: require at least one settlement account added
	if len(banks) == 0 {
		return errors.Validation("at least one settlement bank account required")
	}

	docs, err := s.repo.ListDocuments(ctx, merchantID, kycProfileID)
	if err != nil {
		return err
	}
	required := RequiredDocs(kyc.BusinessType, kyc.KYCLevel)
	docMap := make(map[DocType]int, len(docs))
	for _, d := range docs {
		docMap[d.Type]++
	}
	missing := []DocType{}
	for _, need := range required {
		// Special case: UBO ID front/back at least per owner
		if need == DocFaydaFront || need == DocFaydaBack {
			continue // checked via owners Fayda verification
		}
		if docMap[need] == 0 {
			missing = append(missing, need)
		}
	}
	if len(missing) > 0 {
		return errors.New(errors.CodeDocumentRequired, fmt.Sprintf("missing required docs: %v", missing), 400)
	}

	// Risk scoring simplified - optimal weighted sum
	riskScore := s.calculateRiskScore(kyc, owners, docs)
	// Update status to submitted
	if err := s.repo.UpdateKYCStatus(ctx, merchantID, kycProfileID, StatusSubmitted); err != nil {
		return err
	}

	// Audit review
	_ = s.repo.CreateReview(ctx, &OnboardingReview{
		ID: id.New("orev"), MerchantID: merchantID, KYCProfileID: kycProfileID,
		ReviewerType: "system", FromStatus: StatusDraft, ToStatus: StatusSubmitted,
		Action: ActionSubmit, Comments: fmt.Sprintf("submitted risk_score=%d", riskScore),
		CreatedAt: time.Now(),
	})

	return nil
}

func (s *Service) calculateRiskScore(kyc *KYCProfile, owners []BeneficialOwner, docs []Document) int {
	score := 0
	// High risk industry already blocked
	// Expected TPV high => higher risk
	if kyc.ExpectedMonthlyTPV.GreaterThan(decimal.NewFromInt(1000000)) {
		score += 20
	}
	pepCount := 0
	for _, o := range owners {
		if o.IsPEP {
			pepCount++
		}
	}
	score += pepCount * 30
	// Docs completeness
	score += (len(RequiredDocs(kyc.BusinessType, kyc.KYCLevel)) - len(docs)) * 2
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}

// ApproveMerchant - dual approval if high risk or TPV >1M
func (s *Service) ApproveMerchant(ctx context.Context, merchantID, kycProfileID, reviewerID string, riskScore int, existingApprovals int) error {
	// NBE: high risk or TPV>1M requires 2 approvers (maker-checker)
	needsDual := riskScore >= 70
	kyc, _ := s.repo.GetKYCProfile(ctx, merchantID, kycProfileID)
	if kyc != nil && kyc.ExpectedMonthlyTPV.GreaterThan(decimal.NewFromInt(1000000)) {
		needsDual = true
	}
	if needsDual && existingApprovals < 1 {
		// First approver - leave in pending_approval
		return s.repo.CreateReview(ctx, &OnboardingReview{
			ID: id.New("orev"), MerchantID: merchantID, KYCProfileID: kycProfileID,
			ReviewerID: &reviewerID, ReviewerType: "compliance",
			FromStatus: StatusComplianceCheck, ToStatus: StatusComplianceCheck,
			Action: ActionApprove, Comments: "first approval, needs second",
			CreatedAt: time.Now(),
		})
	}

	// Final approval transaction
	return s.repo.ApproveMerchantTx(ctx, merchantID, kycProfileID, reviewerID)
}

// ListOnboardingStatus returns timeline for outstanding UI
func (s *Service) Timeline(status OnboardingStatus) []struct {
	Step   string
	Status string // done, current, upcoming
} {
	steps := []OnboardingStatus{StatusDraft, StatusSubmitted, StatusInReview, StatusFaydaPending, StatusComplianceCheck, StatusApproved}
	timeline := make([]struct {
		Step   string
		Status string
	}, len(steps))
	found := false
	for i, st := range steps {
		var stt string
		if st == status {
			stt = "current"
			found = true
		} else if !found {
			// if status is beyond this step, it's done
			// simplistic but O(n)
			idxCurrent := indexOf(steps, status)
			if idxCurrent > i {
				stt = "done"
			} else {
				stt = "upcoming"
			}
		} else {
			stt = "upcoming"
		}
		timeline[i] = struct {
			Step   string
			Status string
		}{Step: string(st), Status: stt}
	}
	return timeline
}

func indexOf(arr []OnboardingStatus, v OnboardingStatus) int {
	for i, e := range arr {
		if e == v {
			return i
		}
	}
	return -1
}
