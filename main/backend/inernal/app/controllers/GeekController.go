package controllers

import (
	"GeekReward/inernal/app/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// GeekController 结构体
type GeekController struct {
	geekService         services.GeekService
	notificationService services.NotificationService
}

// NewGeekController 创建新的 GeekController 实例
func NewGeekController(
	geekService services.GeekService,
	notificationService services.NotificationService,
) *GeekController {
	return &GeekController{
		geekService:         geekService,
		notificationService: notificationService,
	}
}

// GetTopGeeks 获取排名前的极客用户
func (ctl *GeekController) GetTopGeeks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 limit"})
		return
	}

	geeks, err := ctl.geekService.GetTopGeeks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取极客失败"})
		return
	}

	c.JSON(http.StatusOK, geeks)
}

// GetGeekByID 获取指定ID的极客用户
func (ctl *GeekController) GetGeekByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	geek, err := ctl.geekService.GetGeekByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "极客未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取极客失败"})
		}
		return
	}

	c.JSON(http.StatusOK, geek)
}

// ExpressAffection 向特定极客或团队表达好感
func (ctl *GeekController) ExpressAffection(c *gin.Context) {
	geekIDStr := c.Param("id")
	geekID, err := uuid.Parse(geekIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的极客ID"})
		return
	}

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

	if err := ctl.geekService.ExpressAffection(geekID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已向该极客表达好感"})
}

// SendInvitation 向特定极客发出组队邀请
func (ctl *GeekController) SendInvitation(c *gin.Context) {
	// 获取邀请的极客ID从URL参数
	geekIDStr := c.Param("id")
	geekID, err := uuid.Parse(geekIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的极客ID"})
		return
	}

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

	// 创建发送邀请的服务调用
	err = ctl.geekService.SendInvitation(geekID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "邀请发送成功"})
}
