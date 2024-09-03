package services

import (
	"fmt"
	"log"

	"flexyword.io/backend/models"
	"flexyword.io/backend/utils"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	// Hash the user's password
	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		return err
	}

	// Set the user's password to the hashed password
	user.Password = hashedPassword

	fmt.Println("Users password: " , user.Password)

	// Create the user
	err = db.Create(&user).Error

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
	