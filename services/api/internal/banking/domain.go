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
