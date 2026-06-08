package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=1,max=160"`
	Password string `json:"password" binding:"required,min=8,max=200"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type CreateMemoryRequest struct {
	Content string `json:"content" binding:"required,min=1"`
	Source  string `json:"source" binding:"omitempty,oneof=manual siri widget telegram import"`
}

type UpdateMemoryRequest struct {
	Content  string   `json:"content"`
	Summary  string   `json:"summary"`
	Category string   `json:"category" binding:"omitempty,oneof=Work Learning Personal Projects Meetings Ideas"`
	Tags     []string `json:"tags"`
}

type CreateTaskRequest struct {
	Text     string     `json:"text"`
	Title    string     `json:"title"`
	Notes    string     `json:"notes"`
	DueAt    *time.Time `json:"due_at"`
	Priority string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Tags     []string   `json:"tags"`
}

type UpdateTaskRequest struct {
	Title    string     `json:"title"`
	Notes    string     `json:"notes"`
	DueAt    *time.Time `json:"due_at"`
	Priority string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Status   string     `json:"status" binding:"omitempty,oneof=inbox today upcoming overdue completed"`
	Tags     []string   `json:"tags"`
}

type CreateReminderRequest struct {
	TaskID   *string   `json:"task_id"`
	MemoryID *string   `json:"memory_id"`
	FireAt   time.Time `json:"fire_at" binding:"required"`
	Title    string    `json:"title" binding:"required"`
	Body     string    `json:"body"`
}

type UpdateSettingsRequest struct {
	GeminiEnabled        bool `json:"gemini_enabled"`
	TelegramEnabled      bool `json:"telegram_enabled"`
	LinearEnabled        bool `json:"linear_enabled"`
	NotificationsEnabled bool `json:"notifications_enabled"`
	SiriEnabled          bool `json:"siri_enabled"`
}

type TelegramConfigRequest struct {
	BotToken string `json:"bot_token" binding:"required"`
	ChatID   string `json:"chat_id" binding:"required"`
	Enabled  bool   `json:"enabled"`
}

type LinearConfigRequest struct {
	APIKey  string `json:"api_key" binding:"required"`
	TeamID  string `json:"team_id" binding:"required"`
	Enabled bool   `json:"enabled"`
}

type CreateLinearIssueRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateLinearIssueRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StateID     string `json:"state_id"`
}
