package controllers

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
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

	user, err := ctl.userService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
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

	var input dtos.UpdateUserProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	user, err := ctl.userService.UpdateUser(userID.(uint), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
		return
	}

	c.JSON(http.StatusOK, user)
}
