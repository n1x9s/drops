package postgres

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MemoryRepository struct {
	db *gorm.DB
}

func NewMemoryRepository(db *gorm.DB) MemoryRepository {
	return MemoryRepository{db: db}
}

func (r MemoryRepository) Create(ctx context.Context, memory domain.Memory, tagNames []string, embedding []float32) (domain.Memory, error) {
	model := MemoryModel{
		ID:       memory.ID,
		UserID:   memory.UserID,
		Content:  memory.Content,
		Summary:  memory.Summary,
		Category: string(memory.Category),
		Source:   memory.Source,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tags, err := ensureTags(ctx, tx, memory.UserID, tagNames)
		if err != nil {
			return err
		}
		model.Tags = tags
		if err := tx.Create(&model).Error; err != nil {
			return err
		}
		if len(embedding) > 0 {
			return tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "owner_type"}, {Name: "owner_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"embedding", "updated_at"}),
			}).Create(&EmbeddingModel{
				ID:        uuid.New(),
				UserID:    memory.UserID,
				OwnerType: "memory",
				OwnerID:   model.ID,
				Embedding: Vector(embedding),
			}).Error
		}
		return nil
	})
	if err != nil {
		return domain.Memory{}, err
	}
	return r.Get(ctx, memory.UserID, model.ID)
}

func (r MemoryRepository) Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (domain.Memory, error) {
	var model MemoryModel
	if err := r.db.WithContext(ctx).Preload("Tags").Where("user_id = ? AND id = ?", userID, id).First(&model).Error; err != nil {
		return domain.Memory{}, mapError(err)
	}
	return model.toDomain(), nil
}

func (r MemoryRepository) List(ctx context.Context, userID uuid.UUID, filter domain.MemoryFilter) ([]domain.Memory, error) {
	var models []MemoryModel
	query := r.db.WithContext(ctx).Preload("Tags").Where("memories.user_id = ?", userID).Order("memories.created_at DESC")
	if filter.Query != "" {
		needle := "%" + strings.ToLower(filter.Query) + "%"
		query = query.Where("LOWER(memories.content) LIKE ? OR LOWER(memories.summary) LIKE ?", needle, needle)
	}
	if filter.Category != "" {
		query = query.Where("memories.category = ?", filter.Category)
	}
	if filter.Tag != "" {
		query = query.Joins("JOIN memory_tags ON memory_tags.memory_id = memories.id").
			Joins("JOIN tags ON tags.id = memory_tags.tag_id").
			Where("tags.user_id = ? AND tags.name = ?", userID, strings.ToLower(filter.Tag))
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Memory, 0, len(models))
	for _, model := range models {
		out = append(out, model.toDomain())
	}
	return out, nil
}

func (r MemoryRepository) Update(ctx context.Context, memory domain.Memory, tagNames []string) (domain.Memory, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model MemoryModel
		if err := tx.Preload("Tags").Where("user_id = ? AND id = ?", memory.UserID, memory.ID).First(&model).Error; err != nil {
			return mapError(err)
		}
		model.Content = memory.Content
		model.Summary = memory.Summary
		model.Category = string(memory.Category)
		model.Source = memory.Source
		if err := tx.Save(&model).Error; err != nil {
			return err
		}
		if tagNames != nil {
			tags, err := ensureTags(ctx, tx, memory.UserID, tagNames)
			if err != nil {
				return err
			}
			if err := tx.Model(&model).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return domain.Memory{}, err
	}
	return r.Get(ctx, memory.UserID, memory.ID)
}

func (r MemoryRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND owner_type = ? AND owner_id = ?", userID, "memory", id).Delete(&EmbeddingModel{}).Error; err != nil {
			return err
		}
		result := tx.Where("user_id = ? AND id = ?", userID, id).Delete(&MemoryModel{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrNotFound
		}
		return nil
	})
}

func (r MemoryRepository) Similar(ctx context.Context, userID uuid.UUID, id uuid.UUID, limit int) ([]domain.SearchResult, error) {
	var source EmbeddingModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND owner_type = ? AND owner_id = ?", userID, "memory", id).First(&source).Error; err != nil {
		return nil, mapError(err)
	}
	var rows []searchRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT m.id, 'memory' AS type, COALESCE(NULLIF(m.summary, ''), LEFT(m.content, 120)) AS title,
		       LEFT(m.content, 240) AS snippet, 1 - (e.embedding <=> ?::vector) AS score,
		       m.category AS category, m.created_at AS occurred_at
		FROM embeddings e
		JOIN memories m ON m.id = e.owner_id
		WHERE e.user_id = ? AND e.owner_type = 'memory' AND e.owner_id <> ? AND m.deleted_at IS NULL
		ORDER BY e.embedding <=> ?::vector
		LIMIT ?`, Vector(source.Embedding), userID, id, Vector(source.Embedding), limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rowsToResults(rows), nil
}
