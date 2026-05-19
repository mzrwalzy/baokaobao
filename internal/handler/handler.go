package handler

import (
	"baokaobao/internal/middleware"
	"baokaobao/internal/model"
	"baokaobao/internal/pkg/response"
	"baokaobao/internal/pkg/wechat"
	"baokaobao/internal/service"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	svc       *service.Service
	wechatSDK *wechat.WechatSDK
}

// AvatarMaxSize defines the maximum allowed avatar upload size (2MB)
const AvatarMaxSize = 2 * 1024 * 1024

func NewHandler(svc *service.Service, wechatSDK *wechat.WechatSDK) *Handler {
	return &Handler{
		svc:       svc,
		wechatSDK: wechatSDK,
	}
}

func (h *Handler) LoginByWechat(c *gin.Context) {
	var req model.LoginByWechatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := h.svc.Auth.LoginByWechat(req.Code)
	if err != nil {
		zap.S().Errorf("LoginByWechat error: %v", err)
		response.InternalError(c, "登录失败，请稍后重试")
		return
	}

	response.Success(c, result)
}

func (h *Handler) DecryptPhone(c *gin.Context) {
	var req model.DecryptPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetInt64("user_id")
	phone, err := h.svc.Auth.DecryptPhone(userID, req.Code)
	if err != nil {
		zap.S().Errorf("DecryptPhone error: %v", err)
		response.InternalError(c, "获取手机号失败，请稍后重试")
		return
	}

	response.Success(c, model.DecryptPhoneResponse{Phone: phone})
}

func (h *Handler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			claims, err := middleware.GetJWT().ParseToken(tokenString)
			if err == nil {
				h.svc.Auth.Logout(tokenString, claims.ExpiresAt.Time)
			}
		}
	}
	response.Success(c, nil)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	profile, err := h.svc.User.GetProfile(userID)
	if err != nil {
		zap.S().Errorf("GetProfile error: %v", err)
		response.InternalError(c, "获取用户信息失败")
		return
	}
	response.Success(c, profile)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetInt64("user_id")
	if err := h.svc.User.UpdateProfile(userID, req.Nickname, req.AvatarURL); err != nil {
		zap.S().Errorf("UpdateProfile error: %v", err)
		response.InternalError(c, "更新用户信息失败")
		return
	}

	response.Success(c, nil)
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "please upload a file")
		return
	}
	defer file.Close()

	if header.Size > AvatarMaxSize {
		response.BadRequest(c, "file size must be less than 2MB")
		return
	}

	userID := c.GetInt64("user_id")

	sniff := make([]byte, 512)
	n, err := file.Read(sniff)
	if err != nil && !errors.Is(err, io.EOF) {
		response.BadRequest(c, "读取文件失败")
		return
	}

	contentType := http.DetectContentType(sniff[:n])
	ext := ""
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	default:
		response.BadRequest(c, "仅支持 JPG 或 PNG 图片")
		return
	}

	filename := fmt.Sprintf("avatar_%d%s", userID, ext)
	reader := io.MultiReader(bytes.NewReader(sniff[:n]), file)

	url, err := h.svc.User.UploadAvatar(userID, reader, filename)
	if err != nil {
		zap.S().Errorf("UploadAvatar error: %v", err)
		response.InternalError(c, "上传头像失败")
		return
	}

	response.Success(c, gin.H{"url": url})
}

func (h *Handler) ListQuestions(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var bankID int64
	if b := c.Query("bank_id"); b != "" {
		v, err := strconv.ParseInt(b, 10, 64)
		if err != nil {
			response.BadRequest(c, "参数错误：无效的题库ID")
			return
		}
		bankID = v
	}
	questionType := c.Query("type")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}

	if bankID > 0 {
		if _, err := h.svc.Question.EnsureUserCanAccessBank(userID, bankID); err != nil {
			if errors.Is(err, model.ErrBankNotPurchased) {
				response.Forbidden(c, "请先购买该题库")
				return
			}
			if errors.Is(err, model.ErrBankNotFound) {
				response.NotFound(c, "题库不存在或已下线")
				return
			}
			zap.S().Errorf("ListQuestions access check error: %v", err)
			response.InternalError(c, "获取题目列表失败")
			return
		}
	}

	questions, total, err := h.svc.Question.ListPublicQuestions(bankID, questionType, page, pageSize)
	if err != nil {
		zap.S().Errorf("ListQuestions error: %v", err)
		response.InternalError(c, "获取题目列表失败")
		return
	}

	response.Page(c, questions, total, page, pageSize)
}

