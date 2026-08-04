package webhook

import (
	"context"
	"time"

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
	Attempt    int
}

func (r *PgRepository) ListPendingDeliveries(ctx context.Context, limit int) ([]DeliveryRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT wd.id, wd.merchant_id, wd.endpoint_id, wd.event_type, wd.payload, we.url, we.secret_hash, wd.attempt_count
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
	var list []DeliveryRow
	for rows.Next() {
		var d DeliveryRow
		if err := rows.Scan(&d.ID, &d.MerchantID, &d.EndpointID, &d.EventType, &d.Payload, &d.URL, &d.Secret, &d.Attempt); err != nil {
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
		SELECT gen_random_ulid_text(), $1, we.id, $2, $3, $4, 'pending', now()
		FROM webhook_endpoints we
		WHERE we.merchant_id=$1 AND we.status='active'
	`, outbox.MerchantID, outbox.ID, outbox.EventType, outbox.Payload)
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
