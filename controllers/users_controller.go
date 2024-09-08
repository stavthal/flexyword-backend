package controllers

import (
	"fmt"
	"net/http"

	"flexyword.io/backend/models"
	"flexyword.io/backend/services"
	"flexyword.io/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUser(c *gin.Context, db *gorm.DB) {
	var user models.User

	// Bind the JSON body to the user model
	err := c.BindJSON(&user)

	fmt.Println(user)	

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

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
	err = services.CreateUser(db, &user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: Generate a JWT token and return it in the response to sign the user in automatically
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func LoginUser(c *gin.Context, db *gorm.DB) {
	var user models.User

	// Bind the JSON body to the user model
	c.BindJSON(&user)

	// Check if user with that email exists
	var existingUser models.User
	db.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User with that email does not exist",
		})
		return
	}

	// Check if the password is correct
	if !utils.ComparePassword(existingUser.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Convert the user's ID to a string from uuid
	parsedUserId := existingUser.ID.String()

	// Generate a JWT token and return it in the response
	token, err := utils.GenerateJWT(parsedUserId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "jwt_generation_error",
			"error": "An unexpected error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}