package main

import (
	"leetcode_backend/controllers"
	"leetcode_backend/models"
	"leetcode_backend/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:1234@tcp(mysql:3306)/questions_db?charset=utf8mb4&parseTime=True&loc=Local"
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
