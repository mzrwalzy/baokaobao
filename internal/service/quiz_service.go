package service

import (
	"baokaobao/internal/model"
	"baokaobao/internal/repository"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const defaultQuestionScore = 5

type QuizService struct {
	repo            *repository.Repository
	questionService *QuestionService
}

type validatedExamAnswer struct {
	Request  model.SubmitAnswerItem
	Question *model.Question
}

func NewQuizService(repo *repository.Repository, questionService *QuestionService) *QuizService {
	return &QuizService{repo: repo, questionService: questionService}
}

func (s *QuizService) SubmitAnswer(userID, questionID int64, myAnswer string) (*model.SubmitAnswerResponse, error) {
	question, err := s.questionService.GetPublicQuestion(questionID)
	if err != nil {
		return nil, err
	}
	if _, err := s.questionService.EnsureUserCanAccessBank(userID, question.BankID); err != nil {
		return nil, err
	}

	isCorrect := s.compareAnswer(question.Answer, myAnswer, question.Type)
	score := 0
	if isCorrect {
		score = defaultQuestionScore
	}

	answer := &model.UserAnswer{
		UserID:     userID,
		QuestionID: questionID,
		MyAnswer:   myAnswer,
		IsCorrect:  boolToInt8(isCorrect),
		Score:      score,
		AnsweredAt: time.Now(),
	}

	existing, err := s.repo.GetLatestUserAnswer(userID, questionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		existing.MyAnswer = answer.MyAnswer
		existing.IsCorrect = answer.IsCorrect
		existing.Score = answer.Score
		existing.AnsweredAt = answer.AnsweredAt
		if err := s.repo.SaveUserAnswer(existing); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.CreateUserAnswer(answer); err != nil {
			return nil, err
		}
	}

	if !isCorrect {
		if err := s.repo.AddToWrongQuestions(userID, questionID); err != nil {
			zap.S().Errorf("AddToWrongQuestions error: %v", err)
		}
	}

	if err := s.repo.UpdateUserScore(userID); err != nil {
		zap.S().Errorf("UpdateUserScore error: %v", err)
	}

	return &model.SubmitAnswerResponse{
		IsCorrect:     isCorrect,
		Score:         score,
		CorrectAnswer: question.Answer,
		Analysis:      question.Analysis,
	}, nil
}

func (s *QuizService) SubmitExam(userID, bankID int64, answers []model.SubmitAnswerItem, duration int) (*model.ExamSubmitResponse, error) {
	if _, err := s.questionService.EnsureUserCanAccessBank(userID, bankID); err != nil {
		return nil, err
	}

	validAnswers, err := s.validateExamAnswers(bankID, answers)
	if err != nil {
		return nil, err
	}

	var totalScore, correctCount, totalQuestions int

	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		for _, item := range validAnswers {
			ans := item.Request
			question := item.Question

			isCorrect := s.compareAnswer(question.Answer, ans.MyAnswer, question.Type)
			score := 0
			if isCorrect {
				score = defaultQuestionScore
				correctCount++
			}
			totalScore += score
			totalQuestions++

			answer := &model.UserAnswer{
				UserID:     userID,
				QuestionID: ans.QuestionID,
				MyAnswer:   ans.MyAnswer,
				IsCorrect:  boolToInt8(isCorrect),
				Score:      score,
				AnsweredAt: time.Now(),
			}
			if err := tx.Create(answer).Error; err != nil {
				return err
			}

			if !isCorrect {
				// Check if already exists
				var existing model.WrongQuestion
				err := tx.Where("user_id = ? AND question_id = ?", userID, ans.QuestionID).First(&existing).Error
				if err == nil {
					continue // already exists
				}
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				wq := &model.WrongQuestion{
					UserID:     userID,
					QuestionID: ans.QuestionID,
					AddedAt:    time.Now(),
				}
				if err := tx.Create(wq).Error; err != nil {
					return err
				}
			}
		}

		examRecord := &model.ExamRecord{
			UserID:        userID,
			BankID:        bankID,
			TotalScore:    totalScore,
			TotalQuestion: totalQuestions,
			CorrectCount:  correctCount,
			Duration:      duration,
		}
		if err := tx.Create(examRecord).Error; err != nil {
			return err
		}

		// Update user score
		var scoreTotal, scoreQuestions, scoreCorrect int64
		tx.Model(&model.UserAnswer{}).Where("user_id = ?", userID).
			Select("COALESCE(SUM(score), 0), COUNT(*), COALESCE(SUM(is_correct), 0)").
			Row().Scan(&scoreTotal, &scoreQuestions, &scoreCorrect)

		score := &model.Score{
			UserID:        userID,
			TotalScore:    int(scoreTotal),
			TotalQuestion: int(scoreQuestions),
			CorrectCount:  int(scoreCorrect),
			QuizDate:      time.Now(),
		}

		// Check if score record exists
		var existing model.Score
		err := tx.Where("user_id = ?", userID).First(&existing).Error
		if err == nil {
			return tx.Model(&model.Score{}).
				Where("id = ?", existing.ID).
				Updates(map[string]interface{}{
					"total_score":    score.TotalScore,
					"total_question": score.TotalQuestion,
					"correct_count":  score.CorrectCount,
					"quiz_date":      score.QuizDate,
				}).Error
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return tx.Create(score).Error
	})

	if err != nil {
		return nil, err
	}

	return &model.ExamSubmitResponse{
		TotalScore:   totalScore,
		CorrectCount: correctCount,
	}, nil
}

