package handler

import "github.com/gin-gonic/gin"

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "todo"})
}
