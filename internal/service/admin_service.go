package service

import (
	"baokaobao/internal/model"
	"baokaobao/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	repo *repository.Repository
}

func NewAdminService(repo *repository.Repository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) Dashboard() (*model.DashboardStats, error) {
	totalUsers, err := s.repo.CountUsers()
	if err != nil {
		return nil, err
	}
	totalQuestions, err := s.repo.CountQuestions()
	if err != nil {
		return nil, err
	}
	totalAnswers, err := s.repo.CountAnswers()
	if err != nil {
		return nil, err
	}
	todayUsers, err := s.repo.CountTodayUsers()
	if err != nil {
		return nil, err
	}

	return &model.DashboardStats{
		TotalUsers:     totalUsers,
		TotalQuestions: totalQuestions,
		TotalAnswers:   totalAnswers,
		TodayUsers:     todayUsers,
	}, nil
}

func (s *AdminService) ListUsers(page, pageSize int, keyword string) ([]model.User, int64, error) {
	return s.repo.ListUsers(page, pageSize, keyword)
}

func (s *AdminService) GetUser(id int64) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *AdminService) UpdateUserStatus(id int64, status int8) error {
	return s.repo.UpdateUserStatus(id, status)
}

func (s *AdminService) GetUserStats() (*model.UserStatsResponse, error) {
	total, err := s.repo.CountUsers()
	if err != nil {
		return nil, err
	}
	today, err := s.repo.CountTodayUsers()
	if err != nil {
		return nil, err
	}

	return &model.UserStatsResponse{
		Total: total,
		Today: today,
	}, nil
}

func (s *AdminService) GetQuestionStats() (*model.QuestionStatsResponse, error) {
	total, err := s.repo.CountQuestions()
	if err != nil {
		return nil, err
	}

	return &model.QuestionStatsResponse{
		Total: total,
	}, nil
}

func (s *AdminService) CreateAdminUser(username, password, nickname string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.AdminUser{
		Username:     username,
		PasswordHash: string(hash),
		Nickname:     nickname,
		Role:         "admin",
		Status:       1,
	}

	return s.repo.CreateAdmin(admin)
}

func (s *AdminService) GetUserPurchasedBanks(userID int64) ([]model.QuestionBank, error) {
	return s.repo.GetUserPurchasedBanks(userID)
}

func (s *AdminService) ListAdminUsers(page, pageSize int) ([]model.AdminUser, int64, error) {
	return s.repo.ListAdminUsers(page, pageSize)
}

func (s *AdminService) GetAdminUserByID(id int64) (*model.AdminUser, error) {
	return s.repo.GetAdminByID(id)
}

func (s *AdminService) CreateAdminUserWithRole(username, password, nickname, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.AdminUser{
		Username:     username,
		PasswordHash: string(hash),
		Nickname:     nickname,
		Role:         role,
		Status:       1,
	}

	return s.repo.CreateAdmin(admin)
}

func (s *AdminService) UpdateAdminUser(admin *model.AdminUser) error {
	return s.repo.UpdateAdminUser(admin)
}

func (s *AdminService) ResetAdminPassword(id int64, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.ResetAdminPassword(id, string(hash))
}

func (s *AdminService) DeleteAdminUser(id int64) error {
	return s.repo.DeleteAdminUser(id)
}

func (s *AdminService) ListAdminBankPermissions(adminID int64) ([]model.AdminBankPermission, error) {
	return s.repo.ListAdminBankPermissions(adminID)
}

func (s *AdminService) GrantAdminBankAccess(adminID, bankID int64) error {
	return s.repo.GrantAdminBankAccess(adminID, bankID)
}

func (s *AdminService) RevokeAdminBankAccess(adminID, bankID int64) error {
	return s.repo.RevokeAdminBankAccess(adminID, bankID)
}
