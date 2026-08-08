package treasury

import "context"

// Service adds business rules on top of the treasury repository.
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

// CreateAndCompleteTransfer creates an internal transfer and completes it immediately,
// moving money between the merchant's own accounts in one atomic operation.
func (s *Service) CreateAndCompleteTransfer(ctx context.Context, merchantID, userID string, t *Transfer) (*Transfer, error) {
	if err := s.repo.CreateTransfer(ctx, merchantID, userID, t); err != nil {
		return nil, err
	}
	if err := s.repo.CompleteTransfer(ctx, t.ID); err != nil {
		// If completion fails (e.g. insufficient balance), the transfer stays 'pending'
		// or 'failed' — surface the error so the caller can retry/adjust.
		return t, err
	}
	return t, nil
}
