package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/platform/metrics"
	"go.uber.org/zap"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Header("X-Request-ID", id)
		c.Set("request_id", id)
		c.Next()
	}
}

func Logging(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		status := c.Writer.Status()
		metrics.HTTPRequests.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Inc()
		log.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("duration", time.Since(start)),
			zap.String("request_id", c.GetString("request_id")),
		)
	}
}
