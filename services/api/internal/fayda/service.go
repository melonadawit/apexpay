package fayda

import (
	"context"
	"fmt"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/platform/crypto"
	"apexpay/internal/platform/errors"
)

// Verifier interface - mock vs live (Strategy pattern for best practice)
type Verifier interface {
	RequestOTP(ctx context.Context, fin, fan string) (faydaTxID string, err error)
	VerifyOTP(ctx context.Context, faydaTxID, otp string) (demographicsMatch bool, faceScore float64, err error)
	VerifyOfflineQR(ctx context.Context, qrData string) (valid bool, parsed QRData, err error)
	VerifyOIDC(ctx context.Context, code string) (FaydaVerification, error)
}

type Repository interface {
	Create(ctx context.Context, v *FaydaVerification) error
	GetByRequestID(ctx context.Context, requestID string) (*FaydaVerification, error)
	GetByOwner(ctx context.Context, ownerID string) ([]FaydaVerification, error)
	UpdateStatus(ctx context.Context, requestID string, status VerificationStatus, txID string, otpVerified bool) error
	UpdateVerificationResult(ctx context.Context, requestID string, demoMatch bool, faceMatch bool, faceScore float64, encryptedRef string) error
}

type Service struct {
	repo         Repository
	verifier     Verifier
	salt         string // for fin hashing
	partnerCode  string
	storageEncryptKey []byte
}

func NewService(repo Repository, verifier Verifier, salt, partnerCode string, encKey []byte) *Service {
	return &Service{
		repo: repo, verifier: verifier, salt: salt, partnerCode: partnerCode,
		storageEncryptKey: encKey,
	}
}

// Init starts Fayda verification - outstanding modern flow with front/back images already uploaded to MinIO via presigned URLs
func (s *Service) Init(ctx context.Context, req InitRequest) (*FaydaVerification, error) {
	// Validate FIN/FAN format per crypto package
	if req.FIN != "" && !crypto.ValidateFaydaFIN(req.FIN) {
		return nil, errors.New(errors.CodeInvalidFaydaFin, "FIN must be 12 digits", 400)
	}
	if req.FAN != "" && !crypto.ValidateFAN(req.FAN) {
		return nil, errors.Validation("FAN must be 16 alphanumeric")
	}
	if req.FIN == "" && req.FAN == "" {
		return nil, errors.Validation("FIN or FAN required")
	}

	// Privacy: hash FIN immediately, keep last4 only
	finHash := ""
	finLast4 := ""
	if req.FIN != "" {
		finHash = crypto.HashFIN(s.salt, req.FIN)
		finLast4 = crypto.Last4(req.FIN)
	} else {
		// If FAN only, we still request OTP via Fayda - FAN lookup returns FIN hash via NIDP
		// For mock, we generate hash from FAN
		finHash = crypto.HashFIN(s.salt, req.FAN)
		finLast4 = crypto.Last4(req.FAN)
	}

	requestID := id.NewFaydaVerification()
	now := time.Now()

	// Call Fayda partner API to request OTP (Strategy: mock will return fake txID)
	var faydaTxID string
	var err error
	if req.Method == MethodOfflineQR {
		// offline path no OTP
		faydaTxID = "offline_" + requestID
	} else {
		faydaTxID, err = s.verifier.RequestOTP(ctx, req.FIN, req.FAN)
		if err != nil {
			return nil, fmt.Errorf("fayda request OTP failed: %w", err)
		}
	}

	v := &FaydaVerification{
		ID:                 id.New("fayda"),
		MerchantID:         req.MerchantID,
		OwnerID:            &req.OwnerID,
		KYCProfileID:       req.KYCProfileID,
		FinHash:            finHash,
		FinLast4:           finLast4,
		FAN:                req.FAN,
		PartnerCode:        s.partnerCode,
		RequestID:          requestID,
		FaydaTransactionID: faydaTxID,
		VerificationMethod: req.Method,
		OTPRequestedAt:     &now,
		ConsentTimestamp:   now,
		ConsentIP:          req.ConsentIP,
		Status:             StatusOTPSent,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}

	// Audit: never log plain FIN
	return v, nil
}

// ConfirmOTP verifies OTP with Fayda and updates result
func (s *Service) ConfirmOTP(ctx context.Context, req ConfirmOTPRequest) (*FaydaVerification, error) {
	v, err := s.repo.GetByRequestID(ctx, req.RequestID)
	if err != nil {
		return nil, errors.NotFound("verification request not found")
	}
	if v.Status == StatusVerified {
		return v, nil // idempotent
	}
	if v.OTPRequestedAt != nil && time.Since(*v.OTPRequestedAt) > 5*time.Minute {
		_ = s.repo.UpdateStatus(ctx, req.RequestID, StatusExpired, v.FaydaTransactionID, false)
		return nil, errors.New(errors.CodeFaydaOTPFailed, "OTP expired", 400)
	}

	demoMatch, faceScore, err := s.verifier.VerifyOTP(ctx, v.FaydaTransactionID, req.OTP)
	if err != nil {
		_ = s.repo.UpdateStatus(ctx, req.RequestID, StatusFailed, v.FaydaTransactionID, false)
		return nil, errors.New(errors.CodeFaydaOTPFailed, "OTP verification failed", 400)
	}

	faceMatch := faceScore >= 0.85
	now := time.Now()

	// Encrypt response payload (mock)
	encryptedRef := fmt.Sprintf("fayda_responses/%s.enc", v.RequestID) // MinIO ref

	err = s.repo.UpdateVerificationResult(ctx, req.RequestID, demoMatch, faceMatch, faceScore, encryptedRef)
	if err != nil {
		return nil, err
	}
	_ = s.repo.UpdateStatus(ctx, req.RequestID, StatusVerified, v.FaydaTransactionID, true)

	v.DemographicsMatch = &demoMatch
	b := faceMatch
	v.FaceMatch = &b
	v.FaceMatchScore = faceScore
	v.OTPVerified = true
	v.OTPVerifiedAt = &now
	v.Status = StatusVerified
	v.VerifiedAt = &now

	return v, nil
}

// VerifyOfflineQR handles FaydaEncode QR scan from Flutter
func (s *Service) VerifyOfflineQR(ctx context.Context, qrString string, requestID string) (bool, error) {
	valid, _, err := s.verifier.VerifyOfflineQR(ctx, qrString)
	if err != nil {
		return false, err
	}
	if valid {
		_ = s.repo.UpdateStatus(ctx, requestID, StatusVerified, "offline_qr", true)
	}
	return valid, nil
}
