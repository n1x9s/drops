package main

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/n1x9s/second-brain/backend/internal/api/handlers"
	"github.com/n1x9s/second-brain/backend/internal/application/auth"
	"github.com/n1x9s/second-brain/backend/internal/application/integrations"
	"github.com/n1x9s/second-brain/backend/internal/application/memories"
	"github.com/n1x9s/second-brain/backend/internal/application/reminders"
	"github.com/n1x9s/second-brain/backend/internal/application/search"
	"github.com/n1x9s/second-brain/backend/internal/application/settings"
	"github.com/n1x9s/second-brain/backend/internal/application/tasks"
	"github.com/n1x9s/second-brain/backend/internal/config"
	"github.com/n1x9s/second-brain/backend/internal/infrastructure/ai"
	"github.com/n1x9s/second-brain/backend/internal/infrastructure/linear"
	"github.com/n1x9s/second-brain/backend/internal/infrastructure/postgres"
	"github.com/n1x9s/second-brain/backend/internal/infrastructure/security"
	"github.com/n1x9s/second-brain/backend/internal/infrastructure/telegram"
	platformdb "github.com/n1x9s/second-brain/backend/internal/platform/database"
	platformhttp "github.com/n1x9s/second-brain/backend/internal/platform/http"
	"github.com/n1x9s/second-brain/backend/internal/platform/logger"
	"github.com/n1x9s/second-brain/backend/internal/platform/metrics"
	"github.com/n1x9s/second-brain/backend/internal/platform/observability"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.Logging.Level, cfg.App.Env)
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	shutdownTracing, err := observability.InitTracing(ctx, cfg.Observability.OTLPEndpoint)
	if err != nil {
		log.Fatal("init tracing", zap.Error(err))
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = shutdownTracing(shutdownCtx)
	}()

	metrics.Register()

	gormDB, sqlDB, err := platformdb.Open(cfg.Database.DSN)
	if err != nil {
		log.Fatal("open database", zap.Error(err))
	}
	defer func() { _ = sqlDB.Close() }()

	userRepo := postgres.NewUserRepository(gormDB)
	memoryRepo := postgres.NewMemoryRepository(gormDB)
	taskRepo := postgres.NewTaskRepository(gormDB)
	reminderRepo := postgres.NewReminderRepository(gormDB)
	searchRepo := postgres.NewSearchRepository(gormDB)
	settingsRepo := postgres.NewSettingsRepository(gormDB)

	hasher := security.NewPasswordHasher()
	tokenManager := security.NewTokenManager(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	gemini := ai.NewGeminiProvider(cfg.Gemini.APIKey, cfg.Gemini.Model, cfg.Gemini.EmbeddingModel)
	telegramClient := telegram.NewClient(cfg.Telegram.BotToken)
	linearClient := linear.NewClient(cfg.Linear.APIKey)

	authSvc := auth.NewService(userRepo, hasher, tokenManager)
	memorySvc := memories.NewService(memoryRepo, gemini)
	taskSvc := tasks.NewService(taskRepo, gemini)
	reminderSvc := reminders.NewService(reminderRepo)
	searchSvc := search.NewService(searchRepo, gemini)
	settingsSvc := settings.NewService(settingsRepo)
	integrationSvc := integrations.NewService(settingsRepo, telegramClient, linearClient)

	router := platformhttp.NewRouter(log, tokenManager, platformhttp.Handlers{
		Health:       handlers.NewHealthHandler(sqlDB),
		Auth:         handlers.NewAuthHandler(authSvc),
		Memories:     handlers.NewMemoryHandler(memorySvc),
		Tasks:        handlers.NewTaskHandler(taskSvc),
		Reminders:    handlers.NewReminderHandler(reminderSvc),
		Search:       handlers.NewSearchHandler(searchSvc),
		Settings:     handlers.NewSettingsHandler(settingsSvc),
		Integrations: handlers.NewIntegrationHandler(integrationSvc),
	})

	server := &stdhttp.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info("api listening", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != stdhttp.ErrServerClosed {
			log.Fatal("listen", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown", zap.Error(err))
	}
}
