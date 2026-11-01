package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"lending-backend-deep/internal/config"
	openapihandler "lending-backend-deep/internal/openapi"
	healthhandler "lending-backend-deep/internal/rest/handler/healthcheck"
	loanhandler "lending-backend-deep/internal/rest/handler/loanapplication"
	masterdatahandler "lending-backend-deep/internal/rest/handler/masterdata"
	"lending-backend-deep/internal/rest/middleware"
	pgstore "lending-backend-deep/internal/store/pg"
	"lending-backend-deep/internal/worker"
	"lending-backend-deep/resources/gen"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// NewRouter wires the chi router and registers all routes and middlewares.
func NewRouter(cfg *config.Env, logger zerolog.Logger, pool *pgxpool.Pool) http.Handler {
	// -----------------------------------------------------------------
	// Dependencies
	// -----------------------------------------------------------------
	partnerStore := pgstore.NewPartnerStore(pool)
	applicationStore := pgstore.NewLoanApplicationStore(pool)
	masterDataStore := pgstore.NewMasterDataStore(pool)

	// -----------------------------------------------------------------
	// Router
	// -----------------------------------------------------------------
	router := chi.NewRouter()
	router.Use(
		chimw.RequestID,                  // injects X-Request-Id header
		middleware.ZerologLogger(logger), // structured request logging
		middleware.PanicRecovery(logger), // catches panics → JSON 500
	)

	// Health probe — no api-key required
	router.Get("/healthcheck", healthhandler.Healthcheck())

	// Partner routes — api-key middleware is scoped to this group only
	router.Group(func(r chi.Router) {
		r.Use(middleware.PartnerAuth(logger, partnerStore))
		r.Post("/loan-applications/v1", loanhandler.Submit(logger, applicationStore))
		r.Get("/loan-applications/v1/{application-reference}/validation-result",
			loanhandler.ValidationResult(logger, applicationStore))
		r.Get("/master-data/v1/enumerations", masterdatahandler.Enumerations(logger, masterDataStore))
	})

	// OpenAPI spec + Redoc UI (opt-in via OPENAPI_UI_ENABLE=true)
	if cfg.OpenAPIUIEnable {
		router.Get("/openapi.json", openapihandler.SpecHandler(gen.OpenAPISpec, "openapi/openapi.json"))
		router.Get(cfg.OpenAPIUIPath, openapihandler.UIHandler())
		logger.Info().Msgf("redoc ui available at %s%s", cfg.ServerAddr(), cfg.OpenAPIUIPath)
	}

	return router
}

// Serve starts the HTTP server and the background workers.  It blocks until a
// SIGINT or SIGTERM is received and then performs a graceful shutdown.
func Serve(cfg *config.Env, logger zerolog.Logger, pool *pgxpool.Pool) error {
	router := NewRouter(cfg, logger, pool)

	location, err := time.LoadLocation(cfg.BusinessTimezone)
	if err != nil {
		return fmt.Errorf("server.Serve: load BUSINESS_TIMEZONE: %w", err)
	}

	// -----------------------------------------------------------------
	// Background workers — they stop after the HTTP server, once no new
	// work can arrive.
	// -----------------------------------------------------------------
	workerCtx, stopWorkers := context.WithCancel(context.Background())
	defer stopWorkers()

	notificationStore := pgstore.NewRejectionNotificationStore(pool)
	validator := worker.NewValidator(logger,
		pgstore.NewLoanApplicationStore(pool), pgstore.NewCustomerStore(),
		pgstore.NewMasterDataStore(pool), notificationStore, location)
	notifier := worker.NewNotifier(logger, notificationStore, cfg.AppEnv != "local")

	var workers sync.WaitGroup
	workers.Go(func() { validator.Run(workerCtx, cfg.ValidatorPollInterval()) })
	workers.Go(func() { notifier.Run(workerCtx, cfg.NotifierPollInterval()) })

	// -----------------------------------------------------------------
	// HTTP server
	// -----------------------------------------------------------------
	srv := &http.Server{
		Addr:         cfg.ServerAddr(),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout(),
		WriteTimeout: cfg.WriteTimeout(),
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info().Msgf("http server listening on %s", cfg.ServerAddr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()

	// -----------------------------------------------------------------
	// Graceful shutdown
	// -----------------------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-quit:
		logger.Info().Msgf("signal %s received, shutting down gracefully", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout())
	defer cancel()

	err = srv.Shutdown(ctx)
	stopWorkers()
	workers.Wait()
	if err != nil {
		return err
	}

	logger.Info().Msg("http server stopped")
	return nil
}
