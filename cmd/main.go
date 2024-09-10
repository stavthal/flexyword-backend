package main

import (
	"fmt"
	"log"
	"time"

	"flexyword.io/backend/controllers"
	"flexyword.io/backend/db"
	"flexyword.io/backend/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)


var Db *gorm.DB

func init() {
	// Load the environmental variables
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		fmt.Println("Proceeding without loading .env file")
	}

	// Initialize the database connection
	var dbErr error
	Db, dbErr = db.NewConnection()
	if dbErr != nil {
		log.Fatalf("Failed to connect to the database: %v", dbErr)
	}
}



func main() {
	// Create the Gin router
	r := gin.Default()

	// Setup CORS middleware
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "https://your-nuxt-frontend.com"},
        AllowMethods:     []string{"GET","POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to translations API!",
		})
	})

	// Create a route group for translation-related routes
    translationGroup := r.Group("/api/translate")
	translationGroup.Use(middlewares.AuthMiddleware())
    {
        translationGroup.POST("/new", func(c *gin.Context) {
            controllers.TranslatePhrase(c, Db)
        })
		translationGroup.GET("/fetch", func(c *gin.Context) {
			controllers.GetTranslations(c, Db)
		})
		translationGroup.DELETE("/delete/:translation_id", func(c *gin.Context) {
			controllers.DeleteTranslation(c, Db)
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