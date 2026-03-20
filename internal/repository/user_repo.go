package repository

import (
	"baokaobao/internal/model"
	"errors"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetUserByOpenID(openid string) (*model.User, error) {
	var user model.User
	err := r.db.Where("openid = ?", openid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) UpdateUserPhone(userID int64, phone string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("phone", phone).Error
}

func (r *Repository) GetAdminByUsername(username string) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := r.db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *Repository) GetAdminByID(id int64) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := r.db.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *Repository) UpdateAdmin(admin *model.AdminUser) error {
	return r.db.Save(admin).Error
}

func (r *Repository) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	r.db.Model(&model.User{}).Count(&total)
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("id desc").Find(&users).Error
	return users, total, err
}

func (r *Repository) GetUserBankAccess(userID, bankID int64) (*model.UserBankAccess, error) {
	var access model.UserBankAccess
	err := r.db.Where("user_id = ? AND bank_id = ?", userID, bankID).First(&access).Error
	if err != nil {
		return nil, err
	}
	return &access, nil
}

func (r *Repository) CreateUserBankAccess(access *model.UserBankAccess) error {
	return r.db.Create(access).Error
}

func (r *Repository) GetUserPurchasedBanks(userID int64) ([]model.QuestionBank, error) {
	var banks []model.QuestionBank
	err := r.db.Table("question_banks").
		Joins("INNER JOIN user_bank_access ON question_banks.id = user_bank_access.bank_id").
		Where("user_bank_access.user_id = ?", userID).
		Find(&banks).Error
	return banks, err
}

func (r *Repository) GetTodayNewUsers() ([]model.User, error) {
	var users []model.User
	today := time.Now().Format("2006-01-02")
	err := r.db.Where("DATE(created_at) = ?", today).Find(&users).Error
	return users, err
}

func (r *Repository) UpdateUserProfile(userID int64, nickname, avatarURL string) error {
	updates := map[string]interface{}{}
	if nickname != "" {
		updates["nickname"] = nickname
	}
	if avatarURL != "" {
		updates["avatar_url"] = avatarURL
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *Repository) ListQuestionBanks(page, pageSize int) ([]model.QuestionBank, int64, error) {
	var banks []model.QuestionBank
	var total int64

	r.db.Model(&model.QuestionBank{}).Count(&total)
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("id desc").Find(&banks).Error
	return banks, total, err
}

func (r *Repository) GetQuestionBank(id int64) (*model.QuestionBank, error) {
	var bank model.QuestionBank
	err := r.db.First(&bank, id).Error
	if err != nil {
		return nil, err
	}
	return &bank, nil
}

func (r *Repository) CreateQuestionBank(bank *model.QuestionBank) error {
	return r.db.Create(bank).Error
}

func (r *Repository) UpdateQuestionBank(bank *model.QuestionBank) error {
	return r.db.Save(bank).Error
}

func (r *Repository) DeleteQuestionBank(id int64) error {
	return r.db.Delete(&model.QuestionBank{}, id).Error
}

func (r *Repository) ListQuestions(bankID int64, page, pageSize int) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := r.db.Model(&model.Question{})
	if bankID > 0 {
		query = query.Where("bank_id = ?", bankID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id desc").Find(&questions).Error

	for i := range questions {
		r.db.Where("question_id = ?", questions[i].ID).Find(&questions[i].Options)
	}

	return questions, total, err
}

func (r *Repository) GetQuestion(id int64) (*model.Question, error) {
	var question model.Question
	err := r.db.First(&question, id).Error
	if err != nil {
		return nil, err
	}
	r.db.Where("question_id = ?", question.ID).Find(&question.Options)
	return &question, nil
}

func (r *Repository) CreateQuestion(question *model.Question) error {
	return r.db.Create(question).Error
}

func (r *Repository) UpdateQuestion(question *model.Question) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(question).Error; err != nil {
			return err
		}
		tx.Where("question_id = ?", question.ID).Delete(&model.QuestionOption{})
		for _, opt := range question.Options {
			opt.QuestionID = question.ID
			if err := tx.Create(&opt).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) DeleteQuestion(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("question_id = ?", id).Delete(&model.QuestionOption{})
		return tx.Delete(&model.Question{}, id).Error
	})
}

func (r *Repository) RandomQuestions(bankID int64, count int) ([]model.Question, error) {
	var questions []model.Question
	query := r.db.Model(&model.Question{}).Where("status = 1")
	if bankID > 0 {
		query = query.Where("bank_id = ?", bankID)
	}

	if err := query.Find(&questions).Error; err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	shuffle := rand.Perm(len(questions))
	result := make([]model.Question, 0)

	for i := 0; i < count && i < len(questions); i++ {
		q := questions[shuffle[i]]
		r.db.Where("question_id = ?", q.ID).Find(&q.Options)
		result = append(result, q)
	}

	return result, nil
}

func (r *Repository) CreateUserAnswer(answer *model.UserAnswer) error {
	return r.db.Create(answer).Error
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

	r.db.Model(&model.UserAnswer{}).Where("user_id = ?", userID).Count(&total)
	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("answered_at desc").Find(&answers).Error
	return answers, total, err
}

func (r *Repository) ListWrongQuestions(userID int64, page, pageSize int) ([]model.WrongQuestion, int64, error) {
	var wqs []model.WrongQuestion
	var total int64

	r.db.Model(&model.WrongQuestion{}).Where("user_id = ?", userID).Count(&total)
	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("added_at desc").Find(&wqs).Error

	for i := range wqs {
		r.db.Where("question_id = ?", wqs[i].QuestionID).Find(&wqs[i].Question.Options)
	}

	return wqs, total, err
}

func (r *Repository) CreateExamRecord(record *model.ExamRecord) error {
	return r.db.Create(record).Error
}

func (r *Repository) UpdateUserScore(userID int64) error {
	var totalScore, totalQuestions, correctCount int64

	r.db.Model(&model.UserAnswer{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(score), 0), COUNT(*), COALESCE(SUM(is_correct), 0)").
		Row().Scan(&totalScore, &totalQuestions, &correctCount)

	score := &model.Score{
		UserID:        userID,
		TotalScore:    int(totalScore),
		TotalQuestion: int(totalQuestions),
		CorrectCount:  int(correctCount),
		QuizDate:      time.Now(),
	}

	existing, _ := r.GetUserScore(userID)
	if existing != nil {
		score.ID = existing.ID
		return r.db.Save(score).Error
	}
	return r.db.Create(score).Error
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
	r.db.Model(&model.UserAnswer{}).Where("user_id = ?", userID).Count(&count)
	return count, nil
}

func (r *Repository) CountUserExams(userID int64) (int64, error) {
	var count int64
	r.db.Model(&model.ExamRecord{}).Where("user_id = ?", userID).Count(&count)
	return count, nil
}

func (r *Repository) CountUsers() (int64, error) {
	var count int64
	r.db.Model(&model.User{}).Count(&count)
	return count, nil
}

func (r *Repository) CountQuestions() (int64, error) {
	var count int64
	r.db.Model(&model.Question{}).Count(&count)
	return count, nil
}

func (r *Repository) CountAnswers() (int64, error) {
	var count int64
	r.db.Model(&model.UserAnswer{}).Count(&count)
	return count, nil
}

func (r *Repository) CountTodayUsers() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	r.db.Model(&model.User{}).Where("DATE(created_at) = ?", today).Count(&count)
	return count, nil
}

func (r *Repository) UpdateUserStatus(userID int64, status int8) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (r *Repository) CreateAdmin(admin *model.AdminUser) error {
	return r.db.Create(admin).Error
}
