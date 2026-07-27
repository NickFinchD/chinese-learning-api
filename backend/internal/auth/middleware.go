package auth

import (
	"net/http"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/NickFinchD/chinese-learning-api/internal/response"
	"github.com/NickFinchD/chinese-learning-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {

		tokenString, err := c.Cookie("access_token")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "access token is required",
			})
			return
		}

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

// RequireAdmin gates admin-only routes. It must run after JWTMiddleware
// (which sets user_id). It re-checks is_admin against the database on every
// request, rather than trusting a JWT claim, so revoking someone's admin
// access takes effect immediately instead of waiting up to the token's
// lifetime to expire.
func RequireAdmin(service *Service) gin.HandlerFunc {

	return func(c *gin.Context) {

		userID := GetUserID(c)

		isAdmin, err := service.IsAdmin(c.Request.Context(), userID)

		if err != nil {
			response.Internal(c)
			c.Abort()
			return
		}

		if !isAdmin {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}

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
