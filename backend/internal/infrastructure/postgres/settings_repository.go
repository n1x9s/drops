package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return SettingsRepository{db: db}
}

func (r SettingsRepository) Get(ctx context.Context, userID uuid.UUID) (domain.UserSettings, error) {
	var model UserSettingsModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model).Error; err != nil {
		if mapError(err) == domain.ErrNotFound {
			return r.Upsert(ctx, domain.UserSettings{
				UserID:               userID,
				GeminiEnabled:        true,
				NotificationsEnabled: true,
				SiriEnabled:          true,
			})
		}
		return domain.UserSettings{}, err
	}
	return model.toDomain(), nil
}

func (r SettingsRepository) Upsert(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error) {
	model := UserSettingsModel{
		UserID:               settings.UserID,
		GeminiEnabled:        settings.GeminiEnabled,
		TelegramEnabled:      settings.TelegramEnabled,
		LinearEnabled:        settings.LinearEnabled,
		NotificationsEnabled: settings.NotificationsEnabled,
		SiriEnabled:          settings.SiriEnabled,
	}
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"gemini_enabled", "telegram_enabled", "linear_enabled", "notifications_enabled", "siri_enabled", "updated_at"}),
	}).Create(&model).Error; err != nil {
		return domain.UserSettings{}, err
	}
	return r.Get(ctx, settings.UserID)
}

func (r SettingsRepository) UpsertTelegram(ctx context.Context, cfg domain.TelegramConfig) (domain.TelegramConfig, error) {
	model := TelegramConfigModel{UserID: cfg.UserID, BotToken: cfg.BotToken, ChatID: cfg.ChatID, Enabled: cfg.Enabled}
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"bot_token", "chat_id", "enabled", "updated_at"}),
	}).Create(&model).Error; err != nil {
		return domain.TelegramConfig{}, err
	}
	return r.GetTelegram(ctx, cfg.UserID)
}

func (r SettingsRepository) GetTelegram(ctx context.Context, userID uuid.UUID) (domain.TelegramConfig, error) {
	var model TelegramConfigModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model).Error; err != nil {
		return domain.TelegramConfig{}, mapError(err)
	}
	return domain.TelegramConfig{UserID: model.UserID, BotToken: model.BotToken, ChatID: model.ChatID, Enabled: model.Enabled, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}, nil
}

func (r SettingsRepository) UpsertLinear(ctx context.Context, cfg domain.LinearConfig) (domain.LinearConfig, error) {
	model := LinearConfigModel{UserID: cfg.UserID, APIKey: cfg.APIKey, TeamID: cfg.TeamID, Enabled: cfg.Enabled}
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"api_key", "team_id", "enabled", "updated_at"}),
	}).Create(&model).Error; err != nil {
		return domain.LinearConfig{}, err
	}
	return r.GetLinear(ctx, cfg.UserID)
}

func (r SettingsRepository) GetLinear(ctx context.Context, userID uuid.UUID) (domain.LinearConfig, error) {
	var model LinearConfigModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model).Error; err != nil {
		return domain.LinearConfig{}, mapError(err)
	}
	return domain.LinearConfig{UserID: model.UserID, APIKey: model.APIKey, TeamID: model.TeamID, Enabled: model.Enabled, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}, nil
}
