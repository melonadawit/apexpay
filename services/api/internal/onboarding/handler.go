package onboarding

import (
	"encoding/json"
	"net/http"

	"apexpay/internal/id"
	"apexpay/internal/platform/crypto"
	pkghttp "apexpay/internal/platform/http"
	mw "apexpay/internal/platform/middleware"
	"apexpay/internal/platform/storage"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct {
	svc     *Service
	storage *storage.Client
}

func NewHandler(svc *Service, storage *storage.Client) *Handler {
	return &Handler{svc: svc, storage: storage}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/kyc", h.CreateKYC)
	r.Get("/kyc/{id}", h.GetKYC)
	r.Post("/owners", h.AddOwner)
	r.Post("/bank-accounts", h.AddBankAccount)
	r.Post("/documents/presign", h.PresignDocument)
	r.Post("/documents", h.CreateDocument)
	r.Post("/submit", h.SubmitKYC)
	r.Get("/status", h.Status)
	r.Get("/timeline", h.Timeline)
}

func (h *Handler) CreateKYC(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MerchantID         string `json:"merchant_id"`
		LegalName          string `json:"legal_name"`
		BusinessType       string `json:"business_type"`
		RegistrationNumber string `json:"registration_number"`
		TINNumber          string `json:"tin_number"`
		Industry           string `json:"industry_category"`
		Description        string `json:"business_description"`
		Region             string `json:"region"`
		City               string `json:"city"`
		AddressFull        string `json:"office_address_full"`
		ContactName        string `json:"contact_person_name"`
		ContactRole        string `json:"contact_person_role"`
		ContactEmail       string `json:"contact_email"`
		ContactPhone       string `json:"contact_phone"`
		ExpectedTPV        string `json:"expected_monthly_tpv"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	tpv, _ := decimal.NewFromString(req.ExpectedTPV)
	profile := &KYCProfile{
		MerchantID: req.MerchantID, LegalName: req.LegalName, BusinessType: BusinessType(req.BusinessType),
		RegistrationNumber: req.RegistrationNumber, TINNumber: req.TINNumber, Industry: Industry(req.Industry),
		Description: req.Description, Region: req.Region, City: req.City, AddressFull: req.AddressFull,
		ContactPersonName: req.ContactName, ContactPersonRole: req.ContactRole, ContactEmail: req.ContactEmail, ContactPhone: req.ContactPhone,
		ExpectedMonthlyTPV: tpv, AvgTicketAmount: decimal.NewFromInt(500),
	}
	resp, err := h.svc.CreateKYC(r.Context(), profile)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, resp)
}

func (h *Handler) GetKYC(w http.ResponseWriter, r *http.Request) {
	merchantID := mw.MerchantID(r.Context())
	idStr := chi.URLParam(r, "id")
	p, err := h.svc.repo.GetKYCProfile(r.Context(), merchantID, idStr)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, p)
}

func (h *Handler) AddOwner(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MerchantID          string `json:"merchant_id"`
		KYCProfileID        string `json:"kyc_profile_id"`
		FullName            string `json:"full_name"`
		FullNameAM          string `json:"full_name_am"`
		Role                string `json:"role"`
		OwnershipPercentage string `json:"ownership_percentage"`
		IDType              string `json:"id_type"`
		Phone               string `json:"phone"`
		Email               string `json:"email"`
		IsAuthSignatory     bool   `json:"is_authorized_signatory"`
		IsPEP               bool   `json:"is_pep"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	pct, _ := decimal.NewFromString(req.OwnershipPercentage)
	owner := &BeneficialOwner{
		MerchantID: req.MerchantID, KYCProfileID: req.KYCProfileID,
		FullName: req.FullName, FullNameAM: req.FullNameAM, Role: OwnerRole(req.Role),
		OwnershipPercentage: pct, IDType: req.IDType, Phone: req.Phone, Email: req.Email,
		IsAuthorizedSignatory: req.IsAuthSignatory, IsPEP: req.IsPEP, Nationality: "ET",
	}
	resp, err := h.svc.AddOwner(r.Context(), owner)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, resp)
}

