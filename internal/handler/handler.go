package handler

import (
	"baokaobao/internal/model"
	"baokaobao/internal/pkg/response"
	"baokaobao/internal/pkg/wechat"
	"baokaobao/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := h.svc.Auth.LoginByWechat(req.Code)
	if err != nil {
		zap.S().Errorf("LoginByWechat error: %v", err)
		response.InternalError(c, err.Error())
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
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, model.DecryptPhoneResponse{Phone: phone})
}

func (h *Handler) Logout(c *gin.Context) {
	response.Success(c, nil)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	profile, err := h.svc.User.GetProfile(userID)
	if err != nil {
		response.InternalError(c, err.Error())
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
		response.InternalError(c, err.Error())
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

	if header.Size > 2*1024*1024 {
		response.BadRequest(c, "file size must be less than 2MB")
		return
	}

	userID := c.GetInt64("user_id")

	ext := ".jpg"
	if header.Header.Get("Content-Type") == "image/png" {
		ext = ".png"
	}

	filename := fmt.Sprintf("avatar_%d%s", userID, ext)

	url, err := h.svc.User.UploadAvatar(userID, file, filename)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"url": url})
}

func (h *Handler) ListQuestions(c *gin.Context) {
	bankID := int64(0)
	if b := c.Query("bank_id"); b != "" {
		bankID = parseInt64(b)
	}
	questionType := c.Query("type")
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	questions, total, err := h.svc.Question.ListQuestions(bankID, questionType, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Page(c, questions, total, page, pageSize)
}

func (h *Handler) GetQuestion(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	question, err := h.svc.Question.GetQuestion(id)
	if err != nil {
		response.NotFound(c, "题目不存在")
		return
	}
	response.Success(c, question)
}

func (h *Handler) RandomQuestions(c *gin.Context) {
	bankID := int64(0)
	if b := c.Query("bank_id"); b != "" {
		bankID = parseInt64(b)
	}
	count := parseIntDefault(c.Query("count"), 10)

	questions, err := h.svc.Question.RandomQuestions(bankID, count)
	if err != nil {
		response.InternalError(c, err.Error())
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
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *Handler) GetQuizHistory(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	answers, total, err := h.svc.Quiz.GetHistory(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Page(c, answers, total, page, pageSize)
}

func (h *Handler) GetWrongQuestions(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	wqs, total, err := h.svc.Quiz.GetWrongQuestions(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Page(c, wqs, total, page, pageSize)
}

func (h *Handler) AddToWrongQuestions(c *gin.Context) {
	userID := c.GetInt64("user_id")
	questionID := parseInt64(c.Param("qid"))

	if questionID == 0 {
		response.BadRequest(c, "invalid question id")
		return
	}

	if err := h.svc.Quiz.AddToWrong(userID, questionID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *Handler) GetMyScore(c *gin.Context) {
	userID := c.GetInt64("user_id")
	score, err := h.svc.Score.GetMyScore(userID)
	if err != nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, score)
}

func (h *Handler) GetRanking(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	ranking, err := h.svc.Score.GetRanking(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, ranking)
}

func (h *Handler) GetStats(c *gin.Context) {
	userID := c.GetInt64("user_id")
	stats, err := h.svc.Score.GetStats(userID)
	if err != nil {
		response.InternalError(c, err.Error())
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
	response.Success(c, nil)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	stats, err := h.svc.Admin.Dashboard()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

func (h *Handler) ListAllUsers(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	users, total, err := h.svc.Admin.ListUsers(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Page(c, users, total, page, pageSize)
}

func (h *Handler) GetUserDetail(c *gin.Context) {
	id := parseInt64(c.Param("id"))
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
	id := parseInt64(c.Param("id"))
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.Admin.UpdateUserStatus(id, req.Status); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *Handler) ListQuestionBanks(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	banks, total, err := h.svc.Question.ListQuestionBanks(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Page(c, banks, total, page, pageSize)
}

func (h *Handler) CreateQuestionBank(c *gin.Context) {
	var bank model.QuestionBank
	if err := c.ShouldBindJSON(&bank); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.Question.CreateQuestionBank(&bank); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, bank)
}

func (h *Handler) UpdateQuestionBank(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	var bank model.QuestionBank
	if err := c.ShouldBindJSON(&bank); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	bank.ID = id

	if err := h.svc.Question.UpdateQuestionBank(&bank); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, bank)
}

func (h *Handler) DeleteQuestionBank(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	if err := h.svc.Question.DeleteQuestionBank(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *Handler) ListAllQuestions(c *gin.Context) {
	bankID := int64(0)
	if b := c.Query("bank_id"); b != "" {
		bankID = parseInt64(b)
	}
	questionType := c.Query("type")
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)

	questions, total, err := h.svc.Question.ListQuestions(bankID, questionType, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
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
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, req.Question)
}

func (h *Handler) UpdateQuestion(c *gin.Context) {
	id := parseInt64(c.Param("id"))
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
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, req.Question)
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	if err := h.svc.Question.DeleteQuestion(id); err != nil {
		response.InternalError(c, err.Error())
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
	bankID := parseInt64(bankIDStr)
	if bankID == 0 {
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
	skipFirst := true

	for _, row := range rows {
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
			if qType == "" {
				qType = "single"
			}
			question.Type = qType
		} else {
			question.Type = "single"
		}
		if len(row) > 4 {
			diff := parseIntDefault(row[4], 3)
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

		if err := h.svc.Question.CreateQuestion(question); err == nil {
			imported++
		}
	}

	response.Success(c, gin.H{"imported": imported})
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
		{"以下哪个是华为的操作系统？", "A", "鸿蒙系统是华为自主研发的操作系统", "single", 2, "鸿蒙", "iOS", "Android", ""},
		{"以下哪些是编程语言？", "AB", "Java和Python都是高级编程语言", "multiple", 3, "Java", "Python", "Windows", "Excel"},
		{"Java是一种编程语言。", "true", "判断题直接填true或false", "truefalse", 1, "", "", "", ""},
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
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

func (h *Handler) GetUserStats(c *gin.Context) {
	stats, err := h.svc.Admin.GetUserStats()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

func (h *Handler) GetQuestionStats(c *gin.Context) {
	stats, err := h.svc.Admin.GetQuestionStats()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, stats)
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
