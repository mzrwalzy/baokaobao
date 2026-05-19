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

func (s *QuestionService) ListPublicQuestionBanks(page, pageSize int) ([]model.QuestionBank, int64, error) {
	return s.repo.ListPublicQuestionBanks(page, pageSize)
}

func (s *QuestionService) GetPublicQuestionBank(id int64) (*model.QuestionBank, error) {
	bank, err := s.repo.GetPublicQuestionBank(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrBankNotFound
		}
		return nil, err
	}
	return bank, nil
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

func (s *QuestionService) ListQuestions(bankID int64, questionType string, page, pageSize int) ([]model.Question, int64, error) {
	return s.repo.ListQuestions(bankID, questionType, page, pageSize)
}

func (s *QuestionService) ListPublicQuestions(bankID int64, questionType string, page, pageSize int) ([]model.Question, int64, error) {
	return s.repo.ListPublicQuestions(bankID, questionType, page, pageSize)
}

func (s *QuestionService) GetQuestion(id int64) (*model.Question, error) {
	return s.repo.GetQuestion(id)
}

func (s *QuestionService) GetPublicQuestion(id int64) (*model.Question, error) {
	question, err := s.repo.GetPublicQuestion(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrQuestionNotFound
		}
		return nil, err
	}
	return question, nil
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

func (s *QuestionService) EnsureUserCanAccessBank(userID, bankID int64) (*model.QuestionBank, error) {
	bank, err := s.GetPublicQuestionBank(bankID)
	if err != nil {
		return nil, err
	}

	if bank.Price <= 0 {
		return bank, nil
	}

	hasAccess, err := s.CheckUserBankAccess(userID, bankID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, model.ErrBankNotPurchased
	}

	return bank, nil
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

func (s *QuestionService) GetBankAccessStatus(userID, bankID int64) (*model.QuestionBankAccessResponse, error) {
	bank, err := s.GetPublicQuestionBank(bankID)
	if err != nil {
		return nil, err
	}

	resp := &model.QuestionBankAccessResponse{
		BankID:           bank.ID,
		Price:            bank.Price,
		RequiresPurchase: bank.Price > 0,
		HasAccess:        bank.Price <= 0,
	}

	if bank.Price <= 0 {
		resp.AccessDescription = "免费题库，可直接刷题"
		return resp, nil
	}

	hasAccess, err := s.CheckUserBankAccess(userID, bankID)
	if err != nil {
		return nil, err
	}

	resp.HasAccess = hasAccess
	if hasAccess {
		resp.AccessDescription = "已开通，可直接刷题"
	} else {
		resp.AccessDescription = "未开通，需要购买权限"
	}

	return resp, nil
}
