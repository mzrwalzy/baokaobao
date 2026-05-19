package service

import (
	"baokaobao/internal/model"
	"baokaobao/internal/repository"
)

type FeedbackService struct {
	repo *repository.Repository
}

func NewFeedbackService(repo *repository.Repository) *FeedbackService {
	return &FeedbackService{repo: repo}
}

func (s *FeedbackService) CreateFeedback(userID int64, feedbackType, content, contact string) (*model.Feedback, error) {
	feedback := &model.Feedback{
		UserID:  userID,
		Type:    feedbackType,
		Content: content,
		Contact: contact,
		Status:  0,
	}
	if err := s.repo.CreateFeedback(feedback); err != nil {
		return nil, err
	}
	return feedback, nil
}

func (s *FeedbackService) ListFeedbacks(page, pageSize int, status int8, feedbackType string) ([]model.Feedback, int64, error) {
	return s.repo.ListFeedbacks(page, pageSize, status, feedbackType)
}

func (s *FeedbackService) UpdateFeedbackStatus(id int64, status int8) error {
	return s.repo.UpdateFeedbackStatus(id, status)
}
