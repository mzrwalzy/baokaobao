package model

import "time"

type Feedback struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index;not null" json:"user_id"`
	Type      string    `gorm:"type:varchar(32);default:general" json:"type"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Contact   string    `gorm:"type:varchar(128)" json:"contact"`
	Status    int8      `gorm:"type:tinyint(1);default:0" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Feedback) TableName() string {
	return "feedbacks"
}
