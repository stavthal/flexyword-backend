package models

import "time"

type Translation struct {
	ID          	uint   			`json:"id" gorm:"primaryKey"`
	Language    	string 			`json:"language" gorm:"size:50;not null"`  // Source language
	Phrase	     	string 			`json:"phrase" gorm:"size:100;not null"`  // Phrase with max 100 characters
	Translations 	map[string]string `json:"translations" gorm:"type:json"`  // Map of translations
	CreatedAt 		time.Time 		`json:"created_at" gorm:"autoCreateTime"`
}
