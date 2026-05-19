package model

import "time"

type TokenBlacklist struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Token     string    `gorm:"type:varchar(512);not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (TokenBlacklist) TableName() string {
	return "token_blacklist"
}