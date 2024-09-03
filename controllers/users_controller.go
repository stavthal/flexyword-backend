package controllers

import (
	"net/http"

	"flexyword.io/backend/models"
	"flexyword.io/backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUser(c *gin.Context, db *gorm.DB) {
	var user models.User

	// Bind the JSON body to the user model
	c.BindJSON(&user)

	// Check if user with that email already exists and handle it
	var existingUser models.User
	db.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User with that email already exists",
		})
		return
	}

	// Check if user with that username already exists and handle it
	db.Where("username = ?", user.Username).First(&existingUser)

	if existingUser.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User with that username already exists",
		})
		return
	}

	// Create the user
	services.CreateUser(db, &user)

	// TODO: Generate a JWT token and return it in the response to sign the user in automatically
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}