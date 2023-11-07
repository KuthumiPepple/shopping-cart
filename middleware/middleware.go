package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuthumipepple/shopping-cart/tokens"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no authorization header provided"})
			c.Abort()
			return
		}

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
