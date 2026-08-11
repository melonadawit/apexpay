package payroll

import (
	"time"

	"github.com/shopspring/decimal"
)

// Note: This file reuses payroll package for business banking P0 features per industry parity
// Ideally would be in separate banking package, but for rapid P0 implementation we extend payroll package
// as payroll is part of workforce money OS and current accounts are heart of business finances per product description

// ==================== Current Accounts Real — Partner Bank CBE/Awash/Dashen ====================

type CurrentAccount struct {
	ID               string
	MerchantID       string
	AccountNumber    string // virtual account number e.g., ETB-CBE-1234567890
	AccountName      string
	AccountType      string // current, saving, virtual, escrow, reserve
	Currency         string
	BankCode         string // CBE, AWASH, DASHEN, ABYSSINIA, etc.
	PartnerBankName  string // Commercial Bank of Ethiopia, Awash Bank, Dashen Bank — Ethiopia equivalent of ICICI/Axis/RBL/YES in India
	Status           string // draft, pending_kyc, pending_approval, active, suspended, closed, frozen
	Balance          decimal.Decimal
	AvailableBalance decimal.Decimal
	OverdraftLimit   decimal.Decimal
	IsPrimary        bool
	IsLite           bool // lite interim account until current account active per lite account concept
	IsVirtual        bool // virtual account for collections smart collect
	ChequeBookIssued bool
	DebitCardIssued  bool
	DebitCardType    string // virtual, physical, both
	CreatedBy        *string
	ApprovedBy       *string
	ApprovedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CurrentAccountOpeningRequest struct {
	ID                   string
	MerchantID           string
	KYCProfileID         *string
	PartnerBankCode      string // CBE, AWASH, DASHEN etc.
	AccountType          string // current, saving, virtual
	RequestedAccountName string
	Status               string // draft, submitted, in_review, kyc_pending, approved, rejected, needs_more_info
	RiskScore            int
	RiskTier             string // low, medium, high
	SubmittedAt          *time.Time
	ReviewedAt           *time.Time
	ReviewerID           *string
	RejectionReason      string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type ChequeBook struct {
	ID                string
	CurrentAccountID  string
	MerchantID        string
	ChequeBookNumber  string
	StartChequeNumber int
	EndChequeNumber   int
	TotalCheques      int
	UsedCheques       int
	Status            string // ordered, issued, active, used_up, blocked, cancelled
	IssuedAt          *time.Time
	IssuedBy          *string
	CreatedAt         time.Time
}

type DebitCard struct {
	ID               string
	CurrentAccountID string
	MerchantID       string
	CardNumberMasked string // ****1234
	CardNumberHash   string // sha256 hash for lookup
	CardType         string // virtual, physical, both
	CardNetwork      string // visa, mastercard, verve, ethswitch
	Status           string // ordered, active, blocked, expired, cancelled
	DailyLimit       decimal.Decimal
	MonthlyLimit     decimal.Decimal
	CardholderName   string
	ExpiryMonth      *int
	ExpiryYear       *int
	CvvHash          string
	IsContactless    bool
	CreatedBy        *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ==================== Escrow Accounts Automated Marketplace P2P ====================

type EscrowAccount struct {
	ID               string
	MerchantID       string // marketplace operator
	AgreementID      string
	AccountNumber    string
	AccountName      string
	Amount           decimal.Decimal
	Currency         string
	Status           string // draft, held, released, returned, disputed, expired
	HeldAt           *time.Time
	ReleaseAt        *time.Time
	ReturnAt         *time.Time
	ExpiresAt        *time.Time
	BuyerMerchantID  *string
	SellerMerchantID *string
	OrderID          *string
	OrderAmount      decimal.Decimal
	PlatformFee      decimal.Decimal // e.g., 10% 100 ETB
	SellerAmount     decimal.Decimal // e.g., 90% 900 ETB
	WithholdingTax   decimal.Decimal // e.g., 2% 20 ETB per Ethiopia Income Tax Proclamation
	LedgerBookID     *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type EscrowAgreement struct {
	ID                    string
	MerchantID            string
	AgreementNumber       string
	Title                 string
	Description           string
	BuyerMerchantID       *string
	SellerMerchantID      *string
	Amount                decimal.Decimal
	Currency              string
	PlatformFeePercent    decimal.Decimal   // 10% platform fee
	WithholdingTaxPercent decimal.Decimal   // 2% withholding tax for services per Ethiopia Income Tax Proclamation
	Conditions            []EscrowCondition // JSON [{type: delivery_confirmed, days: 7}, {type: inspection_period, days: 3}]
	AutoRelease           bool
	AutoReleaseAfterDays  int
	Status                string // draft, active, completed, disputed, cancelled
	CreatedBy             *string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type EscrowCondition struct {
	Type string `json:"type"` // delivery_confirmed, inspection_period, etc.
	Days int    `json:"days"`
}

// ==================== Corporate Cards — Collateral-free Credit Cards ====================

type CorporateCard struct {
	ID                   string
	MerchantID           string
	CurrentAccountID     *string
	CardNumberMasked     string // ****1234
	CardNumberHash       string // sha256 hash
	CardType             string // virtual, physical, both
	CardNetwork          string // visa, mastercard, verve, ethswitch
	CardholderName       string
	CardholderEmail      string
	Status               string          // ordered, active, blocked, expired, cancelled, suspended
	CreditLimit          decimal.Decimal // up to 2Cr ETB equivalent (20L-2Cr INR in India)
	AvailableCredit      decimal.Decimal
	DailyLimit           decimal.Decimal
	MonthlyLimit         decimal.Decimal
	CategoryRestrictions []string               // ["SaaS", "Cloud", "Marketing"] etc.
	SpendingControls     map[string]interface{} // {daily_limit: 50000, monthly_limit: 500000, allowed_categories: ["SaaS", "Cloud"], blocked_merchants: []}
	CashbackPercent      decimal.Decimal        // flat 1% cashback per ApexPay
	ForexMarkupPercent   decimal.Decimal        // 2.5% forex markup
	InterestFreeDays     int                    // up to 45-50 day interest-free period
	IsAddon              bool
	ParentCardID         *string
	CreatedBy            *string
	ApprovedBy           *string
	ApprovedAt           *time.Time
	ExpiryMonth          *int
	ExpiryYear           *int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type CorporateCardTransaction struct {
	ID               string
	CardID           string
	MerchantID       string
	Amount           decimal.Decimal
	Currency         string
	MerchantName     string // e.g., AWS, Google Cloud, Facebook Ads
	MerchantCategory string // SaaS, Cloud, Marketing, etc.
	Status           string // pending, approved, declined, reversed
	DeclineReason    string
	CashbackAmount   decimal.Decimal // flat 1% cashback
	ForexFee         decimal.Decimal // 2.5% forex markup if international
	CreatedAt        time.Time
}

// ==================== Payout Links Enhanced — QR + Scan & Pay ====================

type EnhancedPayoutLink struct {
	ID             string
	MerchantID     string
	Amount         decimal.Decimal
	Currency       string
	PublicToken    string // for QR + public link
	QRCodeData     string // QR code data for EthSwitch interoperable QR
	RecipientName  string
	RecipientPhone string
	RecipientEmail string
	Purpose        string // refund, cashback, reward, vendor payment
	Status         string // active, claimed, expired, cancelled, failed
	ExpiresAt      time.Time
	ClaimedAt      *time.Time
	BeneficiaryID  *string // once claimed, beneficiary created
	EscrowBookID   *string // escrow book until claimed
	LedgerBookID   *string
	CreatedBy      *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ==================== Vendor Invoices + Purchase Orders + Petty Cash ====================

type VendorInvoice struct {
	ID                   string
	MerchantID           string
	VendorID             *string
	InvoiceNumber        string
	InvoiceDate          time.Time
	DueDate              *time.Time
	Amount               decimal.Decimal
	Currency             string
	TaxAmount            decimal.Decimal        // VAT 15% TOT 2%/10%
	WithholdingTaxAmount decimal.Decimal        // 2% for services per Ethiopia Income Tax Proclamation
	TotalAmount          decimal.Decimal        // amount + tax - withholding
	Status               string                 // draft, pending_approval, approved, paid, rejected, cancelled
	OCRRaw               map[string]interface{} // {extracted_text, confidence, vendor_name, tin, invoice_number, amount, tax, withholding, etc.}
	FileKey              *string                // MinIO file key for invoice PDF/image
	FileHash             *string
	CreatedBy            *string
	ApprovedBy           *string
	ApprovedAt           *time.Time
	PaidAt               *time.Time
	PayoutID             *string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type PurchaseOrder struct {
	ID         string
	MerchantID string
	VendorID   *string
	PONumber   string
	Amount     decimal.Decimal
	Currency   string
	Status     string // draft, sent, approved, received, cancelled, closed
	CreatedBy  *string
	ApprovedBy *string
	CreatedAt  time.Time
}

type PettyCashBudget struct {
	ID              string
	MerchantID      string
	BudgetName      string
	Amount          decimal.Decimal
	AssignedTo      *string
	Status          string // active, closed, exhausted
	SpentAmount     decimal.Decimal
	RemainingAmount decimal.Decimal
	CreatedBy       *string
	CreatedAt       time.Time
}

type PettyCashExpense struct {
	ID              string
	BudgetID        string
	MerchantID      string
	Amount          decimal.Decimal
	Description     string
	ReceiptFileKey  *string
	ReceiptFileHash *string
	Status          string // pending, approved, rejected, paid
	ApprovedBy      *string
	CreatedBy       *string
	CreatedAt       time.Time
}

// ==================== Tax Payments Automated Pre-filled Forms Challans Inbox ====================

type TaxPayment struct {
	ID               string
	MerchantID       string
	TaxType          string // vat, tot, withholding, paye, pension, corporate_tax, excise, other — VAT 15% TOT 2%/10% Withholding 2% PAYE (income tax) Pension 7%/11%
	Amount           decimal.Decimal
	Currency         string
	PeriodMonth      *int
	PeriodYear       *int
	DueDate          *time.Time
	Status           string  // draft, pending_approval, pending, paid, failed, cancelled
	ChallanFileKey   *string // MinIO file key for challan PDF
	ChallanFileHash  *string
	PaymentReference string // bank payment reference / UTR
	PaidAt           *time.Time
	CreatedBy        *string
	ApprovedBy       *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type TaxAccountant struct {
	ID          string
	MerchantID  string
	UserID      string
	Role        string // compliance, accountant, auditor, viewer
	Permissions map[string]interface{}
	CreatedAt   time.Time
}

// ==================== Bank Account Verification Penny Testing 1 ETB ====================

type BankAccountVerification struct {
	ID                      string
	MerchantID              string
	BankCode                string
	AccountNumberMasked     string
	AccountNumberHash       string
	AccountName             string
	VerificationMethod      string                 // penny_test, micro_deposit, bank_letter, manual
	Amount                  decimal.Decimal        // 1.00 ETB penny test
	ConnectorID             string                 // bank_ips, telebirr, etc.
	Status                  string                 // pending, processing, verified, failed, expired
	VerificationResponse    map[string]interface{} // {beneficiary_name, account_name_match_score, bank_details, etc.}
	BeneficiaryNameReturned string
	MatchScore              decimal.Decimal // fuzzy Levenshtein <3 match score
	VerifiedAt              *time.Time
	ExpiresAt               *time.Time
	CreatedAt               time.Time
}

// ==================== Collections / Smart Collect / Virtual Accounts ====================

type VirtualAccount struct {
	ID                   string
	MerchantID           string
	VirtualAccountNumber string // e.g., VA-CBE-1234567890
	CustomerID           *string
	Purpose              string // e.g., customer collections, vendor payments, payroll
	Status               string // active, inactive, closed
	BankCode             string
	CreatedAt            time.Time
}

type VirtualAccountTransaction struct {
	ID               string
	VirtualAccountID string
	MerchantID       string
	Amount           decimal.Decimal
	Currency         string
	UTR              string // Unique Transaction Reference
	SenderName       string
	SenderAccount    string
	Status           string // pending, matched, unmatched, reconciled, failed
	MatchedInvoiceID *string
	MatchedAt        *time.Time
	CreatedAt        time.Time
}
