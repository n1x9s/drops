package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/api/dto"
	"github.com/n1x9s/second-brain/backend/internal/infrastructure/security"
)

const userIDKey = "user_id"

func Auth(tokens security.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(strings.ToLower(header), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail("unauthorized", "missing bearer token"))
			return
		}
		claims, err := tokens.ParseAccess(strings.TrimSpace(header[7:]))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail("unauthorized", "invalid access token"))
			return
		}
		c.Set(userIDKey, claims.UserID)
		c.Next()
	}
}

func UserID(c *gin.Context) (uuid.UUID, bool) {
	value, ok := c.Get(userIDKey)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := value.(uuid.UUID)
	return id, ok
}
