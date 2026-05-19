package service

import (
	"baokaobao/internal/model"
	"baokaobao/internal/repository"
)

type AuditLogService struct {
	repo *repository.Repository
}

func NewAuditLogService(repo *repository.Repository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

func (s *AuditLogService) CreateAuditLog(adminID int64, adminName, action, target string, targetID int64, detail, ip string) error {
	log := &model.AuditLog{
		AdminID:   adminID,
		AdminName: adminName,
		Action:    action,
		Target:    target,
		TargetID:  targetID,
		Detail:    detail,
		IP:        ip,
	}
	return s.repo.CreateAuditLog(log)
}

func (s *AuditLogService) ListAuditLogs(page, pageSize int, adminName, action string, startTime, endTime *string) ([]model.AuditLog, int64, error) {
	return s.repo.ListAuditLogs(page, pageSize, adminName, action, startTime, endTime)
}
