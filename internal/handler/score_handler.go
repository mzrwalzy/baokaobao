package handler

import "github.com/gin-gonic/gin"

type ScoreHandler struct{}

func NewScoreHandler() *ScoreHandler {
	return &ScoreHandler{}
}

func (h *ScoreHandler) MyScore(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *ScoreHandler) Ranking(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *ScoreHandler) Stats(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}
