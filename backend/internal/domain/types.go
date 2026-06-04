package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Category string

const (
	CategoryWork     Category = "Work"
	CategoryLearning Category = "Learning"
	CategoryPersonal Category = "Personal"
	CategoryProjects Category = "Projects"
	CategoryMeetings Category = "Meetings"
	CategoryIdeas    Category = "Ideas"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

type TaskStatus string

const (
	TaskStatusInbox     TaskStatus = "inbox"
	TaskStatusToday     TaskStatus = "today"
	TaskStatusUpcoming  TaskStatus = "upcoming"
	TaskStatusOverdue   TaskStatus = "overdue"
	TaskStatusCompleted TaskStatus = "completed"
)

type ReminderStatus string

const (
	ReminderPending ReminderStatus = "pending"
	ReminderSent    ReminderStatus = "sent"
	ReminderSkipped ReminderStatus = "skipped"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type Tag struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Memory struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	Summary   string    `json:"summary"`
	Category  Category  `json:"category"`
	Source    string    `json:"source"`
	Tags      []Tag     `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Title         string     `json:"title"`
	Notes         string     `json:"notes"`
	Priority      Priority   `json:"priority"`
	Status        TaskStatus `json:"status"`
	DueAt         *time.Time `json:"due_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	LinearIssueID *string    `json:"linear_issue_id"`
	Tags          []Tag      `json:"tags"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Reminder struct {
	ID          uuid.UUID      `json:"id"`
	UserID      uuid.UUID      `json:"user_id"`
	TaskID      *uuid.UUID     `json:"task_id"`
	MemoryID    *uuid.UUID     `json:"memory_id"`
	FireAt      time.Time      `json:"fire_at"`
	Title       string         `json:"title"`
	Body        string         `json:"body"`
	Status      ReminderStatus `json:"status"`
	DeliveredAt *time.Time     `json:"delivered_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

type TelegramConfig struct {
	UserID    uuid.UUID `json:"user_id"`
	BotToken  string    `json:"bot_token,omitempty"`
	ChatID    string    `json:"chat_id"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LinearConfig struct {
	UserID    uuid.UUID `json:"user_id"`
	APIKey    string    `json:"api_key,omitempty"`
	TeamID    string    `json:"team_id"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserSettings struct {
	UserID               uuid.UUID `json:"user_id"`
	GeminiEnabled        bool      `json:"gemini_enabled"`
	TelegramEnabled      bool      `json:"telegram_enabled"`
	LinearEnabled        bool      `json:"linear_enabled"`
	NotificationsEnabled bool      `json:"notifications_enabled"`
	SiriEnabled          bool      `json:"siri_enabled"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type Enrichment struct {
	Summary  string
	Category Category
	Tags     []string
	Task     *TaskExtraction
}

type TaskExtraction struct {
	Title    string
	DueAt    *time.Time
	Priority Priority
	Tags     []string
}

type SearchResult struct {
	ID         uuid.UUID `json:"id"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Snippet    string    `json:"snippet"`
	Score      float64   `json:"score"`
	Category   string    `json:"category"`
	Tags       []string  `json:"tags"`
	OccurredAt time.Time `json:"occurred_at"`
}

type MemoryRepository interface {
	Create(ctx context.Context, memory Memory, tagNames []string, embedding []float32) (Memory, error)
	Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (Memory, error)
	List(ctx context.Context, userID uuid.UUID, filter MemoryFilter) ([]Memory, error)
	Update(ctx context.Context, memory Memory, tagNames []string) (Memory, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	Similar(ctx context.Context, userID uuid.UUID, id uuid.UUID, limit int) ([]SearchResult, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task Task, tagNames []string, embedding []float32) (Task, error)
	Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (Task, error)
	List(ctx context.Context, userID uuid.UUID, filter TaskFilter) ([]Task, error)
	Update(ctx context.Context, task Task, tagNames []string) (Task, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id uuid.UUID) (User, error)
	StoreRefreshToken(ctx context.Context, token RefreshToken) error
	FindRefreshToken(ctx context.Context, tokenHash string) (RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error
}

type ReminderRepository interface {
	Create(ctx context.Context, reminder Reminder) (Reminder, error)
	List(ctx context.Context, userID uuid.UUID, includeSent bool) ([]Reminder, error)
	MarkSent(ctx context.Context, userID uuid.UUID, id uuid.UUID, at time.Time) error
}

type SearchRepository interface {
	Search(ctx context.Context, userID uuid.UUID, query string, embedding []float32, limit int) ([]SearchResult, error)
}

type SettingsRepository interface {
	Get(ctx context.Context, userID uuid.UUID) (UserSettings, error)
	Upsert(ctx context.Context, settings UserSettings) (UserSettings, error)
	UpsertTelegram(ctx context.Context, cfg TelegramConfig) (TelegramConfig, error)
	GetTelegram(ctx context.Context, userID uuid.UUID) (TelegramConfig, error)
	UpsertLinear(ctx context.Context, cfg LinearConfig) (LinearConfig, error)
	GetLinear(ctx context.Context, userID uuid.UUID) (LinearConfig, error)
}

type LLMProvider interface {
	EnrichMemory(ctx context.Context, text string) (Enrichment, error)
	ExtractTask(ctx context.Context, text string) (TaskExtraction, error)
	Embed(ctx context.Context, text string) ([]float32, error)
}

type TelegramNotifier interface {
	Send(ctx context.Context, chatID string, message string) error
}

type LinearClient interface {
	CreateIssue(ctx context.Context, teamID string, title string, description string) (LinearIssue, error)
	UpdateIssue(ctx context.Context, issueID string, title string, description string, stateID string) (LinearIssue, error)
	ListIssues(ctx context.Context, teamID string) ([]LinearIssue, error)
}

type LinearIssue struct {
	ID          string
	Identifier  string
	Title       string
	Description string
	State       string
	URL         string
}

type MemoryFilter struct {
	Query    string
	Category string
	Tag      string
	Limit    int
	Offset   int
}

type TaskFilter struct {
	Query  string
	Status string
	Tag    string
	Limit  int
	Offset int
}
