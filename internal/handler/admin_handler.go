package handler

import (
	"errors"
	"fmt"
	"strconv"

	"baokaobao/internal/model"
	"baokaobao/internal/pkg/response"
	"baokaobao/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (h *Handler) ListAdminUsers(c *gin.Context) {
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

	admins, total, err := h.svc.Admin.ListAdminUsers(page, pageSize)
	if err != nil {
		zap.S().Errorf("ListAdminUsers error: %v", err)
		response.InternalError(c, "获取管理员列表失败")
		return
	}
	response.Page(c, admins, total, page, pageSize)
}

func (h *Handler) GetAdminUserDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的管理员ID")
		return
	}
	admin, err := h.svc.Admin.GetAdminUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "管理员不存在")
			return
		}
		zap.S().Errorf("GetAdminUserDetail error: %v", err)
		response.InternalError(c, "获取管理员详情失败")
		return
	}
	response.Success(c, admin)
}

func (h *Handler) CreateAdminUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
		Role     string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if req.Role != model.RoleSuperAdmin && req.Role != model.RoleOperator {
		response.BadRequest(c, "无效的角色")
		return
	}

	if err := utils.ValidateAdminPassword(req.Password); err != nil {
		response.BadRequest(c, "密码强度不足: "+err.Error())
		return
	}

	if err := h.svc.Admin.CreateAdminUserWithRole(req.Username, req.Password, req.Nickname, req.Role); err != nil {
		zap.S().Errorf("CreateAdminUser error: %v", err)
		response.InternalError(c, "创建管理员失败")
		return
	}

	h.auditLog(c, "create", "admin_user", 0, fmt.Sprintf("创建管理员 %s role=%s", req.Username, req.Role))
	response.Success(c, nil)
}

func (h *Handler) UpdateAdminUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的管理员ID")
		return
	}
	var req struct {
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
		Status   int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if req.Role != "" && req.Role != model.RoleSuperAdmin && req.Role != model.RoleOperator {
		response.BadRequest(c, "无效的角色")
		return
	}

	admin := &model.AdminUser{
		ID:       id,
		Nickname: req.Nickname,
		Role:     req.Role,
		Status:   req.Status,
	}
	if err := h.svc.Admin.UpdateAdminUser(admin); err != nil {
		zap.S().Errorf("UpdateAdminUser error: %v", err)
		response.InternalError(c, "更新管理员失败")
		return
	}

	h.auditLog(c, "update", "admin_user", id, fmt.Sprintf("更新管理员 ID=%d", id))
	response.Success(c, nil)
}

func (h *Handler) ResetAdminPassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的管理员ID")
		return
	}
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := utils.ValidateAdminPassword(req.Password); err != nil {
		response.BadRequest(c, "密码强度不足: "+err.Error())
		return
	}

	if err := h.svc.Admin.ResetAdminPassword(id, req.Password); err != nil {
		zap.S().Errorf("ResetAdminPassword error: %v", err)
		response.InternalError(c, "重置密码失败")
		return
	}

	h.auditLog(c, "reset_password", "admin_user", id, fmt.Sprintf("重置管理员密码 ID=%d", id))
	response.Success(c, nil)
}

func (h *Handler) DeleteAdminUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的管理员ID")
		return
	}

	currentAdminID := c.GetInt64("user_id")
	if id == currentAdminID {
		response.BadRequest(c, "不能删除自己")
		return
	}

	if err := h.svc.Admin.DeleteAdminUser(id); err != nil {
		zap.S().Errorf("DeleteAdminUser error: %v", err)
		response.InternalError(c, "删除管理员失败")
		return
	}

	h.auditLog(c, "delete", "admin_user", id, fmt.Sprintf("删除管理员 ID=%d", id))
	response.Success(c, nil)
}

func (h *Handler) ListAdminBankPermissions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的管理员ID")
		return
	}
	perms, err := h.svc.Admin.ListAdminBankPermissions(id)
	if err != nil {
		zap.S().Errorf("ListAdminBankPermissions error: %v", err)
		response.InternalError(c, "获取题库权限失败")
		return
	}
	response.Success(c, perms)
}

func (h *Handler) GrantAdminBankAccess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的管理员ID")
		return
	}
	var req struct {
		BankID int64 `json:"bank_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if _, err := h.svc.Question.GetQuestionBank(req.BankID); err != nil {
		if errors.Is(err, model.ErrBankNotFound) {
			response.NotFound(c, "题库不存在")
			return
		}
		zap.S().Errorf("GrantAdminBankAccess error: %v", err)
		response.InternalError(c, "开通权限失败")
		return
	}

	if err := h.svc.Admin.GrantAdminBankAccess(id, req.BankID); err != nil {
		zap.S().Errorf("GrantAdminBankAccess error: %v", err)
		response.InternalError(c, "开通权限失败")
		return
	}

	h.auditLog(c, "grant", "admin_bank_access", id, fmt.Sprintf("授予管理员题库权限 adminID=%d bankID=%d", id, req.BankID))
	response.Success(c, nil)
}

func (h *Handler) RevokeAdminBankAccess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误：无效的管理员ID")
		return
	}
	bankID, err := strconv.ParseInt(c.Param("bankId"), 10, 64)
	if err != nil || bankID <= 0 {
		response.BadRequest(c, "参数错误：无效的题库ID")
		return
	}

	if err := h.svc.Admin.RevokeAdminBankAccess(id, bankID); err != nil {
		zap.S().Errorf("RevokeAdminBankAccess error: %v", err)
		response.InternalError(c, "取消权限失败")
		return
	}

	h.auditLog(c, "revoke", "admin_bank_access", id, fmt.Sprintf("取消管理员题库权限 adminID=%d bankID=%d", id, bankID))
	response.Success(c, nil)
}
