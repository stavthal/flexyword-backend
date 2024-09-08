package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Translation struct {
	gorm.Model
	ID 				  uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Phrase            string    `json:"phrase" gorm:"not null"`
	InputLanguage     string    `json:"input_language" gorm:"not null"`
	OutputLanguages   string    `json:"output_languages" gorm:"type:jsonb;not null"`  // JSON string
	TranslationResult string    `json:"translation_result" gorm:"type:jsonb;not null"` // JSON string
	UserID            uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User              User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// Custom response struct to exclude User and other unnecessary fields
type TranslationResponse struct {
	ID 			      string `json:"id"`
	Phrase            string `json:"phrase"`
	InputLanguage     string `json:"input_language"`
	OutputLanguages   string `json:"output_languages"`
	TranslationResult string `json:"translation_result"`
	CreatedAt         time.Time `json:"created_at"`
}