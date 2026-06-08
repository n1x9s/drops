package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n1x9s/second-brain/backend/internal/api/dto"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Error("panic recovered", zap.Any("panic", recovered), zap.String("path", c.Request.URL.Path))
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.Fail("internal_error", "unexpected server error"))
	})
}
