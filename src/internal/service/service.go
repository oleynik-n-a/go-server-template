package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-server-template/internal/middleware"
	"github.com/go-server-template/internal/models"
	"github.com/go-server-template/internal/repository"
	"go.uber.org/zap"
)

type AuthService interface {
	CreateUser(string, string) (*models.User, error)
	GetUser(string, string) (*models.User, error)
	GetUserByToken(string) (*models.User, error)
	GenerateAccessTokenString(*models.User) (string, error)
	GenerateRefreshToken(*models.User) (*models.RefreshToken, error)
	RevokeToken(string) error
}

type Service struct {
	Router *gin.Engine
	Repo   *repository.Repository

	AuthService
}

func NewService(repo *repository.Repository, logger *zap.Logger) *Service {
	s := &Service{
		Router:      gin.New(),
		Repo:        repo,
		AuthService: NewAuthService(repo.AuthRepository),
	}

	s.attachMiddleware(middleware.LoggerMiddleware(logger))
	s.attachMiddleware(middleware.RecoveryMiddleware())

	return s
}

func (s *Service) attachMiddleware(middleware gin.HandlerFunc) {
	s.Router.Use(middleware)
}
