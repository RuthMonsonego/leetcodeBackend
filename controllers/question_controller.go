package controllers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "leetcode_backend/models"
    "leetcode_backend/repositories"
)

type QuestionController struct {
    Repo *repository.QuestionRepository
}

// Get all questions
func (ctrl *QuestionController) GetQuestions(c *gin.Context) {
    questions, err := ctrl.Repo.GetAllQuestions()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
        return
    }
    c.JSON(http.StatusOK, questions)
}

// Create a new question
func (ctrl *QuestionController) CreateQuestion(c *gin.Context) {
    var question models.Question
    if err := c.ShouldBindJSON(&question); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    if err := ctrl.Repo.CreateQuestion(&question); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
        return
    }
    c.JSON(http.StatusCreated, question)
}

// Get a single question by code
func (ctrl *QuestionController) GetQuestion(c *gin.Context) {
    code, err := strconv.Atoi(c.Param("code"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
        return
    }
    question, err := ctrl.Repo.GetQuestionByCode(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
        return
    }
    c.JSON(http.StatusOK, question)
}

// Update an existing question
func (ctrl *QuestionController) UpdateQuestion(c *gin.Context) {
    code, err := strconv.Atoi(c.Param("code"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
        return
    }

    var updatedQuestion models.Question
    if err := c.ShouldBindJSON(&updatedQuestion); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    updatedQuestion.Code = code
    if err := ctrl.Repo.UpdateQuestion(&updatedQuestion); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question"})
        return
    }
    c.JSON(http.StatusOK, updatedQuestion)
}

// Delete a question by code
func (ctrl *QuestionController) DeleteQuestion(c *gin.Context) {
    code, err := strconv.Atoi(c.Param("code"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
        return
    }

    if err := ctrl.Repo.DeleteQuestion(code); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
        return
    }
    c.JSON(http.StatusNoContent, nil)
}
