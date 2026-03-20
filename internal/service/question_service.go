package service

import (
	"errors"

	"baokaobao/internal/model"
	"baokaobao/internal/repository"

	"gorm.io/gorm"
)

type QuestionService struct {
	repo *repository.Repository
}

func NewQuestionService(repo *repository.Repository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) ListQuestionBanks(page, pageSize int) ([]model.QuestionBank, int64, error) {
	return s.repo.ListQuestionBanks(page, pageSize)
}

func (s *QuestionService) GetQuestionBank(id int64) (*model.QuestionBank, error) {
	return s.repo.GetQuestionBank(id)
}

func (s *QuestionService) CreateQuestionBank(bank *model.QuestionBank) error {
	return s.repo.CreateQuestionBank(bank)
}

func (s *QuestionService) UpdateQuestionBank(bank *model.QuestionBank) error {
	return s.repo.UpdateQuestionBank(bank)
}

func (s *QuestionService) DeleteQuestionBank(id int64) error {
	return s.repo.DeleteQuestionBank(id)
}

func (s *QuestionService) ListQuestions(bankID int64, page, pageSize int) ([]model.Question, int64, error) {
	return s.repo.ListQuestions(bankID, page, pageSize)
}

func (s *QuestionService) GetQuestion(id int64) (*model.Question, error) {
	return s.repo.GetQuestion(id)
}

func (s *QuestionService) CreateQuestion(question *model.Question) error {
	return s.repo.CreateQuestion(question)
}

func (s *QuestionService) UpdateQuestion(question *model.Question) error {
	return s.repo.UpdateQuestion(question)
}

func (s *QuestionService) DeleteQuestion(id int64) error {
	return s.repo.DeleteQuestion(id)
}

func (s *QuestionService) RandomQuestions(bankID int64, count int) ([]model.Question, error) {
	return s.repo.RandomQuestions(bankID, count)
}

func (s *QuestionService) CheckUserBankAccess(userID, bankID int64) (bool, error) {
	_, err := s.repo.GetUserBankAccess(userID, bankID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *QuestionService) GrantBankAccess(userID, bankID int64) error {
	access := &model.UserBankAccess{
		UserID: userID,
		BankID: bankID,
	}
	return s.repo.CreateUserBankAccess(access)
}
