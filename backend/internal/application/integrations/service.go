package integrations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
)

type Service struct {
	settings domain.SettingsRepository
	telegram domain.TelegramNotifier
	linear   domain.LinearClient
}

func NewService(settings domain.SettingsRepository, telegram domain.TelegramNotifier, linear domain.LinearClient) Service {
	return Service{settings: settings, telegram: telegram, linear: linear}
}

func (s Service) SendTelegramTest(ctx context.Context, userID uuid.UUID) error {
	cfg, err := s.settings.GetTelegram(ctx, userID)
	if err != nil {
		return err
	}
	if !cfg.Enabled {
		return domain.ErrProviderEmpty
	}
	return s.telegram.Send(ctx, cfg.ChatID, "Second Brain is connected.")
}

func (s Service) CreateLinearIssue(ctx context.Context, userID uuid.UUID, title string, description string) (domain.LinearIssue, error) {
	cfg, err := s.settings.GetLinear(ctx, userID)
	if err != nil {
		return domain.LinearIssue{}, err
	}
	if !cfg.Enabled {
		return domain.LinearIssue{}, domain.ErrProviderEmpty
	}
	return s.linear.CreateIssue(ctx, cfg.TeamID, title, description)
}

func (s Service) UpdateLinearIssue(ctx context.Context, userID uuid.UUID, issueID string, title string, description string, stateID string) (domain.LinearIssue, error) {
	cfg, err := s.settings.GetLinear(ctx, userID)
	if err != nil {
		return domain.LinearIssue{}, err
	}
	if !cfg.Enabled {
		return domain.LinearIssue{}, domain.ErrProviderEmpty
	}
	return s.linear.UpdateIssue(ctx, issueID, title, description, stateID)
}

func (s Service) ListLinearIssues(ctx context.Context, userID uuid.UUID) ([]domain.LinearIssue, error) {
	cfg, err := s.settings.GetLinear(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, fmt.Errorf("%w: linear", domain.ErrProviderEmpty)
	}
	return s.linear.ListIssues(ctx, cfg.TeamID)
}
