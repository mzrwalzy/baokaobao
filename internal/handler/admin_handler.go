package handler

import "github.com/gin-gonic/gin"

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) Login(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) Logout(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) ListQuestionBanks(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) CreateQuestionBank(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) UpdateQuestionBank(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) DeleteQuestionBank(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) ListQuestions(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) CreateQuestion(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) UpdateQuestion(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) DeleteQuestion(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) ImportQuestions(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) StatsOverview(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) StatsUsers(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AdminHandler) StatsQuestions(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}
