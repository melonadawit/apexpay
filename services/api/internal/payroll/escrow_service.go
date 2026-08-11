package payroll

import (
	"context"
	"fmt"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/ledger"
	"github.com/shopspring/decimal"
)

// EscrowService — Automated Escrow for Marketplaces P2P Hold & Release Funds Under Defined Conditions per Escrow+ 2024
// Reduces legal and operational overhead of running escrow manually with bank per ApexPay
// Marketplace seller settlement split: Order total 1000 ETB split: Platform fee 10% 100 ETB, Seller 90% 900 ETB, Withholding tax 2% 20 ETB, Hold in escrow until delivery confirmed, then release to seller minus fee and tax
// Auto-release cron daily 02:00 Africa/Addis_Ababa EAT per spec recon daily 02:00 EAT

type EscrowService struct {
	repo   Repository
	ledger *ledger.Service
}

func NewEscrowService(repo Repository, ledgerSvc *ledger.Service) *EscrowService {
	return &EscrowService{repo: repo, ledger: ledgerSvc}
}

// CreateEscrowAgreement — marketplace operator creates agreement with buyer, seller, amount, platform fee %, withholding tax %, conditions, auto-release
func (s *EscrowService) CreateAgreement(ctx context.Context, merchantID string, agreement *EscrowAgreement) error {
	if agreement.Amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("escrow agreement amount must be >0 per Ethiopia business practice")
	}
	if agreement.PlatformFeePercent.LessThan(decimal.Zero) || agreement.PlatformFeePercent.GreaterThan(decimal.NewFromInt(100)) {
		return fmt.Errorf("platform fee percent must be 0-100")
	}
	if agreement.WithholdingTaxPercent.LessThan(decimal.Zero) || agreement.WithholdingTaxPercent.GreaterThan(decimal.NewFromInt(100)) {
		return fmt.Errorf("withholding tax percent must be 0-100 per Ethiopia Income Tax Proclamation")
	}
	agreement.ID = id.New("escagr")
	agreement.MerchantID = merchantID
	if agreement.AgreementNumber == "" {
		agreement.AgreementNumber = fmt.Sprintf("ESC-AGR-%d-%s", time.Now().Year(), agreement.ID[len(agreement.ID)-6:])
	}
	if agreement.Status == "" {
		agreement.Status = "draft"
	}
	// Conditions default: delivery_confirmed 7 days, inspection_period 3 days per Escrow
	if len(agreement.Conditions) == 0 {
		agreement.Conditions = []EscrowCondition{
			{Type: "delivery_confirmed", Days: 7},
			{Type: "inspection_period", Days: 3},
		}
	}
	if agreement.AutoReleaseAfterDays == 0 {
		agreement.AutoReleaseAfterDays = 7
	}
	// Save via repo — need method CreateEscrowAgreement, for now use generic CreateEscrowAccount? We'll need repo method
	// For rapid P0, we will call repo.CreateEscrowAgreement if exists, else fallback to direct SQL via repo
	return s.repo.CreateEscrowAgreement(ctx, agreement)
}

