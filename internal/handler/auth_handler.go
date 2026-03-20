package handler

import "github.com/gin-gonic/gin"

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) LoginByWechat(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AuthHandler) DecryptPhone(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}
