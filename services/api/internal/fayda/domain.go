package fayda

import "time"

// Domain types for Fayda verification - implements id.gov.et flows

type VerificationMethod string

const (
	MethodOTP         VerificationMethod = "otp"
	MethodFace        VerificationMethod = "face"
	MethodFingerprint VerificationMethod = "fingerprint"
	MethodOfflineQR   VerificationMethod = "offline_qr"
	MethodOIDC        VerificationMethod = "oidc_esignet"
	MethodDemographic VerificationMethod = "demographic"
)

type VerificationStatus string

const (
	StatusInitiated      VerificationStatus = "initiated"
	StatusOTPSent        VerificationStatus = "otp_sent"
	StatusPendingConsent VerificationStatus = "pending_consent"
	StatusVerified       VerificationStatus = "verified"
	StatusFailed         VerificationStatus = "failed"
	StatusExpired        VerificationStatus = "expired"
)

// FaydaVerification mirrors DB table fayda_verifications
type FaydaVerification struct {
	ID                   string
	MerchantID           string
	OwnerID              *string
	KYCProfileID         string
	FinHash              string
	FinLast4             string
	FAN                  string
	PartnerCode          string
	RequestID            string
	FaydaTransactionID   string
	VerificationMethod   VerificationMethod
	OTPRequestedAt       *time.Time
	OTPVerified          bool
	OTPVerifiedAt        *time.Time
	ConsentTimestamp     time.Time
	ConsentIP            string
	Status               VerificationStatus
	DemographicsMatch    *bool
	DemographicsScore    int
	FaceMatch            *bool
	FaceMatchScore       float64
	ResponseEncryptedRef string
	FailureCode          string
	FailureMessage       string
	FrontDocID           *string
	BackDocID            *string
	SelfieDocID          *string
	VerifiedAt           *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// Request structs for outstanding UI

type InitRequest struct {
	MerchantID   string
	OwnerID      string
	KYCProfileID string
	FIN          string `validate:"omitempty,len=12"` // plain FIN from UI, never stored plain - hashed immediately
	FAN          string `validate:"omitempty,len=16"`
	Method       VerificationMethod
	FrontFileKey string // MinIO key for front image
	BackFileKey  string
	SelfieKey    string
	ConsentIP    string
}

type ConfirmOTPRequest struct {
	RequestID string
	OTP       string `validate:"required,len=6"`
}

// Fayda partner API payloads (per id.gov.et spec simplified)
type FaydaAuthRequest struct {
	PartnerCode   string        `json:"partnerCode"`
	PartnerAPIKey string        `json:"partnerApiKey"`
	UseCase       string        `json:"useCaseDescription"`
	FIN           string        `json:"fin"` // only for live call, not stored
	FAN           string        `json:"fan,omitempty"`
	OTP           string        `json:"otp,omitempty"`
	Demographics  *Demographics `json:"demographics,omitempty"`
}

type Demographics struct {
	Name        []LangValue `json:"name"`
	Gender      []LangValue `json:"gender"`
	DOB         string      `json:"dob"`
	FullAddress []LangValue `json:"fullAddress"`
}

type LangValue struct {
	Language string `json:"language"` // eng, amh
	Value    string `json:"value"`
}

// Offline QR
type QRData struct {
	// Decoded from FaydaEncode QR - contains masked data + signature
	FINLast4  string
	Name      string
	DOB       string
	Signature string // NIDP signature to verify offline
}
