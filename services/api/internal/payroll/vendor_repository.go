package payroll

import (
	"context"
	"encoding/json"

	"github.com/shopspring/decimal"
)

// ==================== Vendor Invoices — OCR-enabled Invoice Capture + TDS Calculation ====================

func (r *PgRepository) CreateVendorInvoice(ctx context.Context, invoice *VendorInvoice) error {
	ocrRaw, _ := json.Marshal(invoice.OCRRaw)
	_, err := r.pool.Exec(ctx, `INSERT INTO vendor_invoices (id, merchant_id, vendor_id, invoice_number, invoice_date, due_date, amount, currency, tax_amount, withholding_tax_amount, total_amount, status, ocr_raw, file_key, file_hash, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		invoice.ID, invoice.MerchantID, invoice.VendorID, invoice.InvoiceNumber, invoice.InvoiceDate, invoice.DueDate, invoice.Amount.String(), invoice.Currency, invoice.TaxAmount.String(), invoice.WithholdingTaxAmount.String(), invoice.TotalAmount.String(), invoice.Status, string(ocrRaw), invoice.FileKey, invoice.FileHash, invoice.CreatedBy)
	return err
}

func (r *PgRepository) GetVendorInvoice(ctx context.Context, merchantID, invoiceID string) (*VendorInvoice, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, vendor_id, invoice_number, invoice_date, due_date, amount::text, currency, tax_amount::text, withholding_tax_amount::text, total_amount::text, status, ocr_raw, file_key, file_hash FROM vendor_invoices WHERE merchant_id=$1 AND id=$2`, merchantID, invoiceID)
	var inv VendorInvoice
	var amount, taxAmt, withholdingAmt, totalAmt, ocrRawStr string
	err := row.Scan(&inv.ID, &inv.MerchantID, &inv.VendorID, &inv.InvoiceNumber, &inv.InvoiceDate, &inv.DueDate, &amount, &inv.Currency, &taxAmt, &withholdingAmt, &totalAmt, &inv.Status, &ocrRawStr, &inv.FileKey, &inv.FileHash)
	if err != nil {
		return nil, err
	}
	inv.Amount, _ = decimal.NewFromString(amount)
	inv.TaxAmount, _ = decimal.NewFromString(taxAmt)
	inv.WithholdingTaxAmount, _ = decimal.NewFromString(withholdingAmt)
	inv.TotalAmount, _ = decimal.NewFromString(totalAmt)
	_ = json.Unmarshal([]byte(ocrRawStr), &inv.OCRRaw)
	return &inv, nil
}

func (r *PgRepository) ListVendorInvoices(ctx context.Context, merchantID string, status *string) ([]VendorInvoice, error) {
	query := `SELECT id, merchant_id, vendor_id, invoice_number, invoice_date, due_date, amount::text, currency, tax_amount::text, withholding_tax_amount::text, total_amount::text, status, ocr_raw FROM vendor_invoices WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if status != nil {
		query += ` AND status=$2`
		args = append(args, *status)
	}
	query += ` ORDER BY invoice_date DESC, created_at DESC`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []VendorInvoice
	for rows.Next() {
		var inv VendorInvoice
		var amount, taxAmt, withholdingAmt, totalAmt, ocrRawStr string
		if err := rows.Scan(&inv.ID, &inv.MerchantID, &inv.VendorID, &inv.InvoiceNumber, &inv.InvoiceDate, &inv.DueDate, &amount, &inv.Currency, &taxAmt, &withholdingAmt, &totalAmt, &inv.Status, &ocrRawStr); err != nil {
			return nil, err
		}
		inv.Amount, _ = decimal.NewFromString(amount)
		inv.TaxAmount, _ = decimal.NewFromString(taxAmt)
		inv.WithholdingTaxAmount, _ = decimal.NewFromString(withholdingAmt)
		inv.TotalAmount, _ = decimal.NewFromString(totalAmt)
		_ = json.Unmarshal([]byte(ocrRawStr), &inv.OCRRaw)
		list = append(list, inv)
	}
	return list, nil
}

func (r *PgRepository) UpdateVendorInvoiceStatus(ctx context.Context, invoiceID, status, approvedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE vendor_invoices SET status=$1, approved_by=$2, approved_at=now(), updated_at=now() WHERE id=$3`, status, approvedBy, invoiceID)
	return err
}

