package repository

import (
	"baokaobao/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) CreateUserAnswer(answer *model.UserAnswer) error {
	return r.db.Create(answer).Error
}

func (r *Repository) GetLatestUserAnswer(userID, questionID int64) (*model.UserAnswer, error) {
	var answer model.UserAnswer
	err := r.db.Where("user_id = ? AND question_id = ?", userID, questionID).
		Order("id desc").
		First(&answer).Error
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *Repository) SaveUserAnswer(answer *model.UserAnswer) error {
	return r.db.Save(answer).Error
}

func (r *Repository) AddToWrongQuestions(userID, questionID int64) error {
	var existing model.WrongQuestion
	err := r.db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	wq := &model.WrongQuestion{
		UserID:     userID,
		QuestionID: questionID,
		AddedAt:    time.Now(),
	}
	return r.db.Create(wq).Error
}

func (r *Repository) ListUserAnswers(userID int64, page, pageSize int) ([]model.UserAnswer, int64, error) {
	var answers []model.UserAnswer
	var total int64

	if err := r.db.Model(&model.UserAnswer{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).
		Preload("Question").
		Preload("Question.Options").
		Offset(offset).
		Limit(pageSize).
		Order("answered_at desc").
		Find(&answers).Error
	return answers, total, err
}

func (r *Repository) ListWrongQuestions(userID int64, page, pageSize int) ([]model.WrongQuestion, int64, error) {
	var wqs []model.WrongQuestion
	var total int64

	if err := r.db.Model(&model.WrongQuestion{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).
		Preload("Question").
		Preload("Question.Options").
		Offset(offset).
		Limit(pageSize).
		Order("added_at desc").
		Find(&wqs).Error

	return wqs, total, err
}

func (r *Repository) CreateExamRecord(record *model.ExamRecord) error {
	return r.db.Create(record).Error
}

func (r *Repository) ListExamRecords(userID int64, page, pageSize int) ([]model.ExamRecord, int64, error) {
	var records []model.ExamRecord
	var total int64

	r.db.Model(&model.ExamRecord{}).Where("user_id = ?", userID).Count(&total)
	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).
		Preload("Bank", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Offset(offset).Limit(pageSize).Order("submitted_at desc").Find(&records).Error
	return records, total, err
}

func (r *Repository) UpdateUserScore(userID int64) error {
	var totalScore, totalQuestions, correctCount int64

	err := r.db.Model(&model.UserAnswer{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(score), 0), COUNT(*), COALESCE(SUM(is_correct), 0)").
		Row().Scan(&totalScore, &totalQuestions, &correctCount)
	if err != nil {
		return err
	}

	score := &model.Score{
		UserID:        userID,
		TotalScore:    int(totalScore),
		TotalQuestion: int(totalQuestions),
		CorrectCount:  int(correctCount),
		QuizDate:      time.Now(),
	}

	existing, err := r.GetUserScore(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(score).Error
		}
		return err
	}
	return r.db.Model(&model.Score{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"total_score":    score.TotalScore,
			"total_question": score.TotalQuestion,
			"correct_count":  score.CorrectCount,
			"quiz_date":      score.QuizDate,
		}).Error
}

func (r *Repository) GetUserScore(userID int64) (*model.Score, error) {
	var score model.Score
	err := r.db.Where("user_id = ?", userID).First(&score).Error
	if err != nil {
		return nil, err
	}
	return &score, nil
}

func (r *Repository) GetRanking(page, pageSize int) ([]model.RankingResponse, error) {
	var results []model.RankingResponse
	offset := (page - 1) * pageSize

	err := r.db.Table("scores").
		Select("scores.id, scores.user_id, users.nickname, users.avatar_url, scores.total_score, scores.correct_count, scores.total_question").
		Joins("LEFT JOIN users ON users.id = scores.user_id").
		Where("scores.total_question > 0").
		Order("scores.total_score DESC, (scores.correct_count * 100.0 / scores.total_question) DESC").
		Offset(offset).Limit(pageSize).
		Scan(&results).Error

	for i := range results {
		if results[i].TotalQuestion > 0 {
			results[i].CorrectRate = float64(results[i].CorrectCount) / float64(results[i].TotalQuestion) * 100
		}
		results[i].Rank = int64(offset + i + 1)
	}

	return results, err
}

func (r *Repository) CountUserAnswers(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.UserAnswer{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *Repository) CountUserExams(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.ExamRecord{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *Repository) CountWrongQuestions(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.WrongQuestion{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *Repository) CountStudyDays(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.UserAnswer{}).
		Where("user_id = ? AND answered_at IS NOT NULL AND answered_at <> ?", userID, time.Time{}).
		Distinct("DATE(answered_at)").
		Count(&count).Error
	return count, err
}