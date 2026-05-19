package repository

import (
	"baokaobao/internal/model"
)

func (r *Repository) CreateFeedback(feedback *model.Feedback) error {
	return r.db.Create(feedback).Error
}

func (r *Repository) ListFeedbacks(page, pageSize int, status int8, feedbackType string) ([]model.Feedback, int64, error) {
	var feedbacks []model.Feedback
	var total int64

	db := r.db.Model(&model.Feedback{})
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	if feedbackType != "" {
		db = db.Where("type = ?", feedbackType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Order("id desc").Find(&feedbacks).Error
	return feedbacks, total, err
}

func (r *Repository) GetFeedback(id int64) (*model.Feedback, error) {
	var feedback model.Feedback
	err := r.db.First(&feedback, id).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *Repository) UpdateFeedbackStatus(id int64, status int8) error {
	return r.db.Model(&model.Feedback{}).Where("id = ?", id).Update("status", status).Error
}
