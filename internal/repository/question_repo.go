package repository

import (
	"baokaobao/internal/model"

	"gorm.io/gorm"
)

func (r *Repository) ListQuestionBanks(page, pageSize int) ([]model.QuestionBank, int64, error) {
	var banks []model.QuestionBank
	var total int64

	if err := r.db.Model(&model.QuestionBank{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("id desc").Find(&banks).Error
	return banks, total, err
}

func (r *Repository) ListPublicQuestionBanks(page, pageSize int) ([]model.QuestionBank, int64, error) {
	var banks []model.QuestionBank
	var total int64

	query := r.db.Model(&model.QuestionBank{}).Where("status = ?", 1)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id desc").Find(&banks).Error
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

func (r *Repository) GetPublicQuestionBank(id int64) (*model.QuestionBank, error) {
	var bank model.QuestionBank
	err := r.db.Where("id = ? AND status = ?", id, 1).First(&bank).Error
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

func (r *Repository) ListQuestions(bankID int64, questionType string, page, pageSize int) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := r.db.Model(&model.Question{}).Preload("Options")
	if bankID > 0 {
		query = query.Where("bank_id = ?", bankID)
	}
	if questionType != "" {
		query = query.Where("type = ?", questionType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id desc").Find(&questions).Error

	return questions, total, err
}

func (r *Repository) ListPublicQuestions(bankID int64, questionType string, page, pageSize int) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := r.db.Model(&model.Question{}).
		Preload("Options").
		Joins("INNER JOIN question_banks ON question_banks.id = questions.bank_id").
		Where("questions.status = ? AND question_banks.status = ?", 1, 1)

	if bankID > 0 {
		query = query.Where("questions.bank_id = ?", bankID)
	}
	if questionType != "" {
		query = query.Where("questions.type = ?", questionType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("questions.id desc").Find(&questions).Error

	return questions, total, err
}

func (r *Repository) GetQuestion(id int64) (*model.Question, error) {
	var question model.Question
	err := r.db.Preload("Options").First(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *Repository) GetPublicQuestion(id int64) (*model.Question, error) {
	var question model.Question
	err := r.db.Model(&model.Question{}).
		Preload("Options").
		Joins("INNER JOIN question_banks ON question_banks.id = questions.bank_id").
		Where("questions.id = ? AND questions.status = ? AND question_banks.status = ?", id, 1, 1).
		First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *Repository) CreateQuestion(question *model.Question) error {
	return r.db.Create(question).Error
}

func (r *Repository) UpdateQuestion(question *model.Question) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("question_id = ?", question.ID).Delete(&model.QuestionOption{}).Error; err != nil {
			return err
		}
		if err := tx.Save(question).Error; err != nil {
			return err
		}
		if len(question.Options) > 0 {
			for i := range question.Options {
				question.Options[i].QuestionID = question.ID
			}
			if err := tx.CreateInBatches(question.Options, 100).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) DeleteQuestion(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("question_id = ?", id).Delete(&model.QuestionOption{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Question{}, id).Error
	})
}

func (r *Repository) RandomQuestions(bankID int64, count int) ([]model.Question, error) {
	var questions []model.Question
	err := r.db.Model(&model.Question{}).
		Preload("Options").
		Joins("INNER JOIN question_banks ON question_banks.id = questions.bank_id").
		Where("questions.bank_id = ? AND questions.status = ? AND question_banks.status = ?", bankID, 1, 1).
		Order("RAND()").
		Limit(count).
		Find(&questions).Error
	return questions, err
}
