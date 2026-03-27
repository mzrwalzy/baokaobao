package model

import (
	"time"
)

type User struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	OpenID    string     `gorm:"column:openid;type:varchar(128);uniqueIndex;not null" json:"openid"`
	UnionID   string     `gorm:"column:unionid;type:varchar(128);index" json:"unionid"`
	Nickname  string     `gorm:"type:varchar(64)" json:"nickname"`
	AvatarURL string     `gorm:"type:varchar(256)" json:"avatar_url"`
	Phone     string     `gorm:"type:varchar(32);uniqueIndex" json:"phone"`
	Status    int8       `gorm:"type:tinyint(1);default:1;index" json:"status"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type AdminUser struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"type:varchar(256);not null" json:"-"`
	Nickname     string     `gorm:"type:varchar(64)" json:"nickname"`
	Role         string     `gorm:"type:varchar(32);default:admin" json:"role"`
	Status       int8       `gorm:"type:tinyint(1);default:1" json:"status"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (AdminUser) TableName() string {
	return "admin_users"
}

type QuestionBank struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(128);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);default:0" json:"price"`
	CoverImage  string    `gorm:"type:varchar(256)" json:"cover_image"`
	Status      int8      `gorm:"type:tinyint(1);default:1" json:"status"`
	QuestionNum int       `gorm:"type:int;default:0" json:"question_num"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (QuestionBank) TableName() string {
	return "question_banks"
}

type Question struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BankID     int64     `gorm:"index;not null" json:"bank_id"`
	Title      string    `gorm:"type:varchar(256);not null" json:"title"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Images     string    `gorm:"type:text" json:"images"`
	Answer     string    `gorm:"type:text;not null" json:"answer"`
	Analysis   string    `gorm:"type:text" json:"analysis"`
	Type       string    `gorm:"type:varchar(32);default:single" json:"type"`
	Difficulty int8      `gorm:"type:tinyint(1);default:2" json:"difficulty"`
	Status     int8      `gorm:"type:tinyint(1);default:1" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Options []QuestionOption `gorm:"foreignKey:QuestionID" json:"options,omitempty"`
}

func (Question) TableName() string {
	return "questions"
}

type QuestionOption struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionID  int64  `gorm:"index;not null" json:"question_id"`
	OptionKey   string `gorm:"type:varchar(8);not null" json:"option_key"`
	OptionValue string `gorm:"type:text;not null" json:"option_value"`
	OptionImage string `gorm:"type:varchar(256)" json:"option_image"`
}

func (QuestionOption) TableName() string {
	return "question_options"
}

type UserAnswer struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"index;not null" json:"user_id"`
	QuestionID int64     `gorm:"index;not null" json:"question_id"`
	MyAnswer   string    `gorm:"type:text;not null" json:"my_answer"`
	IsCorrect  int8      `gorm:"type:tinyint(1);default:0" json:"is_correct"`
	Score      int       `gorm:"type:int;default:0" json:"score"`
	AnsweredAt time.Time `json:"answered_at"`

	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Question Question `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

func (UserAnswer) TableName() string {
	return "user_answers"
}

type WrongQuestion struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"index;not null" json:"user_id"`
	QuestionID int64     `gorm:"index;not null" json:"question_id"`
	AddedAt    time.Time `json:"added_at"`

	Question Question `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

func (WrongQuestion) TableName() string {
	return "wrong_questions"
}

type Score struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64     `gorm:"index;not null" json:"user_id"`
	TotalScore    int       `gorm:"type:int;default:0" json:"total_score"`
	TotalQuestion int       `gorm:"type:int;default:0" json:"total_question"`
	CorrectCount  int       `gorm:"type:int;default:0" json:"correct_count"`
	QuizDate      time.Time `json:"quiz_date"`
	CreatedAt     time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Score) TableName() string {
	return "scores"
}

type ExamRecord struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64     `gorm:"index;not null" json:"user_id"`
	BankID        int64     `gorm:"index" json:"bank_id"`
	TotalScore    int       `gorm:"type:int;default:0" json:"total_score"`
	TotalQuestion int       `gorm:"type:int;default:0" json:"total_question"`
	CorrectCount  int       `gorm:"type:int;default:0" json:"correct_count"`
	Duration      int       `gorm:"type:int;default:0" json:"duration"`
	SubmittedAt   time.Time `json:"submitted_at"`
	CreatedAt     time.Time `json:"created_at"`
}

func (ExamRecord) TableName() string {
	return "exam_records"
}

type UserBankAccess struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index;not null" json:"user_id"`
	BankID    int64     `gorm:"index;not null" json:"bank_id"`
	BoughtAt  time.Time `json:"bought_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserBankAccess) TableName() string {
	return "user_bank_access"
}
