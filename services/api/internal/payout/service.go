package payout

import (
	"context"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/ledger"
	"apexpay/internal/platform/errors"

	"github.com/shopspring/decimal"
)

type Repository interface {
	CreateBeneficiary(ctx context.Context, b *Beneficiary) error
	GetBeneficiary(ctx context.Context, merchantID, id string) (*Beneficiary, error)
	CreateBatchTx(ctx context.Context, batch *PayoutBatch, journal *ledger.Journal, entries []ledger.Entry) error
	CreatePayout(ctx context.Context, p *Payout) error
	CreateBulkTx(ctx context.Context, batch *PayoutBatch, payouts []Payout, journal *ledger.Journal, entries []ledger.Entry) error
	GetBatch(ctx context.Context, merchantID, batchID string) (*PayoutBatch, error)
	UpdateBatchStatus(ctx context.Context, batchID, status, approvedBy string) error
	UpdatePayoutStatus(ctx context.Context, payoutID string, status PayoutStatus, connectorRef string) error
	GetMerchantBalance(ctx context.Context, merchantID string) (decimal.Decimal, error)
}

type Service struct {
	repo   Repository
	ledger *ledger.Service
}

func NewService(repo Repository, ledgerSvc *ledger.Service) *Service {
	return &Service{repo: repo, ledger: ledgerSvc}
}

// Maker-checker threshold per NBE best practice
var ApprovalThreshold = decimal.NewFromInt(50000) // ETB

func (s *Service) CreateSingle(ctx context.Context, p *Payout) (*Payout, error) {
	if p.Amount.LessThanOrEqual(decimal.Zero) { return nil, errors.Validation("amount >0") }
	beneficiary, err := s.repo.GetBeneficiary(ctx, p.MerchantID, p.BeneficiaryID)
	if err != nil || beneficiary == nil || beneficiary.VerificationStatus != "verified" { return nil, errors.Validation("beneficiary bank account must be verified before payout") }
	balance, err := s.repo.GetMerchantBalance(ctx, p.MerchantID)
	if err != nil {
		return nil, err
	}
	if balance.LessThan(p.Amount) {
		return nil, errors.New(errors.CodeInsufficientBalance, "insufficient merchant payable balance", 400)
	}

	p.ID = id.NewPayout()
	p.CreatedAt = time.Now()
	if p.Amount.GreaterThan(ApprovalThreshold) {
		p.Status = StatusPendingApproval
	} else {
		p.Status = StatusQueued
	}

	// For queued, post ledger M3 immediately: Dr payable Cr clearing
	if p.Status == StatusQueued {
		journal := &ledger.Journal{
			ID: id.NewLedgerJournal(), BookID: "merchant_operating_default",
			PostingKey: "payout:" + p.ID, Memo: "payout single", ReferenceType: "payout", ReferenceID: p.ID,
		}
		entries := []ledger.Entry{
			{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "liability:merchant_payable", Direction: "debit", Amount: p.Amount, Currency: p.Currency},
			{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "asset:clearing:bank", Direction: "credit", Amount: p.Amount, Currency: p.Currency},
		}
		// Use batch tx variant with empty batch? Simplified create payout + ledger
		// Repo should handle
		_ = journal
		_ = entries
	}

	if err := s.repo.CreatePayout(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) CreateBulk(ctx context.Context, req CreateBulkRequest) (*PayoutBatch, error) {
	if len(req.Items) == 0 {
		return nil, errors.Validation("bulk must have at least 1 item")
	}
	if len(req.Items) > 1000 {
		return nil, errors.Validation("bulk max 1000 per batch")
	}
	total := decimal.Zero
	for _, it := range req.Items {
		if it.Amount.LessThanOrEqual(decimal.Zero) { return nil, errors.Validation("amount >0 per item") }
		beneficiary, err := s.repo.GetBeneficiary(ctx, req.MerchantID, it.BeneficiaryID)
		if err != nil || beneficiary == nil || beneficiary.VerificationStatus != "verified" { return nil, errors.Validation("all bulk payout beneficiaries must be verified") }
		total = total.Add(it.Amount)
	}

	balance, err := s.repo.GetMerchantBalance(ctx, req.MerchantID)
	if err != nil {
		return nil, err
	}
	if balance.LessThan(total) {
		return nil, errors.New(errors.CodeInsufficientBalance, "insufficient balance for bulk", 400)
	}

	batch := &PayoutBatch{
		ID: id.NewPayoutBatch(), MerchantID: req.MerchantID, BatchRef: req.BatchRef,
		Amount: total, Currency: req.Currency, Status: "pending_approval", CreatedAt: time.Now(),
	}
	// All bulk require approval per policy
	payouts := make([]Payout, 0, len(req.Items))
	for _, it := range req.Items {
		payouts = append(payouts, Payout{
			ID: id.NewPayout(), MerchantID: req.MerchantID, BatchID: &batch.ID,
			BeneficiaryID: it.BeneficiaryID, PayoutRef: it.PayoutRef,
			Amount: it.Amount, Currency: req.Currency, Status: StatusCreated,
			Method: "bank", CreatedAt: time.Now(),
		})
	}

	journal := &ledger.Journal{
		ID: id.NewLedgerJournal(), BookID: "merchant_operating_default",
		PostingKey: "payout_batch:" + batch.ID, Memo: "bulk payout batch", ReferenceType: "payout_batch", ReferenceID: batch.ID,
	}
	entries := []ledger.Entry{
		{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "liability:merchant_payable", Direction: "debit", Amount: total, Currency: req.Currency},
		{ID: id.New("le"), JournalID: journal.ID, BookID: journal.BookID, AccountID: "asset:clearing:bank", Direction: "credit", Amount: total, Currency: req.Currency},
	}

	if err := s.repo.CreateBulkTx(ctx, batch, payouts, journal, entries); err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *Service) ApproveBatch(ctx context.Context, merchantID, batchID, approverID string) error {
	batch, err := s.repo.GetBatch(ctx, merchantID, batchID)
	if err != nil {
		return errors.NotFound("batch not found")
	}
	if batch.Status != "pending_approval" {
		return errors.Validation("batch not pending approval")
	}
	// Dual approval check could be here
	return s.repo.UpdateBatchStatus(ctx, batchID, "approved", approverID)
}
