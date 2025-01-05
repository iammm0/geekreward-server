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
	userService services.UserService
}

// NewUserController 创建新的 UserController 实例
func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetUserInfo 获取用户信息
func (ctl *UserController) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	user, err := ctl.userService.GetUserByID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserInfo 更新用户信息
func (ctl *UserController) UpdateUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var input dtos.UpdateUserProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	user, err := ctl.userService.UpdateUser(uid, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}
