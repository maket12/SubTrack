package main

import (
	"log"

	"github.com/maket12/SubTrack/internal/adapter/in/http"
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
	subHandler := http.NewSubscriptionHandler(
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
	router := http.NewRouter(subHandler).InitRoutes()

	router.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger.yaml"),
	))

	// ======================
	// 8. Run HTTP server
	// ======================
	logger.Info("starting server", slog.String("address", cfg.HTTPAddress))
	if err := router.Run(cfg.HTTPAddress); err != nil {
		logger.Error("failed to run HTTP server", slog.Any("err", err))
		os.Exit(1)
	}
}
