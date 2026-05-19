package service

import (
	"baokaobao/internal/model"
	"baokaobao/internal/pkg/jwt"
	"baokaobao/internal/pkg/wechat"
	"baokaobao/internal/repository"
)

type Service struct {
	Auth      *AuthService
	User      *UserService
	Question  *QuestionService
	Quiz      *QuizService
	Score     *ScoreService
	Admin     *AdminService
	Feedback  *FeedbackService
	AuditLog  *AuditLogService
	BankStats *BankStatsService
}

func NewService(repo *repository.Repository, jwtSDK *jwt.JWT, wechatSDK *wechat.WechatSDK) *Service {
	return &Service{
		Auth:      NewAuthService(repo, jwtSDK, wechatSDK),
		User:      NewUserService(repo),
		Question:  NewQuestionService(repo),
		Quiz:      NewQuizService(repo, NewQuestionService(repo)),
		Score:     NewScoreService(repo),
		Admin:     NewAdminService(repo),
		Feedback:  NewFeedbackService(repo),
		AuditLog:  NewAuditLogService(repo),
		BankStats: NewBankStatsService(repo),
	}
}

func (s *Service) ListUsers(page, pageSize int, keyword string) ([]model.User, int64, error) {
	return s.User.ListUsers(page, pageSize, keyword)
}
