package main

import (
	"github.com/maket12/SubTrack/internal/adapter/in/http"
	adapterdb "github.com/maket12/SubTrack/internal/adapter/out/db"
	"github.com/maket12/SubTrack/internal/app/usecase"

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
	// 1. Setup logger
	// ======================
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(logLevel),
	}))

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		logger.Error("DATABASE_DSN env is missing")
		os.Exit(1)
	}

	// ======================
	// 2. Connect to Postgres
	// ======================
	db, err := adapterdb.NewPostgres(dsn)
	if err != nil {
		logger.Error("failed to connect db: %v", err)
		os.Exit(1)
	}

	// ======================
	// 3. Repositories
	// ======================
	subRepo := adapterdb.NewSubscriptionRepo(db)

	// ======================
	// 4. Usecases
	// ======================
	createUC := &usecase.CreateSubscriptionUC{Subscriptions: subRepo}
	getUC := &usecase.GetSubscriptionUC{Subscriptions: subRepo}
	updateUC := &usecase.UpdateSubscriptionUC{Subscriptions: subRepo}
	deleteUC := &usecase.DeleteSubscriptionUC{Subscriptions: subRepo}
	listUC := &usecase.GetSubscriptionListUC{Subscriptions: subRepo}
	totalSumUC := &usecase.GetTotalSumUC{Subscriptions: subRepo}

	// ======================
	// 5. Handlers (REST)
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
	// 6. Router
	// ======================
	router := http.NewRouter(subHandler).InitRoutes()

	router.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger.yaml"),
	))

	// ======================
	// 7. Run HTTP server
	// ======================
	httpAddr := os.Getenv("HTTP_ADDRESS")
	if httpAddr == "" {
		logger.Error("HTTP_ADDRESS env is missing")
		os.Exit(1)
	}

	logger.Info("starting server on %s\n", httpAddr)
	if err := router.Run(httpAddr); err != nil {
		logger.Error("failed to run HTTP server: %v", err)
		os.Exit(1)
	}
}
