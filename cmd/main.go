package main

import (
    "leetcode_backend/config"
    "leetcode_backend/controllers"
    "leetcode_backend/repositories"
    "github.com/gin-gonic/gin"
)

func main() {
    config.ConnectDatabase()

    questionRepo := &repositories.QuestionRepository{DB: config.DB}
    questionCtrl := controllers.NewQuestionController(questionRepo)

    router := gin.Default()
    router.GET("/questions", questionCtrl.GetQuestions)
    router.POST("/questions", questionCtrl.CreateQuestion)
    router.GET("/questions/:code", questionCtrl.GetQuestion)
    router.PUT("/questions/:code", questionCtrl.UpdateQuestion)
    router.PUT("/questions/:code/test", questionCtrl.ExecuteUserCode)
    router.DELETE("/questions/:code", questionCtrl.DeleteQuestion)

    router.Run(":8080")
}
