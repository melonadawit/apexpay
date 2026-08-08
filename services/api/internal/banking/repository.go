package banking

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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
