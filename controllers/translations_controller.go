package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"flexyword.io/backend/models"
	"flexyword.io/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// EstimateTokens estimates the number of tokens used based on the length of the text.
// This is a rough approximation.
func EstimateTokens(text string) int {
	return len(text) / 4 // Roughly 4 characters per token
}

// Request model
type TranslateRequest struct {
	InputLanguage string   `json:"input_language" binding:"required"`
	Languages     []string `json:"languages" binding:"required"`
	Phrase        string   `json:"phrase" binding:"required"`
}

// TranslatePhrase handles the translation of a phrase into multiple languages
func TranslatePhrase(c *gin.Context, db *gorm.DB) {
	var request TranslateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user ID from the context set by the middleware
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert userId from string to uuid.UUID
	userIdStr, ok := userIdInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Fetch user from the database
	var user models.User
	if err := db.Preload("PricingPlan").First(&user, "id = ?", userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Enforce pricing plan constraints
	plan := user.PricingPlan

	// Maximum number of languages
	if len(request.Languages) > plan.LanguagesLimit {
		c.JSON(http.StatusForbidden, gin.H{"error": "Exceeded the maximum number of languages allowed for your plan"})
		return
	}

	// Maximum phrase length
	if len(request.Phrase) > plan.PhraseLengthLimit {
		c.JSON(http.StatusForbidden, gin.H{"error": "Exceeded the maximum phrase length allowed for your plan"})
		return
	}

	// Check if the user has exceeded the monthly translation limit
	currentMonth := time.Now().Format("2006-01")
	var translationsCount int64
	if err := db.Model(&models.Translation{}).
		Where("user_id = ? AND to_char(created_at, 'YYYY-MM') = ?", user.ID, currentMonth).
		Count(&translationsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count translations"})
		return
	}

	if translationsCount >= int64(plan.TranslationLimit) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You have reached the maximum number of translations for this month"})
		return
	}

	// Initialize OpenAI client
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	translations := make(map[string]string)
	totalTokensUsed := 0

	// Loop through each language and request translation from OpenAI
	for _, lang := range request.Languages {
		prompt := "Translate the following phrase from " + request.InputLanguage + " into " + lang + ": " + request.Phrase + ". Return only the translated phrase."

		resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: "gpt-4-turbo",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
			MaxTokens:   500,  // Adjust max tokens based on expected translation length
			Temperature: 0.2,  // Low temperature for deterministic results
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get translation from OpenAI: " + err.Error()})
			return
		}

		if len(resp.Choices) > 0 {
			translation := strings.TrimSpace(resp.Choices[0].Message.Content)
			translations[lang] = translation
			totalTokensUsed += EstimateTokens(translation)
		} else {
			translations[lang] = ""
		}
	}

	// Check if the total tokens used exceed the plan's token limit
	if user.UsedTokens+totalTokensUsed > plan.TokenLimit {
		c.JSON(http.StatusForbidden, gin.H{"error": "You have reached the token limit for this month"})
		return
	}

	// Update the user's token usage
	user.UsedTokens += totalTokensUsed
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user token usage"})
		return
	}

	// Convert translations map to JSON string
	translationsJSON, err := json.Marshal(translations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode translations: " + err.Error()})
		return
	}

	// Convert output languages array to JSON string
	outputLanguagesJSON, err := json.Marshal(request.Languages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode output languages: " + err.Error()})
		return
	}

	// Store the translation in the database
	translation := models.Translation{
		Phrase:            request.Phrase,
		InputLanguage:     request.InputLanguage,
		OutputLanguages:   string(outputLanguagesJSON),
		TranslationResult: string(translationsJSON),
		UserID:            user.ID,
	}

	if err := services.StoreTranslation(db, &translation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store translation: " + err.Error()})
		return
	}

	// Send the response as JSON
	c.JSON(http.StatusOK, gin.H{
		"phrase":       request.Phrase,
		"translations": translations,
	})
}


// GetTranslations retrieves the translation history for the authenticated user
func GetTranslations(c *gin.Context, db *gorm.DB) {
	// Retrieve the user ID from the context set by the middleware
	userIdInterface, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert userId from string to uuid.UUID
	userIdStr, ok := userIdInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Fetch translations from the database
	var translations []models.TranslationResponse

	if err := db.Model(&models.Translation{}).
		Select("id, phrase, input_language, output_languages, translation_result, created_at").
		Where("user_id = ?", userId).
		Scan(&translations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch translations"})
		return
	}

	// Send the response as JSON
	c.JSON(http.StatusOK, translations)
}

// DeleteTranslation deletes a translation from the database
func DeleteTranslation(c *gin.Context, db *gorm.DB) {
	// Retrieve the user ID from the context set by the middleware
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert userId from string to uuid.UUID
	userIdStr, ok := userIdInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve the translation ID from the URL parameter
	translationIdStr := c.Params.ByName("translation_id")

	fmt.Printf("translationIdStr: %s\n", translationIdStr)
	translationId, err := uuid.Parse(translationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid translation ID"})
		return
	}

	// Fetch the translation from the database
	var translation models.Translation
	if err := db.First(&translation, "id = ? AND user_id = ?", translationId, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Translation not found"})
		return
	}

	// Delete the translation from the database
	if err := db.Delete(&translation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete translation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Translation deleted"})
}
