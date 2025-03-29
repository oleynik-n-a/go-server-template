package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-server-template/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	service.AuthService
}

func NewHandler(s service.AuthService) *Handler {
	return &Handler{AuthService: s}
}

func (h *Handler) Signup(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)

	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.AuthService.CreateUser(req.Email, req.Password)
	if err != nil {
		logger.Error("Internal server error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessTokenString, err := h.AuthService.GenerateAccessTokenString(user)
	if err != nil {
		logger.Error("Internal server error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := h.AuthService.GenerateRefreshToken(user)
	if err != nil {
		logger.Error("Internal server error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("User registered", zap.String("uuid", user.ID.String()), zap.String("email", user.Email))
	c.JSON(http.StatusCreated, gin.H{
		"message":       "user registered",
		"access_token":  accessTokenString,
		"refresh_token": refreshToken.TokenString,
	})
}

func (h *Handler) Signin(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)

	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.AuthService.GetUser(req.Email, req.Password)
	if err != nil {
		logger.Error("Internal server error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessTokenString, err := h.AuthService.GenerateAccessTokenString(user)
	if err != nil {
		logger.Error("Internal server error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := h.AuthService.GenerateRefreshToken(user)
	if err != nil {
		logger.Error("Internal server error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("User logged in", zap.String("user_id", user.ID.String()), zap.String("token_id", refreshToken.ID.String()))

	c.JSON(http.StatusOK, gin.H{
		"message":       "user logged in",
		"access_token":  accessTokenString,
		"refresh_token": refreshToken.TokenString,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.AuthService.GetUserByToken(req.RefreshToken)
	if err != nil {
		logger.Error("Invalid token error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.AuthService.RevokeToken(req.RefreshToken); err != nil {
		logger.Error("Internal server error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessTokenString, err := h.AuthService.GenerateAccessTokenString(user)
	if err != nil {
		logger.Error("Internal server error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := h.AuthService.GenerateRefreshToken(user)
	if err != nil {
		logger.Error("Internal server error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "tokens refreshed",
		"access_token":  accessTokenString,
		"refresh_token": refreshToken.TokenString,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)

	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.AuthService.RevokeToken(req.RefreshToken); err != nil {
		logger.Error("Internal server error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user logged out"})
}
