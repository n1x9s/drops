package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App           AppConfig
	HTTP          HTTPConfig
	Database      DatabaseConfig
	JWT           JWTConfig
	Gemini        GeminiConfig
	Telegram      TelegramConfig
	Linear        LinearConfig
	Observability ObservabilityConfig
	Logging       LoggingConfig
}

type AppConfig struct {
	Env           string
	PublicBaseURL string
}

type HTTPConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	DSN string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type GeminiConfig struct {
	APIKey         string
	Model          string
	EmbeddingModel string
}

type TelegramConfig struct {
	BotToken string
	ChatID   string
}

type LinearConfig struct {
	APIKey string
	TeamID string
}

type ObservabilityConfig struct {
	OTLPEndpoint string
}

type LoggingConfig struct {
	Level string
}

func Load() (Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		var cfgNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &cfgNotFound) {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	return Config{
		App: AppConfig{
			Env:           v.GetString("app.env"),
			PublicBaseURL: v.GetString("app.public_base_url"),
		},
		HTTP: HTTPConfig{
			Host: v.GetString("http.host"),
			Port: v.GetInt("http.port"),
		},
		Database: DatabaseConfig{
			DSN: v.GetString("database.dsn"),
		},
		JWT: JWTConfig{
			AccessSecret:  v.GetString("jwt.access_secret"),
			RefreshSecret: v.GetString("jwt.refresh_secret"),
			AccessTTL:     v.GetDuration("jwt.access_ttl"),
			RefreshTTL:    v.GetDuration("jwt.refresh_ttl"),
		},
		Gemini: GeminiConfig{
			APIKey:         v.GetString("gemini.api_key"),
			Model:          v.GetString("gemini.model"),
			EmbeddingModel: v.GetString("gemini.embedding_model"),
		},
		Telegram: TelegramConfig{
			BotToken: v.GetString("telegram.bot_token"),
			ChatID:   v.GetString("telegram.chat_id"),
		},
		Linear: LinearConfig{
			APIKey: v.GetString("linear.api_key"),
			TeamID: v.GetString("linear.team_id"),
		},
		Observability: ObservabilityConfig{
			OTLPEndpoint: v.GetString("observability.otlp_endpoint"),
		},
		Logging: LoggingConfig{
			Level: v.GetString("logging.level"),
		},
	}, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.env", "development")
	v.SetDefault("app.public_base_url", "http://localhost:8080")
	v.SetDefault("http.host", "0.0.0.0")
	v.SetDefault("http.port", 8080)
	v.SetDefault("database.dsn", "postgres://secondbrain:secondbrain@localhost:5432/secondbrain?sslmode=disable")
	v.SetDefault("jwt.access_secret", "change-me-access-secret")
	v.SetDefault("jwt.refresh_secret", "change-me-refresh-secret")
	v.SetDefault("jwt.access_ttl", "15m")
	v.SetDefault("jwt.refresh_ttl", "720h")
	v.SetDefault("gemini.model", "gemini-2.5-flash")
	v.SetDefault("gemini.embedding_model", "gemini-embedding-001")
	v.SetDefault("logging.level", "debug")
}
