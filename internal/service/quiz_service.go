package service

import (
	"baokaobao/internal/model"
	"baokaobao/internal/repository"
)

type QuizService struct {
	repo *Repository
}

func NewQuizService(repo *Repository) *QuizService {
	return &QuizService{repo: repo}
}

func (s *QuizService) SubmitAnswer(userID, questionID int64, myAnswer string) (*model.SubmitAnswerResponse, error) {
	question, err := s.repo.GetQuestion(questionID)
	if err != nil {
		return nil, err
	}

	isCorrect := s.compareAnswer(question.Answer, myAnswer, question.Type)
	score := 0
	if isCorrect {
		score = 5
	}

	answer := &model.UserAnswer{
		UserID:     userID,
		QuestionID: questionID,
		MyAnswer:   myAnswer,
		IsCorrect:  boolToInt8(isCorrect),
		Score:      score,
	}
	if err := s.repo.CreateUserAnswer(answer); err != nil {
		return nil, err
	}

	if !isCorrect {
		s.repo.AddToWrongQuestions(userID, questionID)
	}

	s.repo.UpdateUserScore(userID)

	return &model.SubmitAnswerResponse{
		IsCorrect:     isCorrect,
		Score:         score,
		CorrectAnswer: question.Answer,
		Analysis:      question.Analysis,
	}, nil
}

func (s *QuizService) SubmitExam(userID, bankID int64, answers []model.SubmitAnswerItem, duration int) (*model.ExamSubmitResponse, error) {
	var totalScore, correctCount, totalQuestions int

	for _, ans := range answers {
		question, err := s.repo.GetQuestion(ans.QuestionID)
		if err != nil {
			continue
		}

		isCorrect := s.compareAnswer(question.Answer, ans.MyAnswer, question.Type)
		score := 0
		if isCorrect {
			score = 5
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
		}
		s.repo.CreateUserAnswer(answer)

		if !isCorrect {
			s.repo.AddToWrongQuestions(userID, ans.QuestionID)
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
	s.repo.CreateExamRecord(examRecord)

	s.repo.UpdateUserScore(userID)

	return &model.ExamSubmitResponse{
		TotalScore:   totalScore,
		CorrectCount: correctCount,
	}, nil
}

func (s *QuizService) GetHistory(userID int64, page, pageSize int) ([]model.UserAnswer, int64, error) {
	return s.repo.ListUserAnswers(userID, page, pageSize)
}

func (s *QuizService) GetWrongQuestions(userID int64, page, pageSize int) ([]model.WrongQuestion, int64, error) {
	return s.repo.ListWrongQuestions(userID, page, pageSize)
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
