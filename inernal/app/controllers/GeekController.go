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
	geekService services.GeekService
}

// NewGeekController 创建新的 GeekController 实例
func NewGeekController(geekService services.GeekService) *GeekController {
	return &GeekController{geekService: geekService}
}

// GetTopGeeks 获取排名前的极客用户
func (ctl *GeekController) GetTopGeeks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	geeks, err := ctl.geekService.GetTopGeeks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch geeks"})
		return
	}

	c.JSON(http.StatusOK, geeks)
}

// GetGeekByID 获取指定ID的极客用户
func (ctl *GeekController) GetGeekByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	geek, err := ctl.geekService.GetGeekByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Geek not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch geek"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid geek ID"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := ctl.geekService.ExpressAffection(geekID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Affection expressed successfully"})
}

// SendInvitation 向特定极客发出组队邀请
func (ctl *GeekController) SendInvitation(c *gin.Context) {
	// 获取邀请的极客ID从URL参数
	geekIDStr := c.Param("id")
	geekID, err := uuid.Parse(geekIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid geek ID"})
		return
	}

	// 获取当前登录用户的ID从上下文
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// 创建发送邀请的服务调用
	err = ctl.geekService.SendInvitation(geekID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent successfully"})
}
