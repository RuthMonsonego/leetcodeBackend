package repositories

import (
	"leetcodeBackend/models"
	"fmt"
	"gorm.io/gorm"
)

type QuestionRepositoryInterface interface {
	GetAllQuestions(questions *[]models.Question) error
	CreateQuestion(question *models.Question) error
	GetQuestionByID(id string) (*models.Question, error)
	UpdateQuestion(id string, question *models.Question) error
	DeleteQuestion(id string) error
}

type QuestionRepository struct {
	db *gorm.DB
}

var _ QuestionRepositoryInterface = (*QuestionRepository)(nil)

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (qr *QuestionRepository) GetAllQuestions(questions *[]models.Question) error {
	return qr.db.Find(questions).Error
}

func (qr *QuestionRepository) CreateQuestion(question *models.Question) error {
	return qr.db.Create(question).Error
}

func (qr *QuestionRepository) GetQuestionByID(id string) (*models.Question, error) {
	var question models.Question
	if err := qr.db.First(&question, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("question not found: %v", err)
	}
	
	fmt.Printf("Found question: %+v\n", question)
	
	return &question, nil
}

func (qr *QuestionRepository) UpdateQuestion(id string, question *models.Question) error {
	return qr.db.Model(&models.Question{}).Where("id = ?", id).Updates(question).Error
}

func (qr *QuestionRepository) DeleteQuestion(id string) error {
	return qr.db.Delete(&models.Question{}, "id = ?", id).Error
}
