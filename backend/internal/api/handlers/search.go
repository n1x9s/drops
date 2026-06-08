package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n1x9s/second-brain/backend/internal/application/search"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"github.com/n1x9s/second-brain/backend/internal/platform/http/middleware"
)

type SearchHandler struct {
	search search.Service
}

func NewSearchHandler(search search.Service) SearchHandler {
	return SearchHandler{search: search}
}

func (h SearchHandler) Search(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		handleError(c, domain.ErrUnauthorized)
		return
	}
	results, err := h.search.Search(c.Request.Context(), userID, c.Query("q"), queryInt(c, "limit", 25))
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, results)
}
