package handler

import "github.com/gin-gonic/gin"

type QuestionHandler struct{}

func NewQuestionHandler() *QuestionHandler {
	return &QuestionHandler{}
}

func (h *QuestionHandler) List(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *QuestionHandler) Get(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *QuestionHandler) Random(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}
