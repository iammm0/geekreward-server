package controllers

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

// UserController 结构体
type UserController struct {
	userService         services.UserService
	notificationService services.NotificationService
}

// NewUserController 创建新的 UserController 实例
func NewUserController(
	userService services.UserService,
	notificationService services.NotificationService,
) *UserController {
	return &UserController{
		userService:         userService,
		notificationService: notificationService,
	}
}

// GetUserInfo 获取用户信息
func (ctl *UserController) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	user, err := ctl.userService.GetUserByID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserInfo 更新用户信息
func (ctl *UserController) UpdateUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	var input dtos.UpdateUserProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的传输模型", "details": err.Error()})
		return
	}

	user, err := ctl.userService.UpdateUser(uid, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息成功"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}
