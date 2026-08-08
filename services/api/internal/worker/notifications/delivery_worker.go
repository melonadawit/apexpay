package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"apexpay/internal/notify"
)

// DeliveryWorker drains unread notifications and delivers them over each user's preferred
// channels (email/SMS), honoring notification_preferences. In-app delivery is implicit
// (the notification row is the inbox entry); this handles the push out to email/SMS.
type DeliveryWorker struct {
	pool   *pgxpool.Pool
	sender notify.Sender
}

func NewDeliveryWorker(pool *pgxpool.Pool, sender notify.Sender) *DeliveryWorker {
	return &DeliveryWorker{pool: pool, sender: sender}
}

// DeliverDue processes notifications created in the last interval for which the recipient
// has enabled the channel. Marks them read after a successful email send (best-effort).
func (w *DeliveryWorker) DeliverDue(ctx context.Context) (int, error) {
	rows, err := w.pool.Query(ctx, `
		SELECT n.id, n.merchant_id, COALESCE(n.user_id,''), n.type, n.title, n.message,
		       COALESCE(u.email,''), COALESCE(u.phone,'')
		FROM notifications n
		LEFT JOIN users u ON u.id = n.user_id
		WHERE n.is_read = false AND n.created_at >= now() - interval '5 minutes'
		ORDER BY n.created_at ASC LIMIT 100`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	delivered := 0
	for rows.Next() {
		var id, merchantID, userID, ntype, title, message, email, phone string
		if err := rows.Scan(&id, &merchantID, &userID, &ntype, &title, &message, &email, &phone); err != nil {
			continue
		}
		if userID != "" {
			if w.shouldEmail(ctx, merchantID, userID, ntype) && email != "" {
				_ = w.sender.SendEmail(ctx, email, "ApexPay: "+title, message)
				delivered++
			}
			if w.shouldSMS(ctx, merchantID, userID, ntype) && phone != "" {
				_ = w.sender.SendSMS(ctx, phone, title+": "+message)
				delivered++
			}
		}
		// Mark read so we don't re-send every 5s.
		_, _ = w.pool.Exec(ctx, `UPDATE notifications SET is_read=true, read_at=now() WHERE id=$1`, id)
	}
	return delivered, rows.Err()
}

func (w *DeliveryWorker) shouldEmail(ctx context.Context, merchantID, userID, ntype string) bool {
	return w.pref(ctx, merchantID, userID, ntype).Email
}

func (w *DeliveryWorker) shouldSMS(ctx context.Context, merchantID, userID, ntype string) bool {
	return w.pref(ctx, merchantID, userID, ntype).SMS
}

func (w *DeliveryWorker) pref(ctx context.Context, merchantID, userID, ntype string) notify.Preference {
	p := notify.Preference{EventType: ntype, Email: true, SMS: false, Push: true, InApp: true}
	_ = w.pool.QueryRow(ctx, `SELECT email, sms FROM notification_preferences WHERE merchant_id=$1 AND user_id=$2 AND event_type=$3`,
		merchantID, userID, ntype).Scan(&p.Email, &p.SMS)
	return p
}

// RunTicker polls for due notifications every 30s.
func (w *DeliveryWorker) RunTicker(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if n, err := w.DeliverDue(ctx); err == nil && n > 0 {
				fmt.Printf("[notify] delivered %d notifications\n", n)
			}
		}
	}
}
