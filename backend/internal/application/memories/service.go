package memories

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
)

type Service struct {
	memories domain.MemoryRepository
	llm      domain.LLMProvider
}

type CreateInput struct {
	Content string
	Source  string
}

type UpdateInput struct {
	ID       uuid.UUID
	Content  string
	Summary  string
	Category string
	Tags     []string
}

func NewService(memories domain.MemoryRepository, llm domain.LLMProvider) Service {
	return Service{memories: memories, llm: llm}
}

func (s Service) Create(ctx context.Context, userID uuid.UUID, input CreateInput) (domain.Memory, error) {
	content := strings.TrimSpace(input.Content)
	if content == "" {
		return domain.Memory{}, domain.ErrInvalidInput
	}

	enrichment := domain.LocalEnrichMemory(content)
	if s.llm != nil {
		if ai, err := s.llm.EnrichMemory(ctx, content); err == nil {
			if ai.Summary != "" {
				enrichment.Summary = ai.Summary
			}
			if ai.Category != "" {
				enrichment.Category = ai.Category
			}
			if len(ai.Tags) > 0 {
				enrichment.Tags = ai.Tags
			}
		}
	}

	var embedding []float32
	if s.llm != nil {
		if vector, err := s.llm.Embed(ctx, content); err == nil {
			embedding = vector
		}
	}

	source := input.Source
	if source == "" {
		source = "manual"
	}

	return s.memories.Create(ctx, domain.Memory{
		ID:       uuid.New(),
		UserID:   userID,
		Content:  content,
		Summary:  enrichment.Summary,
		Category: enrichment.Category,
		Source:   source,
	}, enrichment.Tags, embedding)
}

func (s Service) Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (domain.Memory, error) {
	return s.memories.Get(ctx, userID, id)
}

func (s Service) List(ctx context.Context, userID uuid.UUID, filter domain.MemoryFilter) ([]domain.Memory, error) {
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 50
	}
	return s.memories.List(ctx, userID, filter)
}

func (s Service) Update(ctx context.Context, userID uuid.UUID, input UpdateInput) (domain.Memory, error) {
	current, err := s.memories.Get(ctx, userID, input.ID)
	if err != nil {
		return domain.Memory{}, err
	}
	if strings.TrimSpace(input.Content) != "" {
		current.Content = strings.TrimSpace(input.Content)
	}
	if strings.TrimSpace(input.Summary) != "" {
		current.Summary = strings.TrimSpace(input.Summary)
	}
	if input.Category != "" {
		current.Category = domain.Category(input.Category)
	}
	return s.memories.Update(ctx, current, input.Tags)
}

func (s Service) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return s.memories.Delete(ctx, userID, id)
}

func (s Service) Similar(ctx context.Context, userID uuid.UUID, id uuid.UUID, limit int) ([]domain.SearchResult, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	return s.memories.Similar(ctx, userID, id, limit)
}
