package main

import (
    "leetcode_backend/config"  
    "leetcode_backend/controllers"
    "leetcode_backend/models"
    "leetcode_backend/repositories"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
	"os"
)

func main() {
    config.ConnectDatabase()

	// Retrieve the database user and password from environment variables
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    // Construct DSN (Data Source Name) for MySQL connection
    dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort +")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect to database")
    }

    db.AutoMigrate(&models.Question{})

    questionRepo := &repository.QuestionRepository{DB: db}
    questionCtrl := &controllers.QuestionController{Repo: questionRepo}

    router := gin.Default()

    router.GET("/questions", questionCtrl.GetQuestions)
    router.POST("/questions", questionCtrl.CreateQuestion)
    router.GET("/questions/:code", questionCtrl.GetQuestion)
    router.PUT("/questions/:code", questionCtrl.UpdateQuestion)
    router.DELETE("/questions/:code", questionCtrl.DeleteQuestion)

    router.Run(":8080")
}