func (r *PgRepository) MarkVendorInvoicePaid(ctx context.Context, invoiceID, payoutID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE vendor_invoices SET status='paid', paid_at=now(), payout_id=$2, updated_at=now() WHERE id=$1`, invoiceID, payoutID)
	return err
}

// ==================== Purchase Orders ====================

func (r *PgRepository) CreatePurchaseOrder(ctx context.Context, po *PurchaseOrder) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO purchase_orders (id, merchant_id, vendor_id, po_number, amount, currency, status, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		po.ID, po.MerchantID, po.VendorID, po.PONumber, po.Amount.String(), po.Currency, po.Status, po.CreatedBy)
	return err
}

func (r *PgRepository) ListPurchaseOrders(ctx context.Context, merchantID string) ([]PurchaseOrder, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, vendor_id, po_number, amount::text, currency, status FROM purchase_orders WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []PurchaseOrder
	for rows.Next() {
		var po PurchaseOrder
		var amount string
		if err := rows.Scan(&po.ID, &po.MerchantID, &po.VendorID, &po.PONumber, &amount, &po.Currency, &po.Status); err != nil {
			return nil, err
		}
		po.Amount, _ = decimal.NewFromString(amount)
		list = append(list, po)
	}
	return list, nil
}

// ==================== Petty Cash Budgets ====================

func (r *PgRepository) CreatePettyCashBudget(ctx context.Context, budget *PettyCashBudget) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO petty_cash_budgets (id, merchant_id, budget_name, amount, assigned_to, status, spent_amount, remaining_amount, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		budget.ID, budget.MerchantID, budget.BudgetName, budget.Amount.String(), budget.AssignedTo, budget.Status, budget.SpentAmount.String(), budget.RemainingAmount.String(), budget.CreatedBy)
	return err
}

func (r *PgRepository) ListPettyCashBudgets(ctx context.Context, merchantID string) ([]PettyCashBudget, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, budget_name, amount::text, assigned_to, status, spent_amount::text, remaining_amount::text FROM petty_cash_budgets WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []PettyCashBudget
	for rows.Next() {
		var b PettyCashBudget
		var amount, spent, remaining string
		if err := rows.Scan(&b.ID, &b.MerchantID, &b.BudgetName, &amount, &b.AssignedTo, &b.Status, &spent, &remaining); err != nil {
			return nil, err
		}
		b.Amount, _ = decimal.NewFromString(amount)
		b.SpentAmount, _ = decimal.NewFromString(spent)
		b.RemainingAmount, _ = decimal.NewFromString(remaining)
		list = append(list, b)
	}
	return list, nil
}

func (r *PgRepository) CreatePettyCashExpense(ctx context.Context, expense *PettyCashExpense) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO petty_cash_expenses (id, budget_id, merchant_id, amount, description, receipt_file_key, receipt_file_hash, status, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		expense.ID, expense.BudgetID, expense.MerchantID, expense.Amount.String(), expense.Description, expense.ReceiptFileKey, expense.ReceiptFileHash, expense.Status, expense.CreatedBy)
	return err
}

func (r *PgRepository) ListPettyCashExpenses(ctx context.Context, merchantID, budgetID string) ([]PettyCashExpense, error) {
	query := `SELECT id, budget_id, merchant_id, amount::text, description, receipt_file_key, receipt_file_hash, status FROM petty_cash_expenses WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if budgetID != "" {
		query += ` AND budget_id=$2`
		args = append(args, budgetID)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []PettyCashExpense
	for rows.Next() {
		var e PettyCashExpense
		var amount string
		if err := rows.Scan(&e.ID, &e.BudgetID, &e.MerchantID, &amount, &e.Description, &e.ReceiptFileKey, &e.ReceiptFileHash, &e.Status); err != nil {
			return nil, err
		}
		e.Amount, _ = decimal.NewFromString(amount)
		list = append(list, e)
	}
	return list, nil
}

// ==================== Tax Payments ====================

func (r *PgRepository) CreateTaxPayment(ctx context.Context, tax *TaxPayment) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO tax_payments (id, merchant_id, tax_type, amount, currency, period_month, period_year, due_date, status, payment_reference, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		tax.ID, tax.MerchantID, tax.TaxType, tax.Amount.String(), tax.Currency, tax.PeriodMonth, tax.PeriodYear, tax.DueDate, tax.Status, tax.PaymentReference, tax.CreatedBy)
	return err
}

