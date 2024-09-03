package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string        `json:"username" gorm:"unique;not null"`       // Unique username
	Email         string        `json:"email" gorm:"unique;not null"`          // Unique email
	Password 	  string 	    `json:"password" gorm:"not null"`
	PricingPlan   PricingPlan   `json:"pricing_plan" gorm:"foreignKey:PricingPlanID"` // Reference to the PricingPlan
	PricingPlanID uint          `json:"pricing_plan_id"`                       // Foreign key for PricingPlan
	UsedTokens    int           `json:"used_tokens" gorm:"default:0"`          // Tokens used in the current billing cycle
	Translations  []Translation `json:"translations" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`  // User's translation history
	LastReset     time.Time     `json:"last_reset"`    
}


	// TODO: Add billing address