func (h *Handler) GetQuestion(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的题目ID")
		return
	}
	question, err := h.svc.Question.GetPublicQuestion(id)
	if err != nil {
		response.NotFound(c, "题目不存在")
		return
	}

	if _, err := h.svc.Question.EnsureUserCanAccessBank(userID, question.BankID); err != nil {
		if errors.Is(err, model.ErrBankNotPurchased) {
			response.Forbidden(c, "请先购买该题库")
			return
		}
		if errors.Is(err, model.ErrBankNotFound) {
			response.NotFound(c, "题库不存在或已下线")
			return
		}
		zap.S().Errorf("GetQuestion access check error: %v", err)
		response.InternalError(c, "获取题目详情失败")
		return
	}
	response.Success(c, question)
}

func (h *Handler) RandomQuestions(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var bankID int64
	if b := c.Query("bank_id"); b != "" {
		v, err := strconv.ParseInt(b, 10, 64)
		if err != nil {
			response.BadRequest(c, "参数错误：无效的题库ID")
			return
		}
		bankID = v
	}
	countStr := c.DefaultQuery("count", "10")
	count, _ := strconv.Atoi(countStr)
	if count < 1 {
		count = 10
	}

	if bankID == 0 {
		response.BadRequest(c, "请选择题库")
		return
	}

	if _, err := h.svc.Question.EnsureUserCanAccessBank(userID, bankID); err != nil {
		if errors.Is(err, model.ErrBankNotPurchased) {
			response.Forbidden(c, "请先购买该题库")
			return
		}
		if errors.Is(err, model.ErrBankNotFound) {
			response.NotFound(c, "题库不存在或已下线")
			return
		}
		zap.S().Errorf("RandomQuestions access check error: %v", err)
		response.InternalError(c, "获取随机题目失败")
		return
	}

	questions, err := h.svc.Question.RandomQuestions(bankID, count)
	if err != nil {
		zap.S().Errorf("RandomQuestions error: %v", err)
		response.InternalError(c, "获取随机题目失败")
		return
	}

	response.Success(c, questions)
}

func (h *Handler) SubmitAnswer(c *gin.Context) {
	var req model.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetInt64("user_id")
	result, err := h.svc.Quiz.SubmitAnswer(userID, req.QuestionID, req.MyAnswer)
	if err != nil {
		zap.S().Errorf("SubmitAnswer error: %v", err)
		h.handleQuizError(c, err, "提交答案失败")
		return
	}

	response.Success(c, result)
}

func (h *Handler) SubmitExam(c *gin.Context) {
	var req model.ExamSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetInt64("user_id")
	result, err := h.svc.Quiz.SubmitExam(userID, req.BankID, req.Answers, req.Duration)
	if err != nil {
		zap.S().Errorf("SubmitExam error: %v", err)
		h.handleQuizError(c, err, "提交考试失败")
		return
	}

	response.Success(c, result)
}

func (h *Handler) handleQuizError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, model.ErrBankNotPurchased), errors.Is(err, model.ErrNoAccessToBank):
		response.Forbidden(c, "无权访问该题目")
	case errors.Is(err, model.ErrBankNotFound):
		response.NotFound(c, "题库不存在或已下线")
	case errors.Is(err, model.ErrQuestionNotFound):
		response.NotFound(c, "题目不存在或已下线")
	default:
		response.InternalError(c, fallback)
	}
}

func (h *Handler) GetQuizHistory(c *gin.Context) {
	userID := c.GetInt64("user_id")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}

	answers, total, err := h.svc.Quiz.GetHistory(userID, page, pageSize)
	if err != nil {
		zap.S().Errorf("GetQuizHistory error: %v", err)
		response.InternalError(c, "获取答题历史失败")
		return
	}

	response.Page(c, answers, total, page, pageSize)
}

func (h *Handler) GetWrongQuestions(c *gin.Context) {
	userID := c.GetInt64("user_id")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}

	wqs, total, err := h.svc.Quiz.GetWrongQuestions(userID, page, pageSize)
	if err != nil {
		zap.S().Errorf("GetWrongQuestions error: %v", err)
		response.InternalError(c, "获取错题本失败")
		return
	}

	response.Page(c, wqs, total, page, pageSize)
}

