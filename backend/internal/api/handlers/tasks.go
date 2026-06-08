package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n1x9s/second-brain/backend/internal/api/dto"
	"github.com/n1x9s/second-brain/backend/internal/application/tasks"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"github.com/n1x9s/second-brain/backend/internal/platform/http/middleware"
)

type TaskHandler struct {
	tasks tasks.Service
}

func NewTaskHandler(tasks tasks.Service) TaskHandler {
	return TaskHandler{tasks: tasks}
}

func (h TaskHandler) Create(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	var req dto.CreateTaskRequest
	if !bind(c, &req) {
		return
	}
	task, err := h.tasks.Create(c.Request.Context(), userID, tasks.CreateInput{
		Text:     req.Text,
		Title:    req.Title,
		Notes:    req.Notes,
		DueAt:    req.DueAt,
		Priority: req.Priority,
		Tags:     req.Tags,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusCreated, task)
}

func (h TaskHandler) List(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	filter := domain.TaskFilter{
		Query:  c.Query("q"),
		Status: c.Query("status"),
		Tag:    c.Query("tag"),
		Limit:  queryInt(c, "limit", 50),
		Offset: queryInt(c, "offset", 0),
	}
	items, err := h.tasks.List(c.Request.Context(), userID, filter)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, items)
}

func (h TaskHandler) Get(c *gin.Context) {
	userID, id, ok := userAndID(c)
	if !ok {
		return
	}
	task, err := h.tasks.Get(c.Request.Context(), userID, id)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, task)
}

func (h TaskHandler) Update(c *gin.Context) {
	userID, id, ok := userAndID(c)
	if !ok {
		return
	}
	var req dto.UpdateTaskRequest
	if !bind(c, &req) {
		return
	}
	task, err := h.tasks.Update(c.Request.Context(), userID, tasks.UpdateInput{
		ID:       id,
		Title:    req.Title,
		Notes:    req.Notes,
		DueAt:    req.DueAt,
		Priority: req.Priority,
		Status:   req.Status,
		Tags:     req.Tags,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, task)
}

func (h TaskHandler) Complete(c *gin.Context) {
	userID, id, ok := userAndID(c)
	if !ok {
		return
	}
	task, err := h.tasks.Complete(c.Request.Context(), userID, id)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, task)
}

func (h TaskHandler) Delete(c *gin.Context) {
	userID, id, ok := userAndID(c)
	if !ok {
		return
	}
	if err := h.tasks.Delete(c.Request.Context(), userID, id); err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, gin.H{"deleted": true})
}
