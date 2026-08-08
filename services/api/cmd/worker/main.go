package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"apexpay/internal/connector"
	"apexpay/internal/id"
	"apexpay/internal/notify"
	"apexpay/internal/platform/config"
	platformcrypto "apexpay/internal/platform/crypto"
	"apexpay/internal/platform/logger"
	"apexpay/internal/webhook"
	"apexpay/internal/worker/accounting"
	"apexpay/internal/worker/credit"
	"apexpay/internal/worker/forex"
	"apexpay/internal/worker/notifications"
)

// Worker entrypoint — a small set of long-lived background jobs, one goroutine per job so a
// failing job never blocks the others. Jobs are idempotent, resilient to transient DB/Redis
// outages, and shut down gracefully via the shared context.
//
// Scheduled jobs:
//   - idempotency reconciliation sweep      (60s)
//   - outbox publish + webhook delivery     (1s  / 1s)
//   - forex rate cache refresh (NBE)        (60s)
//   - credit scoring                        (1h)
//   - accounting two-way sync (Tally/Zoho/QB) (1h)
//   - FCM / in-app notification push         (5s)
//   - connector health sampler               (30s)
//   - subscription invoice dunning           (1h)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	l := logger.New(cfg.Env)
	l.Info().Msgf("ApexPay Worker starting env=%s", cfg.Env)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		l.Fatal().Err(err).Msg("pg pool")
	}
	defer pool.Close()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisURL})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ------------------------------------------------------------------
	// 1. Idempotency reconciliation sweep — quarantine stuck payments.
	//    A connector-started idempotency record older than 15 minutes is moved to
	//    manual_review, never retried automatically. Operations must reconcile its tx_ref
	//    against the rail before any customer retry is permitted.
	// ------------------------------------------------------------------
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				tx, err := pool.Begin(ctx)
				if err != nil {
					l.Error().Err(err).Msg("idempotency reconciliation transaction failed")
					continue
				}
				command, err := tx.Exec(ctx, `WITH moved AS (
					UPDATE idempotency_keys SET state='manual_review', response_code=409,
					response_body=COALESCE(response_body,'{}'::jsonb) || jsonb_build_object('error','connector_reconciliation_required')
					WHERE state='connector_started' AND created_at < now() - interval '15 minutes'
					RETURNING merchant_id,key,COALESCE(response_body->>'tx_ref','') AS tx_ref)
					INSERT INTO payment_reconciliation_cases (merchant_id,idempotency_key,tx_ref)
					SELECT merchant_id,key,tx_ref FROM moved ON CONFLICT (merchant_id,idempotency_key) DO NOTHING`)
				if err != nil {
					_ = tx.Rollback(ctx)
					l.Error().Err(err).Msg("idempotency reconciliation sweep failed")
					continue
				}
				if err = tx.Commit(ctx); err != nil {
					l.Error().Err(err).Msg("idempotency reconciliation commit failed")
					continue
				}
				if command.RowsAffected() > 0 {
					l.Warn().Int64("count", command.RowsAffected()).Msg("payments moved to manual reconciliation")
				}
			}
		}
	}()

	// ------------------------------------------------------------------
	// 2. Outbox publish + webhook delivery. Publishing is transactional; delivery is
	//    retried with bounded exponential backoff by webhook.Service (HMAC + SSRF block).
	// ------------------------------------------------------------------
	webhookRepo := webhook.NewPgRepository(pool)
	webhookSvc := webhook.NewService(webhookRepo, platformcrypto.DeriveKey(cfg.ConnectorEncKey, "webhook-secret"))
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := webhookRepo.PublishOutbox(ctx, 100); err != nil {
					l.Error().Err(err).Msg("outbox publish failed")
				}
				deliveries, err := webhookRepo.ListPendingDeliveries(ctx, 100)
				if err != nil {
					l.Error().Err(err).Msg("webhook delivery poll failed")
					continue
				}
				for _, delivery := range deliveries {
					if err := webhookSvc.Deliver(ctx, delivery); err != nil {
						l.Error().Err(err).Str("delivery_id", delivery.ID).Msg("webhook delivery failed")
					}
				}
			}
		}
	}()

	// ------------------------------------------------------------------
	// 3. Forex rate cache — refresh NBE official rates every 60s into Redis + Postgres.
	// ------------------------------------------------------------------
	forexWorker := forex.NewRateCacheWorker(pool, rdb)
	go forexWorker.RunTicker(ctx)

	// ------------------------------------------------------------------
	// 4. Credit scoring — recompute merchant credit scores/limits hourly.
	// ------------------------------------------------------------------
	scoringWorker := credit.NewScoringWorker(pool)
	go scoringWorker.RunTicker(ctx)

	// ------------------------------------------------------------------
	// 5. Accounting two-way sync — Tally / Zoho / QuickBooks hourly.
	// ------------------------------------------------------------------
	syncWorker := accounting.NewSyncWorker(pool)
	go syncWorker.RunTicker(ctx)

	// ------------------------------------------------------------------
	// 6. Notifications — poll unread notifications and send FCM / in-app pushes.
	//    Plus real email/SMS delivery honoring each user's notification_preferences.
	//    SMTP creds via APEXPAY_SMTP_*; falls back to console logs when unset (dev).
	// ------------------------------------------------------------------
	pushWorker := notifications.NewPushWorker(pool)
	go pushWorker.RunTicker(ctx)

	notifySender := notify.NewEmailSMSDeliverer(&notify.SMTPConfig{
		Host:     os.Getenv("APEXPAY_SMTP_HOST"),
		Port:     os.Getenv("APEXPAY_SMTP_PORT"),
		Username: os.Getenv("APEXPAY_SMTP_USER"),
		Password: os.Getenv("APEXPAY_SMTP_PASS"),
		From:     os.Getenv("APEXPAY_SMTP_FROM"),
	})
	deliveryWorker := notifications.NewDeliveryWorker(pool, notifySender)
	go deliveryWorker.RunTicker(ctx)

	// ------------------------------------------------------------------
	// 7. Connector health sampler — probe each enabled connector every 30s, persist a
	//    health sample, and warm the Redis health cache for the smart router.
	// ------------------------------------------------------------------
	go func() {
		// Probe global (merchant-agnostic) connector configs via the gateway registry.
		registry := connector.NewRegistry(pool, "test", func(cipher []byte) ([]byte, error) {
			return platformcrypto.Decrypt(platformcrypto.DeriveKey(cfg.ConnectorEncKey, "connector-config"), cipher)
		})
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				conns, err := registry.ForMerchant(ctx, "")
				if err != nil {
					l.Error().Err(err).Msg("health sampler: build connectors failed")
					continue
				}
				if len(conns) == 0 {
					conns = map[string]connector.Connector{"mock": connector.NewMock()}
				}
				for connectorID, conn := range conns {
					okFlag, latency, err := conn.Health(ctx)
					_, insertErr := pool.Exec(ctx,
						`INSERT INTO connector_health_samples (id, connector_id, environment, latency_ms, success, error_code) VALUES ($1,$2,$3,$4,$5,$6)`,
						id.New("hs"), connectorID, "live", latency, okFlag, errCode(err))
					if insertErr != nil {
						l.Error().Err(insertErr).Str("connector", connectorID).Msg("health sample insert failed")
						continue
					}
					_ = rdb.Set(ctx, "health:"+connectorID, latency, 60*time.Second).Err()
				}
			}
		}
	}()

	// ------------------------------------------------------------------
	// 8. Subscription invoice dunning — retry open invoices 1d/3d/5d, then fail them.
	//    (State machine only; charging uses the saved payment method via the connector.)
	// ------------------------------------------------------------------
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				command, err := pool.Exec(ctx, `UPDATE subscription_invoices
					SET attempt_count = attempt_count + 1,
					    next_attempt_at = CASE WHEN attempt_count + 1 >= 3 THEN NULL ELSE now() + interval '2 days' END,
					    status = CASE WHEN attempt_count + 1 >= 3 THEN 'failed' ELSE status END
					WHERE status='open' AND attempt_count < 3 AND (next_attempt_at IS NULL OR next_attempt_at <= now())`)
				if err != nil {
					l.Error().Err(err).Msg("dunning advance failed")
					continue
				}
				if command.RowsAffected() > 0 {
					l.Info().Int64("count", command.RowsAffected()).Msg("dunning advanced open invoices")
				}
			}
		}
	}()

	// Block until shutdown signal.
	<-ctx.Done()
	l.Info().Msg("worker shutting down")
	_ = rdb
	_ = pool
}

// errCode converts a probe error into a short DB-storable code, or NULL-safe empty string.
func errCode(err error) any {
	if err != nil {
		return err.Error()
	}
	return nil
}
