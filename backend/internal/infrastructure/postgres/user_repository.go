package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	model := UserModel{ID: user.ID, Email: user.Email, Name: user.Name, PasswordHash: user.PasswordHash}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return domain.User{}, domain.ErrConflict
		}
		return domain.User{}, err
	}
	return model.toDomain(), nil
}

func (r UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		return domain.User{}, mapError(err)
	}
	return model.toDomain(), nil
}

func (r UserRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return domain.User{}, mapError(err)
	}
	return model.toDomain(), nil
}

func (r UserRepository) StoreRefreshToken(ctx context.Context, token domain.RefreshToken) error {
	model := RefreshTokenModel{ID: token.ID, UserID: token.UserID, TokenHash: token.TokenHash, ExpiresAt: token.ExpiresAt}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r UserRepository) FindRefreshToken(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	var model RefreshTokenModel
	if err := r.db.WithContext(ctx).
		Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now().UTC()).
		First(&model).Error; err != nil {
		return domain.RefreshToken{}, mapError(err)
	}
	return model.toDomain(), nil
}

func (r UserRepository) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).Model(&RefreshTokenModel{}).Where("id = ?", tokenID).Update("revoked_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