func (h *Handler) AddToWrongQuestions(c *gin.Context) {
	userID := c.GetInt64("user_id")
	questionID, err := strconv.ParseInt(c.Param("qid"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的问题ID")
		return
	}

	if questionID == 0 {
		response.BadRequest(c, "invalid question id")
		return
	}

	if err := h.svc.Quiz.AddToWrong(userID, questionID); err != nil {
		zap.S().Errorf("AddToWrongQuestions error: %v", err)
		response.InternalError(c, "添加错题失败")
		return
	}

	response.Success(c, nil)
}

func (h *Handler) GetMyScore(c *gin.Context) {
	userID := c.GetInt64("user_id")
	score, err := h.svc.Score.GetMyScore(userID)
	if err != nil {
		response.Success(c, model.Score{
			UserID:        userID,
			TotalScore:    0,
			TotalQuestion: 0,
			CorrectCount:  0,
		})
		return
	}
	response.Success(c, score)
}

func (h *Handler) GetRanking(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}

	ranking, err := h.svc.Score.GetRanking(page, pageSize)
	if err != nil {
		zap.S().Errorf("GetRanking error: %v", err)
		response.InternalError(c, "获取排行榜失败")
		return
	}

	response.Success(c, ranking)
}

func (h *Handler) GetStats(c *gin.Context) {
	userID := c.GetInt64("user_id")
	stats, err := h.svc.Score.GetStats(userID)
	if err != nil {
		zap.S().Errorf("GetStats error: %v", err)
		response.InternalError(c, "获取统计数据失败")
		return
	}
	response.Success(c, stats)
}

func (h *Handler) AdminLogin(c *gin.Context) {
	var req model.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := h.svc.Auth.AdminLogin(&req)
	if err != nil {
		zap.S().Errorf("AdminLogin error: %v", err)
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	response.Success(c, result)
}

func (h *Handler) AdminLogout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			claims, err := middleware.GetJWT().ParseToken(tokenString)
			if err == nil {
				h.svc.Auth.Logout(tokenString, claims.ExpiresAt.Time)
			}
		}
	}
	response.Success(c, nil)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	stats, err := h.svc.Admin.Dashboard()
	if err != nil {
		zap.S().Errorf("GetDashboard error: %v", err)
		response.InternalError(c, "获取仪表盘失败")
		return
	}
	response.Success(c, stats)
}

func (h *Handler) ListAllUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}
	keyword := c.Query("keyword")

	users, total, err := h.svc.Admin.ListUsers(page, pageSize, keyword)
	if err != nil {
		zap.S().Errorf("ListAllUsers error: %v", err)
		response.InternalError(c, "获取用户列表失败")
		return
	}

	response.Page(c, users, total, page, pageSize)
}

func (h *Handler) GetUserDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的用户ID")
		return
	}
	user, err := h.svc.Admin.GetUser(id)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	banks, _ := h.svc.Admin.GetUserPurchasedBanks(id)

	response.Success(c, gin.H{
		"user":  user,
		"banks": banks,
	})
}

func (h *Handler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的用户ID")
		return
	}
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.Admin.UpdateUserStatus(id, req.Status); err != nil {
		zap.S().Errorf("UpdateUserStatus error: %v", err)
		response.InternalError(c, "更新用户状态失败")
		return
	}

	response.Success(c, nil)
}

func (h *Handler) ListQuestionBanks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}

	banks, total, err := h.svc.Question.ListQuestionBanks(page, pageSize)
	if err != nil {
		zap.S().Errorf("ListQuestionBanks error: %v", err)
		response.InternalError(c, "获取题库列表失败")
		return
	}

	response.Page(c, banks, total, page, pageSize)
}

func (h *Handler) ListPublicQuestionBanks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}

	banks, total, err := h.svc.Question.ListPublicQuestionBanks(page, pageSize)
	if err != nil {
		zap.S().Errorf("ListPublicQuestionBanks error: %v", err)
		response.InternalError(c, "获取题库列表失败")
		return
	}

	response.Page(c, banks, total, page, pageSize)
}

func (h *Handler) GetPublicQuestionBankDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的题库ID")
		return
	}

	bank, err := h.svc.Question.GetPublicQuestionBank(id)
	if err != nil {
		if errors.Is(err, model.ErrBankNotFound) {
			response.NotFound(c, "题库不存在或已下线")
			return
		}
		zap.S().Errorf("GetPublicQuestionBankDetail error: %v", err)
		response.InternalError(c, "获取题库详情失败")
		return
	}

	response.Success(c, bank)
}

