package postgres

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return TaskRepository{db: db}
}

func (r TaskRepository) Create(ctx context.Context, task domain.Task, tagNames []string, embedding []float32) (domain.Task, error) {
	model := TaskModel{
		ID:            task.ID,
		UserID:        task.UserID,
		Title:         task.Title,
		Notes:         task.Notes,
		Priority:      string(task.Priority),
		Status:        string(task.Status),
		DueAt:         task.DueAt,
		CompletedAt:   task.CompletedAt,
		LinearIssueID: task.LinearIssueID,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tags, err := ensureTags(ctx, tx, task.UserID, tagNames)
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
				UserID:    task.UserID,
				OwnerType: "task",
				OwnerID:   model.ID,
				Embedding: Vector(embedding),
			}).Error
		}
		return nil
	})
	if err != nil {
		return domain.Task{}, err
	}
	return r.Get(ctx, task.UserID, model.ID)
}

func (r TaskRepository) Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (domain.Task, error) {
	var model TaskModel
	if err := r.db.WithContext(ctx).Preload("Tags").Where("user_id = ? AND id = ?", userID, id).First(&model).Error; err != nil {
		return domain.Task{}, mapError(err)
	}
	return model.toDomain(), nil
}

func (r TaskRepository) List(ctx context.Context, userID uuid.UUID, filter domain.TaskFilter) ([]domain.Task, error) {
	var models []TaskModel
	query := r.db.WithContext(ctx).Preload("Tags").Where("tasks.user_id = ?", userID).Order("tasks.created_at DESC")
	if filter.Query != "" {
		needle := "%" + strings.ToLower(filter.Query) + "%"
		query = query.Where("LOWER(tasks.title) LIKE ? OR LOWER(tasks.notes) LIKE ?", needle, needle)
	}
	if filter.Status != "" {
		query = query.Where("tasks.status = ?", filter.Status)
	}
	if filter.Tag != "" {
		query = query.Joins("JOIN task_tags ON task_tags.task_id = tasks.id").
			Joins("JOIN tags ON tags.id = task_tags.tag_id").
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
	out := make([]domain.Task, 0, len(models))
	for _, model := range models {
		out = append(out, model.toDomain())
	}
	return out, nil
}

func (r TaskRepository) Update(ctx context.Context, task domain.Task, tagNames []string) (domain.Task, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model TaskModel
		if err := tx.Preload("Tags").Where("user_id = ? AND id = ?", task.UserID, task.ID).First(&model).Error; err != nil {
			return mapError(err)
		}
		model.Title = task.Title
		model.Notes = task.Notes
		model.Priority = string(task.Priority)
		model.Status = string(task.Status)
		model.DueAt = task.DueAt
		model.CompletedAt = task.CompletedAt
		model.LinearIssueID = task.LinearIssueID
		if err := tx.Save(&model).Error; err != nil {
			return err
		}
		if tagNames != nil {
			tags, err := ensureTags(ctx, tx, task.UserID, tagNames)
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
		return domain.Task{}, err
	}
	return r.Get(ctx, task.UserID, task.ID)
}

func (r TaskRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND owner_type = ? AND owner_id = ?", userID, "task", id).Delete(&EmbeddingModel{}).Error; err != nil {
			return err
		}
		result := tx.Where("user_id = ? AND id = ?", userID, id).Delete(&TaskModel{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrNotFound
		}
		return nil
	})
}
