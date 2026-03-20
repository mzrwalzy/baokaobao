package handler

import (
	"baokaobao/internal/model"
	"baokaobao/internal/pkg/wechat"
	"baokaobao/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	svc       *service.Service
	wechatSDK *wechat.WechatSDK
}

func NewHandler(svc *service.Service, wechatSDK *wechat.WechatSDK) *Handler {
	return &Handler{
		svc:       svc,
		wechatSDK: wechatSDK,
	}
}

func (h *Handler) LoginByWechat(c *gin.Context) {
	var req model.LoginByWechatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, "参数错误")
		return
	}

	result, err := h.svc.Auth.LoginByWechat(req.Code)
	if err != nil {
		zap.S().Errorf("LoginByWechat error: %v", err)
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, result)
}

func (h *Handler) DecryptPhone(c *gin.Context) {
	var req model.DecryptPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetInt64("user_id")
	phone, err := h.svc.Auth.DecryptPhone(userID, req.Code)
	if err != nil {
		zap.S().Errorf("DecryptPhone error: %v", err)
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, model.DecryptPhoneResponse{Phone: phone})
}

func (h *Handler) Logout(c *gin.Context) {
	model.Success(c, nil)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	profile, err := h.svc.User.GetProfile(userID)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, profile)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetInt64("user_id")
	if err := h.svc.User.UpdateProfile(userID, req.Nickname, req.AvatarURL); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	model.Success(c, gin.H{"url": "https://example.com/avatar.jpg"})
}

func (h *Handler) ListQuestions(c *gin.Context) {
	bankID := int64(0)
	if b := c.Query("bank_id"); b != "" {
		bankID = parseInt64(b)
	}
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	questions, total, err := h.svc.Question.ListQuestions(bankID, page, pageSize)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Page(c, questions, total, page, pageSize)
}

func (h *Handler) GetQuestion(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	question, err := h.svc.Question.GetQuestion(id)
	if err != nil {
		model.NotFound(c, "题目不存在")
		return
	}
	model.Success(c, question)
}

func (h *Handler) RandomQuestions(c *gin.Context) {
	bankID := int64(0)
	if b := c.Query("bank_id"); b != "" {
		bankID = parseInt64(b)
	}
	count := parseIntDefault(c.Query("count"), 10)

	questions, err := h.svc.Question.RandomQuestions(bankID, count)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, questions)
}

func (h *Handler) SubmitAnswer(c *gin.Context) {
	var req model.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetInt64("user_id")
	result, err := h.svc.Quiz.SubmitAnswer(userID, req.QuestionID, req.MyAnswer)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, result)
}

func (h *Handler) GetQuizHistory(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	answers, total, err := h.svc.Quiz.GetHistory(userID, page, pageSize)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Page(c, answers, total, page, pageSize)
}

func (h *Handler) GetWrongQuestions(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	wqs, total, err := h.svc.Quiz.GetWrongQuestions(userID, page, pageSize)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Page(c, wqs, total, page, pageSize)
}

func (h *Handler) AddToWrongQuestions(c *gin.Context) {
	model.Success(c, nil)
}

func (h *Handler) GetMyScore(c *gin.Context) {
	userID := c.GetInt64("user_id")
	score, err := h.svc.Score.GetMyScore(userID)
	if err != nil {
		model.Success(c, nil)
		return
	}
	model.Success(c, score)
}

func (h *Handler) GetRanking(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	ranking, err := h.svc.Score.GetRanking(page, pageSize)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, ranking)
}

func (h *Handler) GetStats(c *gin.Context) {
	userID := c.GetInt64("user_id")
	stats, err := h.svc.Score.GetStats(userID)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, stats)
}

func (h *Handler) AdminLogin(c *gin.Context) {
	var req model.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, "参数错误")
		return
	}

	result, err := h.svc.Auth.AdminLogin(&req)
	if err != nil {
		zap.S().Errorf("AdminLogin error: %v", err)
		model.Unauthorized(c, "用户名或密码错误")
		return
	}

	model.Success(c, result)
}

func (h *Handler) AdminLogout(c *gin.Context) {
	model.Success(c, nil)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	stats, err := h.svc.Admin.Dashboard()
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, stats)
}

func (h *Handler) ListAllUsers(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	users, total, err := h.svc.Admin.ListUsers(page, pageSize)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Page(c, users, total, page, pageSize)
}

func (h *Handler) GetUserDetail(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	user, err := h.svc.Admin.GetUser(id)
	if err != nil {
		model.NotFound(c, "用户不存在")
		return
	}
	model.Success(c, user)
}

func (h *Handler) UpdateUserStatus(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.Admin.UpdateUserStatus(id, req.Status); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

func (h *Handler) ListQuestionBanks(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	banks, total, err := h.svc.Question.ListQuestionBanks(page, pageSize)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Page(c, banks, total, page, pageSize)
}

func (h *Handler) CreateQuestionBank(c *gin.Context) {
	var bank model.QuestionBank
	if err := c.ShouldBindJSON(&bank); err != nil {
		model.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.Question.CreateQuestionBank(&bank); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, bank)
}

func (h *Handler) UpdateQuestionBank(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	var bank model.QuestionBank
	if err := c.ShouldBindJSON(&bank); err != nil {
		model.BadRequest(c, "参数错误")
		return
	}
	bank.ID = id

	if err := h.svc.Question.UpdateQuestionBank(&bank); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, bank)
}

func (h *Handler) DeleteQuestionBank(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	if err := h.svc.Question.DeleteQuestionBank(id); err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, nil)
}

func (h *Handler) ListAllQuestions(c *gin.Context) {
	bankID := int64(0)
	if b := c.Query("bank_id"); b != "" {
		bankID = parseInt64(b)
	}
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	questions, total, err := h.svc.Question.ListQuestions(bankID, page, pageSize)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Page(c, questions, total, page, pageSize)
}

func (h *Handler) CreateQuestion(c *gin.Context) {
	var req struct {
		model.Question
		Options []model.QuestionOption `json:"options"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, "参数错误")
		return
	}

	req.Question.Options = req.Options
	if err := h.svc.Question.CreateQuestion(&req.Question); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, req.Question)
}

func (h *Handler) UpdateQuestion(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	var req struct {
		model.Question
		Options []model.QuestionOption `json:"options"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, "参数错误")
		return
	}
	req.Question.ID = id
	req.Question.Options = req.Options

	if err := h.svc.Question.UpdateQuestion(&req.Question); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, req.Question)
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	if err := h.svc.Question.DeleteQuestion(id); err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, nil)
}

func (h *Handler) ImportQuestions(c *gin.Context) {
	model.Success(c, gin.H{"imported": 0})
}

func (h *Handler) GetStatsOverview(c *gin.Context) {
	stats, err := h.svc.Admin.Dashboard()
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, stats)
}

func (h *Handler) GetUserStats(c *gin.Context) {
	model.Success(c, gin.H{"total": 0, "today": 0})
}

func (h *Handler) GetQuestionStats(c *gin.Context) {
	model.Success(c, gin.H{"total": 0})
}

func parseInt64(s string) int64 {
	var n int64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int64(c-'0')
		}
	}
	return n
}

func parseIntDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	if n == 0 {
		return defaultVal
	}
	return n
}
