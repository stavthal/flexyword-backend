package models

import (
	"gorm.io/gorm"
)

type Translation struct {
	gorm.Model
	Language    	string 				`json:"language" gorm:"size:50;not null"`  // Source language
	Phrase	     	string 				`json:"phrase" gorm:"size:100;not null"`  // Phrase with max 100 characters
	Translations 	map[string]string 	`json:"translations" gorm:"type:json;serializer:json" `  // Map of translations

}
