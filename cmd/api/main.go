package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/dubininme/xm-assessment/internal/config"
	deliveryHttp "github.com/dubininme/xm-assessment/internal/delivery/http"
	"github.com/dubininme/xm-assessment/internal/delivery/http/handler"
	"github.com/dubininme/xm-assessment/internal/delivery/http/middleware"
	"github.com/dubininme/xm-assessment/internal/domain/company"
	"github.com/dubininme/xm-assessment/internal/domain/events"
	"github.com/dubininme/xm-assessment/internal/infra/auth"
	"github.com/dubininme/xm-assessment/internal/infra/kafka"
	"github.com/dubininme/xm-assessment/internal/infra/outbox"
	"github.com/dubininme/xm-assessment/internal/infra/postgres"
	"github.com/dubininme/xm-assessment/pkg/logger"
)

func main() {
	log := logger.NewLogger()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Error("config initialization failed", "error", err)
		panic(err)
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Add logger to context
	ctx = logger.WithLogger(ctx, log)

	db, err := postgres.Connect(ctx, cfg.Db)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		panic(err)
	}
	defer func() { _ = db.Close() }()

	companyRepo := postgres.NewCompanyRepo(db)
	outboxRepo := postgres.NewOutboxRepo(db)
	txManager := postgres.NewTxManager(db)
	dbChecker := postgres.NewDBHealthChecker(db)

	kafkaProducer := kafka.NewProducer(cfg.Kafka.BrokersList(), cfg.Kafka.Topic)
	defer func() { _ = kafkaProducer.Close() }()

	outboxProcessor := outbox.NewProcessor(
		outboxRepo,
		kafkaProducer,
		txManager,
		cfg.Outbox.BatchSize,
		cfg.Outbox.Interval,
	)

	processorErrCh := make(chan error, 1)
	processorCtx, cancelProcessor := context.WithCancel(ctx)
	defer cancelProcessor()

	go func() {
		log.Info("starting outbox processor")
		if err := outboxProcessor.Start(processorCtx); err != nil && !errors.Is(err, context.Canceled) {
			processorErrCh <- fmt.Errorf("outbox processor failed: %w", err)
		}
		close(processorErrCh)
	}()

	httpHandler := initRouter(cfg, companyRepo, outboxRepo, txManager, dbChecker)
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpHandler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("starting HTTP server", "addr", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("server failed to start: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Error("server error", "error", err)
			cancelProcessor()
			panic(err)
		}
	case err := <-processorErrCh:
		if err != nil {
			log.Error("outbox processor error", "error", err)
			panic(err)
		}
	case <-ctx.Done():
		log.Info("shutting down server gracefully")

		// Stop outbox processor first
		cancelProcessor()
		if err := <-processorErrCh; err != nil {
			log.Error("outbox processor error during shutdown", "error", err)
		}

		// Then shutdown HTTP server
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("server forced shutdown", "error", err)
			panic(err)
		}

		if err := <-errCh; err != nil {
			log.Error("server error during shutdown", "error", err)
			panic(err)
		}

		log.Info("server exited gracefully")
		return
	}
}

func initRouter(cfg *config.AppConfig, companyRepo company.CompanyRepository, publisher events.EventsPublisher, txManager company.TxManager, dbChecker handler.HealthChecker) http.Handler {

	cService := company.NewCompanyService(companyRepo, publisher, txManager)
	cHandler := handler.NewCompanyHandler(cService)
	healthHandler := handler.NewHealthHandler(dbChecker)

	jwtService := auth.NewJWTService(cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(jwtService)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	router := deliveryHttp.NewRouter(cHandler, healthHandler, authHandler, authMiddleware)
	return router
}
