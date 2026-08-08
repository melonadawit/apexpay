package dispute

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"apexpay/internal/id"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// Create opens a dispute for a payment.
func (r *Repository) Create(ctx context.Context, merchantID, userID string, d *Dispute) error {
	d.ID = id.New("disp")
	if d.Currency == "" {
		d.Currency = "ETB"
	}
	if d.Status == "" {
		d.Status = "open"
	}
	ev, _ := json.Marshal(d.Evidence)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO disputes (id, merchant_id, payment_id, amount, currency, reason_code, status, evidence, created_by)
		VALUES ($1,$2,$3,$4::numeric,$5,$6,$7,$8::jsonb,$9)`,
		d.ID, merchantID, d.PaymentID, d.Amount, d.Currency, d.ReasonCode, d.Status, string(ev), userID)
	return err
}

// List returns disputes for a merchant.
func (r *Repository) List(ctx context.Context, merchantID, status string, limit int) ([]Dispute, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query := `SELECT id, COALESCE(payment_id,''), amount::text, currency, reason_code, status, evidence::text,
		COALESCE(resolution,''), COALESCE(fee,0)::text,
		to_char(created_at AT TIME ZONE 'Africa/Addis_Ababa','YYYY-MM-DD"T"HH24:MI:SS')
		FROM disputes WHERE merchant_id=$1`
	args := []interface{}{merchantID}
	if status != "" {
		query += ` AND status=$2`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC LIMIT $` + itoa(len(args)+1)
	args = append(args, limit)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Dispute{}
	for rows.Next() {
		var d Dispute
		var ev string
		if err := rows.Scan(&d.ID, &d.PaymentID, &d.Amount, &d.Currency, &d.ReasonCode, &d.Status,
			&ev, &d.Resolution, &d.Fee, &d.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(ev), &d.Evidence)
		list = append(list, d)
	}
	return list, rows.Err()
}

// SubmitEvidence appends evidence and moves status to evidence_submitted.
func (r *Repository) SubmitEvidence(ctx context.Context, merchantID, disputeID string, items []EvidenceItem) error {
	// Load current evidence, append, and update.
	var cur string
	err := r.pool.QueryRow(ctx, `SELECT COALESCE(evidence,'[]')::text FROM disputes WHERE merchant_id=$1 AND id=$2`, merchantID, disputeID).Scan(&cur)
	if err != nil {
		return ErrNotFound
	}
	var existing []EvidenceItem
	_ = json.Unmarshal([]byte(cur), &existing)
	existing = append(existing, items...)
	ev, _ := json.Marshal(existing)
	_, err = r.pool.Exec(ctx, `UPDATE disputes SET evidence=$1::jsonb, status='evidence_submitted', updated_at=now() WHERE merchant_id=$2 AND id=$3`,
		string(ev), merchantID, disputeID)
	return err
}

// Decide resolves a dispute (won/lost).
func (r *Repository) Decide(ctx context.Context, merchantID, disputeID, decision, resolution string, fee string) error {
	if decision != "won" && decision != "lost" {
		return errors.New("decision must be won or lost")
	}
	ct, err := r.pool.Exec(ctx, `UPDATE disputes SET status=$1, resolution=$2, fee=$3::numeric, decided_at=now(), updated_at=now()
		WHERE merchant_id=$4 AND id=$5 AND status IN ('open','evidence_submitted')`,
		decision, resolution, fee, merchantID, disputeID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func itoa(n int) string { return strconv.Itoa(n) }

func nowStr() string        { return time.Now().In(tzEAT()).Format(time.RFC3339) }
func tzEAT() *time.Location { return time.FixedZone("EAT", 3*3600) }
