package model

import "time"

type AdminBankPermission struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID   int64     `gorm:"index;not null" json:"admin_id"`
	BankID    int64     `gorm:"index;not null" json:"bank_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (AdminBankPermission) TableName() string {
	return "admin_bank_permissions"
}
