package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"gorm.io/gorm"
)

type UserModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"size:320;uniqueIndex;not null"`
	Name         string    `gorm:"size:160;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshTokenModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;index;not null"`
	User      UserModel  `gorm:"foreignKey:UserID"`
	TokenHash string     `gorm:"size:128;uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"index;not null"`
	RevokedAt *time.Time `gorm:"index"`
	CreatedAt time.Time
}

type TagModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_tags_user_name;not null"`
	Name      string    `gorm:"size:64;uniqueIndex:idx_tags_user_name;not null"`
	CreatedAt time.Time
}

type MemoryModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;index;not null"`
	User      UserModel  `gorm:"foreignKey:UserID"`
	Content   string     `gorm:"type:text;not null"`
	Summary   string     `gorm:"type:text;not null"`
	Category  string     `gorm:"size:32;index;not null"`
	Source    string     `gorm:"size:32;not null"`
	Tags      []TagModel `gorm:"many2many:memory_tags;joinForeignKey:MemoryID;joinReferences:TagID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TaskModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID  `gorm:"type:uuid;index;not null"`
	User          UserModel  `gorm:"foreignKey:UserID"`
	Title         string     `gorm:"size:500;not null"`
	Notes         string     `gorm:"type:text;not null"`
	Priority      string     `gorm:"size:24;index;not null"`
	Status        string     `gorm:"size:32;index;not null"`
	DueAt         *time.Time `gorm:"index"`
	CompletedAt   *time.Time
	LinearIssueID *string    `gorm:"size:64"`
	Tags          []TagModel `gorm:"many2many:task_tags;joinForeignKey:TaskID;joinReferences:TagID;constraint:OnDelete:CASCADE"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type ReminderModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID  `gorm:"type:uuid;index;not null"`
	TaskID      *uuid.UUID `gorm:"type:uuid;index"`
	MemoryID    *uuid.UUID `gorm:"type:uuid;index"`
	FireAt      time.Time  `gorm:"index;not null"`
	Title       string     `gorm:"size:300;not null"`
	Body        string     `gorm:"type:text;not null"`
	Status      string     `gorm:"size:32;index;not null"`
	DeliveredAt *time.Time
	CreatedAt   time.Time
}

type TelegramConfigModel struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	BotToken  string    `gorm:"type:text;not null"`
	ChatID    string    `gorm:"size:128;not null"`
	Enabled   bool      `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LinearConfigModel struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	APIKey    string    `gorm:"type:text;not null"`
	TeamID    string    `gorm:"size:128;not null"`
	Enabled   bool      `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserSettingsModel struct {
	UserID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	GeminiEnabled        bool      `gorm:"not null;default:true"`
	TelegramEnabled      bool      `gorm:"not null;default:false"`
	LinearEnabled        bool      `gorm:"not null;default:false"`
	NotificationsEnabled bool      `gorm:"not null;default:true"`
	SiriEnabled          bool      `gorm:"not null;default:true"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type EmbeddingModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	OwnerType string    `gorm:"size:32;uniqueIndex:idx_embeddings_owner;not null"`
	OwnerID   uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_embeddings_owner;index;not null"`
	Embedding Vector    `gorm:"type:vector(768);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string {
	return "users"
}

func (RefreshTokenModel) TableName() string {
	return "refresh_tokens"
}

func (TagModel) TableName() string {
	return "tags"
}

func (MemoryModel) TableName() string {
	return "memories"
}

func (TaskModel) TableName() string {
	return "tasks"
}

func (ReminderModel) TableName() string {
	return "reminders"
}

func (TelegramConfigModel) TableName() string {
	return "telegram_configs"
}

func (LinearConfigModel) TableName() string {
	return "linear_configs"
}

func (UserSettingsModel) TableName() string {
	return "user_settings"
}

func (EmbeddingModel) TableName() string {
	return "embeddings"
}

func (m UserModel) toDomain() domain.User {
	return domain.User{ID: m.ID, Email: m.Email, Name: m.Name, PasswordHash: m.PasswordHash, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func (m RefreshTokenModel) toDomain() domain.RefreshToken {
	return domain.RefreshToken{ID: m.ID, UserID: m.UserID, TokenHash: m.TokenHash, ExpiresAt: m.ExpiresAt, RevokedAt: m.RevokedAt, CreatedAt: m.CreatedAt}
}

func (m MemoryModel) toDomain() domain.Memory {
	return domain.Memory{ID: m.ID, UserID: m.UserID, Content: m.Content, Summary: m.Summary, Category: domain.Category(m.Category), Source: m.Source, Tags: tagsToDomain(m.Tags), CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func (m TaskModel) toDomain() domain.Task {
	return domain.Task{ID: m.ID, UserID: m.UserID, Title: m.Title, Notes: m.Notes, Priority: domain.Priority(m.Priority), Status: domain.TaskStatus(m.Status), DueAt: m.DueAt, CompletedAt: m.CompletedAt, LinearIssueID: m.LinearIssueID, Tags: tagsToDomain(m.Tags), CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func (m ReminderModel) toDomain() domain.Reminder {
	return domain.Reminder{ID: m.ID, UserID: m.UserID, TaskID: m.TaskID, MemoryID: m.MemoryID, FireAt: m.FireAt, Title: m.Title, Body: m.Body, Status: domain.ReminderStatus(m.Status), DeliveredAt: m.DeliveredAt, CreatedAt: m.CreatedAt}
}

func (m UserSettingsModel) toDomain() domain.UserSettings {
	return domain.UserSettings{UserID: m.UserID, GeminiEnabled: m.GeminiEnabled, TelegramEnabled: m.TelegramEnabled, LinearEnabled: m.LinearEnabled, NotificationsEnabled: m.NotificationsEnabled, SiriEnabled: m.SiriEnabled, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func tagsToDomain(tags []TagModel) []domain.Tag {
	out := make([]domain.Tag, 0, len(tags))
	for _, tag := range tags {
		out = append(out, domain.Tag{ID: tag.ID, UserID: tag.UserID, Name: tag.Name, CreatedAt: tag.CreatedAt})
	}
	return out
}
