package service

import (
	"fmt"
	"os"
	"time"

	"github.com/go-server-template/internal/models"
	"github.com/go-server-template/internal/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenSecret  = "JWT_ACCESS_SECRET"
	RefreshTokenSecret = "JWT_REFRESH_SECRET"
	salt               = "PASSWORD_SALT"

	accessTokenLifeTime  = time.Minute * 15
	refreshTokenLifeTime = time.Hour * 24 * 30
)

type AuthServiceImpl struct {
	repo repository.AuthRepository
}

// TODO: make cleanup and standardize errors & checks

func NewAuthService(repo repository.AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{repo: repo}
}

func (s *AuthServiceImpl) CreateUser(email, password string) (*models.User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateUser(email, hashedPassword)
}

func (s *AuthServiceImpl) GetUser(email, password string) (*models.User, error) {
	user, err := s.repo.GetUser(email)
	if err != nil {
		return nil, err
	}

	if err := validateUserPassword(user, password); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthServiceImpl) GetUserByToken(token string) (*models.User, error) {
	refreshToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv(RefreshTokenSecret)), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || !refreshToken.Valid {
		return nil, fmt.Errorf("empty claims")
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return nil, fmt.Errorf("token expired")
	}

	return s.repo.GetUserByToken(token)
}

func (s *AuthServiceImpl) GenerateAccessTokenString(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(accessTokenLifeTime).Unix(),
	})

	accessToken, err := token.SignedString([]byte(os.Getenv(AccessTokenSecret)))

	return accessToken, err
}

func (s *AuthServiceImpl) GenerateRefreshToken(user *models.User) (*models.RefreshToken, error) {
	expTime := time.Now().Add(refreshTokenLifeTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expTime.Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv(RefreshTokenSecret)))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.repo.CreateToken(user, tokenString, expTime)

	return refreshToken, err
}

func (s *AuthServiceImpl) RevokeToken(token string) error {
	_, err := s.GetUserByToken(token)
	if err != nil {
		return err
	}

	return s.repo.RevokeToken(token)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password+os.Getenv(salt)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func validateUserPassword(user *models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+os.Getenv(salt)))
}