func (h *Handler) GetQuestionBankAccessStatus(c *gin.Context) {
	userID := c.GetInt64("user_id")
	bankID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的题库ID")
		return
	}

	status, err := h.svc.Question.GetBankAccessStatus(userID, bankID)
	if err != nil {
		if errors.Is(err, model.ErrBankNotFound) {
			response.NotFound(c, "题库不存在或已下线")
			return
		}
		zap.S().Errorf("GetQuestionBankAccessStatus error: %v", err)
		response.InternalError(c, "获取题库权限失败")
		return
	}

	response.Success(c, status)
}

func (h *Handler) CreateQuestionBank(c *gin.Context) {
	var bank model.QuestionBank
	if err := c.ShouldBindJSON(&bank); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.Question.CreateQuestionBank(&bank); err != nil {
		zap.S().Errorf("CreateQuestionBank error: %v", err)
		response.InternalError(c, "创建题库失败")
		return
	}

	response.Success(c, bank)
}

func (h *Handler) UpdateQuestionBank(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的题库ID")
		return
	}
	var bank model.QuestionBank
	if err := c.ShouldBindJSON(&bank); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	bank.ID = id

	if err := h.svc.Question.UpdateQuestionBank(&bank); err != nil {
		zap.S().Errorf("UpdateQuestionBank error: %v", err)
		response.InternalError(c, "更新题库失败")
		return
	}

	response.Success(c, bank)
}

func (h *Handler) DeleteQuestionBank(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的题库ID")
		return
	}
	if err := h.svc.Question.DeleteQuestionBank(id); err != nil {
		zap.S().Errorf("DeleteQuestionBank error: %v", err)
		response.InternalError(c, "删除题库失败")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) ListAllQuestions(c *gin.Context) {
	var bankID int64
	if b := c.Query("bank_id"); b != "" {
		v, err := strconv.ParseInt(b, 10, 64)
		if err != nil {
			response.BadRequest(c, "参数错误：无效的题库ID")
			return
		}
		bankID = v
	}
	questionType := c.Query("type")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}

	questions, total, err := h.svc.Question.ListQuestions(bankID, questionType, page, pageSize)
	if err != nil {
		zap.S().Errorf("ListAllQuestions error: %v", err)
		response.InternalError(c, "获取题目列表失败")
		return
	}

	response.Page(c, questions, total, page, pageSize)
}

