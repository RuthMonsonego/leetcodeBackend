package controllers

import (
    "leetcode_backend/models"
    "leetcode_backend/config"
    "net/http"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
    router.GET("/questions", getQuestions)
    router.POST("/questions", createQuestion)
    router.GET("/questions/:id", getQuestion)
    router.PUT("/questions/:id", updateQuestion)
    router.DELETE("/questions/:id", deleteQuestion)
}

func getQuestions(c *gin.Context) {
    var questions []models.Question
    if err := config.DB.Find(&questions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, questions)
}

func createQuestion(c *gin.Context) {
    var question models.Question
    if err := c.ShouldBindJSON(&question); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := config.DB.Create(&question).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, question)
}

func getQuestion(c *gin.Context) {
    id := c.Param("id")
    var question models.Question
    if err := config.DB.First(&question, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, question)
}

func updateQuestion(c *gin.Context) {
    id := c.Param("id")
    var question models.Question
    if err := config.DB.First(&question, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    if err := c.ShouldBindJSON(&question); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := config.DB.Save(&question).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, question)
}

func deleteQuestion(c *gin.Context) {
    id := c.Param("id")
    if err := config.DB.Delete(&models.Question{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}
