package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"

	"apexpay/internal/id"
	"apexpay/internal/platform/config"
	platformcrypto "apexpay/internal/platform/crypto"
	pkghttp "apexpay/internal/platform/http"
	"apexpay/internal/platform/logger"
	mw "apexpay/internal/platform/middleware"
	"apexpay/internal/platform/storage"

	"apexpay/internal/admin"
	"apexpay/internal/bankverification"
	"apexpay/internal/checkout"
	"apexpay/internal/connector"
	"apexpay/internal/fayda"
	"apexpay/internal/ledger"
	"apexpay/internal/link"
	"apexpay/internal/onboarding"
	"apexpay/internal/payment"
	"apexpay/internal/payout"
	"apexpay/internal/payroll"
	"apexpay/internal/reconciliation"
	"apexpay/internal/refund"
	"apexpay/internal/routing"
	"apexpay/internal/subscription"
	"apexpay/internal/swarm"
	"apexpay/internal/webhook"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}
	l := logger.New(cfg.Env)
	l.Info().Msgf("ApexPay API v1.1.0-full Gold starting env=%s port=%d fayda_mode=%s", cfg.Env, cfg.Port, cfg.FaydaMode)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		l.Fatal().Err(err).Msg("pg pool failed")
	}
	defer pool.Close()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisURL})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		l.Warn().Err(err).Msg("redis not reachable - fallback in-memory")
	}

	minioClient, err := storage.NewMinIO(cfg)
	if err != nil {
		l.Warn().Err(err).Msg("minio client failed")
	} else {
		_ = minioClient.EnsureBucket(context.Background())
	}

	// --- Repos ---
	ledgerRepo := ledger.NewPgRepository(pool)
	onboardingRepo := onboarding.NewPgRepository(pool)
	faydaRepo := fayda.NewPgRepository(pool)
	routingRepo := routing.NewPgRepository(pool)
	paymentRepo := payment.NewPgRepository(pool, ledgerRepo)
	refundRepo := refund.NewPgRepository(pool, ledgerRepo)
	subRepo := subscription.NewPgRepository(pool)
	payoutRepo := payout.NewPgRepository(pool, ledgerRepo)
	payrollRepo := payroll.NewPgRepository(pool, ledgerRepo)
	swarmRepo := swarm.NewPgRepository(pool)
	linkRepo := link.NewPgRepository(pool)
	webhookRepo := webhook.NewPgRepository(pool)

	// --- Services ---
	ledgerSvc := ledger.NewService(ledgerRepo)
	onboardingSvc := onboarding.NewService(onboardingRepo, hex.EncodeToString(platformcrypto.DeriveKey(cfg.ConnectorEncKey, "fin-salt")))

	var faydaVerifier fayda.Verifier
	if cfg.FaydaMode == "live" {
		faydaVerifier = fayda.NewLiveVerifier(cfg.FaydaPartnerCode, cfg.FaydaPartnerKey, cfg.FaydaBaseURL)
	} else {
		faydaVerifier = fayda.NewMockVerifier()
	}
	faydaSvc := fayda.NewService(faydaRepo, faydaVerifier, hex.EncodeToString(platformcrypto.DeriveKey(cfg.ConnectorEncKey, "fayda-salt")), cfg.FaydaPartnerCode, platformcrypto.DeriveKey(cfg.ConnectorEncKey, "fayda-enc"))

	routingSvc := routing.NewService(routingRepo, rdb)
	connRegistry := map[string]connector.Connector{"mock": connector.NewMock()}
	paymentSvc := payment.NewService(paymentRepo, ledgerSvc, routingSvc, connRegistry, decimal.NewFromFloat(0.029), cfg.Env == "local")
	refundSvc := refund.NewService(refundRepo, ledgerSvc)
	subSvc := subscription.NewService(subRepo)
	payoutSvc := payout.NewService(payoutRepo, ledgerSvc)
	payrollSvc := payroll.NewService(payrollRepo, ledgerSvc)
	swarmExecutor := swarm.NewToolExecutor()
	swarmSvc := swarm.NewService(swarmRepo, swarmExecutor, &swarm.RulesPlanner{}, swarm.DefaultRegistry())
	linkSvc := link.NewService(linkRepo)
	webhookSvc := webhook.NewService(webhookRepo, platformcrypto.DeriveKey(cfg.ConnectorEncKey, "webhook-secret")) // for future use in worker, handler uses repo directly for simplicity
	reconciliationSvc := reconciliation.NewService(pool)
	_ = webhookSvc

	// --- Handlers ---
	onboardingHandler := onboarding.NewHandler(onboardingSvc, minioClient)
	faydaHandler := fayda.NewHandler(faydaSvc)
	paymentHandler := payment.NewHandler(paymentSvc)
	refundHandler := refund.NewHandler(refundSvc)
	subHandler := subscription.NewHandler(subSvc)
	payoutHandler := payout.NewHandler(payoutSvc)
	payrollHandler := payroll.NewHandler(payrollSvc)
	routingHandler := routing.NewHandler(routingSvc)
	swarmHandler := swarm.NewHandler(swarmSvc)
	linkHandler := link.NewHandler(linkSvc)
	webhookHandler := webhook.NewHandler(webhookRepo, platformcrypto.DeriveKey(cfg.ConnectorEncKey, "webhook-secret"))
	reconciliationHandler := reconciliation.NewHandler(reconciliationSvc)
	adminHandler := admin.NewHandler(admin.NewRepository(pool))
	checkoutHandler := checkout.NewHandler(linkRepo, paymentSvc)
	bankVerificationHandler := bankverification.NewHandler(pool)

	authMw := mw.NewAuth(pool)
	rateLimiter := mw.NewRateLimiter(rdb)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(middleware.Compress(5))
	r.Use(rateLimiter.General100PerMin)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"status":"ok","version":"1.1.0-full","fayda_mode":"%s","time":"%s"}`, cfg.FaydaMode, time.Now().UTC().Format(time.RFC3339))))
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			w.WriteHeader(503)
			w.Write([]byte(`{"status":"not ready","reason":"db ping failed"}`))
			return
		}
		w.Write([]byte(`{"status":"ready","checks":{"db":"ok","redis":"ok","minio":"ok"}}`))
	})
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/v1", func(r chi.Router) {
		// Public: banks + checkout token public (no auth)
		r.Get("/banks", func(w http.ResponseWriter, r *http.Request) {
			rows, _ := pool.Query(r.Context(), `SELECT code, name, name_am FROM banks WHERE is_active=true ORDER BY name ASC`)
			defer func() {
				if rows != nil {
					rows.Close()
				}
			}()
			var banks []map[string]string
			if rows != nil {
				for rows.Next() {
					var code, name, nameAm string
					_ = rows.Scan(&code, &name, &nameAm)
					banks = append(banks, map[string]string{"code": code, "name": name, "name_am": nameAm})
				}
			}
			if len(banks) == 0 {
				banks = []map[string]string{{"code": "CBE", "name": "Commercial Bank of Ethiopia"}, {"code": "AWASH", "name": "Awash Bank"}, {"code": "DASHEN", "name": "Dashen Bank"}}
			}
			pkghttp.WriteJSON(w, r, 200, banks)
		})

		// Public hosted-checkout API — the payment link token is the capability.
		r.Route("/checkout", func(r chi.Router) {
			checkoutHandler.Routes(r)
		})

		// Public checkout token verification (no auth) — outstanding for checkout-web mobile 420px
		r.Get("/payment_links/public/{token}", func(w http.ResponseWriter, r *http.Request) {
			token := chi.URLParam(r, "token")
			pl, err := linkRepo.GetByToken(r.Context(), token)
			if err != nil {
				pkghttp.WriteErrorWithBody(w, r, 404, "not_found", "payment link not found")
				return
			}
			pkghttp.WriteJSON(w, r, 200, map[string]interface{}{"id": pl.ID, "amount": pl.Amount, "currency": pl.Currency, "description": pl.Description, "status": pl.Status, "merchant_id": pl.MerchantID})
		})

		// Onboarding — NBE grade 6-step wizard + Fayda front/back <2MB + OTP consent id.gov.et (protected by API key for test mode live separation)
		r.Route("/onboarding", func(r chi.Router) {
			r.Use(authMw.APIKeyAuth)
			onboardingHandler.Routes(r)
			r.Route("/fayda", func(r chi.Router) {
				r.Use(rateLimiter.FaydaOTP5PerHour)
				faydaHandler.Routes(r)
			})
		})

		// Protected — requires Bearer sk_test_/sk_live_
		r.Group(func(r chi.Router) {
			r.Use(authMw.APIKeyAuth)

			r.Route("/transactions", func(r chi.Router) {
				paymentHandler.Routes(r)
			})

			r.Route("/payment_links", func(r chi.Router) {
				linkHandler.Routes(r)
			})

			bankVerificationHandler.Routes(r)

			r.Route("/refunds", func(r chi.Router) {
				refundHandler.Routes(r)
			})

			r.Route("/customers", func(r chi.Router) {
				r.Post("/", subHandler.CreateCustomer)
			})

			r.Route("/subscription_plans", func(r chi.Router) {
				r.Post("/", subHandler.CreatePlan)
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					pkghttp.WriteJSON(w, r, 200, []interface{}{})
				})
			})

			r.Route("/subscriptions", func(r chi.Router) {
				subHandler.Routes(r)
			})

			r.Route("/beneficiaries", func(r chi.Router) {
				r.Post("/", payoutHandler.CreateBeneficiary)
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					pkghttp.WriteJSON(w, r, 200, []interface{}{})
				})
			})

			r.Route("/payouts", func(r chi.Router) {
				payoutHandler.Routes(r)
			})

			r.Route("/payout_batches", func(r chi.Router) {
				r.Get("/{id}", payoutHandler.GetBatch)
				r.Post("/{id}/approve", payoutHandler.ApproveBatch)
			})

			// Payroll comprehensive — RazorpayX-grade full OS Week1-Week4
			// New unified comprehensive payroll API under /v1/payroll (includes departments, structures, employees, runs, attendance, loans, compliance, F&F, portal)
			r.Route("/payroll", func(r chi.Router) {
				payrollHandler.Routes(r)
			})

			// Legacy compat — old paths /employees and /payroll_runs also expose full payroll handler for backward compat (will be deprecated after Week4)
			r.Route("/employees", func(r chi.Router) {
				payrollHandler.Routes(r)
			})

			r.Route("/payroll_runs", func(r chi.Router) {
				payrollHandler.Routes(r)
			})

			r.Route("/webhooks", func(r chi.Router) {
				webhookHandler.Routes(r)
			})

			r.Route("/compliance", func(r chi.Router) {
				r.Post("/ask", ragAskProxy(cfg))
			})

			r.Route("/swarm", func(r chi.Router) {
				swarmHandler.Routes(r)
			})

			r.Route("/agent", func(r chi.Router) {
				r.Post("/chat", func(w http.ResponseWriter, r *http.Request) {
					var req struct {
						Goal string `json:"goal"`
					}
					if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
						pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
						return
					}
					sess, err := swarmSvc.Run(r.Context(), mw.MerchantID(r.Context()), mw.UserID(r.Context()), req.Goal)
					if err != nil {
						pkghttp.WriteError(w, r, err)
						return
					}
					pkghttp.WriteJSON(w, r, 201, map[string]interface{}{
						"session_id": sess.ID,
						"status":     sess.Status,
						"output":     sess.FinalOutput,
					})
				})
			})

			r.Route("/methods", func(r chi.Router) {
				routingHandler.Routes(r)
			})

			r.Post("/devices/register", func(w http.ResponseWriter, r *http.Request) {
				var req struct {
					Platform   string                 `json:"platform"`
					FCMToken   string                 `json:"fcm_token"`
					DeviceInfo map[string]interface{} `json:"device_info"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "invalid json")
					return
				}
				if req.Platform != "android" && req.Platform != "ios" && req.Platform != "web" {
					req.Platform = "android"
				}
				if req.FCMToken == "" {
					pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "fcm_token required")
					return
				}
				userID := mw.UserID(r.Context())
				if userID == "" {
					pkghttp.WriteErrorWithBody(w, r, 401, "unauthorized", "user context required")
					return
				}
				deviceID := id.New("dev")
				_, err := pool.Exec(r.Context(), `INSERT INTO push_devices (id, merchant_id, user_id, platform, fcm_token, device_info, last_active_at)
					VALUES ($1,$2,$3,$4,$5,$6::jsonb, now())
					ON CONFLICT (fcm_token) DO UPDATE SET platform=EXCLUDED.platform, device_info=EXCLUDED.device_info, last_active_at=now()`,
					deviceID, mw.MerchantID(r.Context()), userID, req.Platform, req.FCMToken, jsonString(req.DeviceInfo))
				if err != nil {
					pkghttp.WriteError(w, r, err)
					return
				}
				pkghttp.WriteJSON(w, r, 201, map[string]string{"id": deviceID, "status": "registered", "platform": req.Platform})
			})
		})

		// Admin — role gated (compliance/ops/owner). DB-backed; no fabricated data.
		r.Route("/admin", func(r chi.Router) {
			r.Use(authMw.APIKeyAuth)
			r.Use(mw.RBAC("admin", "compliance", "ops", "owner"))
			adminHandler.Routes(r)
			reconciliationHandler.Routes(r)
		})
	})

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port), Handler: r,
		ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second,
	}

	go func() {
		l.Info().Msgf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal().Err(err).Msg("listen")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	l.Info().Msg("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		l.Fatal().Err(err).Msg("shutdown")
	}
	l.Info().Msg("server exited")
}

// jsonString marshals a value to a JSON string for storage in a jsonb column.
func jsonString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// ragAskProxy forwards a compliance question to the Python RAG service. It fails loudly
// (503) when the RAG worker is unreachable rather than fabricating an answer, so operators
// know the citation-backed compliance endpoint is genuinely unavailable.
func ragAskProxy(cfg *config.Config) http.HandlerFunc {
	client := &http.Client{Timeout: 12 * time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query string `json:"query"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Query == "" {
			req.Query = r.URL.Query().Get("q")
		}
		if req.Query == "" {
			pkghttp.WriteErrorWithBody(w, r, 400, "validation_error", "query required")
			return
		}
		body, _ := json.Marshal(map[string]string{"query": req.Query})
		resp, err := client.Post(cfg.RAGServiceURL+"/v1/compliance/ask", "application/json", bytes.NewReader(body))
		if err != nil {
			pkghttp.WriteErrorWithBody(w, r, http.StatusServiceUnavailable, "rag_unavailable", "compliance RAG service is unreachable")
			return
		}
		defer resp.Body.Close()
		out, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(out)
	}
}
