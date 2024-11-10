package repository

import (
    "gorm.io/gorm"
    "leetcode_backend/models"
)

type QuestionRepository struct {
    DB *gorm.DB
}

func (repo *QuestionRepository) GetAllQuestions() ([]models.Question, error) {
    var questions []models.Question
    err := repo.DB.Find(&questions).Error
    return questions, err
}

func (repo *QuestionRepository) CreateQuestion(question *models.Question) error {
    return repo.DB.Create(question).Error
}

func (repo *QuestionRepository) GetQuestionByCode(code int) (*models.Question, error) {
    var question models.Question
    err := repo.DB.First(&question, code).Error
    return &question, err
}

func (repo *QuestionRepository) UpdateQuestion(question *models.Question) error {
    return repo.DB.Save(question).Error
}

func (repo *QuestionRepository) DeleteQuestion(code int) error {
    return repo.DB.Delete(&models.Question{}, code).Error
}
