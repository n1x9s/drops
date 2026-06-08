package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n1x9s/second-brain/backend/internal/api/dto"
	"github.com/n1x9s/second-brain/backend/internal/domain"
)

func respond(c *gin.Context, status int, data any) {
	c.JSON(status, dto.OK(data))
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, dto.Fail("invalid_input", "request is invalid"))
	case errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, dto.Fail("unauthorized", "authentication failed"))
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, dto.Fail("not_found", "resource not found"))
	case errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, dto.Fail("conflict", "resource already exists"))
	case errors.Is(err, domain.ErrProviderEmpty):
		c.JSON(http.StatusBadRequest, dto.Fail("provider_not_configured", "integration is not configured"))
	default:
		c.JSON(http.StatusInternalServerError, dto.Fail("internal_error", err.Error()))
	}
}

func bind(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail("validation_error", err.Error()))
		return false
	}
	return true
}
