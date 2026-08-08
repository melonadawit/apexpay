package notify

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// List returns the user's preferences for a merchant (defaults all on for unconfigured types).
func (r *Repository) List(ctx context.Context, merchantID, userID string) ([]Preference, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT event_type, email, sms, push, in_app FROM notification_preferences
		WHERE merchant_id=$1 AND user_id=$2`, merchantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	configured := map[string]bool{}
	var list []Preference
	for rows.Next() {
		var p Preference
		if err := rows.Scan(&p.EventType, &p.Email, &p.SMS, &p.Push, &p.InApp); err != nil {
			return nil, err
		}
		configured[p.EventType] = true
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Defaults for unconfigured event types.
	for _, et := range EventTypes {
		if !configured[et] {
			list = append(list, Preference{EventType: et, Email: true, SMS: false, Push: true, InApp: true})
		}
	}
	return list, nil
}

// Upsert sets a preference for a user.
func (r *Repository) Upsert(ctx context.Context, merchantID, userID string, p *Preference) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO notification_preferences (merchant_id, user_id, event_type, email, sms, push, in_app)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (merchant_id, user_id, event_type) DO UPDATE SET
			email=EXCLUDED.email, sms=EXCLUDED.sms, push=EXCLUDED.push, in_app=EXCLUDED.in_app`,
		merchantID, userID, p.EventType, p.Email, p.SMS, p.Push, p.InApp)
	return err
}
