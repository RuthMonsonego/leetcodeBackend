package repositories

import (
    "gorm.io/gorm"
    "leetcode_backend/models"
)

// QuestionRepository manages CRUD operations for the Question model.
type QuestionRepository struct {
    DB *gorm.DB
}

// GetAllQuestions retrieves all questions from the database.
func (repo *QuestionRepository) GetAllQuestions() ([]models.Question, error) {
    var questions []models.Question
    err := repo.DB.Find(&questions).Error
    return questions, err
}

// CreateQuestion adds a new question to the database.
func (repo *QuestionRepository) CreateQuestion(question *models.Question) error {
    return repo.DB.Create(question).Error
}

// GetQuestionByCode finds a question by its code.
func (repo *QuestionRepository) GetQuestionByCode(code int) (*models.Question, error) {
    var question models.Question
    err := repo.DB.First(&question, code).Error
    return &question, err
}

// UpdateQuestion updates an existing question's details.
func (repo *QuestionRepository) UpdateQuestion(question *models.Question) error {
    return repo.DB.Save(question).Error
}

// DeleteQuestion removes a question from the database by its code.
func (repo *QuestionRepository) DeleteQuestion(code int) error {
    return repo.DB.Delete(&models.Question{}, code).Error
}
