package service

import (
	"baokaobao/internal/model"
	"baokaobao/internal/repository"
)

type ScoreService struct {
	repo *repository.Repository
}

func NewScoreService(repo *repository.Repository) *ScoreService {
	return &ScoreService{repo: repo}
}

func (s *ScoreService) GetMyScore(userID int64) (*model.Score, error) {
	return s.repo.GetUserScore(userID)
}

func (s *ScoreService) GetRanking(page, pageSize int) ([]model.RankingResponse, error) {
	return s.repo.GetRanking(page, pageSize)
}

func (s *ScoreService) GetStats(userID int64) (*model.StatsResponse, error) {
	score, err := s.repo.GetUserScore(userID)
	if err != nil {
		return &model.StatsResponse{}, nil
	}

	totalAnswers, _ := s.repo.CountUserAnswers(userID)
	totalExams, _ := s.repo.CountUserExams(userID)

	correctRate := 0.0
	if totalAnswers > 0 {
		correctRate = float64(score.CorrectCount) / float64(score.TotalQuestion) * 100
	}

	return &model.StatsResponse{
		TotalScore:     score.TotalScore,
		TotalQuestions: score.TotalQuestion,
		CorrectCount:   score.CorrectCount,
		CorrectRate:    correctRate,
		TotalExams:     int(totalExams),
	}, nil
}
