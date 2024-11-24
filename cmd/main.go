package main

import (
    "time"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "leetcodeBackend/config"
    "leetcodeBackend/controllers"
    "leetcodeBackend/repositories"
    "leetcodeBackend/services"
)

func main() {
    config.ConnectDatabase()

    questionRepository := repositories.NewQuestionRepository(config.DB)

    executionService, err := services.NewExecutionService(questionRepository)
    if err != nil {
        panic(err)
    }

    questionController := controllers.NewQuestionController(questionRepository, executionService)

    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
    r.GET("/questions", questionController.GetQuestions)
    r.POST("/questions", questionController.CreateQuestion)
    r.GET("/questions/:code", questionController.GetQuestion)
    r.PUT("/questions/:code", questionController.UpdateQuestion)
    r.PUT("/questions/:code/test", questionController.ExecuteUserCode)
    r.DELETE("/questions/:code", questionController.DeleteQuestion)

    r.Run(":8080")
}