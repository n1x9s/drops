package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n1x9s/second-brain/backend/internal/api/dto"
	"github.com/n1x9s/second-brain/backend/internal/application/integrations"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"github.com/n1x9s/second-brain/backend/internal/platform/http/middleware"
)

type IntegrationHandler struct {
	integrations integrations.Service
}

func NewIntegrationHandler(integrations integrations.Service) IntegrationHandler {
	return IntegrationHandler{integrations: integrations}
}

func (h IntegrationHandler) TelegramTest(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	if err := h.integrations.SendTelegramTest(c.Request.Context(), userID); err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, gin.H{"sent": true})
}

func (h IntegrationHandler) CreateLinearIssue(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	var req dto.CreateLinearIssueRequest
	if !bind(c, &req) {
		return
	}
	issue, err := h.integrations.CreateLinearIssue(c.Request.Context(), userID, req.Title, req.Description)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusCreated, issue)
}

func (h IntegrationHandler) UpdateLinearIssue(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	var req dto.UpdateLinearIssueRequest
	if !bind(c, &req) {
		return
	}
	issue, err := h.integrations.UpdateLinearIssue(c.Request.Context(), userID, c.Param("id"), req.Title, req.Description, req.StateID)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, issue)
}

func (h IntegrationHandler) ListLinearIssues(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	issues, err := h.integrations.ListLinearIssues(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, issues)
}
