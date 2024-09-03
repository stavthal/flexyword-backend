package models

import (
	"gorm.io/gorm"
)

type Translation struct {
	gorm.Model
	Phrase            string `json:"phrase" gorm:"not null"`
	InputLanguage     string `json:"input_language" gorm:"not null"`
	OutputLanguages   string `json:"output_languages" gorm:"type:jsonb;not null"`  // JSON string
	TranslationResult string `json:"translation_result" gorm:"type:jsonb;not null"` // JSON string
	UserID            uint   `json:"user_id" gorm:"not null"`
}
