package reminders

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
)

type Service struct {
	reminders domain.ReminderRepository
}

type CreateInput struct {
	TaskID   *uuid.UUID
	MemoryID *uuid.UUID
	FireAt   time.Time
	Title    string
	Body     string
}

func NewService(reminders domain.ReminderRepository) Service {
	return Service{reminders: reminders}
}

func (s Service) Create(ctx context.Context, userID uuid.UUID, input CreateInput) (domain.Reminder, error) {
	if input.FireAt.IsZero() || strings.TrimSpace(input.Title) == "" {
		return domain.Reminder{}, domain.ErrInvalidInput
	}
	return s.reminders.Create(ctx, domain.Reminder{
		ID:       uuid.New(),
		UserID:   userID,
		TaskID:   input.TaskID,
		MemoryID: input.MemoryID,
		FireAt:   input.FireAt.UTC(),
		Title:    strings.TrimSpace(input.Title),
		Body:     strings.TrimSpace(input.Body),
		Status:   domain.ReminderPending,
	})
}

func (s Service) List(ctx context.Context, userID uuid.UUID, includeSent bool) ([]domain.Reminder, error) {
	return s.reminders.List(ctx, userID, includeSent)
}
