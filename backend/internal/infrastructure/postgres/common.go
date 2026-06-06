package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrNotFound
	}
	return err
}

func ensureTags(ctx context.Context, tx *gorm.DB, userID uuid.UUID, names []string) ([]TagModel, error) {
	normalized := make([]string, 0, len(names))
	seen := map[string]struct{}{}
	for _, name := range names {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		normalized = append(normalized, name)
	}

	tags := make([]TagModel, 0, len(normalized))
	for _, name := range normalized {
		tag := TagModel{ID: uuid.New(), UserID: userID, Name: name}
		if err := tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "name"}},
			DoNothing: true,
		}).Create(&tag).Error; err != nil {
			return nil, err
		}
		if err := tx.WithContext(ctx).Where("user_id = ? AND name = ?", userID, name).First(&tag).Error; err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func tagNames(tags []domain.Tag) []string {
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		out = append(out, tag.Name)
	}
	return out
}
