package middlewares

import (
	"net/http"
	"strings"

	"flexyword.io/backend/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware function that checks the JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// The token is usually sent as "Bearer <token>", so split the header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Verify the token using the VerifyJWT function
		token := tokenParts[1]
		claims, err := utils.VerifyJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store the userId in the context
		userId := (*claims)["userId"]
		c.Set("userId", userId)

		// Continue to the next handler
		c.Next()
	}
}
