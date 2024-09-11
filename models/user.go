package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            uuid.UUID     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username      string        `json:"username" gorm:"unique;not null"`       // Unique username
	Email         string        `json:"email" gorm:"unique;not null"`          // Unique email
	Password 	  string 	    `json:"password" gorm:"not null"`
	FirstName	  string		`json:"first_name" gorm:"not null"`
	LastName	  string		`json:"last_name" gorm:"not null"`
	PricingPlan   PricingPlan   `json:"pricing_plan" gorm:"foreignKey:PricingPlanID"` // Reference to the PricingPlan
	PricingPlanID uint          `json:"pricing_plan_id"`                       // Foreign key for PricingPlan
	UsedTokens    int           `json:"used_tokens" gorm:"default:0"`          // Tokens used in the current billing cycle
	Translations  []Translation `json:"translations" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`  // User's translation history
	LastReset     time.Time     `json:"last_reset"`							// Last time the token usage was reset
	BillingAddress string       `json:"billing_address" gorm:"null"`		// User's billing address
}

type UserResponse struct {
	ID            uuid.UUID     `json:"id"`
	Username      string        `json:"username"`
	Email         string        `json:"email"`
	FirstName	  string		`json:"first_name"`
	LastName	  string		`json:"last_name"`
	PricingPlan   PricingPlan   `json:"pricing_plan"`
	UsedTokens    int           `json:"used_tokens"`
	Translations  []TranslationResponse `json:"translations"`
	BillingAddress string       `json:"billing_address"`
}