func (s *QuizService) validateExamAnswers(bankID int64, answers []model.SubmitAnswerItem) ([]validatedExamAnswer, error) {
	validAnswers := make([]validatedExamAnswer, 0, len(answers))

	for _, ans := range answers {
		question, err := s.questionService.GetPublicQuestion(ans.QuestionID)
		if err != nil {
			return nil, err
		}
		if question.BankID != bankID {
			return nil, model.ErrNoAccessToBank
		}

		validAnswers = append(validAnswers, validatedExamAnswer{
			Request:  ans,
			Question: question,
		})
	}

	return validAnswers, nil
}

func (s *QuizService) GetHistory(userID int64, page, pageSize int) ([]model.UserAnswer, int64, error) {
	return s.repo.ListUserAnswers(userID, page, pageSize)
}

func (s *QuizService) GetWrongQuestions(userID int64, page, pageSize int) ([]model.WrongQuestion, int64, error) {
	return s.repo.ListWrongQuestions(userID, page, pageSize)
}

func (s *QuizService) AddToWrong(userID, questionID int64) error {
	return s.repo.AddToWrongQuestions(userID, questionID)
}

func (s *QuizService) GetExamRecords(userID int64, page, pageSize int) ([]model.ExamRecord, int64, error) {
	return s.repo.ListExamRecords(userID, page, pageSize)
}

func (s *QuizService) GetMyPurchasedBanks(userID int64) ([]model.QuestionBank, error) {
	return s.repo.GetUserPurchasedBanks(userID)
}

func (s *QuizService) compareAnswer(correct, given, qType string) bool {
	switch qType {
	case "single", "truefalse", "judge":
		return correct == given
	case "multiple":
		return compareMultipleChoice(correct, given)
	default:
		return correct == given
	}
}

func compareMultipleChoice(correct, given string) bool {
	correctSet := stringToSet(correct)
	givenSet := stringToSet(given)

	if len(correctSet) != len(givenSet) {
		return false
	}

	for k := range correctSet {
		if _, ok := givenSet[k]; !ok {
			return false
		}
	}
	return true
}

func stringToSet(s string) map[rune]bool {
	set := make(map[rune]bool)
	for _, c := range s {
		set[c] = true
	}
	return set
}

func boolToInt8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}
