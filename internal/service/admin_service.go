package service

import (
	"baokaobao/internal/model"
	"baokaobao/internal/repository"
)

type AdminService struct {
	repo *Repository
}

func NewAdminService(repo *Repository) *AdminService {
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
