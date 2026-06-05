package settings

import (
	"context"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
)

type Service struct {
	settings domain.SettingsRepository
}

func NewService(settings domain.SettingsRepository) Service {
	return Service{settings: settings}
}

func (s Service) Get(ctx context.Context, userID uuid.UUID) (domain.UserSettings, error) {
	return s.settings.Get(ctx, userID)
}

func (s Service) Upsert(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error) {
	return s.settings.Upsert(ctx, settings)
}

func (s Service) UpsertTelegram(ctx context.Context, cfg domain.TelegramConfig) (domain.TelegramConfig, error) {
	return s.settings.UpsertTelegram(ctx, cfg)
}

func (s Service) GetTelegram(ctx context.Context, userID uuid.UUID) (domain.TelegramConfig, error) {
	return s.settings.GetTelegram(ctx, userID)
}

func (s Service) UpsertLinear(ctx context.Context, cfg domain.LinearConfig) (domain.LinearConfig, error) {
	return s.settings.UpsertLinear(ctx, cfg)
}

func (s Service) GetLinear(ctx context.Context, userID uuid.UUID) (domain.LinearConfig, error) {
	return s.settings.GetLinear(ctx, userID)
}
