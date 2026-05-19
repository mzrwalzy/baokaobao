package repository

import (
	"baokaobao/internal/model"
)

type BankStat struct {
	BankID            int64   `json:"bank_id"`
	BankName          string  `json:"bank_name"`
	QuestionCount     int64   `json:"question_count"`
	PurchasedCount    int64   `json:"purchased_count"`
	TotalAnswers      int64   `json:"total_answers"`
	TotalExams        int64   `json:"total_exams"`
	AvgCorrectRate    float64 `json:"avg_correct_rate"`
}

type BankDetailStat struct {
	BankID          int64           `json:"bank_id"`
	BankName        string          `json:"bank_name"`
	TotalQuestions  int64           `json:"total_questions"`
	TotalUsers      int64           `json:"total_users"`
	TotalAnswers    int64           `json:"total_answers"`
	TotalExams      int64           `json:"total_exams"`
	AvgCorrectRate  float64         `json:"avg_correct_rate"`
	DifficultyDist  []DifficultyDist `json:"difficulty_dist"`
	DailyTrend      []DailyTrend     `json:"daily_trend"`
}

type DifficultyDist struct {
	Difficulty int8  `json:"difficulty"`
	Count      int64 `json:"count"`
}

type DailyTrend struct {
	Date         string  `json:"date"`
	AnswerCount  int64   `json:"answer_count"`
	CorrectRate  float64 `json:"correct_rate"`
}

func (r *Repository) GetBankStatsList(page, pageSize int, sortBy, sortOrder string) ([]BankStat, int64, error) {
	var stats []BankStat
	var total int64

	// Count total banks
	if err := r.db.Model(&model.QuestionBank{}).Where("status = ?", 1).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.QuestionBank{}).
		Select(`
			question_banks.id as bank_id,
			question_banks.name as bank_name,
			COUNT(DISTINCT questions.id) as question_count,
			COUNT(DISTINCT user_bank_access.user_id) as purchased_count,
			COUNT(DISTINCT user_answers.id) as total_answers,
			COUNT(DISTINCT exam_records.id) as total_exams,
			IFNULL(AVG(user_answers.is_correct) * 100, 0) as avg_correct_rate
		`).
		Joins("LEFT JOIN questions ON questions.bank_id = question_banks.id AND questions.status = ?", 1).
		Joins("LEFT JOIN user_bank_access ON user_bank_access.bank_id = question_banks.id").
		Joins("LEFT JOIN user_answers ON user_answers.question_id = questions.id").
		Joins("LEFT JOIN exam_records ON exam_records.bank_id = question_banks.id").
		Where("question_banks.status = ?", 1).
		Group("question_banks.id, question_banks.name")

	// Apply sorting
	orderClause := "avg_correct_rate DESC"
	if sortBy != "" && sortOrder != "" {
		validSortColumns := map[string]bool{
			"question_count": true, "purchased_count": true,
			"total_answers": true, "total_exams": true, "avg_correct_rate": true,
		}
		if validSortColumns[sortBy] {
			orderDir := "DESC"
			if sortOrder == "ascending" {
				orderDir = "ASC"
			}
			orderClause = sortBy + " " + orderDir
		}
	}
	query = query.Order(orderClause)

	err := query.Offset(offset).Limit(pageSize).Scan(&stats).Error
	return stats, total, err
}

func (r *Repository) GetBankStats(bankID int64) (*BankDetailStat, error) {
	var stat BankDetailStat

	// Basic info
	err := r.db.Model(&model.QuestionBank{}).
		Select("question_banks.id as bank_id, question_banks.name as bank_name, COUNT(DISTINCT questions.id) as total_questions").
		Joins("LEFT JOIN questions ON questions.bank_id = question_banks.id AND questions.status = ?", 1).
		Where("question_banks.id = ? AND question_banks.status = ?", bankID, 1).
		Group("question_banks.id, question_banks.name").
		Scan(&stat).Error
	if err != nil {
		return nil, err
	}

	// Total users
	if err := r.db.Model(&model.UserBankAccess{}).Where("bank_id = ?", bankID).Count(&stat.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Total answers
	if err := r.db.Model(&model.UserAnswer{}).
		Joins("INNER JOIN questions ON questions.id = user_answers.question_id").
		Where("questions.bank_id = ?", bankID).
		Count(&stat.TotalAnswers).Error; err != nil {
		return nil, err
	}

	// Total exams
	if err := r.db.Model(&model.ExamRecord{}).Where("bank_id = ?", bankID).Count(&stat.TotalExams).Error; err != nil {
		return nil, err
	}

	// Avg correct rate
	var avgRate float64
	if err := r.db.Raw(`
		SELECT IFNULL(AVG(user_answers.is_correct) * 100, 0)
		FROM user_answers
		INNER JOIN questions ON questions.id = user_answers.question_id
		WHERE questions.bank_id = ?
	`, bankID).Scan(&avgRate).Error; err != nil {
		return nil, err
	}
	stat.AvgCorrectRate = avgRate

	// Difficulty distribution
	if err := r.db.Raw(`
		SELECT questions.difficulty, COUNT(*) as count
		FROM questions
		WHERE questions.bank_id = ? AND questions.status = 1
		GROUP BY questions.difficulty
		ORDER BY questions.difficulty
	`, bankID).Scan(&stat.DifficultyDist).Error; err != nil {
		return nil, err
	}

	// Daily trend (last 30 days)
	if err := r.db.Raw(`
		SELECT 
			DATE(user_answers.answered_at) as date,
			COUNT(*) as answer_count,
			IFNULL(AVG(user_answers.is_correct) * 100, 0) as correct_rate
		FROM user_answers
		INNER JOIN questions ON questions.id = user_answers.question_id
		WHERE questions.bank_id = ? AND user_answers.answered_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
		GROUP BY DATE(user_answers.answered_at)
		ORDER BY date
	`, bankID).Scan(&stat.DailyTrend).Error; err != nil {
		return nil, err
	}

	return &stat, nil
}
