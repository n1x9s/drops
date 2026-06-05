package tasks

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
)

type Service struct {
	tasks domain.TaskRepository
	llm   domain.LLMProvider
}

type CreateInput struct {
	Text     string
	Title    string
	Notes    string
	DueAt    *time.Time
	Priority string
	Tags     []string
}

type UpdateInput struct {
	ID       uuid.UUID
	Title    string
	Notes    string
	DueAt    *time.Time
	Priority string
	Status   string
	Tags     []string
}

func NewService(tasks domain.TaskRepository, llm domain.LLMProvider) Service {
	return Service{tasks: tasks, llm: llm}
}

func (s Service) Create(ctx context.Context, userID uuid.UUID, input CreateInput) (domain.Task, error) {
	extraction := domain.TaskExtraction{
		Title:    strings.TrimSpace(input.Title),
		DueAt:    input.DueAt,
		Priority: normalizePriority(input.Priority),
		Tags:     input.Tags,
	}
	if input.Text != "" {
		extraction = domain.LocalExtractTask(input.Text)
		if s.llm != nil {
			if ai, err := s.llm.ExtractTask(ctx, input.Text); err == nil {
				if ai.Title != "" {
					extraction.Title = ai.Title
				}
				if ai.DueAt != nil {
					extraction.DueAt = ai.DueAt
				}
				if ai.Priority != "" {
					extraction.Priority = ai.Priority
				}
				if len(ai.Tags) > 0 {
					extraction.Tags = ai.Tags
				}
			}
		}
	}
	if extraction.Title == "" {
		return domain.Task{}, domain.ErrInvalidInput
	}
	if extraction.Priority == "" {
		extraction.Priority = PriorityOrDefault(input.Priority)
	}

	var embedding []float32
	if s.llm != nil {
		if vector, err := s.llm.Embed(ctx, extraction.Title+" "+input.Notes); err == nil {
			embedding = vector
		}
	}

	return s.tasks.Create(ctx, domain.Task{
		ID:       uuid.New(),
		UserID:   userID,
		Title:    extraction.Title,
		Notes:    strings.TrimSpace(input.Notes),
		Priority: extraction.Priority,
		Status:   domain.TaskStatusInbox,
		DueAt:    extraction.DueAt,
	}, extraction.Tags, embedding)
}

func (s Service) Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (domain.Task, error) {
	return s.tasks.Get(ctx, userID, id)
}

func (s Service) List(ctx context.Context, userID uuid.UUID, filter domain.TaskFilter) ([]domain.Task, error) {
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 50
	}
	return s.tasks.List(ctx, userID, filter)
}

func (s Service) Update(ctx context.Context, userID uuid.UUID, input UpdateInput) (domain.Task, error) {
	current, err := s.tasks.Get(ctx, userID, input.ID)
	if err != nil {
		return domain.Task{}, err
	}
	if strings.TrimSpace(input.Title) != "" {
		current.Title = strings.TrimSpace(input.Title)
	}
	if input.Notes != "" {
		current.Notes = strings.TrimSpace(input.Notes)
	}
	if input.DueAt != nil {
		current.DueAt = input.DueAt
	}
	if input.Priority != "" {
		current.Priority = normalizePriority(input.Priority)
	}
	if input.Status != "" {
		current.Status = domain.TaskStatus(input.Status)
	}
	return s.tasks.Update(ctx, current, input.Tags)
}

func (s Service) Complete(ctx context.Context, userID uuid.UUID, id uuid.UUID) (domain.Task, error) {
	now := time.Now().UTC()
	task, err := s.tasks.Get(ctx, userID, id)
	if err != nil {
		return domain.Task{}, err
	}
	task.Status = domain.TaskStatusCompleted
	task.CompletedAt = &now
	return s.tasks.Update(ctx, task, nil)
}

func (s Service) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return s.tasks.Delete(ctx, userID, id)
}

func PriorityOrDefault(value string) domain.Priority {
	priority := normalizePriority(value)
	if priority == "" {
		return domain.PriorityMedium
	}
	return priority
}

func normalizePriority(value string) domain.Priority {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "low":
		return domain.PriorityLow
	case "high":
		return domain.PriorityHigh
	case "urgent":
		return domain.PriorityUrgent
	case "medium", "":
		return domain.PriorityMedium
	default:
		return domain.PriorityMedium
	}
}