func (h *Handler) CreateQuestion(c *gin.Context) {
	var req struct {
		model.Question
		Options []model.QuestionOption `json:"options"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	req.Question.Options = req.Options
	if err := h.svc.Question.CreateQuestion(&req.Question); err != nil {
		zap.S().Errorf("CreateQuestion error: %v", err)
		response.InternalError(c, "创建题目失败")
		return
	}

	response.Success(c, req.Question)
}

func (h *Handler) UpdateQuestion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的题目ID")
		return
	}
	var req struct {
		model.Question
		Options []model.QuestionOption `json:"options"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	req.Question.ID = id
	req.Question.Options = req.Options

	if err := h.svc.Question.UpdateQuestion(&req.Question); err != nil {
		zap.S().Errorf("UpdateQuestion error: %v", err)
		response.InternalError(c, "更新题目失败")
		return
	}

	response.Success(c, req.Question)
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的题目ID")
		return
	}
	if err := h.svc.Question.DeleteQuestion(id); err != nil {
		zap.S().Errorf("DeleteQuestion error: %v", err)
		response.InternalError(c, "删除题目失败")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) ImportQuestions(c *gin.Context) {
	bankIDStr := c.PostForm("bank_id")
	if bankIDStr == "" {
		response.BadRequest(c, "请选择题库")
		return
	}
	bankID, err := strconv.ParseInt(bankIDStr, 10, 64)
	if err != nil || bankID == 0 {
		response.BadRequest(c, "无效的题库ID")
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请上传文件")
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		response.BadRequest(c, "文件格式错误，请上传Excel文件(.xlsx)")
		return
	}
	defer f.Close()

	rows, err := f.GetRows("题目导入模板")
	if err != nil {
		rows, err = f.GetRows("Sheet1")
		if err != nil {
			response.BadRequest(c, "读取Excel失败")
			return
		}
	}

	if len(rows) < 2 {
		response.BadRequest(c, "Excel中没有数据")
		return
	}

	imported := 0
	failed := 0
	skipFirst := true
	rowNum := 0

	for _, row := range rows {
		rowNum++
		if skipFirst || len(row) < 3 {
			skipFirst = false
			continue
		}

		question := &model.Question{
			BankID: bankID,
		}

		if len(row) > 0 {
			question.Content = row[0]
		}
		if len(row) > 1 {
			question.Answer = row[1]
		}
		if len(row) > 2 {
			question.Analysis = row[2]
		}
		if len(row) > 3 {
			qType := row[3]
			switch qType {
			case "单选":
				qType = "single"
			case "多选":
				qType = "multiple"
			case "判断":
				qType = "truefalse"
			case "":
				qType = "single"
			}
			question.Type = qType
		} else {
			question.Type = "single"
		}
		if len(row) > 4 && row[4] != "" {
			diff, err := strconv.Atoi(row[4])
			if err != nil || diff < 1 || diff > 5 {
				diff = 3
			}
			question.Difficulty = int8(diff)
		} else {
			question.Difficulty = 3
		}

		if len(row) > 5 && row[5] != "" {
			question.Options = append(question.Options, model.QuestionOption{OptionKey: "A", OptionValue: row[5]})
		}
		if len(row) > 6 && row[6] != "" {
			question.Options = append(question.Options, model.QuestionOption{OptionKey: "B", OptionValue: row[6]})
		}
		if len(row) > 7 && row[7] != "" {
			question.Options = append(question.Options, model.QuestionOption{OptionKey: "C", OptionValue: row[7]})
		}
		if len(row) > 8 && row[8] != "" {
			question.Options = append(question.Options, model.QuestionOption{OptionKey: "D", OptionValue: row[8]})
		}

		if question.Content == "" {
			continue
		}

		if err := h.svc.Question.CreateQuestion(question); err != nil {
			zap.S().Warnf("Import question failed at row %d: %v", rowNum, err)
			failed++
		} else {
			imported++
		}
	}

	response.Success(c, gin.H{"imported": imported, "failed": failed})
}

func (h *Handler) DownloadQuestionTemplate(c *gin.Context) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "题目导入模板"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	headers := []string{"题目内容", "正确答案", "解析", "类型", "难度", "选项A", "选项B", "选项C", "选项D"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	f.SetColWidth(sheetName, "A", "A", 50)
	f.SetColWidth(sheetName, "B", "B", 15)
	f.SetColWidth(sheetName, "C", "C", 30)
	f.SetColWidth(sheetName, "D", "D", 15)
	f.SetColWidth(sheetName, "E", "E", 10)
	f.SetColWidth(sheetName, "F", "F", 30)
	f.SetColWidth(sheetName, "G", "G", 30)
	f.SetColWidth(sheetName, "H", "H", 30)
	f.SetColWidth(sheetName, "I", "I", 30)

	examples := [][]interface{}{
		{"以下哪个是华为的操作系统？", "A", "鸿蒙系统是华为自主研发的操作系统", "单选", 2, "鸿蒙", "iOS", "Android", ""},
		{"以下哪些是编程语言？", "AB", "Java和Python都是高级编程语言", "多选", 3, "Java", "Python", "Windows", "Excel"},
		{"Java是一种编程语言。", "正确", "判断题填写正确或错误", "判断", 1, "", "", "", ""},
	}

	for rowIdx, row := range examples {
		for colIdx, val := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheetName, cell, val)
		}
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=question_template.xlsx")
	c.Header("File-Name", "question_template.xlsx")

	if err := f.Write(c.Writer); err != nil {
		response.InternalError(c, "生成模板失败")
		return
	}
}

func (h *Handler) GetStatsOverview(c *gin.Context) {
	stats, err := h.svc.Admin.Dashboard()
	if err != nil {
		zap.S().Errorf("GetStatsOverview error: %v", err)
		response.InternalError(c, "获取统计概览失败")
		return
	}
	response.Success(c, stats)
}

func (h *Handler) GetUserStats(c *gin.Context) {
	stats, err := h.svc.Admin.GetUserStats()
	if err != nil {
		zap.S().Errorf("GetUserStats error: %v", err)
		response.InternalError(c, "获取用户统计失败")
		return
	}
	response.Success(c, stats)
}

