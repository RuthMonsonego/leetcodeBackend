package controllers

import (
    "net/http"
    "strconv"
    "log"
    "github.com/gin-gonic/gin"
    "leetcode_backend/models"
    "leetcode_backend/repositories"
    "leetcode_backend/services"
)

type QuestionController struct {
    Repo             *repositories.QuestionRepository
    ExecutionService *services.ExecutionService 
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

// NewQuestionController creates a new QuestionController with the given repository and execution service
func NewQuestionController(repo *repositories.QuestionRepository) *QuestionController {
    executionService, err := services.NewExecutionService()
    if err != nil {
        log.Fatal("Failed to initialize execution service:", err)
    }
    return &QuestionController{
        Repo: repo,
        ExecutionService: executionService,
    }
}

// ExecuteUserCode executes user-submitted code and returns the output
func (ctrl *QuestionController) ExecuteUserCode(c *gin.Context) {
    var req services.ExecuteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    output, err := ctrl.ExecutionService.Execute(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"output": output})
}