func (r *PgRepository) ListTaxPayments(ctx context.Context, merchantID string, taxType *string, status *string) ([]TaxPayment, error) {
	query := `SELECT id, merchant_id, tax_type, amount::text, currency, period_month, period_year, due_date, status, challan_file_key, challan_file_hash, payment_reference FROM tax_payments WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if taxType != nil {
		query += ` AND tax_type=$2`
		args = append(args, *taxType)
	}
	if status != nil {
		query += ` AND status=$3`
		args = append(args, *status)
	}
	query += ` ORDER BY due_date ASC, created_at DESC`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []TaxPayment
	for rows.Next() {
		var t TaxPayment
		var amount string
		if err := rows.Scan(&t.ID, &t.MerchantID, &t.TaxType, &amount, &t.Currency, &t.PeriodMonth, &t.PeriodYear, &t.DueDate, &t.Status, &t.ChallanFileKey, &t.ChallanFileHash, &t.PaymentReference); err != nil {
			return nil, err
		}
		t.Amount, _ = decimal.NewFromString(amount)
		list = append(list, t)
	}
	return list, nil
}

func (r *PgRepository) UpdateTaxPaymentStatus(ctx context.Context, taxID, status, challanFileKey, paymentReference string) error {
	_, err := r.pool.Exec(ctx, `UPDATE tax_payments SET status=$1, challan_file_key=$2, payment_reference=$3, updated_at=now() WHERE id=$4`, status, challanFileKey, paymentReference, taxID)
	return err
}

func (r *PgRepository) MarkTaxPaymentPaid(ctx context.Context, taxID, challanFileKey, paymentReference string) error {
	_, err := r.pool.Exec(ctx, `UPDATE tax_payments SET status='paid', challan_file_key=$1, payment_reference=$2, paid_at=now(), updated_at=now() WHERE id=$3`, challanFileKey, paymentReference, taxID)
	return err
}

// ==================== Bank Account Verification Penny Testing 1 ETB ====================

func (r *PgRepository) CreateBankVerification(ctx context.Context, v *BankAccountVerification) error {
	response, _ := json.Marshal(v.VerificationResponse)
	_, err := r.pool.Exec(ctx, `INSERT INTO bank_account_verifications (id, merchant_id, bank_code, account_number_masked, account_number_hash, account_name, verification_method, amount, connector_id, status, verification_response, beneficiary_name_returned, match_score, expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		v.ID, v.MerchantID, v.BankCode, v.AccountNumberMasked, v.AccountNumberHash, v.AccountName, v.VerificationMethod, v.Amount.String(), v.ConnectorID, v.Status, string(response), v.BeneficiaryNameReturned, v.MatchScore.String(), v.ExpiresAt)
	return err
}

