package models

import (
	"log"

	"gorm.io/gorm"
)

type PricingPlan struct {
	gorm.Model
	Name              string  `json:"name" gorm:"unique;not null"`          // Plan name (e.g., "Freemium", "Standard", "Premium")
	TranslationLimit  int     `json:"translation_limit"`                    // Max number of translations per month
	LanguagesLimit    int     `json:"languages_limit"`                      // Max number of languages per translation
	PhraseLengthLimit int     `json:"phrase_length_limit"`                  // Max length of phrase in characters
	TokenLimit        int     `json:"token_limit"`                          // Max tokens allowed per month
	AdvancedFeatures  bool    `json:"advanced_features"`                    // Whether advanced features are enabled
	PrioritySupport   bool    `json:"priority_support"`                     // Whether priority support is provided
	PricePerMonth     float64 `json:"price_per_month"`                      // Monthly price for the plan
}

// SeedPricingPlans seeds the database with default pricing plans or updates existing ones
func (p *PricingPlan) SeedPricingPlans(db *gorm.DB) error {
	// Ensure the PricingPlan table exists
	if !db.Migrator().HasTable(&PricingPlan{}) {
		log.Println("PricingPlan table does not exist, migrating the table...")
		if err := db.AutoMigrate(&PricingPlan{}); err != nil {
			log.Printf("Failed to migrate PricingPlan table: %v", err)
			return err
		}
		log.Println("PricingPlan table created successfully.")
	} else {
		// If the table already exists, auto-migrate to ensure schema is up to date
		if err := db.AutoMigrate(&PricingPlan{}); err != nil {
			log.Printf("Failed to migrate PricingPlan table: %v", err)
			return err
		}
	}

	plans := []PricingPlan{
		{
			Name:              "Freemium",
			TranslationLimit:  10,
			LanguagesLimit:    2,
			PhraseLengthLimit: 100,
			TokenLimit:        1000,
			AdvancedFeatures:  false,
			PrioritySupport:   false,
			PricePerMonth:     0.00,
		},
		{
			Name:              "Standard",
			TranslationLimit:  100,
			LanguagesLimit:    5,
			PhraseLengthLimit: 250,
			TokenLimit:        5000,
			AdvancedFeatures:  true,
			PrioritySupport:   false,
			PricePerMonth:     9.99,
		},
		{
			Name:              "Premium",
			TranslationLimit:  1000,
			LanguagesLimit:    10,
			PhraseLengthLimit: 500,
			TokenLimit:        50000,
			AdvancedFeatures:  true,
			PrioritySupport:   true,
			PricePerMonth:     29.99,
		},
	}

	for _, plan := range plans {
		var existingPlan PricingPlan
		if err := db.Where("name = ?", plan.Name).First(&existingPlan).Error; err == nil {
			// Plan exists - update the existing record with new details
			existingPlan.TranslationLimit = plan.TranslationLimit
			existingPlan.LanguagesLimit = plan.LanguagesLimit
			existingPlan.PhraseLengthLimit = plan.PhraseLengthLimit
			existingPlan.TokenLimit = plan.TokenLimit
			existingPlan.AdvancedFeatures = plan.AdvancedFeatures
			existingPlan.PrioritySupport = plan.PrioritySupport
			existingPlan.PricePerMonth = plan.PricePerMonth

			if err := db.Save(&existingPlan).Error; err != nil {
				log.Printf("Failed to update pricing plan %s: %v", plan.Name, err)
				return err
			}

			log.Printf("Updated pricing plan: %s", plan.Name)
		} else {
			// Plan does not exist - create a new record
			if err := db.Create(&plan).Error; err != nil {
				log.Printf("Failed to create pricing plan %s: %v", plan.Name, err)
				return err
			}

			log.Printf("Created pricing plan: %s", plan.Name)
		}
	}

	return nil
}
