package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n1x9s/second-brain/backend/internal/api/dto"
	"github.com/n1x9s/second-brain/backend/internal/application/settings"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"github.com/n1x9s/second-brain/backend/internal/platform/http/middleware"
)

type SettingsHandler struct {
	settings settings.Service
}

func NewSettingsHandler(settings settings.Service) SettingsHandler {
	return SettingsHandler{settings: settings}
}

func (h SettingsHandler) Get(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	settings, err := h.settings.Get(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, settings)
}

func (h SettingsHandler) Update(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	var req dto.UpdateSettingsRequest
	if !bind(c, &req) {
		return
	}
	settings, err := h.settings.Upsert(c.Request.Context(), domain.UserSettings{
		UserID:               userID,
		GeminiEnabled:        req.GeminiEnabled,
		TelegramEnabled:      req.TelegramEnabled,
		LinearEnabled:        req.LinearEnabled,
		NotificationsEnabled: req.NotificationsEnabled,
		SiriEnabled:          req.SiriEnabled,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, settings)
}

func (h SettingsHandler) UpsertTelegram(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	var req dto.TelegramConfigRequest
	if !bind(c, &req) {
		return
	}
	cfg, err := h.settings.UpsertTelegram(c.Request.Context(), domain.TelegramConfig{UserID: userID, BotToken: req.BotToken, ChatID: req.ChatID, Enabled: req.Enabled})
	if err != nil {
		handleError(c, err)
		return
	}
	cfg.BotToken = ""
	respond(c, http.StatusOK, cfg)
}

func (h SettingsHandler) UpsertLinear(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	var req dto.LinearConfigRequest
	if !bind(c, &req) {
		return
	}
	cfg, err := h.settings.UpsertLinear(c.Request.Context(), domain.LinearConfig{UserID: userID, APIKey: req.APIKey, TeamID: req.TeamID, Enabled: req.Enabled})
	if err != nil {
		handleError(c, err)
		return
	}
	cfg.APIKey = ""
	respond(c, http.StatusOK, cfg)
}
