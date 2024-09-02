package main

import (
	"flexyword.io/backend/controllers"
	"flexyword.io/backend/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


var Db = db.NewConnection()


func main() {

	// Load the environmental variables
	err := godotenv.Load(".env")



	if err != nil {
		panic("Error loading .env file")
	}


	// Create the Gin router
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// Create a route group for translation-related routes
    translationGroup := r.Group("/api/translate")
    {
        translationGroup.POST("/", func(c *gin.Context) {
            controllers.TranslatePhrase(c, Db)
        })
    }


	r.Run(":8080") // listen and serve on 
}