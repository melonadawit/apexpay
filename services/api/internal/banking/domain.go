package banking

// Banking domain read models. All money-adjacent values are string-rendered decimals
// (numeric scanned as text) to keep JSON safe and avoid float drift.

type CurrentAccount struct {
	ID               string `json:"id"`
	AccountNumber    string `json:"account_number"`
	AccountName      string `json:"account_name"`
	AccountType      string `json:"account_type"`
	Currency         string `json:"currency"`
	BankCode         string `json:"bank_code"`
	PartnerBankName  string `json:"partner_bank_name"`
	Status           string `json:"status"`
	Balance          string `json:"balance"`
	AvailableBalance string `json:"available_balance"`
	OverdraftLimit   string `json:"overdraft_limit"`
	IsPrimary        bool   `json:"is_primary"`
	IsLite           bool   `json:"is_lite"`
	IsVirtual        bool   `json:"is_virtual"`
	ChequeBookIssued bool   `json:"cheque_book_issued"`
	DebitCardIssued  bool   `json:"debit_card_issued"`
}

type CreditLine struct {
	ID              string `json:"id"`
	CreditLimit     string `json:"credit_limit"`
	AvailableCredit string `json:"available_credit"`
	UtilizedCredit  string `json:"utilized_credit"`
	InterestRate    string `json:"interest_rate"`
	Status          string `json:"status"`
	CreditScore     *int   `json:"credit_score,omitempty"`
}

type ForexRequest struct {
	ID                string `json:"id"`
	FromCurrency      string `json:"from_currency"`
	ToCurrency        string `json:"to_currency"`
	FromAmount        string `json:"from_amount"`
	ToAmount          string `json:"to_amount"`
	RateUsed          string `json:"rate_used"`
	ForexFeePercent   string `json:"forex_fee_percent"`
	ForexFeeAmount    string `json:"forex_fee_amount"`
	Purpose           string `json:"purpose"`
	Status            string `json:"status"`
	NBEApprovalStatus string `json:"nbe_approval_status,omitempty"`
	CreatedAt         string `json:"created_at"`
}

type VirtualAccount struct {
	ID                   string `json:"id"`
	VirtualAccountNumber string `json:"virtual_account_number"`
	CustomerID           string `json:"customer_id"`
	Purpose              string `json:"purpose"`
	Status               string `json:"status"`
	BankCode             string `json:"bank_code"`
	CreatedAt            string `json:"created_at"`
}

type Notification struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	IsRead    bool   `json:"is_read"`
	ActionURL string `json:"action_url"`
	CreatedAt string `json:"created_at"`
}

type CorporateCard struct {
	ID               string `json:"id"`
	CardNumberMasked string `json:"card_number_masked"`
	CardType         string `json:"card_type"`
	CardNetwork      string `json:"card_network"`
	CardholderName   string `json:"cardholder_name"`
	CardholderEmail  string `json:"cardholder_email"`
	Status           string `json:"status"`
	CreditLimit      string `json:"credit_limit"`
	AvailableCredit  string `json:"available_credit"`
	ForexMarkup      string `json:"forex_markup_percent"`
	CashbackPercent  string `json:"cashback_percent"`
}

type EscrowAccount struct {
	ID            string `json:"id"`
	AgreementID   string `json:"agreement_id"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Amount        string `json:"amount"`
	Status        string `json:"status"`
	OrderID       string `json:"order_id,omitempty"`
	PlatformFee   string `json:"platform_fee"`
	SellerAmount  string `json:"seller_amount"`
}

type SupportTicket struct {
	ID         string `json:"id"`
	Subject    string `json:"subject"`
	Priority   string `json:"priority"`
	Status     string `json:"status"`
	AssignedTo string `json:"assigned_to,omitempty"`
	CreatedAt  string `json:"created_at"`
}

type RelationshipManager struct {
	ID         string `json:"id"`
	RMUserID   string `json:"rm_user_id"`
	Status     string `json:"status"`
	AssignedAt string `json:"assigned_at"`
}

type BankVerification struct {
	ID            string `json:"id"`
	BankCode      string `json:"bank_code"`
	AccountMasked string `json:"account_number_masked"`
	AccountName   string `json:"account_name"`
	Method        string `json:"verification_method"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

type ForexRate struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	Rate         string `json:"rate"`
	BuyRate      string `json:"buy_rate"`
	SellRate     string `json:"sell_rate"`
	Source       string `json:"source"`
	LastUpdated  string `json:"last_updated_at"`
}

// ---- Create request types (action-oriented modules) ----

type VendorInvoice struct {
	ID                   string  `json:"id"`
	VendorID             string  `json:"vendor_id,omitempty"`
	InvoiceNumber        string  `json:"invoice_number"`
	InvoiceDate          string  `json:"invoice_date"`
	DueDate              string  `json:"due_date,omitempty"`
	Amount               string  `json:"amount"`
	Currency             string  `json:"currency"`
	TaxAmount            string  `json:"tax_amount"`
	WithholdingTaxAmount string  `json:"withholding_tax_amount"`
	TotalAmount          string  `json:"total_amount"`
	Status               string  `json:"status"`
	OCRConfidence        float64 `json:"ocr_confidence"`
	VendorName           string  `json:"vendor_name,omitempty"`
	FileKey              string  `json:"file_key,omitempty"`
	CreatedAt            string  `json:"created_at"`
}

type PettyCashBudget struct {
	ID              string `json:"id"`
	BudgetName      string `json:"budget_name"`
	Amount          string `json:"amount"`
	AssignedTo      string `json:"assigned_to,omitempty"`
	Status          string `json:"status"`
	SpentAmount     string `json:"spent_amount"`
	RemainingAmount string `json:"remaining_amount"`
	CreatedAt       string `json:"created_at"`
}

type PettyCashExpense struct {
	ID          string `json:"id"`
	BudgetID    string `json:"budget_id"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
	ReceiptKey  string `json:"receipt_file_key,omitempty"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type TaxPayment struct {
	ID               string `json:"id"`
	TaxType          string `json:"tax_type"`
	Amount           string `json:"amount"`
	Currency         string `json:"currency"`
	PeriodMonth      *int   `json:"period_month,omitempty"`
	PeriodYear       *int   `json:"period_year,omitempty"`
	DueDate          string `json:"due_date,omitempty"`
	Status           string `json:"status"`
	PaymentReference string `json:"payment_reference,omitempty"`
	PaidAt           string `json:"paid_at,omitempty"`
	CreatedAt        string `json:"created_at"`
}

type PayoutLink struct {
	ID             string `json:"id"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	PublicToken    string `json:"public_token"`
	RecipientName  string `json:"recipient_name,omitempty"`
	RecipientPhone string `json:"recipient_phone,omitempty"`
	RecipientEmail string `json:"recipient_email,omitempty"`
	Purpose        string `json:"purpose,omitempty"`
	Status         string `json:"status"`
	ExpiresAt      string `json:"expires_at"`
	CreatedAt      string `json:"created_at"`
}
