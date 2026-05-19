package repository

import (
	"baokaobao/internal/model"
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