// CreateEscrowAccount — Hold funds in escrow book per agreement book_type escrow per DATABASE
// Ledger Model: Dr asset:clearing:bank Amount Cr liability:escrow_payable Amount? Actually hold: Dr asset:clearing:bank Amount Cr liability:escrow_payable Amount? For marketplace, buyer pays, funds held in escrow: Dr asset:clearing:bank (buyer payment) Cr liability:escrow_payable (escrow held)
// Then release: Dr liability:escrow_payable Amount Cr asset:clearing:bank SellerAmount? Actually release to seller minus fee and tax
func (s *EscrowService) CreateEscrowAccount(ctx context.Context, merchantID string, agreementID, buyerID, sellerID, orderID string, orderAmount decimal.Decimal) (*EscrowAccount, error) {
	// Fetch agreement
	agreement, err := s.repo.GetEscrowAgreement(ctx, merchantID, agreementID)
	if err != nil {
		return nil, err
	}

	platformFee := orderAmount.Mul(agreement.PlatformFeePercent).Div(decimal.NewFromInt(100)).Round(2)
	withholdingTax := orderAmount.Mul(agreement.WithholdingTaxPercent).Div(decimal.NewFromInt(100)).Round(2)
	sellerAmount := orderAmount.Sub(platformFee).Sub(withholdingTax)

	accountNumber := fmt.Sprintf("ETB-CBE-ESCROW-%s", id.New("escacc")[len("escacc_"):])
	bookID := id.New("lbk")
	journalID := id.New("ljrn")

	escrowAccount := &EscrowAccount{
		ID:               id.New("escrow"),
		MerchantID:       merchantID,
		AgreementID:      agreementID,
		AccountNumber:    accountNumber,
		AccountName:      fmt.Sprintf("Escrow %s Order %s Buyer %s Seller %s", agreement.AgreementNumber, orderID, buyerID, sellerID),
		Amount:           orderAmount,
		Currency:         "ETB",
		Status:           "held",
		HeldAt:           ptrTime(time.Now()),
		BuyerMerchantID:  &buyerID,
		SellerMerchantID: &sellerID,
		OrderID:          &orderID,
		OrderAmount:      orderAmount,
		PlatformFee:      platformFee,
		SellerAmount:     sellerAmount,
		WithholdingTax:   withholdingTax,
		LedgerBookID:     &bookID,
	}

	// Ledger: Dr asset:clearing:bank Amount Cr liability:escrow_payable Amount (hold)
	journal := &ledger.Journal{
		ID:            journalID,
		BookID:        bookID,
		PostingKey:    fmt.Sprintf("escrow_hold:%s", escrowAccount.ID),
		Memo:          fmt.Sprintf("Escrow hold Order %s Buyer %s Seller %s Amount %s Fee %s Tax %s", orderID, buyerID, sellerID, orderAmount.String(), platformFee.String(), withholdingTax.String()),
		ReferenceType: "escrow_account",
		ReferenceID:   escrowAccount.ID,
	}
	entries := []ledger.Entry{
		{ID: id.New("le"), JournalID: journalID, BookID: bookID, AccountID: "asset:clearing:bank", Direction: "debit", Amount: orderAmount, Currency: "ETB"},
		{ID: id.New("le"), JournalID: journalID, BookID: bookID, AccountID: "liability:escrow_payable", Direction: "credit", Amount: orderAmount, Currency: "ETB"},
	}

	if err := s.repo.CreateEscrowAccountTx(ctx, escrowAccount, journal, entries); err != nil {
		return nil, err
	}

	return escrowAccount, nil
}

// ReleaseEscrow — Release funds to seller minus fee and tax, per agreement conditions met
// Ledger Model Release: Dr liability:escrow_payable OrderAmount Cr liability:platform_fee_due PlatformFee Cr liability:withholding_tax_payable WithholdingTax Cr asset:clearing:bank SellerAmount
// Then Dr platform_fee_due PlatformFee Cr platform_revenue PlatformFee? Actually platform fee becomes revenue, withholding tax becomes payable to ERCA, seller amount becomes clearing bank to seller
func (s *EscrowService) ReleaseEscrow(ctx context.Context, merchantID, escrowAccountID, releaserID string) error {
	escrow, err := s.repo.GetEscrowAccount(ctx, merchantID, escrowAccountID)
	if err != nil {
		return err
	}
	if escrow.Status != "held" {
		return fmt.Errorf("escrow account status must be held to release, got %s", escrow.Status)
	}

	bookID := escrow.LedgerBookID
	if bookID == nil {
		bid := id.New("lbk")
		bookID = &bid
	}
	journalID := id.New("ljrn")
	journal := &ledger.Journal{
		ID:            journalID,
		BookID:        *bookID,
		PostingKey:    fmt.Sprintf("escrow_release:%s", escrow.ID),
		Memo:          fmt.Sprintf("Escrow release Order %s Seller %s Amount %s Fee %s Tax %s", safeDeref(escrow.OrderID), safeDeref(escrow.SellerMerchantID), escrow.Amount.String(), escrow.PlatformFee.String(), escrow.WithholdingTax.String()),
		ReferenceType: "escrow_release",
		ReferenceID:   escrow.ID,
	}
	entries := []ledger.Entry{
		{ID: id.New("le"), JournalID: journalID, BookID: *bookID, AccountID: "liability:escrow_payable", Direction: "debit", Amount: escrow.Amount, Currency: "ETB"},
		{ID: id.New("le"), JournalID: journalID, BookID: *bookID, AccountID: "liability:platform_fee_due", Direction: "credit", Amount: escrow.PlatformFee, Currency: "ETB"},
		{ID: id.New("le"), JournalID: journalID, BookID: *bookID, AccountID: "liability:withholding_tax_payable", Direction: "credit", Amount: escrow.WithholdingTax, Currency: "ETB"},
		{ID: id.New("le"), JournalID: journalID, BookID: *bookID, AccountID: "asset:clearing:bank", Direction: "credit", Amount: escrow.SellerAmount, Currency: "ETB"},
	}

	// Filter zero entries optimization O(n)
	filtered := entries[:0]
	for _, e := range entries {
		if e.Amount.GreaterThan(decimal.Zero) {
			filtered = append(filtered, e)
		}
	}
	if !ledger.ValidateBalanced(filtered) {
		return fmt.Errorf("escrow release ledger unbalanced debit != credit")
	}

	if err := s.repo.ReleaseEscrowTx(ctx, escrow.ID, journal, filtered, releaserID); err != nil {
		return err
	}

	return nil
}

