package repository

import (
	"baokaobao/internal/model"
	"time"
)

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

func (r *Repository) ListUsers(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	if keyword != "" {
		query = query.Where("nickname LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id desc").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	// 填充已购题库数量
	if len(users) > 0 {
		userIDs := make([]int64, len(users))
		for i, u := range users {
			userIDs[i] = u.ID
		}
		var counts []struct {
			UserID int64
			Count  int64
		}
		r.db.Raw(`
			SELECT user_id, COUNT(*) as count
			FROM user_bank_access
			WHERE user_id IN ?
			GROUP BY user_id
		`, userIDs).Scan(&counts)
		countMap := make(map[int64]int64)
		for _, c := range counts {
			countMap[c.UserID] = c.Count
		}
		for i := range users {
			users[i].BankCount = countMap[users[i].ID]
		}
	}

	return users, total, nil
}

func (r *Repository) GetTodayNewUsers() ([]model.User, error) {
	var users []model.User
	today := time.Now().Format("2006-01-02")
	err := r.db.Where("DATE(created_at) = ?", today).Find(&users).Error
	return users, err
}

func (r *Repository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *Repository) CountQuestions() (int64, error) {
	var count int64
	err := r.db.Model(&model.Question{}).Count(&count).Error
	return count, err
}

func (r *Repository) CountAnswers() (int64, error) {
	var count int64
	err := r.db.Model(&model.UserAnswer{}).Count(&count).Error
	return count, err
}

func (r *Repository) CountTodayUsers() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&model.User{}).Where("DATE(created_at) = ?", today).Count(&count).Error
	return count, err
}

func (r *Repository) UpdateUserStatus(userID int64, status int8) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (r *Repository) CreateAdmin(admin *model.AdminUser) error {
	return r.db.Create(admin).Error
}

func (r *Repository) ListAdminUsers(page, pageSize int) ([]model.AdminUser, int64, error) {
	var admins []model.AdminUser
	var total int64

	if err := r.db.Model(&model.AdminUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("id desc").Find(&admins).Error
	return admins, total, err
}

func (r *Repository) UpdateAdminUser(admin *model.AdminUser) error {
	return r.db.Model(&model.AdminUser{}).Where("id = ?", admin.ID).Updates(map[string]interface{}{
		"nickname": admin.Nickname,
		"role":     admin.Role,
		"status":   admin.Status,
	}).Error
}

func (r *Repository) ResetAdminPassword(id int64, passwordHash string) error {
	return r.db.Model(&model.AdminUser{}).Where("id = ?", id).Update("password_hash", passwordHash).Error
}

func (r *Repository) DeleteAdminUser(id int64) error {
	return r.db.Delete(&model.AdminUser{}, id).Error
}

func (r *Repository) ListAdminBankPermissions(adminID int64) ([]model.AdminBankPermission, error) {
	var perms []model.AdminBankPermission
	err := r.db.Where("admin_id = ?", adminID).Find(&perms).Error
	return perms, err
}

func (r *Repository) GrantAdminBankAccess(adminID, bankID int64) error {
	perm := &model.AdminBankPermission{
		AdminID: adminID,
		BankID:  bankID,
	}
	return r.db.Create(perm).Error
}

func (r *Repository) RevokeAdminBankAccess(adminID, bankID int64) error {
	return r.db.Where("admin_id = ? AND bank_id = ?", adminID, bankID).Delete(&model.AdminBankPermission{}).Error
}

func (r *Repository) AddToBlacklist(token string, expiresAt time.Time) error {
	bl := &model.TokenBlacklist{
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return r.db.Create(bl).Error
}

func (r *Repository) IsTokenBlacklisted(token string) (bool, error) {
	var count int64
	err := r.db.Model(&model.TokenBlacklist{}).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		Count(&count).Error
	return count > 0, err
}