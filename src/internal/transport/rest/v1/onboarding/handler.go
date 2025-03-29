package onboarding

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-server-template/internal/models"
	"github.com/go-server-template/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) Test(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.Logger)
	user := c.MustGet("user").(*models.User)

	logger.Info("User is here", zap.String("email", user.Email))
	c.JSON(http.StatusCreated, gin.H{"message": "Test worked"})
}
