package repository

import (
	"baokaobao/internal/model"
)

func (r *Repository) CreateAuditLog(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *Repository) ListAuditLogs(page, pageSize int, adminName, action string, startTime, endTime *string) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	db := r.db.Model(&model.AuditLog{})
	if adminName != "" {
		db = db.Where("admin_name LIKE ?", "%"+adminName+"%")
	}
	if action != "" {
		db = db.Where("action = ?", action)
	}
	if startTime != nil && *startTime != "" {
		db = db.Where("created_at >= ?", *startTime)
	}
	if endTime != nil && *endTime != "" {
		db = db.Where("created_at <= ?", *endTime)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Order("id desc").Find(&logs).Error
	return logs, total, err
}
