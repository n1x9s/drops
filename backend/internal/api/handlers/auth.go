package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n1x9s/second-brain/backend/internal/api/dto"
	appauth "github.com/n1x9s/second-brain/backend/internal/application/auth"
)

type AuthHandler struct {
	auth appauth.Service
}

func NewAuthHandler(auth appauth.Service) AuthHandler {
	return AuthHandler{auth: auth}
}

func (h AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !bind(c, &req) {
		return
	}
	session, err := h.auth.Register(c.Request.Context(), req.Email, req.Name, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusCreated, session)
}

func (h AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !bind(c, &req) {
		return
	}
	session, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, session)
}

func (h AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if !bind(c, &req) {
		return
	}
	session, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, session)
}

func (h AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshRequest
	if !bind(c, &req) {
		return
	}
	if err := h.auth.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		handleError(c, err)
		return
	}
	respond(c, http.StatusOK, gin.H{"logged_out": true})
}
