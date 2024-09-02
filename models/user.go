package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	BillingAddress string `json:"billing_address" gorm:"nullable"`

	// TODO: Add a field to store subscription tier and billing details
}