func (r *PgRepository) GetBankVerification(ctx context.Context, merchantID, verificationID string) (*BankAccountVerification, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, bank_code, account_number_masked, account_number_hash, account_name, verification_method, amount::text, connector_id, status, verification_response, beneficiary_name_returned, match_score::text, verified_at, expires_at FROM bank_account_verifications WHERE merchant_id=$1 AND id=$2`, merchantID, verificationID)
	var v BankAccountVerification
	var amount, matchScore, responseStr string
	err := row.Scan(&v.ID, &v.MerchantID, &v.BankCode, &v.AccountNumberMasked, &v.AccountNumberHash, &v.AccountName, &v.VerificationMethod, &amount, &v.ConnectorID, &v.Status, &responseStr, &v.BeneficiaryNameReturned, &matchScore, &v.VerifiedAt, &v.ExpiresAt)
	if err != nil {
		return nil, err
	}
	v.Amount, _ = decimal.NewFromString(amount)
	v.MatchScore, _ = decimal.NewFromString(matchScore)
	_ = json.Unmarshal([]byte(responseStr), &v.VerificationResponse)
	return &v, nil
}

func (r *PgRepository) ListBankVerifications(ctx context.Context, merchantID string, status *string) ([]BankAccountVerification, error) {
	query := `SELECT id, merchant_id, bank_code, account_number_masked, account_number_hash, account_name, verification_method, amount::text, connector_id, status FROM bank_account_verifications WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if status != nil {
		query += ` AND status=$2`
		args = append(args, *status)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []BankAccountVerification
	for rows.Next() {
		var v BankAccountVerification
		var amount string
		if err := rows.Scan(&v.ID, &v.MerchantID, &v.BankCode, &v.AccountNumberMasked, &v.AccountNumberHash, &v.AccountName, &v.VerificationMethod, &amount, &v.ConnectorID, &v.Status); err != nil {
			return nil, err
		}
		v.Amount, _ = decimal.NewFromString(amount)
		list = append(list, v)
	}
	return list, nil
}

func (r *PgRepository) UpdateBankVerificationStatus(ctx context.Context, verificationID, status string, beneficiaryName string, matchScore decimal.Decimal, response map[string]interface{}) error {
	respJSON, _ := json.Marshal(response)
	_, err := r.pool.Exec(ctx, `UPDATE bank_account_verifications SET status=$1, beneficiary_name_returned=$2, match_score=$3, verification_response=$4, verified_at=now() WHERE id=$5`,
		status, beneficiaryName, matchScore.String(), string(respJSON), verificationID)
	return err
}

// ==================== Virtual Accounts Smart Collect ====================

func (r *PgRepository) CreateVirtualAccount(ctx context.Context, va *VirtualAccount) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO virtual_accounts (id, merchant_id, virtual_account_number, customer_id, purpose, status, bank_code) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		va.ID, va.MerchantID, va.VirtualAccountNumber, va.CustomerID, va.Purpose, va.Status, va.BankCode)
	return err
}

func (r *PgRepository) ListVirtualAccounts(ctx context.Context, merchantID string) ([]VirtualAccount, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, virtual_account_number, customer_id, purpose, status, bank_code FROM virtual_accounts WHERE merchant_id=$1 ORDER BY created_at DESC`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []VirtualAccount
	for rows.Next() {
		var va VirtualAccount
		if err := rows.Scan(&va.ID, &va.MerchantID, &va.VirtualAccountNumber, &va.CustomerID, &va.Purpose, &va.Status, &va.BankCode); err != nil {
			return nil, err
		}
		list = append(list, va)
	}
	return list, nil
}

func (r *PgRepository) CreateVirtualAccountTransaction(ctx context.Context, txn *VirtualAccountTransaction) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO virtual_account_transactions (id, virtual_account_id, merchant_id, amount, currency, utr, sender_name, sender_account, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		txn.ID, txn.VirtualAccountID, txn.MerchantID, txn.Amount.String(), txn.Currency, txn.UTR, txn.SenderName, txn.SenderAccount, txn.Status)
	return err
}

func (r *PgRepository) ListVirtualAccountTransactions(ctx context.Context, merchantID, virtualAccountID string, status *string) ([]VirtualAccountTransaction, error) {
	query := `SELECT id, virtual_account_id, merchant_id, amount::text, currency, utr, sender_name, sender_account, status FROM virtual_account_transactions WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if virtualAccountID != "" {
		query += ` AND virtual_account_id=$2`
		args = append(args, virtualAccountID)
	}
	if status != nil {
		query += ` AND status=$3`
		args = append(args, *status)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []VirtualAccountTransaction
	for rows.Next() {
		var txn VirtualAccountTransaction
		var amount string
		if err := rows.Scan(&txn.ID, &txn.VirtualAccountID, &txn.MerchantID, &amount, &txn.Currency, &txn.UTR, &txn.SenderName, &txn.SenderAccount, &txn.Status); err != nil {
			return nil, err
		}
		txn.Amount, _ = decimal.NewFromString(amount)
		list = append(list, txn)
	}
	return list, nil
}
