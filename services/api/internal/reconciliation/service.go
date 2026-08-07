package reconciliation

import (
	"context"
	"fmt"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Case struct {
	MerchantID string `json:"merchant_id"`
	IdempotencyKey string `json:"idempotency_key"`
	TxRef string `json:"tx_ref"`
	RequestHash string `json:"request_hash"`
	ConnectorState string `json:"connector_state"`
	CreatedAt string `json:"created_at"`
	Status string `json:"status"`
	ReviewerID *string `json:"reviewer_id,omitempty"`
	ReviewerNote *string `json:"reviewer_note,omitempty"`
}

type Service struct{ pool *pgxpool.Pool }
func NewService(pool *pgxpool.Pool) *Service { return &Service{pool: pool} }

func (s *Service) ListOpen(ctx context.Context, limit int) ([]Case, error) {
	if limit < 1 || limit > 100 { limit = 50 }
	rows, err := s.pool.Query(ctx, `SELECT c.merchant_id,c.idempotency_key,COALESCE(c.tx_ref,''),i.request_hash,i.state,c.created_at::text,c.status,c.reviewer_id,c.reviewer_note
		FROM payment_reconciliation_cases c JOIN idempotency_keys i ON i.merchant_id=c.merchant_id AND i.key=c.idempotency_key
		WHERE c.status IN ('open','requires_connector_investigation') ORDER BY c.created_at ASC LIMIT $1`, limit)
	if err != nil { return nil, err }; defer rows.Close()
	var cases []Case
	for rows.Next() { var c Case; if err := rows.Scan(&c.MerchantID,&c.IdempotencyKey,&c.TxRef,&c.RequestHash,&c.ConnectorState,&c.CreatedAt,&c.Status,&c.ReviewerID,&c.ReviewerNote); err != nil{return nil,err}; cases=append(cases,c) }
	return cases, rows.Err()
}

func (s *Service) Decide(ctx context.Context, merchantID, key, decision, reviewerID, note string) error {
	if decision != "confirmed_paid" && decision != "confirmed_not_paid" && decision != "requires_connector_investigation" { return fmt.Errorf("invalid reconciliation decision") }
	tx, err := s.pool.Begin(ctx); if err != nil{return err}; defer func(){_ = tx.Rollback(ctx)}()
	command, err := tx.Exec(ctx, `UPDATE payment_reconciliation_cases SET status=$3,reviewer_id=$4,reviewer_note=$5,decided_at=now(),updated_at=now()
		WHERE merchant_id=$1 AND idempotency_key=$2 AND status IN ('open','requires_connector_investigation')`,merchantID,key,decision,reviewerID,note)
	if err != nil{return err}; if command.RowsAffected()!=1{return fmt.Errorf("open reconciliation case not found")}
	if decision == "confirmed_not_paid" { _,err=tx.Exec(ctx,`UPDATE idempotency_keys SET state='retry_authorized',response_code=409,response_body=COALESCE(response_body,'{}'::jsonb)||jsonb_build_object('resolution','confirmed_not_paid') WHERE merchant_id=$1 AND key=$2 AND state='manual_review'`,merchantID,key); if err!=nil{return err} }
	_,err=tx.Exec(ctx,`INSERT INTO audit_logs (id,merchant_id,actor_type,actor_id,action,resource_type,resource_id,data) VALUES ($1,$2,'operations',$3,$4,'payment_reconciliation',$5,jsonb_build_object('decision',$4,'note',$6))`,id.New("audit"),merchantID,reviewerID,decision,key,note)
	if err!=nil{return err}
	_,err=tx.Exec(ctx,`INSERT INTO outbox_events (id,merchant_id,aggregate_type,aggregate_id,event_type,payload) VALUES ($1,$2,'payment_reconciliation',$3,'payment.reconciliation.'||$4,jsonb_build_object('decision',$4,'idempotency_key',$3))`,id.NewOutbox(),merchantID,key,decision)
	if err!=nil{return err}; return tx.Commit(ctx)
}
