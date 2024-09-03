package controllers

import (
	"context"
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
	InputLanguage	string 			`json:"input_language" binding:"required"`
	Languages    	[]string 		`json:"languages" binding:"required"`
	Phrase	     	string 			`json:"phrase" binding:"required"`
}

type RatesLimitError struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

// TranslatePhrase handles the translation of a phrase into multiple languages
func TranslatePhrase(c *gin.Context, db *gorm.DB) {
	var request TranslateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY")) 


	translations := make(map[string]string)

	// Loop through each language and request translation from OpenAI
	for _, lang := range request.Languages {
		prompt := "Translate the following phrase from " + request.InputLanguage + " into " + lang + ": " + request.Phrase

		resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: "gpt-3.5-turbo",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
			MaxTokens:   60,
			Temperature: 0.2,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		} else {
			if len(resp.Choices) > 0 {
				translations[lang] = resp.Choices[0].Message.Content
			} else {
				translations[lang] = ""
			}
		}

		
	}


	translation := models.Translation {
		Phrase: request.Phrase,
		Translations: translations,
		Language: request.InputLanguage,
	}

	// Store the translation using the translation service
	err := services.StoreTranslation(db, &translation)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
  
	// Send the response as JSON
	c.JSON(http.StatusOK, gin.H{
		"phrase":       request.Phrase,
		"translations": translations,
	})
}


// TODO: Add a handler for listing all translations of a user
// TODO: Add a handler that translates a single phrase and outputs the results into i18n friendly files
// TODO: Add a route handler for premium users to use the application with softer limits
// TODO: Add a middleware to use only for authenticated users
// TODO: Add a way to limit how many requests a user can send each day
// TODO: Add a way to save the new translations to a database, bound to the user's id

