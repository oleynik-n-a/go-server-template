package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-server-template/internal/repository"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

const (
	AccessTokenSecret     = "JWT_ACCESS_SECRET"
	AccessTokenHeaderName = "Authorization"
)

func ValidationMiddleware(repo repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := c.MustGet("logger").(*zap.Logger)

		accessTokenString := c.Request.Header.Get(AccessTokenHeaderName)
		if accessTokenString == "" {
			logger.Warn("Invalid access token", zap.String("token", accessTokenString))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		accessToken, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv(AccessTokenSecret)), nil
		})

		claims, ok := accessToken.Claims.(jwt.MapClaims)
		if !ok || !accessToken.Valid {
			logger.Error("Internal server error", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if claims["exp"].(float64) < float64(time.Now().Unix()) {
			logger.Error("Token expired", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		user, err := repo.GetUser(claims["email"].(string))
		if err != nil {
			logger.Error("Internal server error", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.ID.String() != claims["id"].(string) {
			logger.Error("Id mismatch", zap.Error(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
