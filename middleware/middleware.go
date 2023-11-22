package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuthumipepple/shopping-cart/tokens"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		clientToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, errMsg := tokens.ValidateToken(clientToken)
		if errMsg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
