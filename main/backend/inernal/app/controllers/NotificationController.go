package controllers

import (
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// NotificationController 结构体
type NotificationController struct {
	notificationService services.NotificationService
}

// NewNotificationController 创建新的 NotificationController 实例
func NewNotificationController(notificationService services.NotificationService) *NotificationController {
	return &NotificationController{notificationService: notificationService}
}

// CreateNotification (可选) - 公开接口，用于测试或管理员操作，
// 实际上生产环境中，通知多由业务逻辑内部调用
func (ctl *NotificationController) CreateNotification(c *gin.Context) {
	var input struct {
		UserID      uuid.UUID      `json:"user_id" binding:"required"`
		ActorID     *uuid.UUID     `json:"actor_id"`
		Type        string         `json:"type" binding:"required"`
		Title       string         `json:"title" binding:"required"`
		Description string         `json:"description" binding:"required"`
		RelatedID   *uuid.UUID     `json:"related_id"`
		RelatedType string         `json:"related_type"`
		Metadata    map[string]any `json:"metadata"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的传输模型", "details": err.Error()})
		return
	}

	notification := &tables.Notification{
		UserID:      input.UserID,
		ActorID:     input.ActorID,
		Type:        input.Type,
		Title:       input.Title,
		Description: input.Description,
		RelatedID:   input.RelatedID,
		RelatedType: input.RelatedType,
		Metadata:    input.Metadata,
	}

	if err := ctl.notificationService.CreateNotification(notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "通知创建成功"})
}

// GetUserNotifications 获取用户的所有通知
func (ctl *NotificationController) GetUserNotifications(c *gin.Context) {
	// 获取当前登录用户的ID从上下文
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 断言 userID 为 uuid.UUID 类型
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID类型"})
		return
	}

	notifications, err := ctl.notificationService.GetUserNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取通知失败"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkNotificationAsRead 标记通知为已读
func (ctl *NotificationController) MarkNotificationAsRead(c *gin.Context) {
	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的通知ID"})
		return
	}

	if err := ctl.notificationService.MarkNotificationAsRead(notificationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "将通知标记为已读失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "通知已被标记为已读"})
}

// DeleteNotification 删除通知
func (ctl *NotificationController) DeleteNotification(c *gin.Context) {
	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的通知ID"})
		return
	}

	if err := ctl.notificationService.DeleteNotification(notificationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "通知删除成功"})
}
