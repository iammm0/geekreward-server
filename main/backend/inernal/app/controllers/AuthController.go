package controllers

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// AuthController 结构体
type AuthController struct {
	authService         services.AuthService
	notificationService services.NotificationService
}

// NewAuthController 创建新的 AuthController 实例
func NewAuthController(
	authService services.AuthService,
	notificationService services.NotificationService,
) *AuthController {
	return &AuthController{
		authService:         authService,
		notificationService: notificationService,
	}
}

// Register 用户注册处理函数
func (ctl *AuthController) Register(c *gin.Context) {
	var input dtos.RegisterInput

	// 1. 绑定表单数据(含文本字段)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的传输模型", "details": err.Error()})
		return
	}

	// 2. 获取头像文件 (可选)
	file, err := c.FormFile("profilePicture")
	if err != nil {
		if err.Error() != "http: no such file" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "头像文件上传出错", "details": err.Error()})
			return
		}
		// 如果确实没上传文件 => file = nil
		file = nil
	}

	// 3. 如果上传了头像, 效仿 attachment 做随机命名 & 保存
	var avatarURL string
	if file != nil {
		uploadDir := "./uploads/avatars"
		// 确保目录存在
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建上传目录"})
			return
		}
		ext := filepath.Ext(file.Filename)
		newFilename := uuid.New().String() + ext
		filePath := filepath.Join(uploadDir, newFilename)

		// 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "头像保存失败", "details": err.Error()})
			return
		}

		// 生成可访问URL, 例如 "/uploads/avatars/xxx.jpg"
		avatarURL = "/uploads/avatars/" + newFilename
	}

	// 4. 将头像URL放入 input.ProfilePicture
	input.ProfilePicture = avatarURL

	// 5. 调用 Service 层
	user, serviceErr := ctl.authService.Register(input)
	if serviceErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": serviceErr.Error()})
		return
	}

	// 创建通用通知
	notification := &tables.Notification{
		UserID:      user.ID,
		Type:        "Welcome",
		Title:       "欢迎加入 GeekReward!",
		Description: "感谢您注册 GeekReward 平台，祝您使用愉快！",
		IsRead:      false,
	}

	if err := ctl.notificationService.CreateNotification(notification); err != nil {
		// 记录错误但不影响主流程
		log.Println("发送欢迎通知失败:", err)
	}

	// 6. 返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"message":         "用户注册成功",
		"user_id":         user.ID.String(),
		"profile_picture": user.ProfilePicture,
	})
}

// Login 用户登录处理函数
func (ctl *AuthController) Login(c *gin.Context) {
	var input dtos.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "存在无效的 input 字段", "details": err.Error()})
		return
	}

	token, user, err := ctl.authService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证"})
		return
	}

	// 创建通用通知
	notification := &tables.Notification{
		UserID:      user.ID,
		Type:        "Welcome",
		Title:       "欢迎登录 GeekReward!",
		Description: "您已成功登录 GeekReward！",
		IsRead:      false,
	}

	if err := ctl.notificationService.CreateNotification(notification); err != nil {
		// 记录错误但不影响主流程
		log.Println("发送登录欢迎通知失败:", err)
	}

	// 返回 token & user 信息给前端
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
