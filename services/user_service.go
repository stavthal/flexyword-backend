package services

import (
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

	// Assign Freemium plan if no plan is set
	if user.PricingPlanID == 0 {
		var freemiumPlan models.PricingPlan
		if err := db.Where("name = ?", "Freemium").First(&freemiumPlan).Error; err != nil {
			log.Println("Error fetching Freemium plan:", err)
			return err
		}
		user.PricingPlanID = freemiumPlan.ID
	}

	// Create the user
	if err := db.Create(&user).Error; err != nil {
		log.Println("Error creating user:", err)
		return err
	}

	return nil
}

	