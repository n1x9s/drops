package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/api/dto"
	"github.com/n1x9s/second-brain/backend/internal/application/reminders"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"github.com/n1x9s/second-brain/backend/internal/platform/http/middleware"
)

type ReminderHandler struct {
	reminders reminders.Service
}

func NewReminderHandler(reminders reminders.Service) ReminderHandler {
	return ReminderHandler{reminders: reminders}
}

func (h ReminderHandler) Create(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	var req dto.CreateReminderRequest
	if !bind(c, &req) {
		return
	}
	taskID, err := parseOptionalUUID(req.TaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail("invalid_task_id", "task_id must be a UUID"))
		return
	}
	memoryID, err := parseOptionalUUID(req.MemoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail("invalid_memory_id", "memory_id must be a UUID"))
		return
	}
	reminder, err := h.reminders.Create(c.Request.Context(), userID, reminders.CreateInput{
		TaskID:   taskID,
		MemoryID: memoryID,
		FireAt:   req.FireAt,
		Title:    req.Title,
		Body:     req.Body,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusCreated, reminder)
}

func (h ReminderHandler) List(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	items, err := h.reminders.List(c.Request.Context(), userID, c.Query("include_sent") == "true")
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, items)
}

func parseOptionalUUID(raw *string) (*uuid.UUID, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(*raw)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
