package migrations

import (
	"baokaobao/internal/model"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.AdminUser{},
		&model.QuestionBank{},
		&model.Question{},
		&model.QuestionOption{},
		&model.UserAnswer{},
		&model.WrongQuestion{},
		&model.Score{},
		&model.ExamRecord{},
		&model.UserBankAccess{},
	)
}

func CreateIndexes(db *gorm.DB) error {
	return nil
}
