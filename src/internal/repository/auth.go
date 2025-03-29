package repository

import (
	"fmt"
	"time"

	"github.com/go-server-template/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthRepositoryImpl struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) CreateUser(email, password string) (*models.User, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (email, password) values ($1, $2) RETURNING id", usersTable)
	err := r.db.QueryRow(query, email, password).Scan(&id)
	return models.NewUser(id, email, password), err
}

func (r *AuthRepositoryImpl) GetUser(email string) (*models.User, error) {
	var id uuid.UUID
	var password string
	query := fmt.Sprintf("SELECT id, password FROM %s WHERE email=$1", usersTable)
	err := r.db.QueryRow(query, email).Scan(&id, &password)
	return models.NewUser(id, email, password), err
}

// TODO: fix code style
func (r *AuthRepositoryImpl) GetUserByToken(token string) (*models.User, error) {
	query := fmt.Sprintf(
		"SELECT u.* FROM %s u JOIN %s rt ON u.id = rt.user_id WHERE rt.token = $1 AND rt.revoked = FALSE;",
		usersTable, tokensTable,
	)

	var id uuid.UUID
	var email, password string

	err := r.db.QueryRow(query, token).Scan(&id, &email, &password)
	user := models.NewUser(id, email, password)
	return user, err
}

func (r *AuthRepositoryImpl) CreateToken(user *models.User, tokenString string, expiresAt time.Time) (*models.RefreshToken, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (user_id, token, expires_at) values ($1, $2, $3) RETURNING id", tokensTable)
	err := r.db.QueryRow(query, user.ID, tokenString, expiresAt).Scan(&id)
	return models.NewRefreshToken(id, user.ID, tokenString, expiresAt), err
}

func (r *AuthRepositoryImpl) RevokeToken(token string) error {
	query := fmt.Sprintf("UPDATE %s SET revoked=$1 WHERE token=$2", tokensTable)
	_, err := r.db.Exec(query, true, token)
	return err
}
