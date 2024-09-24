package controllers

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthController 结构体
type AuthController struct {
	authService services.AuthService
}

// NewAuthController 创建新的 AuthController 实例
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register 用户注册处理函数
func (ctl *AuthController) Register(c *gin.Context) {
	var input dtos.RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// 获取文件上传
	file, _ := c.FormFile("profilePicture")

	// 注册用户
	user, err := ctl.authService.Register(input, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user": user})
}

// Login 用户登录处理函数
func (ctl *AuthController) Login(c *gin.Context) {
	var input dtos.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	token, err := ctl.authService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
