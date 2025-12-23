package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	adapterhttp "github.com/maket12/SubTrack/internal/adapter/in/http"
	adapterdb "github.com/maket12/SubTrack/internal/adapter/out/db"
	"github.com/maket12/SubTrack/internal/app/usecase"
	"github.com/maket12/SubTrack/internal/config"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"log/slog"
	"os"
)

func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func main() {
	// ======================
	// 1. Load config
	// ======================
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ======================
	// 2. Setup logger
	// ======================

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	// ======================
	// 3. Connect to Postgres
	// ======================
	db, err := adapterdb.NewPostgres(cfg.DatabaseDSN)
	if err != nil {
		logger.Error("failed to connect db", slog.Any("err", err))
		os.Exit(1)
	}

	// ======================
	// 4. Repositories
	// ======================
	subRepo := adapterdb.NewSubscriptionRepo(db)

	// ======================
	// 5. Usecases
	// ======================
	createUC := &usecase.CreateSubscriptionUC{Subscriptions: subRepo}
	getUC := &usecase.GetSubscriptionUC{Subscriptions: subRepo}
	updateUC := &usecase.UpdateSubscriptionUC{Subscriptions: subRepo}
	deleteUC := &usecase.DeleteSubscriptionUC{Subscriptions: subRepo}
	listUC := &usecase.GetSubscriptionListUC{Subscriptions: subRepo}
	totalSumUC := &usecase.GetTotalSumUC{Subscriptions: subRepo}

	// ======================
	// 6. Handlers (REST)
	// ======================
	subHandler := adapterhttp.NewSubscriptionHandler(
		logger,
		createUC,
		getUC,
		updateUC,
		deleteUC,
		listUC,
		totalSumUC,
	)

	// ======================
	// 7. Router
	// ======================
	router := adapterhttp.NewRouter(subHandler).InitRoutes()

	router.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger.yaml"),
	))

	// ======================
	// 8. Run HTTP server
	// ======================

	srv := &http.Server{
		Addr:    cfg.HTTPAddress,
		Handler: router,
	}

	go func() {
		logger.Info("starting server", slog.String("address", cfg.HTTPAddress))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.Any("err", err))
	}

	if err := db.Close(); err != nil {
		logger.Error("failed to close database", slog.Any("err", err))
	}

	logger.Info("server exited properly")
}
