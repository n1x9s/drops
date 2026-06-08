package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) HealthHandler {
	return HealthHandler{db: db}
}

func (h HealthHandler) Live(c *gin.Context) {
	respond(c, http.StatusOK, gin.H{"status": "live"})
}

func (h HealthHandler) Ready(c *gin.Context) {
	if err := h.db.PingContext(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "error": gin.H{"code": "database_unavailable", "message": err.Error()}})
		return
	}
	respond(c, http.StatusOK, gin.H{"status": "ready"})
}
