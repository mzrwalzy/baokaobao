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
	totalUsers, _ := s.repo.CountUsers()
	totalQuestions, _ := s.repo.CountQuestions()
	totalAnswers, _ := s.repo.CountAnswers()
	todayUsers, _ := s.repo.CountTodayUsers()

	return &model.DashboardStats{
		TotalUsers:     totalUsers,
		TotalQuestions: totalQuestions,
		TotalAnswers:   totalAnswers,
		TodayUsers:     todayUsers,
	}, nil
}

func (s *AdminService) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	return s.repo.ListUsers(page, pageSize)
}

func (s *AdminService) GetUser(id int64) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *AdminService) UpdateUserStatus(id int64, status int8) error {
	return s.repo.UpdateUserStatus(id, status)
}

func (s *AdminService) GetUserStats() (*model.UserStatsResponse, error) {
	total, _ := s.repo.CountUsers()
	today, _ := s.repo.CountTodayUsers()

	return &model.UserStatsResponse{
		Total: total,
		Today: today,
	}, nil
}

func (s *AdminService) GetQuestionStats() (*model.QuestionStatsResponse, error) {
	total, _ := s.repo.CountQuestions()

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
