package handler

import "github.com/gin-gonic/gin"

type QuizHandler struct{}

func NewQuizHandler() *QuizHandler {
	return &QuizHandler{}
}

func (h *QuizHandler) Submit(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *QuizHandler) History(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *QuizHandler) WrongQuestions(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *QuizHandler) AddWrong(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}
