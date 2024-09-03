package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"flexyword.io/backend/models"
	"flexyword.io/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

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

	// Initialize OpenAI client
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	translations := make(map[string]string)

	// Loop through each language and request translation from OpenAI
	for _, lang := range request.Languages {
		prompt := "Translate the following phrase from " + request.InputLanguage + " into " + lang + ": " + request.Phrase

		resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: "gpt-3.5-turbo", // Use a stable, available model. Adjust if needed.
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
			MaxTokens:   4096, // Adjust max tokens based on expected translation length
			Temperature: 0.2, // Low temperature for deterministic results
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get translation from OpenAI: " + err.Error()})
			return
		}

		if len(resp.Choices) > 0 {
			translations[lang] = resp.Choices[0].Message.Content
		} else {
			translations[lang] = ""
		}
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
		OutputLanguages:   string(outputLanguagesJSON), // Store the languages array as JSON string
		TranslationResult: string(translationsJSON),    // Store the translations as JSON string
		UserID:            c.GetUint("userId"),         // Assume userId is stored in context by AuthMiddleware
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
