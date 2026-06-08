package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/api/dto"
	"github.com/n1x9s/second-brain/backend/internal/application/memories"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"github.com/n1x9s/second-brain/backend/internal/platform/http/middleware"
)

type MemoryHandler struct {
	memories memories.Service
}

func NewMemoryHandler(memories memories.Service) MemoryHandler {
	return MemoryHandler{memories: memories}
}

func (h MemoryHandler) Create(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	var req dto.CreateMemoryRequest
	if !bind(c, &req) {
		return
	}
	memory, err := h.memories.Create(c.Request.Context(), userID, memories.CreateInput{Content: req.Content, Source: req.Source})
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusCreated, memory)
}

func (h MemoryHandler) List(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	filter := domain.MemoryFilter{
		Query:    c.Query("q"),
		Category: c.Query("category"),
		Tag:      c.Query("tag"),
		Limit:    queryInt(c, "limit", 50),
		Offset:   queryInt(c, "offset", 0),
	}
	items, err := h.memories.List(c.Request.Context(), userID, filter)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, items)
}

func (h MemoryHandler) Get(c *gin.Context) {
	userID, id, ok := userAndID(c)
	if !ok {
		return
	}
	memory, err := h.memories.Get(c.Request.Context(), userID, id)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, memory)
}

func (h MemoryHandler) Update(c *gin.Context) {
	userID, id, ok := userAndID(c)
	if !ok {
		return
	}
	var req dto.UpdateMemoryRequest
	if !bind(c, &req) {
		return
	}
	memory, err := h.memories.Update(c.Request.Context(), userID, memories.UpdateInput{
		ID:       id,
		Content:  req.Content,
		Summary:  req.Summary,
		Category: req.Category,
		Tags:     req.Tags,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, memory)
}

func (h MemoryHandler) Delete(c *gin.Context) {
	userID, id, ok := userAndID(c)
	if !ok {
		return
	}
	if err := h.memories.Delete(c.Request.Context(), userID, id); err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, gin.H{"deleted": true})
}

func (h MemoryHandler) Similar(c *gin.Context) {
	userID, id, ok := userAndID(c)
	if !ok {
		return
	}
	items, err := h.memories.Similar(c.Request.Context(), userID, id, queryInt(c, "limit", 10))
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, items)
}

func userAndID(c *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail("invalid_id", "id must be a UUID"))
		return uuid.Nil, uuid.Nil, false
	}
	return userID, id, true
}

func queryInt(c *gin.Context, key string, fallback int) int {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
