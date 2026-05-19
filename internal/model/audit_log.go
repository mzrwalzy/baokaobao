package model

import "time"

// AuditLog 记录后台管理操作日志
type AuditLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID    int64     `gorm:"index;not null" json:"admin_id"`
	AdminName  string    `gorm:"type:varchar(64)" json:"admin_name"`
	Action     string    `gorm:"type:varchar(32);not null" json:"action"`
	Target     string    `gorm:"type:varchar(64)" json:"target"`
	TargetID   int64     `gorm:"index" json:"target_id"`
	Detail     string    `gorm:"type:text" json:"detail"`
	IP         string    `gorm:"type:varchar(64)" json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
