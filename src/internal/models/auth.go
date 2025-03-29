package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id" db:"id" binding:"required,uuid"`
	Email    string    `json:"email" db:"email" binding:"required,email"`
	Password string    `json:"password" db:"password" binding:"required,min=8,max=16"`
}

type RefreshToken struct {
	ID          uuid.UUID `json:"id" binding:"required,uuid"`
	UserId      uuid.UUID `json:"user_id" binding:"required,uuid"`
	TokenString string    `json:"token" binding:"required"`
	ExpiresAt   time.Time `json:"expires_at" binding:"required"`
	Revoked     bool      `json:"revoked" binding:"required"`
}

func NewUser(id uuid.UUID, email, password string) *User {
	return &User{ID: id, Email: email, Password: password}
}

func NewRefreshToken(id uuid.UUID, userId uuid.UUID, token string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{ID: id, UserId: userId, TokenString: token, ExpiresAt: expiresAt, Revoked: false}
}
