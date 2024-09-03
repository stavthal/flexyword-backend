package models

import "gorm.io/gorm"

type Translation struct {
	gorm.Model
	Phrase           string `json:"phrase" gorm:"not null"`
	InputLanguage    string `json:"input_language" gorm:"not null"`
	OutputLanguages  string `json:"output_languages" gorm:"type:text;not null"` // JSON encoded array of languages
	TranslationResult string `json:"translation_result" gorm:"type:text;not null"` // JSON encoded translation results
	UserID           uint   `json:"user_id" gorm:"not null"` // Foreign key to the User model
}
