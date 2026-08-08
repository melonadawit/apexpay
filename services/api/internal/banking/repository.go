package banking

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"apexpay/internal/id"
)

// Repository provides read models for the merchant banking modules.
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) CurrentAccounts(ctx context.Context, merchantID string) ([]CurrentAccount, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, account_number, account_name, account_type, currency, bank_code, partner_bank_name, status,
		       balance::text, available_balance::text, overdraft_limit::text,
		       is_primary, is_lite, is_virtual, cheque_book_issued, debit_card_issued
		FROM current_accounts WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []CurrentAccount{}
	for rows.Next() {
		var a CurrentAccount
		if err := rows.Scan(&a.ID, &a.AccountNumber, &a.AccountName, &a.AccountType, &a.Currency,
			&a.BankCode, &a.PartnerBankName, &a.Status, &a.Balance, &a.AvailableBalance,
			&a.OverdraftLimit, &a.IsPrimary, &a.IsLite, &a.IsVirtual, &a.ChequeBookIssued, &a.DebitCardIssued); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *Repository) CreditLines(ctx context.Context, merchantID string) ([]CreditLine, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, credit_limit::text, available_credit::text, utilized_credit::text, interest_rate::text, status, credit_score
		FROM credit_lines WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []CreditLine{}
	for rows.Next() {
		var c CreditLine
		if err := rows.Scan(&c.ID, &c.CreditLimit, &c.AvailableCredit, &c.UtilizedCredit,
			&c.InterestRate, &c.Status, &c.CreditScore); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *Repository) ForexRates(ctx context.Context) ([]ForexRate, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT from_currency, to_currency, rate::text, buy_rate::text, sell_rate::text, source,
		       to_char(last_updated_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM forex_rates ORDER BY to_currency`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []ForexRate{}
	for rows.Next() {
		var f ForexRate
		if err := rows.Scan(&f.FromCurrency, &f.ToCurrency, &f.Rate, &f.BuyRate, &f.SellRate, &f.Source, &f.LastUpdated); err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, rows.Err()
}

func (r *Repository) ForexRequests(ctx context.Context, merchantID string) ([]ForexRequest, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, from_currency, to_currency, from_amount::text, to_amount::text, rate_used::text,
		       forex_fee_percent::text, forex_fee_amount::text, purpose, status,
		       COALESCE(nbe_approval_status,''),
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM forex_requests WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []ForexRequest{}
	for rows.Next() {
		var f ForexRequest
		if err := rows.Scan(&f.ID, &f.FromCurrency, &f.ToCurrency, &f.FromAmount, &f.ToAmount,
			&f.RateUsed, &f.ForexFeePercent, &f.ForexFeeAmount, &f.Purpose, &f.Status,
			&f.NBEApprovalStatus, &f.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, rows.Err()
}

func (r *Repository) VirtualAccounts(ctx context.Context, merchantID string) ([]VirtualAccount, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, virtual_account_number, COALESCE(customer_id,''), COALESCE(purpose,''), status, bank_code,
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM virtual_accounts WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []VirtualAccount{}
	for rows.Next() {
		var v VirtualAccount
		if err := rows.Scan(&v.ID, &v.VirtualAccountNumber, &v.CustomerID, &v.Purpose, &v.Status, &v.BankCode, &v.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

func (r *Repository) Notifications(ctx context.Context, merchantID string, limit int) ([]Notification, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, type, COALESCE(title,''), COALESCE(message,''), is_read, COALESCE(action_url,''),
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM notifications WHERE merchant_id=$1 ORDER BY created_at DESC LIMIT $2`, merchantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Notification{}
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.Type, &n.Title, &n.Message, &n.IsRead, &n.ActionURL, &n.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

func (r *Repository) CorporateCards(ctx context.Context, merchantID string) ([]CorporateCard, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, card_number_masked, card_type, COALESCE(card_network,''), cardholder_name,
		       COALESCE(cardholder_email,''), status, credit_limit::text, available_credit::text,
		       forex_markup_percent::text, cashback_percent::text
		FROM corporate_cards WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []CorporateCard{}
	for rows.Next() {
		var c CorporateCard
		if err := rows.Scan(&c.ID, &c.CardNumberMasked, &c.CardType, &c.CardNetwork, &c.CardholderName,
			&c.CardholderEmail, &c.Status, &c.CreditLimit, &c.AvailableCredit, &c.ForexMarkup, &c.CashbackPercent); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *Repository) EscrowAccounts(ctx context.Context, merchantID string) ([]EscrowAccount, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE(agreement_id,''), COALESCE(account_number,''), COALESCE(account_name,''),
		       amount::text, status, COALESCE(order_id,''), platform_fee::text, seller_amount::text
		FROM escrow_accounts WHERE buyer_merchant_id=$1 OR seller_merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []EscrowAccount{}
	for rows.Next() {
		var e EscrowAccount
		if err := rows.Scan(&e.ID, &e.AgreementID, &e.AccountNumber, &e.AccountName, &e.Amount,
			&e.Status, &e.OrderID, &e.PlatformFee, &e.SellerAmount); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *Repository) SupportTickets(ctx context.Context, merchantID string) ([]SupportTicket, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE(subject,''), COALESCE(priority,''), status, COALESCE(assigned_to,''),
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM support_tickets WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []SupportTicket{}
	for rows.Next() {
		var t SupportTicket
		if err := rows.Scan(&t.ID, &t.Subject, &t.Priority, &t.Status, &t.AssignedTo, &t.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *Repository) RelationshipManagers(ctx context.Context, merchantID string) ([]RelationshipManager, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, rm_user_id, status,
		       to_char(assigned_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM relationship_managers WHERE merchant_id=$1 ORDER BY assigned_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []RelationshipManager{}
	for rows.Next() {
		var r RelationshipManager
		if err := rows.Scan(&r.ID, &r.RMUserID, &r.Status, &r.AssignedAt); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, rows.Err()
}

func (r *Repository) BankVerifications(ctx context.Context, merchantID string) ([]BankVerification, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, bank_code, account_number_masked, account_name, verification_method, status,
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM bank_account_verifications WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []BankVerification{}
	for rows.Next() {
		var v BankVerification
		if err := rows.Scan(&v.ID, &v.BankCode, &v.AccountMasked, &v.AccountName, &v.Method, &v.Status, &v.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

// ---- Vendor invoices ----

func (r *Repository) VendorInvoices(ctx context.Context, merchantID string) ([]VendorInvoice, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE(vendor_id,''), invoice_number,
		       to_char(invoice_date,'YYYY-MM-DD'), COALESCE(to_char(due_date,'YYYY-MM-DD'),''),
		       amount::text, currency, tax_amount::text, withholding_tax_amount::text, total_amount::text, status,
		       COALESCE((ocr_raw->>'confidence')::float8,0), COALESCE(ocr_raw->>'vendor_name',''), COALESCE(file_key,''),
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM vendor_invoices WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []VendorInvoice{}
	for rows.Next() {
		var v VendorInvoice
		if err := rows.Scan(&v.ID, &v.VendorID, &v.InvoiceNumber, &v.InvoiceDate, &v.DueDate,
			&v.Amount, &v.Currency, &v.TaxAmount, &v.WithholdingTaxAmount, &v.TotalAmount, &v.Status,
			&v.OCRConfidence, &v.VendorName, &v.FileKey, &v.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

func (r *Repository) CreateVendorInvoice(ctx context.Context, merchantID, userID string, in *VendorInvoice) error {
	id := id.New("vinv")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO vendor_invoices (id, merchant_id, vendor_id, invoice_number, invoice_date, due_date, amount, currency,
			tax_amount, withholding_tax_amount, total_amount, status, ocr_raw, file_key, created_by)
		VALUES ($1,$2,$3,$4,$5::date, NULLIF($6,'')::date, $7::numeric, $8, $9::numeric, $10::numeric, $11::numeric, $12, $13::jsonb, $14, $15)`,
		id, merchantID, nilString(in.VendorID), in.InvoiceNumber, in.InvoiceDate, in.DueDate,
		in.Amount, in.Currency, in.TaxAmount, in.WithholdingTaxAmount, in.TotalAmount, in.Status,
		`{"confidence":`+jsonFloat(in.OCRConfidence)+`,"vendor_name":"`+escapeJSON(in.VendorName)+`"}`, in.FileKey, userID)
	return err
}

// ---- Petty cash ----

func (r *Repository) PettyCashBudgets(ctx context.Context, merchantID string) ([]PettyCashBudget, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, budget_name, amount::text, COALESCE(assigned_to,''), status, spent_amount::text, remaining_amount::text,
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM petty_cash_budgets WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []PettyCashBudget{}
	for rows.Next() {
		var b PettyCashBudget
		if err := rows.Scan(&b.ID, &b.BudgetName, &b.Amount, &b.AssignedTo, &b.Status, &b.SpentAmount, &b.RemainingAmount, &b.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, rows.Err()
}

func (r *Repository) CreatePettyCashBudget(ctx context.Context, merchantID, userID string, in *PettyCashBudget) error {
	id := id.New("pcb")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO petty_cash_budgets (id, merchant_id, budget_name, amount, assigned_to, status, spent_amount, remaining_amount, created_by)
		VALUES ($1,$2,$3,$4::numeric, NULLIF($5,''), 'active', 0, $4::numeric, $6)`,
		id, merchantID, in.BudgetName, in.Amount, in.AssignedTo, userID)
	return err
}

func (r *Repository) PettyCashExpenses(ctx context.Context, merchantID string) ([]PettyCashExpense, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, budget_id, amount::text, description, COALESCE(receipt_file_key,''), status,
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM petty_cash_expenses WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []PettyCashExpense{}
	for rows.Next() {
		var e PettyCashExpense
		if err := rows.Scan(&e.ID, &e.BudgetID, &e.Amount, &e.Description, &e.ReceiptKey, &e.Status, &e.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *Repository) CreatePettyCashExpense(ctx context.Context, merchantID, userID string, in *PettyCashExpense) error {
	id := id.New("pce")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO petty_cash_expenses (id, budget_id, merchant_id, amount, description, receipt_file_key, status, created_by)
		VALUES ($1,$2,$3,$4::numeric,$5,$6,'pending',$7)`,
		id, in.BudgetID, merchantID, in.Amount, in.Description, nilString(in.ReceiptKey), userID)
	return err
}

// ---- Tax payments ----

func (r *Repository) TaxPayments(ctx context.Context, merchantID string) ([]TaxPayment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tax_type, amount::text, currency, period_month, period_year, COALESCE(to_char(due_date,'YYYY-MM-DD'),''), status,
		       COALESCE(payment_reference,''), COALESCE(to_char(paid_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS'),''),
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM tax_payments WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []TaxPayment{}
	for rows.Next() {
		var t TaxPayment
		var pm, py *int
		if err := rows.Scan(&t.ID, &t.TaxType, &t.Amount, &t.Currency, &pm, &py, &t.DueDate,
			&t.Status, &t.PaymentReference, &t.PaidAt, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.PeriodMonth, t.PeriodYear = pm, py
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *Repository) CreateTaxPayment(ctx context.Context, merchantID, userID string, in *TaxPayment) error {
	id := id.New("tax")
	_, err := r.pool.Exec(ctx, `
		INSERT INTO tax_payments (id, merchant_id, tax_type, amount, currency, period_month, period_year, due_date, status, created_by)
		VALUES ($1,$2,$3,$4::numeric,$5,$6,$7,NULLIF($8,'')::date,$9,$10)`,
		id, merchantID, in.TaxType, in.Amount, in.Currency, in.PeriodMonth, in.PeriodYear, in.DueDate, in.Status, userID)
	return err
}

// ---- Payout links (enhanced) ----

func (r *Repository) PayoutLinks(ctx context.Context, merchantID string) ([]PayoutLink, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, amount::text, currency, public_token, COALESCE(recipient_name,''), COALESCE(recipient_phone,''),
		       COALESCE(recipient_email,''), COALESCE(purpose,''), status,
		       to_char(expires_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS'),
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM payout_links_enhanced WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []PayoutLink{}
	for rows.Next() {
		var p PayoutLink
		if err := rows.Scan(&p.ID, &p.Amount, &p.Currency, &p.PublicToken, &p.RecipientName,
			&p.RecipientPhone, &p.RecipientEmail, &p.Purpose, &p.Status, &p.ExpiresAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *Repository) CreatePayoutLink(ctx context.Context, merchantID, userID string, in *PayoutLink) error {
	id := id.New("plink")
	token := "pl_" + id[5:]
	_, err := r.pool.Exec(ctx, `
		INSERT INTO payout_links_enhanced (id, merchant_id, amount, currency, public_token, recipient_name, recipient_phone, recipient_email, purpose, status, expires_at, created_by)
		VALUES ($1,$2,$3::numeric,$4,$5,$6,$7,$8,$9,'active',now()+interval '7 days',$10)`,
		id, merchantID, in.Amount, in.Currency, token, nilString(in.RecipientName), nilString(in.RecipientPhone), nilString(in.RecipientEmail), nilString(in.Purpose), userID)
	in.ID = id
	in.PublicToken = token
	in.Status = "active"
	return err
}

// ---- small helpers ----

func nilString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func jsonFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func escapeJSON(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

type AccountingIntegration struct {
	ID             string `json:"id"`
	Provider       string `json:"provider"`
	Status         string `json:"status"`
	LastSyncAt     string `json:"last_sync_at,omitempty"`
	LastSyncStatus string `json:"last_sync_status,omitempty"`
	LastSyncError  string `json:"last_sync_error,omitempty"`
	CreatedAt      string `json:"created_at"`
}

func (r *Repository) AccountingIntegrations(ctx context.Context, merchantID string) ([]AccountingIntegration, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, provider, status,
		       COALESCE(to_char(last_sync_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS'),''),
		       COALESCE(last_sync_status,''), COALESCE(last_sync_error,''),
		       to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM accounting_integrations WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []AccountingIntegration{}
	for rows.Next() {
		var a AccountingIntegration
		if err := rows.Scan(&a.ID, &a.Provider, &a.Status, &a.LastSyncAt, &a.LastSyncStatus, &a.LastSyncError, &a.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}
