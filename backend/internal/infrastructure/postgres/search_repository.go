package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"gorm.io/gorm"
)

type SearchRepository struct {
	db *gorm.DB
}

type searchRow struct {
	ID         uuid.UUID
	Type       string
	Title      string
	Snippet    string
	Score      float64
	Category   string
	OccurredAt time.Time
}

func NewSearchRepository(db *gorm.DB) SearchRepository {
	return SearchRepository{db: db}
}

func (r SearchRepository) Search(ctx context.Context, userID uuid.UUID, query string, embedding []float32, limit int) ([]domain.SearchResult, error) {
	if len(embedding) > 0 {
		return r.vectorSearch(ctx, userID, embedding, limit)
	}
	return r.textSearch(ctx, userID, query, limit)
}

func (r SearchRepository) vectorSearch(ctx context.Context, userID uuid.UUID, embedding []float32, limit int) ([]domain.SearchResult, error) {
	var rows []searchRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT m.id, 'memory' AS type, COALESCE(NULLIF(m.summary, ''), LEFT(m.content, 120)) AS title,
		       LEFT(m.content, 240) AS snippet, 1 - (e.embedding <=> ?::vector) AS score,
		       m.category AS category, m.created_at AS occurred_at
		FROM embeddings e
		JOIN memories m ON m.id = e.owner_id
		WHERE e.user_id = ? AND e.owner_type = 'memory' AND m.deleted_at IS NULL
		UNION ALL
		SELECT t.id, 'task' AS type, t.title AS title, LEFT(t.notes, 240) AS snippet,
		       1 - (e.embedding <=> ?::vector) AS score, t.status AS category, t.created_at AS occurred_at
		FROM embeddings e
		JOIN tasks t ON t.id = e.owner_id
		WHERE e.user_id = ? AND e.owner_type = 'task' AND t.deleted_at IS NULL
		ORDER BY score DESC
		LIMIT ?`, Vector(embedding), userID, Vector(embedding), userID, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rowsToResults(rows), nil
}

func (r SearchRepository) textSearch(ctx context.Context, userID uuid.UUID, query string, limit int) ([]domain.SearchResult, error) {
	needle := "%" + strings.ToLower(query) + "%"
	var rows []searchRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, 'memory' AS type, COALESCE(NULLIF(summary, ''), LEFT(content, 120)) AS title,
		       LEFT(content, 240) AS snippet, 0.5 AS score, category, created_at AS occurred_at
		FROM memories
		WHERE user_id = ? AND deleted_at IS NULL AND (LOWER(content) LIKE ? OR LOWER(summary) LIKE ?)
		UNION ALL
		SELECT id, 'task' AS type, title, LEFT(notes, 240) AS snippet, 0.5 AS score, status AS category, created_at AS occurred_at
		FROM tasks
		WHERE user_id = ? AND deleted_at IS NULL AND (LOWER(title) LIKE ? OR LOWER(notes) LIKE ?)
		ORDER BY occurred_at DESC
		LIMIT ?`, userID, needle, needle, userID, needle, needle, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rowsToResults(rows), nil
}

func rowsToResults(rows []searchRow) []domain.SearchResult {
	results := make([]domain.SearchResult, 0, len(rows))
	for _, row := range rows {
		results = append(results, domain.SearchResult{
			ID:         row.ID,
			Type:       row.Type,
			Title:      row.Title,
			Snippet:    row.Snippet,
			Score:      row.Score,
			Category:   row.Category,
			OccurredAt: row.OccurredAt,
		})
	}
	return results
}
