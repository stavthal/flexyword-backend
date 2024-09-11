package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	// Retrieve the user ID from the context set by the middleware
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return uuid.UUID{}, fmt.Errorf("unauthorized")
	}

	// Convert userId from string to uuid.UUID
	userIdStr, ok := userIdInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return uuid.UUID{}, fmt.Errorf("invalid user ID")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return uuid.UUID{}, err
	}

	return userId, nil
}
