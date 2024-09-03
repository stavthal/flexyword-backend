package main

import (
	"log"

	"flexyword.io/backend/controllers"
	"flexyword.io/backend/db"
	"flexyword.io/backend/models"
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

	// AutoMigrate the schema
    err = Db.AutoMigrate(&models.User{}, &models.Translation{})
    if err != nil {
        log.Fatalf("Error migrating database schema: %v", err)
    }


	// Create the Gin router
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to translations API!",
		})
	})

	// Create a route group for translation-related routes
    translationGroup := r.Group("/api/translate")
    {
        translationGroup.POST("/", func(c *gin.Context) {
            controllers.TranslatePhrase(c, Db)
        })
    }

	usersGroup := r.Group("/api/users")
	{
		usersGroup.POST("/register", func(c *gin.Context) {
			controllers.RegisterUser(c, Db)
		})
		usersGroup.POST("/login", func(c *gin.Context) {
			controllers.LoginUser(c, Db)
		})
	}



	
	r.Run(":8080") // listen and serve on 
}