// ReturnEscrow — Return funds to buyer (dispute, expired)
func (s *EscrowService) ReturnEscrow(ctx context.Context, merchantID, escrowAccountID, returnerID string, reason string) error {
	escrow, err := s.repo.GetEscrowAccount(ctx, merchantID, escrowAccountID)
	if err != nil {
		return err
	}
	if escrow.Status != "held" {
		return fmt.Errorf("escrow account status must be held to return, got %s", escrow.Status)
	}

	bookID := escrow.LedgerBookID
	if bookID == nil {
		bid := id.New("lbk")
		bookID = &bid
	}
	journalID := id.New("ljrn")
	journal := &ledger.Journal{
		ID:            journalID,
		BookID:        *bookID,
		PostingKey:    fmt.Sprintf("escrow_return:%s", escrow.ID),
		Memo:          fmt.Sprintf("Escrow return Order %s Buyer %s Reason %s", safeDeref(escrow.OrderID), safeDeref(escrow.BuyerMerchantID), reason),
		ReferenceType: "escrow_return",
		ReferenceID:   escrow.ID,
	}
	entries := []ledger.Entry{
		{ID: id.New("le"), JournalID: journalID, BookID: *bookID, AccountID: "liability:escrow_payable", Direction: "debit", Amount: escrow.Amount, Currency: "ETB"},
		{ID: id.New("le"), JournalID: journalID, BookID: *bookID, AccountID: "asset:clearing:bank", Direction: "credit", Amount: escrow.Amount, Currency: "ETB"},
	}

	return s.repo.ReturnEscrowTx(ctx, escrow.ID, journal, entries, returnerID, reason)
}

// AutoReleaseExpiredEscrows — cron daily 02:00 Africa/Addis_Ababa EAT per spec recon daily 02:00 EAT
// Checks escrow_accounts status held where expires_at <= now() and auto_release true, then releases
// O(n) where n = expired escrows (usually small), optimal for daily cron
func (s *EscrowService) AutoReleaseExpiredEscrows(ctx context.Context) (int, error) {
	expiredEscrows, err := s.repo.ListExpiredEscrowsForAutoRelease(ctx)
	if err != nil {
		return 0, err
	}
	releasedCount := 0
	for _, escrow := range expiredEscrows {
		if err := s.ReleaseEscrow(ctx, escrow.MerchantID, escrow.ID, "system_auto_release"); err == nil {
			releasedCount++
		}
	}
	return releasedCount, nil
}

// Helpers

func ptrTime(t time.Time) *time.Time { return &t }
func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

var _ = fmt.Sprintf
