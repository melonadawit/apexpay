package webhook

import (
	"context"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct{ pool *pgxpool.Pool }

func NewPgRepository(pool *pgxpool.Pool) *PgRepository { return &PgRepository{pool: pool} }

type DeliveryRow struct {
	ID         string
	MerchantID string
	EndpointID string
	EventType  string
	Payload    []byte
	URL        string
	Secret     string
	EncryptedSecret []byte
	Attempt    int
}


// PublishOutbox creates delivery rows and marks each event published in one transaction.
func (r *PgRepository) PublishOutbox(ctx context.Context, limit int) (int, error) {
	tx, err := r.pool.Begin(ctx); if err != nil { return 0, err }; defer func(){ _ = tx.Rollback(ctx) }()
	rows, err := tx.Query(ctx, `SELECT id,merchant_id,event_type,payload FROM outbox_events WHERE published_at IS NULL ORDER BY created_at LIMIT $1 FOR UPDATE SKIP LOCKED`, limit)
	if err != nil{return 0,err}; defer rows.Close(); count:=0
	for rows.Next(){ var eventID,merchantID,eventType string; var payload []byte; if err:=rows.Scan(&eventID,&merchantID,&eventType,&payload);err!=nil{return 0,err}
		eps,err:=tx.Query(ctx,`SELECT id FROM webhook_endpoints WHERE merchant_id=$1 AND status='active' AND (events ? $2 OR events ? '*')`,merchantID,eventType);if err!=nil{return 0,err}
		for eps.Next(){var endpointID string;if err:=eps.Scan(&endpointID);err!=nil{eps.Close();return 0,err};_,err=tx.Exec(ctx,`INSERT INTO webhook_deliveries (id,merchant_id,endpoint_id,outbox_event_id,event_type,payload,status,next_attempt_at) VALUES ($1,$2,$3,$4,$5,$6,'pending',now())`,id.New("wd"),merchantID,endpointID,eventID,eventType,payload);if err!=nil{eps.Close();return 0,err}}
		eps.Close(); if _,err=tx.Exec(ctx,`UPDATE outbox_events SET published_at=now() WHERE id=$1`,eventID);err!=nil{return 0,err};count++
	}
	if err:=rows.Err();err!=nil{return 0,err};return count,tx.Commit(ctx)
}

func (r *PgRepository) ListPendingDeliveries(ctx context.Context, limit int) ([]Delivery, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT wd.id, wd.merchant_id, wd.endpoint_id, wd.event_type, wd.payload, we.url, we.secret_hash, we.secret_encrypted, wd.attempt_count
		FROM webhook_deliveries wd
		JOIN webhook_endpoints we ON we.id = wd.endpoint_id
		WHERE wd.status IN ('pending','failed') AND (wd.next_attempt_at IS NULL OR wd.next_attempt_at <= now())
		ORDER BY wd.created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Delivery
	for rows.Next() {
		var d Delivery
		if err := rows.Scan(&d.ID, &d.MerchantID, &d.EndpointID, &d.EventType, &d.Payload, &d.URL, &d.Secret, &d.EncryptedSecret, &d.AttemptCount); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *PgRepository) MarkSuccess(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE webhook_deliveries SET status='success', attempt_count=attempt_count+1, updated_at=now() WHERE id=$1`, id)
	return err
}

func (r *PgRepository) MarkFailed(ctx context.Context, id string, statusCode int, errMsg string, nextAttempt time.Time) error {
	// After 10 attempts mark dead
	_, err := r.pool.Exec(ctx, `
		UPDATE webhook_deliveries
		SET status=CASE WHEN attempt_count>=10 THEN 'dead' ELSE 'failed' END,
		    attempt_count=attempt_count+1,
		    last_status_code=$2,
		    last_error=$3,
		    next_attempt_at=$4,
		    updated_at=now()
		WHERE id=$1
	`, id, statusCode, errMsg, nextAttempt)
	return err
}

func (r *PgRepository) ListPendingOutbox(ctx context.Context, limit int) ([]OutboxRow, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, aggregate_type, aggregate_id, event_type, payload FROM outbox_events WHERE published_at IS NULL ORDER BY created_at ASC LIMIT $1 FOR UPDATE SKIP LOCKED`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []OutboxRow
	for rows.Next() {
		var o OutboxRow
		var payload []byte
		if err := rows.Scan(&o.ID, &o.MerchantID, &o.AggregateType, &o.AggregateID, &o.EventType, &payload); err != nil {
			return nil, err
		}
		o.Payload = payload
		list = append(list, o)
	}
	return list, nil
}

func (r *PgRepository) MarkOutboxPublished(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE outbox_events SET published_at=now() WHERE id=$1`, id)
	return err
}

func (r *PgRepository) CreateDeliveriesForOutbox(ctx context.Context, outbox OutboxRow) error {
	// For each active endpoint for merchant, create delivery
	_, err := r.pool.Exec(ctx, `
		INSERT INTO webhook_deliveries (id, merchant_id, endpoint_id, outbox_event_id, event_type, payload, status, next_attempt_at)
		SELECT $5 || we.id, $1, we.id, $2, $3, $4, 'pending', now()
		FROM webhook_endpoints we
		WHERE we.merchant_id=$1 AND we.status='active'
	`, outbox.MerchantID, outbox.ID, outbox.EventType, outbox.Payload, "wd_")
	return err
}

type OutboxRow struct {
	ID            string
	MerchantID    string
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte
}
