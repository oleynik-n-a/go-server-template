package repository

import (
	"time"

	"github.com/go-server-template/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	CreateUser(string, string) (*models.User, error)
	GetUser(string) (*models.User, error)
	GetUserByToken(token string) (*models.User, error)
	CreateToken(*models.User, string, time.Time) (*models.RefreshToken, error)
	RevokeToken(string) error
}

type Repository struct {
	AuthRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{AuthRepository: NewAuthRepository(db)}
}