func (h *Handler) AddBankAccount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MerchantID    string `json:"merchant_id"`
		AccountName   string `json:"account_name"`
		AccountNumber string `json:"account_number"`
		BankCode      string `json:"bank_code"`
		BankName      string `json:"bank_name"`
		IsDefault     bool   `json:"is_settlement_default"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	// Hash account number for privacy per DATABASE
	hash := crypto.HashFIN("salt", req.AccountNumber) // reuse hash fn
	masked := crypto.MaskAccount(req.AccountNumber)
	acc := &BankAccount{
		ID: id.New("bank"), MerchantID: req.MerchantID, AccountName: req.AccountName,
		AccountNumberMasked: masked, AccountNumberHash: hash,
		BankCode: req.BankCode, BankName: req.BankName,
		IsSettlementDefault: req.IsDefault, VerificationStatus: "pending",
	}
	if err := h.svc.repo.CreateBankAccount(r.Context(), acc); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, acc)
}

func (h *Handler) PresignDocument(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MerchantID  string `json:"merchant_id"`
		DocType     string `json:"doc_type"`
		FileName    string `json:"file_name"`
		ContentType string `json:"content_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	objectKey := storage.ObjectKey(req.MerchantID, req.DocType, id.NewDocument(), "pdf")
	url, err := h.storage.PresignedPutURL(r.Context(), objectKey, 15*60*1000000000) // 15m
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"object_key": objectKey, "upload_url": url, "expires_in": "15m"})
}

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MerchantID   string `json:"merchant_id"`
		KYCProfileID string `json:"kyc_profile_id"`
		DocType      string `json:"doc_type"`
		FileKey      string `json:"file_key"`
		FileHash     string `json:"file_hash"`
		MimeType     string `json:"mime_type"`
		SizeBytes    int    `json:"file_size_bytes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	doc := &Document{
		ID: id.NewDocument(), MerchantID: req.MerchantID, KYCProfileID: req.KYCProfileID,
		Type: DocType(req.DocType), FileKey: req.FileKey, FileHash: req.FileHash, MimeType: req.MimeType, SizeBytes: req.SizeBytes, Status: "uploaded",
	}
	if err := h.svc.repo.CreateDocument(r.Context(), doc); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 201, doc)
}

func (h *Handler) SubmitKYC(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MerchantID   string `json:"merchant_id"`
		KYCProfileID string `json:"kyc_profile_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
		return
	}
	if err := h.svc.SubmitKYC(r.Context(), req.MerchantID, req.KYCProfileID); err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	pkghttp.WriteJSON(w, r, 200, map[string]string{"status": "submitted", "message": "KYC submitted for compliance review, risk scoring pending"})
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	merchantID := r.URL.Query().Get("merchant_id")
	if merchantID == "" {
		// try context
		if authenticatedMerchantID := mw.MerchantID(r.Context()); authenticatedMerchantID != "" {
			merchantID = authenticatedMerchantID
		}
	}
	profile, err := h.svc.repo.GetLatestKYCProfile(r.Context(), merchantID)
	if err != nil {
		pkghttp.WriteError(w, r, err)
		return
	}
	owners, _ := h.svc.repo.ListOwners(r.Context(), merchantID, profile.ID)
	banks, _ := h.svc.repo.ListBankAccounts(r.Context(), merchantID)
	docs, _ := h.svc.repo.ListDocuments(r.Context(), merchantID, profile.ID)
	checks, _ := h.svc.repo.ListComplianceChecks(r.Context(), merchantID, profile.ID)

	timeline := h.svc.Timeline(profile.OnboardingStatus)
	pkghttp.WriteJSON(w, r, 200, map[string]interface{}{
		"merchant_id": merchantID, "kyc_profile": profile, "owners": owners, "banks": banks, "documents": docs, "compliance_checks": checks, "timeline": timeline, "progress": len(timeline),
	})
}

func (h *Handler) Timeline(w http.ResponseWriter, r *http.Request) {
	statusStr := r.URL.Query().Get("status")
	tl := h.svc.Timeline(OnboardingStatus(statusStr))
	pkghttp.WriteJSON(w, r, 200, tl)
}
