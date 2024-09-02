package services

import (
	"flexyword.io/backend/models"
	"gorm.io/gorm"
)

func StoreTranslation(db *gorm.DB, translation *models.Translation) error {
	return db.Create(translation).Error
}