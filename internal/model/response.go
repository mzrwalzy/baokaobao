package model

type UserProfileResponse struct {
	ID        int64  `json:"id"`
	OpenID    string `json:"openid"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Phone     string `json:"phone"`
}

type QuestionDetailResponse struct {
	ID         int64            `json:"id"`
	BankID     int64            `json:"bank_id"`
	Title      string           `json:"title"`
	Content    string           `json:"content"`
	Images     []string         `json:"images"`
	Type       string           `json:"type"`
	Difficulty int8             `json:"difficulty"`
	Options    []OptionResponse `json:"options"`
}

type OptionResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Image string `json:"image,omitempty"`
}

type SubmitAnswerRequest struct {
	QuestionID int64  `json:"question_id" binding:"required"`
	MyAnswer   string `json:"my_answer" binding:"required"`
}

type SubmitAnswerResponse struct {
	IsCorrect     bool   `json:"is_correct"`
	Score         int    `json:"score"`
	CorrectAnswer string `json:"correct_answer,omitempty"`
	Analysis      string `json:"analysis,omitempty"`
}

type ExamSubmitRequest struct {
	BankID   int64              `json:"bank_id" binding:"required"`
	Duration int                `json:"duration"`
	Answers  []SubmitAnswerItem `json:"answers" binding:"required"`
}

type SubmitAnswerItem struct {
	QuestionID int64  `json:"question_id" binding:"required"`
	MyAnswer   string `json:"my_answer" binding:"required"`
}

type ExamSubmitResponse struct {
	TotalScore   int `json:"total_score"`
	CorrectCount int `json:"correct_count"`
}

type RankingResponse struct {
	Rank          int64   `json:"rank"`
	UserID        int64   `json:"user_id"`
	Nickname      string  `json:"nickname"`
	AvatarURL     string  `json:"avatar_url"`
	TotalScore    int     `json:"total_score"`
	CorrectCount  int     `json:"correct_count"`
	TotalQuestion int     `json:"total_question"`
	CorrectRate   float64 `json:"correct_rate"`
}

type StatsResponse struct {
	TotalScore     int     `json:"total_score"`
	TotalQuestions int     `json:"total_questions"`
	CorrectCount   int     `json:"correct_count"`
	CorrectRate    float64 `json:"correct_rate"`
	TotalExams     int     `json:"total_exams"`
}

type QuestionBankResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CoverImage  string  `json:"cover_image"`
	QuestionNum int     `json:"question_num"`
	Status      int8    `json:"status"`
}

type DashboardStats struct {
	TotalUsers     int64 `json:"total_users"`
	TotalQuestions int64 `json:"total_questions"`
	TotalAnswers   int64 `json:"total_answers"`
	TodayUsers     int64 `json:"today_users"`
}

type UserStatsResponse struct {
	Total int64 `json:"total"`
	Today int64 `json:"today"`
}

type QuestionStatsResponse struct {
	Total int64 `json:"total"`
}
