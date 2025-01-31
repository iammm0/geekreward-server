package controllers

import (
	"GeekReward/inernal/app/services"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AttachmentController 处理附件相关请求
type AttachmentController struct {
	notificationService services.NotificationService
}

// NewAttachmentController 创建新的 AttachmentController
func NewAttachmentController(
	notificationService services.NotificationService,
) *AttachmentController {
	return &AttachmentController{
		notificationService: notificationService,
	}
}

// UploadAttachment 处理单个文件上传的示例
func (ctl *AttachmentController) UploadAttachment(c *gin.Context) {
	// 从表单中获取文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file is received"})
		return
	}

	// 可选：指定存储目录
	uploadDir := "./uploads"
	// 确保目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	// 使用UUID重命名文件，避免重名
	ext := filepath.Ext(file.Filename)                // 如 .jpg, .png
	newFilename := uuid.New().String() + ext          // 形如 123e4567-e89b-12d3-a456-426655440000.png
	filePath := filepath.Join(uploadDir, newFilename) // ./uploads/xxx.png

	// 保存到本地
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// 生成可访问的URL (视情况添加域名/ip)
	// 假设后端已配置: router.Static("/uploads", "./uploads") => /uploads 可以访问本地 ./uploads
	fileURL := fmt.Sprintf("/uploads/%s", newFilename)

	c.JSON(http.StatusOK, gin.H{
		"message":  "upload success",
		"file_url": fileURL,
	})
}
