package search

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
)

type Service struct {
	search domain.SearchRepository
	llm    domain.LLMProvider
}

func NewService(search domain.SearchRepository, llm domain.LLMProvider) Service {
	return Service{search: search, llm: llm}
}

func (s Service) Search(ctx context.Context, userID uuid.UUID, query string, limit int) ([]domain.SearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, domain.ErrInvalidInput
	}
	if limit <= 0 || limit > 100 {
		limit = 25
	}

	var embedding []float32
	if s.llm != nil {
		if vector, err := s.llm.Embed(ctx, query); err == nil {
			embedding = vector
		}
	}
	return s.search.Search(ctx, userID, query, embedding, limit)
}
