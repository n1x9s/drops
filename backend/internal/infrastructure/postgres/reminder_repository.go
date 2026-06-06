package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"gorm.io/gorm"
)

type ReminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return ReminderRepository{db: db}
}

func (r ReminderRepository) Create(ctx context.Context, reminder domain.Reminder) (domain.Reminder, error) {
	model := ReminderModel{
		ID:       reminder.ID,
		UserID:   reminder.UserID,
		TaskID:   reminder.TaskID,
		MemoryID: reminder.MemoryID,
		FireAt:   reminder.FireAt,
		Title:    reminder.Title,
		Body:     reminder.Body,
		Status:   string(reminder.Status),
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Reminder{}, err
	}
	return model.toDomain(), nil
}

func (r ReminderRepository) List(ctx context.Context, userID uuid.UUID, includeSent bool) ([]domain.Reminder, error) {
	var models []ReminderModel
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("fire_at ASC")
	if !includeSent {
		query = query.Where("status = ?", string(domain.ReminderPending))
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Reminder, 0, len(models))
	for _, model := range models {
		out = append(out, model.toDomain())
	}
	return out, nil
}

func (r ReminderRepository) MarkSent(ctx context.Context, userID uuid.UUID, id uuid.UUID, at time.Time) error {
	result := r.db.WithContext(ctx).Model(&ReminderModel{}).
		Where("user_id = ? AND id = ?", userID, id).
		Updates(map[string]any{"status": string(domain.ReminderSent), "delivered_at": at})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
