package controllers

import (
	"fmt"
	"leetcodeBackend/models"
	"leetcodeBackend/repositories"
	"leetcodeBackend/services"
	"net/http"
	"github.com/gin-gonic/gin"
)

type QuestionController struct {
	questionRepo      repositories.QuestionRepositoryInterface
	executionService *services.ExecutionService
} 

func NewQuestionController(questionRepo repositories.QuestionRepositoryInterface, executionService *services.ExecutionService) *QuestionController {
	return &QuestionController{
			questionRepo:      questionRepo,
			executionService: executionService,
	}
}

func (qc *QuestionController) GetQuestions(c *gin.Context) {
	var questions []models.Question
	if err := qc.questionRepo.GetAllQuestions(&questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

func (qc *QuestionController) CreateQuestion(c *gin.Context) {
	var newQuestion models.Question
	if err := c.ShouldBindJSON(&newQuestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := qc.questionRepo.CreateQuestion(&newQuestion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error creating question: %s", err)})
		return
	}

	c.JSON(http.StatusCreated, newQuestion)
}

func (qc *QuestionController) GetQuestion(c *gin.Context) {
	id := c.Param("code")
	question, err := qc.questionRepo.GetQuestionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Question not found: %v", err)})
		return
	}
	c.JSON(http.StatusOK, question)
}

func (qc *QuestionController) UpdateQuestion(c *gin.Context) {
	id := c.Param("code")
	var updatedQuestion models.Question
	if err := c.ShouldBindJSON(&updatedQuestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := qc.questionRepo.UpdateQuestion(id, &updatedQuestion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating question: %s", err)})
		return
	}

	c.JSON(http.StatusOK, updatedQuestion)
}

func (qc *QuestionController) DeleteQuestion(c *gin.Context) {
	id := c.Param("code")
	if err := qc.questionRepo.DeleteQuestion(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error deleting question: %s", err)})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (qc *QuestionController) ExecuteUserCode(c *gin.Context) {
    var request struct {
        UserExecutionCode string        `json:"user_execution_code" binding:"required"`
        Language         string        `json:"language" binding:"required,oneof=go python"`
        QuestionCode     int          `json:"question_code" binding:"required"`
        Arguments       []interface{} `json:"arguments" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    imageName, err := qc.executionService.CreateDockerImage(
        request.UserExecutionCode,
        request.Language,
        request.QuestionCode,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error creating Docker image: %v", err)})
        return
    }

    output, err := qc.executionService.RunOnKubernetes(imageName, request.Arguments)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error running code: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"output": output})
}