func (h *Handler) GetQuestionStats(c *gin.Context) {
	stats, err := h.svc.Admin.GetQuestionStats()
	if err != nil {
		zap.S().Errorf("GetQuestionStats error: %v", err)
		response.InternalError(c, "获取题目统计失败")
		return
	}
	response.Success(c, stats)
}

func (h *Handler) auditLog(c *gin.Context, action, target string, targetID int64, detail string) {
	adminID := c.GetInt64("user_id")
	adminName, _ := c.Get("openid")
	adminNameStr, _ := adminName.(string)

	// Persist to database
	if err := h.svc.AuditLog.CreateAuditLog(adminID, adminNameStr, action, target, targetID, detail, c.ClientIP()); err != nil {
		zap.S().Errorf("persist audit log failed: %v", err)
	}

	zap.S().Infow("admin audit log",
		"admin_id", adminID,
		"action", action,
		"target", target,
		"target_id", targetID,
		"detail", detail,
		"ip", c.ClientIP(),
	)
}

func (h *Handler) ListAuditLogs(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}
	adminName := c.Query("admin_name")
	action := c.Query("action")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	var startPtr, endPtr *string
	if startTime != "" {
		startPtr = &startTime
	}
	if endTime != "" {
		endPtr = &endTime
	}

	logs, total, err := h.svc.AuditLog.ListAuditLogs(page, pageSize, adminName, action, startPtr, endPtr)
	if err != nil {
		zap.S().Errorf("ListAuditLogs error: %v", err)
		response.InternalError(c, "获取审计日志失败")
		return
	}
	response.Page(c, logs, total, page, pageSize)
}

func (h *Handler) GetBankStatsList(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	stats, total, err := h.svc.BankStats.GetBankStatsList(page, pageSize, sortBy, sortOrder)
	if err != nil {
		zap.S().Errorf("GetBankStatsList error: %v", err)
		response.InternalError(c, "获取题库统计失败")
		return
	}
	response.Page(c, stats, total, page, pageSize)
}

func (h *Handler) GetBankStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的题库ID")
		return
	}

	stat, err := h.svc.BankStats.GetBankStats(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "题库不存在")
			return
		}
		zap.S().Errorf("GetBankStats error: %v", err)
		response.InternalError(c, "获取题库统计详情失败")
		return
	}
	response.Success(c, stat)
}

func (h *Handler) ListFeedbacks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}

	status := int8(-1)
	if s := c.Query("status"); s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			status = int8(v)
		}
	}
	feedbackType := c.Query("type")

	feedbacks, total, err := h.svc.Feedback.ListFeedbacks(page, pageSize, status, feedbackType)
	if err != nil {
		zap.S().Errorf("ListFeedbacks error: %v", err)
		response.InternalError(c, "获取反馈列表失败")
		return
	}
	response.Page(c, feedbacks, total, page, pageSize)
}

func (h *Handler) UpdateFeedbackStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的反馈ID")
		return
	}
	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.Feedback.UpdateFeedbackStatus(id, req.Status); err != nil {
		zap.S().Errorf("UpdateFeedbackStatus error: %v", err)
		response.InternalError(c, "更新反馈状态失败")
		return
	}
	response.Success(c, nil)
}

func (h *Handler) CreateFeedback(c *gin.Context) {
	var req struct {
		Type    string `json:"type"`
		Content string `json:"content" binding:"required"`
		Contact string `json:"contact"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetInt64("user_id")
	feedback, err := h.svc.Feedback.CreateFeedback(userID, req.Type, req.Content, req.Contact)
	if err != nil {
		zap.S().Errorf("CreateFeedback error: %v", err)
		response.InternalError(c, "提交反馈失败")
		return
	}
	response.Success(c, feedback)
}

func (h *Handler) GetExamRecords(c *gin.Context) {
	userID := c.GetInt64("user_id")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 20
	}

	records, total, err := h.svc.Quiz.GetExamRecords(userID, page, pageSize)
	if err != nil {
		zap.S().Errorf("GetExamRecords error: %v", err)
		response.InternalError(c, "获取考试记录失败")
		return
	}
	response.Page(c, records, total, page, pageSize)
}

func (h *Handler) GetMyPurchasedBanks(c *gin.Context) {
	userID := c.GetInt64("user_id")
	banks, err := h.svc.Quiz.GetMyPurchasedBanks(userID)
	if err != nil {
		zap.S().Errorf("GetMyPurchasedBanks error: %v", err)
		response.InternalError(c, "获取已购题库失败")
		return
	}
	response.Success(c, banks)
}
