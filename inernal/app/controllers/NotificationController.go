package controllers

import (
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// NotificationController 结构体
type NotificationController struct {
	notificationService services.NotificationService
}

// NewNotificationController 创建新的 NotificationController 实例
func NewNotificationController(notificationService services.NotificationService) *NotificationController {
	return &NotificationController{notificationService: notificationService}
}

// CreateNotification 创建新的通知
func (ctl *NotificationController) CreateNotification(c *gin.Context) {
	var input struct {
		UserID      uint   `json:"user_id" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	notification := &tables.Notification{
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
	}

	if err := ctl.notificationService.CreateNotification(notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification created successfully"})
}

// GetUserNotifications 获取用户的所有通知
func (ctl *NotificationController) GetUserNotifications(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	notifications, err := ctl.notificationService.GetUserNotifications(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkNotificationAsRead 标记通知为已读
func (ctl *NotificationController) MarkNotificationAsRead(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	if err := ctl.notificationService.MarkNotificationAsRead(uint(notificationID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// DeleteNotification 删除通知
func (ctl *NotificationController) DeleteNotification(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	if err := ctl.notificationService.DeleteNotification(uint(notificationID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}
