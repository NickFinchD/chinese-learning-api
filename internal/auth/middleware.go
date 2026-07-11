package auth

import (
	"net/http"
	"strings"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "authorization header is required",
			})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseToken(tokenString, cfg.JWT.Secret)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func GetUserID(c *gin.Context) int64 {
	value, exists := c.Get("user_id")

	if !exists {
		return 0
	}

	userID, ok := value.(int64)
	if !ok {
		return 0
	}

	return userID
}
